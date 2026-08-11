/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"time"

	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/sql"
)

// KpiWindowReader is aggregator's read-only view of the KPI zone.
type KpiWindowReader interface {
	// WindowsInRange returns kpi_windows rows for a KPI with window_id in
	// [fromID, toID).
	WindowsInRange(orgID, kpiKey string, fromID, toID int64) ([]schema.KpiWindow, error)
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

func (r *repo) WindowsInRange(orgID, kpiKey string, fromID, toID int64) ([]schema.KpiWindow, error) {
	rows := []schema.KpiWindow{}

	err := r.db.GetGormDb().
		Where("org_id = ? AND kpi_key = ? AND window_id >= ? AND window_id < ?",
			orgID, kpiKey, fromID, toID).
		Order("window_id").
		Find(&rows).Error
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
