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
	Services          GrpcEndpoints       `mapstructure:"services"`
	Descriptions      ServiceDescriptions `mapstructure:"descriptions"`
	Http              HttpEndpoints       `mapstructure:"http"`
	Metrics           config.Metrics      `mapstructure:"metrics"`
	Auth              *config.Auth        `mapstructure:"auth"`
	Server            rest.HttpConfig
}

type GrpcEndpoints struct {
	Timeout          time.Duration
	Controller       string
	Configurator     string
	Software         string
	State            string
	SiteController   string
	OperationMonitor string
}

// ServiceDescriptions holds a human-readable description per gRPC service,
// returned by GET /status so consumers know which features are affected
// when a service is unavailable. Overridable via env vars
// (DESCRIPTIONS_<SERVICE>) without a code change.
type ServiceDescriptions struct {
	Controller       string
	Configurator     string
	Software         string
	State            string
	SiteController   string
	OperationMonitor string
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
			Timeout:          3 * time.Second,
			Controller:       "controller:9090",
			Software:         "software:9090",
			Configurator:     "configurator:9090",
			State:            "state:9090",
			SiteController:   "site-controller:9090",
			OperationMonitor: "operation-monitor:9090",
		},
		Descriptions: ServiceDescriptions{
			Controller:       "Node controller: node lifecycle control operations",
			Configurator:     "Node configuration: applying and versioning node configs",
			Software:         "Node software: software and app rollout to nodes",
			State:            "Node state: tracking node state transitions",
			SiteController:   "Site controller: site-level control operations",
			OperationMonitor: "Operation monitoring: tracking long-running node operations",
		},
		Http: HttpEndpoints{
			Timeout: 3 * time.Second,
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		Metrics: *config.DefaultMetrics(),
		Auth:    config.LoadAuthHostConfig("auth"),
	}
}
