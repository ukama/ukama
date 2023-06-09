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
	From              string `default:"hello@dev.ukama.com" validation:"required"`
	Host              string `default:"localhost" validation:"required"`
	Port              int    `default:"25" validation:"required"`
	Password          string
	Username          string
}
type SmtpConfig struct {
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
			Port: 8080,
			Cors: defaultCors,
		},
		Service: config.LoadServiceHostConfig(name),
	}
}
