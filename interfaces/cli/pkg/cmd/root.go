/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ukama/ukama/interfaces/cli/cmd/version"
	"github.com/ukama/ukama/interfaces/cli/pkg"
	"github.com/ukama/ukama/interfaces/cli/pkg/cmd/config"
	"github.com/ukama/ukama/interfaces/cli/pkg/cmd/node"
)

func RootCommand() *cobra.Command {
	var cfgFile string

	var rootCmd = &cobra.Command{
		Use:   pkg.CliName,
		Short: "Ukama CLI",
	}

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file (default is $HOME/.ukama.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
	rootCmd.Version = version.Version
	cfMgr := config.NewConfMgr(cfgFile, rootCmd.OutOrStdout(), rootCmd.ErrOrStderr())

	// top level commands

	rootCmd.AddCommand(node.NewNodeCommand(cfMgr))

	return rootCmd
}
