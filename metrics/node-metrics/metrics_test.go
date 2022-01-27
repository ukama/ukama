package nodemetrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetQuery(t *testing.T) {

	t.Run("NoRate", func(t *testing.T) {
		m := metricQuery{Metric: "memory", NeedRate: false}
		r := m.getQuery("ND12", "")

		assert.Equal(t, "avg(memory {nodeid='ND12'}) without (job, instance)", r)
	})

	t.Run("NeedRate", func(t *testing.T) {
		m := metricQuery{Metric: "memory", NeedRate: true}
		r := m.getQuery("ND12", "")

		assert.Equal(t, "avg(rate(memory {nodeid='ND12'}[1h])) without (job, instance)", r)
	})

}
