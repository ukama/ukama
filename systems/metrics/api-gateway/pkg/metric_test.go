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
