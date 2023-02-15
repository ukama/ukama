package pkg

import (
	"strings"
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
	Key               string
	Service           *config.Service
	DataPlan          *config.Service
	SubsRegistry      *config.GrpcService
	SimPool           *config.GrpcService
	TestAgent         *config.GrpcService
	OperatorAgent     *config.Service
}

func NewConfig(name string) *Config {
	// Sanitize name
	name = strings.ReplaceAll(name, "-", "_")

	return &Config{
		DB: &config.Database{
			DbName: name,
		},

		Service: config.LoadServiceHostConfig(name),

		MsgClient: &config.MsgClient{
			Timeout: 5 * time.Second,
		},

		DataPlan: &config.Service{
			Host: `default:"data-plan"`,
			Port: `default:"8080"`,
			Uri:  `default:"data-plan:8080"`,
		},

		SubsRegistry: &config.GrpcService{
			Timeout: 2 * time.Second,
			Host:    `default:"subscriber-registry:9090"`,
		},

		SimPool: &config.GrpcService{
			Timeout: 2 * time.Second,
			Host:    `default:"sim-pool:9090"`,
		},

		TestAgent: &config.GrpcService{
			Timeout: 2 * time.Second,
			Host:    `default:"test-agent:9090"`,
		},

		OperatorAgent: &config.Service{
			Host: `default:"operator-agent"`,
			Port: `default:"8080"`,
			Uri:  `default:"operator-agent:8080"`,
		},
	}
}
