package internal

import "github.com/ukama/ukamaX/common/config"

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                config.Database
	Grpc              Grpc
}

type Grpc struct {
	Port int
}

func NewConfig() *Config {
	return &Config{
		DB: config.Database{
			Host:       "localhost",
			Password:   "Pass2020!",
			DbName:     "registry",
			Username:   "postgres",
			Port:       5432,
			SslEnabled: false,
		},
		Grpc: Grpc{
			Port: 9090,
		},
	}
}
