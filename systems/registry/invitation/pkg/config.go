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

	"github.com/ukama/ukama/systems/common/config"
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig     `mapstructure:",squash"`
	DB                   *uconf.Database   `default:"{}"`
	Grpc                 *uconf.Grpc       `default:"{}"`
	Queue                *uconf.Queue      `default:"{}"`
	Timeout              time.Duration     `default:"3s"`
	MsgClient            *config.MsgClient `default:"{}"`
	Service              *uconf.Service
	OrgName              string
	TemplateName         string `default:"member-invite"`
	InvitationExpiryTime uint   `default:"24"`
	AuthLoginbaseURL     string `default:"http://localhost:4455/auth/login"`
	Http                 HttpServices
}

type HttpServices struct {
	NucleusClient string `defaut:"api-gateway-nucleus:8080"`
	InitClient    string `defaut:"api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: name,
		},
		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 7 * time.Second,
		},
	}
}
