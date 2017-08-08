package insightslib

import (
	"fmt"

	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/docker/docker/daemon/logger"
)

type insightsLogger struct {
	appinsights.TelemetryClient
}

const (
	insightsDriverName = "insights"
	insightsURLKey     = "insights-url"
	insightsTokenKey   = "insights-key"
)

func New(info logger.Info) (logger.Logger, error) {
	insightsToken, ok := info.Config[insightsTokenKey]
	if !ok {
		return nil, fmt.Errorf("%s: %s is expected", insightsDriverName, insightsTokenKey)
	}

	client := appinsights.NewTelemetryClient(insightsToken)
	return insightsLogger{
		TelemetryClient: client,
	}, nil
}

func (l insightsLogger) Name() string {
	return insightsDriverName
}

func (l insightsLogger) Log(msg *logger.Message) error {
	l.TrackTrace(string(msg.Line))
	return nil
}

func (l insightsLogger) Close() error {
	return nil
}
