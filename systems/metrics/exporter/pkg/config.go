package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
)

const (
	LABEL_ORG        = "org"
	LABEL_NETWROK    = "network"
	LABEL_NODE       = "node"
	LABEL_SUBSCRIBER = "susbscriber"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Grpc              *config.Grpc      `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	Queue             *config.Queue     `default:"{}"`
	Service           *config.Service   `default:"{}"`
	ExporterHost      string            `default:"localhost"`
	Org               string            `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
	IsMsgBus          bool              `default:"true"`
	MetricConfig      []MetricConfig    `default:"{}"`
	Metrics           *config.Metrics   `default:"{}"`
}

type MetricConfig struct {
	Name          string
	Event         string
	Type          string
	Units         string
	Labels        map[string]string
	DynamicLabels []string
	Details       string
	Buckets       []float64
}

func NewConfig(name string) *Config {
	return &Config{
		Grpc: &config.Grpc{
			Port: 9092,
		},
		Metrics: &config.Metrics{
			Port: 10251,
		},
		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout:        5 * time.Second,
			ListenerRoutes: []string{"event.cloud.simmanager.sim.usage"},
		},
		MetricConfig: []MetricConfig{
			{
				Name:          "sim_usage",
				Event:         "event.cloud.simmanager.sim.usage", //"event.cloud.cdr.sim.usage"}
				Type:          "histogram",
				Units:         "bytes",
				Labels:        map[string]string{"name": "usage"},
				DynamicLabels: []string{"sim", "org", "network", "subscriber", "sim_type"},
				Details:       "Data Usage of the sim",
				Buckets:       []float64{1024, 10240, 102400, 1024000, 10240000, 102400000},
			},
			{
				Name:          "sim_usage_duration",
				Event:         "event.cloud.simmanager.sim.duration", //"event.cloud.cdr.sim.usage"}
				Type:          "histogram",
				Units:         "seconds",
				Labels:        map[string]string{"name": "usage_duration"},
				DynamicLabels: []string{"sim", "org", "network", "subscriber", "sim_type"},
				Details:       "Data Usage durations",
				Buckets:       []float64{60, 300, 600, 1200, 1800, 2700, 3600, 7200, 18000},
			},
		},
	}
}
