/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ukama/ukama/utils/ukamacli/cmd/version"

	homedir "github.com/mitchellh/go-homedir"
)

var rootCmd = &cobra.Command{
	Use:   "ukamacli",
	Short: "Ukama's management CLI tool",
	Long: `ukamacli - short for ukama command line interface - manages
various resources of a given Ukama saas.

ukamacli is a CLI tool that can be used to build, deploy,
install,  start, restart, stop or uninstall various ukama
resources.`,

	Version: version.Version,
}

var cfgFile string = ""

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config", "", "config file (default is $HOME/.ukamacli.yaml)")

	replacer := strings.NewReplacer("-", "_")

	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("UKAMACLI")

	versionTemplate := `{{printf "%s: %s - version: %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".ukamacli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file", viper.ConfigFileUsed())
	}
}
