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

	"github.com/spf13/cobra"

	"github.com/ukama/msgcli/internal/push"
	"github.com/ukama/msgcli/util"
)

var (
	eventScope = util.EnumParam{
		Values: []string{"local", "global"},
	}
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push events",
	Long: `The push command pushes events to a running service throught it's associated
message client.`,

	Aliases:      []string{"p"},
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		org, err := cmd.Flags().GetString("org")
		if err != nil {
			return err
		}

		route, err := cmd.Flags().GetString("route")
		if err != nil {
			return err
		}

		msg, err := cmd.Flags().GetString("message")
		if err != nil {
			return err
		}

		return push.Run(org, route, msg)
	},
}

func init() {
	eventsCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringP("org", "o", "ukama-testorg", "name of the org to send the event to (default \"ukama-testorg\")")

	pushCmd.Flags().StringP("route", "r", "", "route for the event (should match \"system.service.object.action\")")
	// pushCmd.MarkFlagRequired("route")

	pushCmd.Flags().StringP("message", "m", "", "message for the event (should be in json format)")
	pushCmd.Flags().VarP(&eventScope, "scope", "s",
		fmt.Sprintf("event scope. Must match one of the following: %q (default \"local\" )",
			eventScope.Values))
}
