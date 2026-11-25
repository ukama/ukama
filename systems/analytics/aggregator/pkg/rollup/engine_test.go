/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rollup

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/schema"
)

const (
	testOrg = "ukama"
	testKpi = "REVENUE"
)

var testGrid = schema.Grid{W: time.Minute}

type rollupID struct {
	scope string
	span  string
	start int64
}

type fakeStore struct {
	windows     []schema.KpiWindow
	rollups     map[rollupID]schema.KpiRollup
	fromWindows []time.Time
	fromDaily   []rollupID
	failDaily   bool
}

func newFakeStore() *fakeStore {
	return &fakeStore{rollups: map[rollupID]schema.KpiRollup{}}
}

func (f *fakeStore) addWindow(scope string, at time.Time, value float64) {
	f.windows = append(f.windows, schema.KpiWindow{
		KpiKey: testKpi, OrgID: testOrg, Scope: scope, WindowID: testGrid.WindowAt(at).ID,
		Value: value, Sum: value, Count: 1, Min: value, Max: value, Unit: "cents",
	})
}

func (f *fakeStore) AggregateWindows(_, _ string, _, _ int64) ([]db.WindowAgg, error) {
	return nil, nil
}

func (f *fakeStore) WindowBounds(_, _ string) (int64, int64, bool, error) {
	if len(f.windows) == 0 {
		return 0, 0, false, nil
	}

	lo, hi := f.windows[0].WindowID, f.windows[0].WindowID
	for _, w := range f.windows {
		lo, hi = min(lo, w.WindowID), max(hi, w.WindowID)
	}

	return lo, hi, true, nil
}

type fold struct {
	sum, count, min, max, last float64
	order                      int64
	seen                       bool
}

func (a *fold) add(sum, count, lo, hi, last float64, order int64) {
	if !a.seen {
		a.min, a.max, a.order = lo, hi, order-1
		a.seen = true
	}

	a.sum += sum
	a.count += count
	a.min, a.max = min(a.min, lo), max(a.max, hi)

	if order > a.order {
		a.order, a.last = order, last
	}
}

func (f *fakeStore) store(folds map[string]*fold, span string, mode db.ValueMode, start, end time.Time,
	isPartial bool, now time.Time) {
	for scope, a := range folds {
		value := a.sum

		switch mode {
		case db.ValueAvg:
			value = a.sum / a.count
		case db.ValueLast:
			value = a.last
		}

		f.rollups[rollupID{scope, span, start.Unix()}] = schema.KpiRollup{
			KpiKey: testKpi, OrgID: testOrg, Scope: scope, Span: span, SpanStart: start, SpanEnd: end,
			Op: schema.RollupRowOp, Value: value, Sum: a.sum, Count: a.count, Min: a.min, Max: a.max,
			Last: a.last, IsPartial: isPartial, ComputedAt: now,
		}
	}
}

func (f *fakeStore) RollupFromWindows(_, _ string, mode db.ValueMode, start, end time.Time,
	fromID, toID int64, isPartial bool, now time.Time) error {
	if f.failDaily {
		return errors.New("db down")
	}

	f.fromWindows = append(f.fromWindows, start)
	folds := map[string]*fold{}

	for _, w := range f.windows {
		if w.WindowID < fromID || w.WindowID >= toID {
			continue
		}

		if folds[w.Scope] == nil {
			folds[w.Scope] = &fold{}
		}

		folds[w.Scope].add(w.Sum, w.Count, w.Min, w.Max, w.Value, w.WindowID)
	}

	f.store(folds, SpanDaily, mode, start, end, isPartial, now)

	return nil
}

func (f *fakeStore) RollupFromDaily(_, _, span string, mode db.ValueMode, start, end time.Time,
	isPartial bool, now time.Time) error {
	f.fromDaily = append(f.fromDaily, rollupID{span: span, start: start.Unix()})
	folds := map[string]*fold{}

	for id, r := range f.rollups {
		if id.span != SpanDaily || r.SpanStart.Before(start) || !r.SpanStart.Before(end) {
			continue
		}

		if folds[id.scope] == nil {
			folds[id.scope] = &fold{}
		}

		folds[id.scope].add(r.Sum, r.Count, r.Min, r.Max, r.Last, r.SpanStart.Unix())
	}

	f.store(folds, span, mode, start, end, isPartial, now)

	return nil
}

func (f *fakeStore) SpanStates(_, _, span string) ([]db.SpanState, error) {
	byStart := map[int64]db.SpanState{}

	for id, r := range f.rollups {
		if id.span != span {
			continue
		}

		st, ok := byStart[id.start]
		if !ok {
			st = db.SpanState{SpanStart: r.SpanStart, ComputedAt: r.ComputedAt}
		}

		st.IsPartial = st.IsPartial || r.IsPartial
		if r.ComputedAt.Before(st.ComputedAt) {
			st.ComputedAt = r.ComputedAt
		}

		byStart[id.start] = st
	}

	out := make([]db.SpanState, 0, len(byStart))
	for _, st := range byStart {
		out = append(out, st)
	}

	return out, nil
}

func newTestEngine(t *testing.T, store *fakeStore, kind string, keep time.Duration) *Engine {
	t.Helper()

	e, err := NewEngine(testGrid, []schema.KpiSpec{{Kpi: testKpi, Kind: kind}}, store, store,
		testOrg, "UTC", time.Minute, time.Second, keep)
	require.NoError(t, err)

	return e
}

func day(daysAgo int) time.Time {
	now := time.Now().UTC()

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -daysAgo)
}

func spanStart(t *testing.T, span string, at time.Time) int64 {
	t.Helper()

	s, err := SpanStart(span, at, time.UTC)
	require.NoError(t, err)

	return s.Unix()
}

func TestOnKpiComputed_CoalescesEventsPerDay(t *testing.T) {
	store := newFakeStore()
	e := newTestEngine(t, store, schema.KindFlow, 0)
	d := day(1)

	for i := 0; i < 500; i++ {
		at := d.Add(time.Duration(i) * time.Minute)
		store.addWindow("{}", at, 1)
		e.OnKpiComputed(testKpi, testGrid.WindowAt(at).ID)
	}

	assert.Empty(t, store.fromWindows, "events only queue; nothing runs until a flush")

	e.FlushOnce()

	assert.Equal(t, []time.Time{d}, store.fromWindows, "500 events for one day fold it once")
	assert.Len(t, store.fromDaily, 2, "one weekly and one monthly refresh")

	e.FlushOnce()

	assert.Len(t, store.fromWindows, 1, "a drained queue does no work")
}

func TestOnKpiComputed_IgnoresUnknownKpi(t *testing.T) {
	store := newFakeStore()
	e := newTestEngine(t, store, schema.KindFlow, 0)

	e.OnKpiComputed("NOT_DEPLOYED", testGrid.WindowAt(day(0)).ID)
	e.FlushOnce()

	assert.Empty(t, store.fromWindows)
}

func TestFlush_WeeklyAndMonthlyFoldDailyRows(t *testing.T) {
	store := newFakeStore()
	e := newTestEngine(t, store, schema.KindFlow, 0)

	// Two days in the same ISO week and month, anchored on that week's Monday.
	monday := time.Unix(spanStart(t, SpanWeekly, day(7)), 0).UTC()
	if monday.Day() > 27 {
		monday = monday.AddDate(0, 0, -7)
	}

	tuesday := monday.AddDate(0, 0, 1)

	store.addWindow("{}", monday.Add(time.Hour), 100)
	store.addWindow("{}", monday.Add(2*time.Hour), 300)
	store.addWindow("{}", tuesday.Add(time.Hour), 50)

	for _, w := range store.windows {
		e.OnKpiComputed(testKpi, w.WindowID)
	}

	e.FlushOnce()

	daily := store.rollups[rollupID{"{}", SpanDaily, monday.Unix()}]
	assert.Equal(t, 400.0, daily.Sum)
	assert.Equal(t, 300.0, daily.Last)

	for _, span := range []string{SpanWeekly, SpanMonthly} {
		row, ok := store.rollups[rollupID{"{}", span, spanStart(t, span, monday)}]
		require.True(t, ok, span)
		assert.Equal(t, 450.0, row.Sum, span)
		assert.Equal(t, 450.0, row.Value, span)
		assert.Equal(t, 3.0, row.Count, span)
		assert.Equal(t, 50.0, row.Min, span)
		assert.Equal(t, 300.0, row.Max, span)
		assert.Equal(t, 50.0, row.Last, "%s takes the latest day's level", span)
	}
}

func TestBackfill_SkipsSettledSpansAndRedoesTheRest(t *testing.T) {
	store := newFakeStore()
	e := newTestEngine(t, store, schema.KindFlow, 0)

	settled, cutShort, missing := day(4), day(3), day(2)
	for _, d := range []time.Time{settled, cutShort, missing} {
		store.addWindow("{}", d.Add(time.Hour), 10)
	}

	final := schema.KpiRollup{Sum: 10, Count: 1, Min: 10, Max: 10, Last: 10}

	row := final
	row.SpanStart, row.ComputedAt = settled, settled.AddDate(0, 0, 1).Add(2*time.Hour)
	store.rollups[rollupID{"{}", SpanDaily, settled.Unix()}] = row

	// Rolled up one minute after the day closed: late windows may be missing.
	row = final
	row.SpanStart, row.ComputedAt = cutShort, cutShort.AddDate(0, 0, 1).Add(time.Minute)
	store.rollups[rollupID{"{}", SpanDaily, cutShort.Unix()}] = row

	assert.True(t, e.backfillOnce())

	assert.NotContains(t, store.fromWindows, settled)
	assert.Contains(t, store.fromWindows, cutShort)
	assert.Contains(t, store.fromWindows, missing)
	assert.Contains(t, store.fromWindows, day(0), "today is always partial")

	_, ok := store.rollups[rollupID{"{}", SpanMonthly, spanStart(t, SpanMonthly, missing)}]
	assert.True(t, ok, "the month holding a recomputed day is refreshed")

	select {
	case <-e.Ready():
	default:
		t.Fatal("a clean backfill must open the pruner gate")
	}

	assert.True(t, e.backfillOnce())

	for _, d := range []time.Time{cutShort, missing} {
		n := 0

		for _, got := range store.fromWindows {
			if got.Equal(d) {
				n++
			}
		}

		assert.Equal(t, 1, n, "a second pass leaves %s alone", d)
	}
}

func TestBackfill_FrozenSpanIsFinalOnceItExists(t *testing.T) {
	store := newFakeStore()
	e := newTestEngine(t, store, schema.KindFlow, 48*time.Hour)

	frozen := day(5)
	store.addWindow("{}", frozen.Add(time.Hour), 10)
	store.rollups[rollupID{"{}", SpanDaily, frozen.Unix()}] = schema.KpiRollup{
		SpanStart: frozen, IsPartial: true, ComputedAt: frozen,
	}

	assert.True(t, e.backfillOnce())
	assert.NotContains(t, store.fromWindows, frozen, "its windows may be pruned")
}

func TestBackfill_ErrorKeepsPrunerGateClosed(t *testing.T) {
	store := newFakeStore()
	store.failDaily = true
	e := newTestEngine(t, store, schema.KindFlow, 0)

	store.addWindow("{}", day(1).Add(time.Hour), 10)

	assert.False(t, e.backfillOnce())

	select {
	case <-e.Ready():
		t.Fatal("pruning must wait for a clean backfill")
	default:
	}
}

func TestValueMode(t *testing.T) {
	assert.Equal(t, db.ValueSum, valueMode(schema.KpiSpec{Kind: schema.KindFlow}))
	assert.Equal(t, db.ValueLast, valueMode(schema.KpiSpec{Kind: schema.KindGauge}))
	assert.Equal(t, db.ValueAvg, valueMode(schema.KpiSpec{Kind: schema.KindGauge, ScopeAgg: schema.ScopeAggAvg}))
}

func TestNewEngine_RejectsZeroFlushInterval(t *testing.T) {
	store := newFakeStore()

	_, err := NewEngine(testGrid, nil, store, store, testOrg, "UTC", time.Minute, 0, 0)
	assert.Error(t, err)
}
