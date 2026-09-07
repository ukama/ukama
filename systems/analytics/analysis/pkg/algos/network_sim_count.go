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

// NetworkSimCountNetworksInput is the zero-fill input every spec using this
// algo must declare; every OTHER input is a sim-count series to add up.
const NetworkSimCountNetworksInput = "networks"

// NetworkSimCount sums every input other than "networks" (per-network SIM
// count gauges) into one value per network:
//
//	ACTIVE_CUSTOMERS = active_sims
//	CUSTOMERS        = active_sims + inactive_sims
//
// "networks" (registry.network.getAll) zero-fills networks with no series.
func NetworkSimCount(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	networks, ok := in[NetworkSimCountNetworksInput]
	if !ok {
		return nil, fmt.Errorf("%s: missing input %q", spec.Kpi, NetworkSimCountNetworksInput)
	}

	if len(in) < 2 {
		return nil, fmt.Errorf("%s: no sim-count input declared", spec.Kpi)
	}

	// network -> summed count across every declared sim-count dataset
	counts := map[string]float64{}

	for name, rows := range in {
		if name == NetworkSimCountNetworksInput {
			continue
		}

		for _, row := range rows {
			networkID := str(row["network_id"])
			if networkID == "" {
				continue
			}

			counts[networkID] += sampleValue(row["value"])
		}
	}

	return zeroFilled(networks, func(networkID string) float64 {
		return counts[networkID]
	}), nil
}
