package pkg

import (
	"github.com/gin-contrib/cors"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Server            rest.HttpConfig
	Service           *config.Service
	R                 *rest.RestClient
	Mailer            *config.Mailer
}

func NewConfig(name string) *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},
		Server: rest.HttpConfig{
			Port: 8085,
			Cors: defaultCors,
		},
		Service: config.LoadServiceHostConfig(name),
		Mailer:  config.LoadMailerHostConfig(name),
	}
}
