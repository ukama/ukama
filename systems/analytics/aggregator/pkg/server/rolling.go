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

	pb "github.com/ukama/ukama/systems/analytics/aggregator/pb/gen"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// Rolling-window span tokens for the /kpis/values endpoint. Unlike the
// calendar spans (daily/weekly/monthly) that the rollup engine precomputes
// into kpi_rollups, these are computed ON READ by aggregating the exact
// kpi_windows components over a trailing duration ending "now". This gives a
// true "last 24h / last 7 days / last 30 days" for every op — SUM of sums,
// weighted AVG (sum/count, never avg-of-avg), MIN/MAX, and LAST = the most
// recent window's value.
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
// span, so every op is exact.
type windowAgg struct {
	sum, count, min, max float64
	lastWindow           int64
	lastValue            float64
}

// getKpisRolling serves the values endpoint for a rolling-window token by
// aggregating kpi_windows over [now-lookback, now) (including the in-progress
// window, so results are always is_partial), with a trend vs the immediately
// preceding window [now-2*lookback, now-lookback).
func (s *AggregatorServer) getKpisRolling(req *pb.GetKpisRequest) (*pb.GetKpisResponse, error) {
	lookback := rollingLookback[strings.ToLower(req.Span)]
	now := time.Now().UTC()

	// currTo includes the current (partial) window; prevTo == currFrom so the
	// previous window abuts without overlapping.
	currFrom := s.grid.WindowAt(now.Add(-lookback)).ID
	currTo := s.grid.WindowAt(now).ID + 1
	prevFrom := s.grid.WindowAt(now.Add(-2 * lookback)).ID
	prevTo := currFrom

	from := now.Add(-lookback).Format(time.RFC3339)
	to := now.Format(time.RFC3339)
	computedAt := now.Format(time.RFC3339)

	values := make([]*pb.KpiValue, 0)

	// Unknown keys are skipped rather than failing the batch — see knownKpis.
	for _, kpi := range s.knownKpis(req.Keys) {
		op, err := s.resolveOp(kpi, req.Op)
		if err != nil {
			return nil, err
		}

		curr, err := s.aggregateWindows(kpi.Kpi, currFrom, currTo)
		if err != nil {
			return nil, err
		}

		prev, err := s.aggregateWindows(kpi.Kpi, prevFrom, prevTo)
		if err != nil {
			return nil, err
		}

		for scope, agg := range curr {
			if !scopeMatches(scope, req.Scope) {
				continue
			}

			value, ok := rollingOpValue(op, agg)
			if !ok {
				continue
			}

			values = append(values, &pb.KpiValue{
				Kpi:        kpi.Kpi,
				Value:      value,
				Span:       strings.ToLower(req.Span),
				Op:         op,
				From:       from,
				To:         to,
				Type:       kpi.Output.Type,
				Unit:       kpi.Output.Unit,
				Symbol:     kpi.Output.Symbol,
				IsPartial:  true,
				Scope:      schema.ParseScope(scope),
				Trend:      rollingTrend(op, agg, prev[scope]),
				ComputedAt: computedAt,
			})
		}
	}

	return &pb.GetKpisResponse{Values: values}, nil
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

// rollingOpValue applies an op to aggregated components — identical semantics
// to the rollup engine's opValue (weighted AVG, never avg-of-avg).
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
	case "DELTA":
		// span usage from a cumulative counter; clamp resets to 0
		if d := a.max - a.min; d > 0 {
			return d, true
		}

		return 0, true
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

// rollingTrend builds the period-over-period trend for a rolling read by
// comparing the current window aggregate against the preceding one, using the
// same up|down|flat|new|na semantics as the rollup engine's applyTrend.
func rollingTrend(op string, curr, prev *windowAgg) *pb.Trend {
	t := &pb.Trend{Direction: "new"}

	if prev == nil {
		return t
	}

	cv, ok1 := rollingOpValue(op, curr)
	pv, ok2 := rollingOpValue(op, prev)
	if !ok1 || !ok2 {
		return t
	}

	changeAbs := cv - pv

	t.Direction = ""
	t.HasPrevious = true
	t.PrevValue = pv
	t.ChangeAbs = changeAbs

	if pv == 0 {
		if cv == 0 {
			t.Direction = "flat"
		} else {
			t.Direction = "na" // percent undefined; change_abs still served
		}

		return t
	}

	changePct := changeAbs / math.Abs(pv) * 100
	t.ChangePct = changePct

	switch {
	case changePct == 0:
		t.Direction = "flat"
	case changePct > 0:
		t.Direction = "up"
	default:
		t.Direction = "down"
	}

	return t
}
