/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	uconf "github.com/ukama/ukama/systems/common/config"
	enpkg "github.com/ukama/ukama/systems/notification/event-notify/pkg"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
	Service          *uconf.Service
	OrgName          string
	OrgId            string
	EventNotifyHost  string `default:"localhost:9090"`
	Http             HttpServices
}

type HttpServices struct {
	InitClient string `defaut:"api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: enpkg.ServiceName,
		},
		Service: uconf.LoadServiceHostConfig(name),
	}
}
