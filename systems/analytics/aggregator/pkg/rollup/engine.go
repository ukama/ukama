/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rollup

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// MinWindowRetention keeps every window a last_30d trend (60 days) reads.
const MinWindowRetention = 63 * 24 * time.Hour

// Windows may still land this long after a span closes, so an earlier rollup is not final.
const daySettle = time.Hour

// Engine materializes calendar-span rollups from kpi_windows components.
//
// One row per (kpi, scope, span, span_start) — schema.RollupRowOp — carrying
// the folded components (Sum/Count/Min/Max/Last). Every aggregation,
// group_by fold and trend is computed at READ time from these components;
// nothing is precomputed per op and no trend state is stored. Recomputation
// is idempotent (upserts).
type Engine struct {
	grid    schema.Grid
	kpis    map[string]schema.KpiSpec
	windows db.KpiWindowReader
	rollups db.RollupWriter
	org     string
	loc     *time.Location
	sweep   time.Duration
	flush   time.Duration
	keep    time.Duration
	stop    chan struct{}

	// work serializes statements so overlapping upserts cannot deadlock.
	work sync.Mutex

	dirtyMu sync.Mutex
	dirty   map[spanKey]struct{}

	ready     chan struct{}
	readyOnce sync.Once
}

type spanKey struct {
	kpi   string
	span  string
	start int64
}

// keep is the kpi_windows retention (0 = forever). Spans that start before
// now-keep may have lost windows, so they are never recomputed once they exist.
func NewEngine(grid schema.Grid, kpis []schema.KpiSpec, windows db.KpiWindowReader,
	rollups db.RollupWriter, org, timezone string, sweep, flush, keep time.Duration) (*Engine, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("loading rollup timezone %q: %w", timezone, err)
	}

	if flush <= 0 {
		return nil, fmt.Errorf("rollup flush interval must be positive, got %s", flush)
	}

	byKey := map[string]schema.KpiSpec{}
	for _, k := range kpis {
		byKey[k.Kpi] = k
	}

	return &Engine{
		grid:    grid,
		kpis:    byKey,
		windows: windows,
		rollups: rollups,
		org:     org,
		loc:     loc,
		sweep:   sweep,
		flush:   flush,
		keep:    keep,
		stop:    make(chan struct{}),
		dirty:   map[spanKey]struct{}{},
		ready:   make(chan struct{}),
	}, nil
}

// OnKpiComputed is the event fast path: queue the day containing the window.
func (e *Engine) OnKpiComputed(kpiKey string, windowID int64) {
	if _, ok := e.kpis[kpiKey]; !ok {
		log.Warnf("kpi.computed for unknown kpi %q (spec not deployed here?)", kpiKey)

		return
	}

	day, err := SpanStart(SpanDaily, e.grid.Window(windowID).Start, e.loc)
	if err != nil {
		return
	}

	e.dirtyMu.Lock()
	e.dirty[spanKey{kpi: kpiKey, span: SpanDaily, start: day.Unix()}] = struct{}{}
	e.dirtyMu.Unlock()
}

func (e *Engine) FlushOnce() {
	e.dirtyMu.Lock()
	batch := e.dirty
	e.dirty = map[spanKey]struct{}{}
	e.dirtyMu.Unlock()

	parents := map[spanKey]struct{}{}

	for key := range batch {
		day := time.Unix(key.start, 0)

		if err := e.RecomputeSpan(e.kpis[key.kpi], SpanDaily, day); err != nil {
			log.Errorf("recompute %s %s @ %s: %v", key.kpi, SpanDaily, day.UTC(), err)
		}

		for _, parent := range e.parentsOf(key.kpi, day) {
			parents[parent] = struct{}{}
		}
	}

	for key := range parents {
		start := time.Unix(key.start, 0)

		if err := e.RecomputeSpan(e.kpis[key.kpi], key.span, start); err != nil {
			log.Errorf("recompute %s %s @ %s: %v", key.kpi, key.span, start.UTC(), err)
		}
	}
}

func (e *Engine) parentsOf(kpiKey string, day time.Time) []spanKey {
	out := make([]spanKey, 0, len(Spans))

	for _, span := range Spans {
		if span == SpanDaily {
			continue
		}

		start, err := SpanStart(span, day, e.loc)
		if err != nil {
			continue
		}

		out = append(out, spanKey{kpi: kpiKey, span: span, start: start.Unix()})
	}

	return out
}

func (e *Engine) flushLoop() {
	ticker := time.NewTicker(e.flush)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.FlushOnce()
		case <-e.stop:
			return
		}
	}
}

// StartSweeper first backfills the spans downtime left missing or unfinished,
// then periodically recomputes the current spans to cover lost events and keep
// partial spans fresh.
func (e *Engine) StartSweeper() {
	go e.flushLoop()

	clean := e.backfillOnce()

	ticker := time.NewTicker(e.sweep)
	defer ticker.Stop()

	e.sweepOnce()

	for {
		select {
		case <-ticker.C:
			if !clean {
				clean = e.backfillOnce()
			}

			e.sweepOnce()
		case <-e.stop:
			return
		}
	}
}

func (e *Engine) Stop() {
	close(e.stop)
}

// Ready closes after the first backfill that finished without errors: from
// then on every span holding kpi_windows has a rollup, so pruning is safe.
func (e *Engine) Ready() <-chan struct{} {
	return e.ready
}

func (e *Engine) frozen(spanStart, now time.Time) bool {
	return e.keep > 0 && spanStart.Before(now.Add(-e.keep))
}

// backfillOnce materializes the spans of each KPI's kpi_windows history that
// have no rollup, are still partial, or were rolled up before they settled.
func (e *Engine) backfillOnce() bool {
	now := time.Now().UTC()
	clean := true

	for _, kpi := range e.kpis {
		minID, _, ok, err := e.windows.WindowBounds(e.org, kpi.Kpi)
		if err != nil {
			log.Errorf("backfill: window bounds for %s: %v", kpi.Kpi, err)

			clean = false

			continue
		}

		if !ok {
			continue // no windows yet
		}

		start := e.grid.Window(minID).Start
		touched := map[spanKey]struct{}{}
		done := 0

		for _, span := range Spans {
			n, err := e.backfillSpans(kpi, span, start, now, touched)
			if err != nil {
				log.Errorf("backfill: %s %s: %v", kpi.Kpi, span, err)

				clean = false
			}

			done += n
		}

		log.Infof("backfill: %s re-materialized %d spans from %s", kpi.Kpi, done, start)
	}

	if clean {
		e.readyOnce.Do(func() { close(e.ready) })
	}

	return clean
}

func (e *Engine) backfillSpans(kpi schema.KpiSpec, span string, start, now time.Time,
	touched map[spanKey]struct{}) (int, error) {
	states, err := e.rollups.SpanStates(e.org, kpi.Kpi, span)
	if err != nil {
		return 0, err
	}

	byStart := make(map[int64]db.SpanState, len(states))
	for _, st := range states {
		byStart[st.SpanStart.Unix()] = st
	}

	var firstErr error

	done := 0
	t := start

	for {
		spanStart, err := SpanStart(span, t, e.loc)
		if err != nil {
			return done, err
		}

		spanEnd, err := SpanEnd(span, spanStart)
		if err != nil {
			return done, err
		}

		key := spanKey{kpi: kpi.Kpi, span: span, start: spanStart.Unix()}
		st, exists := byStart[key.start]
		_, wasTouched := touched[key]

		if e.needsBackfill(spanStart, spanEnd, now, st, exists, wasTouched) {
			if err := e.recompute(kpi, span, spanStart); err != nil {
				if firstErr == nil {
					firstErr = fmt.Errorf("recompute @ %s: %w", spanStart, err)
				}
			} else {
				done++
			}

			if span == SpanDaily {
				for _, parent := range e.parentsOf(kpi.Kpi, spanStart) {
					touched[parent] = struct{}{}
				}
			}
		}

		if !spanEnd.Before(now) {
			return done, firstErr
		}

		t = spanEnd
	}
}

// A frozen span is final once it exists: its windows may be partly pruned.
func (e *Engine) needsBackfill(spanStart, spanEnd, now time.Time, st db.SpanState, exists, touched bool) bool {
	if !exists {
		return true
	}

	if e.frozen(spanStart, now) {
		return false
	}

	return touched || st.IsPartial || st.ComputedAt.Before(spanEnd.Add(daySettle))
}

// sweepOnce recomputes the current AND previous span for every KPI: current
// keeps partial spans fresh; previous closes out spans whose final
// kpi.computed event was lost (clears stale is_partial rows). Spans lists
// daily first, so weekly and monthly fold days that were just refreshed.
func (e *Engine) sweepOnce() {
	now := time.Now().UTC()

	for _, kpi := range e.kpis {
		for _, span := range Spans {
			if err := e.RecomputeSpan(kpi, span, now); err != nil {
				log.Errorf("sweeper: recompute %s %s: %v", kpi.Kpi, span, err)
			}

			spanStart, err := SpanStart(span, now, e.loc)
			if err != nil {
				continue
			}

			prevStart, err := PrevSpanStart(span, spanStart)
			if err != nil {
				continue
			}

			if err := e.RecomputeSpan(kpi, span, prevStart); err != nil {
				log.Errorf("sweeper: recompute previous %s %s: %v", kpi.Kpi, span, err)
			}
		}
	}
}

// RecomputeSpan rebuilds the components row for every scope of the span
// containing t.
func (e *Engine) RecomputeSpan(kpi schema.KpiSpec, span string, t time.Time) error {
	spanStart, err := SpanStart(span, t, e.loc)
	if err != nil {
		return err
	}

	if e.frozen(spanStart, time.Now().UTC()) {
		log.Debugf("skipping %s %s @ %s: span starts before kpi_windows retention", kpi.Kpi, span, spanStart)

		return nil
	}

	return e.recompute(kpi, span, spanStart)
}

func (e *Engine) recompute(kpi schema.KpiSpec, span string, spanStart time.Time) error {
	spanEnd, err := SpanEnd(span, spanStart)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	isPartial := spanEnd.After(now)
	mode := valueMode(kpi)

	e.work.Lock()
	defer e.work.Unlock()

	if span == SpanDaily {
		fromID, toID := e.windowRange(spanStart, spanEnd)

		return e.rollups.RollupFromWindows(e.org, kpi.Kpi, mode, spanStart, spanEnd, fromID, toID, isPartial, now)
	}

	return e.rollups.RollupFromDaily(e.org, kpi.Kpi, span, mode, spanStart, spanEnd, isPartial, now)
}

func (e *Engine) windowRange(spanStart, spanEnd time.Time) (int64, int64) {
	return e.grid.WindowAt(spanStart.UTC()).ID, e.grid.WindowAt(spanEnd.UTC()).ID
}

// valueMode picks the KPI's kind-default aggregation cached on the row: flows
// sum over the span; ratio gauges take the weighted average; additive
// gauges take the latest level.
func valueMode(kpi schema.KpiSpec) db.ValueMode {
	if kpi.Kind == schema.KindGauge {
		if kpi.ScopeAgg == schema.ScopeAggAvg {
			return db.ValueAvg
		}

		return db.ValueLast
	}

	return db.ValueSum
}
