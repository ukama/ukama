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
// Formula: Site Uptime % = U / (U + (D - I)) x 100
//   U = windows where the site's service is up
//   D = windows where it is down
//   I = intentionally-off windows (operator switched the radio off) —
//       excluded from the denominator
//
// Per-window node classification (cnodes are excluded; only tnode/anode):
//   tnode UP:  cellular.available && radio.available && radio.state == on
//   anode UP:  radio.available && radio.state == on
//   INTENTIONAL: radio.available && radio.state == off (radio works but is
//       switched off; for tnodes this overrides cellular state)
//   DOWN: everything else — radio.available == false, cellular fault on a
//       tnode, or the health probe unreachable (unreachable == true)
//
// Per-window site classification:
//   UP          — every relevant node is UP
//   INTENTIONAL — not up, and every non-up node is intentionally off
//   DOWN        — any node in unplanned down
//
// Emitted components make the aggregator's weighted AVG equal the formula
// exactly at every span: Sum = 100 when up else 0; Count = 1 for up/down
// windows, 0 for intentional windows (excluded from the denominator).
// Query with op=AVG and span daily/weekly/monthly.
//
// Inputs:
//   health — node.health.interfaces (state): cellular_available,
//            radio_available, radio_state, unreachable + node lineage
//   nodes  — registry.node.list: the expected tnode/anode set per site; a
//            node with no health row yet counts as down.
func SiteUptime(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, err := classifySites(in, "SITE_UPTIME")
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(sites))

	for siteID, agg := range sites {
		scope := map[string]string{"network_id": agg.networkID, "site_id": siteID}

		switch agg.state() {
		case stateUp:
			results = append(results, Result{Scope: scope, Value: 100, Sum: 100, Count: 1, Min: 100, Max: 100})
		case stateDown:
			results = append(results, Result{Scope: scope, Value: 0, Sum: 0, Count: 1, Min: 0, Max: 0})
		default: // intentionally off only: excluded from the denominator
			results = append(results, Result{Scope: scope, Value: 0, Sum: 0, Count: 0, Min: 0, Max: 0})
		}
	}

	return results, nil
}

// NetworkUptime (NETWORK_UPTIME @ scope network_id), percentage. Pools site
// time-units across the network per the formula
//
//	Network Uptime % = SUM(Ui) / (SUM(Ui) + SUM(Di - Ii)) x 100
//
// Per window each network emits Sum = 100 x sites_up and Count = sites_up +
// sites_unplanned_down (intentionally-off sites excluded), so the
// aggregator's weighted AVG equals the pooled formula exactly at every span
// — NOT an average of per-site percentages.
func NetworkUptime(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, err := classifySites(in, "NETWORK_UPTIME")
	if err != nil {
		return nil, err
	}

	type netAgg struct{ up, down float64 }

	networks := map[string]*netAgg{}

	for _, agg := range sites {
		n, ok := networks[agg.networkID]
		if !ok {
			n = &netAgg{}
			networks[agg.networkID] = n
		}

		switch agg.state() {
		case stateUp:
			n.up++
		case stateDown:
			n.down++
		} // intentional sites join neither numerator nor denominator
	}

	results := make([]Result, 0, len(networks))

	for networkID, n := range networks {
		scope := map[string]string{"network_id": networkID}
		denominator := n.up + n.down

		value := 0.0
		if denominator > 0 {
			value = n.up / denominator * 100
		}

		results = append(results, Result{
			Scope: scope,
			Value: value,
			Sum:   n.up * 100,
			Count: denominator, // 0 when all sites intentionally off -> window excluded
			Min:   value,
			Max:   value,
		})
	}

	return results, nil
}

// siteAgg accumulates one site's node states for a window.
type siteAgg struct {
	networkID   string
	up          bool
	intentional bool
	down        bool
}

func (a *siteAgg) state() int {
	switch {
	case a.up:
		return stateUp
	case a.down:
		return stateDown
	default:
		return stateIntentional
	}
}

// classifySites builds the per-window site classification shared by
// SITE_UPTIME and NETWORK_UPTIME from the health + nodes inputs.
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
		nodeType := strings.ToLower(str(node["type"]))
		if nodeType != "tnode" && nodeType != "anode" {
			continue // cnodes excluded from uptime
		}

		siteID, networkID := str(node["site_id"]), str(node["network_id"])
		if siteID == "" || networkID == "" {
			continue // unattached node
		}

		agg, ok := sites[siteID]
		if !ok {
			agg = &siteAgg{networkID: networkID, up: true}
			sites[siteID] = agg
		}

		switch classifyNode(nodeType, healthByNode[str(node["node_id"])]) {
		case stateUp:
			// keeps agg.up
		case stateIntentional:
			agg.up = false
			agg.intentional = true
		default:
			agg.up = false
			agg.down = true
		}
	}

	return sites, nil
}

// indexHealthByNode keys the node.health.interfaces rows of a window by node
// id. Shared by the uptime algos and SITES_DEGRADED.
func indexHealthByNode(health []map[string]interface{}) map[string]map[string]interface{} {
	out := make(map[string]map[string]interface{}, len(health))
	for _, h := range health {
		out[str(h["node_id"])] = h
	}

	return out
}

// Node/window classification states.
const (
	stateUp = iota
	stateIntentional
	stateDown
)

// classifyNode implements the per-node rules. A missing health row (probe
// never ran) and unreachable probes both count as down.
func classifyNode(nodeType string, h map[string]interface{}) int {
	if h == nil || asBool(h["unreachable"]) {
		return stateDown
	}

	radioAvailable := asBool(h["radio_available"])
	radioOn := strings.EqualFold(str(h["radio_state"]), "on")

	// Radio works but is switched off => operator intent.
	if radioAvailable && !radioOn {
		return stateIntentional
	}

	if !radioAvailable {
		return stateDown
	}

	// Radio available and on. Tnodes additionally require cellular.
	if nodeType == "tnode" && !asBool(h["cellular_available"]) {
		return stateDown
	}

	return stateUp
}
