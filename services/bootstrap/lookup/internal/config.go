package internal

import "github.com/ukama/ukama/services/common/config"

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                config.Database
}

var ServiceConf *Config

// NewConfig creates new config with default values. Those values will be overridden by Viper
func NewConfig() *Config {
	return &Config{
		DB: config.Database{
			Host:       "localhost",
			Password:   "Pass2020!",
			Username:   "postgres",
			DbName:     ServiceName,
			SslEnabled: false,
			Port:       5432,
		},
	}
}
