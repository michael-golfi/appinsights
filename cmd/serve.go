// Copyright © 2017 Michael Golfi <michael.golfi@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/michael.golfi/appinsights/handler"
)

var level string
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
