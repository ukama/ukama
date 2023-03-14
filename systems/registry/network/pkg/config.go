package pkg

import (
	"github.com/ukama/ukama/systems/common/config"
	pmetric "github.com/ukama/ukama/systems/common/pushgatewayMetrics"
)


const (
	NumberOfNetwork = "number_of_network"
	GaugeType           = "gauge"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database `default:"{}"`
	Grpc              *config.Grpc     `default:"{}"`
	OrgHost           string           `default:"org:9090"`
	Metrics           *config.Metrics  `default:"{}"`
	Org               string `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
	PushMetricHost    string `default:"http://localhost:9091"`
}
var NetworkMetric = []pmetric.MetricConfig{{
	Name:   NumberOfNetwork,
	Type:   GaugeType,
	Labels: map[string]string{"network": "", "org": ""},
	Value:  0,
},
}