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
	"strconv"
	"testing"
	"time"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// uptimeWindowSeconds is the W the uptime tests measure against.
const uptimeWindowSeconds = 60

// uptimeWindow cannot be the zero Window: its length is the denominator.
func uptimeWindow() schema.Window {
	start := time.Unix(1787257920, 0).UTC()

	return schema.Window{
		ID:    start.Unix() / uptimeWindowSeconds,
		Start: start,
		End:   start.Add(uptimeWindowSeconds * time.Second),
	}
}

func node(id, nodeType, siteID string, connectivity interface{}) map[string]interface{} {
	return map[string]interface{}{
		"node_id":      id,
		"type":         nodeType,
		"site_id":      siteID,
		"network_id":   "net-1",
		"connectivity": connectivity,
		"state":        "Operational",
	}
}

// nodeState is node() with an explicit registry state.
func nodeState(id, nodeType, siteID string, connectivity, state interface{}) map[string]interface{} {
	n := node(id, nodeType, siteID, connectivity)
	n["state"] = state

	return n
}

func health(id string, cellular, radio bool, state string) map[string]interface{} {
	return map[string]interface{}{
		"node_id":            id,
		"cellular_available": cellular,
		"radio_available":    radio,
		"radio_state":        state,
	}
}

// gained is an uptime row as /v1/last serialises it with fn=increase: a
// [ts, "value"] pair carrying the seconds gained over the window.
func gained(id string, seconds string) map[string]interface{} {
	return map[string]interface{}{
		"node_id":    id,
		"site_id":    "site-a",
		"network_id": "net-1",
		"value":      []interface{}{1787257964.364, seconds},
	}
}

// alive: up for the whole window.
func alive(id string) map[string]interface{} {
	return gained(id, "60")
}

// stalled: stopped reporting, so the pushgateway keeps serving its last value
// and increase() over the window is 0.
func stalled(id string) map[string]interface{} {
	return gained(id, "0")
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

func orEmpty(rows []map[string]interface{}) []map[string]interface{} {
	if rows == nil {
		return none()
	}

	return rows
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

func TestSiteUptimePerWindow(t *testing.T) {
	tnode := []map[string]interface{}{node("t1", "tnode", "site-a", "Online")}

	tests := []struct {
		name    string
		nodes   []map[string]interface{}
		health  []map[string]interface{}
		com     []map[string]interface{}
		ctl     []map[string]interface{}
		wantSum float64
	}{
		// the counter carried the window
		{
			name:    "a node up for the whole window is at 100",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", true, true, "on")},
			com:     []map[string]interface{}{alive("t1")},
			wantSum: 100,
		},
		{
			name:    "increase extrapolating past the window is clamped to 100",
			nodes:   tnode,
			com:     []map[string]interface{}{gained("t1", "63.4")},
			wantSum: 100,
		},
		{
			name:    "a node up for half the window is at 50",
			nodes:   tnode,
			com:     []map[string]interface{}{gained("t1", "30")},
			wantSum: 50,
		},

		// no gain: stalled or never reported
		{
			name:    "a stalled counter is down",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", true, true, "on")},
			com:     []map[string]interface{}{stalled("t1")},
			wantSum: 0,
		},
		{
			name:    "a node with no series at all is down — a missed KPI is downtime",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", true, true, "on")},
			com:     []map[string]interface{}{alive("some-other-node")},
			wantSum: 0,
		},
		{
			name:    "a NaN reading is down",
			nodes:   tnode,
			com:     []map[string]interface{}{gained("t1", "NaN")},
			wantSum: 0,
		},
		{
			name:    "an unparseable reading is down",
			nodes:   tnode,
			com:     []map[string]interface{}{gained("t1", "unavailable")},
			wantSum: 0,
		},
		{
			name:    "a negative reading is down",
			nodes:   tnode,
			com:     []map[string]interface{}{gained("t1", "-12")},
			wantSum: 0,
		},

		// the health report overrides a healthy-looking counter
		{
			name:    "service unavailable is down however alive the counter looks",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", false, true, "on")},
			com:     []map[string]interface{}{alive("t1")},
			wantSum: 0,
		},
		{
			name:    "radio unavailable is down",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", true, false, "")},
			com:     []map[string]interface{}{alive("t1")},
			wantSum: 0,
		},
		{
			name:    "radio available but switched off is still up — only availability is read",
			nodes:   tnode,
			health:  []map[string]interface{}{health("t1", true, true, "off")},
			com:     []map[string]interface{}{alive("t1")},
			wantSum: 100,
		},
		{
			name:    "an unreachable probe reports no flags, so the counter decides",
			nodes:   tnode,
			health:  []map[string]interface{}{{"node_id": "t1", "unreachable": true}},
			com:     []map[string]interface{}{alive("t1")},
			wantSum: 100,
		},
		{
			name:    "a null flag is not a false one",
			nodes:   tnode,
			health:  []map[string]interface{}{{"node_id": "t1", "radio_available": nil}},
			com:     []map[string]interface{}{alive("t1")},
			wantSum: 100,
		},
		{
			name:    "the anode carries no cellular field, so it can never trip on service",
			nodes:   []map[string]interface{}{node("a1", "anode", "site-a", "Online")},
			health:  []map[string]interface{}{{"node_id": "a1", "radio_available": true}},
			ctl:     []map[string]interface{}{alive("a1")},
			wantSum: 100,
		},

		// the site takes its least available node
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
			com:     []map[string]interface{}{alive("t1")},
			ctl:     []map[string]interface{}{alive("a1")},
			wantSum: 0,
		},
		{
			name: "the site is only as available as its weakest node",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("c1", "cnode", "site-a", "Online"),
			},
			com:     []map[string]interface{}{alive("t1"), gained("c1", "15")},
			wantSum: 25,
		},
		{
			name: "a stalled cnode takes an otherwise healthy site down",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("c1", "cnode", "site-a", "Online"),
			},
			health:  []map[string]interface{}{health("t1", true, true, "on")},
			com:     []map[string]interface{}{alive("t1"), stalled("c1")},
			wantSum: 0,
		},
		{
			name: "an hnode on the site is ignored even when it is stalled",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				node("h1", "hnode", "site-a", "Offline"),
			},
			health:  []map[string]interface{}{health("t1", true, true, "on")},
			com:     []map[string]interface{}{alive("t1"), stalled("h1")},
			wantSum: 100,
		},
		{
			name:    "a site carrying only an hnode has no node to judge and reads 100",
			nodes:   []map[string]interface{}{node("h1", "hnode", "site-a", "Online")},
			wantSum: 100,
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
			in := inputs(tt.nodes, orEmpty(tt.health), orEmpty(tt.com), orEmpty(tt.ctl))

			results, err := SiteUptime(uptimeWindow(), in, schema.KpiSpec{})
			if err != nil {
				t.Fatalf("SiteUptime: %v", err)
			}

			got := siteResult(t, results, "site-a")
			if math.Abs(got.Sum-tt.wantSum) > 1e-9 || got.Count != 1 {
				t.Errorf("got sum=%v count=%v, want sum=%v count=1",
					got.Sum, got.Count, tt.wantSum)
			}
		})
	}
}

// A short window grid still measures a range wide enough for increase() to see
// several scrapes.
func TestSiteUptimeLookbackIsFloored(t *testing.T) {
	tests := []struct {
		name       string
		windowSecs int64
		wantDenom  float64
	}{
		{name: "10s grid floors to 60s", windowSecs: 10, wantDenom: 60},
		{name: "60s grid is its own denominator", windowSecs: 60, wantDenom: 60},
		{name: "5m grid is its own denominator", windowSecs: 300, wantDenom: 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Unix(1787257920, 0).UTC()
			win := schema.Window{
				Start: start,
				End:   start.Add(time.Duration(tt.windowSecs) * time.Second),
			}

			if got := lookbackSeconds(win); got != tt.wantDenom {
				t.Fatalf("lookbackSeconds = %v, want %v", got, tt.wantDenom)
			}

			// A node that gained the full denominator reads 100.
			in := inputs(
				[]map[string]interface{}{node("t1", "tnode", "site-a", "Online")},
				[]map[string]interface{}{health("t1", true, true, "on")},
				[]map[string]interface{}{gained("t1", strconv.FormatFloat(tt.wantDenom, 'f', -1, 64))},
				none())

			results, err := SiteUptime(win, in, schema.KpiSpec{})
			if err != nil {
				t.Fatalf("SiteUptime: %v", err)
			}

			if v := siteResult(t, results, "site-a").Value; v != 100 {
				t.Errorf("site value = %v, want 100", v)
			}
		})
	}
}

// The anode reports on the ctl series, the tnode and cnode on com.
func TestSiteUptimeReadsBothUptimeSeries(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("a1", "anode", "site-a", "Online"),
	}

	in := inputs(nodes,
		[]map[string]interface{}{health("t1", true, true, "on")},
		[]map[string]interface{}{alive("t1")},
		[]map[string]interface{}{stalled("a1")})

	results, err := SiteUptime(uptimeWindow(), in, schema.KpiSpec{})
	if err != nil {
		t.Fatalf("SiteUptime: %v", err)
	}

	if v := siteResult(t, results, "site-a").Value; v != 0 {
		t.Errorf("site value = %v, want 0 — the anode's ctl counter stalled", v)
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

			if _, err := SiteUptime(uptimeWindow(), in, schema.KpiSpec{}); err == nil {
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
		[]map[string]interface{}{alive("t1"), alive("t2")}, none())

	results, err := SiteUptime(uptimeWindow(), in, schema.KpiSpec{})
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
	// 3 sites: two up all window, one whose tnode stalled.
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
	com := []map[string]interface{}{alive("t1"), stalled("t2"), alive("t3")}

	results, err := NetworkUptime(uptimeWindow(), inputs(nodes, healthRows, com, none()),
		schema.KpiSpec{})
	if err != nil {
		t.Fatalf("NetworkUptime: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 network row, got %d", len(results))
	}

	got := results[0]

	// 2 up of 3; every site is in the denominator.
	if math.Abs(got.Value-200.0/3.0) > 1e-9 {
		t.Errorf("value = %v, want %v", got.Value, 200.0/3.0)
	}

	if got.Sum != 200 || got.Count != 3 {
		t.Errorf("components sum=%v count=%v, want 200 / 3 (weighted AVG must reproduce the formula)",
			got.Sum, got.Count)
	}
}

// A partly-available site contributes its fraction, not a whole site.
func TestNetworkUptimeCarriesPartialSites(t *testing.T) {
	nodes := []map[string]interface{}{
		node("t1", "tnode", "site-a", "Online"),
		node("t2", "tnode", "site-b", "Online"),
	}
	com := []map[string]interface{}{alive("t1"), gained("t2", "15")}

	results, err := NetworkUptime(uptimeWindow(), inputs(nodes, none(), com, none()),
		schema.KpiSpec{})
	if err != nil {
		t.Fatalf("NetworkUptime: %v", err)
	}

	got := results[0]
	if math.Abs(got.Sum-125) > 1e-9 || got.Count != 2 {
		t.Errorf("components sum=%v count=%v, want 125 / 2", got.Sum, got.Count)
	}

	if math.Abs(got.Value-62.5) > 1e-9 {
		t.Errorf("value = %v, want 62.5", got.Value)
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

		// state must be Operational, not merely connected
		{
			name:  "online but still configuring is not counted",
			nodes: []map[string]interface{}{nodeState("t1", "tnode", "site-a", "Online", "Configuring")},
			want:  0,
		},
		{
			name:  "online but only Ready is not counted",
			nodes: []map[string]interface{}{nodeState("t1", "tnode", "site-a", "Online", "Ready")},
			want:  0,
		},
		{
			name:  "a faulty node is not counted",
			nodes: []map[string]interface{}{nodeState("t1", "tnode", "site-a", "Online", "Faulty")},
			want:  0,
		},
		{
			name:  "an offboarded node is not counted",
			nodes: []map[string]interface{}{nodeState("t1", "tnode", "site-a", "Online", "Offboarded")},
			want:  0,
		},
		{
			name:  "a missing state is not operational",
			nodes: []map[string]interface{}{nodeState("t1", "tnode", "site-a", "Online", nil)},
			want:  0,
		},
		{
			name:  "numeric enum 2 is operational",
			nodes: []map[string]interface{}{nodeState("t1", "tnode", "site-a", "Online", float64(2))},
			want:  1,
		},
		{
			name: "one non-operational node disqualifies the whole site",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				nodeState("a1", "anode", "site-a", "Online", "Updating"),
			},
			want: 0,
		},
		{
			name: "an hnode in a bad state is still ignored",
			nodes: []map[string]interface{}{
				node("t1", "tnode", "site-a", "Online"),
				nodeState("h1", "hnode", "site-a", "Offline", "Faulty"),
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
