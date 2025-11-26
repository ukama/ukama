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
	mb "github.com/ukama/ukama/systems/common/msgbus"
)

type ListenerConfig struct {
	ExecutionRetryCount int64           // max retries count
	RetryPeriodSec      int             // how long request waits after failure to try again
	Threads             int             // how many go routines run message processor
	Exchange            string          // exchange
	Routes              []mb.RoutingKey // routes
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Queue             config.Queue
	Listener          ListenerConfig
	Metrics           *config.Metrics
	OrgName           string `default:"ukama"`
	TimeoutSeconds    int    `default:"3"`
	DevicePort        int    `default:"0"`
	Net               string `default:"nns:9090"`
	Registry          string `default:"api-gateway-registry:8080"`
	Http              HttpServices
}

type HttpServices struct {
	InitClient string `default:"api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		Queue: config.Queue{
			Uri: "amqp://guest:guest@rabbitmq:5672/",
		},
		Listener: ListenerConfig{
			ExecutionRetryCount: 3,
			RetryPeriodSec:      30,
			Threads:             3,
			Routes:              []mb.RoutingKey{"request.cloud.local.*.*.*.nodefeeder.publish"},
			Exchange:            "amq.topic",
		},
		Metrics: config.DefaultMetrics(),
		//request.cloud.local.ukamaorg.messaging.eventgenerator.nodefeeder.publish
	}
}
