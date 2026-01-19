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

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Metrics          *uconf.Metrics   `default:"{}"`
	PushGateway      string           `default:"http://localhost:9091"`
	Timeout          time.Duration    `default:"3s"`
	Queue            *uconf.Queue     `default:"{}"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	SiteHost         string           `default:"site:9090"`
	Http             HttpServices
	Service          *uconf.Service
	OrgName          string
	OrgId            string
}

type HttpServices struct {
	InventoryClient string `default:"http://api-gateway-inventory:8080"`
}

const (
	NumberOfNodes        = "number_of_nodes"
	NumberOfOnlineNodes  = "online_node_count"
	NumberOfOfflineNodes = "offline_node_count"
	GaugeType            = "gauge"
)

var NodeMetric = []metric.MetricConfig{
	{
		Name:  NumberOfNodes,
		Type:  GaugeType,
		Value: 0,
	},
	{
		Name:  NumberOfOnlineNodes,
		Type:  GaugeType,
		Value: 0,
	},
	{
		Name:  NumberOfOfflineNodes,
		Type:  GaugeType,
		Value: 0,
	},
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Host:    "msgclient-registry:9095",
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.node.health.capps.store",
				"event.cloud.local.{{ .Org}}.node.state.node.transition",
				"event.cloud.local.{{ .Org}}.registry.site.site.create",
				"event.cloud.local.{{ .Org}}.inventory.component.node.added",
			},
		},
	}
}
