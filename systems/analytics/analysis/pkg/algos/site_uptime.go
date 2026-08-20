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

	"github.com/ukama/ukama/systems/analytics/schema"
)

// Node types that make up a site and are judged by the uptime KPIs. Any other
// type attached to a site (hnode) is ignored rather than allowed to drag the
// site down — the same rule SITES_ONLINE uses.
const (
	nodeTypeTower      = "tnode"
	nodeTypeAmplifier  = "anode"
	nodeTypeController = "cnode"
)

// SiteUptime (SITE_UPTIME @ scope network_id+site_id), percentage.
//
// v3 judges a node on TWO independent signals, and it has to pass both:
//
//	LIVENESS — the node is powering the metrics path. The sanitizer
//	           republishes each node's system uptime counter with its
//	           node_id/site/network stamped on it (com_node_uptime for the
//	           tnode and the cnode, ctl_node_uptime for the anode), so a
//	           reading exists for every node type. A node that stopped
//	           pushing does not vanish — Prometheus turns the series into a
//	           staleness marker, which arrives as the literal NaN, and that
//	           is what reads as down here. No sample at all (the series aged
//	           out of the pull's lookback) is down for the same reason.
//	SERVICE + RADIO — the interfaces the site actually serves with, from the
//	           node-gateway health probe: service is the tnode's
//	           interfaces.cellular.available (the anode reports
//	           "cellular": null and carries no service), radio is
//	           interfaces.radio.available on tnode and anode alike.
//
// Per-window node state:
//
//	tnode UP: uptime reported && cellular.available && radio.available
//	anode UP: uptime reported && radio.available
//	cnode UP: uptime reported          (it is not probed: node.health.
//	          interfaces only fans out over tnode/anode, and the cnode
//	          carries neither a cellular nor a radio interface)
//	DOWN:     a NaN/absent uptime reading, either flag false, an unreachable
//	          probe, or no health row at all
//
// radio.state is deliberately NOT read. Uptime tracks interface
// AVAILABILITY — a radio that is available but switched off is still up, and
// there is no planned/unplanned distinction to carry.
//
// Per-window site state:
//
//	UP   — the site has at least one tnode/anode/cnode and every one of them
//	       is up
//	DOWN — any of them is down, or the site has none of the three
//
// v2 judged only the tnode and the anode, and only on the health probe: a
// node whose gateway probe answered from a stale cache read up while the box
// itself was dark, and a site was never held to its cnode. v3 requires
// liveness on every site node type, matching what sites_online@v3 calls a
// site.
//
// Every window counts: Sum = 100 when up else 0, Count = 1 always. The
// aggregator's weighted AVG is then exactly
//
//	uptime % = up_windows / total_windows x 100
//
// at daily/weekly/monthly. Query with op=AVG for the timeseries.
//
// Inputs:
//
//	nodes      — registry.node.list: site membership (EVERY node type, so a
//	             site is never invisible) and the expected node set. Registry
//	             stays authoritative for membership rather than the metric's
//	             own site label: a site whose nodes are ALL dark still has to
//	             produce a row, and a dark node has no metric to carry a
//	             label on.
//	health     — node.health.interfaces (state): cellular_available,
//	             radio_available, unreachable
//	com_uptime — metrics.com_uptime.last: tnode + cnode liveness
//	ctl_uptime — metrics.ctl_uptime.last: anode liveness
func SiteUptime(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, err := classifySites(in, "SITE_UPTIME")
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(sites))

	for siteID, agg := range sites {
		scope := map[string]string{"network_id": agg.networkID, "site_id": siteID}

		value := 0.0
		if agg.isUp() {
			value = 100
		}

		results = append(results, Result{
			Scope: scope,
			Value: value,
			Sum:   value,
			Count: 1,
			Min:   value,
			Max:   value,
		})
	}

	return results, nil
}

// NetworkUptime (NETWORK_UPTIME @ scope network_id), percentage: the uptime of
// all the sites under the network.
//
//	network uptime % = sites_up / sites_total x 100
//
// Per window each network emits Sum = 100 x sites_up and Count = sites_total,
// so the aggregator's weighted AVG equals the pooled formula exactly at every
// span — NOT an average of per-site percentages. Pooling this way weights each
// site by the site-time it actually lost: one site down for a day and
// twenty-four sites down for an hour each cost the same.
//
// Inputs and per-site rule: identical to SITE_UPTIME.
func NetworkUptime(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, err := classifySites(in, "NETWORK_UPTIME")
	if err != nil {
		return nil, err
	}

	type netAgg struct{ up, total float64 }

	networks := map[string]*netAgg{}

	for _, agg := range sites {
		n, ok := networks[agg.networkID]
		if !ok {
			n = &netAgg{}
			networks[agg.networkID] = n
		}

		n.total++

		if agg.isUp() {
			n.up++
		}
	}

	results := make([]Result, 0, len(networks))

	for networkID, n := range networks {
		value := 0.0
		if n.total > 0 {
			value = n.up / n.total * 100
		}

		results = append(results, Result{
			Scope: map[string]string{"network_id": networkID},
			Value: value,
			Sum:   n.up * 100,
			Count: n.total,
			Min:   value,
			Max:   value,
		})
	}

	return results, nil
}

// siteAgg accumulates one site's node states for a window. judged counts the
// nodes that make up a site (tnode/anode/cnode); a site with none of them has
// nothing to serve with.
type siteAgg struct {
	networkID string
	judged    int
	upNodes   int
	downNodes int
}

// isUp: the site has something to serve with, and none of it is down.
func (a *siteAgg) isUp() bool {
	return a.judged > 0 && a.downNodes == 0
}

// classifySites builds the per-window site state shared by SITE_UPTIME and
// NETWORK_UPTIME from the health + uptime + nodes inputs.
func classifySites(in Datasets, kpi string) (map[string]*siteAgg, error) {
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
	reportingByNode, uptimeKnown := indexUptimeByNode(com, ctl)

	sites := map[string]*siteAgg{}

	for _, node := range nodes {
		siteID, networkID := str(node["site_id"]), str(node["network_id"])
		if siteID == "" || networkID == "" {
			continue // unattached node
		}

		// Every node type registers the site, so a site made only of hnodes
		// still produces a row (classified down) instead of vanishing.
		agg, ok := sites[siteID]
		if !ok {
			agg = &siteAgg{networkID: networkID}
			sites[siteID] = agg
		}

		nodeType := strings.ToLower(str(node["type"]))
		if !isSiteNodeType(nodeType) {
			continue // not part of what makes a site
		}

		agg.judged++

		nodeID := str(node["node_id"])

		if nodeIsUp(nodeType, healthByNode[nodeID], reportingByNode[nodeID], uptimeKnown) {
			agg.upNodes++
		} else {
			agg.downNodes++
		}
	}

	return sites, nil
}

// isSiteNodeType reports whether a node type is one of the three that make up
// a site and is therefore judged by the uptime KPIs.
func isSiteNodeType(nodeType string) bool {
	return nodeType == nodeTypeTower ||
		nodeType == nodeTypeAmplifier ||
		nodeType == nodeTypeController
}

// indexHealthByNode keys the node.health.interfaces rows of a window by node
// id.
func indexHealthByNode(health []map[string]interface{}) map[string]map[string]interface{} {
	out := make(map[string]map[string]interface{}, len(health))
	for _, h := range health {
		out[str(h["node_id"])] = h
	}

	return out
}

// indexUptimeByNode folds the com (tnode + cnode) and ctl (anode) uptime
// series into one node_id -> is-reporting map.
//
// The second return value says whether the uptime datasets carried ANY row at
// all. Empty means the metrics path delivered nothing this window — which is
// a statement about the pipeline, not about the nodes — so the caller drops
// the liveness term rather than reporting every site in the org as down. A
// node that is genuinely dead is NOT absent: it arrives with a NaN reading.
func indexUptimeByNode(sets ...[]map[string]interface{}) (map[string]bool, bool) {
	out := map[string]bool{}
	seen := false

	for _, set := range sets {
		for _, row := range set {
			seen = true

			nodeID := str(row["node_id"])
			if nodeID == "" {
				continue // a series without its node label is unattributable
			}

			// A node should appear in exactly one of the two series; if it
			// somehow appears in both, any live reading counts.
			out[nodeID] = out[nodeID] || uptimeReported(row["value"])
		}
	}

	return out, seen
}

// nodeIsUp judges one node from the signals its type actually has.
//
// uptimeKnown = false means the uptime datasets were empty for the window
// (see indexUptimeByNode): the liveness term is dropped and the rule degrades
// to the v2 health-only one, with the cnode left unjudged because nothing
// probes it.
func nodeIsUp(nodeType string, h map[string]interface{}, reporting, uptimeKnown bool) bool {
	if uptimeKnown && !reporting {
		return false
	}

	if nodeType == nodeTypeController {
		// The cnode carries neither a cellular nor a radio interface and is
		// not probed by node.health.interfaces, so liveness is its only
		// signal — already checked above.
		return true
	}

	// A missing health row (never probed) and an unreachable probe both count
	// as down: neither can confirm the interfaces are there.
	if h == nil || asBool(h["unreachable"]) {
		return false
	}

	if !asBool(h["radio_available"]) {
		return false
	}

	// Cellular is the service interface, and it is tnode-only.
	if nodeType == nodeTypeTower && !asBool(h["cellular_available"]) {
		return false
	}

	return true
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
