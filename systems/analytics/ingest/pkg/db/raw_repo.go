/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"time"

	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/sql"
)

type RawRepo interface {
	// InsertWindowed appends windowed records; duplicates (same dedup key)
	// are no-ops.
	InsertWindowed(records []schema.RawRecord) error
	// UpsertSnapshot applies change-log semantics for a full_snapshot pull of
	// one dataset window: rows are stored only when an entity's content hash
	// changed (or it re-appeared), and tombstones are stored for entities
	// that disappeared. Returns number of change rows written.
	UpsertSnapshot(orgID, datasetKey string, windowID int64, current []schema.RawRecord) (int, error)
	// StateAsOf returns the latest non-deleted row per entity with
	// window_id <= windowID.
	StateAsOf(orgID, datasetKey string, windowID int64) ([]schema.RawRecord, error)
	// WindowRows returns the rows belonging to exactly windowID.
	WindowRows(orgID, datasetKey string, windowID int64) ([]schema.RawRecord, error)
}

type rawRepo struct {
	db sql.Db
}

func NewRawRepo(db sql.Db) RawRepo {
	return &rawRepo{db: db}
}

func (r *rawRepo) InsertWindowed(records []schema.RawRecord) error {
	if len(records) == 0 {
		return nil
	}

	return r.db.GetGormDb().
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&records).Error
}

func (r *rawRepo) UpsertSnapshot(orgID, datasetKey string, windowID int64, current []schema.RawRecord) (int, error) {
	previous, err := r.StateAsOf(orgID, datasetKey, windowID)
	if err != nil {
		return 0, err
	}

	prevHash := make(map[string]string, len(previous))
	for _, p := range previous {
		prevHash[p.EntityKey] = p.ContentHash
	}

	changes := make([]schema.RawRecord, 0)
	seen := make(map[string]bool, len(current))

	for _, c := range current {
		seen[c.EntityKey] = true

		if h, ok := prevHash[c.EntityKey]; !ok || h != c.ContentHash {
			changes = append(changes, c)
		}
	}

	// Tombstones for entities that disappeared from the source.
	for _, p := range previous {
		if !seen[p.EntityKey] {
			changes = append(changes, schema.RawRecord{
				OrgID:      orgID,
				DatasetKey: datasetKey,
				WindowID:   windowID,
				EntityKey:  p.EntityKey,
				DedupKey:   p.EntityKey + ":deleted",
				Deleted:    true,
				EventTime:  p.EventTime,
				Payload:    "{}",
				Fields:     p.Fields,
				IngestedAt: time.Now().UTC(),
			})
		}
	}

	if len(changes) == 0 {
		return 0, nil
	}

	err = r.db.GetGormDb().
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&changes).Error
	if err != nil {
		return 0, err
	}

	return len(changes), nil
}

func (r *rawRepo) StateAsOf(orgID, datasetKey string, windowID int64) ([]schema.RawRecord, error) {
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

func (r *rawRepo) WindowRows(orgID, datasetKey string, windowID int64) ([]schema.RawRecord, error) {
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
