package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Administrative configuration
	level   string
	cfgFile string

	// Application Insights Configuration
	insightsURL                  string
	insightsToken                string
	insightsInsecureSkipVerify   string
	insightsGzipCompression      string
	insightsGzipCompressionLevel string
	insightsVerifyConnection     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "appinsights",
	Short: "appinsights is a docker logging plugin",
	Long: `appinsights is a docker logging plugin that streams logs to json files 
in local disk storage and also streams to Microsoft App Insights.
The plugin supports the docker log command and also supports buffering and retry
for remote App Insights in the case of network disconnection`,
}

// Execute is the entrypoint to the command execution context.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&insightsURL, "url", "", "https://dc.services.visualstudio.com", "The URL for App Insights")
	rootCmd.PersistentFlags().StringVarP(&insightsToken, "key", "k", "", "Insights Instrumentation Key")
	rootCmd.PersistentFlags().StringVarP(&insightsInsecureSkipVerify, "insecure-skip-verify", "", "false", "Skip verifying the SSL certificate")
	rootCmd.PersistentFlags().StringVarP(&insightsGzipCompression, "compress", "c", "false", "Enable GZip compression")
	rootCmd.PersistentFlags().StringVarP(&insightsGzipCompressionLevel, "compress-level", "", "0", "GZip compression level")
	rootCmd.PersistentFlags().StringVarP(&insightsVerifyConnection, "verify-connection", "", "false", "Verify the connection to App Insights on start")
}
