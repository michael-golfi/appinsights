package insights

import (
	"time"

	"github.com/docker/docker/daemon/logger"
)

const (
	verbose = 0
	//Info = 1
	//Debug = 1
	//Error = 2
	critical = 4
)

type envelope struct {
	Name string            `json:"name"`
	Time string            `json:"time"`
	Ikey string            `json:"iKey"`
	Tags map[string]string `json:"tags"`
	Data data              `json:"data"`
}

type data struct {
	BaseType string      `json:"baseType"`
	BaseData messageData `json:"baseData"`
}

type messageData struct {
	Version       int               `json:"ver"`
	Properties    map[string]string `json:"properties"`
	Message       string            `json:"message"`
	SeverityLevel int               `json:"severityLevel"`
}

func (l *insightsLogger) createInsightsMessage(msg *logger.Message) *envelope {
	message := *l.nullMessage
	message.Time = msg.Timestamp.Format(time.RFC3339)

	props := make(map[string]string)
	for _, attr := range msg.Attrs {
		props[attr.Key] = attr.Value
	}

	props["source"] = msg.Source
	message.Data = data{
		BaseType: "MessageData",
		BaseData: messageData{
			Version:       2,
			Message:       string(msg.Line),
			SeverityLevel: verbose,
			Properties:    props,
		},
	}
	return &message
}
