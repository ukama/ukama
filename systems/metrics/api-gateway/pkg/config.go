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
	Auth              *config.Auth   `mapstructure:"auth"`
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
	//Health metrics
	"temperature_trx": Metric{false, "trx_sensors_tempsensor1_temperature", ""},
	"temperature_com": Metric{false, "com_sensors_tempsensor1_temperature_microprocessor", ""},

	"temperature_ctl": Metric{false, "ctl_sensors_tempsensor_microprocessor", ""},
	"temperature_rfe": Metric{false, "rfe_sensors_tempsensor_pa", ""},

	"temperature_S1_trx_hn": Metric{false, "trx_sensors_tempsensor1_temperature", ""},
	"temperature_S2_trx_hn": Metric{false, "trx_sensors_tempsensor2_temperature", ""},
	"temperature_S1_rfe_hn": Metric{false, "rfe_sensors_tempsensor1_pa1", ""},
	"temperature_S2_rfe_hn": Metric{false, "rfe_sensors_tempsensor2_pa2", ""},

	//Uptime Metrics
	"uptime_trx": Metric{false, "trx_generic_system_uptime_seconds ", ""},
	"uptime_com": Metric{false, "com_generic_system_uptime_seconds ", ""},
	"uptime_ctl": Metric{false, "ctl_generic_system_uptime_seconds ", ""},

	//Subscribers Metrics
	"subscribers_active":   Metric{false, "trx_lte_core_active_ue", ""},
	"subscribers_attached": Metric{false, "trx_lte_core_subscribers", ""},

	//Radio Metrics
	//Power Metrics (TX, RX, PA)
	"tx_power": Metric{false, "rfe_sensor_adc_tx_power", ""},
	"rx_power": Metric{false, "rfe_sensor_adc_rx_power", ""},
	"pa_power": Metric{false, "rfe_sensor_adc_pa_power", ""},

	//Resources Metrics
	//Memory Metrics (TRX, COM, CTL)
	"memory_trx_total": Metric{false, "trx_memory_ddr_total", ""},
	"memory_trx_used":  Metric{false, "trx_memory_ddr_used", ""},
	"memory_trx_free":  Metric{false, "trx_memory_ddr_free", ""},

	"memory_com_total": Metric{false, "com_memory_ddr_total", ""},
	"memory_com_used":  Metric{false, "com_memory_ddr_used", ""},
	"memory_com_free":  Metric{false, "com_memory_ddr_free", ""},

	"memory_ctl_total": Metric{false, "ctl_memory_ddr_total", ""},
	"memory_ctl_used":  Metric{false, "ctl_memory_ddr_used", ""},
	"memory_ctl_free":  Metric{false, "ctl_memory_ddr_free", ""},

	//CPU Metrics (TRX, COM, CTL)
	"cpu_trx_usage":    Metric{false, "trx_soc_cpu_usage", ""},
	"cpu_trx_c0_usage": Metric{false, "trx_soc_cpu_core0_usage", ""},
	"cpu_trx_c1_usage": Metric{false, "trx_soc_cpu_core1_usage", ""},
	"cpu_trx_c2_usage": Metric{false, "trx_soc_cpu_core2_usage", ""},
	"cpu_trx_c3_usage": Metric{false, "trx_soc_cpu_core3_usage", ""},

	"cpu_com_usage":    Metric{false, "com_soc_cpu_usage", ""},
	"cpu_com_c0_usage": Metric{false, "com_soc_cpu_core0_usage", ""},
	"cpu_com_c1_usage": Metric{false, "com_soc_cpu_core1_usage", ""},
	"cpu_com_c2_usage": Metric{false, "com_soc_cpu_core2_usage", ""},
	"cpu_com_c3_usage": Metric{false, "com_soc_cpu_core3_usage", ""},

	"cpu_ctl_total": Metric{false, "ctl_soc_cpu_usage", ""},
	"cpu_ctl_used":  Metric{false, "ctl_soc_cpu_core0_usage", ""},

	//DISK Metrics (TRX, COM, CTL)
	"disk_trx_total": Metric{false, "trx_storage_emmc_total", ""},
	"disk_trx_used":  Metric{false, "trx_storage_emmc_used", ""},
	"disk_trx_free":  Metric{false, "trx_storage_emmc_free", ""},

	"disk_com_total": Metric{false, "com_storage_emmc_total", ""},
	"disk_com_used":  Metric{false, "com_storage_emmc_used", ""},
	"disk_com_free":  Metric{false, "com_storage_emmc_free", ""},

	"disk_ctl_total": Metric{false, "ctl_storage_emmc_total", ""},
	"disk_ctl_used":  Metric{false, "ctl_storage_emmc_used", ""},
	"disk_ctl_free":  Metric{false, "ctl_storage_emmc_free", ""},

	//Power Level
	"power_level": Metric{false, "trx_sensors_powermanagement_power", ""},
	//rawQueries:
	//uptime: { query: "min(trx_generic_system_uptime_seconds{ {{ `{{.Filter}}` }} })  < min(ctl_generic_system_uptime_seconds{ {{ `{{.Filter}}` }} })" }
	//live-status: { query: "(count(last_over_time(trx_generic_system_uptime_seconds{ {{ `{{.Filter}}` }} }[1m]))/sum(last_over_time(node_count{node_type!='amplifier', {{ `{{.Filter}}` }} }[1m])))*100" }
	//"live-nodes: { query: "count(last_over_time(trx_generic_system_uptime_seconds{ {{ `{{.Filter}}` }} }[1m]))" }
}

type Kratos struct {
	_url string
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
		Auth: config.LoadAuthHostConfig("auth"),
	}
}
