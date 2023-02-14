package pkg

import (
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig      `mapstructure:",squash"`
	DB                     *config.Database  `default:"{}"`
	Grpc                   *config.Grpc      `default:"{}"`
	Queue                  *config.Queue     `default:"{}"`
	Metrics                *config.Metrics   `default:"{}"`
	Timeout                time.Duration     `default:"3s"`
	MsgClient              *config.MsgClient `default:"{}"`
	Key                    string
	Service                *config.Service
	PackageHost            string `default:"package:9090"`
	SubscriberRegistryHost string `default:"subscriber-registry:9090"`
	SimPoolHost            string `default:"sim-pool:9090"`
	TestAgentHost          string `default:"test-agent:9090"`
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
