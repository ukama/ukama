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
	metric "github.com/ukama/ukama/systems/common/metrics"
)

const (
	NumberOfActiveOrgs            = "number_of_active_org"
	NumberOfInactiveOrgs          = "number_of_inactive_org"
	NumberOfActiveMembersOfOrgs   = "number_of_active_org_members"
	NumberOfInactiveMembersOfOrgs = "number_of_inactive_org_members"
	NumberOfActiveUsers           = "number_of_active_users"
	NumberOfInactiveUsers         = "number_of_inactive_users"
	GaugeType                     = "gauge"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Service           *config.Service   `default:"{}"`
	Queue             *config.Queue     `default:"{}"`
	MsgClient         *config.MsgClient `default:"{}"`
	UserHost          string            `default:"http://user:8080"`
	OrchestratorHost  string            `default:"http://orchestrator:8080"`
	OrgName           string            `default:"ukama"`
	OwnerId           string
	OrgId             string
	Pushgateway       string `default:"http://localhost:9091"`
	InitClientHost    string `default:"http://ukama.initclient:8080"`
}

var OrgMetrics = []metric.MetricConfig{
	{
		Name:   NumberOfActiveOrgs,
		Type:   GaugeType,
		Labels: map[string]string{"state": "active"},
		Value:  0,
	},
	{
		Name:   NumberOfInactiveOrgs,
		Type:   GaugeType,
		Labels: map[string]string{"state": "inactive"},
		Value:  0,
	},
	{
		Name:   NumberOfActiveUsers,
		Type:   GaugeType,
		Labels: map[string]string{"state": "active", "service": "org"},
		Value:  0,
	},
	{
		Name:   NumberOfInactiveUsers,
		Type:   GaugeType,
		Labels: map[string]string{"state": "inactive", "service": "org"},
		Value:  0,
	},
	{
		Name:   NumberOfActiveMembersOfOrgs,
		Type:   GaugeType,
		Labels: map[string]string{"state": "active"},
		Value:  0,
	},
	{
		Name:   NumberOfInactiveMembersOfOrgs,
		Type:   GaugeType,
		Labels: map[string]string{"state": "inactive"},
		Value:  0,
	},
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: name,
		},
		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 7 * time.Second,
		},
	}
}
