package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	Queue             *config.Queue     `default:"{}"`
	Service           *config.Service   `default:"{}"`
	SimTokenKey       string            `default:"11111111111111111111111111111111"`
	AsrHost           string            `default:"localhost"`
	NetworkHost       string            `default:"http://localhost:8085"`
	PCRFHost          string            `default:"http://localhost:8085"`
	FactoryHost       string            `default:"http://localhost:8085"`
	Org               string            `default:"880f7c63-eb57-461a-b514-248ce91e9b3e"`
}

type SimManager struct {
	Host string
	Name string
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: ServiceName,
		},

		Grpc: &config.Grpc{
			Port: 9090,
		},
		SimTokenKey: "11111111111111111111111111111111",
		Service:     config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout:        5 * time.Second,
			ListenerRoutes: []string{"event.cloud.lookup.organization.create"},
		},
	}
}
