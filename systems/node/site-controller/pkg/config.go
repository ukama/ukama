/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	DB                *uconf.Database `default:"{}"`
	uconf.BaseConfig `mapstructure:",squash"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Service          *uconf.Service
	MsgClient        *uconf.MsgClient `default:"{}"`
	OrgName          string           `default:"ukama"`
	Services         GrpcEndpoints
}

type GrpcEndpoints struct {
	Timeout    time.Duration `default:"3s"`
	Controller string        `default:"controller:9090"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.node.health.health.create",
				"event.cloud.local.{{ .Org}}.node.state.state.create",
				"event.cloud.local.{{ .Org}}.registry.node.node.create",
				"event.cloud.local.{{ .Org}}.registry.node.node.update",
			},
		},
		Services: GrpcEndpoints{
			Timeout:    3 * time.Second,
			Controller: "controller:9090",
		},
	}
}
