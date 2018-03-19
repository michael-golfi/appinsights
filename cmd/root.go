package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/michael.golfi/appinsights/constants"
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
	rootCmd.PersistentFlags().StringVarP(&constants.Endpoint, constants.EndpointKey, "", constants.Endpoint, "The URL for App Insights")
	rootCmd.PersistentFlags().StringVarP(&constants.Token, constants.TokenKey, "k", constants.Token, "Insights Instrumentation Key")
	rootCmd.PersistentFlags().StringVarP(&constants.InsecureSkipVerifyStr, constants.InsecureSkipVerifyKey, "", constants.InsecureSkipVerifyStr, "Skip verifying the SSL certificate")
	rootCmd.PersistentFlags().StringVarP(&constants.GzipCompressionStr, constants.GzipCompressionKey, "c", constants.GzipCompressionStr, "Enable GZip compression")
	rootCmd.PersistentFlags().StringVarP(&constants.GzipCompressionLevelStr, constants.GzipCompressionLevelKey, "", constants.GzipCompressionLevelStr, "GZip compression level")
	rootCmd.PersistentFlags().StringVarP(&constants.VerifyConnectionStr, constants.VerifyConnectionKey, "", constants.VerifyConnectionStr, "Verify the connection to App Insights on start")
	rootCmd.PersistentFlags().StringVarP(&constants.BatchSizeStr, constants.BatchSizeKey, "", constants.BatchSizeStr, "Message Batch Size")
	rootCmd.PersistentFlags().StringVarP(&constants.BatchIntervalStr, constants.BatchIntervalKey, "", constants.BatchIntervalStr, "Message Batch Interval")
}
