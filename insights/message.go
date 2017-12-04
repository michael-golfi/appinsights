package insights

import (
	"time"

	"github.com/docker/docker/daemon/logger"
	"github.com/fatih/structs"
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
	Version       int                    `json:"ver"`
	Properties    map[string]interface{} `json:"properties"`
	Message       string                 `json:"message"`
	SeverityLevel int                    `json:"severityLevel"`
}

func (l *insightsLogger) createInsightsMessage(msg *logger.Message) *envelope {
	message := *l.nullMessage
	message.Time = time.Now().UTC().Format(time.RFC3339)

	ctx := structs.Map(l.logCtx)

	message.Data = data{
		BaseType: "MessageData",
		BaseData: messageData{
			Version:       2,
			Message:       string(msg.Line),
			SeverityLevel: verbose,
			Properties:    ctx,
		},
	}
	return &message
}
