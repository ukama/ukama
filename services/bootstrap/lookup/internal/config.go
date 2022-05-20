package internal

import (
	
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

		Server: config.DefaultHTTPConfig(),
		ServiceRouter: "http://localhost:8091",
		ApiIf: config.ServiceApiIf{
			Name: ServiceName,
			P: []config.Route{
				{
					"ping": ServiceName, "path": "/ping",
				},
				{
					"node": "*", "looking_for": "org_credentials", "path": "/orgs/node",
				},
				{
					"node": "*", "looking_to": "add_node", "org": "*", "path": "/orgs/node",
				},
				{
					"org": "*", "looking_to": "add_org", "path": "/orgs",
				},
			},

			F: config.Forward{
				Ip:   "localhost",
				Port: 8080,
			},
		},

		DB: config.DefaultDatabaseName(ServiceName),
	}
}
