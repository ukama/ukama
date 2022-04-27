package internal

import (
	"github.com/gin-contrib/cors"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Server            rest.HttpConfig
	Metrics           config.Metrics
	ApiIf             config.ServiceApiIf
	ServiceRouter     string
	DB                config.Database
}

var ServiceConf *Config

// NewConfig creates new config with default values. Those values will be overridden by Viper
func NewConfig() *Config {
	return &Config{
		Server: rest.HttpConfig{
			Port: 8087,
			Cors: cors.Config{
				AllowOrigins: []string{"http://localhost", "https://localhost", "*"},
			},
		},
		ServiceRouter: "http://localhost:8091",
		ApiIf: config.ServiceApiIf{
			Name: ServiceName,
			P: []config.Route{
				{
					"ping": ServiceName, "path": "/ping",
				},
				{
					"node": "UK-SA2154-HNODE-A1-0001", "looking_for": "org_credentials", "path": "/orgs/node",
				},
				{
					"node": "UK-SA2154-HNODE-A1-0001", "looking_to": "add_node", "org": "fundme", "path": "/orgs/node",
				},
				{
					"org": "UK-SA2154-HNODE-A1-0001", "looking_for": "add_org", "path": "/orgs/",
				},
			},

			F: config.Forward{
				Ip:   "localhost",
				Port: 8087,
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
