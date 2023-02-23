package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_MetricTypeFromString(t *testing.T) {
	inputs := []struct {
		sType string
		mType MetricType
	}{
		{sType: "gauge", mType: MetricGuage},
		{sType: "counter", mType: MetricCounter},
		{sType: "histogram", mType: MetricHistogram},
		{sType: "summary", mType: MetricSummary},
		{sType: "unknown", mType: MetricUnkown},
	}

	for _, m := range inputs {
		result := MetricTypeFromString(m.sType)
		assert.Equal(t, result, m.mType)
	}

}

func TestMetrics_NewMetrics(t *testing.T) {
	tC := InitTestConfig()

	for _, cfg := range tC.MetricConfig {
		m := NewMetrics(cfg.Name, cfg.Type)
		assert.NotNil(t, m)
		err := m.InitializeMetric(cfg.Name, cfg, nil)
		if m.Type == MetricUnkown {
			assert.Contains(t, err.Error(), "not supported")
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestMetrics_SetMetrics(t *testing.T) {
	tC := InitTestConfig()

	for _, cfg := range tC.MetricConfig {
		m := NewMetrics(cfg.Name, cfg.Type)
		assert.NotNil(t, m)
		err := m.InitializeMetric(cfg.Name, cfg, nil)
		if m.Type == MetricUnkown {
			assert.Contains(t, err.Error(), "not supported")
		} else {
			assert.NoError(t, err)
		}

		err = m.SetMetric(1, nil)
		if m.Type == MetricUnkown {
			assert.Contains(t, err.Error(), "not supported")
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestMetrics_SetupMetrics(t *testing.T) {
	tC := InitTestConfig()
	nm := NewMetricsCollector(tC.MetricConfig)

	t.Run("SetUpNewMetric", func(t *testing.T) {
		m, err := SetUpMetric(tC.MetricConfig[0].Event, nm, tC.MetricConfig[0].Labels, tC.MetricConfig[0].Name, nil)
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.Equal(t, tC.MetricConfig[0].Name, m.Name)
	})

	t.Run("SetUpNewMetric_failure", func(t *testing.T) {
		_, err := SetUpMetric(tC.MetricConfig[0].Event, nm, tC.MetricConfig[0].Labels, tC.MetricConfig[0].Name, nil)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "already exist")
		}
	})

}
