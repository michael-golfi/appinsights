package insights

import (
	"time"

	"encoding/json"
	"log"

	ai "github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/docker/docker/daemon/logger"
)

func (l *insightsLogger) createInsightsMessage(msg *logger.Message) *ai.Envelope {
	ctx, err := mapLogCtx(l.logCtx)
	if err != nil {
		log.Println(err)
	}

	ctx["Source"] = msg.Source
	for _, attr := range msg.Attrs {
		ctx[attr.Key] = attr.Value
	}

	return &ai.Envelope{
		Name:       "Microsoft.ApplicationInsights.MessageData",
		IKey:       l.instrumentationKey,
		SampleRate: 100.0,
		Time:       time.Now().UTC().Format(time.RFC3339),
		Data: &ai.Data{
			Base: ai.Base{
				BaseType: "MessageData",
			},
			BaseData: &ai.MessageData{
				Ver:           2,
				Message:       string(msg.Line),
				SeverityLevel: ai.Verbose,
				Properties:    ctx,
			},
		},
	}
}

func mapLogCtx(logCtx logger.Info) (map[string]string, error) {
	out := make(map[string]string, 5)
	out["ContainerID"] = logCtx.ContainerID
	out["ContainerName"] = logCtx.ContainerName
	out["ContainerEntrypoint"] = logCtx.ContainerEntrypoint
	out["ContainerImageID"] = logCtx.ContainerImageID
	out["ContainerImageName"] = logCtx.ContainerImageName
	out["LogPath"] = logCtx.LogPath
	out["DaemonName"] = logCtx.DaemonName
	out["ContainerCreated"] = logCtx.ContainerCreated.Format(time.RFC3339)

	args, err := json.Marshal(logCtx.ContainerArgs)
	if err != nil {
		return nil, err
	}

	env, err := json.Marshal(logCtx.ContainerEnv)
	if err != nil {
		return nil, err
	}

	conf, err := json.Marshal(logCtx.Config)
	if err != nil {
		return nil, err
	}

	labels, err := json.Marshal(logCtx.ContainerLabels)
	if err != nil {
		return nil, err
	}

	out["ContainerArgs"] = string(args)
	out["ContainerEnv"] = string(env)
	out["Config"] = string(conf)
	out["ContainerLabels"] = string(labels)
	return out, nil
}
