package pkg

import (
	"github.com/gin-contrib/cors"
	"github.com/ukama/ukamaX/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Kratos            Kratos `mapstructure:"kratos"`
	Port              int
	BypassAuthMode    bool
	Cors              cors.Config
	Services          GrpcEndpoints  `mapstructure:"services"`
	Metrics           config.Metrics `mapstructure:"metrics"`
}

type Kratos struct {
	Url string
}

type GrpcEndpoints struct {
	TimeoutSeconds int `mapstructure:"timeoutSeconds"`
	Registry       string
	Hss            string
}

func NewConfig() *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true

	return &Config{
		Kratos: Kratos{
			"http://kratos",
		},
		Port:           8080,
		BypassAuthMode: false,
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},

		Services: GrpcEndpoints{
			TimeoutSeconds: 5,
			Registry:       "registry:9090",
			Hss:            "hss:9090",
		},
		Cors:    defaultCors,
		Metrics: *config.DefaultMetrics(),
	}
}
