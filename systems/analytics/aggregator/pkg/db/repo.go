/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"fmt"
	"time"

	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/sql"
)

type ValueMode int

const (
	ValueSum ValueMode = iota
	ValueAvg
	ValueLast
)

type WindowAgg struct {
	Scope        string
	Sum          float64
	Count        float64
	Min          float64
	Max          float64
	Last         float64
	LastWindowID int64
	Unit         string
	Symbol       string
	ComputedAt   time.Time
}

type SpanState struct {
	SpanStart  time.Time
	IsPartial  bool
	ComputedAt time.Time
}

// KpiWindowReader is aggregator's read-only view of the KPI zone.
type KpiWindowReader interface {
	AggregateWindows(orgID, kpiKey string, fromID, toID int64) ([]WindowAgg, error)
	// WindowBounds returns the min/max window_id a KPI has rows for
	// (ok=false when it has none) — drives the boot backfill.
	WindowBounds(orgID, kpiKey string) (minID, maxID int64, ok bool, err error)
}

// RawStateReader is aggregator's read-only view of the ingest zone's
// change-log state — used by performance reports for entity attributes.
type RawStateReader interface {
	// LatestState returns the latest non-deleted row per entity for a
	// dataset (state as of now).
	LatestState(orgID, datasetKey string) ([]schema.RawRecord, error)
}

// RollupRepo owns the rollup zone. Reads return the single components row
// per (scope, span_start) — schema.RollupRowOp; there is no per-op storage.
type RollupRepo interface {
	Upsert(rows []schema.KpiRollup) error
	// Latest returns the newest span row per scope for a KPI (optionally
	// scope-filtered).
	Latest(orgID, kpiKey, span string, scopeFilter map[string]string) ([]schema.KpiRollup, error)
	// Range returns rows ordered by span_start within [from, to).
	Range(orgID, kpiKey, span string, from, to time.Time, scopeFilter map[string]string) ([]schema.KpiRollup, error)
}

type RollupWriter interface {
	RollupFromWindows(orgID, kpiKey string, mode ValueMode, spanStart, spanEnd time.Time,
		fromID, toID int64, isPartial bool, now time.Time) error
	RollupFromDaily(orgID, kpiKey, span string, mode ValueMode, spanStart, spanEnd time.Time,
		isPartial bool, now time.Time) error
	SpanStates(orgID, kpiKey, span string) ([]SpanState, error)
}

// RetentionRepo prunes the KPI zone.
type RetentionRepo interface {
	PruneWindowsBefore(orgID string, beforeWindow int64, batch int) (int64, error)
	PruneKpiLedgerBefore(orgID string, beforeWindow int64, batch int) (int64, error)
}

// spanDaily mirrors rollup.SpanDaily; rollup imports this package.
const spanDaily = "daily"

type repo struct {
	db sql.Db
}

func NewRepo(db sql.Db) *repo {
	return &repo{db: db}
}

func (r *repo) LatestState(orgID, datasetKey string) ([]schema.RawRecord, error) {
	rows := []schema.RawRecord{}

	err := r.db.GetGormDb().Raw(`
		SELECT * FROM (
			SELECT DISTINCT ON (entity_key) *
			FROM raw_records
			WHERE org_id = ? AND dataset_key = ?
			ORDER BY entity_key, window_id DESC, id DESC
		) latest
		WHERE latest.deleted = false`,
		orgID, datasetKey).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// MAX over [order key, col] arrays picks the latest row's col without sorting each group.
const (
	latestWindowValue = `(MAX(ARRAY[window_id::float8, "value"]))[2]`
	latestWindowText  = `(MAX(ARRAY[lpad(window_id::text, 20, '0'), %s]))[2]`
	latestDailyLast   = `(MAX(ARRAY[extract(epoch FROM span_start)::float8, "last"]))[2]`
	latestDailyText   = `(MAX(ARRAY[lpad(extract(epoch FROM span_start)::bigint::text, 20, '0'), %s]))[2]`
)

const rollupInsert = `
	INSERT INTO kpi_rollups
		(kpi_key, org_id, scope, span, span_start, span_end, op, "value",
		 "sum", "count", "min", "max", "last", value_type, unit, symbol, is_partial, computed_at)
	SELECT a.kpi_key, a.org_id, a.scope, ?, ?, ?, ?, %s,
		a."sum", a."count", a."min", a."max", a."last", a.value_type, a.unit, a.symbol, ?, ?
	FROM (%s) a
	ON CONFLICT (kpi_key, org_id, scope, span, span_start, op) DO UPDATE SET
		span_end = EXCLUDED.span_end, "value" = EXCLUDED."value", "sum" = EXCLUDED."sum",
		"count" = EXCLUDED."count", "min" = EXCLUDED."min", "max" = EXCLUDED."max",
		"last" = EXCLUDED."last", value_type = EXCLUDED.value_type, unit = EXCLUDED.unit,
		symbol = EXCLUDED.symbol, is_partial = EXCLUDED.is_partial, computed_at = EXCLUDED.computed_at`

var windowsFold = `
	SELECT kpi_key, org_id, scope,
		SUM("sum") AS "sum", SUM("count") AS "count", MIN("min") AS "min", MAX("max") AS "max",
		` + latestWindowValue + ` AS "last",
		` + fmt.Sprintf(latestWindowText, "value_type") + ` AS value_type,
		` + fmt.Sprintf(latestWindowText, "unit") + ` AS unit,
		` + fmt.Sprintf(latestWindowText, "symbol") + ` AS symbol
	FROM kpi_windows
	WHERE org_id = ? AND kpi_key = ? AND window_id >= ? AND window_id < ?
	GROUP BY kpi_key, org_id, scope`

var dailyFold = `
	SELECT kpi_key, org_id, scope,
		SUM("sum") AS "sum", SUM("count") AS "count", MIN("min") AS "min", MAX("max") AS "max",
		` + latestDailyLast + ` AS "last",
		` + fmt.Sprintf(latestDailyText, "value_type") + ` AS value_type,
		` + fmt.Sprintf(latestDailyText, "unit") + ` AS unit,
		` + fmt.Sprintf(latestDailyText, "symbol") + ` AS symbol
	FROM kpi_rollups
	WHERE org_id = ? AND kpi_key = ? AND span = ? AND op = ? AND span_start >= ? AND span_start < ?
	GROUP BY kpi_key, org_id, scope`

func valueExpr(mode ValueMode) string {
	switch mode {
	case ValueAvg:
		return `CASE WHEN a."count" = 0 THEN 0 ELSE a."sum" / a."count" END`
	case ValueLast:
		return `a."last"`
	default:
		return `a."sum"`
	}
}

func (r *repo) AggregateWindows(orgID, kpiKey string, fromID, toID int64) ([]WindowAgg, error) {
	rows := []WindowAgg{}

	err := r.db.GetGormDb().Raw(`
		SELECT scope,
			SUM("sum") AS "sum", SUM("count") AS "count", MIN("min") AS "min", MAX("max") AS "max",
			`+latestWindowValue+` AS "last", MAX(window_id) AS last_window_id,
			`+fmt.Sprintf(latestWindowText, "unit")+` AS unit,
			`+fmt.Sprintf(latestWindowText, "symbol")+` AS symbol,
			MAX(computed_at) AS computed_at
		FROM kpi_windows
		WHERE org_id = ? AND kpi_key = ? AND window_id >= ? AND window_id < ?
		GROUP BY scope`,
		orgID, kpiKey, fromID, toID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *repo) RollupFromWindows(orgID, kpiKey string, mode ValueMode, spanStart, spanEnd time.Time,
	fromID, toID int64, isPartial bool, now time.Time) error {
	return r.db.GetGormDb().Exec(fmt.Sprintf(rollupInsert, valueExpr(mode), windowsFold),
		spanDaily, spanStart.UTC(), spanEnd.UTC(), schema.RollupRowOp, isPartial, now.UTC(),
		orgID, kpiKey, fromID, toID).Error
}

func (r *repo) RollupFromDaily(orgID, kpiKey, span string, mode ValueMode, spanStart, spanEnd time.Time,
	isPartial bool, now time.Time) error {
	return r.db.GetGormDb().Exec(fmt.Sprintf(rollupInsert, valueExpr(mode), dailyFold),
		span, spanStart.UTC(), spanEnd.UTC(), schema.RollupRowOp, isPartial, now.UTC(),
		orgID, kpiKey, spanDaily, schema.RollupRowOp, spanStart.UTC(), spanEnd.UTC()).Error
}

func (r *repo) SpanStates(orgID, kpiKey, span string) ([]SpanState, error) {
	rows := []SpanState{}

	err := r.db.GetGormDb().Raw(`
		SELECT span_start, bool_or(is_partial) AS is_partial, MIN(computed_at) AS computed_at
		FROM kpi_rollups
		WHERE org_id = ? AND kpi_key = ? AND span = ? AND op = ?
		GROUP BY span_start`,
		orgID, kpiKey, span, schema.RollupRowOp).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *repo) Upsert(rows []schema.KpiRollup) error {
	if len(rows) == 0 {
		return nil
	}

	return r.db.GetGormDb().Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "kpi_key"}, {Name: "org_id"}, {Name: "scope"},
			{Name: "span"}, {Name: "span_start"}, {Name: "op"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"span_end", "value", "sum", "count", "min", "max", "last",
			"value_type", "unit", "symbol", "is_partial", "computed_at",
		}),
	}).Create(&rows).Error
}

func (r *repo) Latest(orgID, kpiKey, span string, scopeFilter map[string]string) ([]schema.KpiRollup, error) {
	rows := []schema.KpiRollup{}

	q := r.db.GetGormDb().Raw(`
		SELECT DISTINCT ON (scope) *
		FROM kpi_rollups
		WHERE org_id = ? AND kpi_key = ? AND span = ? AND op = ?
		ORDER BY scope, span_start DESC`,
		orgID, kpiKey, span, schema.RollupRowOp)

	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	return filterScope(rows, scopeFilter), nil
}

func (r *repo) Range(orgID, kpiKey, span string, from, to time.Time, scopeFilter map[string]string) ([]schema.KpiRollup, error) {
	rows := []schema.KpiRollup{}

	err := r.db.GetGormDb().
		Where("org_id = ? AND kpi_key = ? AND span = ? AND op = ? AND span_start >= ? AND span_start < ?",
			orgID, kpiKey, span, schema.RollupRowOp, from, to).
		Order("span_start").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return filterScope(rows, scopeFilter), nil
}

func (r *repo) WindowBounds(orgID, kpiKey string) (int64, int64, bool, error) {
	var bounds struct {
		MinID *int64
		MaxID *int64
	}

	err := r.db.GetGormDb().Model(&schema.KpiWindow{}).
		Select("MIN(window_id) AS min_id, MAX(window_id) AS max_id").
		Where("org_id = ? AND kpi_key = ?", orgID, kpiKey).
		Scan(&bounds).Error
	if err != nil {
		return 0, 0, false, err
	}

	if bounds.MinID == nil || bounds.MaxID == nil {
		return 0, 0, false, nil
	}

	return *bounds.MinID, *bounds.MaxID, true, nil
}

func (r *repo) PruneWindowsBefore(orgID string, beforeWindow int64, batch int) (int64, error) {
	res := r.db.GetGormDb().Exec(`
		DELETE FROM kpi_windows WHERE id IN (
			SELECT id FROM kpi_windows
			WHERE org_id = ? AND window_id < ?
			LIMIT ?)`,
		orgID, beforeWindow, batch)

	return res.RowsAffected, res.Error
}

func (r *repo) PruneKpiLedgerBefore(orgID string, beforeWindow int64, batch int) (int64, error) {
	res := r.db.GetGormDb().Exec(`
		DELETE FROM window_ledgers WHERE id IN (
			SELECT id FROM window_ledgers
			WHERE org_id = ? AND kind = ? AND window_id < ?
			LIMIT ?)`,
		orgID, schema.LedgerKindKpi, beforeWindow, batch)

	return res.RowsAffected, res.Error
}

// filterScope keeps rows whose canonical scope JSON contains all requested
// key/value pairs. Done in Go to stay portable (scope is a varchar column).
//
// REVENUE / PAID_CUSTOMERS are network-scoped (payments attributed to a
// network via the paying SIM; unresolvable SIMs go to an org bucket). A
// per-network read returns only rows whose network_id matches — the org
// bucket (empty scope) is org-only and never bleeds into a network number.
func filterScope(rows []schema.KpiRollup, filter map[string]string) []schema.KpiRollup {
	if len(filter) == 0 {
		return rows
	}

	out := make([]schema.KpiRollup, 0, len(rows))

	for _, row := range rows {
		scope := schema.ParseScope(row.Scope)

		match := true
		for k, v := range filter {
			if scope[k] != v {
				match = false

				break
			}
		}

		if match {
			out = append(out, row)
		}
	}

	return out
}
