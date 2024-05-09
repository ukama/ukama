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
	DB                string
	Bridge            BrdigeConfig
	Server            rest.HttpConfig
	HttpServices      HttpEndpoints  `mapstructure:"httpServices"`
	Auth              *config.Auth   `mapstructure:"auth"`
	Metrics           config.Metrics `mapstructure:"metrics"`
	SyncPeriod        time.Duration  `default:"10s"`
}

type BrdigeConfig struct {
	Name            string `default:"br0"`
	Ip              string `default:"10.10.10.1"`
	NetType         string
	Period          time.Duration `default:"2s"`
	Management      string        `default:"/usr/local/var/run/openvswitch"`
	SessionIdleTime time.Duration `default:"60s"`
}

type HttpEndpoints struct {
	Timeout time.Duration
	Policy  string
	Noded   string `default:"localhost:9090"`
}

func NewConfig(name string) *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		BaseConfig: config.BaseConfig{
			DebugMode: true,
		},
		Bridge: BrdigeConfig{
			Name:            "br0",
			Period:          2 * time.Second,
			Management:      "/usr/local/var/run/openvswitch",
			SessionIdleTime: 60 * time.Second,
		},
		DB: name,
		Server: rest.HttpConfig{
			Port: 8090,
			Cors: defaultCors,
		},

		HttpServices: HttpEndpoints{
			Timeout: 3 * time.Second,
			Policy:  "http://localhost:8087",
		},
		Metrics: *config.DefaultMetrics(),
		//Auth:    config.LoadAuthHostConfig("auth"),
		Auth: &config.Auth{
			BypassAuthMode: true,
		},
		SyncPeriod: 5 * time.Second,
	}
}
