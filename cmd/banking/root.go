// Copyright © 2015 Michael Wagner <mitch.wagna@gmail.com>
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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitch000001/go-hbci/client"
	"github.com/mitch000001/go-hbci/domain"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "banking",
	Short: "banking is a CLI to access bank data via HBCI",
	Long:  `banking is a CLI to access bank data via HBCI`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var cfgFile string
var url string
var UserID string
var BLZ string
var PIN string

var account domain.InternationalAccountConnection
var clientConfig client.Config
var hbciClient *client.Client
var debug bool

func init() {
	cobra.OnInitialize(
		initConfig,
		initClient,
	)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.banking.yaml)")
	rootCmd.PersistentFlags().StringVar(&url, "hbci.url", "", "the URL to the bank institute")
	rootCmd.PersistentFlags().StringVar(&UserID, "userID", "", "the account ID to authenticate with")
	rootCmd.PersistentFlags().StringVar(&BLZ, "blz", "", "the identifier for the bank institute")
	rootCmd.PersistentFlags().StringVar(&PIN, "pin", "", "the pin for the provided account")
	viper.BindPFlag("userID", rootCmd.PersistentFlags().Lookup("userID"))
	viper.BindPFlag("blz", rootCmd.PersistentFlags().Lookup("blz"))
	rootCmd.MarkPersistentFlagRequired("pin")

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug logging (very verbose)")
}

func initClient() {
	var missingFlags []string
	userID := viper.GetString("userID")
	blz := viper.GetString("blz")
	if userID == "" {
		missingFlags = append(missingFlags, `"userID"`)
	}
	if blz == "" {
		missingFlags = append(missingFlags, `"blz"`)
	}
	if len(missingFlags) != 0 {
		fmt.Printf("Error: required flag(s) %s not set\n", strings.Join(missingFlags, ", "))
		os.Exit(1)
	}
	clientConfig = client.Config{
		URL:                url,
		AccountID:          userID,
		BankID:             blz,
		PIN:                PIN,
		EnableDebugLogging: debug,
		ClientSystemID:     viper.GetString("client_system_id"),
		SecurityFunction:   viper.GetString("security_function"),
	}
	c, err := client.New(clientConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hbciClient = c
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

		// Search config in home directory with name ".banking" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".banking")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.WriteConfigAs(filepath.Join(home, ".banking.yaml"))
	}
}
