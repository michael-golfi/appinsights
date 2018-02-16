package cmd

import (
	"fmt"
	"os"

	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/michael.golfi/appinsights/handler"
)

var logLevels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start serving the plugin",
	Long: `This plugin streams log messages to Microsoft App Insights and
	it allows the usage of docker log. This server starts serving on the
	Unix socket: /run/docker/plugins/appinsights.sock`,
	Run: func(cmd *cobra.Command, args []string) {
		if logLevel, exists := logLevels[level]; exists {
			logrus.SetLevel(logLevel)
		} else {
			fmt.Fprintln(os.Stderr, "Invalid log level: ", logLevel)
			os.Exit(1)
		}

		h := sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
		handler.Handle(&h, handler.NewDriver())
		if err := h.ServeUnix("appinsights", 0); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&level, "verbose", "v", "info", "Sets log level: [info, debug, warn, error]")
}
