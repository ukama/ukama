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
// node that is still pushing its system uptime counter.
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

// uptimeInputs builds the four inputs the uptime algos share. Nodes with no
// explicit uptime row are given a live one, so a test only has to spell out
// the liveness it actually cares about.
func uptimeInputs(nodes, healthRows []map[string]interface{},
	uptimeRows ...map[string]interface{}) Datasets {
	given := map[string]bool{}
	rows := append([]map[string]interface{}{}, uptimeRows...)

	for _, r := range uptimeRows {
		given[str(r["node_id"])] = true
	}

	for _, n := range nodes {
		id := str(n["node_id"])
		if !given[id] {
			rows = append(rows, up(id, "542344.69"))
		}
	}

	// The split between the two series does not matter to the algo — it folds
	// them into one map — so the fixtures put every row in com_uptime and
	// keep a dedicated case for the ctl (anode) series.
	return Datasets{
		"nodes":      nodes,
		"health":     healthRows,
		"com_uptime": rows,
		"ctl_uptime": []map[string]interface{}{},
	}
}

// noUptime is the degraded shape: both uptime datasets empty, which means the
// metrics path delivered nothing rather than that every node died.
func noUptime(nodes, healthRows []map[string]interface{}) Datasets {
	return Datasets{
		"nodes":      nodes,
		"health":     healthRows,
		"com_uptime": []map[string]interface{}{},
		"ctl_uptime": []map[string]interface{}{},
	}
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
		uptime    []map[string]interface{}
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
			name:      "a healthy node that stopped pushing uptime is down",
			nodes:     []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health:    []map[string]interface{}{health("t1", true, true, "on")},
			uptime:    []map[string]interface{}{dead("t1")},
			wantSum:   0,
			wantCount: 1,
		},
		{
			name:      "a zero uptime reading still counts as reporting",
			nodes:     []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
			health:    []map[string]interface{}{health("t1", true, true, "on")},
			uptime:    []map[string]interface{}{up("t1", "0")},
			wantSum:   100,
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
			name: "every node up makes the site up",
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
			name:      "a site with only a live cnode is up — v3 judges the cnode too",
			nodes:     []map[string]interface{}{node("c1", "cnode", "site-a", "Online")},
			health:    []map[string]interface{}{},
			wantSum:   100,
			wantCount: 1,
		},
		{
			name:      "a site with only a dead cnode is down, not missing",
			nodes:     []map[string]interface{}{node("c1", "cnode", "site-a", "Online")},
			health:    []map[string]interface{}{},
			uptime:    []map[string]interface{}{dead("c1")},
			wantSum:   0,
			wantCount: 1,
		},
		{
			name: "a dead cnode takes an otherwise healthy site down",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("c1", "cnode", "site-a", "Online"),
			},
			health:    []map[string]interface{}{health("t1", true, true, "on")},
			uptime:    []map[string]interface{}{dead("c1")},
			wantSum:   0,
			wantCount: 1,
		},
		{
			name: "an hnode on the site is still ignored",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("h1", "hnode", "site-a", "Offline"),
			},
			health:    []map[string]interface{}{health("t1", true, true, "on")},
			uptime:    []map[string]interface{}{dead("h1")},
			wantSum:   100,
			wantCount: 1,
		},
		{
			name:      "a site with only an hnode has nothing to serve with",
			nodes:     []map[string]interface{}{node("h1", "hnode", "site-a", "Online")},
			health:    []map[string]interface{}{},
			wantSum:   0,
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := uptimeInputs(tt.nodes, tt.health, tt.uptime...)

			results, err := SiteUptime(schema.Window{}, in, schema.KpiSpec{})
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

// A node whose series is missing entirely — it aged out of the pull's
// lookback rather than going stale — is down for the same reason a NaN is:
// nothing confirms it is alive.
func TestSiteUptimeAbsentSeriesIsDown(t *testing.T) {
	in := Datasets{
		"nodes":      []map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
		"health":     []map[string]interface{}{health("t1", true, true, "on")},
		"com_uptime": []map[string]interface{}{up("some-other-node", "1")},
		"ctl_uptime": []map[string]interface{}{},
	}

	results, err := SiteUptime(schema.Window{}, in, schema.KpiSpec{})
	if err != nil {
		t.Fatalf("SiteUptime: %v", err)
	}

	if v := siteResult(t, results, "site-a").Value; v != 0 {
		t.Errorf("site value = %v, want 0 — t1 has no uptime series", v)
	}
}

// The anode's liveness arrives on the ctl series, the tnode's and cnode's on
// the com series; the algo folds both.
func TestSiteUptimeReadsBothUptimeSeries(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("a1", "anode", "site-a", "Online"),
	}
	healthRows := []map[string]interface{}{
		health("t1", true, true, "on"),
		{"node_id": "a1", "radio_available": true},
	}

	in := Datasets{
		"nodes":      nodes,
		"health":     healthRows,
		"com_uptime": []map[string]interface{}{up("t1", "542344.69")},
		"ctl_uptime": []map[string]interface{}{dead("a1")},
	}

	results, err := SiteUptime(schema.Window{}, in, schema.KpiSpec{})
	if err != nil {
		t.Fatalf("SiteUptime: %v", err)
	}

	if v := siteResult(t, results, "site-a").Value; v != 0 {
		t.Errorf("site value = %v, want 0 — the anode's ctl series is NaN", v)
	}
}

// Empty uptime datasets mean the metrics path delivered nothing, not that
// every node in the org died: the liveness term is dropped and the health
// rule alone decides, so a broken sanitizer cannot zero the whole KPI.
func TestSiteUptimeDegradesWhenUptimeDatasetIsEmpty(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("c1", "cnode", "site-a", "Online"),
		node("t2", "tnode", "site-b", "Online"),
	}
	healthRows := []map[string]interface{}{
		health("t1", true, true, "on"),
		health("t2", true, false, ""), // radio unavailable -> still down
	}

	results, err := SiteUptime(schema.Window{}, noUptime(nodes, healthRows), schema.KpiSpec{})
	if err != nil {
		t.Fatalf("SiteUptime: %v", err)
	}

	if v := siteResult(t, results, "site-a").Value; v != 100 {
		t.Errorf("site-a value = %v, want 100 — no uptime signal must not read as down", v)
	}

	if v := siteResult(t, results, "site-b").Value; v != 0 {
		t.Errorf("site-b value = %v, want 0 — the health rule still applies", v)
	}
}

func TestSiteUptimeMissingInputs(t *testing.T) {
	full := Datasets{
		"nodes":      {},
		"health":     {},
		"com_uptime": {},
		"ctl_uptime": {},
	}

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

// Every site must emit a row, including one that carries nothing that can
// serve — a silently absent row is what made this KPI look broken before.
func TestSiteUptimeEmitsEverySite(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("h1", "hnode", "site-b", "Online"),
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
		t.Errorf("site-b (hnode only) value = %v, want 0", v)
	}
}

func TestNetworkUptimePoolsSites(t *testing.T) {
	// 3 sites: one up, one whose radio is switched off (still up), one whose
	// tnode stopped pushing uptime.
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("t2", "tnode", "site-b", "Online"),
		node("t3", "tnode", "site-c", "Online"),
	}
	healthRows := []map[string]interface{}{
		health("t1", true, true, "on"),  // up
		health("t2", true, true, "on"),  // healthy, but dark below -> down
		health("t3", true, true, "off"), // available, just switched off -> up
	}

	in := uptimeInputs(nodes, healthRows, dead("t2"))

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

// A site with nothing that can serve drags the network percentage down rather
// than being quietly left out.
func TestNetworkUptimeCountsEverySite(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("h1", "hnode", "site-b", "Online"),
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
