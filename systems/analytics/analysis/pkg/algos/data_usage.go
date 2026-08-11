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

// DataUsage (DATA_USAGE @ scope network_id+site_id+package_id+
// sim_package_id+iccid): bytes consumed IN the window, per Prometheus series
// of the metrics system's cumulative data_usage counter
// (metrics.data_usage.last). package_id is the CATALOG package (cdr's
// `dataplan` label, renamed at ingest to match the rest of analytics);
// sim_package_id is the sim's assignment instance.
//
// The source counter is cumulative, and /v1/last carries no in-response
// baseline — so the KPI stores per-window INCREMENTS, not snapshots:
//
//	increment = clamp₀(counter now − counter as of the previous window)
//
// using the same dataset read twice: mode state (as of this window) and mode
// state_prev (as of the previous one). Increments telescope, so any span
// rollup is an exact SUM — and, unlike counter snapshots read with a max−min
// DELTA, they stay exact under ANY read-time filter or group_by fold across
// series. Rows carry Value = Sum = increment, Count = 1.
//
// Semantics at the edges:
//   - First appearance of a series: no baseline yet → increment 0 (the
//     pre-history consumption is unknowable; counting the whole counter
//     would double-count what was consumed before the pipeline watched).
//   - Idle series: pushgateway re-exposure keeps them in /v1/last →
//     explicit 0 rows (zero-fill), never a gap.
//   - Counter reset (cdr restart): cur < prev → clamp to 0. The residual
//     consumed before the reset is lost — accepted approximation.
func DataUsage(win schema.Window, in Datasets, spec schema.KpiSpec) ([]Result, error) {
	usage, ok := in["usage"]
	if !ok {
		return nil, fmt.Errorf("DATA_USAGE: missing input 'usage'")
	}

	prev, ok := in["usage_prev"]
	if !ok {
		return nil, fmt.Errorf("DATA_USAGE: missing input 'usage_prev'")
	}

	baseline := make(map[string]float64, len(prev))

	for _, rec := range prev {
		key, ok := seriesKey(rec)
		if !ok {
			continue
		}

		baseline[key] = sampleValue(rec["value"])
	}

	results := make([]Result, 0, len(usage))

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
		}
		// first appearance = baseline only, increment 0

		results = append(results, Result{
			Scope: map[string]string{
				"network_id":     str(rec["network_id"]),
				"site_id":        str(rec["site_id"]),
				"package_id":     str(rec["package_id"]),
				"sim_package_id": str(rec["sim_package_id"]),
				"iccid":          str(rec["iccid"]),
			},
			Value: increment,
			Sum:   increment,
			Count: 1,
			Min:   increment,
			Max:   increment,
		})
	}

	return results, nil
}

// seriesKey identifies one Prometheus series across windows — the same
// composite the ingest spec uses as its entity key
// (iccid|sim_package_id|site_id).
func seriesKey(rec map[string]interface{}) (string, bool) {
	iccid, simPkg, site := str(rec["iccid"]), str(rec["sim_package_id"]), str(rec["site_id"])
	if iccid == "" || simPkg == "" {
		return "", false
	}

	return strings.Join([]string{iccid, simPkg, site}, "|"), true
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
