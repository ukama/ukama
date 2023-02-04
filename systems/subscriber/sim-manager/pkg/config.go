package pkg

import (
	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig      `mapstructure:",squash"`
	DB                     *config.Database  `default:"{}"`
	Grpc                   *config.Grpc      `default:"{}"`
	MsgClient              *config.MsgClient `default:"{}"`
	Queue                  *config.Queue     `default:"{}"`
	Service                *config.Service
	Key                    string
	Metrics                *config.Metrics `default:"{}"`
	PackageHost            string          `default:"package:9090"`
	SubscriberRegistryHost string          `default:"subscriber-registry:9090"`
	SimPoolHost            string          `default:"sim-pool:9090"`
	TestAgentHost          string          `default:"test-agent:9090"`
}
