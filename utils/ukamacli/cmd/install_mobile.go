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

	"github.com/ukama/ukama/utils/ukamacli/internal/resource/mobile"
)

var installMobileCmd = &cobra.Command{
	Use:   "mobile",
	Short: "Install mobile",
	Long:  `The install mobile command installs the given Ukama mobiie app.`,

	Aliases:      []string{"m"},
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return mobile.Install()
	},
}

func init() {
	installCmd.AddCommand(installMobileCmd)
}
