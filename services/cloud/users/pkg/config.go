package pkg

import (
	"time"

	"github.com/ukama/ukama/services/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                config.Database
	Grpc              config.Grpc
	HssHost           string
	SimManager        SimManager
	SimTokenKey       string
	Queue             config.Queue
	KratoAdminUrl     string
}

type SimManager struct {
	Host    string
	Name    string
	Timeout time.Duration
}

func NewConfig() *Config {
	return &Config{
		DB: config.Database{
			Host:       "localhost",
			Password:   "Pass2020!",
			DbName:     ServiceName,
			Username:   "postgres",
			Port:       5432,
			SslEnabled: false,
		},
		Grpc: config.Grpc{
			Port: 9090,
		},
		HssHost: "localhost:9090",
		SimManager: SimManager{
			Host:    "localhost:9090",
			Name:    "SimManager",
			Timeout: 3 * time.Second,
		},
		KratoAdminUrl:"localhost",
		SimTokenKey: "11111111111111111111111111111111",
	}
}
