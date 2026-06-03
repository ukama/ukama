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
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
	Service          *uconf.Service
	Timeout          time.Duration `default:"3s"`
	PushGateway      string        `default:"http://localhost:9091"`
	OrgName          string
	OrgId            string

	// SimLowStockThreshold is the number of available SIMs below which
	// the sim pool low_stock KPI is raised.
	SimLowStockThreshold uint32 `default:"50"`

	// Thresholds used to derive network flags from the latest rollups/samples.
	NetworkLatencyThresholdMs float64 `default:"100"`
	BatteryCriticalPercent    float64 `default:"20"`
	TelemetryFreshSeconds     int64   `default:"600"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			// Analytics is the single read service behind api-gateway.
			// Collector owns writes and migrations.
			DbName: SystemName,
		},
		Service: uconf.LoadServiceHostConfig(name),
	}
}
