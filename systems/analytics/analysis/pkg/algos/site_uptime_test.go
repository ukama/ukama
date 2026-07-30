/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"math"
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

// uptimeInputs builds the two inputs the uptime algos share.
func uptimeInputs(nodes, healthRows []map[string]interface{}) Datasets {
	return Datasets{"nodes": nodes, "health": healthRows}
}

// siteResult finds the row for one site.
func siteResult(t *testing.T, results []Result, siteID string) Result {
	t.Helper()

	for _, r := range results {
		if r.Scope["site_id"] == siteID {
			return r
		}
	}

	t.Fatalf("no result row for site %q (got %d rows)", siteID, len(results))

	return Result{}
}

func TestSiteUptimeServiceAndRadio(t *testing.T) {
	tests := []struct {
		name      string
		nodes     []map[string]interface{}
		health    []map[string]interface{}
		wantSum   float64
		wantCount float64
	}{
		{
			name:      "service and radio available is up",
			nodes:     []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health:    []map[string]interface{}{health("t1", true, true, "on")},
			wantSum:   100,
			wantCount: 1,
		},
		{
			name:      "service unavailable is down",
			nodes:     []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health:    []map[string]interface{}{health("t1", false, true, "on")},
			wantSum:   0,
			wantCount: 1,
		},
		{
			name:      "radio unavailable is down",
			nodes:     []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health:    []map[string]interface{}{health("t1", true, false, "")},
			wantSum:   0,
			wantCount: 1,
		},
		{
			name:      "radio available but switched off is still up — only availability is read",
			nodes:     []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health:    []map[string]interface{}{health("t1", true, true, "off")},
			wantSum:   100,
			wantCount: 1,
		},
		{
			name:      "unreachable probe is down",
			nodes:     []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health:    []map[string]interface{}{{"node_id": "t1", "unreachable": true}},
			wantSum:   0,
			wantCount: 1,
		},
		{
			name:      "no health row is down",
			nodes:     []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health:    []map[string]interface{}{},
			wantSum:   0,
			wantCount: 1,
		},
		{
			name:  "anode needs only radio — no cellular key of its own",
			nodes: []map[string]interface{}{node("a1", "anode", "site-a", "Online")},
			health: []map[string]interface{}{
				{"node_id": "a1", "radio_available": true, "radio_state": "on"},
			},
			wantSum:   100,
			wantCount: 1,
		},
		{
			name: "one node down takes the whole site down",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("a1", "anode", "site-a", "Online"),
			},
			health: []map[string]interface{}{
				health("t1", true, false, ""),  // radio unavailable -> down
				health("a1", true, true, "on"), // up
			},
			wantSum:   0,
			wantCount: 1,
		},
		{
			name: "every bearer up makes the site up",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("a1", "anode", "site-a", "Online"),
			},
			health: []map[string]interface{}{
				health("t1", true, true, "on"),
				health("a1", true, true, "off"),
			},
			wantSum:   100,
			wantCount: 1,
		},
		{
			name:      "site with only a cnode is down, not missing",
			nodes:     []map[string]interface{}{node("c1", "cnode", "site-a", "Online")},
			health:    []map[string]interface{}{},
			wantSum:   0,
			wantCount: 1,
		},
		{
			name: "a cnode does not affect a healthy site",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("c1", "cnode", "site-a", "Offline"),
			},
			health:    []map[string]interface{}{health("t1", true, true, "on")},
			wantSum:   100,
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := SiteUptime(schema.Window{}, uptimeInputs(tt.nodes, tt.health), schema.KpiSpec{})
			if err != nil {
				t.Fatalf("SiteUptime: %v", err)
			}

			got := siteResult(t, results, "site-a")
			if got.Sum != tt.wantSum || got.Count != tt.wantCount {
				t.Errorf("got sum=%v count=%v, want sum=%v count=%v",
					got.Sum, got.Count, tt.wantSum, tt.wantCount)
			}
		})
	}
}

// Every site must emit a row, including one that carries no service or radio —
// a silently absent row is what made this KPI look broken before.
func TestSiteUptimeEmitsEverySite(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("c1", "cnode", "site-b", "Online"),
	}
	healthRows := []map[string]interface{}{health("t1", true, true, "on")}

	results, err := SiteUptime(schema.Window{}, uptimeInputs(nodes, healthRows), schema.KpiSpec{})
	if err != nil {
		t.Fatalf("SiteUptime: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected a row for both sites, got %d", len(results))
	}

	if v := siteResult(t, results, "site-a").Value; v != 100 {
		t.Errorf("site-a value = %v, want 100", v)
	}

	if v := siteResult(t, results, "site-b").Value; v != 0 {
		t.Errorf("site-b (cnode only) value = %v, want 0", v)
	}
}

func TestNetworkUptimePoolsSites(t *testing.T) {
	// 3 sites: one up, one down, one whose radio is switched off (still up).
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("t2", "tnode", "site-b", "Online"),
		node("t3", "tnode", "site-c", "Online"),
	}
	healthRows := []map[string]interface{}{
		health("t1", true, true, "on"),  // up
		health("t2", true, false, ""),   // down
		health("t3", true, true, "off"), // available, just switched off -> up
	}

	results, err := NetworkUptime(schema.Window{}, uptimeInputs(nodes, healthRows), schema.KpiSpec{})
	if err != nil {
		t.Fatalf("NetworkUptime: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 network row, got %d", len(results))
	}

	got := results[0]

	// 2 up of 3 sites; every site is in the denominator.
	if math.Abs(got.Value-200.0/3.0) > 1e-9 {
		t.Errorf("value = %v, want %v", got.Value, 200.0/3.0)
	}

	if got.Sum != 200 || got.Count != 3 {
		t.Errorf("components sum=%v count=%v, want 200 / 3 (weighted AVG must reproduce the formula)",
			got.Sum, got.Count)
	}
}

// A site with no tnode/anode has no service and no radio, so it drags the
// network percentage down rather than being quietly left out.
func TestNetworkUptimeCountsEverySite(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("c1", "cnode", "site-b", "Online"),
	}
	healthRows := []map[string]interface{}{health("t1", true, true, "on")}

	results, err := NetworkUptime(schema.Window{}, uptimeInputs(nodes, healthRows), schema.KpiSpec{})
	if err != nil {
		t.Fatalf("NetworkUptime: %v", err)
	}

	if results[0].Value != 50 || results[0].Count != 2 {
		t.Errorf("value=%v count=%v, want 50 / 2", results[0].Value, results[0].Count)
	}
}

func TestSitesOnlineAllNodes(t *testing.T) {
	tests := []struct {
		name  string
		nodes []map[string]interface{}
		want  float64
	}{
		{
			name: "every node online counts the site",
			nodes: []map[string]interface{}{
				node("c1", "cnode", "site-a", "Online"),
				node("t1", "tnode", "site-a", "Online"),
				node("a1", "anode", "site-a", "Online"),
			},
			want: 1,
		},
		{
			name: "one offline node disqualifies the site",
			nodes: []map[string]interface{}{
				node("c1", "cnode", "site-a", "Online"),
				node("t1", "tnode", "site-a", "Offline"),
			},
			want: 0,
		},
		{
			name:  "a site with no cnode still counts if its nodes are online",
			nodes: []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			want:  1,
		},
		{
			name:  "an offline cnode disqualifies the site",
			nodes: []map[string]interface{}{node("c1", "cnode", "site-a", "Offline")},
			want:  0,
		},
		{
			name:  "numeric enum 1 is online",
			nodes: []map[string]interface{}{node("t1", "tnode", "site-a", float64(1))},
			want:  1,
		},
		{
			name:  "undefined connectivity is not online",
			nodes: []map[string]interface{}{node("t1", "tnode", "site-a", float64(0))},
			want:  0,
		},
		{
			name: "an hnode attached to the site is ignored",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("h1", "hnode", "site-a", "Offline"),
			},
			want: 1,
		},
		{
			name:  "a site with only an hnode is not online",
			nodes: []map[string]interface{}{node("h1", "hnode", "site-a", "Online")},
			want:  0,
		},
		{
			name: "two sites, one fully online",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("t2", "tnode", "site-b", "Online"),
				node("a2", "anode", "site-b", "Offline"),
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := Datasets{
				"nodes":    tt.nodes,
				"networks": []map[string]interface{}{{"network_id": "net-1"}},
			}

			results, err := SitesOnline(schema.Window{}, in, schema.KpiSpec{})
			if err != nil {
				t.Fatalf("SitesOnline: %v", err)
			}

			if len(results) != 1 {
				t.Fatalf("expected 1 network row, got %d", len(results))
			}

			if results[0].Value != tt.want {
				t.Errorf("sites online = %v, want %v", results[0].Value, tt.want)
			}
		})
	}
}
