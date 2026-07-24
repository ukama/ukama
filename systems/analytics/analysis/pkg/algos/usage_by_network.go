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
	"strconv"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// UsageByNetwork (USAGE_BY_NETWORK @ scope network_id): the network's
// CUMULATIVE data usage (bytes) as of the window — the sum, over the network's
// sims, of each sim's lifetime usage counter.
//
// The usage input (subscriber.usage.getBySim) is a full_snapshot of the
// per-sim cumulative counter from /v1/usages (no from/to: the endpoint returns
// the lifetime counter, which always 200s — the windowed CDR path errors on
// idle windows). Read with mode: state, so each window carries every sim's
// latest counter as of that window.
//
// The per-window value is a SNAPSHOT, not an increment. Data consumed over a
// span is derived at read time as DELTA = max - min of this cumulative across
// the span's windows (see the aggregator's DELTA op). To make that exact, each
// row sets Value = Sum = Min = Max = the network's cumulative total and
// Count = 1. Networks with no sims/usage zero-fill (total 0) so the console
// shows "0 B", never "—".
//
// Caveat: DELTA assumes the counter is monotonic across the span. A mid-span
// reset (rollover / re-provision) is approximated (negative deltas clamp to 0).
func UsageByNetwork(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	usage, ok := in["usage"]
	if !ok {
		return nil, fmt.Errorf("USAGE_BY_NETWORK: missing input 'usage'")
	}

	networks, ok := in["networks"]
	if !ok {
		return nil, fmt.Errorf("USAGE_BY_NETWORK: missing input 'networks'")
	}

	// Sum each sim's cumulative usage counter into its network.
	byNetwork := map[string]float64{}

	for _, rec := range usage {
		networkID := str(rec["network_id"])
		if networkID == "" {
			continue
		}

		byNetwork[networkID] += sumNumeric(rec["usage"])
	}

	// One row per network (zero-fill), carrying the cumulative total as of this
	// window with Min=Max=Value so the read-time DELTA op yields max-min.
	results := make([]Result, 0, len(networks))

	for _, network := range networks {
		networkID := str(network["network_id"])
		if networkID == "" {
			continue
		}

		total := byNetwork[networkID] // 0 when the network has no sims/usage

		results = append(results, Result{
			Scope: map[string]string{"network_id": networkID},
			Value: total,
			Sum:   total,
			Count: 1,
			Min:   total,
			Max:   total,
		})
	}

	return results, nil
}

// sumNumeric adds up every numeric leaf in a decoded JSON value — the usage
// object from /v1/usages has dynamic keys (per imsi/iccid) with byte counts as
// values, sometimes serialized as strings.
func sumNumeric(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case string:
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return f
		}

		return 0
	case map[string]interface{}:
		total := 0.0
		for _, val := range t {
			total += sumNumeric(val)
		}

		return total
	case []interface{}:
		total := 0.0
		for _, val := range t {
			total += sumNumeric(val)
		}

		return total
	default:
		return 0
	}
}
