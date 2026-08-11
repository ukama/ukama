/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"math"
	"sort"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// Read-time scope folding.
//
// Fine-grained KPIs (DATA_USAGE is scoped per network×site×package×
// assignment×iccid series) would otherwise answer "usage of network X" with
// N per-series rows. Folding aggregates rows' persisted COMPONENTS
// (Sum/Count/Min/Max/Last) — the same exact fold the rollup engine applies
// across windows, so AVG stays weighted (Σsum/Σcount), never an average of
// averages, and LAST folds as the sum of each member's latest level (the
// additive-gauge cross-scope semantics). The grain is always
// (filter keys ∩ kpi scope) ∪ group_by — the Query planner and the legacy
// adapters both fold through here.

// componentOps can be folded across scopes exactly.
var componentOps = map[string]bool{
	"SUM": true, "AVG": true, "MIN": true, "MAX": true, "COUNT": true,
}

// validateScopeKeys rejects filter keys that no requested KPI is scoped by:
// a typo'd key would otherwise silently match nothing and read as
// "value is 0/absent" instead of "you asked a malformed question".
func validateScopeKeys(kpis []schema.KpiSpec, filter map[string]string) error {
	for key := range filter {
		known := false

		for _, kpi := range kpis {
			for _, sc := range kpi.Scope {
				if sc == key {
					known = true

					break
				}
			}
		}

		if !known {
			return status.Errorf(codes.InvalidArgument,
				"scope filter key %q is not a scope dimension of any requested kpi", key)
		}
	}

	return nil
}

// validateGroupBy requires every group_by key to be a scope dimension of
// every requested KPI (folding by a key a KPI does not carry is undefined).
func validateGroupBy(kpis []schema.KpiSpec, groupBy []string) error {
	for _, key := range groupBy {
		for _, kpi := range kpis {
			found := false

			for _, sc := range kpi.Scope {
				if sc == key {
					found = true

					break
				}
			}

			if !found {
				return status.Errorf(codes.InvalidArgument,
					"group_by key %q is not a scope dimension of kpi %s (scope: %s)",
					key, kpi.Kpi, strings.Join(kpi.Scope, ","))
			}
		}
	}

	return nil
}

// implicitGroupKeys derives the fold grain from a scope filter: the filter
// keys this KPI is actually scoped by, sorted. Keys outside the KPI's scope
// are dropped (the filter already yields no rows for that KPI). Returns nil
// when nothing applies — the read stays ungrouped.
func implicitGroupKeys(kpi schema.KpiSpec, filter map[string]string) []string {
	keys := make([]string, 0, len(filter))

	for key := range filter {
		for _, sc := range kpi.Scope {
			if sc == key {
				keys = append(keys, key)

				break
			}
		}
	}

	if len(keys) == 0 {
		return nil
	}

	sort.Strings(keys)

	return keys
}

// projectScope reduces a canonical scope JSON to the requested keys,
// returning the canonical JSON of the projection (the group key).
func projectScope(scope string, keys []string) string {
	full := schema.ParseScope(scope)
	projected := make(map[string]string, len(keys))

	for _, k := range keys {
		projected[k] = full[k]
	}

	return schema.CanonicalScope(projected)
}

// scopeGroup accumulates components of the rollup rows folded into one
// group_by bucket.
type scopeGroup struct {
	scope      string // canonical projection JSON
	spanStart  time.Time
	spanEnd    time.Time
	sum        float64
	count      float64
	min        float64
	max        float64
	lastSum    float64 // Σ row.Last — each member's latest level (additive gauges)
	isPartial  bool
	meta       schema.KpiRollup // first folded row (units/type metadata)
	computedAt time.Time
}

func (g *scopeGroup) fold(row schema.KpiRollup) {
	g.sum += row.Sum
	g.count += row.Count
	g.min = math.Min(g.min, row.Min)
	g.max = math.Max(g.max, row.Max)
	g.lastSum += row.Last
	g.isPartial = g.isPartial || row.IsPartial

	if row.ComputedAt.After(g.computedAt) {
		g.computedAt = row.ComputedAt
	}
}

func (g *scopeGroup) value(op string) (float64, bool) {
	switch strings.ToUpper(op) {
	case "SUM":
		return g.sum, true
	case "COUNT":
		return g.count, true
	case "AVG":
		if g.count == 0 {
			return 0, false
		}

		return g.sum / g.count, true
	case "MIN":
		return g.min, true
	case "MAX":
		return g.max, true
	case "LAST":
		// Sum of each member's latest level — the additive-gauge fold
		// ("current total across scopes").
		return g.lastSum, true
	default:
		return 0, false
	}
}

// groupRollups folds rollup rows into one bucket per (group key, span_start).
// Latest-value reads then keep, per group key, only the newest span_start —
// Latest returns each scope's newest span row, and scopes refresh at
// different times, so folding across mixed span_starts would mix periods.
func groupRollups(rows []schema.KpiRollup, groupBy []string) map[string]map[time.Time]*scopeGroup {
	out := map[string]map[time.Time]*scopeGroup{}

	for _, row := range rows {
		key := projectScope(row.Scope, groupBy)

		bySpan, ok := out[key]
		if !ok {
			bySpan = map[time.Time]*scopeGroup{}
			out[key] = bySpan
		}

		g, ok := bySpan[row.SpanStart]
		if !ok {
			g = &scopeGroup{
				scope:     key,
				spanStart: row.SpanStart,
				spanEnd:   row.SpanEnd,
				min:       math.Inf(1),
				max:       math.Inf(-1),
				meta:      row,
			}
			bySpan[row.SpanStart] = g
		}

		g.fold(row)
	}

	return out
}

// latestGroups reduces the (group, span_start) buckets to each group's
// newest span, sorted by group key for deterministic output.
func latestGroups(grouped map[string]map[time.Time]*scopeGroup) []*scopeGroup {
	out := make([]*scopeGroup, 0, len(grouped))

	for _, bySpan := range grouped {
		var newest *scopeGroup

		for _, g := range bySpan {
			if newest == nil || g.spanStart.After(newest.spanStart) {
				newest = g
			}
		}

		out = append(out, newest)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].scope < out[j].scope })

	return out
}
