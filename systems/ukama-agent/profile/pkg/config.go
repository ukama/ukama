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
	NetworkHost       string            `default:"http://localhost:8085"`
	PCRFHost          string            `default:"http://localhost:8085"`
	FactoryHost       string            `default:"http://localhost:8085"`
	Org               string            `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
	IsMsgBus          bool              `default:"false"`
	NodePolicyPath    string            `default:"/v1/epc/pcrf/subscriber"`
	PolicyCheckPeriod time.Duration     `default:"10s"`
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
	}
}
