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
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database   `default:"{}"`
	Grpc             *uconf.Grpc       `default:"{}"`
	Queue            *uconf.Queue      `default:"{}"`
	Metrics          *uconf.Metrics    `default:"{}"`
	Timeout          time.Duration     `default:"3s"`
	HeathCheck       HeathCheckRoutine `default:"{}"`
	System           string            `default:"init"`
	OrgName          string
	MasterOrgName    string
	Shovel           Shovel
	MsgBus           MsgBus
}

type MsgBus struct {
	ManagementUri string
	User          string
	Password      string
}

type Shovel struct {
	SrcProtocol     string `json:"src_protocol" default:"amqp091"`
	DestProtocol    string `default:"amqp091" json:"src-protocol"`
	SrcExchange     string `default:"amq.topic" json:"src-exchange"`
	SrcExchangeKey  string `json:"src-exchange-key,omitempty"`
	DestExchange    string `default:"amq.topic" json:"dest-exchange,omitempty"`
	DestExchangeKey string `json:"dest-exchange-key,omitempty"`
	DestQueue       string `json:"dest-queue,omitempty"`
	SrcQueue        string `json:"src-queue,omitempty"`
	SrcUri          string `json:"src-uri"`
	DestUri         string `json:"dest-uri"`
}

type HeathCheckRoutine struct {
	Period      time.Duration `default:"60s"`
	AllowedMiss uint32        `default:"3"`
}

func NewConfig() *Config {
	return &Config{
		Grpc: &uconf.Grpc{
			Port: 9095,
		},
	}
}
