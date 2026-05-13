package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database `mapstructure:"db"`
	Grpc              *config.Grpc     `mapstructure:"grpc"`
	Services          GrpcEndpoints    `mapstructure:"services"`
}

type GrpcEndpoints struct {
	Timeout    time.Duration `mapstructure:"timeout"`
	Controller string        `mapstructure:"controller"`
}

func NewConfig(name string) *Config {
	db := config.DefaultDatabaseName(name)
	return &Config{
		BaseConfig: config.BaseConfig{DebugMode: false},
		DB:         &db,
		Grpc:       &config.Grpc{},
		Services:   GrpcEndpoints{Timeout: 3 * time.Second, Controller: "controller:9090"},
	}
}
