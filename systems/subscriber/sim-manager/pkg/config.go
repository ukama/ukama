package pkg

import (
	"strings"
	"time"

	"github.com/ukama/ukama/systems/common/config"
	metric "github.com/ukama/ukama/systems/common/pushgatewayMetrics"
)

const (
	NumberOfSubscribers = "number_of_subscribers"
	ActiveCount         = "active_sim_count"
	InactiveCount       = "inactive_sim_count"
	TerminatedCount     = "terminated_sim_count"
	GaugeType           = "gauge"
)

type MetricConfig struct {
	Name    string
	Type    string
	Labels  map[string]string
	Details string
	Buckets []float64
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Queue             *config.Queue     `default:"{}"`
	Metrics           *config.Metrics   `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	Service           *config.Service
	Key               string
	DataPlan          string `default:"http://data-plan:8080"`
	Registry          string `default:"registry:9090"`
	SimPool           string `default:"sim:9090"`
	TestAgent         string `default:"test-agent:9090"`
	OperatorAgent     string `default:"http://operator-agent:8080"`
	OrgHost           string `default:"http://registry-api-gw:8080"`
	Org               string `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
	PushMetricHost    string `default:"http://localhost:9091"`
}

func NewConfig(name string) *Config {
	// Sanitize name
	name = strings.ReplaceAll(name, "-", "_")
	return &Config{
		DB: &config.Database{
			DbName: name,
		},
		Service: config.LoadServiceHostConfig(name),

		MsgClient: &config.MsgClient{
			Timeout: 5 * time.Second,
		},
	}
}

type SimMetrics struct {
	Name   string
	Type   string
	Labels map[string]string
	Value  float64
}

var SimMetric = []metric.MetricConfig{{
	Name:   NumberOfSubscribers,
	Type:   GaugeType,
	Labels: map[string]string{"network": "", "org": ""},
	Value:  0,
},
	{
		Name:   ActiveCount,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
	{
		Name:   InactiveCount,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
	{
		Name:   TerminatedCount,
		Type:   GaugeType,
		Labels: map[string]string{"network": "", "org": ""},
		Value:  0,
	},
}
