package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database `default:"{}"`
	Grpc              *config.Grpc     `default:"{}"`
	Service           *config.Service  `default:"{}"`
	Queue            *config.Queue     `default:"{}"`
	MsgClient        *config.MsgClient `default:"{}"`
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
		MsgClient: &config.MsgClient{
			Timeout: 7 * time.Second,
		},
	}
}
