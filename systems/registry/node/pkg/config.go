package pkg

import (
	"time"

	uconf "github.com/ukama/ukama/systems/common/config"
	metric "github.com/ukama/ukama/systems/common/metrics"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Metrics          *uconf.Metrics   `default:"{}"`
	PushGateway      string           `default:"http://localhost:9091"`
	Timeout          time.Duration    `default:"3s"`
	Queue            *uconf.Queue     `default:"{}"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Service          *uconf.Service
	OrgHost          string `default:"org:9090"`
	NetworkHost      string `default:"network:9090"`
	Org              string `default:"org"`
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

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 5 * time.Second,
			ListenerRoutes: []string{
				"event.cloud.mesh.node.online",
				"event.cloud.mesh.node.offline",
			},
		},
	}
}
