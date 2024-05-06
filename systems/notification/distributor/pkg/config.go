/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"github.com/ukama/ukama/systems/common/config"
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
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: enpkg.ServiceName,
		},
		Service: uconf.LoadServiceHostConfig(name),
	}
}
