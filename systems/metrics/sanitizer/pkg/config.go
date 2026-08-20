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
	RawNodeActiveSubscribers = "trx_lte_core_active_ue"
	RawComNodeUptimeSeconds  = "com_generic_system_uptime_seconds"
	RawCtlNodeUptimeSeconds  = "ctl_generic_system_uptime_seconds"
)

const (
	NodeActiveSubscribers = "active_subscribers"
	ComNodeUptime         = "com_node_uptime"
	CtlNodeUptime         = "ctl_node_uptime"

	GaugeType = "gauge"
)

const (
	NodeIdLabel  = "node_id"
	SiteIdLabel  = "site"
	NetworkLabel = "network"
)

type SanitizedMetric struct {
	Name          string
	SkipUnchanged bool
}

var SanitizedMetrics = map[string]SanitizedMetric{
	RawNodeActiveSubscribers: {Name: NodeActiveSubscribers, SkipUnchanged: true},
	RawComNodeUptimeSeconds:  {Name: ComNodeUptime},
	RawCtlNodeUptimeSeconds:  {Name: CtlNodeUptime},
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Grpc              *config.Grpc      `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	Queue             *config.Queue     `default:"{}"`
	Service           *config.Service   `default:"{}"`
	PushGateway       string            `default:"http://localhost:9091"`
	IsMsgBus          bool              `default:"true"`
	Metrics           *config.Metrics   `default:"{}"`
	Org               string            `default:""`
	Http              HttpServices
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
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.registry.node.node.assign",
				"event.cloud.local.{{ .Org}}.registry.node.node.release",
			},
		},
	}
}

var NodeMetrics = []pmetric.MetricConfig{
	{
		Name:    NodeActiveSubscribers,
		Type:    GaugeType,
		Details: "Number of active subscribers per node, with site and network labels appended",
		Labels:  map[string]string{NodeIdLabel: "", SiteIdLabel: "", NetworkLabel: ""},
		Value:   0,
	},
	{
		Name:    ComNodeUptime,
		Type:    GaugeType,
		Details: "System uptime in seconds reported by com nodes, with site and network labels appended",
		Labels:  map[string]string{NodeIdLabel: "", SiteIdLabel: "", NetworkLabel: ""},
		Value:   0,
	},
	{
		Name:    CtlNodeUptime,
		Type:    GaugeType,
		Details: "System uptime in seconds reported by ctl nodes, with site and network labels appended",
		Labels:  map[string]string{NodeIdLabel: "", SiteIdLabel: "", NetworkLabel: ""},
		Value:   0,
	},
}
