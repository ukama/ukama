/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Grpc             *uconf.Grpc   `default:"{}"`
	Queue            *uconf.Queue  `default:"{}"`
	Timeout          time.Duration `default:"20s"`
	Service          *uconf.Service
	MsgClient        *uconf.MsgClient `default:"{}"`
	OrgName          string           `default:"ukama"`
	Lookup           string           `default:"lookup:9090"`
	Http             HttpServices
	DNSMap           map[string]string
}

type HttpServices struct {
	InitClient string `default:"api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout:        7 * time.Second,
			ListenerRoutes: []string{},
		},
		DNSMap: map[string]string{
			"ukama": "dev.bootstrap.ukama.com",
		},
	}
}
