/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package internal

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig  `mapstructure:",squash"`
	DB                *uconf.Database  `default:"{}"`
	Grpc              *uconf.Grpc      `default:"{}"`
	Metrics           *uconf.Metrics   `default:"{}"`
	Timeout           time.Duration    `default:"3s"`
	Queue             *uconf.Queue     `default:"{}"`
	MsgClient         *uconf.MsgClient `default:"{}"`
	Service           *uconf.Service
	System            string `default:"testing"`
	OrgName           string
	SchedulerInterval time.Duration `default:"10s"`
	StripeKey         string
	PawapayKey        string
	PaymentsHost      string `default:"http://payments-api-gateway-payments-1:8080"`
	WebhooksHost      string `default:"http://webhooks-api-gateway-webhooks-1:8080"`
	PawapayHost       string `default:"https://api.sandbox.pawapay.cloud"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},

		Service: uconf.LoadServiceHostConfig(name),

		MsgClient: &uconf.MsgClient{
			Host:    "msg-client-testing:9095",
			Timeout: 5 * time.Second,
		},
	}
}
