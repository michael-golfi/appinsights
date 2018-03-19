package insights

import (
	"testing"
	"github.com/docker/docker/daemon/logger"
	"github.com/stretchr/testify/require"
	"github.com/docker/docker/api/types/backend"
	"gitlab.com/michael.golfi/appinsights/constants"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights/contracts"
)

func TestNewAppInsightsLogger(t *testing.T) {
	goodInfo := logger.Info{
		Config: map[string]string{
			constants.TokenKey: "some token",
		},
	}

	client, err := New(goodInfo)
	require.NoError(t, err)
	require.Equal(t, constants.DriverName, client.Name())
	require.NoError(t, client.Close())

	noToken := logger.Info{
		Config: map[string]string{},
	}

	client, err = New(noToken)
	require.Nil(t, client)
	require.Error(t, err)
}

func TestInflateTraceMessage(t *testing.T) {
	goodInfo := logger.Info{
		Config: map[string]string{
			constants.TokenKey: "some token",
		},
	}

	insightsLog := insightsLogger{
		logCtx: goodInfo,
	}

	msg := logger.NewMessage()
	msg.Source = "Test"
	msg.Line = []byte("Some Message")
	msg.Attrs = []backend.LogAttr{
		{Key: "Hello", Value: "World"},
	}

	trace := insightsLog.createInsightsMessage(msg)
	data, ok := trace.Data.(contracts.Data)
	require.True(t, ok)

	val, ok := data.BaseData.(contracts.MessageData)
	require.True(t, ok)
	require.Equal(t, string(msg.Line), val.Message)
	require.Equal(t, msg.Source, val.Properties["Source"])
	require.Equal(t, "World", val.Properties["Hello"])
}
