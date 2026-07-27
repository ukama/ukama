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
	evt "github.com/ukama/ukama/systems/common/events"
)

type MailerConfig struct {
	Host     string `default:"localhost"`
	Port     int    `default:"587"`
	Username string `default:""`
	Password string `default:""`
	From     string `default:"hello@dev.ukama.com"`
}

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Timeout          time.Duration    `default:"50s"`
	TemplatesPath    string           `default:"templates"`
	Service          *uconf.Service
	Mailer           *MailerConfig
	OrgName          string
	OrgId            string
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 7 * time.Second,
			ListenerRoutes: []string{
				evt.EventRoutingKey[evt.EventInviteCreate],
				evt.EventRoutingKey[evt.EventSimAllocate],
				evt.EventRoutingKey[evt.EventSimAddPackage],
			}},
	}
}
