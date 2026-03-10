/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
	evt "github.com/ukama/ukama/systems/common/events"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"20s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	OrgName          string           `default:"ukama"`
	WimsiHost        string           `default:"http://wimsi:8080"`
	Health           string           `default:"health:9090"`
	Apps             []*App           `default:"[]"`
	NodeGwIP         string           `default:"0.0.0.0"`
	Service          *uconf.Service
}

type App struct {
	Name        string
	Space       string
	Notes       string
	MetricsKeys []string
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 7 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.global.{{ .Org}}.hub.distributor.app.chunkready",
				evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline],
			},
		},
		Apps: []*App{
			{
				Name:        "backhaul",
				Space:       "system",
				Notes:       "Backhaul software",
				MetricsKeys: []string{"backhaul_software_cpu", "backhaul_software_memory"},
			},
			{
				Name:        "core",
				Space:       "system",
				Notes:       "Core software",
				MetricsKeys: []string{"core_software_cpu", "core_software_memory"},
			},
			{
				Name:        "metricsd",
				Space:       "system",
				Notes:       "Metrics software",
				MetricsKeys: []string{"metrics_software_cpu", "metrics_software_memory"},
			},
			{
				Name:        "switch",
				Space:       "system",
				Notes:       "Switch software",
				MetricsKeys: []string{"switch_software_cpu", "switch_software_memory"},
			},
		},
	}
}
