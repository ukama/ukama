/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/sql"
)

// inFlightLease: an in_flight entry older than this is considered abandoned
// and becomes re-claimable.
const inFlightLease = 15 * time.Minute

// RawReader is analysis' read-only view of the ingest zone.
type RawReader interface {
	StateAsOf(orgID, datasetKey string, windowID int64) ([]schema.RawRecord, error)
	WindowRows(orgID, datasetKey string, windowID int64) ([]schema.RawRecord, error)
}

// KpiRepo writes the KPI zone.
type KpiRepo interface {
	Upsert(rows []schema.KpiWindow) error
}

// LedgerRepo gives analysis its view of the shared ledger: read dataset
// progress, own the "kpi" kind.
type LedgerRepo interface {
	DatasetStatus(orgID, datasetKey string, windowID int64) (string, error)
	ClaimKpi(orgID, kpiKey string, windowID int64) (bool, error)
	MarkKpi(orgID, kpiKey string, windowID int64, status, detail string) error
	KpiStatus(orgID, kpiKey string, windowID int64) (string, error)
}

type ErrorRepo interface {
	Insert(orgID, kpiKey string, windowID int64, errMsg string) error
}

type repo struct {
	db sql.Db
}

func NewRepo(db sql.Db) *repo {
	return &repo{db: db}
}

func (r *repo) StateAsOf(orgID, datasetKey string, windowID int64) ([]schema.RawRecord, error) {
	rows := []schema.RawRecord{}

	err := r.db.GetGormDb().Raw(`
		SELECT * FROM (
			SELECT DISTINCT ON (entity_key) *
			FROM raw_records
			WHERE org_id = ? AND dataset_key = ? AND window_id <= ?
			ORDER BY entity_key, window_id DESC, id DESC
		) latest
		WHERE latest.deleted = false`,
		orgID, datasetKey, windowID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *repo) WindowRows(orgID, datasetKey string, windowID int64) ([]schema.RawRecord, error) {
	rows := []schema.RawRecord{}

	err := r.db.GetGormDb().
		Where("org_id = ? AND dataset_key = ? AND window_id = ? AND deleted = false",
			orgID, datasetKey, windowID).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *repo) Upsert(rows []schema.KpiWindow) error {
	if len(rows) == 0 {
		return nil
	}

	return r.db.GetGormDb().Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "kpi_key"}, {Name: "org_id"}, {Name: "scope"}, {Name: "window_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"value", "sum", "count", "min", "max",
			"value_type", "unit", "symbol", "algo_version", "computed_at",
		}),
	}).Create(&rows).Error
}

func (r *repo) DatasetStatus(orgID, datasetKey string, windowID int64) (string, error) {
	return r.status(orgID, schema.LedgerKindDataset, datasetKey, windowID)
}

func (r *repo) KpiStatus(orgID, kpiKey string, windowID int64) (string, error) {
	return r.status(orgID, schema.LedgerKindKpi, kpiKey, windowID)
}

func (r *repo) status(orgID, kind, key string, windowID int64) (string, error) {
	entry := schema.WindowLedger{}

	err := r.db.GetGormDb().
		Where("org_id = ? AND kind = ? AND key = ? AND window_id = ?", orgID, kind, key, windowID).
		First(&entry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}

		return "", err
	}

	return entry.Status, nil
}

func (r *repo) ClaimKpi(orgID, kpiKey string, windowID int64) (bool, error) {
	claimed := false

	err := r.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		entry := schema.WindowLedger{
			OrgID:    orgID,
			Kind:     schema.LedgerKindKpi,
			Key:      kpiKey,
			WindowID: windowID,
			Status:   schema.StatusInFlight,
			Attempt:  1,
		}

		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&entry)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected > 0 {
			claimed = true

			return nil
		}

		existing := schema.WindowLedger{}

		err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("org_id = ? AND kind = ? AND key = ? AND window_id = ?",
				orgID, schema.LedgerKindKpi, kpiKey, windowID).
			First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}

			return err
		}

		reclaimable := existing.Status == schema.StatusFailed ||
			existing.Status == schema.StatusDirty ||
			(existing.Status == schema.StatusInFlight &&
				time.Since(existing.UpdatedAt) > inFlightLease)

		if reclaimable {
			err = tx.Model(&existing).Updates(map[string]interface{}{
				"status":  schema.StatusInFlight,
				"attempt": existing.Attempt + 1,
			}).Error
			if err != nil {
				return err
			}

			claimed = true
		}

		return nil
	})

	return claimed, err
}

func (r *repo) MarkKpi(orgID, kpiKey string, windowID int64, status, detail string) error {
	return r.db.GetGormDb().Model(&schema.WindowLedger{}).
		Where("org_id = ? AND kind = ? AND key = ? AND window_id = ?",
			orgID, schema.LedgerKindKpi, kpiKey, windowID).
		Updates(map[string]interface{}{"status": status, "detail": detail}).Error
}

func (r *repo) Insert(orgID, kpiKey string, windowID int64, errMsg string) error {
	return r.db.GetGormDb().Create(&schema.AnalysisError{
		OrgID:     orgID,
		KpiKey:    kpiKey,
		WindowID:  windowID,
		Error:     errMsg,
		CreatedAt: time.Now().UTC(),
	}).Error
}
