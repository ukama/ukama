package pkg

import (
	bootstrap "github.com/ukama/ukama/services/bootstrap/client"
	"github.com/ukama/ukama/services/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                config.Database
	Grpc              config.Grpc
	BootstrapAuth     bootstrap.AuthConfig
	BootstrapUrl      string
	Queue             config.Queue
	Debug             bootstrap.DebugConf
}

func NewConfig() *Config {
	return &Config{
		DB: config.DefaultDatabaseName(ServiceName),
		Grpc: config.Grpc{
			Port: 9090,
		},
		Queue: config.Queue{
			Uri: "amqp://guest:guest@localhost:5672",
		},
	}
}
