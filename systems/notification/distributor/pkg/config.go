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

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Grpc             *uconf.Grpc   `default:"{}"`
	Timeout          time.Duration `default:"3s"`
	Service          *uconf.Service
	OrgName          string
	OrgId            string
}

func NewConfig(name string) *Config {
	return &Config{
		Service: uconf.LoadServiceHostConfig(name),
	}
}
