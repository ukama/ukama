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

type Routes map[string]string

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           config.Metrics
	Server            rest.HttpConfig
	ApiIf             ServiceApiIf
	RouterService     string
}

type Pattern struct {
	MustRoutes     Routes `json:"all"`
	OptionalRoutes Routes `json:"any"`
}

type Forward struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type ServiceApiIf struct {
	Name string  `json:"name"`
	P    Pattern `json:"pattern"`
	F    Forward `json:"forward"`
}

func NewConfig() *Config {

	optRoute := make(Routes)
	optRoute["oabc"] = "oxyz"
	optRoute["o123"] = "o789"

	mustRoute := make(Routes)
	mustRoute["mabc"] = "mxyz"
	mustRoute["m123"] = "m789"

	return &Config{
		Server: rest.HttpConfig{
			Port: 8085,
			Cors: cors.Config{
				AllowOrigins: []string{"http://localhost", "https://localhost", "*"},
			},
		},

		RouterService: "http://localhost:8090",
		ApiIf: ServiceApiIf{
			Name: "bootsrap",
			P: Pattern{
				MustRoutes:     nil,
				OptionalRoutes: nil,
			},
			F: Forward{
				Ip:   "http://localhost",
				Port: 8095,
			},
		},
	}
}
