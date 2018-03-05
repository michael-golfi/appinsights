package handler

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/plugins/logdriver"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/jsonfilelog"
	protoio "github.com/gogo/protobuf/io"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tonistiigi/fifo"
	"gitlab.com/michael.golfi/appinsights/insights"
)

// Driver maintains a mutex for synchronizing map access for tracking logpairs and maintains the logging interface for each container
type Driver struct {
	logs   *logPairMap
	idx    *logPairMap
	logger logger.Logger
}

type logPair struct {
	fileLog      logger.Logger
	stream io.ReadCloser
	aiLog     logger.Logger
	info   logger.Info
}

// NewDriver creates a driver which initializes the logpairs for each container
func NewDriver() *Driver {
	return &Driver{
		logs: newLogPairMap(),
		idx: newLogPairMap(),
	}
}

func (d *Driver) StartLogging(file string, logCtx logger.Info) error {
	if _, exists := d.logs.Load(file); exists {
		return fmt.Errorf("logger for %q already exists", file)
	}

	logCtx.Config["tag"] = ""

	if logCtx.LogPath == "" {
		logCtx.LogPath = filepath.Join("/var/log/docker", logCtx.ContainerID)
	}
	if err := os.MkdirAll(filepath.Dir(logCtx.LogPath), 0755); err != nil {
		return errors.Wrap(err, "error setting up logger dir")
	}

	l, err := jsonfilelog.New(logCtx)
	if err != nil {
		return errors.Wrap(err, "error creating jsonfile logger")
	}

	sl, err := insights.New(logCtx)
	if err != nil {
		return errors.Wrap(err, "error creating appinsights logger")
	}

	logrus.WithField("id", logCtx.ContainerID).WithField("file", file).WithField("logpath", logCtx.LogPath).Debugf("Start logging")
	f, err := fifo.OpenFifo(context.Background(), file, syscall.O_RDONLY, 0700)
	if err != nil {
		return errors.Wrapf(err, "error opening logger fifo: %q", file)
	}

	lf := &logPair{l, f, sl, logCtx}
	d.logs.Store(file, lf)
	d.idx.Store(logCtx.ContainerID, lf)
	go consumeLog(lf)
	return nil
}

func (d *Driver) StopLogging(file string) error {
	logrus.WithField("file", file).Debugf("Stop logging")
	if lf, ok := d.logs.Load(file); ok {
		if err := lf.stream.Close(); err != nil {
			logrus.WithField("file", file).Errorf("Could not stop logging: %s", file)
			return err
		}
		d.logs.Delete(file)
	}
	return nil
}

func consumeLog(lf *logPair) {
	dec := protoio.NewUint32DelimitedReader(lf.stream, binary.BigEndian, 1e6)
	defer dec.Close()
	var buf logdriver.LogEntry
	for {
		if err := dec.ReadMsg(&buf); err != nil {
			if err == io.EOF {
				logrus.WithField("id", lf.info.ContainerID).WithError(err).Debug("shutting down log logger")
				lf.stream.Close()
				return
			}
			dec = protoio.NewUint32DelimitedReader(lf.stream, binary.BigEndian, 1e6)
		}

		var msg logger.Message
		msg.Line = buf.Line
		msg.Source = buf.Source
		msg.Partial = buf.Partial
		msg.Timestamp = time.Unix(0, buf.TimeNano)

		if err := lf.fileLog.Log(&msg); err != nil {
			logrus.WithField("id", lf.info.ContainerID).WithError(err).WithField("message", msg).Error("error writing log message")
			continue
		}

		if err := lf.aiLog.Log(&msg); err != nil {
			logrus.WithField("id", lf.info.ContainerID).WithError(err).WithField("message", msg).Error("error writing log message")
			continue
		}

		buf.Reset()
	}
}

func (d *Driver) ReadLogs(info logger.Info, config logger.ReadConfig) (io.ReadCloser, error) {
	lf, exists := d.idx.Load(info.ContainerID)
	if !exists {
		return nil, fmt.Errorf("logger does not exist for %s", info.ContainerID)
	}

	r, w := io.Pipe()
	lr, ok := lf.fileLog.(logger.LogReader)
	if !ok {
		return nil, fmt.Errorf("logger does not support reading")
	}

	go func() {
		watcher := lr.ReadLogs(config)

		enc := protoio.NewUint32DelimitedWriter(w, binary.BigEndian)
		defer enc.Close()
		defer watcher.Close()

		var buf logdriver.LogEntry
		for {
			select {
			case msg, ok := <-watcher.Msg:
				if !ok {
					w.Close()
					return
				}

				buf.Line = msg.Line
				buf.Partial = msg.Partial
				buf.TimeNano = msg.Timestamp.UnixNano()
				buf.Source = msg.Source

				if err := enc.WriteMsg(&buf); err != nil {
					w.CloseWithError(err)
					return
				}
			case err := <-watcher.Err:
				w.CloseWithError(err)
				return
			}

			buf.Reset()
		}
	}()

	return r, nil
}