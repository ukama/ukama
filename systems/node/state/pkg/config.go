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
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
	Queue            *uconf.Queue    `default:"{}"`
	Timeout          time.Duration   `default:"20s"`
	Service          *uconf.Service
	MsgClient        *config.MsgClient `default:"{}"`
	OrgName          string            `default:"ukama"`
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
				"event.cloud.local.{{ .Org}}.registry.node.node.assign",
				"event.cloud.local.{{ .Org}}.registry.node.node.release",
				"event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
				"event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
				"event.cloud.local.{{ .Org}}.node.notify.notification.store",
				"event.cloud.local.{{ .Org}}.node.notify.notification.config.ready",
				"event.cloud.local.{{ .Org}}.registry.node.node.create",
			},
		},
	}

}
