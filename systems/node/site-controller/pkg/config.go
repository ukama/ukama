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
	uconf.BaseConfig 	`mapstructure:",squash"`
	DB               	*uconf.Database  `default:"{}"`
	Grpc             	*uconf.Grpc      `default:"{}"`
	Queue            	*uconf.Queue     `default:"{}"`
	Timeout          	time.Duration    `default:"20s"`
	MsgClient        	*uconf.MsgClient `default:"{}"`
	HealthHost          string        `default:"health:9090"`
	OrgName             string
	Service             *uconf.Service
	Http                HttpServices
	ReconcileInterval    time.Duration `default:"30s"`
	ReconcileMaxRetries  int           `default:"3"`
	ComponentSyncDelay   time.Duration `default:"30s"`
}
type HttpServices struct {
	InitClient    string `default:"api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Timeout: 3 * time.Second,
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 7 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.registry.site.site.create",
				"event.cloud.local.ukama.node.health.report.store",
			},
		},
	}
}
