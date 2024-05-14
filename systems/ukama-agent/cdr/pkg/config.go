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

	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	Queue             *config.Queue     `default:"{}"`
	Service           *config.Service   `default:"{}"`
	AsrHost           string            `default:"asr:9090"`
	OrgName           string            `default:"ukama"`
	OrgId             string            `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
	IsMsgBus          bool              `default:"true"`
}

type SimManager struct {
	Host string
	Name string
}

type GrpcEndPoints struct {
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: ServiceName,
		},

		Grpc: &config.Grpc{
			Port: 9090,
		},

		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.create",
				"event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.update",
			},
		},
		AsrHost: "asr:9090",
	}
}
