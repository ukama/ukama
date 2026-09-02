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

// NetworkSimCount sums the metrics system's per-network SIM-count gauges into
// one value per network. Which gauges are summed is the spec's choice, not the
// algo's — every input other than "networks" is added:
//
//	ACTIVE_CUSTOMERS = active_sims
//	CUSTOMERS        = active_sims + inactive_sims
//
// The series are levels (sim-manager recounts from its DB and pushes an
// absolute count per network), so the window value is the count as observed,
// not an increment.
//
// Inputs:
//
//	networks — registry.network.getAll (network_id): zero-fill, so a network
//	           with no series at all emits 0 rather than a gap.
//	others   — a metrics.*_sims.last dataset (network_id, value). Series
//	           carrying no network label are unattributable and dropped.
//
// Scope is network_id alone: the underlying gauges carry no other dimension.
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
