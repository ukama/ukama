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
	agent "github.com/ukama/ukama/systems/common/rest/client/ukamaagent"
	cenums "github.com/ukama/ukama/testing/common/enums"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/clients"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	DsimfactoryHost  string           `default:"dsimfactory:9090"`
	Service          *uconf.Service
	OrgName          string
	Http             HttpServices
	RoutineConfig    RoutineConfig
}

type RoutineConfig struct {
	Min      float64 `default:"10"`
	Normal   float64 `default:"20"`
	Max      float64 `default:"40"`
	Interval uint64  `default:"1"`
}

type HttpServices struct {
	AgentNodeGateway string `default:"http://node-gateway-ukama-agent:8080"`
	InitClient       string `default:"api-gateway-init:8080"`
}

type WMessage struct {
	Iccid     string           `json:"iccid"`
	Imsi      string           `json:"imsi"`
	Expiry    string           `json:"expiry"`
	Status    bool             `json:"status"`
	Profile   cenums.Profile   `json:"profile"`
	NodeId    string           `json:"node_id"`
	Scenario  cenums.SCENARIOS `json:"scenario"`
	CDRClient clients.CDRClient
	Agent     agent.UkamaAgentClient `json:"agent"`
}

func NewConfig(name string) *Config {
	return &Config{
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activepackage",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate",
				"event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
				"event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.deactivate",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activate",
			},
		},
	}
}
