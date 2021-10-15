package pkg

import (
	"github.com/ukama/ukamaX/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                config.Database
	Grpc              Grpc
	BootstrapAuth     Auth
	BootstrapUrl      string
	DeviceGatewayHost string // should be an IP
}

type Grpc struct {
	Port int
}

type Auth struct {
	ClientId     string
	ClientSecret string
	Audience     string
	GrantType    string
	Auth0Host    string
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
		Grpc: Grpc{
			Port: 9090,
		},
		BootstrapAuth: Auth{
			Audience:  "bootstrap.ukama.com",
			GrantType: "client_credentials",
		},
	}
}
