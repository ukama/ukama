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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/analytics/schema"
)

const (
	oneGiB = 1024 * 1024 * 1024
	oneMiB = 1024 * 1024
)

// window covering all of 2026-07-21 UTC. Flow KPIs count on first appearance
// in state, so the bounds only matter for gauges.
func testWindow() schema.Window {
	return schema.Window{
		ID:    1,
		Start: time.Date(2026, 7, 21, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2026, 7, 22, 0, 0, 0, 0, time.UTC),
	}
}

// resultFor returns the value for a network's scope row.
func resultFor(results []Result, networkID string) (float64, bool) {
	for _, r := range results {
		if r.Scope["network_id"] == networkID {
			return r.Value, true
		}
	}

	return 0, false
}

func TestDataSold_SumsAllowanceBytesPerNetwork(t *testing.T) {
	in := Datasets{
		"sim_packages": {
			// two 1 GB assignments on net-a, new this window
			{"sim_package_id": "a1", "package_id": "pkg-gb", "network_id": "net-a"},
			{"sim_package_id": "a2", "package_id": "pkg-gb", "network_id": "net-a"},
			// one 500 MB assignment on net-b, new this window
			{"sim_package_id": "b1", "package_id": "pkg-mb", "network_id": "net-b"},
			// already known one window ago — must NOT be counted again
			{"sim_package_id": "old", "package_id": "pkg-gb", "network_id": "net-a"},
		},
		"sim_packages_prev": {
			{"sim_package_id": "old"},
		},
		"packages": {
			{"package_id": "pkg-gb", "network_id": "net-a", "data_volume": 1.0, "data_unit": "gigabytes"},
			{"package_id": "pkg-mb", "network_id": "net-b", "data_volume": 500.0, "data_unit": "megabytes"},
		},
		"networks": {
			{"network_id": "net-a"},
			{"network_id": "net-b"},
			{"network_id": "net-c"}, // no sales -> zero-filled
		},
	}

	results, err := DataSold(testWindow(), in, schema.KpiSpec{})
	assert.NoError(t, err)

	a, ok := resultFor(results, "net-a")
	assert.True(t, ok)
	assert.Equal(t, float64(2*oneGiB), a, "two 1 GB sales in-window")

	b, ok := resultFor(results, "net-b")
	assert.True(t, ok)
	assert.Equal(t, float64(500*oneMiB), b, "one 500 MB sale in-window")

	c, ok := resultFor(results, "net-c")
	assert.True(t, ok)
	assert.Equal(t, float64(0), c, "no sales -> zero-filled")
}

// An org-level package (empty network_id) applies to any network.
func TestDataSold_OrgLevelPackageAppliesToNetwork(t *testing.T) {
	in := Datasets{
		"sim_packages": {
			{"sim_package_id": "a1", "package_id": "pkg-org", "network_id": "net-a"},
		},
		"sim_packages_prev": {},
		"packages": {
			{"package_id": "pkg-org", "network_id": "", "data_volume": 2.0, "data_unit": "gigabytes"},
		},
		"networks": {
			{"network_id": "net-a"},
		},
	}

	results, err := DataSold(testWindow(), in, schema.KpiSpec{})
	assert.NoError(t, err)

	a, ok := resultFor(results, "net-a")
	assert.True(t, ok)
	assert.Equal(t, float64(2*oneGiB), a)
}

// An absent/unknown data unit is treated as bytes rather than dropping to 0.
func TestDataSold_UnknownUnitTreatedAsBytes(t *testing.T) {
	in := Datasets{
		"sim_packages": {
			{"sim_package_id": "a1", "package_id": "pkg-x", "network_id": "net-a"},
		},
		"sim_packages_prev": {},
		"packages": {
			{"package_id": "pkg-x", "network_id": "net-a", "data_volume": 1048576.0},
		},
		"networks": {
			{"network_id": "net-a"},
		},
	}

	results, err := DataSold(testWindow(), in, schema.KpiSpec{})
	assert.NoError(t, err)

	a, ok := resultFor(results, "net-a")
	assert.True(t, ok)
	assert.Equal(t, float64(1048576), a, "raw volume counted as bytes")
}

// A sale counts once even when the assignment first appears long after it was
// created.
func TestPackageSales_CountsLateObservationExactlyOnce(t *testing.T) {
	assignment := map[string]interface{}{
		"sim_package_id": "late-1",
		"package_id":     "pkg-a",
		"network_id":     "net-a",
		// created hours before the window that finally observes it
		"created_at": "2026-07-20T03:00:00Z",
	}

	// Window that first sees it: counted.
	first, err := PackageSales(testWindow(), Datasets{
		"sim_packages":      {assignment},
		"sim_packages_prev": {},
	}, schema.KpiSpec{})
	assert.NoError(t, err)
	assert.Len(t, first, 1)
	assert.Equal(t, float64(1), first[0].Value)

	// Every later window already knows it: not counted again.
	later, err := PackageSales(testWindow(), Datasets{
		"sim_packages":      {assignment},
		"sim_packages_prev": {assignment},
	}, schema.KpiSpec{})
	assert.NoError(t, err)
	assert.Empty(t, later)
}

// A content change (queued package activating) is not a new sale.
func TestPackageSales_ContentChangeIsNotASale(t *testing.T) {
	results, err := PackageSales(testWindow(), Datasets{
		"sim_packages": {
			{"sim_package_id": "p1", "package_id": "pkg-a", "network_id": "net-a", "is_active": true},
		},
		"sim_packages_prev": {
			{"sim_package_id": "p1", "package_id": "pkg-a", "network_id": "net-a", "is_active": false},
		},
	}, schema.KpiSpec{})
	assert.NoError(t, err)
	assert.Empty(t, results)
}

// _entity_key from the change-log wins over any mapped id field.
func TestPackageSales_PrefersChangeLogEntityKey(t *testing.T) {
	results, err := PackageSales(testWindow(), Datasets{
		"sim_packages": {
			{FieldEntityKey: "e1", "package_id": "pkg-a", "network_id": "net-a"},
		},
		"sim_packages_prev": {
			{FieldEntityKey: "e1"},
		},
	}, schema.KpiSpec{})
	assert.NoError(t, err)
	assert.Empty(t, results, "same entity key -> already counted")
}

// Missing the lag-1 baseline is a spec error, not a silent zero.
func TestPackageSales_RequiresPrevBaseline(t *testing.T) {
	_, err := PackageSales(testWindow(), Datasets{
		"sim_packages": {},
	}, schema.KpiSpec{})
	assert.Error(t, err)
}
