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

	pmetric "github.com/ukama/ukama/systems/common/metrics"
)

const (
	DataUsage = "data_usage"
	GaugeType = "gauge"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Metrics           *config.Metrics   `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	PushGateway       string            `default:"http://localhost:9091"`
	Queue             *config.Queue     `default:"{}"`
	Service           *config.Service   `default:"{}"`
	AsrHost           string            `default:"asr:9090"`
	IsMsgBus          bool              `default:"true"`
	OrgName           string
	OrgId             string
	Http              HttpServices
}

type HttpServices struct {
	InitClient string `default:"api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		AsrHost: "asr:9090",
		DB: &config.Database{
			DbName: ServiceName,
		},

		Grpc: &config.Grpc{
			Port: 9090,
		},

		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.create",
				"event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.update",
			},
		},
	}
}

var UsageMetrics = []pmetric.MetricConfig{
	{
		Name:   DataUsage,
		Type:   GaugeType,
		Labels: map[string]string{"package": "", "dataplan": "", "network": "", "site": "", "iccid": "", "session": ""},
		Value:  0,
	},
}
