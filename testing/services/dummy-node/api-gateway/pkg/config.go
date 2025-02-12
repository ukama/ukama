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
	Services          GrpcEndpoints `mapstructure:"services"`
	HttpServices      HttpEndpoints `mapstructure:"httpServices"`
}

type GrpcEndpoints struct {
	Timeout  time.Duration
	Node     string
	Dmetrics string
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
			Timeout:  10 * time.Second,
			Node:     "dnode:9090",
			Dmetrics: "dmetrics:9090",
		},
		HttpServices: HttpEndpoints{
			Timeout:     3 * time.Second,
			NodeMetrics: "http://localhost",
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
	}
}
