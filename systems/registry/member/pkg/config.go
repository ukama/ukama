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
	NumberOfMembers         = "number_of_members"
	NumberOfActiveMembers   = "active_members"
	NumberOfInactiveMembers = "inactive_members"
	GaugeType               = "gauge"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	PushGateway      string           `default:"http://localhost:9091"`
	Service          *uconf.Service
	Http             HttpServices
	OwnerId          string
	OrgId            string
	OrgName          string
	MasterOrgName    string
}

type HttpServices struct {
	NucleusClient string `default:"api-gateway-nucleus:8080"`
}

var MemberMetric = []metric.MetricConfig{
	{
		Name:   NumberOfMembers,
		Type:   GaugeType,
		Labels: map[string]string{"org": ""},
		Value:  0,
	},
	{
		Name:   NumberOfActiveMembers,
		Type:   GaugeType,
		Labels: map[string]string{"org": ""},
		Value:  0,
	},
	{
		Name:   NumberOfInactiveMembers,
		Type:   GaugeType,
		Labels: map[string]string{"org": ""},
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
			Timeout:        5 * time.Second,
			ListenerRoutes: []string{"event.cloud.local.{{ .Org}}.registry.invitation.invitation.update"},
		},
	}
}
