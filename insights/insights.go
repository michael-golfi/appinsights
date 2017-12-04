package insights

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/docker/docker/daemon/logger"
	"github.com/sirupsen/logrus"
)

const (
	insightsDriverName            = "insights"
	insightsURLKey                = "insights-url"
	insightsTokenKey              = "insights-key"
	insightsInsecureSkipVerifyKey = "insights-insecureskipverify"
	insightsGzipCompressionKey    = "insights-gzip"
	insightsVerifyConnectionKey   = "insights-verify-connection"
)

const (
	// How often do we send messages (if we are not reaching batch size)
	defaultPostMessagesFrequency = 5 * time.Second
	// How big can be batch of messages
	defaultPostMessagesBatchSize = 1000
	// Maximum number of messages we can store in buffer
	defaultBufferMaximum = 10 * defaultPostMessagesBatchSize
	// Number of messages allowed to be queued in the channel
	defaultStreamChannelSize = 4 * defaultPostMessagesBatchSize

	batchSendTimeout = 30 * time.Second
)

const (
	envVarPostMessagesFrequency = "INSIGHTS_LOGGING_DRIVER_POST_MESSAGES_FREQUENCY"
	envVarPostMessagesBatchSize = "INSIGHTS_LOGGING_DRIVER_POST_MESSAGES_BATCH_SIZE"
	envVarBufferMaximum         = "INSIGHTS_LOGGING_DRIVER_BUFFER_MAX"
	envVarStreamChannelSize     = "INSIGHTS_LOGGING_DRIVER_CHANNEL_SIZE"
)

type insightsLogger struct {
	client    *http.Client
	transport *http.Transport

	url                string
	instrumentationKey string
	nullMessage        *envelope
	// TODO - Support more message types than only trace
	//nullTrace   *insightsMessage

	// http compression
	gzipCompression      bool
	gzipCompressionLevel int

	// Advanced options
	postMessagesFrequency time.Duration
	postMessagesBatchSize int
	bufferMaximum         int

	// For synchronization between background worker and logger.
	// We use channel to send messages to worker go routine.
	// All other variables for blocking Close call before we flush all messages to HEC
	stream     chan *envelope
	lock       sync.RWMutex
	closed     bool
	closedCond *sync.Cond

	logCtx logger.Info
}

func init() {
	if err := logger.RegisterLogDriver(insightsDriverName, New); err != nil {
		logrus.Fatal(err)
	}
	if err := logger.RegisterLogOptValidator(insightsDriverName, ValidateLogOpt); err != nil {
		logrus.Fatal(err)
	}
}

// New creates splunk logger driver using configuration passed in context
func New(info logger.Info) (logger.Logger, error) {
	/*hostname, err := info.Hostname()
	if err != nil {
		return nil, fmt.Errorf("%s: cannot access hostname to set source field", insightsDriverName)
	}*/

	// Parse and validate URL
	insightsURL, err := parseURL(info)
	if err != nil {
		return nil, err
	}

	// Instrumentation Token is required parameter
	insightsToken, ok := info.Config[insightsTokenKey]
	if !ok {
		return nil, fmt.Errorf("%s: %s is expected", insightsDriverName, insightsTokenKey)
	}

	insecureSkipVerify, err := getInsecureSkipVerify(info)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}

	// Set GZip
	gzipCompression, err := getGzipCompression(info)
	if err != nil {
		return nil, err
	}

	// Set up Transport and client
	gzipCompressionLevel, err := getGzipCompressionLevel(info)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	nullMessage := &envelope{
		Name: fmt.Sprintf("Microsoft.ApplicationInsights.MessageData", insightsToken),
		Ikey: insightsToken,
	}

	var (
		postMessagesFrequency = getAdvancedOptionDuration(envVarPostMessagesFrequency, defaultPostMessagesFrequency)
		postMessagesBatchSize = getAdvancedOptionInt(envVarPostMessagesBatchSize, defaultPostMessagesBatchSize)
		bufferMaximum         = getAdvancedOptionInt(envVarBufferMaximum, defaultBufferMaximum)
		streamChannelSize     = getAdvancedOptionInt(envVarStreamChannelSize, defaultStreamChannelSize)
	)

	logger := &insightsLogger{
		client:             client,
		url:                insightsURL.String(),
		instrumentationKey: insightsToken,

		nullMessage:           nullMessage,
		gzipCompression:       gzipCompression,
		gzipCompressionLevel:  gzipCompressionLevel,
		stream:                make(chan *envelope, streamChannelSize),
		postMessagesFrequency: postMessagesFrequency,
		postMessagesBatchSize: postMessagesBatchSize,
		bufferMaximum:         bufferMaximum,

		logCtx: info,
	}

	// By default we verify connection, but we allow use to skip that
	verifyConnection := true
	if verifyConnectionStr, ok := info.Config[insightsVerifyConnectionKey]; ok {
		var err error
		verifyConnection, err = strconv.ParseBool(verifyConnectionStr)
		if err != nil {
			return nil, err
		}
	}
	if verifyConnection {
		err = verifyInsightsConnection(logger)
		if err != nil {
			return nil, err
		}
	}

	go logger.worker()
	return logger, nil
}

func (l *insightsLogger) Name() string {
	return insightsDriverName
}

func (l *insightsLogger) Log(msg *logger.Message) error {
	message := l.createInsightsMessage(msg)
	logger.PutMessage(msg)
	return l.queueMessageAsync(message)
}
