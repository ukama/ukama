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
	ServiceRouter     string
	GitUser           string
	GitPass           string
	Docker            Docker
	VmImage           string
	BuilderImage      string
	BuilderCmd        []string
	RabbitUri         string
	RepoServerUrl     string
	Namespace         string
	SecRef            string
	CmRef             string
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

		ServiceRouter: "http://localhost:8081",
		ApiIf: config.ServiceApiIf{
			Name: ServiceName,
			P: []config.Route{
				{
					"ping": ServiceName, "path": "/ping",
				},
				{
					"node": "*", "looking_for": "info", "path": "/node",
				},
				{
					"looking_to": "create_node", "type": "*", "count": "*", "path": "/node",
				},
				{
					"node": "*", "looking_to": "delete", "path": "/node",
				},
			},
			F: config.Forward{
				Ip:   "localhost",
				Port: 8080,
				Path: "/",
			},
		},
	}
}
