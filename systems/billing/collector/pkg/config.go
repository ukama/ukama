package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Queue             *config.Queue     `default:"{}"`
	Metrics           *config.Metrics   `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	Service           *config.Service
	System            string `default:"billing"`
	LagoHost          string `default:"localhost"`
	LagoPort          uint   `default:"3000"`
	LagoAPIKey        string
	OrgName           string
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: name,
		},

		Service: config.LoadServiceHostConfig(name),

		MsgClient: &config.MsgClient{
			Host:    "msg-client-billing:9095",
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.local.{{ .Org}}.operator.cdr.sim.usage",
				"event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.create",
				"event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.update",
				"event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.delete",
				"event.cloud.local.{{ .Org}}.subscriber.simmanager.package.activate",
			},
		},
	}
}
