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
	AsrHost           string            `default:"localhost"`
	DataplanHost      string            `default:"http://localhost:8085"`
	NetworkHost       string            `default:"http://localhost:8085"`
	FactoryHost       string            `default:"http://localhost:8085"`
	OrgName           string            `default:"ukama"`
	OrgId             string            `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
	IsMsgBus          bool              `default:"true"`
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

		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout:        5 * time.Second,
			ListenerRoutes: []string{},
		},
		DataplanHost: "http://192.168.0.14:8085",
		NetworkHost:  "http://192.168.0.14:8085",
		FactoryHost:  "http://192.168.0.14:8085",
	}
}
