package pkg

import (
	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Grpc              *config.Grpc    `default:"{}"`
	Metrics           *config.Metrics `default:"{}"`
}
