/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
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
	MsgClient        *uconf.MsgClient `default:"{}"`
	Service          *uconf.Service
	Http             HttpServices
	Window           *WindowConfig `default:"{}"`
	Engine           *EngineConfig `default:"{}"`
	PushGateway      string        `default:"http://localhost:9091"`
	OrgName          string
	OrgId            string
}

type HttpServices struct {
	InitClient string `default:"api-gateway-init:8080"`
}

// WindowConfig is the shared pipeline grid — MUST be identical across
// ingest, analysis and aggregator (env: WINDOW_W, WINDOW_L).
type WindowConfig struct {
	W time.Duration `default:"5m"`
	L time.Duration `default:"10m"`
}

type EngineConfig struct {
	// TickInterval is how often the scheduler looks for eligible windows.
	TickInterval time.Duration `default:"30s"`
	// CatchupWindows bounds how far back a cold/behind ingest reaches.
	CatchupWindows int64 `default:"12"`
	// SpecsDir holds the source specs (*.yaml).
	SpecsDir string `default:"configs/sources"`
	// RetryCount is the per-pull retry budget within an eligibility period.
	RetryCount int `default:"3"`
	// ResolverTTL is the initclient resolution cache TTL.
	ResolverTTL time.Duration `default:"10m"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: SystemName,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			// ingest only publishes; it consumes no routes.
			ListenerRoutes: []string{},
		},
	}
}
