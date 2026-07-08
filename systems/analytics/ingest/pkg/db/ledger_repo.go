/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
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
// (worker crashed between Claim and Mark) and becomes re-claimable.
const inFlightLease = 15 * time.Minute

type LedgerRepo interface {
	// Claim atomically claims (kind, key, org, window) for processing.
	// Returns false when another worker already claimed or completed it.
	// Failed and dirty entries are re-claimable. This is the persistent
	// "pulled list per run" — it survives restarts and replicas.
	Claim(orgID, kind, key string, windowID int64) (bool, error)
	// Mark sets the terminal status for a claimed entry.
	Mark(orgID, kind, key string, windowID int64, status, detail string) error
	// Status returns the entry status, or "" when absent.
	Status(orgID, kind, key string, windowID int64) (string, error)
	// LastWithStatus returns the newest window id for (kind, key) with the
	// given status, or -1 when none.
	LastWithStatus(orgID, kind, key, status string) (int64, error)
	// WindowsWithStatus lists window ids for (kind, key) with the given
	// status within [from, to].
	WindowsWithStatus(orgID, kind, key, status string, from, to int64) ([]int64, error)
}

type ledgerRepo struct {
	db sql.Db
}

func NewLedgerRepo(db sql.Db) LedgerRepo {
	return &ledgerRepo{db: db}
}

func (l *ledgerRepo) Claim(orgID, kind, key string, windowID int64) (bool, error) {
	claimed := false

	err := l.db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		entry := schema.WindowLedger{
			OrgID:    orgID,
			Kind:     kind,
			Key:      key,
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

		// Row exists: re-claim only failed/dirty entries, with row lock so
		// concurrent replicas cannot double-claim.
		existing := schema.WindowLedger{}

		err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("org_id = ? AND kind = ? AND key = ? AND window_id = ?", orgID, kind, key, windowID).
			First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil // locked by another worker
			}

			return err
		}

		reclaimable := existing.Status == schema.StatusFailed ||
			existing.Status == schema.StatusDirty ||
			(existing.Status == schema.StatusInFlight &&
				time.Since(existing.UpdatedAt) > inFlightLease)

		if reclaimable {
			err = tx.Model(&existing).
				Updates(map[string]interface{}{
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

func (l *ledgerRepo) Mark(orgID, kind, key string, windowID int64, status, detail string) error {
	return l.db.GetGormDb().Model(&schema.WindowLedger{}).
		Where("org_id = ? AND kind = ? AND key = ? AND window_id = ?", orgID, kind, key, windowID).
		Updates(map[string]interface{}{"status": status, "detail": detail}).Error
}

func (l *ledgerRepo) Status(orgID, kind, key string, windowID int64) (string, error) {
	entry := schema.WindowLedger{}

	err := l.db.GetGormDb().
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

func (l *ledgerRepo) LastWithStatus(orgID, kind, key, status string) (int64, error) {
	var res *int64

	err := l.db.GetGormDb().Model(&schema.WindowLedger{}).
		Select("max(window_id)").
		Where("org_id = ? AND kind = ? AND key = ? AND status = ?", orgID, kind, key, status).
		Scan(&res).Error
	if err != nil {
		return -1, err
	}

	if res == nil {
		return -1, nil
	}

	return *res, nil
}

func (l *ledgerRepo) WindowsWithStatus(orgID, kind, key, status string, from, to int64) ([]int64, error) {
	ids := []int64{}

	err := l.db.GetGormDb().Model(&schema.WindowLedger{}).
		Where("org_id = ? AND kind = ? AND key = ? AND status = ? AND window_id BETWEEN ? AND ?",
			orgID, kind, key, status, from, to).
		Order("window_id").
		Pluck("window_id", &ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}
