package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetQuery(t *testing.T) {

	t.Run("NoRate", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: false}
		r := m.getQuery(NewFilter().WithNodeId("ND12"), "1h", "avg")

		assert.Equal(t, "avg(memory {nodeid='ND12'}) without (job,instance,receive,tenant_id)", r)
	})

	t.Run("NeedRate", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: true}
		r := m.getQuery(NewFilter().WithNodeId("ND12"), "1h", "avg")

		assert.Equal(t, "avg(rate(memory {nodeid='ND12'}[1h])) without (job,instance,receive,tenant_id)", r)
	})

	t.Run("RateInterval", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: true, RateInterval: "1m"}
		r := m.getQuery(NewFilter().WithNodeId("ND12"), "1h", "avg")

		assert.Equal(t, "avg(rate(memory {nodeid='ND12'}[1m])) without (job,instance,receive,tenant_id)", r)
	})

	t.Run("AggregateFunc", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: false}
		r := m.getQuery(NewFilter().WithNodeId("ND12"), "1h", "sum")

		assert.Equal(t, "sum(memory {nodeid='ND12'}) without (job,instance,receive,tenant_id)", r)
	})

}

func TestGetAggregateQuery(t *testing.T) {
	t.Run("AggregateOrg", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: true}
		r := m.getAggregateQuery(NewFilter().WithOrg("org1"), "sum")

		assert.Equal(t, "sum(memory {org='org1'}) without (job,instance,receive,tenant_id,nodeid,network)", r)
	})

	t.Run("AggregateNet", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: true}
		r := m.getAggregateQuery(NewFilter().WithNetwork("org1", "net1"), "sum")

		assert.Equal(t, "sum(memory {org='org1',network='net1'}) without (job,instance,receive,tenant_id,nodeid)", r)
	})

	t.Run("AggregateFunctionConf", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: true, AggregateFunc: "max"}
		r := m.getAggregateQuery(NewFilter().WithNetwork("org1", "net1"), "sum")

		assert.Equal(t, "max(memory {org='org1',network='net1'}) without (job,instance,receive,tenant_id,nodeid)", r)
	})
}

func TestGetLatestQuery(t *testing.T) {
	t.Run("LatestNodeId", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: true}
		r := m.getLatestQuery(NewFilter().WithNodeId("node-id"))

		assert.Equal(t, "memory {nodeid='node-id'}", r)
	})
}
