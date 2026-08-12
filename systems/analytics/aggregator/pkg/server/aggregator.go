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
	"fmt"
	"sort"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/analytics/aggregator/pb/gen"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/performance"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/rollup"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// AggregatorServer is the KPI read API over the components rollups plus the
// performance-report composer. Query (query.go) is the primary read path;
// GetKpis/GetKpiTimeSeries/GetKpiBreakdown are thin adapters over the same
// planner — one semantics, two request shapes.
type AggregatorServer struct {
	org      string
	kpis     []schema.KpiSpec
	byKey    map[string]schema.KpiSpec
	rollups  db.RollupRepo
	composer *performance.Composer
	// grid + windows back the rolling-range reads, which fold raw
	// kpi_windows components instead of calendar rollups.
	grid    schema.Grid
	windows db.KpiWindowReader
	// loc anchors calendar range tokens (today/this_week/this_month) —
	// same timezone the rollup engine uses for span boundaries.
	loc *time.Location
	pb.UnimplementedAggregatorServiceServer
}

func NewAggregatorServer(org string, kpis []schema.KpiSpec, rollups db.RollupRepo,
	composer *performance.Composer, grid schema.Grid, windows db.KpiWindowReader,
	timezone string) (*AggregatorServer, error) {
	byKey := map[string]schema.KpiSpec{}
	for _, k := range kpis {
		byKey[k.Kpi] = k
	}

	loc := time.UTC

	if timezone != "" {
		parsed, err := time.LoadLocation(timezone)
		if err != nil {
			return nil, fmt.Errorf("loading aggregator timezone %q: %w", timezone, err)
		}

		loc = parsed
	}

	return &AggregatorServer{
		org:      org,
		kpis:     kpis,
		byKey:    byKey,
		rollups:  rollups,
		composer: composer,
		grid:     grid,
		windows:  windows,
		loc:      loc,
	}, nil
}

func (s *AggregatorServer) ListReports(ctx context.Context, req *pb.ListReportsRequest) (*pb.ListReportsResponse, error) {
	specs := s.composer.List()

	infos := make([]*pb.ReportInfo, 0, len(specs))

	for _, r := range specs {
		columns := make([]string, 0, len(r.Columns))
		for _, c := range r.Columns {
			columns = append(columns, c.Name)
		}

		attributes := make([]string, 0, len(r.Resource.Attributes))
		for _, a := range r.Resource.Attributes {
			attributes = append(attributes, a.Name)
		}

		infos = append(infos, &pb.ReportInfo{
			Report:     r.Report,
			Title:      r.Title,
			Resource:   r.Resource.Dataset,
			Columns:    columns,
			Attributes: attributes,
		})
	}

	return &pb.ListReportsResponse{Reports: infos}, nil
}

func (s *AggregatorServer) GetPerformanceReport(ctx context.Context, req *pb.GetPerformanceReportRequest) (*pb.GetPerformanceReportResponse, error) {
	// Reports accept the rolling filter tokens (last_24h/7d/30d) as well as the
	// calendar spans: the composer maps a rolling span to a precise trailing
	// window and falls back to the config window for anything else.
	span := strings.ToLower(req.Span)
	if !isRollingSpan(span) {
		validated, err := validateSpan(span)
		if err != nil {
			return nil, err
		}

		span = validated
	}

	report, err := s.composer.Compose(req.Report, span, req.Scope, int(req.Top))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "composing report: %v", err)
	}

	rows := make([]*pb.ReportRow, 0, len(report.Rows))

	for _, row := range report.Rows {
		cells := make([]*pb.ReportCell, 0, len(row.Cells))

		for _, cell := range row.Cells {
			pbCell := &pb.ReportCell{
				Column:    cell.Column,
				Value:     cell.Value,
				Unit:      cell.Unit,
				Symbol:    cell.Symbol,
				Format:    cell.Format,
				IsPartial: cell.IsPartial,
				Trend:     &pb.Trend{Direction: cell.Trend},
			}

			if cell.PrevValue != nil {
				pbCell.Trend.HasPrevious = true
				pbCell.Trend.PrevValue = *cell.PrevValue
			}

			if cell.ChangeAbs != nil {
				pbCell.Trend.ChangeAbs = *cell.ChangeAbs
			}

			if cell.ChangePct != nil {
				pbCell.Trend.ChangePct = *cell.ChangePct
			}

			if !cell.ComputedAt.IsZero() {
				pbCell.ComputedAt = cell.ComputedAt.UTC().Format(time.RFC3339)
			}

			cells = append(cells, pbCell)
		}

		rows = append(rows, &pb.ReportRow{
			EntityId:   row.EntityID,
			Attributes: row.Attributes,
			Cells:      cells,
			Status:     row.Status,
		})
	}

	return &pb.GetPerformanceReportResponse{
		Report: report.Report,
		Title:  report.Title,
		Span:   report.Span,
		Rows:   rows,
	}, nil
}

func (s *AggregatorServer) ListKpis(ctx context.Context, req *pb.ListKpisRequest) (*pb.ListKpisResponse, error) {
	infos := make([]*pb.KpiInfo, 0, len(s.kpis))

	for _, k := range s.kpis {
		infos = append(infos, &pb.KpiInfo{
			Kpi:               k.Kpi,
			Domain:            k.Domain,
			Scope:             k.Scope,
			Type:              k.Output.Type,
			Unit:              k.Output.Unit,
			Symbol:            k.Output.Symbol,
			PositiveDirection: k.PositiveDirection,
			Kind:              k.Kind,
			ScopeAgg:          k.ScopeAgg,
		})
	}

	sort.Slice(infos, func(i, j int) bool { return infos[i].Kpi < infos[j].Kpi })

	return &pb.ListKpisResponse{Kpis: infos}, nil
}

// --- legacy adapters over the Query planner ---

// legacyAgg maps this endpoint family's op param to a planner aggregation:
// empty, LAST and DELTA resolve to the KPI's kind default; component ops
// pass through as overrides.
func legacyAgg(kpi schema.KpiSpec, op string) (string, error) {
	op = strings.ToUpper(strings.TrimSpace(op))

	switch op {
	case "", "LAST", "DELTA":
		return kpi.DefaultReadOp(), nil
	}

	if !componentOps[op] {
		return "", status.Errorf(codes.InvalidArgument,
			"unknown op %s (use SUM, AVG, MIN, MAX, COUNT — or omit it: the KPI's kind picks the right one)", op)
	}

	return op, nil
}

// legacyGrain keeps this endpoint family's row shape: an unfiltered,
// ungrouped read returns one row per full KPI scope; anything else folds
// to filter ∪ group_by, exactly like Query.
func legacyGrain(kpi schema.KpiSpec, filter map[string]string, groupBy []string) []string {
	if len(filter) == 0 && len(groupBy) == 0 {
		return kpi.Scope
	}

	return grainKeys(kpi, filter, groupBy)
}

// legacyRangeToken maps a legacy span to a Query range token.
var legacyRangeToken = map[string]string{
	rollup.SpanDaily:   "today",
	rollup.SpanWeekly:  "this_week",
	rollup.SpanMonthly: "this_month",
}

func (s *AggregatorServer) legacyRange(span string) (queryRange, string, error) {
	span = strings.ToLower(span)

	if isRollingSpan(span) {
		rng, err := s.resolveQueryRange(span, "", "", time.Now().UTC())

		return rng, span, err
	}

	validated, err := validateSpan(span)
	if err != nil {
		return queryRange{}, "", err
	}

	rng, err := s.resolveQueryRange(legacyRangeToken[validated], "", "", time.Now().UTC())

	return rng, validated, err
}

func (s *AggregatorServer) GetKpis(ctx context.Context, req *pb.GetKpisRequest) (*pb.GetKpisResponse, error) {
	kpis, err := s.knownKpis(req.Keys)
	if err != nil {
		return nil, err
	}

	if err := validateScopeKeys(kpis, req.Scope); err != nil {
		return nil, err
	}

	if err := validateGroupBy(kpis, req.GroupBy); err != nil {
		return nil, err
	}

	rng, span, err := s.legacyRange(req.Span)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	values := make([]*pb.KpiValue, 0)

	for _, kpi := range kpis {
		op, err := legacyAgg(kpi, req.Op)
		if err != nil {
			return nil, err
		}

		rows, err := s.queryTotal(kpi, op, legacyGrain(kpi, req.Scope, req.GroupBy), req.Scope, rng, now)
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			values = append(values, queryRowToKpiValue(row, span, op, now))
		}
	}

	return &pb.GetKpisResponse{Values: values}, nil
}

func (s *AggregatorServer) GetKpiTimeSeries(ctx context.Context, req *pb.GetKpiTimeSeriesRequest) (*pb.GetKpiTimeSeriesResponse, error) {
	span, err := validateSpan(req.Span)
	if err != nil {
		return nil, err
	}

	from, to, err := parseRange(req.From, req.To)
	if err != nil {
		return nil, err
	}

	kpis, err := s.knownKpis(req.Keys)
	if err != nil {
		return nil, err
	}

	if err := validateScopeKeys(kpis, req.Scope); err != nil {
		return nil, err
	}

	if err := validateGroupBy(kpis, req.GroupBy); err != nil {
		return nil, err
	}

	gran := map[string]string{
		rollup.SpanDaily:   granDay,
		rollup.SpanWeekly:  granWeek,
		rollup.SpanMonthly: granMonth,
	}[span]

	rng := queryRange{from: from, to: to}
	now := time.Now().UTC()
	values := make([]*pb.KpiValue, 0)

	for _, kpi := range kpis {
		op, err := legacyAgg(kpi, req.Op)
		if err != nil {
			return nil, err
		}

		rows, err := s.querySeries(kpi, op, legacyGrain(kpi, req.Scope, req.GroupBy), req.Scope, gran, rng, now)
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			for _, pt := range row.Points {
				values = append(values, &pb.KpiValue{
					Kpi:        row.Kpi,
					Value:      pt.Value,
					Span:       span,
					Op:         op,
					From:       pt.From,
					To:         pt.To,
					Type:       row.Type,
					Unit:       row.Unit,
					Symbol:     row.Symbol,
					IsPartial:  pt.IsPartial,
					Scope:      row.Dims,
					Trend:      pt.Trend,
					ComputedAt: now.Format(time.RFC3339),
				})
			}
		}
	}

	return &pb.GetKpiTimeSeriesResponse{Values: values}, nil
}

func (s *AggregatorServer) GetKpiBreakdown(ctx context.Context, req *pb.GetKpiBreakdownRequest) (*pb.GetKpiBreakdownResponse, error) {
	kpi, ok := s.lookupKpi(req.Key)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown kpi %q (deployed: %s)",
			req.Key, strings.Join(s.kpiKeys(), ","))
	}

	if req.By == "" {
		return nil, status.Error(codes.InvalidArgument, "breakdown dimension 'by' is required")
	}

	if err := validateGroupBy([]schema.KpiSpec{kpi}, []string{req.By}); err != nil {
		return nil, err
	}

	if err := validateScopeKeys([]schema.KpiSpec{kpi}, req.Scope); err != nil {
		return nil, err
	}

	op, err := legacyAgg(kpi, req.Op)
	if err != nil {
		return nil, err
	}

	rng, span, err := s.legacyRange(req.Span)
	if err != nil {
		return nil, err
	}

	rows, err := s.queryTotal(kpi, op, grainKeys(kpi, req.Scope, []string{req.By}), req.Scope, rng, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	out := make([]*pb.BreakdownRow, 0, len(rows))

	var from, to string

	for _, row := range rows {
		if len(row.Points) == 0 {
			continue
		}

		pt := row.Points[0]
		from, to = pt.From, pt.To

		out = append(out, &pb.BreakdownRow{
			ScopeValue: row.Dims[req.By],
			Value:      pt.Value,
			Trend:      pt.Trend,
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Value > out[j].Value })

	if req.Top > 0 && int(req.Top) < len(out) {
		out = out[:req.Top]
	}

	return &pb.GetKpiBreakdownResponse{
		Kpi:  kpi.Kpi,
		Span: span,
		Op:   op,
		From: from,
		To:   to,
		Rows: out,
	}, nil
}

// queryRowToKpiValue flattens a single-point Query row into the legacy
// KpiValue shape.
func queryRowToKpiValue(row *pb.QueryRow, span, op string, now time.Time) *pb.KpiValue {
	v := &pb.KpiValue{
		Kpi:        row.Kpi,
		Span:       span,
		Op:         op,
		Type:       row.Type,
		Unit:       row.Unit,
		Symbol:     row.Symbol,
		Scope:      row.Dims,
		ComputedAt: now.UTC().Format(time.RFC3339),
	}

	if len(row.Points) > 0 {
		pt := row.Points[0]
		v.Value = pt.Value
		v.From = pt.From
		v.To = pt.To
		v.IsPartial = pt.IsPartial
		v.Trend = pt.Trend
	}

	return v
}

// lookupKpi resolves a requested key to its spec. Keys are case-insensitive:
// specs declare them upper-case (DATA_USAGE) while callers commonly send the
// lower-case form.
func (s *AggregatorServer) lookupKpi(key string) (schema.KpiSpec, bool) {
	if kpi, ok := s.byKey[key]; ok {
		return kpi, true
	}

	kpi, ok := s.byKey[strings.ToUpper(strings.TrimSpace(key))]

	return kpi, ok
}

// knownKpis filters requested keys down to the ones this deployment actually
// has a spec for. An unknown key is SKIPPED, not fatal: a caller asks for one
// list of keys to fill a whole tile row, so failing the request over a KPI
// that is not deployed here would blank every other tile alongside
// it. Consumers already degrade a missing value to "—", which is the intended
// contract. Single-key endpoints keep the NotFound — there, the unknown key is
// the entire answer.
//
// If NO key resolves there is nothing to degrade to, and every later check
// (scope validation above all) would report the empty spec set as if the
// caller's filters were wrong — so that case is an explicit error.
func (s *AggregatorServer) knownKpis(keys []string) ([]schema.KpiSpec, error) {
	out := make([]schema.KpiSpec, 0, len(keys))

	for _, key := range keys {
		kpi, ok := s.lookupKpi(key)
		if !ok {
			log.Warnf("skipping unknown kpi %q — no spec deployed in this configs/kpis (check the key)", key)

			continue
		}

		out = append(out, kpi)
	}

	if len(out) == 0 {
		return nil, status.Errorf(codes.NotFound,
			"no known kpi among %q — deployed keys: %s",
			strings.Join(keys, ","), strings.Join(s.kpiKeys(), ","))
	}

	return out, nil
}

// kpiKeys lists the deployed KPI keys, sorted, for error messages.
func (s *AggregatorServer) kpiKeys() []string {
	keys := make([]string, 0, len(s.byKey))
	for k := range s.byKey {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func validateSpan(span string) (string, error) {
	span = strings.ToLower(span)
	if span == "" {
		span = rollup.SpanDaily
	}

	for _, s := range rollup.Spans {
		if s == span {
			return span, nil
		}
	}

	return "", status.Errorf(codes.InvalidArgument, "unknown span %q", span)
}

func parseRange(fromStr, toStr string) (time.Time, time.Time, error) {
	now := time.Now().UTC()

	from := now.AddDate(0, -1, 0)
	to := now.AddDate(0, 0, 1)

	var err error

	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, status.Errorf(codes.InvalidArgument, "invalid from: %v", err)
		}
	}

	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			return time.Time{}, time.Time{}, status.Errorf(codes.InvalidArgument, "invalid to: %v", err)
		}
	}

	return from, to, nil
}
