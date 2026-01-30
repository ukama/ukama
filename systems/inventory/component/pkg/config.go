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
	uconf.BaseConfig     `mapstructure:",squash"`
	DB                   *uconf.Database  `default:"{}"`
	Grpc                 *uconf.Grpc      `default:"{}"`
	Queue                *uconf.Queue     `default:"{}"`
	Timeout              time.Duration    `default:"3s"`
	MsgClient            *uconf.MsgClient `default:"{}"`
	Service              *uconf.Service
	RepoUrl              string               `default:""`
	Token                string               `default:""`
	OwnerId              string               `default:""`
	SchedulerInterval    time.Duration        `default:"1m"`
	OrgName              string               `default:"ukama"`
	Username             string               `default:"ukama"`
	ComponentEnvironment string               `default:"production"`
	RepoPath             string               `default:"/temp/git/networks"`
	PushGateway          string               `default:"http://localhost:9091"`
	FactoryUrl           string               `default:"http://api-gateway-factory:8080"`
	NodeComponentDetails NodeComponentDetails `default:"{}"`
}

type NodeComponentDetails struct {
	ImagesURL     string                  `default:""`
	Specification string                  `default:""`
	Warranty      uint32                  `default:"1"`
	Managed       string                  `default:"true"`
	Category      string                  `default:"access"`
	Manufacturer  string                  `default:"Ukama Inc"`
	Inventory     string                  `default:"ukma-access"`
	DatasheetURL  string                  `default:"http://www.ukama.com/datasheet"`
}

var NetworkMetric = []metric.MetricConfig{}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
		},
	}
}
