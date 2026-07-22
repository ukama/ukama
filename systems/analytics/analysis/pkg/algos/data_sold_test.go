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

// window covering all of 2026-07-21 UTC.
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
			// two 1 GB sales on net-a inside the window
			{"package_id": "pkg-gb", "network_id": "net-a", "start_date": "2026-07-21T09:00:00Z"},
			{"package_id": "pkg-gb", "network_id": "net-a", "start_date": "2026-07-21T18:00:00Z"},
			// one 500 MB sale on net-b inside the window
			{"package_id": "pkg-mb", "network_id": "net-b", "start_date": "2026-07-21T12:00:00Z"},
			// a sale OUTSIDE the window — must be ignored
			{"package_id": "pkg-gb", "network_id": "net-a", "start_date": "2026-07-20T23:59:00Z"},
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
			{"package_id": "pkg-org", "network_id": "net-a", "start_date": "2026-07-21T09:00:00Z"},
		},
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
			{"package_id": "pkg-x", "network_id": "net-a", "start_date": "2026-07-21T09:00:00Z"},
		},
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
