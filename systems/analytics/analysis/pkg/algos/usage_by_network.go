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
	"math"
	"strconv"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// UsageByNetwork (USAGE_BY_NETWORK @ scope network_id): total data usage
// (bytes) per network for the window, summed over the per-sim usage records.
//
// Inputs:
//   usage    — subscriber.usage.getBySim (mode: window): one record per
//              active sim per window; the "usage" field is the dynamic
//              usage object from /v1/usages, whose numeric values are the
//              consumed bytes.
//   networks — registry.network.getAll (network_id) — zero-fill.
//
// Components: Sum = total bytes, Count = number of sim usage records,
// Min/Max = smallest/largest per-sim usage — this makes daily SUM exact and
// AVG a true per-window weighted average in the aggregator.
func UsageByNetwork(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	usage, ok := in["usage"]
	if !ok {
		return nil, fmt.Errorf("USAGE_BY_NETWORK: missing input 'usage'")
	}

	networks, ok := in["networks"]
	if !ok {
		return nil, fmt.Errorf("USAGE_BY_NETWORK: missing input 'networks'")
	}

	type agg struct {
		sum, min, max float64
		count         float64
	}

	byNetwork := map[string]*agg{}

	for _, rec := range usage {
		networkID := str(rec["network_id"])
		if networkID == "" {
			continue
		}

		bytes := sumNumeric(rec["usage"])

		a, ok := byNetwork[networkID]
		if !ok {
			a = &agg{min: math.Inf(1), max: math.Inf(-1)}
			byNetwork[networkID] = a
		}

		a.sum += bytes
		a.count++
		a.min = math.Min(a.min, bytes)
		a.max = math.Max(a.max, bytes)
	}

	results := make([]Result, 0, len(networks))

	for _, network := range networks {
		networkID := str(network["network_id"])
		if networkID == "" {
			continue
		}

		a, ok := byNetwork[networkID]
		if !ok {
			results = append(results, Result{Scope: map[string]string{"network_id": networkID}})

			continue
		}

		results = append(results, Result{
			Scope: map[string]string{"network_id": networkID},
			Value: a.sum,
			Sum:   a.sum,
			Count: a.count,
			Min:   a.min,
			Max:   a.max,
		})
	}

	return results, nil
}

// sumNumeric adds up every numeric leaf in a decoded JSON value — the usage
// object from /v1/usages has dynamic keys (per imsi/sim) with byte counts as
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
