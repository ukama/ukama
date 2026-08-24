/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package performance

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// spanDaily is the rollup span the report aggregates over its rolling window.
// Daily granularity keeps the read cheap (~one row per scope per day) while
// still covering the full config window.
const spanDaily = "daily"

// Composer builds resource performance reports at read time: entity rows
// from the raw zone's change-log state + KPI cells aggregated over a rolling
// report window (config, default 8 weeks) of daily rollups. No new state.
type Composer struct {
	org     string
	reports map[string]schema.ReportSpec
	state   db.RawStateReader
	rollups db.RollupRepo
	windows db.KpiWindowReader
	grid    schema.Grid
	window  time.Duration // fallback window when the read span isn't a rolling one
}

// Cell is one KPI value attached to an entity row.
type Cell struct {
	Column     string
	Value      float64
	Unit       string
	Symbol     string
	Format     string
	IsPartial  bool
	Trend      string
	PrevValue  *float64
	ChangeAbs  *float64
	ChangePct  *float64
	ComputedAt time.Time
	hasData    bool
}

// Row is one entity in the report.
type Row struct {
	EntityID   string
	Attributes map[string]string
	Cells      []Cell
	Status     string
}

// Report is the composed result.
type Report struct {
	Report string
	Title  string
	Span   string
	Rows   []Row
}

func NewComposer(org string, reports []schema.ReportSpec, state db.RawStateReader,
	rollups db.RollupRepo, windows db.KpiWindowReader, grid schema.Grid,
	window time.Duration) *Composer {
	byKey := map[string]schema.ReportSpec{}
	for _, r := range reports {
		byKey[r.Report] = r
	}

	return &Composer{
		org:     org,
		reports: byKey,
		state:   state,
		rollups: rollups,
		windows: windows,
		grid:    grid,
		window:  window,
	}
}

func (c *Composer) List() []schema.ReportSpec {
	out := make([]schema.ReportSpec, 0, len(c.reports))
	for _, r := range c.reports {
		out = append(out, r)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Report < out[j].Report })

	return out
}

// Compose builds the report for a span with optional scope filters (e.g.
// network_id) and top-N truncation.
func (c *Composer) Compose(report, span string, scopeFilter map[string]string, top int) (*Report, error) {
	spec, ok := c.reports[report]
	if !ok {
		return nil, fmt.Errorf("unknown report %q", report)
	}

	entities, err := c.entities(spec, scopeFilter)
	if err != nil {
		return nil, err
	}

	// Cells are aggregated over a trailing window of DAILY rollups. The window
	// follows the UI filter (last 24h / 7 days / 30 days) so the report matches
	// the rest of the page; an empty/unknown span falls back to the configured
	// default. Trend compares the current window against the one before it.
	window := c.window
	rolling := false
	if d, ok := rollingWindow(span); ok {
		window = d
		rolling = true
	}

	now := time.Now().UTC()
	from := now.Add(-window)
	prevFrom := now.Add(-2 * window)

	// Rolling spans aggregate the raw kpi_windows over the exact window (precise
	// last 24h / 7 days / 30 days, matching the values endpoint); a non-rolling
	// span falls back to daily rollups over the configured window.
	var fromID, toID, prevFromID int64
	if rolling {
		fromID = c.grid.WindowAt(from).ID
		toID = c.grid.WindowAt(now).ID + 1
		prevFromID = c.grid.WindowAt(prevFrom).ID
	}

	columnRows := make([]map[string][]schema.KpiRollup, len(spec.Columns))
	columnPrev := make([]map[string]float64, len(spec.Columns))

	for i, col := range spec.Columns {
		var curr, prev []schema.KpiRollup

		var err error

		if rolling {
			curr, err = c.windowRollups(col.Kpi, fromID, toID, scopeFilter)
		} else {
			curr, err = c.rollups.Range(c.org, col.Kpi, spanDaily, from, now, scopeFilter)
		}

		if err != nil {
			return nil, fmt.Errorf("column %s: %w", col.Name, err)
		}

		if rolling {
			prev, err = c.windowRollups(col.Kpi, prevFromID, fromID, scopeFilter)
		} else {
			prev, err = c.rollups.Range(c.org, col.Kpi, spanDaily, prevFrom, from, scopeFilter)
		}

		if err != nil {
			return nil, fmt.Errorf("column %s (previous window): %w", col.Name, err)
		}

		columnRows[i] = groupByEntity(curr, spec.RowScope)

		columnPrev[i] = map[string]float64{}
		for id, rs := range groupByEntity(prev, spec.RowScope) {
			if v, ok := foldValue(col.Op, rs); ok {
				columnPrev[i][id] = v
			}
		}
	}

	rows := make([]Row, 0, len(entities))

	for entityID, fields := range entities {
		row := Row{
			EntityID:   entityID,
			Attributes: map[string]string{},
			Cells:      make([]Cell, 0, len(spec.Columns)),
		}

		ruleValues := map[string]interface{}{}

		for _, attr := range spec.Resource.Attributes {
			v := fields[attr.Field]
			row.Attributes[attr.Name] = fmt.Sprintf("%v", orEmpty(v))
			ruleValues[attr.Name] = v
		}

		for i, col := range spec.Columns {
			prevValue, hasPrev := columnPrev[i][entityID]
			cell := composeCell(col, columnRows[i][entityID], prevValue, hasPrev)
			row.Cells = append(row.Cells, cell)
			ruleValues[col.Name] = cell.Value
		}

		status, err := EvaluateStatus(spec.Status, ruleValues)
		if err != nil {
			return nil, err
		}
		row.Status = status

		rows = append(rows, row)
	}

	sortRows(rows, spec)

	if top > 0 && top < len(rows) {
		rows = rows[:top]
	}

	return &Report{Report: spec.Report, Title: spec.Title, Span: reportWindowLabel(window), Rows: rows}, nil
}

// entities loads the latest state of the resource dataset, applying the
// org-level/network match rule when a network filter is present.
func (c *Composer) entities(spec schema.ReportSpec, scopeFilter map[string]string) (map[string]map[string]interface{}, error) {
	records, err := c.state.LatestState(c.org, spec.Resource.Dataset)
	if err != nil {
		return nil, fmt.Errorf("loading resource %s: %w", spec.Resource.Dataset, err)
	}

	networkID := scopeFilter["network_id"]

	out := make(map[string]map[string]interface{}, len(records))

	for _, rec := range records {
		fields := map[string]interface{}{}
		if err := json.Unmarshal([]byte(rec.Fields), &fields); err != nil {
			return nil, fmt.Errorf("resource %s entity %s fields: %w", spec.Resource.Dataset, rec.EntityKey, err)
		}

		// Org-level entities (empty match field) apply to every network.
		if networkID != "" && spec.Resource.NetworkMatch != "" {
			match := fmt.Sprintf("%v", orEmpty(fields[spec.Resource.NetworkMatch]))
			if match != "" && match != networkID {
				continue
			}
		}

		out[rec.EntityKey] = fields
	}

	return out, nil
}

// groupByEntity buckets rollup rows (across scopes and days in the window) by
// the report's row-scope entity id.
func groupByEntity(rows []schema.KpiRollup, rowScope string) map[string][]schema.KpiRollup {
	out := map[string][]schema.KpiRollup{}

	for _, r := range rows {
		if id := schema.ParseScope(r.Scope)[rowScope]; id != "" {
			out[id] = append(out[id], r)
		}
	}

	return out
}

// foldValue reduces all of an entity's rows in the window (over days and
// scopes) to a single value per the column op, from the rows' COMPONENTS:
// SUM/COUNT fold exactly; AVG is weighted (Σsum/Σcount, never a mean of
// daily values); MIN/MAX take the extremes; LAST sums, per scope, the
// latest row's level (additive-gauge semantics).
func foldValue(op string, rows []schema.KpiRollup) (float64, bool) {
	if len(rows) == 0 {
		return 0, false
	}

	switch strings.ToUpper(op) {
	case "MIN":
		v := rows[0].Min
		for _, r := range rows[1:] {
			if r.Min < v {
				v = r.Min
			}
		}

		return v, true
	case "MAX":
		v := rows[0].Max
		for _, r := range rows[1:] {
			if r.Max > v {
				v = r.Max
			}
		}

		return v, true
	case "AVG":
		sum, count := 0.0, 0.0
		for _, r := range rows {
			sum += r.Sum
			count += r.Count
		}

		if count == 0 {
			return 0, false
		}

		return sum / count, true
	case "LAST":
		latestPerScope := map[string]schema.KpiRollup{}
		for _, r := range rows {
			if cur, ok := latestPerScope[r.Scope]; !ok || r.SpanStart.After(cur.SpanStart) {
				latestPerScope[r.Scope] = r
			}
		}

		v := 0.0
		for _, r := range latestPerScope {
			v += r.Last
		}

		return v, true
	case "COUNT":
		v := 0.0
		for _, r := range rows {
			v += r.Count
		}

		return v, true
	default: // SUM
		v := 0.0
		for _, r := range rows {
			v += r.Sum
		}

		return v, true
	}
}

// composeCell builds one report cell from an entity's rows in the current
// window, with a trend vs the same entity's value in the preceding window.
func composeCell(col schema.ReportColumn, rows []schema.KpiRollup, prevValue float64, hasPrev bool) Cell {
	cell := Cell{Column: col.Name, Format: col.Format}

	if len(rows) == 0 {
		return cell // zero cell: entity with no KPI data still lists
	}

	cell.hasData = true
	cell.Unit = rows[0].Unit
	cell.Symbol = rows[0].Symbol
	cell.Value, _ = foldValue(col.Op, rows)

	for _, r := range rows {
		if r.IsPartial {
			cell.IsPartial = true
		}

		if r.ComputedAt.After(cell.ComputedAt) {
			cell.ComputedAt = r.ComputedAt
		}
	}

	if !hasPrev {
		cell.Trend = "new"

		return cell
	}

	prev := prevValue
	changeAbs := cell.Value - prev
	cell.PrevValue = &prev
	cell.ChangeAbs = &changeAbs

	if prev == 0 {
		if cell.Value == 0 {
			cell.Trend = "flat"
		} else {
			cell.Trend = "na" // percent undefined; change_abs still served
		}

		return cell
	}

	changePct := changeAbs / math.Abs(prev) * 100
	cell.ChangePct = &changePct

	switch {
	case changePct == 0:
		cell.Trend = "flat"
	case changePct > 0:
		cell.Trend = "up"
	default:
		cell.Trend = "down"
	}

	return cell
}

// windowRollups reads raw kpi_windows for a KPI over [fromID,toID), keeps the
// rows matching the report's scope filter, and adapts each to a KpiRollup so
// the existing group/fold path aggregates it. This gives the report the same
// precise trailing-window aggregation the values endpoint uses for rolling
// spans (exact totals, not a coarse daily bucket).
func (c *Composer) windowRollups(kpiKey string, fromID, toID int64, filter map[string]string) ([]schema.KpiRollup, error) {
	rows, err := c.windows.WindowsInRange(c.org, kpiKey, fromID, toID)
	if err != nil {
		return nil, err
	}

	out := make([]schema.KpiRollup, 0, len(rows))
	for _, w := range rows {
		if !windowScopeMatches(w.Scope, filter) {
			continue
		}

		out = append(out, schema.KpiRollup{
			Scope:      w.Scope,
			Value:      w.Value,
			Sum:        w.Sum,
			Count:      w.Count,
			Min:        w.Min,
			Max:        w.Max,
			Last:       w.Value, // a window's value IS its level at that moment
			SpanStart:  c.grid.Window(w.WindowID).Start,
			Unit:       w.Unit,
			Symbol:     w.Symbol,
			ComputedAt: w.ComputedAt,
		})
	}

	return out, nil
}

// windowScopeMatches reports whether a kpi_windows scope satisfies the report's
// scope filter (e.g. network_id). An empty filter matches everything.
func windowScopeMatches(scope string, filter map[string]string) bool {
	if len(filter) == 0 {
		return true
	}

	m := schema.ParseScope(scope)
	for k, v := range filter {
		if m[k] != v {
			return false
		}
	}

	return true
}

// rollingWindow maps the UI filter span (last_24h / last_7d / last_30d) to the
// trailing report-window duration. Returns false for an empty/unknown span so
// the caller keeps the configured default.
func rollingWindow(span string) (time.Duration, bool) {
	switch strings.ToLower(span) {
	case "last_24h":
		return 24 * time.Hour, true
	case "last_7d":
		return 7 * 24 * time.Hour, true
	case "last_30d":
		return 30 * 24 * time.Hour, true
	default:
		return 0, false
	}
}

// reportWindowLabel renders the report window as a short token for the API
// response, e.g. 1344h -> "8w", 72h -> "3d".
func reportWindowLabel(w time.Duration) string {
	if weeks := int(w.Hours()) / (24 * 7); weeks > 0 {
		return fmt.Sprintf("%dw", weeks)
	}

	return fmt.Sprintf("%dd", int(w.Hours())/24)
}

func sortRows(rows []Row, spec schema.ReportSpec) {
	column := spec.Sort.Column

	idx := -1
	for i, col := range spec.Columns {
		if col.Name == column {
			idx = i

			break
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		if idx >= 0 {
			a, b := rows[i].Cells[idx].Value, rows[j].Cells[idx].Value
			if a != b {
				if spec.Sort.Desc {
					return a > b
				}

				return a < b
			}
		}

		return rows[i].EntityID < rows[j].EntityID // stable fallback
	})
}

func orEmpty(v interface{}) interface{} {
	if v == nil {
		return ""
	}

	return v
}
