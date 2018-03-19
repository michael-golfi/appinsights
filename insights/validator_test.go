package insights

import (
	"testing"
	"gitlab.com/michael.golfi/appinsights/constants"
	"github.com/stretchr/testify/require"
	"github.com/docker/docker/daemon/logger"
	"time"
)

func copyConfig(src logger.Info) logger.Info {
	dest := make(map[string]string, len(src.Config))
	for k,v := range src.Config {
		dest[k] = v
	}
	return logger.Info { Config: dest }
}

func TestInitializeEnv(t *testing.T) {
	allValuesEmpty := logger.Info{
		Config: map[string]string{
			constants.EndpointKey:             "",
			constants.TokenKey:                "",
			constants.InsecureSkipVerifyKey:   "",
			constants.GzipCompressionKey:      "",
			constants.GzipCompressionLevelKey: "",
			constants.VerifyConnectionKey:     "",
			constants.BatchSizeKey:            "",
			constants.BatchIntervalKey:        "",
		},
	}
	err := InitializeEnv(allValuesEmpty)
	require.NoError(t, err)


	emptyConfig := logger.Info { Config: map[string]string{} }
	err = InitializeEnv(emptyConfig)
	require.Error(t, err)

	badKey := copyConfig(allValuesEmpty)
	badKey.Config["Some Weird Key"] = ""
	err = InitializeEnv(badKey)
	require.Error(t, err)

	noToken := copyConfig(allValuesEmpty)
	delete(noToken.Config, constants.TokenKey)
	err = InitializeEnv(noToken)
	require.Error(t, err)
}

func TestValidateLogOpt(t *testing.T) {
	allSuccess := make(map[string]string)
	allSuccess[constants.EndpointKey] = ""
	allSuccess[constants.TokenKey] = ""
	allSuccess[constants.InsecureSkipVerifyKey] = ""
	allSuccess[constants.GzipCompressionKey] = ""
	allSuccess[constants.GzipCompressionLevelKey] = ""
	allSuccess[constants.VerifyConnectionKey] = ""
	allSuccess[constants.BatchSizeKey] = ""
	allSuccess[constants.BatchIntervalKey] = ""
	err := validateLogOpt(allSuccess)
	require.NoError(t, err)

	oneParam := make(map[string]string)
	oneParam[constants.EndpointKey] = ""
	err = validateLogOpt(oneParam)
	require.NoError(t, err)

	invalidParam := make(map[string]string)
	invalidParam["some param"] = ""
	err = validateLogOpt(invalidParam)
	require.Error(t, err)

	empty := make(map[string]string)
	err = validateLogOpt(empty)
	require.Error(t, err)
}

func TestGetAdvancedOption(t *testing.T) {
	key := "key"
	val := "val"
	def := "default"

	info := logger.Info{
		Config: map[string]string{
			key: val,
		},
	}

	res := getAdvancedOption(info, key, def)
	require.Equal(t, val, res)

	delete(info.Config, "key")
	res = getAdvancedOption(info, key, def)
	require.Equal(t, def, res)
}

func TestGetAdvancedOptionDuration(t *testing.T) {
	key := "key"
	valStr := "5s"
	val := 5 * time.Second
	def := 6 * time.Second

	info := logger.Info{
		Config: map[string]string{
			key: valStr,
		},
	}

	res := getAdvancedOptionDuration(info, key, def)
	require.Equal(t, val, res)

	delete(info.Config, "key")
	res = getAdvancedOptionDuration(info, key, def)
	require.Equal(t, def, res)

	info.Config[key] = "bad val"
	res = getAdvancedOptionDuration(info, key, def)
	require.Equal(t, def, res)
}

func TestGetAdvancedOptionInt(t *testing.T) {
	key := "key"
	valStr := "5"

	info := logger.Info{
		Config: map[string]string{
			key: valStr,
		},
	}

	res := getAdvancedOptionInt(info, key, 6)
	require.Equal(t, 5, res)

	delete(info.Config, "key")
	res = getAdvancedOptionInt(info, key, 6)
	require.Equal(t, 6, res)

	info.Config[key] = "bad val"
	res = getAdvancedOptionInt(info, key, 6)
	require.Equal(t, 6, res)
}

func TestGetAdvancedOptionBool(t *testing.T) {
	key := "key"
	valStr := "true"

	info := logger.Info{
		Config: map[string]string{
			key: valStr,
		},
	}

	res := getAdvancedOptionBool(info, key, false)
	require.Equal(t, true, res)

	delete(info.Config, "key")
	res = getAdvancedOptionBool(info, key, false)
	require.Equal(t, false, res)

	info.Config[key] = "bad val"
	res = getAdvancedOptionBool(info, key, false)
	require.Equal(t, false, res)
}
