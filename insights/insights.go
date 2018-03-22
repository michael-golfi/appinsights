package insights

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/docker/docker/daemon/logger"
	"github.com/sirupsen/logrus"
	"gitlab.com/michael.golfi/appinsights/constants"
)

type insightsLogger struct {
	client                *http.Client
	transport             *http.Transport
	url                   string
	instrumentationKey    string
	gzipCompression       bool
	gzipCompressionLevel  int
	postMessagesFrequency time.Duration
	postMessagesBatchSize int
	bufferMaximum         int
	sendTimeout           time.Duration
	// For synchronization between background worker and logger.
	// We use channel to send messages to worker go routine.
	// All other variables for blocking Close call before we flush all messages to HEC
	stream     chan *contracts.Envelope
	lock       sync.RWMutex
	closed     bool
	closedCond *sync.Cond
	logCtx     logger.Info
}

func init() {
	if err := logger.RegisterLogDriver(constants.DriverName, New); err != nil {
		logrus.Fatal(err)
	}
	if err := logger.RegisterLogOptValidator(constants.DriverName, validateLogOpt); err != nil {
		logrus.Fatal(err)
	}
}

// New creates appinsights logger driver using configuration passed in context
func New(info logger.Info) (logger.Logger, error) {
	if err := InitializeEnv(info); err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: constants.InsecureSkipVerify,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
	}

	if constants.VerifyConnection {
		err := verifyInsightsConnection(constants.Endpoint)
		if err != nil {
			return nil, err
		}
	}

	insightsLogger := &insightsLogger{
		client:                client,
		transport:             transport,
		url:                   constants.Endpoint,
		instrumentationKey:    constants.Token,
		gzipCompression:       constants.GzipCompression,
		gzipCompressionLevel:  constants.GzipCompressionLevel,
		stream:                make(chan *contracts.Envelope, constants.StreamChannelSize),
		postMessagesFrequency: constants.BatchInterval,
		postMessagesBatchSize: constants.BatchSize,
		bufferMaximum:         constants.BufferMaximum,
		sendTimeout:           constants.SendTimeout,
		logCtx:                info,
	}

	go insightsLogger.worker()
	return insightsLogger, nil
}

func (l *insightsLogger) Name() string {
	return constants.DriverName
}

func (l *insightsLogger) Log(msg *logger.Message) error {
	message := l.createInsightsMessage(msg)
	logger.PutMessage(msg)
	return l.queueMessageAsync(message)
}
