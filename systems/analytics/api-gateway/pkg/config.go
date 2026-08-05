/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"time"

	"github.com/gin-contrib/cors"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Server            rest.HttpConfig
	Services          GrpcEndpoints       `mapstructure:"services"`
	Descriptions      ServiceDescriptions `mapstructure:"descriptions"`
	Metrics           config.Metrics      `mapstructure:"metrics"`
	Auth              *config.Auth        `mapstructure:"auth"`
}

type GrpcEndpoints struct {
	Timeout    time.Duration
	Aggregator string
	Ingest     string
	Analysis   string
}

// ServiceDescriptions holds a human-readable description per gRPC service,
// returned by GET /status so consumers know which features are affected
// when a service is unavailable. Overridable via env vars
// (DESCRIPTIONS_<SERVICE>) without a code change.
type ServiceDescriptions struct {
	Aggregator string
	Ingest     string
	Analysis   string
}

func NewConfig() *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		Services: GrpcEndpoints{
			Timeout:    20 * time.Second,
			Aggregator: "aggregator:9090",
			Ingest:     "ingest:9090",
			Analysis:   "analysis:9090",
		},
		Descriptions: ServiceDescriptions{
			Aggregator: "Analytics aggregation: aggregating platform analytics",
			Ingest:     "Analytics ingestion: receiving and storing incoming analytics data",
			Analysis:   "Analytics analysis: running analysis over collected analytics data",
		},
		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		Metrics: *config.DefaultMetrics(),
		Auth:    config.LoadAuthHostConfig("auth"),
	}
}
