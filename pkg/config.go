package pkg

import (
	cors "github.com/gin-contrib/cors"
	"github.com/ukama/openIoR/services/common/config"
	"github.com/ukama/openIoR/services/common/rest"
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
			Port: 8085,
			Cors: cors.Config{
				AllowOrigins: []string{"http://localhost", "https://localhost", "*"},
			},
		},

		ServiceRouter: "http://localhost:8090",
		ApiIf: config.ServiceApiIf{
			Name: "bootsrap",
			P: config.Pattern{
				Routes: []config.Route{
					{
						"node": "*", "looking_for": "validation", "Path": "/nodes",
					},
				},
			},
			F: config.Forward{
				Ip:   "http://localhost",
				Port: 8095,
			},
		},
	}
}
