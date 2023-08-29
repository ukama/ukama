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
	Services          GrpcEndpoints  `mapstructure:"services"`
	HttpServices      HttpEndpoints  `mapstructure:"httpServices"`
	Metrics           config.Metrics `mapstructure:"metrics"`
	Auth              *config.Auth   `mapstructure:"auth"`
}

type HttpEndpoints struct {
	Timeout     time.Duration
	NodeMetrics string
}

type GrpcEndpoints struct {
	Timeout time.Duration
	Mailer  string
	Notify  string
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
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},

		HttpServices: HttpEndpoints{
			Timeout:     3 * time.Second,
			NodeMetrics: "http://localhost",
		},

		Metrics: *config.DefaultMetrics(),

		Auth: config.LoadAuthHostConfig("auth"),
	}
}
