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

func TestGetAggregateQuery(t *testing.T) {
	t.Run("AggregateNet", func(t *testing.T) {
		m := Metric{Metric: "memory", NeedRate: true}
		r := m.getAggregateQuery(NewFilter().WithNetwork("net1"), "sum")

		assert.Equal(t, "sum(memory {network='net1'}) without (job,instance,receive,tenant_id,node_id)", r)
	})
}
