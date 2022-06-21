package pkg

import (
	bootstrap "github.com/ukama/ukama/services/bootstrap/client"
	"github.com/ukama/ukama/services/common/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database      `default:"{}"`
	Grpc              *config.Grpc          `default:"{}"`
	BootstrapAuth     *bootstrap.AuthConfig `default:"{}"`
	BootstrapUrl      string
	// this host will be sent to bootstrap service
	DeviceGatewayHost string        // should be an IP
	Queue             *config.Queue `default:"{}"`
	// debugMode shoule be enabled to allow bypyssing bootstrap
	Debug   *bootstrap.DebugConf `default:"{}"`
	Metrics *config.Metrics      `default:"{}"`
}
