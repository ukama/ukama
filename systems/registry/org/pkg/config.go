package pkg

import (
	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database `default:"{}"`
	Grpc              *config.Grpc     `default:"{}"`
	Metrics           *config.Metrics  `default:"{}"`
	Service           *config.Service  `default:"{}"`
	Users             string           `default:"users:9090"`
	OrgName           string           `default:"ukama"`
	OrgOwnerUUID      string
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: name,
		},

		Service: config.LoadServiceHostConfig(name),

		Metrics: &config.Metrics{
			Port: 10251,
		},
	}
}
