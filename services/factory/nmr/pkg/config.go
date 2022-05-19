package pkg

import (
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           config.Metrics
	Server            rest.HttpConfig
	ApiIf             config.ServiceApiIf
	ServiceRouter     string
	DB                config.Database
}

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
					"node": "*", "looking_for": "info", "path": "/node",
				},
				{
					"node": "*", "looking_to": "update", "path": "/node",
				},
				{
					"node": "*", "looking_to": "delete", "path": "/node",
				},
				{
					"node": "*", "looking_for": "status_info", "path": "/node/status",
				},
				{
					"node": "*", "looking_to": "status_update", "status": "*", "path": "/node/status",
				},
				{
					"node": "*", "looking_for": "mfg_status_info", "path": "/node/mfg_status",
				},
				{
					"node": "*", "looking_to": "mfg_status_update", "mfg_status": "*", "path": "/node/mfg_status",
				},
				{
					"node": "*", "looking_for": "list", "path": "/node/all",
				},
				{
					"module": "*", "looking_for": "info", "path": "/module",
				},
				{
					"module": "*", "looking_to": "update", "path": "/module",
				},
				{
					"module": "*", "looking_to": "delete", "path": "/module",
				},
				{
					"module": "*", "looking_to": "allocate", "node": "*", "path": "/module/",
				},
				{
					"module": "*", "looking_for": "mfg_status_info", "status": "*", "path": "/module/status",
				},
				{
					"module": "*", "looking_to": "mfg_status_update", "status": "*", "path": "/module/status",
				},
				{
					"module": "*", "looking_for": "field_info", "field": "*", "path": "/module/field",
				},
				{
					"module": "*", "looking_to": "field_update", "field": "*", "path": "/module/field",
				},
				{
					"module": "*", "looking_for": "mfg_info", "data": "*", "path": "/module/data",
				},
				{
					"module": "*", "looking_for": "list", "path": "/module/all",
				},
				{
					"module": "*", "looking_to": "bootstrap_cert_delete", "path": "/module/bootstrapcerts",
				},
			},
			F: config.Forward{
				Ip:   "localhost",
				Port: 8085,
				Path: "/",
			},
		},
		DB: config.DefaultDatabaseName(ServiceName),
	}
}
