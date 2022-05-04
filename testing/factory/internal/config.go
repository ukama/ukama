package internal

import (
	cors "github.com/gin-contrib/cors"
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
	Kubeconfig        string
	RabbitUri         string
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
		Server: rest.HttpConfig{
			Port: 8086,
			Cors: cors.Config{
				AllowOrigins: []string{"http://localhost", "https://localhost", "*"},
			},
		},

		//ServiceRouter: "http://192.168.0.14:8091",
		ServiceRouter: "http://localhost:8091",
		ApiIf: config.ServiceApiIf{
			Name: ServiceName,
			P: []config.Route{
				{
					"ping": ServiceName, "path": "/ping",
				},
				{
					"node": "*", "looking_for": "info", "path": "/node/",
				},
				{
					"looking_to": "create_node", "type": "*", "count": "*", "path": "/node/",
				},
				{
					"node": "*", "looking_to": "delete", "path": "/node/",
				},
			},
			F: config.Forward{
				//Ip:   "192.168.0.27",
				Ip:   "localhost",
				Port: 8086,
				Path: "/",
			},
		},
	}
}
