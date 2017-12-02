package insights

import (
	"fmt"
	"time"

	"github.com/docker/docker/daemon/logger"
)

type telemetry struct {
	Time         string
	Properties   map[string]string
	Context      map[string]string
	TagOverrides map[string]string
}

type metricTelemetry struct {
	telemetry
	Name   string
	Value  string
	Count  int
	Min    int
	Max    int
	StdDev int
}

type traceTelemetry struct {
	telemetry

	// Keep for now...
	Hostname string

	Message  string
	Severity string
}

func (d *metricTelemetry) String() string {
	return ""
}

func (d *traceTelemetry) String() string {
	return ""
}

func (l *insightsLogger) createInsightsMessage(msg *logger.Message) *traceTelemetry {
	message := *l.nullMessage
	message.Time = fmt.Sprintf("%f", float64(msg.Timestamp.UnixNano())/float64(time.Second))
	return &message
}
