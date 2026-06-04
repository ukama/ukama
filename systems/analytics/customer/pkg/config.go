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
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	PushGateway      string           `default:"http://localhost:9091"`
	Service          *uconf.Service
	OrgName          string
	OrgId            string

	// SimLowStockThreshold is the number of available SIMs below which
	// the sim pool low_stock KPI is raised.
	SimLowStockThreshold uint32 `default:"50"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			// The customer service is a read-only consumer of the shared
			// analytics database owned by the collector service. It never
			// migrates or writes to it.
			DbName: "analytics",
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			// Read-only analytics service: registered with msgclient for
			// lifecycle only; it does not listen to any routes.
			ListenerRoutes: []string{},
		},
	}
}
