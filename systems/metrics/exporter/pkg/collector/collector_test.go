package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg"
)

type TestConfig struct {
	MetricConfig []pkg.MetricConfig
	Metrics      *config.Metrics
}

func InitTestConfig() *TestConfig {
	t := &TestConfig{}
	t.MetricConfig = []pkg.MetricConfig{
		{
			Name:    "subscriber_simusage",
			Event:   "event.cloud.simmanager.sim.usage", //"event.cloud.cdr.sim.usage"}
			Type:    "histogram",
			Units:   "bytes",
			Labels:  map[string]string{"name": "usage"},
			Details: "Data Usage of the sim",
			Buckets: []float64{1024, 10240, 102400, 1024000, 10240000, 102400000},
		},
		{
			Name:    "subscriber_simusage_duration",
			Event:   "event.cloud.simmanager.sim.duration", //"event.cloud.cdr.sim.usage"}
			Type:    "histogram",
			Units:   "seconds",
			Labels:  map[string]string{"name": "usage_duration"},
			Details: "Data Usage durations",
			Buckets: []float64{60, 300, 600, 1200, 1800, 2700, 3600, 7200, 18000},
		},
		{
			Name:    "subscriber_simcount",
			Event:   "event.cloud.simmanager.sim.count", //"event.cloud.cdr.sim.usage"}
			Type:    "counter",
			Units:   "",
			Labels:  map[string]string{"name": "simcount"},
			Details: "Counter test",
		},
		{
			Name:    "subscriber_activesim",
			Event:   "event.cloud.simmanager.sim.activecount", //"event.cloud.cdr.sim.usage"}
			Type:    "gauge",
			Units:   "",
			Labels:  map[string]string{"name": "active_simcount"},
			Details: "Gauge test",
		},
		{
			Name:    "subscriber_simusage_request",
			Event:   "event.cloud.simmanager.sim.request", //"event.cloud.cdr.sim.usage"}
			Type:    "summary",
			Units:   "seconds",
			Labels:  map[string]string{"name": "sim_usage"},
			Details: "Summary test",
			Buckets: []float64{60, 300, 600, 1200, 1800, 2700, 3600, 7200, 18000},
		},
		{
			Name:    "unkown",
			Event:   "event.cloud.simmanager.sim.unkown", //"event.cloud.cdr.sim.usage"}
			Type:    "unkown",
			Units:   "",
			Labels:  map[string]string{"name": "unkown"},
			Details: "Data Usage of the sim",
		},
	}

	t.Metrics = &config.Metrics{
		Port: 10251,
	}

	return t
}

func TestCollector_NewMetricCollector(t *testing.T) {
	tC := InitTestConfig()
	nm := NewMetricsCollector(tC.MetricConfig)
	assert.NotNil(t, nm)
}

func TestCollector_GetConfigForEvent(t *testing.T) {
	tC := InitTestConfig()
	nm := NewMetricsCollector(tC.MetricConfig)

	t.Run("GetConfigSuccess", func(t *testing.T) {
		k, err := nm.GetConfigForEvent(tC.MetricConfig[1].Event)
		assert.NoError(t, err)
		if assert.NotNil(t, k) {
			assert.Equal(t, k.Name, tC.MetricConfig[1].Name)
		}
	})

	t.Run("GetConfigFailure", func(t *testing.T) {
		_, err := nm.GetConfigForEvent("UnkownEvent")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "event UnkownEvent not supported")
		}
	})

}

func TestCollector_GetMetric(t *testing.T) {
	tC := InitTestConfig()
	nm := NewMetricsCollector(tC.MetricConfig)
	m := NewMetrics(tC.MetricConfig[0].Name, tC.MetricConfig[0].Type)
	t.Run("AddMetricsSuccess", func(t *testing.T) {
		err := m.InitializeMetric(tC.MetricConfig[1].Name, tC.MetricConfig[1], nil)
		assert.NoError(t, err)
		err = nm.AddMetrics(tC.MetricConfig[1].Name, *m)
		assert.NoError(t, err)
	})

	t.Run("AddMetricsFailure_AlreadyRegistered", func(t *testing.T) {
		err := m.InitializeMetric(tC.MetricConfig[1].Name, tC.MetricConfig[1], nil)
		assert.NoError(t, err)
		err = nm.AddMetrics(tC.MetricConfig[1].Name, *m)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "already exist")
		}
	})

}
