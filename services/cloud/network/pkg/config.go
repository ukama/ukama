package pkg

import (
	bootstrap "github.com/ukama/ukama/services/bootstrap/client"
	"github.com/ukama/ukama/services/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                config.Database
	Grpc              config.Grpc
	BootstrapAuth     bootstrap.AuthConfig
	BootstrapUrl      string
	Queue             config.Queue
	Debug             bootstrap.DebugConf
}
