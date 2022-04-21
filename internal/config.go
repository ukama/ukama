package internal

import (
	"github.com/ukama/openIoR/services/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	ApiIf             config.ServiceApiIf
	ServiceRouter     string
	DB                config.Database
}

var ServiceConf *Config

// NewConfig creates new config with default values. Those values will be overridden by Viper
func NewConfig() *Config {
	return &Config{
		ApiIf: config.ServiceApiIf{
			Name: "lookup",
			P: config.Pattern{
				Routes: []config.Route{
					{
						"node": "*", "looking_for": "info", "Path": "/node/",
					},
					{
						"node": "*", "looking_to": "update", "Path": "/node/",
					},
					{
						"node": "*", "looking_to": "delete", "Path": "/node/",
					},
					{
						"node": "*", "looking_for": "info", "status": "*", "Path": "/node/status",
					},
					{
						"node": "*", "looking_to": "status", "status": "*", "Path": "/node/status",
					},
					{
						"node": "*", "looking_for": "list", "Path": "/node/all",
					},
				},
			},
			F: config.Forward{
				Ip:   "http://localhost",
				Port: 8095,
			},
		},

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
