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
	Auth              *config.Auth  `mapstructure:"auth"`
}

type GrpcEndpoints struct {
	Timeout     time.Duration
	Dsubscriber string
	Dsimfactory string
	Dcontroller string
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
			Timeout:     3 * time.Second,
			Dsubscriber: "dsubscriber:9090",
			Dsimfactory: "dsimfactory:9090",
			Dcontroller: "dcontroller:9090",
		},
		HttpServices: HttpEndpoints{
			Timeout: 3 * time.Second,
		},
		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
	}
}
