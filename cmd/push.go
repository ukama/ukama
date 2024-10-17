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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ukama/msgcli/internal/push"
	"github.com/ukama/msgcli/util"
)

const (
	defaultOrg         = "ukamatestorg"
	defaultScope       = "local"
	defaultClusterURL  = "http://localhost:15672"
	defaultVhost       = "%2F"
	defaultExchange    = "amq.topic"
	defaultClusterUsr  = "guest"
	defaultClusterPswd = "guest"
)

var (
	eventScope = &util.EnumParam{
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
		// unbound vars
		clusterURL := viper.GetString("cluster-URL")
		if clusterURL == "" {
			clusterURL = defaultClusterURL
		}

		clusterUsr := viper.GetString("cluster-usr")
		if clusterUsr == "" {
			clusterUsr = defaultClusterUsr
		}

		clusterPswd := viper.GetString("cluster-pswd")
		if clusterPswd == "" {
			clusterPswd = defaultClusterPswd
		}

		vHost := viper.GetString("vhost")
		if vHost == "" {
			vHost = defaultVhost
		}

		exchange := viper.GetString("exchange")
		if exchange == "" {
			exchange = defaultExchange
		}

		// bound vars
		org := viper.GetString("default-org")
		scope := viper.GetString("default-scope")

		route, err := cmd.Flags().GetString("route")
		if err != nil {
			return err
		}

		msg, err := cmd.Flags().GetString("message")
		if err != nil {
			return err
		}

		cfg := &util.PushConfig{
			ClusterURL:   clusterURL,
			ClusterUsr:   clusterUsr,
			ClusterPswd:  clusterPswd,
			Vhost:        vHost,
			Exchange:     exchange,
			OutputFormat: defaultOutputFormat,
		}

		return push.Run(org, scope, route, msg, os.Stdout, cfg)
	},
}

func init() {
	eventsCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringP("org", "o", defaultOrg, "name of the org to send the event to")

	eventScope.Set(defaultScope)
	pushCmd.Flags().VarP(eventScope, "scope", "s",
		fmt.Sprintf("event scope. Must match one of the following: %q", eventScope.Values))

	pushCmd.Flags().StringP("route", "r", "", "route for the event (should match \"system.service.object.action\")")
	pushCmd.MarkFlagRequired("route")
	pushCmd.Flags().StringP("message", "m", "", "message for the event (should be in json format)")

	viper.BindPFlag("default-org", pushCmd.Flags().Lookup("org"))
	viper.BindPFlag("default-scope", pushCmd.Flags().Lookup("scope"))
}
