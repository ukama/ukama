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
	"github.com/ukama/ukama/systems/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           *config.Metrics
	Server            rest.HttpConfig
	Services          GrpcEndpoints `mapstructure:"services"`
	Storage           MinioConfig
	Service           *config.Service
	Queue             *config.Queue     `default:"{}"`
	MsgClient         *config.MsgClient `default:"{}"`
	OrgName           string
	OrgId             string
	IsGlobal          bool
	PushGateway       string
	Grpc              *config.Grpc `default:"{}"`
}

type GrpcEndpoints struct {
	Timeout time.Duration
	Chunker string
}

type MinioConfig struct {
	TimeoutSecond         int
	Endpoint              string
	AccessKey             string
	SecretKey             string
	BucketSuffix          string
	Region                string
	SkipBucketCreation    bool
	ArtifactTypeBucketMap map[string]string
}

type ChunkerConfig struct {
	Host          string
	TimeoutSecond int
}

func NewConfig(name string) *Config {
	return &Config{
		Metrics: config.DefaultMetrics(),
		Storage: MinioConfig{
			Endpoint:      "localhost:9000",
			AccessKey:     "minioadmin",
			SecretKey:     "minioadmin",
			BucketSuffix:  "local-test",
			TimeoutSecond: 3,
			ArtifactTypeBucketMap: map[string]string{
				"cappart": "capp",
				"certart": "cert",
			},
		},

		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 5 * time.Second,
		},

		Services: GrpcEndpoints{
			Timeout: 600 * time.Second,
			Chunker: "distributor:9090",
		},

		Grpc: &config.Grpc{
			Port:       9090,
			MaxMsgSize: 209715200,
		},
	}
}
