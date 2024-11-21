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
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
	Queue            *uconf.Queue    `default:"{}"`
	Timeout          time.Duration   `default:"20s"`
	Service          *uconf.Service
	MsgClient        *uconf.MsgClient `default:"{}"`
	OrgName          string
	OrgId            string
	ConfigPath       string
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
				"event.cloud.local.{{ .Org}}.node.notify.notification.store",
				"event.cloud.local.{{ .Org}}.node.state.node.force",
				evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline],
				evt.NodeStateEventRoutingKey[evt.NodeStateEventOffline],
				evt.NodeStateEventRoutingKey[evt.NodeStateEventAssign],
				evt.NodeStateEventRoutingKey[evt.NodeStateEventRelease],
			},
		},
	}
}
