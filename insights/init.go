package insights

import (
	"compress/gzip"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/daemon/logger"
	"github.com/sirupsen/logrus"
)

// ValidateLogOpt looks for all supported by splunk driver options
func ValidateLogOpt(cfg map[string]string) error {
	for key := range cfg {
		switch key {
		case insightsURLKey:
		case insightsTokenKey:
		case insightsInsecureSkipVerifyKey:
		case insightsGzipCompressionKey:
		case insightsVerifyConnectionKey:
		default:
			return fmt.Errorf("unknown log opt '%s' for %s log driver", key, insightsDriverName)
		}
	}
	return nil
}

func getAdvancedOptionDuration(envName string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(envName)
	if valueStr == "" {
		return defaultValue
	}
	parsedValue, err := time.ParseDuration(valueStr)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to parse value of %s as duration. Using default %v. %v", envName, defaultValue, err))
		return defaultValue
	}
	return parsedValue
}

func getAdvancedOptionInt(envName string, defaultValue int) int {
	valueStr := os.Getenv(envName)
	if valueStr == "" {
		return defaultValue
	}
	parsedValue, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to parse value of %s as integer. Using default %d. %v", envName, defaultValue, err))
		return defaultValue
	}
	return int(parsedValue)
}

// allow users to trust them with skipping verification
func getInsecureSkipVerify(info logger.Info) (bool, error) {
	if insecureSkipVerifyStr, ok := info.Config[insightsInsecureSkipVerifyKey]; ok {
		insecureSkipVerify, err := strconv.ParseBool(insecureSkipVerifyStr)
		if err != nil {
			logrus.Error(fmt.Sprintf("Failed to parse value of insecureSkipVerify as boolean. %v", err))
			return false, err
		}
		return insecureSkipVerify, nil
	}
	return false, nil
}

func getGzipCompression(info logger.Info) (bool, error) {
	if gzipCompressionStr, ok := info.Config[insightsGzipCompressionKey]; ok {
		gzipCompression, err := strconv.ParseBool(gzipCompressionStr)
		if err != nil {
			return false, err
		}
		return gzipCompression, nil
	}
	return false, nil
}

func getGzipCompressionLevel(info logger.Info) (int, error) {
	gzipCompressionLevel := gzip.DefaultCompression
	if gzipCompressionLevelStr, ok := info.Config[insightsGzipCompressionLevelKey]; ok {
		gzipCompressionLevel64, err := strconv.ParseInt(gzipCompressionLevelStr, 10, 32)
		if err != nil {
			logrus.Error(fmt.Sprintf("Failed to parse value of gzipCompressionLevel as integer. %v", err))
			return gzip.DefaultCompression, err
		}
		gzipCompressionLevel = int(gzipCompressionLevel64)
		if gzipCompressionLevel < gzip.DefaultCompression || gzipCompressionLevel > gzip.BestCompression {
			err := fmt.Errorf("not supported level '%s' for %s (supported values between %d and %d)", gzipCompressionLevelStr, insightsGzipCompressionLevelKey, gzip.DefaultCompression, gzip.BestCompression)
			return gzip.DefaultCompression, err
		}
	}
	return gzipCompressionLevel, nil
}
