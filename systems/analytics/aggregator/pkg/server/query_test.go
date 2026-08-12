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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	pb "github.com/ukama/ukama/systems/analytics/aggregator/pb/gen"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/performance"
	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/rollup"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// mustServer builds an AggregatorServer with the UTC calendar for tests.
func mustServer(t *testing.T, org string, kpis []schema.KpiSpec, rollups db.RollupRepo,
	composer *performance.Composer, grid schema.Grid, windows db.KpiWindowReader) *AggregatorServer {
	t.Helper()

	srv, err := NewAggregatorServer(org, kpis, rollups, composer, grid, windows, "UTC")
	if err != nil {
		t.Fatalf("building server: %v", err)
	}

	return srv
}

// stubRollups serves canned kpi_rollups rows, honoring the kpi/span/op/
// span_start-range/scope-filter selection the planner relies on.
type stubRollups struct{ rows []schema.KpiRollup }

func (s stubRollups) Upsert([]schema.KpiRollup) error { return nil }

func (s stubRollups) Get(string, string, string, string, time.Time, string) (*schema.KpiRollup, error) {
	return nil, nil
}

func (s stubRollups) ScopesSeen(string, string, string) ([]string, error) { return nil, nil }

func (s stubRollups) Latest(org, kpi, span string, filter map[string]string) ([]schema.KpiRollup, error) {
	out := []schema.KpiRollup{}

	for _, r := range s.rows {
		if r.KpiKey == kpi && r.Span == span && r.Op == schema.RollupRowOp && scopeMatches(r.Scope, filter) {
			out = append(out, r)
		}
	}

	return out, nil
}

func (s stubRollups) Range(org, kpi, span string, from, to time.Time, filter map[string]string) ([]schema.KpiRollup, error) {
	out := []schema.KpiRollup{}

	for _, r := range s.rows {
		if r.KpiKey == kpi && r.Span == span && r.Op == schema.RollupRowOp &&
			!r.SpanStart.Before(from) && r.SpanStart.Before(to) &&
			scopeMatches(r.Scope, filter) {
			out = append(out, r)
		}
	}

	return out, nil
}

func flowUsageSpec() schema.KpiSpec {
	s := dataUsageSpec()
	s.Kind = schema.KindFlow

	return s
}

func monthlyRow(kpi, scope string, start time.Time, sum, count float64) schema.KpiRollup {
	return schema.KpiRollup{
		KpiKey: kpi, Scope: scope, Span: rollup.SpanMonthly, Op: schema.RollupRowOp,
		SpanStart: start, SpanEnd: start.AddDate(0, 1, 0),
		Value: sum, Sum: sum, Count: count, Min: sum, Max: sum, Last: sum,
	}
}

func TestQuery_FlowTotal_CalendarRange(t *testing.T) {
	now := time.Now().UTC()

	thisMonth, err := rollup.SpanStart(rollup.SpanMonthly, now, time.UTC)
	assert.NoError(t, err)

	prevMonth, err := rollup.PrevSpanStart(rollup.SpanMonthly, thisMonth)
	assert.NoError(t, err)

	srv := mustServer(t, "org", []schema.KpiSpec{flowUsageSpec()},
		stubRollups{rows: []schema.KpiRollup{
			monthlyRow("DATA_USAGE", seriesScope("net-a", "icc1"), thisMonth, 300, 3),
			monthlyRow("DATA_USAGE", seriesScope("net-a", "icc2"), thisMonth, 200, 2),
			monthlyRow("DATA_USAGE", seriesScope("net-b", "icc3"), thisMonth, 7, 1),
			monthlyRow("DATA_USAGE", seriesScope("net-a", "icc1"), prevMonth, 100, 4),
		}},
		nil, schema.Grid{W: 5 * time.Minute}, stubWindows{})

	// "How much data did network net-a use this month?" — no op, no span
	// tokens, no group_by. One number.
	resp, err := srv.Query(context.TODO(), &pb.QueryRequest{
		Kpis:   []string{"DATA_USAGE"},
		Filter: map[string]string{"network_id": "net-a"},
		Range:  "this_month",
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Rows, 1, "filter folds to one row")

	row := resp.Rows[0]
	assert.Equal(t, "DATA_USAGE", row.Kpi)
	assert.Equal(t, map[string]string{"network_id": "net-a"}, row.Dims)
	assert.Len(t, row.Points, 1)
	assert.Equal(t, float64(500), row.Points[0].Value, "series fold: 300+200")
	assert.Equal(t, "up", row.Points[0].Trend.Direction, "vs previous month's 100")
	assert.Equal(t, float64(100), row.Points[0].Trend.PrevValue)
	assert.Equal(t, "bytes", row.Unit)
}

func TestQuery_GroupByBreakdown_TopN(t *testing.T) {
	now := time.Now().UTC()

	thisMonth, err := rollup.SpanStart(rollup.SpanMonthly, now, time.UTC)
	assert.NoError(t, err)

	srv := mustServer(t, "org", []schema.KpiSpec{flowUsageSpec()},
		stubRollups{rows: []schema.KpiRollup{
			monthlyRow("DATA_USAGE", seriesScope("net-a", "icc1"), thisMonth, 300, 3),
			monthlyRow("DATA_USAGE", seriesScope("net-a", "icc2"), thisMonth, 200, 2),
			monthlyRow("DATA_USAGE", seriesScope("net-b", "icc3"), thisMonth, 7, 1),
		}},
		nil, schema.Grid{W: 5 * time.Minute}, stubWindows{})

	// "Top network by usage this month" — breakdown is just group_by+sort+top.
	resp, err := srv.Query(context.TODO(), &pb.QueryRequest{
		Kpis:    []string{"DATA_USAGE"},
		GroupBy: []string{"network_id"},
		Range:   "this_month",
		Sort:    "-value",
		Top:     1,
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Rows, 1)
	assert.Equal(t, "net-a", resp.Rows[0].Dims["network_id"])
	assert.Equal(t, float64(500), resp.Rows[0].Points[0].Value)
}

func TestQuery_GaugeSum_RollingRange(t *testing.T) {
	gauge := schema.KpiSpec{
		Kpi:      "ACTIVE_CUSTOMERS",
		Kind:     schema.KindGauge,
		ScopeAgg: schema.ScopeAggSum,
		Scope:    []string{"network_id"},
		Output:   schema.OutputSpec{Type: "int", Unit: "count"},
	}

	netA := schema.CanonicalScope(map[string]string{"network_id": "net-a"})
	netB := schema.CanonicalScope(map[string]string{"network_id": "net-b"})

	srv := mustServer(t, "org", []schema.KpiSpec{gauge}, stubRollups{},
		nil, schema.Grid{W: 5 * time.Minute},
		stubWindows{rows: []schema.KpiWindow{
			// per network: older window then latest — the gauge's current
			// level is the LATEST value, not a sum over time.
			{KpiKey: "ACTIVE_CUSTOMERS", Scope: netA, WindowID: 1, Value: 9, Sum: 9, Count: 1},
			{KpiKey: "ACTIVE_CUSTOMERS", Scope: netA, WindowID: 2, Value: 10, Sum: 10, Count: 1},
			{KpiKey: "ACTIVE_CUSTOMERS", Scope: netB, WindowID: 2, Value: 5, Sum: 5, Count: 1},
		}},
	)

	// "How many active customers right now?" — org total = sum of each
	// network's latest level (10 + 5), NOT a sum over windows.
	resp, err := srv.Query(context.TODO(), &pb.QueryRequest{
		Kpis:  []string{"ACTIVE_CUSTOMERS"},
		Range: "last_24h",
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Rows, 1)
	assert.Equal(t, float64(15), resp.Rows[0].Points[0].Value)
	assert.Empty(t, resp.Rows[0].Dims, "no filter/group_by -> org-wide total")
}

func TestQuery_SeriesDaily(t *testing.T) {
	day1 := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)
	day2 := day1.AddDate(0, 0, 1)

	daily := func(scope string, start time.Time, sum float64) schema.KpiRollup {
		return schema.KpiRollup{
			KpiKey: "DATA_USAGE", Scope: scope, Span: rollup.SpanDaily, Op: schema.RollupRowOp,
			SpanStart: start, SpanEnd: start.AddDate(0, 0, 1),
			Value: sum, Sum: sum, Count: 1, Min: sum, Max: sum, Last: sum,
		}
	}

	srv := mustServer(t, "org", []schema.KpiSpec{flowUsageSpec()},
		stubRollups{rows: []schema.KpiRollup{
			daily(seriesScope("net-a", "icc1"), day1, 100),
			daily(seriesScope("net-a", "icc2"), day1, 50),
			daily(seriesScope("net-a", "icc1"), day2, 400),
		}},
		nil, schema.Grid{W: 5 * time.Minute}, stubWindows{})

	resp, err := srv.Query(context.TODO(), &pb.QueryRequest{
		Kpis:        []string{"DATA_USAGE"},
		Filter:      map[string]string{"network_id": "net-a"},
		From:        day1.Format(time.RFC3339),
		To:          day2.AddDate(0, 0, 1).Format(time.RFC3339),
		Granularity: "day",
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Rows, 1, "one row for the network, points per day")

	points := resp.Rows[0].Points
	assert.Len(t, points, 2)
	assert.Equal(t, float64(150), points[0].Value, "day 1: both series folded")
	assert.Equal(t, "new", points[0].Trend.Direction)
	assert.Equal(t, float64(400), points[1].Value)
	assert.Equal(t, "up", points[1].Trend.Direction, "vs the previous bucket")
}

func TestQuery_Validation(t *testing.T) {
	srv := mustServer(t, "org", []schema.KpiSpec{flowUsageSpec()},
		stubRollups{}, nil, schema.Grid{W: 5 * time.Minute}, stubWindows{})

	cases := []*pb.QueryRequest{
		{}, // no kpis
		{Kpis: []string{"DATA_USAGE"}, Range: "nope"},  // unknown range token
		{Kpis: []string{"DATA_USAGE"}, Agg: "last"},    // LAST is not a caller-facing agg
		{Kpis: []string{"DATA_USAGE"}, Agg: "median"},  // unknown agg
		{Kpis: []string{"DATA_USAGE"}, From: "2026-x"}, // bad custom range
		{Kpis: []string{"DATA_USAGE"}, Granularity: "hour"},
		{Kpis: []string{"DATA_USAGE"}, Filter: map[string]string{"node_id": "x"}},
		{Kpis: []string{"DATA_USAGE"}, GroupBy: []string{"subscriber_id"}},
		{Kpis: []string{"DATA_USAGE"}, Range: "today", From: "2026-08-01T00:00:00Z", To: "2026-08-02T00:00:00Z"},
	}

	for i, req := range cases {
		_, err := srv.Query(context.TODO(), req)
		assert.Error(t, err, "case %d should be rejected", i)
	}
}

func TestQuery_KpiKeysAreCaseInsensitive(t *testing.T) {
	srv := mustServer(t, "org", []schema.KpiSpec{flowUsageSpec()},
		stubRollups{}, nil, schema.Grid{W: 5 * time.Minute}, stubWindows{})

	// A lower-case key used to resolve to NO spec, and the empty spec set then
	// made every scope filter look invalid — the caller was told their filter
	// was wrong when the key was what did not match.
	for _, key := range []string{"DATA_USAGE", "data_usage", "Data_Usage", " data_usage "} {
		_, err := srv.Query(context.TODO(), &pb.QueryRequest{
			Kpis:   []string{key},
			Filter: map[string]string{"iccid": "8910303228701258322"},
		})
		assert.NoError(t, err, "key %q should resolve", key)
	}
}

func TestQuery_EveryScopeDimensionIsAcceptedAsAFilter(t *testing.T) {
	srv := mustServer(t, "org", []schema.KpiSpec{flowUsageSpec()},
		stubRollups{}, nil, schema.Grid{W: 5 * time.Minute}, stubWindows{})

	for _, dim := range flowUsageSpec().Scope {
		_, err := srv.Query(context.TODO(), &pb.QueryRequest{
			Kpis:   []string{"data_usage"},
			Filter: map[string]string{dim: "x"},
		})
		assert.NoError(t, err, "filter %q is a DATA_USAGE scope dimension", dim)

		_, err = srv.Query(context.TODO(), &pb.QueryRequest{
			Kpis:    []string{"data_usage"},
			GroupBy: []string{dim},
		})
		assert.NoError(t, err, "group_by %q is a DATA_USAGE scope dimension", dim)
	}
}

func TestQuery_UnknownKpiIsNotBlamedOnTheFilter(t *testing.T) {
	srv := mustServer(t, "org", []schema.KpiSpec{flowUsageSpec()},
		stubRollups{}, nil, schema.Grid{W: 5 * time.Minute}, stubWindows{})

	_, err := srv.Query(context.TODO(), &pb.QueryRequest{
		Kpis:   []string{"NO_SUCH_KPI"},
		Filter: map[string]string{"iccid": "x"},
	})
	assert.ErrorContains(t, err, "no known kpi")

	// A known key alongside an unknown one still answers.
	_, err = srv.Query(context.TODO(), &pb.QueryRequest{
		Kpis:   []string{"NO_SUCH_KPI", "DATA_USAGE"},
		Filter: map[string]string{"iccid": "x"},
	})
	assert.NoError(t, err)
}
