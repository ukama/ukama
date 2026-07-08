/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package algos

import (
	"fmt"
	"strings"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// ActiveCustomers (ACTIVE_CUSTOMERS @ scope network_id): number of
// subscribers per network with at least one SIM in "active" status.
//
// Inputs:
//   subscribers — subscriber.registry.getByNetwork (subscriber_id,
//                 network_id, sims[] with status)
//   networks    — registry.network.getAll (network_id) — zero-fill.
func ActiveCustomers(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	subscribers, ok := in["subscribers"]
	if !ok {
		return nil, fmt.Errorf("ACTIVE_CUSTOMERS: missing input 'subscribers'")
	}

	networks, ok := in["networks"]
	if !ok {
		return nil, fmt.Errorf("ACTIVE_CUSTOMERS: missing input 'networks'")
	}

	active := map[string]int{} // network -> active subscriber count

	for _, sub := range subscribers {
		networkID := str(sub["network_id"])
		if networkID == "" {
			continue
		}

		if hasActiveSim(sub["sims"]) {
			active[networkID]++
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
			float64(active[networkID]),
		))
	}

	return results, nil
}

func hasActiveSim(sims interface{}) bool {
	arr, ok := sims.([]interface{})
	if !ok {
		return false
	}

	for _, s := range arr {
		sim, ok := s.(map[string]interface{})
		if !ok {
			continue
		}

		if strings.EqualFold(str(sim["status"]), "active") {
			return true
		}
	}

	return false
}
