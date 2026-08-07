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
	Timeout   time.Duration
	Report    string
	Collector string
}

// ServiceDescriptions holds a human-readable description per gRPC service,
// returned by GET /status so consumers know which features are affected
// when a service is unavailable. Overridable via env vars
// (DESCRIPTIONS_<SERVICE>) without a code change.
type ServiceDescriptions struct {
	Report    string
	Collector string
}

type HttpEndpoints struct {
	Timeout time.Duration
	Files   string
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
			Timeout:   3 * time.Second,
			Report:    "billing-report:9090",
			Collector: "collector:9090",
		},
		Descriptions: ServiceDescriptions{
			Report:    "Billing reports: generating billing reports and invoices",
			Collector: "Billing collection: collecting usage and billing events for invoicing",
		},
		Http: HttpEndpoints{
			Timeout: 3 * time.Second,
			Files:   `http://billing-report:3000`,
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		Metrics: *config.DefaultMetrics(),
		Auth:    config.LoadAuthHostConfig("auth"),
	}
}
