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
	Server            rest.HttpConfig
	Services          GrpcEndpoints       `mapstructure:"services"`
	Descriptions      ServiceDescriptions `mapstructure:"descriptions"`
	Http              HttpEndpoints       `mapstructure:"http"`
	Metrics           config.Metrics      `mapstructure:"metrics"`
	Auth              *config.Auth        `mapstructure:"auth"`
}

type GrpcEndpoints struct {
	Timeout  time.Duration
	Package  string
	Baserate string
	Rate     string
}

// ServiceDescriptions holds a human-readable description per gRPC service,
// returned by GET /status so consumers know which features are affected
// when a service is unavailable. Overridable via env vars
// (DESCRIPTIONS_<SERVICE>) without a code change.
type ServiceDescriptions struct {
	Package  string
	Baserate string
	Rate     string
}

type HttpEndpoints struct {
	Timeout time.Duration
}

func NewConfig() *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},
		Services: GrpcEndpoints{
			Timeout:  5 * time.Second,
			Package:  "package:9090",
			Baserate: "baserate:9090",
			Rate:     "rate:9090",
		},
		Descriptions: ServiceDescriptions{
			Package:  "Data packages: creating and managing data plans and packages",
			Baserate: "Base rates: managing operator base rates",
			Rate:     "Rates: computing subscriber rates and markups",
		},
		Http: HttpEndpoints{
			Timeout: 5 * time.Second,
		},
		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		Metrics: *config.DefaultMetrics(),
		Auth:    config.LoadAuthHostConfig("auth"),
	}
}
