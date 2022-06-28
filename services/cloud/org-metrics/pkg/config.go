package pkg

import (
	"time"

	"github.com/ukama/ukama/services/common/config"
)

type ServerConfig struct {
	Port int `default:"10251"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           *config.Metrics `default:"{}"`
	Server            *ServerConfig   `default:"{}"`
	Network           *NetworkConf    `default:"{}"`
}

type NetworkConf struct {
	config.GrpcService `mapstructure:",squash"`
	PollInterval       time.Duration `default:"1m"`
}
