package pkg

import (
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Grpc              *config.Grpc    `default:"{}"`
	Metrics           *config.Metrics `default:"{}"`
	Service           *config.Service
	Timeout           time.Duration `default:"3s"`
}

func NewConfig(name string) *Config {
	// Sanitize name
	name = strings.ReplaceAll(name, "-", "_")

	return &Config{
		Service: config.LoadServiceHostConfig(name),
	}
}
