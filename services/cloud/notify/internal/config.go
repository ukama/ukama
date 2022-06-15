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
	NodeCmd           []string
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
					"looking_to": "post_notification", "path": "/",
				},
				{
					"looking_to": "delete_notification", "path": "/",
				},
				{
					"looking_for": "list_notification", "path": "/",
				},
				{
					"node": "*", "looking_to": "notification", "type": "*", "path": "/",
				},
				{
					"node": "*", "looking_to": "delete_notification", "type": "*", "path": "/",
				},
				{
					"node": "*", "looking_for": "list_notification", "type": "*", "path": "/",
				},
				{
					"service": "*", "looking_to": "notification", "type": "*", "path": "/",
				},
				{
					"service": "*", "looking_to": "delete_notification", "type": "*", "path": "/",
				},
				{
					"service": "*", "looking_for": "list_notification", "type": "*", "path": "/",
				},
			},
			F: config.DefaultForwardConfig(),
		},
		DB: config.DefaultDatabaseName(ServiceName),
	}
}
