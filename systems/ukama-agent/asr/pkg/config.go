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
	config.BaseConfig    `mapstructure:",squash"`
	DB                   *config.Database  `default:"{}"`
	Grpc                 *config.Grpc      `default:"{}"`
	Timeout              time.Duration     `default:"3s"`
	MsgClient            *config.MsgClient `default:"{}"`
	Queue                *config.Queue     `default:"{}"`
	Service              *config.Service   `default:"{}"`
	FactoryHost          string            `default:"http://localhost:8085"`
	Reroute              string            `default:"http://localhost:8085"`
	CDRHost              string            `default:"cdr:9090"`
	IsMsgBus             bool              `default:"true"`
	Period               time.Duration     `default:"3s"`
	Monitor              bool              `default:"true"`
	AllowedTimeOfService int64             `default:"259200"` // 72 hours = 86400 *3 seconds
	OrgName              string
	OrgId                string
	Http                 HttpServices
}

type HttpServices struct {
	InitClient string `default:"api-gateway-init:8080"`
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
				"event.cloud.local.*.ukamaagent.cdr.cdr.create",
			},
		},
	}
}
