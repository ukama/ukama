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
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/performance"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/rollup"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// AggregatorServer is the generic KPI read API over kpi_rollups plus the
// performance-report composer. No per-KPI/per-report code: adding either
// requires no changes here.
type AggregatorServer struct {
	org      string
	kpis     []schema.KpiSpec
	byKey    map[string]schema.KpiSpec
	rollups  db.RollupRepo
	composer *performance.Composer
	// grid + windows back the rolling-window read path (last_24h/7d/30d),
	// which aggregates kpi_windows on read instead of reading precomputed
	// calendar-span rollups. See rolling.go.
	grid    schema.Grid
	windows db.KpiWindowReader
	pb.UnimplementedAggregatorServiceServer
}

func NewAggregatorServer(org string, kpis []schema.KpiSpec, rollups db.RollupRepo,
	composer *performance.Composer, grid schema.Grid, windows db.KpiWindowReader) *AggregatorServer {
	byKey := map[string]schema.KpiSpec{}
	for _, k := range kpis {
		byKey[k.Kpi] = k
	}

	return &AggregatorServer{
		org:      org,
		kpis:     kpis,
		byKey:    byKey,
		rollups:  rollups,
		composer: composer,
		grid:     grid,
		windows:  windows,
	}
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
			RollupOps:         k.RollupOps,
			Type:              k.Output.Type,
			Unit:              k.Output.Unit,
			Symbol:            k.Output.Symbol,
			PositiveDirection: k.PositiveDirection,
		})
	}

	sort.Slice(infos, func(i, j int) bool { return infos[i].Kpi < infos[j].Kpi })

	return &pb.ListKpisResponse{Kpis: infos}, nil
}

func (s *AggregatorServer) GetKpis(ctx context.Context, req *pb.GetKpisRequest) (*pb.GetKpisResponse, error) {
	// Rolling windows (last_24h/7d/30d) are computed on read from kpi_windows;
	// calendar spans (daily/weekly/monthly) read precomputed rollups below.
	if isRollingSpan(req.Span) {
		return s.getKpisRolling(req)
	}

	span, err := validateSpan(req.Span)
	if err != nil {
		return nil, err
	}

	values := make([]*pb.KpiValue, 0)

	for _, key := range req.Keys {
		kpi, op, err := s.resolveKpiOp(key, req.Op)
		if err != nil {
			return nil, err
		}

		rows, err := s.rollups.Latest(s.org, kpi.Kpi, span, op, req.Scope)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "reading rollups: %v", err)
		}

		for _, row := range rows {
			values = append(values, toKpiValue(row))
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

	values := make([]*pb.KpiValue, 0)

	for _, key := range req.Keys {
		kpi, op, err := s.resolveKpiOp(key, req.Op)
		if err != nil {
			return nil, err
		}

		rows, err := s.rollups.Range(s.org, kpi.Kpi, span, op, from, to, req.Scope)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "reading rollups: %v", err)
		}

		for _, row := range rows {
			values = append(values, toKpiValue(row))
		}
	}

	return &pb.GetKpiTimeSeriesResponse{Values: values}, nil
}

func (s *AggregatorServer) GetKpiBreakdown(ctx context.Context, req *pb.GetKpiBreakdownRequest) (*pb.GetKpiBreakdownResponse, error) {
	span, err := validateSpan(req.Span)
	if err != nil {
		return nil, err
	}

	kpi, op, err := s.resolveKpiOp(req.Key, req.Op)
	if err != nil {
		return nil, err
	}

	if req.By == "" {
		return nil, status.Error(codes.InvalidArgument, "breakdown dimension 'by' is required")
	}

	rows, err := s.rollups.Latest(s.org, kpi.Kpi, span, op, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "reading rollups: %v", err)
	}

	out := make([]*pb.BreakdownRow, 0, len(rows))

	var from, to string

	for _, row := range rows {
		scope := schema.ParseScope(row.Scope)

		v, ok := scope[req.By]
		if !ok {
			continue
		}

		from = row.SpanStart.UTC().Format(time.RFC3339)
		to = row.SpanEnd.UTC().Format(time.RFC3339)

		out = append(out, &pb.BreakdownRow{
			ScopeValue: v,
			Value:      row.Value,
			Trend:      toTrend(row),
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

func (s *AggregatorServer) resolveKpiOp(key, op string) (schema.KpiSpec, string, error) {
	kpi, ok := s.byKey[key]
	if !ok {
		return schema.KpiSpec{}, "", status.Errorf(codes.NotFound, "unknown kpi %q", key)
	}

	allowed := map[string]bool{}
	for _, o := range kpi.RollupOps {
		allowed[strings.ToUpper(o)] = true
	}

	if op == "" {
		// Default op: LAST when allowed (most intuitive "current value"),
		// otherwise the spec's first op.
		if allowed["LAST"] {
			return kpi, "LAST", nil
		}

		return kpi, strings.ToUpper(kpi.RollupOps[0]), nil
	}

	op = strings.ToUpper(op)
	if !allowed[op] {
		return schema.KpiSpec{}, "", status.Errorf(codes.InvalidArgument,
			"op %s not allowed for kpi %s (allowed: %s)", op, key, strings.Join(kpi.RollupOps, ","))
	}

	return kpi, op, nil
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

func toKpiValue(row schema.KpiRollup) *pb.KpiValue {
	return &pb.KpiValue{
		Kpi:        row.KpiKey,
		Value:      row.Value,
		Span:       row.Span,
		Op:         row.Op,
		From:       row.SpanStart.UTC().Format(time.RFC3339),
		To:         row.SpanEnd.UTC().Format(time.RFC3339),
		Type:       row.ValueType,
		Unit:       row.Unit,
		Symbol:     row.Symbol,
		IsPartial:  row.IsPartial,
		Scope:      schema.ParseScope(row.Scope),
		Trend:      toTrend(row),
		ComputedAt: row.ComputedAt.UTC().Format(time.RFC3339),
	}
}

func toTrend(row schema.KpiRollup) *pb.Trend {
	t := &pb.Trend{Direction: row.Trend}

	if row.PrevValue != nil {
		t.HasPrevious = true
		t.PrevValue = *row.PrevValue
	}

	if row.ChangeAbs != nil {
		t.ChangeAbs = *row.ChangeAbs
	}

	if row.ChangePct != nil {
		t.ChangePct = *row.ChangePct
	}

	return t
}
