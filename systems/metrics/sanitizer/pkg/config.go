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
	NodeActiveSubscribers = "active_subscribers_per_node"
	GaugeType             = "gauge"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Grpc              *config.Grpc `default:"{}"`
	Http              HttpServices
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	Queue             *config.Queue     `default:"{}"`
	Service           *config.Service   `default:"{}"`
	PushGatewayHost   string            `default:"http://localhost:9091"`
	IsMsgBus          bool              `default:"true"`
	Metrics           *config.Metrics   `default:"{}"`
	Org               string            `default:""`
	OrgName           string
}

type HttpServices struct {
	InitClient string `default:"api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		Grpc: &config.Grpc{
			Port: 9090,
		},
		Metrics: &config.Metrics{
			Port: 10251,
		},
		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 5 * time.Second,
		},
	}
}

var NodeActiveSubscribersMetric = []pmetric.MetricConfig{
	{
		Name:   NodeActiveSubscribers,
		Type:   GaugeType,
		Labels: map[string]string{"nodeid": "", "site": "", "network": ""},
		Value:  0,
	},
}
