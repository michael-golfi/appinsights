package insights

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/pkg/urlutil"
	"github.com/sirupsen/logrus"
)

func parseURL(info logger.Info) (*url.URL, error) {
	insightsURLStr, ok := info.Config[insightsURLKey]
	if !ok {
		return nil, fmt.Errorf("%s: %s is expected", insightsDriverName, insightsURLKey)
	}

	insightsURL, err := url.Parse(insightsURLStr)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse %s as url value in %s", insightsDriverName, insightsURLStr, insightsURLKey)
	}

	if !urlutil.IsURL(insightsURLStr) ||
		!insightsURL.IsAbs() ||
		(insightsURL.Path != "" && insightsURL.Path != "/") ||
		insightsURL.RawQuery != "" ||
		insightsURL.Fragment != "" {
		return nil, fmt.Errorf("%s: expected format scheme://dns_name_or_ip:port for %s", insightsDriverName, insightsURLKey)
	}

	// REVIEW
	//	insightsURL.Path = "/services/collector/event/1.0"

	return insightsURL, nil
}

func verifyInsightsConnection(l *insightsLogger) error {
	req, err := http.NewRequest(http.MethodOptions, l.url, nil)
	if err != nil {
		return err
	}
	res, err := l.client.Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if res.StatusCode != http.StatusOK {
		var body []byte
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: failed to verify connection - %s - %s", insightsDriverName, res.Status, body)
	}
	return nil
}

// REVIEW
func (l *insightsLogger) queueMessageAsync(message *traceTelemetry) error {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if l.closedCond != nil {
		return fmt.Errorf("%s: driver is closed", insightsDriverName)
	}
	l.stream <- message
	return nil
}

// REVIEW
func (l *insightsLogger) postMessages(messages []*traceTelemetry, lastChance bool) []*traceTelemetry {
	messagesLen := len(messages)

	ctx, cancel := context.WithTimeout(context.Background(), batchSendTimeout)
	defer cancel()

	for i := 0; i < messagesLen; i += l.postMessagesBatchSize {
		upperBound := i + l.postMessagesBatchSize
		if upperBound > messagesLen {
			upperBound = messagesLen
		}

		if err := l.tryPostMessages(ctx, messages[i:upperBound]); err != nil {
			logrus.WithError(err).WithField("module", "logger/splunk").Warn("Error while sending logs")
			if messagesLen-i >= l.bufferMaximum || lastChance {
				// If this is last chance - print them all to the daemon log
				if lastChance {
					upperBound = messagesLen
				}
				// Not all sent, but buffer has got to its maximum, let's log all messages
				// we could not send and return buffer minus one batch size
				for j := i; j < upperBound; j++ {
					if jsonEvent, err := json.Marshal(messages[j]); err != nil {
						logrus.Error(err)
					} else {
						logrus.Error(fmt.Errorf("Failed to send a message '%s'", string(jsonEvent)))
					}
				}
				return messages[upperBound:messagesLen]
			}
			// Not all sent, returning buffer from where we have not sent messages
			return messages[i:messagesLen]
		}
	}
	// All sent, return empty buffer
	return messages[:0]
}

// REVIEW
func (l *insightsLogger) tryPostMessages(ctx context.Context, messages []*traceTelemetry) error {
	if len(messages) == 0 {
		return nil
	}
	var buffer bytes.Buffer
	var writer io.Writer
	var gzipWriter *gzip.Writer
	var err error
	// If gzip compression is enabled - create gzip writer with specified compression
	// level. If gzip compression is disabled, use standard buffer as a writer
	if l.gzipCompression {
		gzipWriter, err = gzip.NewWriterLevel(&buffer, l.gzipCompressionLevel)
		if err != nil {
			return err
		}
		writer = gzipWriter
	} else {
		writer = &buffer
	}
	for _, message := range messages {
		jsonEvent, err := json.Marshal(message)
		if err != nil {
			return err
		}
		if _, err := writer.Write(jsonEvent); err != nil {
			return err
		}
	}
	// If gzip compression is enabled, tell it, that we are done
	if l.gzipCompression {
		err = gzipWriter.Close()
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest("POST", l.url, bytes.NewBuffer(buffer.Bytes()))
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", l.auth)
	// Tell if we are sending gzip compressed body
	if l.gzipCompression {
		req.Header.Set("Content-Encoding", "gzip")
	}
	res, err := l.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		var body []byte
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: failed to send event - %s - %s", insightsDriverName, res.Status, body)
	}
	io.Copy(ioutil.Discard, res.Body)
	return nil
}
