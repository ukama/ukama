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
	Queue             config.Queue
}

var ServiceConfig *Config

func NewConfig() *Config {

	return &Config{
		Server: config.DefaultHTTPConfig(),

		ServiceRouter: "http://localhost:8091",

		Queue: config.Queue{
			Uri: "amqp://guest:guest@localhost:5672",
		},

		ApiIf: config.ServiceApiIf{
			Name: ServiceName,
			P: []config.Route{
				{
					"ping": ServiceName, "path": "/ping",
				},
				{
					"looking_to": "post_notification", "path": "/notification",
				},
				{
					"looking_to": "delete_notification", "path": "/notification",
				},
				{
					"looking_for": "list_notification", "path": "/notification/list",
				},
				{
					"node": "*", "looking_to": "notification", "type": "*", "path": "/notification/node",
				},
				{
					"node": "*", "looking_to": "delete_notification", "type": "*", "path": "/notification/node",
				},
				{
					"node": "*", "looking_for": "list_notification", "count": "*", "path": "/notification/node/list",
				},
				{
					"service": "*", "looking_to": "notification", "type": "*", "path": "/notification/service",
				},
				{
					"service": "*", "looking_to": "delete_notification", "type": "*", "path": "/notification/service",
				},
				{
					"service": "*", "looking_for": "list_notification", "count": "*", "path": "/notification/service/list",
				},
			},
			F: config.DefaultForwardConfig(),
		},
		DB: config.DefaultDatabaseName(ServiceName),
	}
}
