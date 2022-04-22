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
		ServiceRouter: "http://192.168.0.14:8091",
		ApiIf: config.ServiceApiIf{
			Name: "lookup",
			P: config.Pattern{
				Routes: []config.Route{
					{
						"node": "*", "looking_for": "node", "org": "*", "Path": "/orgs/node",
					},
					{
						"node": "*", "looking_to": "add_node", "org": "*", "Path": "/orgs/node",
					},
					{
						"org": "*", "looking_for": "add_org", "Path": "/orgs/",
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
