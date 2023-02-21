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
	Service           *config.Service
	Key               string
	DataPlan          string `default:"data-plan:8080"`
	SubsRegistry      string `default:"registry:9091"`
	SimPool           string `default:"sim-pool:9090"`
	TestAgent         string `default:"test-agent:9093"`
	OperatorAgent     string `default:"operator-agent:8080"`
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
	}
}
