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

// stubWindows returns the same rows regardless of range, so the current and
// previous windows aggregate identically (trend = flat) — enough to exercise
// the read path deterministically without a DB or a clock.
type stubWindows struct{ rows []schema.KpiWindow }

func (s stubWindows) WindowsInRange(_ string, _ string, _ int64, _ int64) ([]schema.KpiWindow, error) {
	return s.rows, nil
}

func (s stubWindows) WindowBounds(_ string, _ string) (int64, int64, bool, error) {
	if len(s.rows) == 0 {
		return 0, 0, false, nil
	}

	minID, maxID := s.rows[0].WindowID, s.rows[0].WindowID
	for _, r := range s.rows[1:] {
		if r.WindowID < minID {
			minID = r.WindowID
		}
		if r.WindowID > maxID {
			maxID = r.WindowID
		}
	}

	return minID, maxID, true, nil
}

func TestRollingOpValue(t *testing.T) {
	a := &windowAgg{sum: 300, count: 3, min: 50, max: 200, lastValue: 200}

	v, ok := rollingOpValue("SUM", a)
	assert.True(t, ok)
	assert.Equal(t, float64(300), v)

	v, ok = rollingOpValue("AVG", a) // weighted: sum/count, never avg-of-avg
	assert.True(t, ok)
	assert.Equal(t, float64(100), v)

	v, _ = rollingOpValue("COUNT", a)
	assert.Equal(t, float64(3), v)
	v, _ = rollingOpValue("MIN", a)
	assert.Equal(t, float64(50), v)
	v, _ = rollingOpValue("MAX", a)
	assert.Equal(t, float64(200), v)
	v, _ = rollingOpValue("LAST", a)
	assert.Equal(t, float64(200), v)

	_, ok = rollingOpValue("AVG", &windowAgg{count: 0})
	assert.False(t, ok, "AVG with no observations is undefined")
}

func TestScopeMatches(t *testing.T) {
	orgWide := schema.CanonicalScope(nil)
	netA := schema.CanonicalScope(map[string]string{"network_id": "net-a"})

	assert.True(t, scopeMatches(netA, nil), "no filter matches everything")
	assert.False(t, scopeMatches(orgWide, map[string]string{"network_id": "net-a"}),
		"org bucket (empty scope) is org-only, excluded from a network filter")
	assert.True(t, scopeMatches(netA, map[string]string{"network_id": "net-a"}))
	assert.False(t, scopeMatches(netA, map[string]string{"network_id": "net-b"}))
}

func TestGetKpisRolling_SumNetworkScope(t *testing.T) {
	netA := schema.CanonicalScope(map[string]string{"network_id": "net-a"})
	srv := mustServer(t,

		"org",
		[]schema.KpiSpec{{
			Kpi:    "REVENUE",
			Scope:  []string{"network_id"},
			Output: schema.OutputSpec{Type: "int", Unit: "cents", Symbol: "$"},
		}},
		nil, // rollups (unused on the rolling path)
		nil, // composer (unused)
		schema.Grid{W: 5 * time.Minute},
		stubWindows{rows: []schema.KpiWindow{
			{KpiKey: "REVENUE", Scope: netA, WindowID: 1, Sum: 100, Count: 1, Min: 100, Max: 100, Value: 100},
			{KpiKey: "REVENUE", Scope: netA, WindowID: 2, Sum: 200, Count: 1, Min: 200, Max: 200, Value: 200},
		}},
	)

	resp, err := srv.GetKpis(context.TODO(), &pb.GetKpisRequest{
		Keys:  []string{"REVENUE"},
		Span:  "last_7d",
		Scope: map[string]string{"network_id": "net-a"},
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Values, 1)

	v := resp.Values[0]
	assert.Equal(t, "REVENUE", v.Kpi)
	assert.Equal(t, float64(300), v.Value, "default op SUM over the window")
	assert.Equal(t, "SUM", v.Op)
	assert.Equal(t, "last_7d", v.Span)
	assert.True(t, v.IsPartial)
	assert.Equal(t, "net-a", v.Scope["network_id"], "network-scoped row")
	// current == previous (stub returns same rows) -> flat
	assert.Equal(t, "flat", v.Trend.Direction)
}
