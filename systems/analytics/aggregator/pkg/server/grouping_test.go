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
	"github.com/ukama/ukama/systems/analytics/schema"
)

func seriesScope(network, iccid string) string {
	return schema.CanonicalScope(map[string]string{
		"network_id": network,
		"iccid":      iccid,
	})
}

func dataUsageSpec() schema.KpiSpec {
	return schema.KpiSpec{
		Kpi:    "DATA_USAGE",
		Kind:   schema.KindFlow,
		Scope:  []string{"network_id", "site_id", "package_id", "sim_package_id", "iccid"},
		Output: schema.OutputSpec{Type: "float", Unit: "bytes", Symbol: "B"},
	}
}

func TestGroupRollups(t *testing.T) {
	spanStart := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	spanEnd := spanStart.AddDate(0, 0, 1)

	rows := []schema.KpiRollup{
		{Scope: seriesScope("net-a", "icc1"), SpanStart: spanStart, SpanEnd: spanEnd, Op: schema.RollupRowOp, Value: 100, Sum: 100, Count: 2, Min: 40, Max: 60, Last: 100},
		{Scope: seriesScope("net-a", "icc2"), SpanStart: spanStart, SpanEnd: spanEnd, Op: schema.RollupRowOp, Value: 50, Sum: 50, Count: 2, Min: 20, Max: 30, Last: 50},
		{Scope: seriesScope("net-b", "icc3"), SpanStart: spanStart, SpanEnd: spanEnd, Op: schema.RollupRowOp, Value: 7, Sum: 7, Count: 1, Min: 7, Max: 7, Last: 7},
	}

	groups := latestGroups(groupRollups(rows, []string{"network_id"}))
	assert.Len(t, groups, 2)

	byNet := map[string]*scopeGroup{}
	for _, g := range groups {
		byNet[schema.ParseScope(g.scope)["network_id"]] = g
	}

	sum, ok := byNet["net-a"].value("SUM")
	assert.True(t, ok)
	assert.Equal(t, float64(150), sum, "components fold across series")

	avg, ok := byNet["net-a"].value("AVG")
	assert.True(t, ok)
	assert.Equal(t, float64(37.5), avg, "weighted AVG: sum-of-sums / sum-of-counts")

	max, _ := byNet["net-a"].value("MAX")
	assert.Equal(t, float64(60), max)

	last, ok := byNet["net-a"].value("LAST")
	assert.True(t, ok)
	assert.Equal(t, float64(150), last,
		"LAST folds as the sum of row values (additive-gauge fold, used by Query)")
}

func TestGroupRollups_KeepsNewestSpanPerGroup(t *testing.T) {
	older := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)
	newer := older.AddDate(0, 0, 1)

	rows := []schema.KpiRollup{
		// icc1 last rolled up yesterday; icc2 today. Folding across the two
		// span_starts would mix periods — only the newest span may count.
		{Scope: seriesScope("net-a", "icc1"), SpanStart: older, Op: schema.RollupRowOp, Sum: 999, Count: 1},
		{Scope: seriesScope("net-a", "icc2"), SpanStart: newer, Op: schema.RollupRowOp, Sum: 50, Count: 1},
	}

	groups := latestGroups(groupRollups(rows, []string{"network_id"}))
	assert.Len(t, groups, 1)

	sum, _ := groups[0].value("SUM")
	assert.Equal(t, float64(50), sum, "only the newest span_start folds")
	assert.Equal(t, newer, groups[0].spanStart)
}

func TestGetKpisRolling_GroupByFoldsSeries(t *testing.T) {
	srv := mustServer(t,

		"org",
		[]schema.KpiSpec{dataUsageSpec()},
		nil, nil,
		schema.Grid{W: 5 * time.Minute},
		stubWindows{rows: []schema.KpiWindow{
			{KpiKey: "DATA_USAGE", Scope: seriesScope("net-a", "icc1"), WindowID: 1, Sum: 100, Count: 1, Min: 100, Max: 100, Value: 100},
			{KpiKey: "DATA_USAGE", Scope: seriesScope("net-a", "icc2"), WindowID: 1, Sum: 40, Count: 1, Min: 40, Max: 40, Value: 40},
			{KpiKey: "DATA_USAGE", Scope: seriesScope("net-b", "icc3"), WindowID: 1, Sum: 7, Count: 1, Min: 7, Max: 7, Value: 7},
		}},
	)

	resp, err := srv.GetKpis(context.TODO(), &pb.GetKpisRequest{
		Keys:    []string{"DATA_USAGE"},
		Span:    "last_24h",
		Op:      "SUM",
		GroupBy: []string{"network_id"},
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Values, 2, "one folded row per network")

	byNet := map[string]float64{}
	for _, v := range resp.Values {
		byNet[v.Scope["network_id"]] = v.Value
	}

	assert.Equal(t, float64(140), byNet["net-a"])
	assert.Equal(t, float64(7), byNet["net-b"])
}

func TestGetKpisRolling_GroupByWithFilter(t *testing.T) {
	srv := mustServer(t,

		"org",
		[]schema.KpiSpec{dataUsageSpec()},
		nil, nil,
		schema.Grid{W: 5 * time.Minute},
		stubWindows{rows: []schema.KpiWindow{
			{KpiKey: "DATA_USAGE", Scope: seriesScope("net-a", "icc1"), WindowID: 1, Sum: 100, Count: 1},
			{KpiKey: "DATA_USAGE", Scope: seriesScope("net-b", "icc3"), WindowID: 1, Sum: 7, Count: 1},
		}},
	)

	resp, err := srv.GetKpis(context.TODO(), &pb.GetKpisRequest{
		Keys:    []string{"DATA_USAGE"},
		Span:    "last_24h",
		Op:      "SUM",
		Scope:   map[string]string{"network_id": "net-a"},
		GroupBy: []string{"network_id"},
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Values, 1, "filter applies before the fold")
	assert.Equal(t, float64(100), resp.Values[0].Value)
}

func TestGetKpisRolling_FilterImpliesGrouping(t *testing.T) {
	srv := mustServer(t,

		"org",
		[]schema.KpiSpec{dataUsageSpec()},
		nil, nil,
		schema.Grid{W: 5 * time.Minute},
		stubWindows{rows: []schema.KpiWindow{
			{KpiKey: "DATA_USAGE", Scope: seriesScope("net-a", "icc1"), WindowID: 1, Sum: 100, Count: 1, Min: 100, Max: 100, Value: 100},
			{KpiKey: "DATA_USAGE", Scope: seriesScope("net-a", "icc2"), WindowID: 1, Sum: 40, Count: 1, Min: 40, Max: 40, Value: 40},
			{KpiKey: "DATA_USAGE", Scope: seriesScope("net-b", "icc3"), WindowID: 1, Sum: 7, Count: 1, Min: 7, Max: 7, Value: 7},
		}},
	)

	// No group_by: the filter itself is the requested grain — ONE row for
	// network net-a, its series folded.
	resp, err := srv.GetKpis(context.TODO(), &pb.GetKpisRequest{
		Keys:  []string{"DATA_USAGE"},
		Span:  "last_24h",
		Op:    "SUM",
		Scope: map[string]string{"network_id": "net-a"},
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Values, 1, "filter keys become the fold grain")
	assert.Equal(t, float64(140), resp.Values[0].Value)
	assert.Equal(t, "net-a", resp.Values[0].Scope["network_id"])
	// stub returns identical rows for current and previous period -> flat,
	// proving the trend comparison is computed at the folded grain.
	assert.Equal(t, "flat", resp.Values[0].Trend.Direction)

	// Filtering down to a single series keeps the classic per-row shape:
	// full scope, untouched by the fold machinery.
	resp, err = srv.GetKpis(context.TODO(), &pb.GetKpisRequest{
		Keys:  []string{"DATA_USAGE"},
		Span:  "last_24h",
		Op:    "SUM",
		Scope: map[string]string{"network_id": "net-b", "iccid": "icc3"},
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Values, 1)
	assert.Equal(t, float64(7), resp.Values[0].Value)
	assert.Equal(t, "icc3", resp.Values[0].Scope["iccid"], "single-member group passes through with full scope")
}

func TestGetKpis_RejectsMalformedReads(t *testing.T) {
	srv := mustServer(t,
		"org", []schema.KpiSpec{dataUsageSpec()}, nil, nil,
		schema.Grid{W: 5 * time.Minute}, stubWindows{})

	// Unknown scope filter key: 400, never a silent empty result.
	_, err := srv.GetKpis(context.TODO(), &pb.GetKpisRequest{
		Keys:  []string{"DATA_USAGE"},
		Span:  "last_24h",
		Scope: map[string]string{"node_id": "x"},
	})
	assert.Error(t, err)

	// Unknown group_by key.
	_, err = srv.GetKpis(context.TODO(), &pb.GetKpisRequest{
		Keys:    []string{"DATA_USAGE"},
		Span:    "last_24h",
		GroupBy: []string{"subscriber_id"},
	})
	assert.Error(t, err)
}
