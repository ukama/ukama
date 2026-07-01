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

	"github.com/gin-contrib/cors"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	DB         string
	Bridge     BrdigeConfig
	Server     rest.HttpConfig
	Auth       *config.Auth   `mapstructure:"auth"`
	Metrics    config.Metrics `mapstructure:"metrics"`
	SyncPeriod time.Duration  `default:"10s"`
}

type BrdigeConfig struct {
	Name            string
	Ip              string
	NetType         string
	Period          time.Duration `default:"2s"`
	Management      string
	SessionIdleTime time.Duration `default:"60s"`
}

func NewConfig(name string) *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		BaseConfig: config.BaseConfig{
			DebugMode: true,
		},
		DB: DefaultDBPath,
		Bridge: BrdigeConfig{
			Period:          2 * time.Second,
			SessionIdleTime: 60 * time.Second,
		},
		Server: rest.HttpConfig{
			Port: 0,
			Cors: defaultCors,
		},
		Metrics: *config.DefaultMetrics(),
		Auth: &config.Auth{
			BypassAuthMode: true,
		},
		SyncPeriod: 5 * time.Second,
	}
}
