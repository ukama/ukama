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
	"math"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/schema"
)

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
	rollups db.RollupRepo
	org     string
	loc     *time.Location
	sweep   time.Duration
	stop    chan struct{}
}

func NewEngine(grid schema.Grid, kpis []schema.KpiSpec, windows db.KpiWindowReader,
	rollups db.RollupRepo, org, timezone string, sweep time.Duration) (*Engine, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("loading rollup timezone %q: %w", timezone, err)
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
		stop:    make(chan struct{}),
	}, nil
}

// OnKpiComputed is the event fast path: recompute every span containing the
// window.
func (e *Engine) OnKpiComputed(kpiKey string, windowID int64) {
	kpi, ok := e.kpis[kpiKey]
	if !ok {
		log.Warnf("kpi.computed for unknown kpi %q (spec not deployed here?)", kpiKey)

		return
	}

	win := e.grid.Window(windowID)

	for _, span := range Spans {
		if err := e.RecomputeSpan(kpi, span, win.Start); err != nil {
			log.Errorf("recompute %s %s @ %s: %v", kpiKey, span, win.Start, err)
		}
	}
}

// StartSweeper first backfills every span covered by existing kpi_windows
// (heals stale rows and gaps after downtime — cheap and idempotent), then
// periodically recomputes the current spans to cover lost events and keep
// partial spans fresh.
func (e *Engine) StartSweeper() {
	e.backfillOnce()

	ticker := time.NewTicker(e.sweep)
	defer ticker.Stop()

	e.sweepOnce()

	for {
		select {
		case <-ticker.C:
			e.sweepOnce()
		case <-e.stop:
			return
		}
	}
}

func (e *Engine) Stop() {
	close(e.stop)
}

// backfillOnce re-materializes all spans intersecting each KPI's kpi_windows
// history. Runs at boot: after this, every historical span carries the
// current components-row shape regardless of what an older engine wrote.
func (e *Engine) backfillOnce() {
	now := time.Now().UTC()

	for _, kpi := range e.kpis {
		minID, _, ok, err := e.windows.WindowBounds(e.org, kpi.Kpi)
		if err != nil {
			log.Errorf("backfill: window bounds for %s: %v", kpi.Kpi, err)

			continue
		}

		if !ok {
			continue // no windows yet
		}

		start := e.grid.Window(minID).Start

		for _, span := range Spans {
			t := start

			for {
				spanStart, err := SpanStart(span, t, e.loc)
				if err != nil {
					break
				}

				if err := e.RecomputeSpan(kpi, span, spanStart); err != nil {
					log.Errorf("backfill: recompute %s %s @ %s: %v", kpi.Kpi, span, spanStart, err)
				}

				next, err := SpanEnd(span, spanStart)
				if err != nil || !next.Before(now) {
					break
				}

				t = next
			}
		}

		log.Infof("backfill: %s rollups re-materialized from %s", kpi.Kpi, start)
	}
}

// sweepOnce recomputes the current AND previous span for every KPI: current
// keeps partial spans fresh; previous closes out spans whose final
// kpi.computed event was lost (clears stale is_partial rows).
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

	spanEnd, err := SpanEnd(span, spanStart)
	if err != nil {
		return err
	}

	fromID, toID := e.windowRange(spanStart, spanEnd)

	rows, err := e.windows.WindowsInRange(e.org, kpi.Kpi, fromID, toID)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return nil // nothing observed in this span yet
	}

	now := time.Now().UTC()
	isPartial := spanEnd.After(now)

	// Aggregate components per scope.
	type agg struct {
		sum, count, min, max float64
		lastWindow           int64
		lastValue            float64
		meta                 schema.KpiWindow
	}

	byScope := map[string]*agg{}

	for _, row := range rows {
		a, ok := byScope[row.Scope]
		if !ok {
			a = &agg{min: math.Inf(1), max: math.Inf(-1), lastWindow: -1}
			byScope[row.Scope] = a
		}

		a.sum += row.Sum
		a.count += row.Count
		a.min = math.Min(a.min, row.Min)
		a.max = math.Max(a.max, row.Max)

		if row.WindowID > a.lastWindow {
			a.lastWindow = row.WindowID
			a.lastValue = row.Value
		}

		a.meta = row
	}

	upserts := make([]schema.KpiRollup, 0, len(byScope))

	for scope, a := range byScope {
		upserts = append(upserts, schema.KpiRollup{
			KpiKey:     kpi.Kpi,
			OrgID:      e.org,
			Scope:      scope,
			Span:       span,
			SpanStart:  spanStart.UTC(),
			SpanEnd:    spanEnd.UTC(),
			Op:         schema.RollupRowOp,
			Value:      defaultValue(kpi, a.sum, a.count, a.lastValue),
			Sum:        a.sum,
			Count:      a.count,
			Min:        a.min,
			Max:        a.max,
			Last:       a.lastValue,
			ValueType:  a.meta.ValueType,
			Unit:       a.meta.Unit,
			Symbol:     a.meta.Symbol,
			IsPartial:  isPartial,
			ComputedAt: now,
		})
	}

	return e.rollups.Upsert(upserts)
}

func (e *Engine) windowRange(spanStart, spanEnd time.Time) (int64, int64) {
	return e.grid.WindowAt(spanStart.UTC()).ID, e.grid.WindowAt(spanEnd.UTC()).ID
}

// defaultValue caches the KPI's kind-default aggregation on the row: flows
// sum over the span; ratio gauges take the weighted average; additive
// gauges take the latest level.
func defaultValue(kpi schema.KpiSpec, sum, count, last float64) float64 {
	if kpi.Kind == schema.KindGauge {
		if kpi.ScopeAgg == schema.ScopeAggAvg {
			if count == 0 {
				return 0
			}

			return sum / count
		}

		return last
	}

	return sum
}
