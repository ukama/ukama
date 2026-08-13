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
	"strings"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// DataUsage stores bytes consumed IN the window, per Prometheus series of the
// metrics system's cumulative data_usage counter:
//
//	increment = clamp₀(counter now − counter as of the previous window)
//
// reading metrics.data_usage.last twice — mode state (as of this window) and
// mode state_prev (as of the previous one). Increments telescope, so any span
// rollup is an exact SUM under any read-time filter or group_by fold. Rows
// carry Value = Sum = increment, Count = 1.
//
// Edges:
//   - params.first_value "count" counts a series' first observed value as
//     consumption; "baseline" (the default) records it as baseline only.
//   - Idle series stay in /v1/last and emit explicit 0 rows, never a gap.
//   - Counter reset (cur < prev) clamps to 0; the pre-reset residual is lost.
//
// package_id is the CATALOG package (cdr's `dataplan` label, renamed at
// ingest); sim_package_id is the sim's assignment instance.
func DataUsage(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	usage, ok := in["usage"]
	if !ok {
		return nil, fmt.Errorf("DATA_USAGE: missing input 'usage'")
	}

	prev, ok := in["usage_prev"]
	if !ok {
		return nil, fmt.Errorf("DATA_USAGE: missing input 'usage_prev'")
	}

	countFirst := spec.Params["first_value"] == "count"

	baseline := make(map[string]float64, len(prev))

	for _, rec := range prev {
		key, ok := seriesKey(rec)
		if !ok {
			continue
		}

		baseline[key] = sampleValue(rec["value"])
	}

	// Increments are per SERIES (session included) but summed into the KPI's
	// scope: several series share one scope, and separate Results would
	// collide on the (kpi, scope, window) unique index.
	type scopeAgg struct {
		scope map[string]string
		total float64
	}

	byScope := make(map[string]*scopeAgg, len(usage))
	order := make([]string, 0, len(usage))

	for _, rec := range usage {
		key, ok := seriesKey(rec)
		if !ok {
			continue // a series without its identity labels is unattributable
		}

		cur := sampleValue(rec["value"])

		increment := 0.0

		if base, seen := baseline[key]; seen {
			if d := cur - base; d > 0 {
				increment = d
			}
			// d < 0 = counter reset → clamp to 0
		} else if countFirst && cur > 0 {
			increment = cur
		}

		scope := map[string]string{
			"network_id":     str(rec["network_id"]),
			"site_id":        str(rec["site_id"]),
			"package_id":     str(rec["package_id"]),
			"sim_package_id": str(rec["sim_package_id"]),
			"iccid":          str(rec["iccid"]),
		}

		scopeKey := schema.CanonicalScope(scope)

		agg, seen := byScope[scopeKey]
		if !seen {
			agg = &scopeAgg{scope: scope}
			byScope[scopeKey] = agg
			order = append(order, scopeKey)
		}

		agg.total += increment
	}

	results := make([]Result, 0, len(byScope))

	for _, scopeKey := range order {
		agg := byScope[scopeKey]

		results = append(results, Result{
			Scope: agg.scope,
			Value: agg.total,
			Sum:   agg.total,
			Count: 1,
			Min:   agg.total,
			Max:   agg.total,
		})
	}

	return results, nil
}

// seriesKey identifies one Prometheus series across windows, matching the
// ingest spec's entity key. session_id is optional: series pushed without the
// label share the empty component.
func seriesKey(rec map[string]interface{}) (string, bool) {
	iccid, simPkg := str(rec["iccid"]), str(rec["sim_package_id"])
	if iccid == "" || simPkg == "" {
		return "", false
	}

	return strings.Join([]string{
		iccid, simPkg, str(rec["site_id"]), str(rec["session_id"]),
	}, "|"), true
}

// sampleValue extracts the numeric value from a mapped Prometheus sample.
// /v1/last serializes an instant vector, so `value` is the pair
// [unix_ts, "value-as-string"]; be liberal and also accept a bare number or
// numeric string.
func sampleValue(v interface{}) float64 {
	switch t := v.(type) {
	case []interface{}:
		if len(t) == 0 {
			return 0
		}

		return sampleValue(t[len(t)-1])
	case float64:
		return t
	case string:
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return f
		}

		return 0
	default:
		return 0
	}
}
