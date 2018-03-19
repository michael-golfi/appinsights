package insights

import (
	"fmt"
	"gitlab.com/michael.golfi/appinsights/constants"
	"time"
	"github.com/docker/docker/daemon/logger"
	"strconv"
	"github.com/sirupsen/logrus"
)

func InitializeEnv(info logger.Info) error {
	if err := validateLogOpt(info.Config); err != nil {
		return err
	}

	// Instrumentation Token is required parameter
	token, ok := info.Config[constants.TokenKey]
	if !ok {
		return fmt.Errorf("%s: %s is expected", constants.DriverName, constants.TokenKey)
	}

	// Merge configurations
	var (
		endpoint             = getAdvancedOption(info, constants.EndpointKey, constants.Endpoint)
		skipVerify           = getAdvancedOptionBool(info, constants.InsecureSkipVerifyKey, constants.InsecureSkipVerify)
		gzipCompression      = getAdvancedOptionBool(info, constants.GzipCompressionKey, constants.GzipCompression)
		gzipCompressionLevel = getAdvancedOptionInt(info, constants.GzipCompressionLevelKey, constants.GzipCompressionLevel)
		verifyConnection     = getAdvancedOptionBool(info, constants.VerifyConnectionKey, constants.VerifyConnection)
		batchSize            = getAdvancedOptionInt(info, constants.BatchSizeKey, constants.BatchSize)
		batchInterval        = getAdvancedOptionDuration(info, constants.BatchIntervalKey, constants.BatchInterval)
	)

	constants.Endpoint = endpoint
	constants.Token = token
	constants.InsecureSkipVerify = skipVerify
	constants.GzipCompression = gzipCompression
	constants.GzipCompressionLevel = gzipCompressionLevel
	constants.VerifyConnection = verifyConnection
	constants.BatchSize = batchSize
	constants.BatchInterval = batchInterval
	return nil
}

func validateLogOpt(cfg map[string]string) error {
	if len(cfg) == 0 {
		return fmt.Errorf("configuration cannot be empty")
	}

	for key := range cfg {
		switch key {
		case constants.EndpointKey:
		case constants.TokenKey:
		case constants.InsecureSkipVerifyKey:
		case constants.GzipCompressionKey:
		case constants.GzipCompressionLevelKey:
		case constants.VerifyConnectionKey:
		case constants.BatchSizeKey:
		case constants.BatchIntervalKey:
		default:
			return fmt.Errorf("unknown log opt '%s' for %s log driver", key, constants.DriverName)
		}
	}
	return nil
}

func getAdvancedOption(info logger.Info, name, def string) string {
	val, ok := info.Config[name]
	if val == "" || !ok {
		return def
	}
	return val
}

func getAdvancedOptionDuration(info logger.Info, name string, def time.Duration) time.Duration {
	val, ok := info.Config[name]
	if val == "" || !ok {
		return def
	}
	parsed, err := time.ParseDuration(val)
	if err != nil {
		logrus.Errorf("failed to parse value of %s as duration. Using default %v. %v", name, def, err)
		return def
	}
	return parsed
}

func getAdvancedOptionInt(info logger.Info, name string, def int) int {
	val, ok := info.Config[name]
	if val == "" || !ok {
		return def
	}
	parsed, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		logrus.Errorf("failed to parse value of %s as duration. Using default %v. %v", name, def, err)
		return def
	}
	return int(parsed)
}

func getAdvancedOptionBool(info logger.Info, name string, def bool) bool {
	val, ok := info.Config[name]
	if val == "" || !ok {
		return def
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		logrus.Errorf("Failed to parse value of %s as duration. Using default %v. %v", name, def, err)
		return def
	}
	return parsed
}
