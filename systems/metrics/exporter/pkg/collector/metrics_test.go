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
		for _, ms := range cfg.Schema {
			m := NewMetrics(ms.Name, ms.Type)
			assert.NotNil(t, m)
			err := m.InitializeMetric(ms)
			if m.Type == MetricUnkown {
				assert.Contains(t, err.Error(), "not supported")
			} else {
				assert.NoError(t, err)
			}
		}
	}
}

func TestMetrics_SetMetrics(t *testing.T) {
	tC := InitTestConfig()

	for _, cfg := range tC.MetricConfig {
		for _, ms := range cfg.Schema {
			m := NewMetrics(ms.Name, ms.Type)
			assert.NotNil(t, m)
			err := m.InitializeMetric(ms)
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
}

func TestMetrics_SetupMetrics(t *testing.T) {
	tC := InitTestConfig()
	nm := NewMetricsCollector(tC.MetricConfig)

	t.Run("SetUpNewMetric", func(t *testing.T) {
		m, err := SetUpMetric(nm, tC.MetricConfig[0].Schema[0])
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.Equal(t, tC.MetricConfig[0].Schema[0].Name, m.Name)
	})

	t.Run("SetUpNewMetric_failure", func(t *testing.T) {
		_, err := SetUpMetric(nm, tC.MetricConfig[0].Schema[0])
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "already exist")
		}
	})

}
