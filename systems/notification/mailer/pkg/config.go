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
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
	Queue            *uconf.Queue    `default:"{}"`
	Timeout          time.Duration   `default:"50s"`
	Service          *uconf.Service
	TemplatesPath    string `default:"templates"`
	MailerHost       string `default:"localhost"`
	MailerPort       int    `default:"587"`
	MailerUsername   string `default:""`
	MailerPassword   string `default:""`
	MailerFrom       string `default:"hello@dev.ukama.com"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
	}
}
