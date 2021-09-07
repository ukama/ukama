package pkg

import (
	"github.com/ukama/ukamaX/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Kratos            Kratos `mapstructure:"kratos"`
	Port              int
	BypassAuthMode    bool
	Services          GrpcEndpoints `mapstructure:"services"`
}

type Kratos struct {
	Url string
}

type GrpcEndpoints struct {
	TimeoutSeconds int `mapstructure:"timeoutSeconds"`
	Registry       string
}

func NewConfig() *Config {
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
		},
	}
}
