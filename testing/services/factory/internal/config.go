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
	config.BaseConfig     `mapstructure:",squash"`
	Metrics               config.Metrics
	Server                rest.HttpConfig
	ApiIf                 config.ServiceApiIf
	ServiceRouter         string
	GitUser               string
	GitPass               string
	Docker                Docker
	BuilderRegCred        string `default:"dregcred"`
	BuilderImage          string
	BuilderCmd            []string
	RabbitUri             string
	VNodeRepoServerUrl    string
	VNodeRepoName         string `default:"virtualnode"`
	Namespace             string
	BackOffLimit          int32 `default:"4"`
	TimeToLive            int32 `default:"60"`
	ActiveDeadLineSeconds int64 `default:"3600"`
	SecRef                string
	CmRef                 string
}

var ServiceConfig *Config

/* Info/List--> Get
   Update --> PUT
   Add --> POST
   delete -> delete
*/
func NewConfig() *Config {

	return &Config{
		Server: rest.DefaultHTTPConfig(),

		ServiceRouter: "http://localhost:8081",
		ApiIf: config.ServiceApiIf{
			Name: ServiceName,
			P: []config.Route{
				{
					"ping": ServiceName, "path": "/ping",
				},
				{
					"node": "*", "looking_for": "fact_node_info", "path": "/node",
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
