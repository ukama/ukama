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
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: name,
		},

		Service: config.LoadServiceHostConfig(name),

		MsgClient: &config.MsgClient{
			Host:           "msg-client-billing-exporter:9095",
			Timeout:        5 * time.Second,
			ListenerRoutes: []string{"event.cloud.cdr.sim.usage"},
		},
	}
}
