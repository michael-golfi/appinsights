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
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/daemon/logger"
	"github.com/spf13/cobra"
	"gitlab.com/michael.golfi/appinsights/insights"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log [message]",
	Short: "Test sending a message to App Insights",
	Long: `Log a message to Microsoft App Insights. 
Use this command to test connectivity.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := insights.New(createLoggerInfo())
		if err != nil {
			panic(err)
		}

		in := bufio.NewReader(os.Stdin)

		if len(args) < 1 {
			for {
				fmt.Print("> ")
				line, _, err := in.ReadLine()
				if err != nil {
					panic(err)
				}

				input := string(line)

				if input == "exit" {
					return
				}

				if input == "" {
					continue
				}

				sendMessage(client, input)
			}
		} else if len(args) == 1 {
			sendMessage(client, args[0])
		} else {
			log.Println("log only accepts zero or one argument.")
		}
	},
}

func sendMessage(client logger.Logger, line string) {
	msg := logger.NewMessage()
	msg.Line = []byte(line)
	client.Log(msg)
}

func createLoggerInfo() logger.Info {
	config := make(map[string]string)
	config["insights-url"] = insightsURL
	config["insights-key"] = insightsToken
	config["insights-insecureskipverify"] = insightsInsecureSkipVerify
	config["insights-gzip"] = insightsGzipCompression
	config["insights-verify-connection"] = insightsVerifyConnection

	return logger.Info{
		Config: config,
	}
}

func init() {
	rootCmd.AddCommand(logCmd)
}
