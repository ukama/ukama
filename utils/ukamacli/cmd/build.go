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
)

const (
	defaultOutputFormat = "json"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build resources",
	Long: `Build the targeted resource.

`,

	Aliases: []string{"b"},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
