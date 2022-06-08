package pkg

import (
	"github.com/gin-contrib/cors"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Kratos            Kratos `mapstructure:"kratos"`
	Server            rest.HttpConfig
	BypassAuthMode    bool
	Services          GrpcEndpoints  `mapstructure:"services"`
	HttpServices      HttpEndpoints  `mapstructure:"httpServices"`
	Metrics           config.Metrics `mapstructure:"metrics"`
}

type Kratos struct {
	Url string
}

type GrpcEndpoints struct {
	TimeoutSeconds int `mapstructure:"timeoutSeconds"`
	Registry       string
	Users          string
}

type HttpEndpoints struct {
	TimeoutSeconds int `mapstructure:"timeoutSeconds"`
	NodeMetrics    string
}

func NewConfig() *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		Kratos: Kratos{
			"http://kratos",
		},
		BypassAuthMode: false,
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},

		Services: GrpcEndpoints{
			TimeoutSeconds: 5,
			Registry:       "network:9090",
			Users:          "users:9090",
		},
		HttpServices: HttpEndpoints{
			TimeoutSeconds: 5,
			NodeMetrics:    "http://localhost",
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		Metrics: *config.DefaultMetrics(),
	}
}
