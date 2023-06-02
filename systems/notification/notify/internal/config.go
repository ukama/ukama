package internal

import (
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           config.Metrics
	Server            rest.HttpConfig
	DB                config.Database
	Queue             config.Queue
}

var ServiceConfig *Config

func NewConfig() *Config {

	return &Config{
		Server: rest.DefaultHTTPConfig(),

		Queue: config.Queue{
			Uri: "amqp://guest:guest@localhost:5672",
		},

		DB: config.DefaultDatabaseName(ServiceName),
	}
}
