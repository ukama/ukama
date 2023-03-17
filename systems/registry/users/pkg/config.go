package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Service           *config.Service   `default:"{}"`
	MsgClient         *config.MsgClient `default:"{}"`
	Org               string            `default:"org:9090"`
	OrgOWnerName      string            `default:"Ukama Root"`
	OrgOWnerEmail     string            `default:"hello@ukama.com"`
	OrgOWnerPhone     string            `default:"0000000000"`
	OrgOWnerUUID      string
	Queue             *config.Queue `default:"{}"`
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
