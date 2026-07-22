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

// A network scope filter must keep the matching network row AND any org-wide
// (empty scope) row, so org-scoped KPIs like REVENUE still render under a
// network-scoped read.
func TestFilterScope_OrgWideRowMatchesNetworkFilter(t *testing.T) {
	orgWide := schema.CanonicalScope(nil)                                   // "{}"
	netA := schema.CanonicalScope(map[string]string{"network_id": "net-a"}) // {"network_id":"net-a"}
	netB := schema.CanonicalScope(map[string]string{"network_id": "net-b"})

	rows := []schema.KpiRollup{
		{KpiKey: "REVENUE", Scope: orgWide, Value: 12640},
		{KpiKey: "MRR", Scope: netA, Value: 3000},
		{KpiKey: "MRR", Scope: netB, Value: 9000},
	}

	got := filterScope(rows, map[string]string{"network_id": "net-a"})

	assert.ElementsMatch(t, []string{orgWide, netA}, scopeValues(got),
		"org-wide row and the matching network row survive; other networks dropped")
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
