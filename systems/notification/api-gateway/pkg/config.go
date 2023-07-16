package pkg

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Server            rest.HttpConfig
	Services          GrpcEndpoints `mapstructure:"services"`
	Auth              *config.Auth  `mapstructure:"auth"`
}

type GrpcEndpoints struct {
	Timeout time.Duration
	Mailer  string
}

func NewConfig() *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},

		Services: GrpcEndpoints{
			Timeout: 5 * time.Second,
			Mailer:  "mailer:9090",
		},

		Server: rest.HttpConfig{
			Port: 8085,
			Cors: defaultCors,
		},
		Auth: config.LoadAuthHostConfig("auth"),
	}
}
