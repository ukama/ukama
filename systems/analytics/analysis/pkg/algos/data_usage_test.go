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
	return sessionRow(iccid, simPkg, pkg, network, site, "", value)
}

func sessionRow(iccid, simPkg, pkg, network, site, session string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"iccid":          iccid,
		"sim_package_id": simPkg,
		"package_id":     pkg,
		"network_id":     network,
		"site_id":        site,
		"session_id":     session,
		"value":          value,
	}
}

func TestDataUsage(t *testing.T) {
	win := schema.Window{ID: 10}
	spec := schema.KpiSpec{Kpi: "DATA_USAGE"}

	t.Run("increment is cur minus prev per series", func(t *testing.T) {
		in := Datasets{
			"usage": {
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

		assert.Equal(t, map[string]string{
			"network_id":     "net1",
			"site_id":        "site1",
			"package_id":     "pkg1",
			"sim_package_id": "sp1",
			"iccid":          "icc1",
		}, byIccid["icc1"].Scope)
	})

	t.Run("first appearance is baseline by default", func(t *testing.T) {
		in := Datasets{
			"usage": {
				usageRow("icc1", "sp1", "pkg1", "net1", "site1", []interface{}{1.7e9, "5000"}),
			},
			"usage_prev": {},
		}

		results, err := DataUsage(win, in, spec)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, 0.0, results[0].Value)
	})

	t.Run("first appearance counts with first_value count", func(t *testing.T) {
		countSpec := spec
		countSpec.Params = map[string]string{"first_value": "count"}

		in := Datasets{
			"usage": {
				usageRow("icc1", "sp1", "pkg1", "net1", "site1", []interface{}{1.7e9, "5000"}),
			},
			"usage_prev": {},
		}

		results, err := DataUsage(win, in, countSpec)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, 5000.0, results[0].Value)
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

	t.Run("concurrent sessions sum into one scope row", func(t *testing.T) {
		// Two sessions of one sim share every scope key: they must fold into
		// ONE row or collide on the (kpi, scope, window) unique index.
		in := Datasets{
			"usage": {
				sessionRow("icc1", "sp1", "pkg1", "net1", "site1", "1", []interface{}{1.7e9, "500"}),
				sessionRow("icc1", "sp1", "pkg1", "net1", "site1", "2", []interface{}{1.7e9, "300"}),
			},
			"usage_prev": {
				sessionRow("icc1", "sp1", "pkg1", "net1", "site1", "1", []interface{}{1.7e9, "200"}),
				sessionRow("icc1", "sp1", "pkg1", "net1", "site1", "2", []interface{}{1.7e9, "100"}),
			},
		}

		results, err := DataUsage(win, in, spec)
		assert.NoError(t, err)
		assert.Len(t, results, 1, "one row per scope, not per session")
		assert.Equal(t, 500.0, results[0].Value, "300 (session 1) + 200 (session 2)")
		assert.NotContains(t, results[0].Scope, "session_id", "session is series identity, not a KPI dimension")
	})

	t.Run("a session restarting its counter does not bleed into another", func(t *testing.T) {
		// session 2 is new (no baseline) while session 1 continues.
		countSpec := spec
		countSpec.Params = map[string]string{"first_value": "count"}

		in := Datasets{
			"usage": {
				sessionRow("icc1", "sp1", "pkg1", "net1", "site1", "1", []interface{}{1.7e9, "500"}),
				sessionRow("icc1", "sp1", "pkg1", "net1", "site1", "2", []interface{}{1.7e9, "70"}),
			},
			"usage_prev": {
				sessionRow("icc1", "sp1", "pkg1", "net1", "site1", "1", []interface{}{1.7e9, "450"}),
			},
		}

		results, err := DataUsage(win, in, countSpec)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		// 50 from session 1's delta + 70 as session 2's first value.
		assert.Equal(t, 120.0, results[0].Value)
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
