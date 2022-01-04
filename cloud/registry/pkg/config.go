package pkg

import (
	"github.com/ukama/ukamaX/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                config.Database
	Grpc              config.Grpc
	BootstrapAuth     Auth
	BootstrapUrl      string
	DeviceGatewayHost string // should be an IP
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
		DB: config.Database{
			Host:       "localhost",
			Password:   "Pass2020!",
			DbName:     "registry",
			Username:   "postgres",
			Port:       5432,
			SslEnabled: false,
		},
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
	}
}
