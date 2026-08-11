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

func usageRow(iccid, simPkg, pkg, network, site string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"iccid":          iccid,
		"sim_package_id": simPkg,
		"package_id":     pkg,
		"network_id":     network,
		"site_id":        site,
		"value":          value,
	}
}

func TestDataUsage(t *testing.T) {
	win := schema.Window{ID: 10}
	spec := schema.KpiSpec{Kpi: "DATA_USAGE"}

	t.Run("increment is cur minus prev per series", func(t *testing.T) {
		in := Datasets{
			"usage": {
				// Prometheus sample pair: [ts, "value-as-string"].
				usageRow("icc1", "sp1", "pkg1", "net1", "site1", []interface{}{1.7e9, "400"}),
				usageRow("icc2", "sp2", "pkg1", "net1", "site1", []interface{}{1.7e9, "1000"}),
			},
			"usage_prev": {
				usageRow("icc1", "sp1", "pkg1", "net1", "site1", []interface{}{1.7e9, "150"}),
				usageRow("icc2", "sp2", "pkg1", "net1", "site1", []interface{}{1.7e9, "1000"}),
			},
		}

		results, err := DataUsage(win, in, spec)
		assert.NoError(t, err)
		assert.Len(t, results, 2)

		byIccid := map[string]Result{}
		for _, r := range results {
			byIccid[r.Scope["iccid"]] = r
		}

		// icc1 consumed 250 in the window; icc2 was idle -> explicit 0 row.
		assert.Equal(t, 250.0, byIccid["icc1"].Value)
		assert.Equal(t, 250.0, byIccid["icc1"].Sum)
		assert.Equal(t, 1.0, byIccid["icc1"].Count)
		assert.Equal(t, 0.0, byIccid["icc2"].Value)

		// Full 5-key scope on every row.
		assert.Equal(t, map[string]string{
			"network_id":     "net1",
			"site_id":        "site1",
			"package_id":     "pkg1",
			"sim_package_id": "sp1",
			"iccid":          "icc1",
		}, byIccid["icc1"].Scope)
	})

	t.Run("first appearance is baseline only", func(t *testing.T) {
		in := Datasets{
			"usage": {
				usageRow("icc1", "sp1", "pkg1", "net1", "site1", []interface{}{1.7e9, "5000"}),
			},
			"usage_prev": {},
		}

		results, err := DataUsage(win, in, spec)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		// The pre-history consumption is unknowable — counting the whole
		// counter would double-count everything consumed before ingest began.
		assert.Equal(t, 0.0, results[0].Value)
	})

	t.Run("counter reset clamps to zero", func(t *testing.T) {
		in := Datasets{
			"usage": {
				usageRow("icc1", "sp1", "pkg1", "net1", "site1", []interface{}{1.7e9, "30"}),
			},
			"usage_prev": {
				usageRow("icc1", "sp1", "pkg1", "net1", "site1", []interface{}{1.7e9, "900"}),
			},
		}

		results, err := DataUsage(win, in, spec)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, 0.0, results[0].Value)
	})

	t.Run("series without identity labels are skipped", func(t *testing.T) {
		in := Datasets{
			"usage": {
				usageRow("", "sp1", "pkg1", "net1", "site1", []interface{}{1.7e9, "100"}),
				usageRow("icc1", "", "pkg1", "net1", "site1", []interface{}{1.7e9, "100"}),
			},
			"usage_prev": {},
		}

		results, err := DataUsage(win, in, spec)
		assert.NoError(t, err)
		assert.Len(t, results, 0)
	})

	t.Run("same iccid across assignments and sites stays distinct", func(t *testing.T) {
		in := Datasets{
			"usage": {
				usageRow("icc1", "sp1", "pkg1", "net1", "siteA", []interface{}{1.7e9, "300"}),
				usageRow("icc1", "sp1", "pkg1", "net1", "siteB", []interface{}{1.7e9, "700"}),
			},
			"usage_prev": {
				usageRow("icc1", "sp1", "pkg1", "net1", "siteA", []interface{}{1.7e9, "100"}),
				usageRow("icc1", "sp1", "pkg1", "net1", "siteB", []interface{}{1.7e9, "100"}),
			},
		}

		results, err := DataUsage(win, in, spec)
		assert.NoError(t, err)
		assert.Len(t, results, 2)

		total := 0.0
		for _, r := range results {
			total += r.Sum
		}

		assert.Equal(t, 800.0, total)
	})

	t.Run("missing inputs error", func(t *testing.T) {
		_, err := DataUsage(win, Datasets{"usage": {}}, spec)
		assert.Error(t, err)

		_, err = DataUsage(win, Datasets{"usage_prev": {}}, spec)
		assert.Error(t, err)
	})
}

func TestSampleValue(t *testing.T) {
	assert.Equal(t, 42.0, sampleValue([]interface{}{1.7e9, "42"}))
	assert.Equal(t, 42.0, sampleValue([]interface{}{1.7e9, 42.0}))
	assert.Equal(t, 42.0, sampleValue("42"))
	assert.Equal(t, 42.0, sampleValue(42.0))
	assert.Equal(t, 0.0, sampleValue(nil))
	assert.Equal(t, 0.0, sampleValue([]interface{}{}))
	assert.Equal(t, 0.0, sampleValue("not-a-number"))
}
