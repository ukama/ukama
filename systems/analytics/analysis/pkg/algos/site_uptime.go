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
	"strings"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// SiteUptime (SITE_UPTIME @ scope network_id+site_id), percentage.
//
// A site is up in a window when its SERVICE and its RADIO are available:
//
//	service — the digital cellular service, reported by the tnode as
//	          interfaces.cellular.available. The anode has no cellular
//	          interface (its payload carries "cellular": null), so service is
//	          only asserted for tnodes.
//	radio   — interfaces.radio.available, reported by tnode and anode alike.
//
// Per-window node state (only tnode/anode carry service and radio; cnodes are
// not probed and are not judged here):
//
//	tnode UP: cellular.available && radio.available
//	anode UP: radio.available
//	DOWN:     either flag false, an unreachable probe, or no health row at all
//
// radio.state is deliberately NOT read. Uptime tracks interface
// AVAILABILITY — a radio that is available but switched off is still up, and
// there is no planned/unplanned distinction to carry.
//
// Per-window site state:
//
//	UP   — the site has at least one tnode/anode and every one of them is up
//	DOWN — any of them is down, or the site has no tnode/anode at all
//	       (nothing to carry service or radio, so it cannot be serving)
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
//	health — node.health.interfaces (state): cellular_available,
//	         radio_available, unreachable + node lineage
//	nodes  — registry.node.list: site membership (EVERY node type, so a site
//	         is never invisible) and the expected tnode/anode set
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
// Inputs: identical to SITE_UPTIME.
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

// siteAgg accumulates one site's node states for a window. bearers counts the
// nodes that actually carry service/radio (tnode/anode); a site with none
// cannot be serving.
type siteAgg struct {
	networkID string
	bearers   int
	upNodes   int
	downNodes int
}

// isUp: the site has something to serve with, and none of it is down.
func (a *siteAgg) isUp() bool {
	return a.bearers > 0 && a.downNodes == 0
}

// classifySites builds the per-window site state shared by SITE_UPTIME and
// NETWORK_UPTIME from the health + nodes inputs.
func classifySites(in Datasets, kpi string) (map[string]*siteAgg, error) {
	health, ok := in["health"]
	if !ok {
		return nil, fmt.Errorf("%s: missing input 'health'", kpi)
	}

	nodes, ok := in["nodes"]
	if !ok {
		return nil, fmt.Errorf("%s: missing input 'nodes'", kpi)
	}

	healthByNode := indexHealthByNode(health)

	sites := map[string]*siteAgg{}

	for _, node := range nodes {
		siteID, networkID := str(node["site_id"]), str(node["network_id"])
		if siteID == "" || networkID == "" {
			continue // unattached node
		}

		// Every node type registers the site, so a site made only of cnodes
		// still produces a row (classified down) instead of vanishing.
		agg, ok := sites[siteID]
		if !ok {
			agg = &siteAgg{networkID: networkID}
			sites[siteID] = agg
		}

		nodeType := strings.ToLower(str(node["type"]))
		if nodeType != "tnode" && nodeType != "anode" {
			continue // carries neither service nor radio
		}

		agg.bearers++

		if nodeIsUp(nodeType, healthByNode[str(node["node_id"])]) {
			agg.upNodes++
		} else {
			agg.downNodes++
		}
	}

	return sites, nil
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

// nodeIsUp reports whether a node's service and radio interfaces are
// available. A missing health row (never probed) and an unreachable probe both
// count as down — neither can confirm the interfaces are there.
func nodeIsUp(nodeType string, h map[string]interface{}) bool {
	if h == nil || asBool(h["unreachable"]) {
		return false
	}

	if !asBool(h["radio_available"]) {
		return false
	}

	// Cellular is the service interface, and it is tnode-only.
	if nodeType == "tnode" && !asBool(h["cellular_available"]) {
		return false
	}

	return true
}
