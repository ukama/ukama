package pkg

import (
	cors "github.com/gin-contrib/cors"
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
		Server: rest.HttpConfig{
			Port: 8086,
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
