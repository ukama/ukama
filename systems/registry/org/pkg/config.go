package pkg

import (
	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database `default:"{}"`
	Grpc              *config.Grpc     `default:"{}"`
	UsersHost         string           `default:"users:9090"`
	Metrics           *config.Metrics  `default:"{}"`
}
