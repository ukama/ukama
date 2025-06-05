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
	NumberOfSims    = "number_of_sims"
	ActiveCount     = "active_sim_count"
	InactiveCount   = "inactive_sim_count"
	TerminatedCount = "terminated_sim_count"
	GaugeType       = "gauge"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Queue             *config.Queue     `default:"{}"`
	Metrics           *config.Metrics   `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	PushMetricHost    string            `default:"http://localhost:9091"`
	SimPool           string            `default:"simpool:9090"`
	Registry          string            `default:"registry:9090"`
	TestAgent         string            `default:"testagent:9090"`
	OperatorAgent     string            `default:"http://operator-agent:8080"`
	Service           *config.Service
	Key               string
	OrgId             string
	OrgName           string
	Http              HttpServices
}

type HttpServices struct {
	InitClient    string `defaut:"api-gateway-init:8080"`
	NucleusClient string `defaut:"api-gateway-nucleus:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: name,
		},
		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate",
				"event.cloud.local.{{ .Org}}.payments.processor.payment.success",
				"event.cloud.local.{{ .Org}}.ukamaagent.asr.activesubscriber.delete",
				"event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.delete",
				"event.cloud.local.{{ .Org}}.ukamaagent.cdr.cdr.create",
				"event.cloud.local.{{ .Org}}.operator.cdr.cdr.create",
				"event.cloud.local.{{ .Org}}.ukamaagent.asr.subscriber.asr_cleanup_completed",
			},
		},
	}
}

var SimMetric = []pmetric.MetricConfig{
	{
		Name:   NumberOfSims,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
	{
		Name:   ActiveCount,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
	{
		Name:   InactiveCount,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
	{
		Name:   TerminatedCount,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
}
