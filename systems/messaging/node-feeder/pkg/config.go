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

type RegistryConfig struct {
	Host           string
	TimeoutSeconds int
}

type NetConfig struct {
	Host           string
	TimeoutSeconds int
}

type DeviceNetworkConfig struct {
	Port           uint // set to 0 to bypass port addition
	TimeoutSeconds uint // timeout for one request to a device
}

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
	Registry          RegistryConfig
	Net               NetConfig
	Device            DeviceNetworkConfig
	Listener          ListenerConfig
	Metrics           *config.Metrics
}

func NewConfig() *Config {

	return &Config{
		Queue: config.Queue{
			Uri: "amqp://guest:guest@rabbitmq:5672/",
		},
		Registry: RegistryConfig{
			Host:           "localhost:9090",
			TimeoutSeconds: 3,
		},
		Net: NetConfig{
			Host:           "localhost:9090",
			TimeoutSeconds: 3,
		},
		Device: DeviceNetworkConfig{
			Port:           0,
			TimeoutSeconds: 3,
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
