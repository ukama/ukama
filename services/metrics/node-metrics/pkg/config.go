package pkg

import (
	cors "github.com/gin-contrib/cors"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/common/rest"
	"time"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           *config.Metrics
	Server            rest.HttpConfig
	NodeMetrics       *NodeMetricsConfig
}

type NodeMetricsConfig struct {
	Metrics             map[string]Metric `json:"metrics"`
	MetricsServer       string
	Timeout             time.Duration
	DefaultRateInterval string
}

var defaultPrometheusMetric = map[string]Metric{
	"cpu":    Metric{false, "trx_soc_cpu_usage", ""},
	"memory": Metric{false, "trx_memory_ddr_used", ""},
	"users":  Metric{false, "trx_lte_core_active_ue", ""},
}

func NewConfig() *Config {
	return &Config{
		Server: rest.HttpConfig{
			Port: 8080,
			Cors: cors.Config{
				AllowOrigins: []string{"http://localhost", "https://localhost"},
			},
		},
		Metrics: config.DefaultMetrics(),
		NodeMetrics: &NodeMetricsConfig{
			Metrics:             defaultPrometheusMetric,
			MetricsServer:       "http://localhost:8080",
			Timeout:             time.Second * 5,
			DefaultRateInterval: "5m",
		},
	}
}
