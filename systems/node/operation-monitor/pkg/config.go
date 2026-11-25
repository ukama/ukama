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
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Timeout          time.Duration    `default:"20s"`
	Service          *uconf.Service
	OrgName          string
	OrgId            string
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
				"event.cloud.local.{{ .Org}}.node.state.node.transition",
				"event.cloud.local.{{ .Org}}.node.health.apps.changed",
				"event.cloud.global.{{ .Org}}.operation.manager.operation.completed",
				"event.cloud.global.{{ .Org}}.operation.manager.operation.failed",
				"event.cloud.global.{{ .Org}}.operation.manager.operation.timeout",
				"event.cloud.global.{{ .Org}}.operation.manager.operation.cancelled",
			},
		},
	}
}
