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

// The node types that make up a site and are judged by the uptime KPIs. Any
// other type attached to a site (hnode) is ignored.
const (
	nodeTypeTower      = "tnode"
	nodeTypeAmplifier  = "anode"
	nodeTypeController = "cnode"
)

// SiteUptime (SITE_UPTIME @ scope network_id+site_id), percentage.
//
// Uptime is 100% by default and comes down only on evidence of failure: a
// window counts against a site when something in it positively reports a
// problem, never when information is merely absent.
//
// The two kinds of evidence:
//
//	LIVENESS — each node pushes a system uptime counter, which the sanitizer
//	           republishes with its node_id/site/network stamped on it
//	           (com_node_uptime for the tnode and the cnode, ctl_node_uptime
//	           for the anode). When a node dies Prometheus turns its series
//	           into a staleness marker, which arrives as the literal NaN.
//	SERVICE + RADIO — interfaces.cellular.available (service; tnode-only, as
//	           the anode reports "cellular": null) and
//	           interfaces.radio.available, from the node-gateway health report.
//
// Per node:
//
//	DOWN ⟺ its uptime reading is NaN
//	       OR a health flag it reports is explicitly false
//	UP   otherwise, including when nothing has reported anything about it
//
// Per site:
//
//	DOWN — any node that makes up the site is down
//	UP   — otherwise, including a site carrying no tnode/anode/cnode
//
// Node lifecycle state and node-gateway reachability are not read: a node is a
// separate device, and its own uptime counter is the authority on whether it
// is alive. radio.state is not read either — uptime tracks interface
// AVAILABILITY, so a radio that is available but switched off is still up, and
// there is no planned/unplanned distinction.
//
// Liveness rests on the staleness marker arriving. A node whose series stops
// being written without Prometheus marking it stale reads up until a NaN
// appears or the series ages past the pull's lookback.
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
//	nodes      — registry.node.list: site membership only
//	health     — node.health.interfaces (state): cellular_available,
//	             radio_available
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

// siteAgg accumulates one site's state for a window. Only downNodes decides
// the outcome: uptime is 100% until something is reported down.
type siteAgg struct {
	networkID string
	downNodes int
}

func (a *siteAgg) isUp() bool { return a.downNodes == 0 }

// classifySites builds the per-window site state shared by SITE_UPTIME and
// NETWORK_UPTIME from the nodes + uptime + health inputs.
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
	uptimeByNode := indexUptimeByNode(com, ctl)

	sites := map[string]*siteAgg{}

	for _, node := range nodes {
		siteID, networkID := str(node["site_id"]), str(node["network_id"])
		if siteID == "" || networkID == "" {
			continue // unattached node
		}

		// Every node type registers the site, so a site made only of hnodes
		// still produces a row.
		agg, ok := sites[siteID]
		if !ok {
			agg = &siteAgg{networkID: networkID}
			sites[siteID] = agg
		}

		nodeType := strings.ToLower(str(node["type"]))
		if !isSiteNodeType(nodeType) {
			continue // not part of what makes a site
		}

		nodeID := str(node["node_id"])

		if nodeIsDown(nodeType, healthByNode[nodeID], uptimeByNode[nodeID]) {
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

// nodeIsDown reports whether this window carries positive evidence that the
// node is not serving. Everything else — including knowing nothing at all
// about it — leaves the node up.
func nodeIsDown(nodeType string, h map[string]interface{}, u uptimeReading) bool {
	// The node's own counter went stale: it is dead.
	if u.present && !u.alive {
		return true
	}

	// The radio is reported, and reported gone.
	if flagIsFalse(h, "radio_available") {
		return true
	}

	// Cellular is the service interface, and it is tnode-only: the anode
	// reports "cellular": null, so the field never reaches the row for it.
	if nodeType == nodeTypeTower && flagIsFalse(h, "cellular_available") {
		return true
	}

	return false
}

// flagIsFalse reports whether a health flag is present AND says false. A
// missing field is not a false one: the ingest mapper writes a field only when
// its path resolves, so an unreachable probe records the bare
// `unreachable: true` marker with no interface flags at all, and the anode's
// null cellular object yields no cellular field.
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

// uptimeReading is one node's liveness evidence for a window. The two fields
// answer different questions: `present` says whether the node has a series at
// all, `alive` says whether that series carries a real reading. A dead node is
// present and not alive (its staleness marker arrives as NaN); a node nothing
// has been heard from yet is simply not present.
type uptimeReading struct {
	present bool
	alive   bool
}

// indexUptimeByNode folds the com (tnode + cnode) and ctl (anode) uptime
// series into one node_id -> reading map.
func indexUptimeByNode(sets ...[]map[string]interface{}) map[string]uptimeReading {
	out := map[string]uptimeReading{}

	for _, set := range sets {
		for _, row := range set {
			nodeID := str(row["node_id"])
			if nodeID == "" {
				continue // a series without its node label is unattributable
			}

			// A node should appear in exactly one of the two series; if it
			// somehow appears in both, any live reading counts.
			prev := out[nodeID]
			out[nodeID] = uptimeReading{
				present: true,
				alive:   prev.alive || uptimeReported(row["value"]),
			}
		}
	}

	return out
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

// uptimeReported reports whether a mapped uptime sample is a real reading.
//
// /v1/last serialises a Prometheus instant vector, so `value` is the
// [unix_ts, "value-as-string"] pair:
//
//	"value": [1787257964.364, "NaN"]       -> dead
//	"value": [1787257964.364, "543529.39"] -> alive
//
// Anything unparseable, infinite or negative reads as dead: a series that
// exists but whose value is not a finite non-negative number of seconds is not
// evidence the node is alive.
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
