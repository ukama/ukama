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
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/rollup"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// AggregatorServer is the generic KPI read API over kpi_rollups. No per-KPI
// code: adding a KPI requires no changes here.
type AggregatorServer struct {
	org     string
	kpis    []schema.KpiSpec
	byKey   map[string]schema.KpiSpec
	rollups db.RollupRepo
	pb.UnimplementedAggregatorServiceServer
}

func NewAggregatorServer(org string, kpis []schema.KpiSpec, rollups db.RollupRepo) *AggregatorServer {
	byKey := map[string]schema.KpiSpec{}
	for _, k := range kpis {
		byKey[k.Kpi] = k
	}

	return &AggregatorServer{
		org:     org,
		kpis:    kpis,
		byKey:   byKey,
		rollups: rollups,
	}
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
