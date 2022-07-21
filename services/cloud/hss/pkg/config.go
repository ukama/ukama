package pkg

import (
	"github.com/ukama/ukama/services/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                config.Database
	Grpc              config.Grpc
	Queue             config.Queue
	Metrics           config.Metrics
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
		Metrics: *config.DefaultMetrics(),
	}
}
