package constants

import "time"

var (
	DriverName = "appinsights"

	// Application Insights Configuration Keys
	EndpointKey             = "endpoint"
	TokenKey                = "token"
	InsecureSkipVerifyKey   = "insecure-skip-verify"
	GzipCompressionKey      = "gzip"
	GzipCompressionLevelKey = "gzip-level"
	VerifyConnectionKey     = "verify-connection"
	BatchSizeKey            = "batch-size"
	BatchIntervalKey        = "batch-interval"

	// Application Insights String Configuration
	Endpoint                = "https://dc.services.visualstudio.com/v2/track"
	Token                   = ""
	VerifyConnectionStr     = "true"
	InsecureSkipVerifyStr   = "false"
	GzipCompressionStr      = "false"
	GzipCompressionLevelStr = "0"
	BatchSizeStr            = "1024"
	BatchIntervalStr        = "5s"

	// Application Insights Configuration
	VerifyConnection     = true
	InsecureSkipVerify   = false
	GzipCompression      = false
	GzipCompressionLevel = 0
	BatchSize            = 1024
	BatchInterval        = 5 * time.Second

	BufferMaximum = 10 * BatchSize
	StreamChannelSize = 4 * BatchSize
	SendTimeout = 30 * time.Second
)
