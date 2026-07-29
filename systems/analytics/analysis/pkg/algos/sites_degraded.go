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

// SitesDegraded (SITES_DEGRADED @ scope network_id), definition v2: sites per
// network carrying at least one fault. v1 only looked at registry
// connectivity, which missed a node that is reachable but whose radio or
// cellular interface is down — the site keeps reporting "fine" while it
// cannot actually serve.
//
// A site is degraded when ANY of its nodes trips one of these:
//
//  1. registry connectivity is not online  — any node type, including cnode
//  2. tnode health: cellular.available == false
//  3. tnode or anode health: radio.available == false
//  4. tnode or anode health probe unreachable, or the node has never reported
//
// Rules 2-4 read node.health.interfaces (the /v1/health/interfaces payload,
// already ingested per tnode/anode every window) — cellular is a tnode-only
// interface, the anode reports "cellular": null, so it is only asserted for
// tnodes. Radio is asserted for both.
//
// Deliberately NOT degraded: a radio that is available but switched off
// (radio.state == "off"). That is operator intent, not a fault, and it
// matches SITE_UPTIME, which excludes intentionally-off windows from
// downtime rather than counting them against the site.
//
// Inputs:
//
//	nodes    — registry.node.list (node_id, site_id, network_id, type,
//	           connectivity): the site membership and registry connectivity
//	health   — node.health.interfaces (cellular_available, radio_available,
//	           radio_state, unreachable) keyed by node_id
//	networks — registry.network.getAll (network_id) — zero-fill, so a network
//	           with no degraded site reads 0 rather than "—"
func SitesDegraded(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, networks, err := groupSites(in, "SITES_DEGRADED")
	if err != nil {
		return nil, err
	}

	health, ok := in["health"]
	if !ok {
		return nil, fmt.Errorf("SITES_DEGRADED: missing input 'health'")
	}

	healthByNode := indexHealthByNode(health)

	results := make([]Result, 0, len(networks))

	for _, network := range networks {
		networkID := str(network["network_id"])
		if networkID == "" {
			continue
		}

		degraded := 0

		for _, nodes := range sites[networkID] {
			if siteIsDegraded(nodes, healthByNode) {
				degraded++
			}
		}

		results = append(results, CountResult(
			map[string]string{"network_id": networkID},
			float64(degraded),
		))
	}

	return results, nil
}

// siteIsDegraded reports whether any node of the site carries a fault. First
// match wins — the KPI counts degraded sites, not faults per site.
func siteIsDegraded(nodes []map[string]interface{}, healthByNode map[string]map[string]interface{}) bool {
	for _, node := range nodes {
		// 1. Registry connectivity, every node type.
		if !isOnline(node["connectivity"]) {
			return true
		}

		// 2-4 are health-report rules. Only tnode/anode are probed
		// (node.yaml filters the fan-out), so a cnode/hnode legitimately has
		// no health row and must not be judged on one.
		nodeType := strings.ToLower(str(node["type"]))
		if nodeType != "tnode" && nodeType != "anode" {
			continue
		}

		h, ok := healthByNode[str(node["node_id"])]
		if !ok || asBool(h["unreachable"]) {
			// Never probed, or the probe could not reach the node: the
			// interfaces cannot be confirmed healthy. SITE_UPTIME already
			// treats both as down, so treating them as degraded keeps the two
			// KPIs telling the same story.
			return true
		}

		// Radio is present on both node types.
		if !asBool(h["radio_available"]) {
			return true
		}

		// Cellular is tnode-only (anode reports "cellular": null, which the
		// ingest mapper leaves absent rather than false).
		if nodeType == "tnode" && !asBool(h["cellular_available"]) {
			return true
		}
	}

	return false
}
