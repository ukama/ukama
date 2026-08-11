/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"sort"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ukama/ukama/systems/analytics/aggregator/pb/gen"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/rollup"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// Query is THE read API: one question shape for latest values, time series
// and top-N breakdowns (see docs/read-api-simplification-plan.md, Phase 2).
//
//	kpis + filters + group_by + range + granularity (+ sort/top)
//
// Invariants:
//   - Aggregation is derived from the KPI's kind (flow → SUM; gauge →
//     latest level, folded across scopes per scope_agg). `agg` is an expert
//     override only.
//   - Filters always fold: the answer's grain is exactly
//     (filter keys ∩ kpi scope) ∪ group_by. With neither, the grain is
//     empty — one row per KPI, totalled over everything.
//   - granularity=total → one point per row; day|week|month → calendar
//     buckets (materialized rollups) intersecting the range.
//   - Trend is the same aggregation over the previous range (total) or the
//     previous bucket (series), computed at the answer's grain.
//
// Planner routing: bucketed reads and calendar-token totals fold
// kpi_rollups; rolling/custom totals fold kpi_windows components.
func (s *AggregatorServer) Query(ctx context.Context, req *pb.QueryRequest) (*pb.QueryResponse, error) {
	if len(req.Kpis) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one kpi is required")
	}

	kpis := s.knownKpis(req.Kpis)

	if err := validateScopeKeys(kpis, req.Filter); err != nil {
		return nil, err
	}

	if err := validateGroupBy(kpis, req.GroupBy); err != nil {
		return nil, err
	}

	gran, err := validateGranularity(req.Granularity)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	rng, err := s.resolveQueryRange(req.Range, req.From, req.To, now)
	if err != nil {
		return nil, err
	}

	rows := make([]*pb.QueryRow, 0)

	for _, kpi := range kpis {
		op, err := resolveAgg(kpi, req.Agg)
		if err != nil {
			return nil, err
		}

		grain := grainKeys(kpi, req.Filter, req.GroupBy)

		var kpiRows []*pb.QueryRow

		if gran == granTotal {
			kpiRows, err = s.queryTotal(kpi, op, grain, req.Filter, rng, now)
		} else {
			kpiRows, err = s.querySeries(kpi, op, grain, req.Filter, gran, rng, now)
		}

		if err != nil {
			return nil, err
		}

		rows = append(rows, kpiRows...)
	}

	sortRows(rows, req.Sort)

	if req.Top > 0 && int(req.Top) < len(rows) {
		rows = rows[:req.Top]
	}

	return &pb.QueryResponse{
		Rows:        rows,
		Granularity: gran,
		From:        rng.from.UTC().Format(time.RFC3339),
		To:          rng.to.UTC().Format(time.RFC3339),
	}, nil
}

// --- request resolution ---

const (
	granTotal = "total"
	granDay   = "day"
	granWeek  = "week"
	granMonth = "month"
)

// granSpan maps a series granularity to the materialized rollup span.
var granSpan = map[string]string{
	granDay:   rollup.SpanDaily,
	granWeek:  rollup.SpanWeekly,
	granMonth: rollup.SpanMonthly,
}

func validateGranularity(g string) (string, error) {
	g = strings.ToLower(g)
	if g == "" {
		g = granTotal
	}

	if g != granTotal {
		if _, ok := granSpan[g]; !ok {
			return "", status.Errorf(codes.InvalidArgument,
				"granularity must be total, day, week or month, got %q", g)
		}
	}

	return g, nil
}

// queryRange is the resolved question period plus how to obtain its
// predecessor for trend.
type queryRange struct {
	from, to time.Time
	// calendarSpan is set for today/this_week/this_month tokens — the total
	// answer then reads that span's materialized rollup rows and trends
	// against the previous span.
	calendarSpan string
	prevFrom     time.Time
	prevTo       time.Time
}

func (s *AggregatorServer) resolveQueryRange(token, fromStr, toStr string, now time.Time) (queryRange, error) {
	token = strings.ToLower(token)

	if token != "" && (fromStr != "" || toStr != "") {
		return queryRange{}, status.Error(codes.InvalidArgument,
			"range token and from/to are mutually exclusive")
	}

	// Custom absolute range.
	if fromStr != "" || toStr != "" {
		if fromStr == "" || toStr == "" {
			return queryRange{}, status.Error(codes.InvalidArgument,
				"custom range needs both from and to")
		}

		from, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return queryRange{}, status.Errorf(codes.InvalidArgument, "invalid from: %v", err)
		}

		to, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return queryRange{}, status.Errorf(codes.InvalidArgument, "invalid to: %v", err)
		}

		if !to.After(from) {
			return queryRange{}, status.Error(codes.InvalidArgument, "to must be after from")
		}

		length := to.Sub(from)

		return queryRange{from: from, to: to, prevFrom: from.Add(-length), prevTo: from}, nil
	}

	if token == "" {
		token = "this_month" // default: the business-reporting period
	}

	// Rolling tokens: a trailing window ending now.
	if d, ok := map[string]time.Duration{
		"last_24h": 24 * time.Hour,
		"last_7d":  7 * 24 * time.Hour,
		"last_30d": 30 * 24 * time.Hour,
	}[token]; ok {
		return queryRange{
			from: now.Add(-d), to: now,
			prevFrom: now.Add(-2 * d), prevTo: now.Add(-d),
		}, nil
	}

	// Calendar tokens: the current span so far, in the org timezone.
	span, ok := map[string]string{
		"today":      rollup.SpanDaily,
		"this_week":  rollup.SpanWeekly,
		"this_month": rollup.SpanMonthly,
	}[token]
	if !ok {
		return queryRange{}, status.Errorf(codes.InvalidArgument,
			"unknown range %q (today|this_week|this_month|last_24h|last_7d|last_30d or from/to)", token)
	}

	start, err := rollup.SpanStart(span, now, s.loc)
	if err != nil {
		return queryRange{}, status.Errorf(codes.Internal, "resolving span start: %v", err)
	}

	prevStart, err := rollup.PrevSpanStart(span, start)
	if err != nil {
		return queryRange{}, status.Errorf(codes.Internal, "resolving previous span: %v", err)
	}

	return queryRange{
		from: start, to: now,
		calendarSpan: span,
		prevFrom:     prevStart, prevTo: start,
	}, nil
}

// resolveAgg picks the aggregation: kind-derived default, or the expert
// `agg` override. Every component op is always computable — rollup rows
// carry the components, nothing is materialized per op.
func resolveAgg(kpi schema.KpiSpec, agg string) (string, error) {
	if agg == "" {
		return kpi.DefaultReadOp(), nil
	}

	op := strings.ToUpper(agg)
	if !componentOps[op] {
		return "", status.Errorf(codes.InvalidArgument,
			"agg must be one of sum, avg, min, max, count; got %q", agg)
	}

	return op, nil
}

// grainKeys is the answer's grain: the filter keys this KPI is scoped by,
// plus the requested group_by dimensions.
func grainKeys(kpi schema.KpiSpec, filter map[string]string, groupBy []string) []string {
	seen := map[string]bool{}
	keys := implicitGroupKeys(kpi, filter)

	for _, k := range keys {
		seen[k] = true
	}

	for _, k := range groupBy {
		if !seen[k] {
			keys = append(keys, k)
			seen[k] = true
		}
	}

	sort.Strings(keys)

	return keys
}

// --- total (one point per row) ---

func (s *AggregatorServer) queryTotal(kpi schema.KpiSpec, op string, grain []string,
	filter map[string]string, rng queryRange, now time.Time) ([]*pb.QueryRow, error) {
	var (
		cur, prev map[string]float64
		partial   map[string]bool
		err       error
	)

	if rng.calendarSpan != "" {
		cur, partial, err = s.totalFromRollups(kpi, op, grain, filter, rng.calendarSpan, rng.from, rng.to)
		if err != nil {
			return nil, err
		}

		prev, _, err = s.totalFromRollups(kpi, op, grain, filter, rng.calendarSpan, rng.prevFrom, rng.prevTo)
	} else {
		cur, err = s.totalFromWindows(kpi, op, grain, filter, rng.from, rng.to)
		if err != nil {
			return nil, err
		}

		prev, err = s.totalFromWindows(kpi, op, grain, filter, rng.prevFrom, rng.prevTo)
		partial = nil // window folds include the in-progress window: always partial
	}

	if err != nil {
		return nil, err
	}

	rows := make([]*pb.QueryRow, 0, len(cur))

	for _, key := range sortedKeys(cur) {
		prevVal, hasPrev := prev[key]

		isPartial := true
		if partial != nil {
			isPartial = partial[key]
		}

		rows = append(rows, &pb.QueryRow{
			Kpi:    kpi.Kpi,
			Dims:   schema.ParseScope(key),
			Type:   kpi.Output.Type,
			Unit:   kpi.Output.Unit,
			Symbol: kpi.Output.Symbol,
			Points: []*pb.QueryPoint{{
				From:      rng.from.UTC().Format(time.RFC3339),
				To:        rng.to.UTC().Format(time.RFC3339),
				Value:     cur[key],
				Trend:     trendFromValues(cur[key], prevVal, hasPrev),
				IsPartial: isPartial,
			}},
		})
	}

	return rows, nil
}

// totalFromRollups folds one calendar span's rollup rows to the grain.
func (s *AggregatorServer) totalFromRollups(kpi schema.KpiSpec, op string, grain []string,
	filter map[string]string, span string, from, to time.Time) (map[string]float64, map[string]bool, error) {
	rows, err := s.rollups.Range(s.org, kpi.Kpi, span, from, to, filter)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "reading rollups: %v", err)
	}

	values := map[string]float64{}
	partial := map[string]bool{}

	for _, g := range latestGroups(groupRollups(rows, grain)) {
		if v, ok := g.value(op); ok {
			values[g.scope] = v
			partial[g.scope] = g.isPartial
		}
	}

	return values, partial, nil
}

// totalFromWindows folds kpi_windows components over [from, to) to the grain.
func (s *AggregatorServer) totalFromWindows(kpi schema.KpiSpec, op string, grain []string,
	filter map[string]string, from, to time.Time) (map[string]float64, error) {
	fromID := s.grid.WindowAt(from.UTC()).ID
	toID := s.grid.WindowAt(to.UTC()).ID + 1

	perScope, err := s.aggregateWindows(kpi.Kpi, fromID, toID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "reading kpi windows: %v", err)
	}

	values := map[string]float64{}

	for key, g := range foldWindowAggs(perScope, filter, grain) {
		if strings.EqualFold(op, "LAST") {
			// Additive-gauge fold: current level = Σ of each member
			// series' latest window value.
			values[key] = g.lastSum

			continue
		}

		if v, ok := rollingOpValue(op, g.agg); ok {
			values[key] = v
		}
	}

	return values, nil
}

// --- series (calendar buckets from materialized rollups) ---

func (s *AggregatorServer) querySeries(kpi schema.KpiSpec, op string, grain []string,
	filter map[string]string, gran string, rng queryRange, now time.Time) ([]*pb.QueryRow, error) {
	span := granSpan[gran]

	// Align the read to the first bucket containing `from` so partial-edge
	// buckets are included rather than silently dropped.
	alignedFrom, err := rollup.SpanStart(span, rng.from, s.loc)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "aligning range: %v", err)
	}

	rows, err := s.rollups.Range(s.org, kpi.Kpi, span, alignedFrom, rng.to, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "reading rollups: %v", err)
	}

	// bucketed[grainKey][spanStart] -> folded group
	grouped := groupRollups(rows, grain)

	out := make([]*pb.QueryRow, 0, len(grouped))

	for _, key := range sortedGroupKeys(grouped) {
		bySpan := grouped[key]

		starts := make([]time.Time, 0, len(bySpan))
		for t := range bySpan {
			starts = append(starts, t)
		}

		sort.Slice(starts, func(i, j int) bool { return starts[i].Before(starts[j]) })

		points := make([]*pb.QueryPoint, 0, len(starts))

		var prevVal float64

		hasPrev := false

		for _, t := range starts {
			g := bySpan[t]

			v, ok := g.value(op)
			if !ok {
				continue
			}

			points = append(points, &pb.QueryPoint{
				From:      g.spanStart.UTC().Format(time.RFC3339),
				To:        g.spanEnd.UTC().Format(time.RFC3339),
				Value:     v,
				Trend:     trendFromValues(v, prevVal, hasPrev),
				IsPartial: g.isPartial,
			})

			prevVal, hasPrev = v, true
		}

		if len(points) == 0 {
			continue
		}

		out = append(out, &pb.QueryRow{
			Kpi:    kpi.Kpi,
			Dims:   schema.ParseScope(key),
			Type:   kpi.Output.Type,
			Unit:   kpi.Output.Unit,
			Symbol: kpi.Output.Symbol,
			Points: points,
		})
	}

	return out, nil
}

// --- shared bits ---

// trendFromValues is the one trend definition: current vs previous, at the
// answer's grain. Same up|down|flat|new|na semantics as the rollup engine.
func trendFromValues(cur, prev float64, hasPrev bool) *pb.Trend {
	if !hasPrev {
		return &pb.Trend{Direction: "new"}
	}

	t := &pb.Trend{
		HasPrevious: true,
		PrevValue:   prev,
		ChangeAbs:   cur - prev,
	}

	switch {
	case prev == 0 && cur == 0:
		t.Direction = "flat"
	case prev == 0:
		t.Direction = "na" // percent undefined; change_abs still served
	default:
		changePct := (cur - prev) / abs(prev) * 100
		t.ChangePct = changePct

		switch {
		case changePct == 0:
			t.Direction = "flat"
		case changePct > 0:
			t.Direction = "up"
		default:
			t.Direction = "down"
		}
	}

	return t
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}

	return v
}

// sortRows orders by each row's latest point value: "-value" (default when
// sorting is requested) descending, "value" ascending.
func sortRows(rows []*pb.QueryRow, token string) {
	token = strings.ToLower(strings.TrimSpace(token))
	if token == "" {
		return // stable spec order: kpi request order, then dims
	}

	asc := token == "value"

	if !asc && token != "-value" {
		return
	}

	last := func(r *pb.QueryRow) float64 {
		if len(r.Points) == 0 {
			return 0
		}

		return r.Points[len(r.Points)-1].Value
	}

	sort.SliceStable(rows, func(i, j int) bool {
		if asc {
			return last(rows[i]) < last(rows[j])
		}

		return last(rows[i]) > last(rows[j])
	})
}

func sortedKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func sortedGroupKeys(m map[string]map[time.Time]*scopeGroup) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
