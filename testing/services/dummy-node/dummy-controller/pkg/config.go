package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/internal/backhaul"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/internal/battery"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/internal/solar"
)

type Config struct {
    DB               *uconf.Database `default:"{}"`
    uconf.BaseConfig `mapstructure:",squash"`
    Grpc             *uconf.Grpc    `default:"{}"`
    Queue            *uconf.Queue    `default:"{}"`
    Timeout          time.Duration   `default:"20s"`
    Service          *uconf.Service
    OrgName          string         `default:"ukama"`
    SolarMetrics     *solar.SolarProvider     `default:"{}"`
    BatteryMetrics   *battery.BatteryProvider `default:"{}"`
    NetworkMetrics   *backhaul.BackhaulProvider `default:"{}"`
}

func NewConfig(name string) *Config {
	return &Config{
		Service: uconf.LoadServiceHostConfig(name),
	
	}
}
