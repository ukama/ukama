package pkg

import (
	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database `default:"{}"`
	Grpc              *config.Grpc     `default:"{}"`
	PackageHost       string           `default:"package:9090"`
	TestAgentHost     string           `default:"test-agent:9090"`
	Metrics           *config.Metrics  `default:"{}"`
}
