package pkg

import (
	uconf "github.com/ukama/ukama/systems/common/config"
	metric "github.com/ukama/ukama/systems/common/metrics"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database `default:"{}"`
	Grpc             *uconf.Grpc     `default:"{}"`
	Metrics          *uconf.Metrics  `default:"{}"`
	PushGateway      string          `default:"http://localhost:9091"`
}

const (
	NumberOfNodes         = "number_of_nodes"
	NumberOfActiveNodes   = "active_node_count"
	NumberOfInactiveNodes = "inactive_node_count"
	GaugeType             = "gauge"
)

var NodeMetric = []metric.MetricConfig{
	{
		Name:  NumberOfNodes,
		Type:  GaugeType,
		Value: 0,
	},
	{
		Name:  NumberOfActiveNodes,
		Type:  GaugeType,
		Value: 0,
	},
	{
		Name:  NumberOfInactiveNodes,
		Type:  GaugeType,
		Value: 0,
	},
}
