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

// siteNodes indexes nodes per network per site.
type siteNodes map[string]map[string][]map[string]interface{} // network -> site -> nodes

// SitesOnline (SITES_ONLINE @ scope network_id), definition v2: a site is
// online when its cnode (control/connectivity node) is online. Sites without
// a cnode are not online.
//
// Inputs:
//   nodes    — registry.node.list (node_id, site_id, network_id, type,
//              connectivity)
//   networks — registry.network.getAll (network_id) — zero-fill.
func SitesOnline(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, networks, err := groupSites(in, "SITES_ONLINE")
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(networks))

	for _, network := range networks {
		networkID := str(network["network_id"])
		if networkID == "" {
			continue
		}

		online := 0

		for _, nodes := range sites[networkID] {
			if siteIsOnline(nodes) {
				online++
			}
		}

		results = append(results, CountResult(
			map[string]string{"network_id": networkID},
			float64(online),
		))
	}

	return results, nil
}

// SitesDegraded (SITES_DEGRADED @ scope network_id): sites with at least one
// offline node.
//
// Inputs: same as SITES_ONLINE.
func SitesDegraded(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sites, networks, err := groupSites(in, "SITES_DEGRADED")
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(networks))

	for _, network := range networks {
		networkID := str(network["network_id"])
		if networkID == "" {
			continue
		}

		degraded := 0

		for _, nodes := range sites[networkID] {
			if siteHasOfflineNode(nodes) {
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

func groupSites(in Datasets, kpi string) (siteNodes, []map[string]interface{}, error) {
	nodes, ok := in["nodes"]
	if !ok {
		return nil, nil, fmt.Errorf("%s: missing input 'nodes'", kpi)
	}

	networks, ok := in["networks"]
	if !ok {
		return nil, nil, fmt.Errorf("%s: missing input 'networks'", kpi)
	}

	sites := siteNodes{}

	for _, node := range nodes {
		networkID := str(node["network_id"])
		siteID := str(node["site_id"])

		if networkID == "" || siteID == "" {
			continue // node not attached to a site yet
		}

		if sites[networkID] == nil {
			sites[networkID] = map[string][]map[string]interface{}{}
		}

		sites[networkID][siteID] = append(sites[networkID][siteID], node)
	}

	return sites, networks, nil
}

// siteIsOnline: the site's cnode must be online (definition v2).
func siteIsOnline(nodes []map[string]interface{}) bool {
	for _, node := range nodes {
		if strings.EqualFold(str(node["type"]), "cnode") && isOnline(node["connectivity"]) {
			return true
		}
	}

	return false
}

// siteHasOfflineNode: any node of the site is not online.
func siteHasOfflineNode(nodes []map[string]interface{}) bool {
	for _, node := range nodes {
		if !isOnline(node["connectivity"]) {
			return true
		}
	}

	return false
}

// isOnline handles both serializations of the connectivity enum: the string
// name ("Online") and the numeric ukama.NodeConnectivity value (1 = online,
// emitted when protobuf enums pass through encoding/json).
func isOnline(v interface{}) bool {
	switch strings.ToLower(str(v)) {
	case "online", "1":
		return true
	default:
		return false
	}
}

func str(v interface{}) string {
	if v == nil {
		return ""
	}

	if s, ok := v.(string); ok {
		return s
	}

	return fmt.Sprintf("%v", v)
}
