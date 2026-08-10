/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQuery(t *testing.T) {

	t.Run("NoRate", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: false}
		r := m.getQuery(NewFilter().WithNodeId("ND12"), "1h", "avg")

		assert.Equal(t, "avg(memory {node_id='ND12'}) without (job,instance,receive,tenant_id)", r)
	})

	t.Run("NeedRate", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: true}
		r := m.getQuery(NewFilter().WithNodeId("ND12"), "1h", "avg")

		assert.Equal(t, "avg(rate(memory {node_id='ND12'}[1h])) without (job,instance,receive,tenant_id)", r)
	})

	t.Run("RateInterval", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: true, RateInterval: "1m"}
		r := m.getQuery(NewFilter().WithNodeId("ND12"), "1h", "avg")

		assert.Equal(t, "avg(rate(memory {node_id='ND12'}[1m])) without (job,instance,receive,tenant_id)", r)
	})

	t.Run("AggregateFunc", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: false}
		r := m.getQuery(NewFilter().WithNodeId("ND12"), "1h", "sum")

		assert.Equal(t, "sum(memory {node_id='ND12'})", r)
	})

}

func TestGetLastQuery(t *testing.T) {
	// Unaggregated: every matching series comes back, exactly as a direct
	// `last_over_time(data_usage[7d])` against Prometheus would return it.
	t.Run("NoFilters", func(t *testing.T) {
		m := Metric{Metric: "data_usage"}
		r := m.getLastQuery(NewFilter(), "7d")

		assert.Equal(t, "last_over_time(data_usage {}[7d])", r)
	})

	t.Run("Filtered", func(t *testing.T) {
		m := Metric{Metric: "data_usage"}
		f := NewFilter().WithNetwork("net-1").WithPackage("pkg-1").WithIccid("8910309414559836625")
		r := m.getLastQuery(f, "24h")

		assert.Equal(t,
			"last_over_time(data_usage {network='net-1',package='pkg-1',iccid='8910309414559836625'}[24h])",
			r)
	})

	// NeedRate is a range-vector concern; the instant KPI query never rates.
	t.Run("IgnoresNeedRate", func(t *testing.T) {
		m := Metric{Metric: "data_usage", NeedRate: true, RateInterval: "1m"}
		r := m.getLastQuery(NewFilter(), "7d")

		assert.Equal(t, "last_over_time(data_usage {}[7d])", r)
	})
}

func TestValidateLookback(t *testing.T) {
	for _, ok := range []string{"30s", "5m", "24h", "7d", "2w", "1y", "1h30m", "500ms"} {
		assert.NoError(t, ValidateLookback(ok), ok)
	}

	for _, bad := range []string{"", "7", "d", "7 d", "-7d", "7x", "7d]) or vector(1) (", "1h30"} {
		assert.Error(t, ValidateLookback(bad), bad)
	}
}

func TestGetFilter(t *testing.T) {
	t.Run("Package", func(t *testing.T) {
		assert.Equal(t, "package='pkg-1'", NewFilter().WithPackage("pkg-1").GetFilter())
	})

	t.Run("Iccid", func(t *testing.T) {
		assert.Equal(t, "iccid='8910300000003540855'",
			NewFilter().WithIccid("8910300000003540855").GetFilter())
	})

	t.Run("Omitted", func(t *testing.T) {
		assert.Equal(t, "site='site-1'",
			NewFilter().WithSite("site-1").WithPackage("").WithIccid("").GetFilter())
	})

	t.Run("CombinedWithOtherLabels", func(t *testing.T) {
		f := NewFilter().WithAny("net-1", "sub-1", "sim-1", "site-1", "ND12", "avg").
			WithPackage("pkg-1").WithIccid("iccid-1")

		assert.Equal(t,
			"node_id='ND12',network='net-1',subscriber='sub-1',sim='sim-1',site='site-1',package='pkg-1',iccid='iccid-1'",
			f.GetFilter())
	})
}

func TestGetAggregateQuery(t *testing.T) {
	t.Run("AggregateNet", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: true}
		r := m.getAggregateQuery(NewFilter().WithNetwork("net1"), "sum")

		assert.Equal(t, "sum(memory {network='net1'}) without (job,instance,receive,tenant_id,node_id)", r)
	})
}
