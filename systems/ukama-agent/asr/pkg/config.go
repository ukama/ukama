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
	AsrHost              string            `default:"localhost"`
	DataplanHost         string            `default:"http://localhost:8085"`
	NetworkHost          string            `default:"http://localhost:8085"`
	FactoryHost          string            `default:"http://localhost:8085"`
	Reroute              string            `default:"http://localhost:8085"`
	CDRHost              string            `default:"http://localhost:8085"`
	OrgName              string            `default:"ukama"`
	OrgId                string            `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
	IsMsgBus             bool              `default:"true"`
	Period               time.Duration     `default:"3s"`
	Monitor              bool              `default:"true"`
	AllowedTimeOfService int64             `default:"259200"` // 72 hours = 86400 *3 seconds
}

type SimManager struct {
	Host string
	Name string
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
		DataplanHost: "http://192.168.0.14:8085",
		NetworkHost:  "http://192.168.0.14:8085",
		FactoryHost:  "http://192.168.0.14:8085",
		CDRHost:      "http://192.168.0.14:8085",
	}
}
