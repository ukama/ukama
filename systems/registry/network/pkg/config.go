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
	metric "github.com/ukama/ukama/systems/common/metrics"
)

const (
	NumberOfNetworks = "number_of_networks"
	NumberOfSites    = "number_of_sites"
	GaugeType        = "gauge"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	PushGateway      string           `default:"http://localhost:9091"`
	Country          string           `default:"usa"`
	Language         string           `default:"en"`
	Currency         string           `default:"usd"`
	Service          *uconf.Service
	Http             HttpServices
	OrgName          string
	OrgId            string
}

type HttpServices struct {
	NucleusClient string `default:"api-gateway-nucleus:8080"`
}

var NetworkMetric = []metric.MetricConfig{
	{
		Name:   NumberOfNetworks,
		Type:   GaugeType,
		Labels: map[string]string{"org": ""},
		Value:  0,
	},
	{
		Name:   NumberOfSites,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
		},
	}
}
