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
	"strings"
	"time"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// Rolling-range reads ("last 24h/7d/30d" and custom from/to totals) fold
// the raw kpi_windows components over the exact trailing window — precise
// where calendar rollups would be bucket-aligned. The Query planner
// (query.go) is the only consumer.
const (
	SpanLast24h = "last_24h"
	SpanLast7d  = "last_7d"
	SpanLast30d = "last_30d"
)

// rollingLookback maps each rolling token to its trailing duration.
var rollingLookback = map[string]time.Duration{
	SpanLast24h: 24 * time.Hour,
	SpanLast7d:  7 * 24 * time.Hour,
	SpanLast30d: 30 * 24 * time.Hour,
}

// isRollingSpan reports whether span is one of the rolling-window tokens.
func isRollingSpan(span string) bool {
	_, ok := rollingLookback[strings.ToLower(span)]

	return ok
}

// windowAgg accumulates kpi_windows components for one scope over a window
// range — the same component set the rollup engine aggregates per calendar
// span, so every read-time aggregation is exact.
type windowAgg struct {
	sum, count, min, max float64
	lastWindow           int64
	lastValue            float64
}

// aggregateWindows reads kpi_windows for a KPI in [fromID, toID) and folds
// their components into one windowAgg per scope.
func (s *AggregatorServer) aggregateWindows(kpiKey string, fromID, toID int64) (map[string]*windowAgg, error) {
	rows, err := s.windows.WindowsInRange(s.org, kpiKey, fromID, toID)
	if err != nil {
		return nil, err
	}

	out := map[string]*windowAgg{}

	for _, row := range rows {
		a, ok := out[row.Scope]
		if !ok {
			a = &windowAgg{min: math.Inf(1), max: math.Inf(-1), lastWindow: -1}
			out[row.Scope] = a
		}

		a.sum += row.Sum
		a.count += row.Count
		a.min = math.Min(a.min, row.Min)
		a.max = math.Max(a.max, row.Max)

		if row.WindowID > a.lastWindow {
			a.lastWindow = row.WindowID
			a.lastValue = row.Value
		}
	}

	return out, nil
}

// rollingOpValue applies an aggregation to folded components — identical
// semantics to the rollup-row reads (weighted AVG, never avg-of-avg).
func rollingOpValue(op string, a *windowAgg) (float64, bool) {
	switch strings.ToUpper(op) {
	case "SUM":
		return a.sum, true
	case "COUNT":
		return a.count, true
	case "AVG":
		if a.count == 0 {
			return 0, false
		}

		return a.sum / a.count, true
	case "MIN":
		return a.min, true
	case "MAX":
		return a.max, true
	case "LAST":
		return a.lastValue, true
	default:
		return 0, false
	}
}

// scopeMatches mirrors the rollup repo's filterScope: with no filter every
// scope matches; otherwise every requested key/value must be present in the
// row's scope. The org bucket (empty scope) is org-only — it never matches a
// per-network filter, so unresolved-SIM revenue can't inflate a network.
func scopeMatches(scope string, filter map[string]string) bool {
	if len(filter) == 0 {
		return true
	}

	m := schema.ParseScope(scope)

	for k, v := range filter {
		if m[k] != v {
			return false
		}
	}

	return true
}

// rollGroup is one fold bucket of windowAggs: the aggregated components
// plus the additive-gauge fold of each member's latest level.
type rollGroup struct {
	agg     *windowAgg
	lastSum float64 // Σ members' latest-window values
}

// foldWindowAggs applies the scope filter, then folds windowAgg components
// into one aggregate per distinct combination of the group keys (keyed by
// the projected canonical scope JSON).
func foldWindowAggs(aggs map[string]*windowAgg, filter map[string]string, groupBy []string) map[string]*rollGroup {
	out := map[string]*rollGroup{}

	for scope, agg := range aggs {
		if !scopeMatches(scope, filter) {
			continue
		}

		key := projectScope(scope, groupBy)

		g, ok := out[key]
		if !ok {
			g = &rollGroup{agg: &windowAgg{min: math.Inf(1), max: math.Inf(-1), lastWindow: -1}}
			out[key] = g
		}

		g.agg.sum += agg.sum
		g.agg.count += agg.count
		g.agg.min = math.Min(g.agg.min, agg.min)
		g.agg.max = math.Max(g.agg.max, agg.max)
		g.lastSum += agg.lastValue
	}

	return out
}
