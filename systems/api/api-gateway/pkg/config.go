package pkg

import (
	"time"

	"github.com/gin-contrib/cors"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database `default:"{}"`
	Server            rest.HttpConfig
	HttpServices      HttpEndpoints  `mapstructure:"httpServices"`
	Metrics           config.Metrics `mapstructure:"metrics"`
	Auth              *config.Auth   `mapstructure:"auth"`
}

type HttpEndpoints struct {
	Timeout     time.Duration
	Network     string
	Package     string
	Subscriber  string
	Sim         string
	Node        string
	NodeMetrics string
}

func NewConfig(name string) *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},

		DB: &config.Database{
			DbName: name,
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},

		HttpServices: HttpEndpoints{
			Timeout:     3 * time.Second,
			Network:     "http://localhost",
			Package:     "http://localhost",
			Subscriber:  "http://localhost",
			Sim:         "http://localhost",
			Node:        "http://localhost",
			NodeMetrics: "http://localhost",
		},

		Metrics: *config.DefaultMetrics(),
		Auth:    config.LoadAuthHostConfig("auth"),
	}
}
