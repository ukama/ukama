package internal

import (
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/rest"
)

type Docker struct {
	User string
	Pass string
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           config.Metrics
	Server            rest.HttpConfig
	ApiIf             config.ServiceApiIf
	DB                config.Database
	ServiceRouter     string
	GitUser           string
	GitPass           string
	Docker            Docker
	NodeImage         string
	BuilderImage      string
	BuilderCmd        []string
	Kubeconfig        string
	Queue             config.Queue
	RepoServerUrl     string
	Namespace         string
}

var ServiceConfig *Config

/* Info/List--> Get
   Update --> PUT
   Add --> POST
   delete -> delete
*/
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
					"node": "*", "looking_for": "vnode_info", "path": "/node",
				},
				{
					"node": "*", "looking_to": "vnode_power_on", "path": "/node/poweron",
				},
				{
					"node": "*", "looking_to": "vnode_power_off", "path": "/node/poweroff",
				},
				{
					"looking_for": "vnode_list", "path": "/list",
				},
			},
			F: config.DefaultForwardConfig(),
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
