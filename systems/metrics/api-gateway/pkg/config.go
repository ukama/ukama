package pkg

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"
)

type NameUpdate struct {
	Required bool   `json:"required" default:"false"`
	Slice    string `json:"slice" default:""`
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Server            rest.HttpConfig
	Services          GrpcEndpoints  `mapstructure:"services"`
	HttpServices      HttpEndpoints  `mapstructure:"httpServices"`
	MetricsServer     config.Metrics `mapstructure:"metrics"`
	MetricsStore      string         `default:"http://localhost:8080"`
	MetricsConfig     *MetricsConfig
}

type Metric struct {
	NeedRate bool   `json:"needRate"`
	Metric   string `json:"metric"`
	// Range vector duration used in Rate func https://prometheus.io/docs/prometheus/latest/querying/basics/#time-durations
	// if NeedRate is false then this field is ignored
	// Example: 1d or 5h, or 30s
	RateInterval string `json:"rateInterval"`

	// consider adding aggregation function as a parameter
}

type MetricsConfig struct {
	Metrics             map[string]Metric `json:"metrics"`
	MetricsServer       string            `default:"http://localhost:9090"`
	Timeout             time.Duration
	DefaultRateInterval string
}

var defaultPrometheusMetric = map[string]Metric{
	"cpu":                Metric{false, "trx_soc_cpu_usage", ""},
	"memory":             Metric{false, "trx_memory_ddr_used", ""},
	"users":              Metric{false, "trx_lte_core_active_ue", ""},
	"sim_usage":          Metric{false, "sim_usage_sum", ""},
	"sim_usage_duration": Metric{false, "sim_usage_duration_sum", ""},
	"sim_count":          Metric{false, "number_of_active_sims", ""},
	"active_sims":        Metric{false, "number_of_active_sims", ""},
	"inactive_sims":      Metric{false, "number_of_inactive_sims", ""},
	"terminated_sims":    Metric{false, "number_of_terminated_sims", ""},
	"active_users":       Metric{false, "number_of_active_users", ""},
	"inactive_users":     Metric{false, "number_of_inactive_users", ""},
	"active_orgs":        Metric{false, "number_of_active_org", ""},
	"inactive_orgs":      Metric{false, "number_of_inactive_org", ""},
	"active_members":     Metric{false, "number_of_active_org_members", ""},
	"inactive_members":   Metric{false, "number_of_inactive_org_members", ""},
	"networks":           Metric{false, "number_of_networks", ""},
	"sites":              Metric{false, "number_of_sites", ""},
}

type Kratos struct {
	Url string
}

type GrpcEndpoints struct {
	Timeout  time.Duration
	Exporter string
}

type HttpEndpoints struct {
	Timeout     time.Duration
	NodeMetrics string
}

func NewConfig() *Config {
	defaultCors := cors.DefaultConfig()
	defaultCors.AllowWildcard = true
	defaultCors.AllowOrigins = []string{"http://localhost", "https://localhost"}

	return &Config{
		BaseConfig: config.BaseConfig{
			DebugMode: false,
		},

		Services: GrpcEndpoints{
			Timeout:  3 * time.Second,
			Exporter: "0.0.0.0:9090",
		},
		HttpServices: HttpEndpoints{
			Timeout:     3 * time.Second,
			NodeMetrics: "http://localhost",
		},

		Server: rest.HttpConfig{
			Port: 8080,
			Cors: defaultCors,
		},
		MetricsServer: *config.DefaultMetrics(),

		MetricsConfig: &MetricsConfig{
			Metrics:             defaultPrometheusMetric,
			Timeout:             time.Second * 5,
			DefaultRateInterval: "5m",
		},
	}
}
