package pkg

import (
	"github.com/ukama/ukama/services/common/config"
	"time"
)

type ServerConfig struct {
	Port int `default:"10251"`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           *config.Metrics `default:"{}"`
	Server            *ServerConfig   `default:"{}"`
	Registry          *RegistryConf   `default:"{}"`
}

type RegistryConf struct {
	config.GrpcService `mapstructure:",squash"`
	PollInterval       time.Duration `default:"1m"`
}
