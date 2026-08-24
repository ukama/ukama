/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package algos

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// minUptimeLookback floors the range the increase is measured over: it has to
// span several scrape intervals, or Prometheus returns a single sample (and
// drops the series) or two identical ones (increase 0), reading a live node as
// down. Below the floor consecutive windows overlap, each reporting the
// trailing 60s. Ingest applies the same floor via {{ atLeast .WindowSeconds
// 60 }}, so the two cannot drift.
const minUptimeLookback = 60 * time.Second

// The node types that make up a site. Any other type attached to a site
// (hnode) is ignored by the uptime KPIs.
const (
	nodeTypeTower      = "tnode"
	nodeTypeAmplifier  = "anode"
	nodeTypeController = "cnode"
)

// SiteUptime (SITE_UPTIME @ network_id+site_id), percentage of the window.
//
//	node % = clamp(increase / L, 0, 1) x 100, forced to 0 by a health flag
//	         reported false
//	site % = min over the site's tnode/anode/cnode
//
// L is the window length floored at minUptimeLookback. Ingest pulls each
// node's uptime counter as `increase` over that same range, so the input
// value is seconds of uptime gained over it: ~L for a node up throughout, 0
// for one whose counter stalled or that never reported.
//
// Emits Sum = the site's percentage, Count = 1, so the aggregator's AVG over
// any span is the mean of its windows. Query with op=AVG.
func SiteUptime(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, err := classifySites(in, lookbackSeconds(win), "SITE_UPTIME")
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(sites))

	for siteID, agg := range sites {
		value := agg.percent()

		results = append(results, Result{
			Scope: map[string]string{"network_id": agg.networkID, "site_id": siteID},
			Value: value,
			Sum:   value,
			Count: 1,
			Min:   value,
			Max:   value,
		})
	}

	return results, nil
}

// NetworkUptime (NETWORK_UPTIME @ network_id), percentage of the window:
//
//	network % = sum(site %) / site count
//
// Emits Sum = the total of its sites' percentages, Count = its site count, so
// the aggregator's AVG spans every (site, window) pair rather than averaging
// per-site averages: each site is weighted by the site-time it lost, and a
// site that joined mid-span counts only for the windows it existed for.
//
// Per-site rule and inputs: identical to SITE_UPTIME.
func NetworkUptime(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, err := classifySites(in, lookbackSeconds(win), "NETWORK_UPTIME")
	if err != nil {
		return nil, err
	}

	type netAgg struct{ total, count float64 }

	networks := map[string]*netAgg{}

	for _, agg := range sites {
		n, ok := networks[agg.networkID]
		if !ok {
			n = &netAgg{}
			networks[agg.networkID] = n
		}

		n.total += agg.percent()
		n.count++
	}

	results := make([]Result, 0, len(networks))

	for networkID, n := range networks {
		value := 0.0
		if n.count > 0 {
			value = n.total / n.count
		}

		results = append(results, Result{
			Scope: map[string]string{"network_id": networkID},
			Value: value,
			Sum:   n.total,
			Count: n.count,
			Min:   value,
			Max:   value,
		})
	}

	return results, nil
}

// lookbackSeconds is the denominator every node's uptime gain is measured
// against: the window length floored at minUptimeLookback, which is the range
// ingest asked Prometheus for.
func lookbackSeconds(win schema.Window) float64 {
	seconds := win.End.Sub(win.Start).Seconds()
	if floor := minUptimeLookback.Seconds(); seconds < floor {
		return floor
	}

	return seconds
}

// siteAgg holds one site's state for a window. The site is only as available
// as its least available node, so worst starts at 100 — a site with no node
// to judge is up — and each node pulls it down in turn.
type siteAgg struct {
	networkID string
	worst     float64
}

func (a *siteAgg) observe(nodePercent float64) {
	if nodePercent < a.worst {
		a.worst = nodePercent
	}
}

func (a *siteAgg) percent() float64 { return a.worst }

// classifySites builds the per-window site state shared by SITE_UPTIME and
// NETWORK_UPTIME.
func classifySites(in Datasets, lookback float64, kpi string) (map[string]*siteAgg, error) {
	health, ok := in["health"]
	if !ok {
		return nil, fmt.Errorf("%s: missing input 'health'", kpi)
	}

	nodes, ok := in["nodes"]
	if !ok {
		return nil, fmt.Errorf("%s: missing input 'nodes'", kpi)
	}

	com, ok := in["com_uptime"]
	if !ok {
		return nil, fmt.Errorf("%s: missing input 'com_uptime'", kpi)
	}

	ctl, ok := in["ctl_uptime"]
	if !ok {
		return nil, fmt.Errorf("%s: missing input 'ctl_uptime'", kpi)
	}

	healthByNode := indexHealthByNode(health)
	gainByNode := indexUptimeGainByNode(com, ctl)

	sites := map[string]*siteAgg{}

	// Registry is the authority on site membership, not the metric's own site
	// label: a dark node has no metric left to carry a label on.
	for _, node := range nodes {
		siteID, networkID := str(node["site_id"]), str(node["network_id"])
		if siteID == "" || networkID == "" {
			continue // unattached node
		}

		// Every node type registers the site, so a site made only of hnodes
		// still produces a row.
		agg, ok := sites[siteID]
		if !ok {
			agg = &siteAgg{networkID: networkID, worst: 100}
			sites[siteID] = agg
		}

		nodeType := strings.ToLower(str(node["type"]))
		if !isSiteNodeType(nodeType) {
			continue
		}

		nodeID := str(node["node_id"])

		agg.observe(nodeUptimePercent(nodeType, healthByNode[nodeID],
			gainByNode[nodeID], lookback))
	}

	return sites, nil
}

func isSiteNodeType(nodeType string) bool {
	return nodeType == nodeTypeTower ||
		nodeType == nodeTypeAmplifier ||
		nodeType == nodeTypeController
}

// nodeUptimePercent is how much of the window this node was serving for. The
// uptime counter sets the ceiling; a health flag reported false overrides it,
// since an unavailable interface means the node was not serving however alive
// its counter looked.
//
// radio.state is not read: uptime tracks interface AVAILABILITY, so a radio
// that is available but switched off is still up. Neither is node lifecycle
// state or node-gateway reachability — the node's own counter is the
// authority on whether it is alive.
func nodeUptimePercent(nodeType string, h map[string]interface{},
	gain float64, lookback float64) float64 {
	if flagIsFalse(h, "radio_available") {
		return 0
	}

	// Cellular is the service interface, and tnode-only: the anode reports
	// "cellular": null, so the field never reaches its row.
	if nodeType == nodeTypeTower && flagIsFalse(h, "cellular_available") {
		return 0
	}

	if lookback <= 0 {
		return 0
	}

	// increase() extrapolates to the edges of its range, so a node up
	// throughout can come back a little over (or under) the lookback.
	ratio := gain / lookback
	if ratio > 1 {
		ratio = 1
	}

	if ratio < 0 {
		ratio = 0
	}

	return ratio * 100
}

// flagIsFalse reports whether a health flag is present AND says false. A
// missing field is not a false one: the ingest mapper writes a field only when
// its path resolves, so an unreachable probe records the bare
// `unreachable: true` marker with no interface flags at all.
func flagIsFalse(h map[string]interface{}, key string) bool {
	if h == nil {
		return false
	}

	v, ok := h[key]
	if !ok || v == nil {
		return false
	}

	if s, isStr := v.(string); isStr && s == "" {
		return false
	}

	return !asBool(v)
}

// indexUptimeGainByNode folds the com (tnode + cnode) and ctl (anode) series
// into one node_id -> seconds-gained map. A node absent from the map gained
// nothing, the same as one whose counter stalled: both read as downtime.
func indexUptimeGainByNode(sets ...[]map[string]interface{}) map[string]float64 {
	out := map[string]float64{}

	for _, set := range sets {
		for _, row := range set {
			nodeID := str(row["node_id"])
			if nodeID == "" {
				continue // a series without its node label is unattributable
			}

			// A node belongs to exactly one of the two series; if it somehow
			// appears in both, the larger gain counts.
			if gain := uptimeGain(row["value"]); gain > out[nodeID] {
				out[nodeID] = gain
			}
		}
	}

	return out
}

func indexHealthByNode(health []map[string]interface{}) map[string]map[string]interface{} {
	out := make(map[string]map[string]interface{}, len(health))
	for _, h := range health {
		out[str(h["node_id"])] = h
	}

	return out
}

// uptimeGain reads the seconds of uptime a node gained over the window out of
// a Prometheus [unix_ts, "value"] sample pair:
//
//	[1787257964.364, "59.87"] -> up for the whole 60s window
//	[1787257964.364, "0"]     -> counter stalled: down
//	[1787257964.364, "NaN"]   -> no usable reading: down
//
// Anything unparseable, infinite or negative reads as no gain at all.
func uptimeGain(v interface{}) float64 {
	seconds, ok := uptimeSeconds(v)
	if !ok || math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds < 0 {
		return 0
	}

	return seconds
}

// uptimeSeconds pulls the numeric half out of a Prometheus sample pair. Be
// liberal on the way in: a bare number or numeric string is accepted too.
func uptimeSeconds(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case []interface{}:
		if len(t) == 0 {
			return 0, false
		}

		return uptimeSeconds(t[len(t)-1])
	case float64:
		return t, true
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0, false
		}

		return f, true
	default:
		return 0, false
	}
}

// uptimeReported reports whether a mapped uptime sample is a real reading.
//
// /v1/last serialises a Prometheus instant vector, so `value` is the
// [unix_ts, "value-as-string"] pair. A node that stopped pushing keeps its
// series — Prometheus turns it into a staleness marker, the sanitizer
// forwards it untouched, and it comes back as the string "NaN":
//
//	"value": [1787257964.364, "NaN"]     -> not reporting
//	"value": [1787257964.364, "543529.39"] -> reporting
//
// Anything unparseable, infinite or negative is treated as not reporting: an
// uptime that is not a finite non-negative number of seconds is not evidence
// the node is alive.
func uptimeReported(v interface{}) bool {
	seconds, ok := uptimeSeconds(v)

	return ok && !math.IsNaN(seconds) && !math.IsInf(seconds, 0) && seconds >= 0
}

// uptimeSeconds pulls the numeric half out of a Prometheus sample pair. Be
// liberal on the way in: a bare number or numeric string is accepted too.
func uptimeSeconds(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case []interface{}:
		if len(t) == 0 {
			return 0, false
		}

		return uptimeSeconds(t[len(t)-1])
	case float64:
		return t, true
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0, false
		}

		return f, true
	default:
		return 0, false
	}
}
