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

type MetricType string

const (
	MetricGuage     MetricType = "guage"
	MetricCounter   MetricType = "counter"
	MetricHistogram MetricType = "histogram"
	MetricSummary   MetricType = "summary"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database  `default:"{}"`
	Grpc              *config.Grpc      `default:"{}"`
	Timeout           time.Duration     `default:"3s"`
	MsgClient         *config.MsgClient `default:"{}"`
	Queue             *config.Queue     `default:"{}"`
	Service           *config.Service   `default:"{}"`
	ExporterHost      string            `default:"localhost"`
	Org               string            `default:"40987edb-ebb6-4f84-a27c-99db7c136100"`
	IsMsgBus          bool              `default:"true"`
	KpiConfig         []KPIConfig       `default:"{}"`
	Metrics           *config.Metrics   `default:"{}"`
}

type KPIConfig struct {
	Name    string
	Event   string
	Type    MetricType
	Units   string
	Labels  map[string]string
	Details string
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: ServiceName,
		},

		Grpc: &config.Grpc{
			Port: 9092,
		},

		Service: config.LoadServiceHostConfig(name),
		MsgClient: &config.MsgClient{
			Timeout:        5 * time.Second,
			ListenerRoutes: []string{"event.cloud.simmanager.sim.usage"},
		},
		KpiConfig: []KPIConfig{
			{
				Name:    "subscriber_simusage",
				Event:   "event.cloud.simmanager.sim.usage", //"event.cloud.cdr.sim.usage"}
				Type:    MetricGuage,
				Units:   "bytes",
				Labels:  map[string]string{"name": "usage"},
				Details: "Data Usage of the sim",
			},
		},
		Metrics: &config.Metrics{
			Port: 10251,
		},
	}
}
