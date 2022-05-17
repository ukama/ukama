package pkg

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/rest"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Metrics           *config.Metrics
	Server            rest.HttpConfig
	NodeMetrics       *NodeMetricsConfig
}

type NodeMetricsConfig struct {
	Metrics             map[string]Metric
	RawQueries          map[string]Query
	MetricsServer       string
	Timeout             time.Duration
	DefaultRateInterval string
}

var defaultPrometheusMetric = map[string]Metric{
	"cpu":    {Metric: "trx_soc_cpu_usage", AggregateFunc: "sum"},
	"memory": {Metric: "trx_memory_ddr_used", AggregateFunc: "sum"},
	"users":  {Metric: "trx_lte_core_active_ue", AggregateFunc: "avg"},
}

var defautQueries = map[string]Query{
	"uptime": {Query: `min(trx_generic_system_uptime_seconds{${.Filter}})  < min(ctl_generic_system_uptime_seconds{${.Filter}})`},
}

type Metrics struct {
	conf *NodeMetricsConfig
}

type Interval struct {
	// Unix time
	Start int64
	// Unix time
	End int64
	// Step in seconds
	Step uint
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
			RawQueries:          defautQueries,
			MetricsServer:       "http://localhost:8080",
			Timeout:             time.Second * 5,
			DefaultRateInterval: "5m",
		},
	}
}
