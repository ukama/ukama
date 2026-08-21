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

// up is an uptime row as /v1/last serialises it: the [ts, "value"] pair of a
// node still pushing its system uptime counter.
func up(id string, seconds string) map[string]interface{} {
	return map[string]interface{}{
		"node_id":    id,
		"site_id":    "site-a",
		"network_id": "net-1",
		"value":      []interface{}{1787257964.364, seconds},
	}
}

// dead is the same row for a node that stopped pushing: Prometheus' staleness
// marker, forwarded by the sanitizer and serialised as the string "NaN".
func dead(id string) map[string]interface{} {
	return up(id, "NaN")
}

// inputs builds the four datasets. Nothing is auto-filled: a test spells out
// only the evidence it wants to exist.
func inputs(nodes, healthRows, comRows, ctlRows []map[string]interface{}) Datasets {
	return Datasets{
		"nodes":      nodes,
		"health":     healthRows,
		"com_uptime": comRows,
		"ctl_uptime": ctlRows,
	}
}

func none() []map[string]interface{} { return []map[string]interface{}{} }

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

func TestSiteUptimeDownNeedsEvidence(t *testing.T) {
	tnode := []map[string]interface{}{node("t1", "tnode", "site-a", "Online")}

	tests := []struct {
		name    string
		nodes   []map[string]interface{}
		health  []map[string]interface{}
		com     []map[string]interface{}
		ctl     []map[string]interface{}
		wantSum float64
	}{
		// nothing reported: up, by default
		{
			name:    "a freshly created site nothing has reported on is up",
			nodes:   tnode,
			wantSum: 100,
		},
		{
			name:    "an unreachable probe is NOT downtime — it reports nothing",
			nodes:   tnode,
			health:  []map[string]interface{}{{"node_id": "t1", "unreachable": true}},
			wantSum: 100,
		},
		{
			name:    "a health row with no interface flags is not downtime",
			nodes:   tnode,
			health:  []map[string]interface{}{{"node_id": "t1"}},
			wantSum: 100,
		},
		{
			name:    "a node with no uptime series yet is up",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", true, true, "on")},
			com:     []map[string]interface{}{up("some-other-node", "1")},
			wantSum: 100,
		},
		{
			name:    "a null flag is not a false one",
			nodes:   tnode,
			health:  []map[string]interface{}{{"node_id": "t1", "radio_available": nil}},
			wantSum: 100,
		},

		// the node's own counter says it died
		{
			name:    "a NaN uptime reading takes the node down",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", true, true, "on")},
			com:     []map[string]interface{}{dead("t1")},
			wantSum: 0,
		},
		{
			name:    "a zero uptime reading is alive — the node just booted",
			nodes:   tnode,
			com:     []map[string]interface{}{up("t1", "0")},
			wantSum: 100,
		},
		{
			name:    "an unparseable reading counts as dead",
			nodes:   tnode,
			com:     []map[string]interface{}{up("t1", "unavailable")},
			wantSum: 0,
		},

		// the health report says an interface is gone
		{
			name:    "service unavailable is down",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", false, true, "on")},
			com:     []map[string]interface{}{up("t1", "542344.69")},
			wantSum: 0,
		},
		{
			name:    "radio unavailable is down",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", true, false, "")},
			com:     []map[string]interface{}{up("t1", "542344.69")},
			wantSum: 0,
		},
		{
			name:    "radio available but switched off is still up — only availability is read",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", true, true, "off")},
			com:     []map[string]interface{}{up("t1", "542344.69")},
			wantSum: 100,
		},
		{
			name:    "the anode carries no cellular field, so it can never trip on service",
			nodes:   []map[string]interface{}{node("a1", "anode", "site-a", "Online")},
			health:  []map[string]interface{}{{"node_id": "a1", "radio_available": true}},
			ctl:     []map[string]interface{}{up("a1", "542337.11")},
			wantSum: 100,
		},

		// the site is the conjunction of its nodes
		{
			name: "one node down takes the whole site down",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("a1", "anode", "site-a", "Online"),
			},
			health: []map[string]interface{}{
				health("t1", true, false, ""),
				{"node_id": "a1", "radio_available": true},
			},
			wantSum: 0,
		},
		{
			name: "a dead cnode takes an otherwise healthy site down",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("c1", "cnode", "site-a", "Online"),
			},
			health:  []map[string]interface{}{health("t1", true, true, "on")},
			com:     []map[string]interface{}{up("t1", "542344.69"), dead("c1")},
			wantSum: 0,
		},
		{
			name: "an hnode on the site is ignored even when it is dead",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("h1", "hnode", "site-a", "Offline"),
			},
			health:  []map[string]interface{}{health("t1", true, true, "on")},
			com:     []map[string]interface{}{up("t1", "542344.69"), dead("h1")},
			wantSum: 100,
		},
		{
			name:    "a site carrying only an hnode has nothing reported against it",
			nodes:   []map[string]interface{}{node("h1", "hnode", "site-a", "Online")},
			wantSum: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := inputs(tt.nodes, orEmpty(tt.health), orEmpty(tt.com), orEmpty(tt.ctl))

			results, err := SiteUptime(schema.Window{}, in, schema.KpiSpec{})
			if err != nil {
				t.Fatalf("SiteUptime: %v", err)
			}

			got := siteResult(t, results, "site-a")
			if got.Sum != tt.wantSum || got.Count != 1 {
				t.Errorf("got sum=%v count=%v, want sum=%v count=1",
					got.Sum, got.Count, tt.wantSum)
			}
		})
	}
}

func orEmpty(rows []map[string]interface{}) []map[string]interface{} {
	if rows == nil {
		return none()
	}

	return rows
}

// A site created moments ago, whose nodes have not been probed and have no
// uptime series, reads 100%.
func TestSiteUptimeNewSiteStartsAt100(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("a1", "anode", "site-a", "Online"),
		node("c1", "cnode", "site-a", "Online"),
	}
	// Every probe unreachable, no metrics at all.
	unreachable := []map[string]interface{}{
		{"node_id": "t1", "unreachable": true},
		{"node_id": "a1", "unreachable": true},
	}

	results, err := SiteUptime(schema.Window{}, inputs(nodes, unreachable, none(), none()),
		schema.KpiSpec{})
	if err != nil {
		t.Fatalf("SiteUptime: %v", err)
	}

	if v := siteResult(t, results, "site-a").Value; v != 100 {
		t.Errorf("new site value = %v, want 100 — nothing has reported a failure", v)
	}
}

// The anode's liveness arrives on the ctl series, the tnode's and cnode's on
// the com series; the algo folds both.
func TestSiteUptimeReadsBothUptimeSeries(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("a1", "anode", "site-a", "Online"),
	}

	in := inputs(nodes,
		[]map[string]interface{}{health("t1", true, true, "on")},
		[]map[string]interface{}{up("t1", "542344.69")},
		[]map[string]interface{}{dead("a1")})

	results, err := SiteUptime(schema.Window{}, in, schema.KpiSpec{})
	if err != nil {
		t.Fatalf("SiteUptime: %v", err)
	}

	if v := siteResult(t, results, "site-a").Value; v != 0 {
		t.Errorf("site value = %v, want 0 — the anode's ctl series is NaN", v)
	}
}

func TestSiteUptimeMissingInputs(t *testing.T) {
	full := Datasets{"nodes": {}, "health": {}, "com_uptime": {}, "ctl_uptime": {}}

	for _, missing := range []string{"nodes", "health", "com_uptime", "ctl_uptime"} {
		t.Run("missing "+missing, func(t *testing.T) {
			in := Datasets{}

			for k, v := range full {
				if k != missing {
					in[k] = v
				}
			}

			if _, err := SiteUptime(schema.Window{}, in, schema.KpiSpec{}); err == nil {
				t.Errorf("expected an error when input %q is absent", missing)
			}
		})
	}
}

// Every site must emit a row.
func TestSiteUptimeEmitsEverySite(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("t2", "tnode", "site-b", "Online"),
	}
	in := inputs(nodes,
		[]map[string]interface{}{health("t1", true, true, "on"), health("t2", true, false, "")},
		none(), none())

	results, err := SiteUptime(schema.Window{}, in, schema.KpiSpec{})
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
		t.Errorf("site-b value = %v, want 0 — its radio is reported gone", v)
	}
}

func TestNetworkUptimePoolsSites(t *testing.T) {
	// 3 sites: one healthy, one whose radio is switched off (still up), one
	// whose tnode died.
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("t2", "tnode", "site-b", "Online"),
		node("t3", "tnode", "site-c", "Online"),
	}
	healthRows := []map[string]interface{}{
		health("t1", true, true, "on"),
		health("t2", true, true, "on"),
		health("t3", true, true, "off"), // available, just switched off -> up
	}

	in := inputs(nodes, healthRows, []map[string]interface{}{dead("t2")}, none())

	results, err := NetworkUptime(schema.Window{}, in, schema.KpiSpec{})
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

// A network with nothing reported against it is at 100%.
func TestNetworkUptimeNewNetworkStartsAt100(t *testing.T) {
	nodes := []map[string]interface{}{node("t1", "tnode", "site-a", "Online")}

	results, err := NetworkUptime(schema.Window{}, inputs(nodes, none(), none(), none()),
		schema.KpiSpec{})
	if err != nil {
		t.Fatalf("NetworkUptime: %v", err)
	}

	if results[0].Value != 100 || results[0].Count != 1 {
		t.Errorf("value=%v count=%v, want 100 / 1", results[0].Value, results[0].Count)
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
