package internal

import "github.com/ukama/openIoR/services/common/config"

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	ApiIf             config.ServiceApiIf
	RouterService     string
	DB                config.Database
}

var ServiceConf *Config

// NewConfig creates new config with default values. Those values will be overridden by Viper
func NewConfig() *Config {
	return &Config{
		ApiIf: config.ServiceApiIf{
			Name: "nmr",
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
						"node": "*", "looking_for": "info", "mfg_status": "*", "Path": "/node/mfg_status",
					},
					{
						"node": "*", "looking_to": "update", "mfg_status": "*", "Path": "/node/mfg_status",
					},
					{
						"node": "*", "looking_for": "list", "Path": "/node/all",
					},
					{
						"module": "*", "looking_for": "info", "Path": "/module/",
					},
					{
						"module": "*", "looking_to": "update", "Path": "/module/",
					},
					{
						"module": "*", "looking_to": "delete", "Path": "/module/",
					},
					{
						"module": "*", "looking_for": "info", "status": "*", "Path": "/module/status",
					},
					{
						"module": "*", "looking_to": "status", "status": "*", "Path": "/module/status",
					},
					{
						"module": "*", "looking_for": "info", "field": "*", "Path": "/module/field",
					},
					{
						"module": "*", "looking_to": "update", "field": "*", "Path": "/module/field",
					},
					{
						"module": "*", "looking_for": "info", "data": "*", "Path": "/module/data",
					},
					{
						"module": "*", "looking_for": "list", "Path": "/module/all",
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
