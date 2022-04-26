package pkg

import (
	"github.com/ukama/ukamaX/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Port              int
	Services          GrpcEndpoints `mapstructure:"services"`
	SwaggerAssets     string
}

type GrpcEndpoints struct {
	TimeoutSeconds int `mapstructure:"timeoutSeconds"`
	Hss            string
}

func NewConfig() *Config {
	return &Config{
		Port: 8080,
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},

		Services: GrpcEndpoints{
			TimeoutSeconds: 5,
			Hss:            "hss:9090",
		},
		SwaggerAssets: "swagger-ui",
	}
}
