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
	CountType = "gauge"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Metrics           *config.Metrics   `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	PushGatewayHost   string            `default:"http://localhost:9091"`
	Queue             *config.Queue     `default:"{}"`
	Service           *config.Service   `default:"{}"`
	AsrHost           string            `default:"asr:9090"`
	IsMsgBus          bool              `default:"true"`
	OrgName           string
	OrgId             string
}

type GrpcEndPoints struct {
}

func NewConfig(name string) *Config {
	return &Config{
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
		Type:   CountType,
		Labels: map[string]string{"package": "", "dataplan": "", "network": ""},
		Value:  0,
	},
}
