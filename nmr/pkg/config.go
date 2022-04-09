package pkg

import (
	cors "github.com/gin-contrib/cors"
	"github.com/ukama/openIoR/services/common/config"
	"github.com/ukama/openIoR/services/common/rest"
)

type Routes map[string]string

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           config.Metrics
	Server            rest.HttpConfig
	ApiIf             ServiceApiIf
	RouterService     string
	DB                config.Database
}

type Pattern struct {
	SRoutes []Routes
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

	return &Config{
		Server: rest.HttpConfig{
			Port: 8085,
			Cors: cors.Config{
				AllowOrigins: []string{"http://localhost", "https://localhost", "*"},
			},
		},

		RouterService: "http://localhost:8090",
		ApiIf: ServiceApiIf{
			Name: "nmr",
			P: Pattern{
				SRoutes: []Routes{
					{
						"node": "*", "looking_for": "node_info", "Path": "/node",
					},
					{
						"node": "*", "looking_to": "node_info", "Path": "/node",
					},
					{
						"node": "*", "looking_for": "node_status", "Path": "/node",
					},
					{
						"node": "*", "looking_to": "node_status", "Path": "/node",
					},
					{
						"node": "*", "looking_for": "node_list", "Path": "/node/all",
					},
					{
						"module": "*", "looking_for": "module_info", "Path": "/module",
					},
					{
						"module": "*", "looking_to": "module_info", "Path": "/module",
					},
					{
						"module": "*", "looking_for": "module_status", "Path": "/module",
					},
					{
						"module": "*", "looking_to": "module_status", "Path": "/module",
					},
					{
						"module": "*", "looking_for": "module_list", "Path": "/module/all",
					},
					{
						"module": "*", "looking_to": "module_data", "Path": "/module/data",
					},
					{
						"module": "*", "looking_for": "module_data", "Path": "/module/data",
					},
				},
			},
			F: Forward{
				Ip:   "http://localhost",
				Port: 8095,
			},
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
