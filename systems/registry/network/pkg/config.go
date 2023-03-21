package pkg

import (
	"time"

	pmetric "github.com/ukama/ukama/systems/common/metrics"

	uconf "github.com/ukama/ukama/systems/common/config"
)
const (
	NumberOfNetwork = "number_of_network"
	GaugeType           = "gauge"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	OrgHost          string           `default:"org:9090"`
	Service          *uconf.Service
	Org               string `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
	PushGatewayHost    string `default:"http://localhost:9091"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
		},
	}
}
var NetworkMetric = []pmetric.MetricConfig{{
	Name:   NumberOfNetwork,
	Type:   GaugeType,
	Labels: map[string]string{"network": "", "org": ""},
	Value:  0,
},
}