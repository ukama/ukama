/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/systems/analytics/schema"
)

func scopeValues(rows []schema.KpiRollup) []string {
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.Scope)
	}

	return out
}

// With per-network attribution, a network filter returns ONLY the matching
// network row. The org bucket (empty scope, e.g. revenue from payments whose
// SIM couldn't be resolved) is org-only and must NOT bleed into a single
// network's number.
func TestFilterScope_OrgBucketExcludedFromNetworkFilter(t *testing.T) {
	orgWide := schema.CanonicalScope(nil)                                   // "{}"
	netA := schema.CanonicalScope(map[string]string{"network_id": "net-a"}) // {"network_id":"net-a"}
	netB := schema.CanonicalScope(map[string]string{"network_id": "net-b"})

	rows := []schema.KpiRollup{
		{KpiKey: "REVENUE", Scope: orgWide, Value: 12640},
		{KpiKey: "REVENUE", Scope: netA, Value: 100},
		{KpiKey: "REVENUE", Scope: netB, Value: 900},
	}

	got := filterScope(rows, map[string]string{"network_id": "net-a"})

	assert.Equal(t, []string{netA}, scopeValues(got),
		"only the matching network row; org bucket and other networks dropped")
}

// An empty filter returns every row unchanged.
func TestFilterScope_NoFilterReturnsAll(t *testing.T) {
	rows := []schema.KpiRollup{
		{Scope: schema.CanonicalScope(nil)},
		{Scope: schema.CanonicalScope(map[string]string{"network_id": "net-a"})},
	}

	got := filterScope(rows, nil)

	assert.Len(t, got, 2)
}

// A network filter drops rows for other networks.
func TestFilterScope_FiltersNonMatchingNetwork(t *testing.T) {
	netA := schema.CanonicalScope(map[string]string{"network_id": "net-a"})
	netB := schema.CanonicalScope(map[string]string{"network_id": "net-b"})

	rows := []schema.KpiRollup{
		{Scope: netA},
		{Scope: netB},
	}

	got := filterScope(rows, map[string]string{"network_id": "net-a"})

	assert.Equal(t, []string{netA}, scopeValues(got))
}
