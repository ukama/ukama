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

	cors "github.com/gin-contrib/cors"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type StoreConfig struct {
	ClientCert    string
	ClientKey     string
	CaCert        string
	SkipVerify    bool
	TrustInsecure bool
	CacheRepair   bool
	ErrorRetry    int
	Uncompressed  bool
}

type ChunkConfig struct {
	N            int
	Stores       []string
	MinChunkSize uint64
	MaxChunkSize uint64
	AvgChunkSize uint64
	CreateIndex  bool
	InFormat     string
	Extension    string
}

// S3Creds holds credentials or references to an S3 credentials file.
type S3Creds struct {
	AccessKey          string
	SecretKey          string
	AwsCredentialsFile string
	AwsProfile         string
	// Having an explicit aws region makes minio slightly faster because it avoids url parsing
	AwsRegion string
}

type SecurityConfig struct {
	Cert      string
	Key       string
	MutualTLS bool
	ClientCA  string
	Auth      string
}
type DistributionConfig struct {
	Address        []string
	LogFile        string
	HTTPTimeout    time.Duration
	HTTPErrorRetry int
	S3Credentials  map[string]S3Creds
	StoreCfg       StoreConfig
	Chunk          ChunkConfig
	Security       SecurityConfig
}

type MinioConfig struct {
	TimeoutSecond      int
	Endpoint           string
	AccessKey          string
	SecretKey          string
	BucketSuffix       string
	Region             string
	SkipBucketCreation bool
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           *config.Metrics
	Server            rest.HttpConfig
	Distribution      DistributionConfig
	Storage           MinioConfig
	Service           *config.Service
	Queue             *config.Queue     `default:"{}"`
	MsgClient         *config.MsgClient `default:"{}"`
	OrgName           string
}

func NewConfig(name string) *Config {
	return &Config{
		Server: rest.HttpConfig{
			Port: 8098,
			Cors: cors.Config{
				AllowOrigins: []string{"http://localhost", "https://localhost", "*"},
			},
		},

		Metrics: config.DefaultMetrics(),
		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout: 5 * time.Second,
		},

		Distribution: DistributionConfig{
			Address:        []string{":8099"},
			LogFile:        "-",
			HTTPTimeout:    20,
			HTTPErrorRetry: 1,

			StoreCfg: StoreConfig{
				ClientCert:    "",
				ClientKey:     "",
				CaCert:        "",
				SkipVerify:    false,
				TrustInsecure: false,
				CacheRepair:   false,
				ErrorRetry:    3,
				Uncompressed:  false,
			},

			Chunk: ChunkConfig{
				N:            10,
				Stores:       []string{"s3+http://localhost:9000/artifact-hub-local-test/chunks?lookup=path"},
				MinChunkSize: 64,
				MaxChunkSize: 256,
				AvgChunkSize: 64,
				CreateIndex:  true,
				InFormat:     "disk",
				Extension:    "tar.gz",
			},

			Security: SecurityConfig{
				Cert:      "",
				Key:       "",
				MutualTLS: false,
				ClientCA:  "",
				Auth:      "",
			},
		},

		Storage: MinioConfig{
			Endpoint:           "localhost:9000",
			AccessKey:          "minioadmin",
			SecretKey:          "minioadmin",
			BucketSuffix:       "local-test",
			TimeoutSecond:      3,
			SkipBucketCreation: true,
		},
	}
}
