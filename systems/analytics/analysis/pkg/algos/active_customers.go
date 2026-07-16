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

// ActiveCustomers (ACTIVE_CUSTOMERS @ scope network_id): number of distinct
// subscribers per network with at least one SIM in "active" status.
//
// Inputs:
//   sims     — subscriber.sim.list (sim_id, subscriber_id, network_id,
//              status): ALL sims; the active filter is applied here. (v3:
//              shares the dataset with the package/usage KPIs instead of a
//              separate source-filtered pull.)
//   networks — registry.network.getAll (network_id) — zero-fill so networks
//              with no active sims still emit 0.
func ActiveCustomers(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	sims, ok := in["sims"]
	if !ok {
		return nil, fmt.Errorf("ACTIVE_CUSTOMERS: missing input 'sims'")
	}

	networks, ok := in["networks"]
	if !ok {
		return nil, fmt.Errorf("ACTIVE_CUSTOMERS: missing input 'networks'")
	}

	// network -> set of distinct subscribers with an active sim
	active := map[string]map[string]bool{}

	for _, sim := range sims {
		networkID := str(sim["network_id"])
		subscriberID := str(sim["subscriber_id"])

		if networkID == "" || subscriberID == "" {
			continue
		}

		if strings.EqualFold(str(sim["status"]), "active") {
			if active[networkID] == nil {
				active[networkID] = map[string]bool{}
			}

			active[networkID][subscriberID] = true
		}
	}

	return zeroFilled(networks, func(networkID string) float64 {
		return float64(len(active[networkID]))
	}), nil
}
