package pkg

import (
	"github.com/ukama/ukamaX/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Grpc              config.Grpc
	EtcdHost          string
	EtcdEnabled       bool
}

func NewConfig() *Config {
	return &Config{
		Grpc: config.Grpc{
			Port: 9090,
		},
		EtcdHost:    "localhost:2379",
		EtcdEnabled: true,
	}
}
