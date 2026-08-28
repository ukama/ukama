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
	Engine           *EngineConfig `default:"{}"`
	PushGateway      string        `default:"http://localhost:9091"`
	OrgName          string
	OrgId            string
}

// WindowConfig is the shared pipeline grid — MUST match ingest/aggregator.
type WindowConfig struct {
	W time.Duration `default:"5m"`
}

type EngineConfig struct {
	// SweepInterval is the ledger sweeper cadence: it recovers any window
	// whose window.ready event was lost.
	SweepInterval time.Duration `default:"60s"`
	// Catchup is how far back the sweeper looks for windows whose inputs
	// became ready late. A duration, so the horizon is independent of Window.W.
	Catchup time.Duration `default:"1h"`
	// CatchupWindows pins an explicit window count. 0 = derive from Catchup.
	CatchupWindows int64 `default:"0"`
	// SpecsDir holds the KPI specs (*.yaml).
	SpecsDir string `default:"configs/kpis"`
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
				schema.WindowReadyRoute,
			},
		},
	}
}
