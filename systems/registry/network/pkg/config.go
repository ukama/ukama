package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
	metric "github.com/ukama/ukama/systems/common/metrics"
)

const (
	NumberOfNetwork = "number_of_network"
	NumberOfSites   = "number_of_sites"
	GaugeType       = "gauge"
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
	PushGatewayHost  string `default:"http://localhost:9091"`
}

var NetworkMetric = []metric.MetricConfig{
	{
		Name:   NumberOfNetwork,
		Type:   GaugeType,
		Labels: map[string]string{"org": ""},
		Value:  0,
	},
	{
		Name:   NumberOfSites,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
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
