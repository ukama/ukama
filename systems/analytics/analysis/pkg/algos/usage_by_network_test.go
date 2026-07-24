/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// Each sim's cumulative counter is summed into its network; the per-window
// value carries that cumulative total with Min=Max=Value so the read-time
// DELTA op (max-min) can recover span usage. Networks with no sims zero-fill.
func TestUsageByNetwork_CumulativePerNetwork(t *testing.T) {
	in := Datasets{
		// usage is read mode: state — one latest cumulative row per sim.
		"usage": {
			{"sim_id": "sim-A", "network_id": "net-a", "usage": map[string]interface{}{"8910": float64(100)}},
			{"sim_id": "sim-B", "network_id": "net-a", "usage": map[string]interface{}{"8911": "50"}},
			// a record with no network lineage is ignored
			{"sim_id": "sim-X", "network_id": "", "usage": map[string]interface{}{"8912": float64(999)}},
		},
		"networks": {{"network_id": "net-a"}, {"network_id": "net-b"}},
	}

	results, err := UsageByNetwork(schema.Window{}, in, schema.KpiSpec{})
	assert.NoError(t, err)

	a := rowFor(results, "net-a")
	assert.NotNil(t, a, "net-a row present")
	assert.Equal(t, float64(150), a.Value, "sum of both sims' cumulative counters")
	assert.Equal(t, a.Value, a.Min, "Min=cumulative so span DELTA works")
	assert.Equal(t, a.Value, a.Max, "Max=cumulative so span DELTA works")

	b := rowFor(results, "net-b")
	assert.NotNil(t, b, "net-b zero-filled (never absent -> never renders as —)")
	assert.Equal(t, float64(0), b.Value)
}

// Missing inputs are a hard error (mirrors the other algos).
func TestUsageByNetwork_MissingInputs(t *testing.T) {
	_, err := UsageByNetwork(schema.Window{}, Datasets{"networks": {}}, schema.KpiSpec{})
	assert.Error(t, err)

	_, err = UsageByNetwork(schema.Window{}, Datasets{"usage": {}}, schema.KpiSpec{})
	assert.Error(t, err)
}
