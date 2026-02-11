/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig     `mapstructure:",squash"`
	EtcdHost          	 string
	DialTimeoutSecond 	 time.Duration
	Grpc                 *uconf.Grpc      `default:"{}"`
	Queue                *uconf.Queue     `default:"{}"`
	Timeout              time.Duration    `default:"3s"`
	MsgClient            *uconf.MsgClient `default:"{}"`
	SchedulerInterval    time.Duration    `default:"1m"`
	OrgName              string           `default:"ukama"`
	PrometheusHost       string           `default:"http://localhost:9079"`
	Service              *uconf.Service
	Http            	 HttpServices
}

type HttpServices struct {
	InitClient string `default:"http://api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		EtcdHost:          "localhost:2379",
		DialTimeoutSecond: 5 * time.Second,
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
		},
	}
}
