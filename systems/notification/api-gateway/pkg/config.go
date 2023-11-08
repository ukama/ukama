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
	Services          GrpcEndpoints  `mapstructure:"services"`
	HttpServices      HttpEndpoints  `mapstructure:"httpServices"`
	Metrics           config.Metrics `mapstructure:"metrics"`
	Auth              *config.Auth   `mapstructure:"auth"`
}

type HttpEndpoints struct {
	Timeout     time.Duration
	NodeMetrics string
}

type GrpcEndpoints struct {
	Timeout time.Duration
	Mailer  string
	Notify  string
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
			Timeout: 5 * time.Second,
			Mailer:  "mailer:9090",
			Notify:  "notification-notify:9090",
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},

		HttpServices: HttpEndpoints{
			Timeout:     3 * time.Second,
			NodeMetrics: "http://localhost",
		},

		Metrics: *config.DefaultMetrics(),

		Auth: config.LoadAuthHostConfig("auth"),
	}
}