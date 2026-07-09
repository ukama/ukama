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

	"github.com/ukama/ukama/systems/analytics/schema"
)

// ActiveCustomers (ACTIVE_CUSTOMERS @ scope network_id): number of distinct
// subscribers per network holding at least one active SIM.
//
// Inputs:
//   active_sims — subscriber.sim.listActive (sim_id, subscriber_id,
//                 network_id) — source-filtered to sim_status=active
//   networks    — registry.network.getAll (network_id) — zero-fill so
//                 networks with no active sims still emit 0.
func ActiveCustomers(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sims, ok := in["active_sims"]
	if !ok {
		return nil, fmt.Errorf("ACTIVE_CUSTOMERS: missing input 'active_sims'")
	}

	networks, ok := in["networks"]
	if !ok {
		return nil, fmt.Errorf("ACTIVE_CUSTOMERS: missing input 'networks'")
	}

	// network -> set of distinct subscribers with an active sim
	subscribers := map[string]map[string]bool{}

	for _, sim := range sims {
		networkID := str(sim["network_id"])
		subscriberID := str(sim["subscriber_id"])

		if networkID == "" || subscriberID == "" {
			continue
		}

		if subscribers[networkID] == nil {
			subscribers[networkID] = map[string]bool{}
		}

		subscribers[networkID][subscriberID] = true
	}

	results := make([]Result, 0, len(networks))

	for _, network := range networks {
		networkID := str(network["network_id"])
		if networkID == "" {
			continue
		}

		results = append(results, CountResult(
			map[string]string{"network_id": networkID},
			float64(len(subscribers[networkID])),
		))
	}

	return results, nil
}
