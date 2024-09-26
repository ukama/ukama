/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "msgcli",
	Short: "Ukama's MsgClient CLI test tool",
	Long: `msgcli - short for message command line interface - scans
and test events for Ukama services.

msgcli is a CLI tool that scans any given Ukama service source
code and displays all the events (routing keys + messages)
that it both listens to and generates, and also sends events to
a running MsgClient instance targeting a specific Ukama service.`,

	Version: "0.1",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(InitConfig)
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.msgcli.yaml)")

	versionTemplate := `{{printf "%s: %s - version: %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}
