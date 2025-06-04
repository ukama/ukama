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
)

type DeletionWorkerConfig struct {
	CheckInterval   time.Duration `default:"10m"`
	DeletionTimeout time.Duration `default:"15m"`
	MaxRetries      int           `default:"3"`
}

type Config struct {
	uconf.BaseConfig  `mapstructure:",squash"`
	DB                *uconf.Database        `default:"{}"`
	Grpc              *uconf.Grpc           `default:"{}"`
	Queue             *uconf.Queue          `default:"{}"`
	Timeout           time.Duration         `default:"10s"`
	MsgClient         *uconf.MsgClient      `default:"{}"`
	SimManagerHost    string                `default:"simmanager:9090"`
	Service           *uconf.Service
	Http              HttpServices
	OrgName           string
	OrgId             string
	DeletionWorker    *DeletionWorkerConfig `default:"{}"`
}

type HttpServices struct {
	NucleusClient string `defaut:"api-gateway-nucleus:8080"`
	InitClient    string `defaut:"api-gateway-init:8080"`
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
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sims.deletion_completed",
				"event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.delete",
				"event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.create",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activate",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.deactivate",
			},
		},
		DeletionWorker: &DeletionWorkerConfig{
			CheckInterval:   10 * time.Minute,
			DeletionTimeout: 15 * time.Minute,
			MaxRetries:      3,
		},
	}
}