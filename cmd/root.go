// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	level                      string
	cfgFile                    string
	insightsURL                string
	insightsToken              string
	insightsInsecureSkipVerify string
	insightsGzipCompression    string
	insightsVerifyConnection   string
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.appinsights.yaml)")

	rootCmd.PersistentFlags().StringVarP(&insightsURL, "url", "", "https://dc.services.visualstudio.com", "The URL for App Insights")
	rootCmd.PersistentFlags().StringVarP(&insightsToken, "key", "k", "", "Insights Instrumentation Key")
	rootCmd.PersistentFlags().StringVarP(&insightsInsecureSkipVerify, "insecure-skip-verify", "", "false", "Skip verifying the SSL certificate")
	rootCmd.PersistentFlags().StringVarP(&insightsGzipCompression, "compress", "c", "false", "Enable GZip compression")
	rootCmd.PersistentFlags().StringVarP(&insightsVerifyConnection, "verify-connection", "", "false", "Verify the connection to App Insights on start")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".newcmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".appinsights")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
