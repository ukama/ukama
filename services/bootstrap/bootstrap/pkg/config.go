package pkg

import (

	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/rest"
)

// type Route {
// 	key string
// 	value interface{}
// }

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           config.Metrics
	Server            rest.HttpConfig
	ApiIf             config.ServiceApiIf
	ServiceRouter     string
}

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
					"node": "*", "looking_for": "validation", "path": "/",
				},
			},
			F: config.Forward{
				Ip:   "localhost",
				Port: 8086,
				Path: "/",
			},
		},
	}
}
