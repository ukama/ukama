package pkg

import (
	"github.com/ukama/ukamaX/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                config.Database
	Grpc              config.Grpc
	Queue             config.Queue
	SimManager        SimManager
	SimTokenKey       string
}

type SimManager struct {
	Host     string
	Name     string
	Disabled bool
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
		Queue: config.Queue{
			Uri: "amqp://guest:guest@localhost:5672",
		},
		Grpc: config.Grpc{
			Port: 9090,
		},
		SimManager: SimManager{
			Host:     "localhost:9090",
			Name:     "SimManager",
			Disabled: false,
		},
		SimTokenKey: "11111111111111111111111111111111",
	}
}
