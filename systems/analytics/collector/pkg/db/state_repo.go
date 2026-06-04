/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"time"

	"github.com/ukama/ukama/systems/common/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StateRepo interface {
	UpsertRefreshState(state *RefreshState) error
	GetRefreshStates() ([]RefreshState, error)
	MarkRollupDirty(rollup string) error
	SetRollupWatermark(rollup string, watermark time.Time) error
	GetRollupStates() ([]RollupState, error)
}

type stateRepo struct {
	Db gormHandle
}

func NewStateRepo(db sql.Db) StateRepo {
	return &stateRepo{
		Db: db,
	}
}

func NewStateRepoWithGorm(db *gorm.DB) StateRepo {
	return &stateRepo{
		Db: gormOnly{db: db},
	}
}

func (r *stateRepo) UpsertRefreshState(state *RefreshState) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "source"}},
		UpdateAll: true,
	}).Create(state)

	return result.Error
}

func (r *stateRepo) GetRefreshStates() ([]RefreshState, error) {
	var states []RefreshState

	result := r.Db.GetGormDb().Order("source asc").Find(&states)
	if result.Error != nil {
		return nil, result.Error
	}

	return states, nil
}

func (r *stateRepo) MarkRollupDirty(rollup string) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "rollup"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"dirty": true}),
	}).Create(&RollupState{
		Rollup: rollup,
		Dirty:  true,
	})

	return result.Error
}

func (r *stateRepo) SetRollupWatermark(rollup string, watermark time.Time) error {
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "rollup"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"watermark": watermark,
			"dirty":     false,
		}),
	}).Create(&RollupState{
		Rollup:    rollup,
		Watermark: watermark,
		Dirty:     false,
	})

	return result.Error
}

func (r *stateRepo) GetRollupStates() ([]RollupState, error) {
	var states []RollupState

	result := r.Db.GetGormDb().Order("rollup asc").Find(&states)
	if result.Error != nil {
		return nil, result.Error
	}

	return states, nil
}
