/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
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
	OrgName          string
	OrgId            string
	RollupScheduler  RollupScheduler
}

type RollupScheduler struct {
	Enabled      bool          `default:"true"`
	Interval     time.Duration `default:"5m"`
	LookbackDays int           `default:"30"`
}

type HttpServices struct {
	RegistryClient   string `default:"api-gateway-registry:8080"`
	SubscriberClient string `default:"api-gateway-subscriber:8080"`
	DataplanClient   string `default:"api-gateway-dataplan:8080"`
	MetricsClient    string `default:"api-gateway-metrics:8080"`
	NodeClient       string `default:"api-gateway-node:8080"`
	InventoryClient  string `default:"api-gateway-inventory:8080"`
	BillingClient    string `default:"api-gateway-billing:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			// The collector is the single writer (and the only service that
			// runs AutoMigrate) on the shared "analytics" database. Business,
			// customer and network services read the same database, so its
			// name is fixed instead of being derived from the service name.
			DbName: "analytics",
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			// ListenerRoutes are provided via deployment config; the full
			// list of consumed events is documented in pkg/server/event.go.
		},
	}
}
