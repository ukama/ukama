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

/* Info/List--> Get
   Update --> PUT
   Add --> POST
   delete -> delete
*/
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
						"node": "*", "looking_for": "info", "Path": "/node/",
					},
					{
						"node": "*", "looking_to": "update", "Path": "/node/",
					},
					{
						"node": "*", "looking_to": "delete", "Path": "/node/",
					},
					{
						"node": "*", "looking_for": "info", "status": "*", "Path": "/node/status",
					},
					{
						"node": "*", "looking_to": "status", "status": "*", "Path": "/node/status",
					},
					{
						"node": "*", "looking_for": "info", "mfg_status": "*", "Path": "/node/mfg_status",
					},
					{
						"node": "*", "looking_to": "update", "mfg_status": "*", "Path": "/node/mfg_status",
					},
					{
						"node": "*", "looking_for": "list", "Path": "/node/all",
					},
					{
						"module": "*", "looking_for": "info", "Path": "/module/",
					},
					{
						"module": "*", "looking_to": "update", "Path": "/module/",
					},
					{
						"module": "*", "looking_to": "delete", "Path": "/module/",
					},
					{
						"module": "*", "looking_for": "info", "status": "*", "Path": "/module/status",
					},
					{
						"module": "*", "looking_to": "status", "status": "*", "Path": "/module/status",
					},
					{
						"module": "*", "looking_for": "info", "field": "*", "Path": "/module/field",
					},
					{
						"module": "*", "looking_to": "update", "field": "*", "Path": "/module/field",
					},
					{
						"module": "*", "looking_for": "info", "data": "*", "Path": "/module/data",
					},
					{
						"module": "*", "looking_for": "list", "Path": "/module/all",
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
