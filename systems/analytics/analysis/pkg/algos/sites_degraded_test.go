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

	"github.com/ukama/ukama/systems/analytics/schema"
)

func node(id, nodeType, siteID string, connectivity interface{}) map[string]interface{} {
	return map[string]interface{}{
		"node_id":      id,
		"type":         nodeType,
		"site_id":      siteID,
		"network_id":   "net-1",
		"connectivity": connectivity,
	}
}

func health(id string, cellular, radio bool, state string) map[string]interface{} {
	return map[string]interface{}{
		"node_id":            id,
		"cellular_available": cellular,
		"radio_available":    radio,
		"radio_state":        state,
	}
}

func degradedCount(t *testing.T, nodes, healthRows []map[string]interface{}) float64 {
	t.Helper()

	in := Datasets{
		"nodes":    nodes,
		"health":   healthRows,
		"networks": []map[string]interface{}{{"network_id": "net-1"}},
	}

	results, err := SitesDegraded(schema.Window{}, in, schema.KpiSpec{})
	if err != nil {
		t.Fatalf("SitesDegraded: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 zero-filled network row, got %d", len(results))
	}

	return results[0].Value
}

func TestSitesDegraded(t *testing.T) {
	tests := []struct {
		name   string
		nodes  []map[string]interface{}
		health []map[string]interface{}
		want   float64
	}{
		{
			name:   "healthy site is not degraded",
			nodes:  []map[string]interface{}{node("t1", "tnode", "site-a", "Online"), node("c1", "cnode", "site-a", "Online")},
			health: []map[string]interface{}{health("t1", true, true, "on")},
			want:   0,
		},
		{
			name:   "offline node degrades the site",
			nodes:  []map[string]interface{}{node("t1", "tnode", "site-a", "Offline")},
			health: []map[string]interface{}{health("t1", true, true, "on")},
			want:   1,
		},
		{
			name:   "offline cnode degrades the site even though it has no health report",
			nodes:  []map[string]interface{}{node("t1", "tnode", "site-a", "Online"), node("c1", "cnode", "site-a", "Offline")},
			health: []map[string]interface{}{health("t1", true, true, "on")},
			want:   1,
		},
		{
			name:   "tnode cellular unavailable degrades the site",
			nodes:  []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health: []map[string]interface{}{health("t1", false, true, "on")},
			want:   1,
		},
		{
			name:   "tnode radio unavailable degrades the site",
			nodes:  []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health: []map[string]interface{}{health("t1", true, false, "")},
			want:   1,
		},
		{
			name:  "anode radio unavailable degrades the site",
			nodes: []map[string]interface{}{node("t1", "tnode", "site-a", "Online"), node("a1", "anode", "site-a", "Online")},
			health: []map[string]interface{}{
				health("t1", true, true, "on"),
				// anode reports "cellular": null -> the mapper leaves the key absent
				{"node_id": "a1", "radio_available": false, "radio_state": ""},
			},
			want: 1,
		},
		{
			name:  "anode with a healthy radio and no cellular key is not degraded",
			nodes: []map[string]interface{}{node("a1", "anode", "site-a", "Online")},
			health: []map[string]interface{}{
				{"node_id": "a1", "radio_available": true, "radio_state": "on"},
			},
			want: 0,
		},
		{
			name:   "radio available but switched off is operator intent, not degraded",
			nodes:  []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health: []map[string]interface{}{health("t1", true, true, "off")},
			want:   0,
		},
		{
			name:  "unreachable health probe degrades the site",
			nodes: []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health: []map[string]interface{}{
				{"node_id": "t1", "unreachable": true},
			},
			want: 1,
		},
		{
			name:   "tnode with no health row at all degrades the site",
			nodes:  []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health: []map[string]interface{}{},
			want:   1,
		},
		{
			name:   "numeric connectivity enum (1 = online) is understood",
			nodes:  []map[string]interface{}{node("t1", "tnode", "site-a", float64(1))},
			health: []map[string]interface{}{health("t1", true, true, "on")},
			want:   0,
		},
		{
			name:   "numeric connectivity enum (2 = offline) degrades the site",
			nodes:  []map[string]interface{}{node("t1", "tnode", "site-a", float64(2))},
			health: []map[string]interface{}{health("t1", true, true, "on")},
			want:   1,
		},
		{
			name: "only the faulty site counts, not every site in the network",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("t2", "tnode", "site-b", "Online"),
			},
			health: []map[string]interface{}{
				health("t1", true, true, "on"),
				health("t2", false, true, "on"),
			},
			want: 1,
		},
		{
			name: "a node not attached to a site is ignored",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("t2", "tnode", "", "Offline"),
			},
			health: []map[string]interface{}{health("t1", true, true, "on")},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := degradedCount(t, tt.nodes, tt.health); got != tt.want {
				t.Errorf("degraded sites = %v, want %v", got, tt.want)
			}
		})
	}
}

// A network with no sites at all must still emit a zero row so the series
// stays continuous.
func TestSitesDegradedZeroFill(t *testing.T) {
	in := Datasets{
		"nodes":    []map[string]interface{}{},
		"health":   []map[string]interface{}{},
		"networks": []map[string]interface{}{{"network_id": "net-1"}, {"network_id": "net-2"}},
	}

	results, err := SitesDegraded(schema.Window{}, in, schema.KpiSpec{})
	if err != nil {
		t.Fatalf("SitesDegraded: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 zero-filled rows, got %d", len(results))
	}

	for _, r := range results {
		if r.Value != 0 || r.Count != 1 {
			t.Errorf("scope %v: got value=%v count=%v, want 0 / 1", r.Scope, r.Value, r.Count)
		}
	}
}

func TestSitesDegradedRequiresHealthInput(t *testing.T) {
	in := Datasets{
		"nodes":    []map[string]interface{}{},
		"networks": []map[string]interface{}{{"network_id": "net-1"}},
	}

	if _, err := SitesDegraded(schema.Window{}, in, schema.KpiSpec{}); err == nil {
		t.Fatal("expected an error when the health input is missing")
	}
}
