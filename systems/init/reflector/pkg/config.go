/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"fmt"
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Timeout          time.Duration    `default:"20s"`
	OrgName          string           `default:"ukama"`
	Service          *uconf.Service
	ServiceConfig    ServiceConfig
}

type ServiceConfig struct {
	Scheme 				string `default:"http"`
	Host   				string `default:"127.0.0.1"`
	Port   				int32  `default:"8088"`
	Seed   				int64  `default:"1"`
	MaxUploadBytes 		int64  `default:"67108864"`
	MaxDownloadBytes 	int64  `default:"67108864"`
}

const ReflectorBasePath = "/reflector"

func NewConfig(name string) *Config {
	return &Config{
		Service: uconf.LoadServiceHostConfig(name),
	}
}

// BaseURL returns the reflector base URL consumed by backhaul discovery.
// Backhaul appends "/v1/ping", "/v1/download/:bytes", and "/v1/upload" itself.
func (c *Config) BaseURL() string {
	return fmt.Sprintf("%s://%s:%d%s",
		c.ServiceConfig.Scheme,
		c.ServiceConfig.Host,
		c.ServiceConfig.Port,
		ReflectorBasePath,
	)
}