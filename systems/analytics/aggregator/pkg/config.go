/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/analytics/schema"
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Service          *uconf.Service
	Window           *WindowConfig `default:"{}"`
	Rollup           *RollupConfig `default:"{}"`
	PushGateway      string        `default:"http://localhost:9091"`
	OrgName          string
	OrgId            string
}

// WindowConfig is the shared pipeline grid — MUST match ingest/analysis.
type WindowConfig struct {
	W time.Duration `default:"5m"`
}

type RollupConfig struct {
	// Timezone spans (day/week/month boundaries) are computed in.
	Timezone string `default:"UTC"`
	// SpecsDir holds the KPI specs (rollup_ops, output metadata).
	SpecsDir string `default:"configs/kpis"`
	// ReportsDir holds the performance report specs.
	ReportsDir string `default:"configs/reports"`
	// SweepInterval recomputes current spans (covers lost events + partial
	// span refresh).
	SweepInterval time.Duration `default:"120s"`
	// FlatThresholdPct: |change_pct| below this reads as a flat trend.
	FlatThresholdPct float64 `default:"1.0"`
	// ReportWindow is the trailing window the performance-report composer
	// aggregates over. Reports need more history than the headline KPI tiles
	// to show stable per-entity stats, so this is a longer, config-driven
	// rolling window (default 8 weeks = 1344h) rather than the UI filter.
	ReportWindow time.Duration `default:"1344h"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: SystemName,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				schema.KpiComputedRoute,
			},
		},
	}
}
