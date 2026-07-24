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
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// Engine recomputes span rollups from kpi_windows components. All ops are
// exact: AVG is weighted (sum of sums / sum of counts) — never an average of
// window averages. Recomputation is idempotent (upserts).
type Engine struct {
	grid    schema.Grid
	kpis    map[string]schema.KpiSpec
	windows db.KpiWindowReader
	rollups db.RollupRepo
	org     string
	loc     *time.Location
	flatPct float64
	sweep   time.Duration
	stop    chan struct{}
}

func NewEngine(grid schema.Grid, kpis []schema.KpiSpec, windows db.KpiWindowReader,
	rollups db.RollupRepo, org, timezone string, flatPct float64, sweep time.Duration) (*Engine, error) {
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
		flatPct: flatPct,
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

// StartSweeper periodically recomputes the current spans for every KPI:
// covers lost events and keeps partial (current) spans fresh.
func (e *Engine) StartSweeper() {
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

// RecomputeSpan rebuilds all rollup rows (every scope, every allowed op) for
// the span containing t, then refreshes the following span's trend (a
// corrected span dirties its successor's comparison).
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

	upserts := make([]schema.KpiRollup, 0, len(byScope)*len(kpi.RollupOps))

	for scope, a := range byScope {
		for _, op := range kpi.RollupOps {
			op = strings.ToUpper(op)

			value, ok := opValue(op, a.sum, a.count, a.min, a.max, a.lastValue)
			if !ok {
				continue
			}

			row := schema.KpiRollup{
				KpiKey:     kpi.Kpi,
				OrgID:      e.org,
				Scope:      scope,
				Span:       span,
				SpanStart:  spanStart.UTC(),
				SpanEnd:    spanEnd.UTC(),
				Op:         op,
				Value:      value,
				ValueType:  a.meta.ValueType,
				Unit:       a.meta.Unit,
				Symbol:     a.meta.Symbol,
				IsPartial:  isPartial,
				ComputedAt: now,
			}

			if err := e.applyTrend(&row, spanStart); err != nil {
				return err
			}

			upserts = append(upserts, row)
		}
	}

	if err := e.rollups.Upsert(upserts); err != nil {
		return err
	}

	// A recomputed span invalidates its successor's trend comparison.
	return e.refreshNextTrend(kpi, span, spanStart)
}

func (e *Engine) windowRange(spanStart, spanEnd time.Time) (int64, int64) {
	return e.grid.WindowAt(spanStart.UTC()).ID, e.grid.WindowAt(spanEnd.UTC()).ID
}

func opValue(op string, sum, count, min, max, last float64) (float64, bool) {
	switch op {
	case "SUM":
		return sum, true
	case "COUNT":
		return count, true
	case "AVG":
		if count == 0 {
			return 0, false
		}

		return sum / count, true
	case "MIN":
		return min, true
	case "MAX":
		return max, true
	case "LAST":
		return last, true
	case "DELTA":
		// span usage from a cumulative counter; clamp resets to 0
		if d := max - min; d > 0 {
			return d, true
		}

		return 0, true
	default:
		return 0, false
	}
}

// applyTrend fills the trend fields by comparing with the previous same-span
// row (same kpi/scope/op): up|down|flat|new|na.
func (e *Engine) applyTrend(row *schema.KpiRollup, spanStart time.Time) error {
	prevStart, err := PrevSpanStart(row.Span, spanStart)
	if err != nil {
		return err
	}

	prev, err := e.rollups.Get(row.OrgID, row.KpiKey, row.Scope, row.Span, prevStart.UTC(), row.Op)
	if err != nil {
		return err
	}

	if prev == nil {
		row.Trend = "new"

		return nil
	}

	prevValue := prev.Value
	changeAbs := row.Value - prevValue

	row.PrevValue = &prevValue
	row.ChangeAbs = &changeAbs

	if prevValue == 0 {
		if row.Value == 0 {
			row.Trend = "flat"
		} else {
			row.Trend = "na" // percent undefined; change_abs still served
		}

		return nil
	}

	changePct := changeAbs / math.Abs(prevValue) * 100
	row.ChangePct = &changePct

	switch {
	case math.Abs(changePct) < e.flatPct:
		row.Trend = "flat"
	case changePct > 0:
		row.Trend = "up"
	default:
		row.Trend = "down"
	}

	return nil
}

// refreshNextTrend re-applies trend fields on the following span's rows so a
// corrected previous value cascades forward.
func (e *Engine) refreshNextTrend(kpi schema.KpiSpec, span string, spanStart time.Time) error {
	nextStart, err := SpanEnd(span, spanStart)
	if err != nil {
		return err
	}

	scopes, err := e.rollups.ScopesSeen(e.org, kpi.Kpi, span)
	if err != nil {
		return err
	}

	updates := make([]schema.KpiRollup, 0)

	for _, scope := range scopes {
		for _, op := range kpi.RollupOps {
			op = strings.ToUpper(op)

			next, err := e.rollups.Get(e.org, kpi.Kpi, scope, span, nextStart.UTC(), op)
			if err != nil {
				return err
			}

			if next == nil {
				continue
			}

			next.PrevValue, next.ChangeAbs, next.ChangePct = nil, nil, nil

			if err := e.applyTrend(next, nextStart); err != nil {
				return err
			}

			updates = append(updates, *next)
		}
	}

	return e.rollups.Upsert(updates)
}
