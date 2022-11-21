package pkg

import (
	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database `default:"{}"`
	Grpc              *config.Grpc     `default:"{}"`
	OrgHost           string           `default:"org:9090"`
	Metrics           *config.Metrics  `default:"{}"`
}
