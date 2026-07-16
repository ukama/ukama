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
	"sort"
	"time"

	"github.com/ukama/ukama/systems/analytics/aggregator/pkg/db"
	"github.com/ukama/ukama/systems/analytics/schema"
)

// Composer builds resource performance reports at read time: entity rows
// from the raw zone's change-log state + KPI cells from the latest
// available rollup rows (same read path as /kpis/values). No new state.
type Composer struct {
	org     string
	reports map[string]schema.ReportSpec
	state   db.RawStateReader
	rollups db.RollupRepo
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
	rollups db.RollupRepo) *Composer {
	byKey := map[string]schema.ReportSpec{}
	for _, r := range reports {
		byKey[r.Report] = r
	}

	return &Composer{
		org:     org,
		reports: byKey,
		state:   state,
		rollups: rollups,
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

	// Latest KPI rows per column, keyed by the row-scope entity id.
	columnRows := make([]map[string][]schema.KpiRollup, len(spec.Columns))

	for i, col := range spec.Columns {
		rows, err := c.rollups.Latest(c.org, col.Kpi, span, col.Op, scopeFilter)
		if err != nil {
			return nil, fmt.Errorf("column %s: %w", col.Name, err)
		}

		byEntity := map[string][]schema.KpiRollup{}

		for _, row := range rows {
			scope := schema.ParseScope(row.Scope)
			if id := scope[spec.RowScope]; id != "" {
				byEntity[id] = append(byEntity[id], row)
			}
		}

		columnRows[i] = byEntity
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
			cell := combine(col, columnRows[i][entityID])
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

	return &Report{Report: spec.Report, Title: spec.Title, Span: span, Rows: rows}, nil
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

// combine merges the (possibly multiple, e.g. one per network when no
// network filter is given) latest rollup rows for an entity into one cell,
// per the column's op. Trend is carried only when exactly one row matched.
func combine(col schema.ReportColumn, rows []schema.KpiRollup) Cell {
	cell := Cell{Column: col.Name, Format: col.Format}

	if len(rows) == 0 {
		return cell // zero cell: entity with no KPI data still lists
	}

	cell.hasData = true
	cell.Unit = rows[0].Unit
	cell.Symbol = rows[0].Symbol

	switch col.Op {
	case "MIN":
		cell.Value = rows[0].Value
		for _, r := range rows[1:] {
			if r.Value < cell.Value {
				cell.Value = r.Value
			}
		}
	case "MAX":
		cell.Value = rows[0].Value
		for _, r := range rows[1:] {
			if r.Value > cell.Value {
				cell.Value = r.Value
			}
		}
	case "AVG":
		sum := 0.0
		for _, r := range rows {
			sum += r.Value
		}
		cell.Value = sum / float64(len(rows)) // simple mean across scopes (documented approximation)
	default: // SUM, COUNT, LAST: additive across scopes
		for _, r := range rows {
			cell.Value += r.Value
		}
	}

	for _, r := range rows {
		if r.IsPartial {
			cell.IsPartial = true
		}

		if r.ComputedAt.After(cell.ComputedAt) {
			cell.ComputedAt = r.ComputedAt
		}
	}

	if len(rows) == 1 {
		cell.Trend = rows[0].Trend
		cell.PrevValue = rows[0].PrevValue
		cell.ChangeAbs = rows[0].ChangeAbs
		cell.ChangePct = rows[0].ChangePct
	}

	return cell
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
