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

// SitesOnline (SITES_ONLINE @ scope network_id): number of sites per network
// with at least one node whose connectivity is "online".
//
// Inputs:
//   nodes    — registry.node.list (node_id, site_id, network_id, connectivity)
//   networks — registry.network.getAll (network_id) — ensures networks with
//              no online sites still emit a 0 value.
func SitesOnline(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	nodes, ok := in["nodes"]
	if !ok {
		return nil, fmt.Errorf("SITES_ONLINE: missing input 'nodes'")
	}

	networks, ok := in["networks"]
	if !ok {
		return nil, fmt.Errorf("SITES_ONLINE: missing input 'networks'")
	}

	onlineSites := map[string]map[string]bool{} // network -> set of online sites

	for _, node := range nodes {
		networkID := str(node["network_id"])
		siteID := str(node["site_id"])

		if networkID == "" || siteID == "" {
			continue // node not attached to a site yet
		}

		if isOnline(node["connectivity"]) {
			if onlineSites[networkID] == nil {
				onlineSites[networkID] = map[string]bool{}
			}

			onlineSites[networkID][siteID] = true
		}
	}

	results := make([]Result, 0, len(networks))

	for _, network := range networks {
		networkID := str(network["network_id"])
		if networkID == "" {
			continue
		}

		results = append(results, CountResult(
			map[string]string{"network_id": networkID},
			float64(len(onlineSites[networkID])),
		))
	}

	return results, nil
}

// isOnline handles both serializations of the connectivity enum: the string
// name ("online") and the numeric ukama.NodeConnectivity value (1 = online,
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
