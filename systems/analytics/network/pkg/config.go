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

	// Thresholds used to derive flags from the latest rollups/samples.
	NetworkLatencyThresholdMs float64 `default:"100"`
	BatteryCriticalPercent    float64 `default:"20"`
	TelemetryFreshSeconds     int64   `default:"600"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			// All analytics services share the single "analytics" database.
			// The collector service owns the schema (it is the only writer
			// and the only service running AutoMigrate); network is read-only.
			DbName: SystemName,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			// Read-only analytics service: no msgbus event handling, so no
			// listener routes are registered. The client is still created
			// for service lifecycle registration (Register/Start).
			ListenerRoutes: []string{},
		},
	}
}
