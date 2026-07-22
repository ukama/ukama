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
	window  time.Duration // trailing report window (config, default 8 weeks)
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
	rollups db.RollupRepo, window time.Duration) *Composer {
	byKey := map[string]schema.ReportSpec{}
	for _, r := range reports {
		byKey[r.Report] = r
	}

	return &Composer{
		org:     org,
		reports: byKey,
		state:   state,
		rollups: rollups,
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

	// Cells are aggregated over a rolling report window (config, default 8
	// weeks) of DAILY rollups — enough history for stable per-entity stats,
	// consistent with the values endpoint's rolling reads (just longer and
	// coarser). Trend compares the current window against the one before it.
	now := time.Now().UTC()
	from := now.Add(-c.window)
	prevFrom := now.Add(-2 * c.window)

	columnRows := make([]map[string][]schema.KpiRollup, len(spec.Columns))
	columnPrev := make([]map[string]float64, len(spec.Columns))

	for i, col := range spec.Columns {
		curr, err := c.rollups.Range(c.org, col.Kpi, spanDaily, col.Op, from, now, scopeFilter)
		if err != nil {
			return nil, fmt.Errorf("column %s: %w", col.Name, err)
		}

		prev, err := c.rollups.Range(c.org, col.Kpi, spanDaily, col.Op, prevFrom, from, scopeFilter)
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

	// span is ignored: the report window is config-driven, not the UI filter.
	// The returned Span reports the actual window used (e.g. "8w").
	_ = span

	return &Report{Report: spec.Report, Title: spec.Title, Span: reportWindowLabel(c.window), Rows: rows}, nil
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

// foldValue reduces all of an entity's daily rows in the window (over days and
// scopes) to a single value per the column op. SUM/COUNT sum exactly; MIN/MAX
// take the extreme; LAST takes the most recent day's value; AVG is a mean
// across the daily values (a documented approximation, as before).
func foldValue(op string, rows []schema.KpiRollup) (float64, bool) {
	if len(rows) == 0 {
		return 0, false
	}

	switch strings.ToUpper(op) {
	case "MIN":
		v := rows[0].Value
		for _, r := range rows[1:] {
			if r.Value < v {
				v = r.Value
			}
		}

		return v, true
	case "MAX":
		v := rows[0].Value
		for _, r := range rows[1:] {
			if r.Value > v {
				v = r.Value
			}
		}

		return v, true
	case "AVG":
		sum := 0.0
		for _, r := range rows {
			sum += r.Value
		}

		return sum / float64(len(rows)), true
	case "LAST":
		v, latest := rows[0].Value, rows[0].SpanStart
		for _, r := range rows[1:] {
			if r.SpanStart.After(latest) {
				latest, v = r.SpanStart, r.Value
			}
		}

		return v, true
	default: // SUM, COUNT: additive across days and scopes
		sum := 0.0
		for _, r := range rows {
			sum += r.Value
		}

		return sum, true
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
