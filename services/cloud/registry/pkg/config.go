package pkg

import (
	"github.com/ukama/ukama/services/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                config.Database
	Grpc              config.Grpc
	BootstrapAuth     Auth
	BootstrapUrl      string
	DeviceGatewayHost string // should be an IP
	Queue             config.Queue
	Debug             Debug
}

type Auth struct {
	ClientId     string
	ClientSecret string
	Audience     string
	GrantType    string
	Auth0Host    string
}

type Debug struct {
	DisableBootstrap bool
}

func NewConfig() *Config {
	return &Config{
		DB: config.DefaultDatabase(),
		Grpc: config.Grpc{
			Port: 9090,
		},
		BootstrapAuth: Auth{
			Audience:  "bootstrap.ukama.com",
			GrantType: "client_credentials",
		},
		Debug: Debug{
			DisableBootstrap: false,
		},
		Queue: config.Queue{
			Uri: "amqp://guest:guest@localhost:5672",
		},
	}
}
