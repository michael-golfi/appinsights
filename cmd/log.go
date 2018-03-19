package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/daemon/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/michael.golfi/appinsights/insights"
	"gitlab.com/michael.golfi/appinsights/constants"
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
			logrus.Error(fmt.Sprintf("Failed to create logging client. %v", err))
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
	config := make(map[string]string, 8)
	config[constants.EndpointKey] = constants.Endpoint
	config[constants.TokenKey] = constants.Token
	config[constants.InsecureSkipVerifyKey] = constants.InsecureSkipVerifyStr
	config[constants.GzipCompressionKey] = constants.GzipCompressionStr
	config[constants.GzipCompressionLevelKey] = constants.GzipCompressionLevelStr
	config[constants.VerifyConnectionKey] = constants.VerifyConnectionStr
	config[constants.BatchSizeKey] = constants.BatchSizeStr
	config[constants.BatchIntervalKey] = constants.BatchIntervalStr

	return logger.Info{
		Config: config,
	}
}

func init() {
	rootCmd.AddCommand(logCmd)
}
