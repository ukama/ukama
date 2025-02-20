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
	Queue             *config.Queue     `default:"{}"`
	Metrics           *config.Metrics   `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	Service           *config.Service
	System            string `default:"billing"`
	LagoHost          string `default:"localhost"`
	LagoPort          uint   `default:"3000"`
	LagoAPIKey        string
	OrgName           string
	OrgId             string
	WebhookUrl        string
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: name,
		},

		Service: config.LoadServiceHostConfig(name),

		MsgClient: &config.MsgClient{
			Host:    "msg-client-billing:9095",
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.dataplan.package.package.create",
				"event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.create",
				"event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.update",
				"event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.delete",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activepackage",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.expirepackage",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.usage",
				"event.cloud.local.{{ .Org}}.operator.cdr.sim.fakeusage",

				// TODO: we need to add the relevant arch in order to support listening
				// global events from Ukama to a local deployed org.
				"event.cloud.local.{{ .Org}}.orchestrator.constructor.org.deploy",
				"event.cloud.local.{{ .Org}}.inventory.accounting.accounting.sync",
			},
		},
	}
}
