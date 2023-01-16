package pkg

import (
	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig      `mapstructure:",squash"`
	DB                     *config.Database `default:"{}"`
	Grpc                   *config.Grpc     `default:"{}"`
	PackageHost            string           `default:"package:9090"`
	SubscriberRegistryHost string           `default:"subscriber-registry:9090"`
	SimPoolHost            string           `default:"sim-pool:9090"`
	TestAgentHost          string           `default:"test-agent:9090"`
	Key                    string
	Metrics                *config.Metrics `default:"{}"`
}
