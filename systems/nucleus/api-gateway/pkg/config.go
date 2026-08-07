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
	Timeout time.Duration
	User    string
	Org     string
}

// ServiceDescriptions holds a human-readable description per gRPC service,
// returned by GET /status so consumers know which features are affected
// when a service is unavailable. Overridable via env vars
// (DESCRIPTIONS_<SERVICE>) without a code change.
type ServiceDescriptions struct {
	User string
	Org  string
}

type HttpEndpoints struct {
	Timeout     time.Duration
	NodeMetrics string
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
			Timeout: 20 * time.Second,
			User:    "user:9090",
			Org:     "org:9090",
		},
		Descriptions: ServiceDescriptions{
			User: "Users: platform user accounts",
			Org:  "Organizations: organization accounts and membership",
		},

		Http: HttpEndpoints{
			Timeout:     20 * time.Second,
			NodeMetrics: "http://localhost",
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		Metrics: *config.DefaultMetrics(),
		Auth:    config.LoadAuthHostConfig("auth"),
	}
}
