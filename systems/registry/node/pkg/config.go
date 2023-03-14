package pkg

import (
	uconf "github.com/ukama/ukama/systems/common/config"
	pmetric "github.com/ukama/ukama/systems/common/pushgatewayMetrics"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
	Metrics          *uconf.Metrics  `default:"{}"`
	PushMetricHost   string          `default:"http://localhost:9091"`
	Org              string          `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
}

const (
	NumberOfNodes         = "number_of_nodes"
	NumberOfActiveNodes   = "active_node_count"
	NumberOfInactiveNodes = "inactive_node_count"
	GaugeType             = "gauge"
)

var NodeMetric = []pmetric.MetricConfig{{
	Name:   NumberOfNodes,
	Type:   GaugeType,
	Labels: map[string]string{"network": "", "org": ""},
	Value:  0,
},
	{
		Name:   NumberOfActiveNodes,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
	{
		Name:   NumberOfInactiveNodes,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
}
