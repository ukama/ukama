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

// Node types that make up a site; any other type (hnode) is ignored.
const (
	nodeTypeTower      = "tnode"
	nodeTypeAmplifier  = "anode"
	nodeTypeController = "cnode"
)

// SiteUptime (SITE_UPTIME @ network_id+site_id): one availability check per
// window.
//
//	node = DOWN if a health flag is reported false or its uptime counter did
//	       not advance during the window; UP otherwise
//	site = 100 if every tnode/anode/cnode is UP, else 0
//
// Emits Sum = 100 or 0, Count = 1: the aggregator's AVG is the share of
// windows the site was up for.
func SiteUptime(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, err := classifySites(in, "SITE_UPTIME")
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

// NetworkUptime (NETWORK_UPTIME @ network_id): sites up / site count x 100,
// with the same per-site check as SITE_UPTIME.
//
// Emits Sum = total of the sites' values, Count = site count, so the
// aggregator's AVG is weighted by (site, window) pairs, not an average of
// per-site averages.
func NetworkUptime(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, err := classifySites(in, "NETWORK_UPTIME")
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

// siteAgg is one site's check for the window: starts at 100 and any node
// reported down takes it to 0.
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

// classifySites runs the per-site check shared by SITE_UPTIME and
// NETWORK_UPTIME.
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
	gainByNode := indexUptimeGainByNode(com, ctl)

	sites := map[string]*siteAgg{}

	// Site membership comes from the registry, not from metric labels: a dark
	// node has no metric to carry a label.
	for _, node := range nodes {
		siteID, networkID := str(node["site_id"]), str(node["network_id"])
		if siteID == "" || networkID == "" {
			continue // unattached node
		}

		// Every node type registers the site, so an hnode-only site still
		// emits a row.
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
		gain, seen := gainByNode[nodeID]

		agg.observe(nodeUptimePercent(nodeType, healthByNode[nodeID], gain, seen))
	}

	return sites, nil
}

func isSiteNodeType(nodeType string) bool {
	return nodeType == nodeTypeTower ||
		nodeType == nodeTypeAmplifier ||
		nodeType == nodeTypeController
}

// nodeUptimePercent is one node's check for the window: 100 (up) or 0 (down).
//
// A health flag reported false is down. Otherwise the node is up if its
// uptime counter advanced (gain > 0) or if it has no series in the window at
// all: uptime is 100 by default and only comes down on evidence, and a series
// that has not appeared yet is not evidence. A series present with no gain is
// a stalled counter, which is.
//
// radio.state, node lifecycle state and node-gateway reachability are not
// read: the node's own counter is the authority.
func nodeUptimePercent(nodeType string, h map[string]interface{},
	gain float64, seen bool) float64 {
	if flagIsFalse(h, "radio_available") {
		return 0
	}

	// Cellular is tnode-only: the anode reports "cellular": null.
	if nodeType == nodeTypeTower && flagIsFalse(h, "cellular_available") {
		return 0
	}

	if !seen || gain > 0 {
		return 100
	}

	return 0
}

// flagIsFalse reports whether a health flag is present and false. A missing
// field is not a false one: an unreachable probe records no flags at all.
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

// indexUptimeGainByNode folds the com (tnode, cnode) and ctl (anode) series
// into one node_id -> seconds-gained map. Membership means the node had a
// series in the window; absence means it reported nothing.
func indexUptimeGainByNode(sets ...[]map[string]interface{}) map[string]float64 {
	out := map[string]float64{}

	for _, set := range sets {
		for _, row := range set {
			nodeID := str(row["node_id"])
			if nodeID == "" {
				continue
			}

			// A node in both series keeps the larger gain.
			gain := uptimeGain(row["value"])
			if prev, ok := out[nodeID]; !ok || gain > prev {
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

// uptimeGain reads the seconds gained out of a Prometheus [unix_ts, "value"]
// sample pair. NaN, infinite, negative or unparseable values read as 0.
func uptimeGain(v interface{}) float64 {
	seconds, ok := uptimeSeconds(v)
	if !ok || math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds < 0 {
		return 0
	}

	return seconds
}

// uptimeSeconds pulls the numeric half out of a sample pair; a bare number or
// numeric string is accepted too.
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
