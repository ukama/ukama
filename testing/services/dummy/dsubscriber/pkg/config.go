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
	cenums "github.com/ukama/ukama/testing/common/enums"
	pb "github.com/ukama/ukama/testing/services/dummy/dsubscriber/pb/gen"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Service          *uconf.Service
	OrgName          string
	Http             HttpServices
}

type HttpServices struct {
	InventoryClient string `defaut:"http://api-gateway-inventory:8080"`
}

type WMessage struct {
	Iccid     string         `json:"iccid"`
	Expiry    string         `json:"expiry"`
	Status    pb.Status      `json:"status"`
	Profile   cenums.Profile `json:"profile"`
	PackageId string         `json:"package_id"`
}

func NewConfig(name string) *Config {
	return &Config{
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activepackage",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate",
			},
		},
	}
}
