/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"errors"
	"time"

	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
)

type IntentRepo interface {
	Add(intent *Intent) error
	Get(operationId uuid.UUID) (*Intent, error)
	List(f IntentFilter) ([]Intent, error)
	Update(operationId uuid.UUID, ch IntentUpdate) (*Intent, error)
}

type IntentFilter struct {
	ResourceKey    string
	Status         *IntentStatus
	DeadlineBefore time.Time
	Limit          int
}

type IntentUpdate struct {
	Status     *IntentStatus
	Armed      *bool
	Dispatched *bool
}

type intentRepo struct {
	db sql.Db
}

func NewIntentRepo(db sql.Db) IntentRepo {
	return &intentRepo{db: db}
}

func (r *intentRepo) Add(intent *Intent) error {
	return r.db.GetGormDb().Create(intent).Error
}

func (r *intentRepo) Get(operationId uuid.UUID) (*Intent, error) {
	var intent Intent
	if err := r.db.GetGormDb().Where("operation_id = ?",
		operationId).First(&intent).Error; err != nil {
		return nil, err
	}
	return &intent, nil
}

func (r *intentRepo) List(f IntentFilter) ([]Intent, error) {
	q := r.db.GetGormDb().Model(&Intent{})
	if f.ResourceKey != "" {
		q = q.Where("resource_key = ?", f.ResourceKey)
	}
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if !f.DeadlineBefore.IsZero() {
		q = q.Where("deadline < ?", f.DeadlineBefore)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	var intents []Intent
	err := q.Order("deadline ASC").Find(&intents).Error
	return intents, err
}

func (r *intentRepo) Update(operationId uuid.UUID, ch IntentUpdate) (*Intent, error) {
	fields := map[string]any{}
	if ch.Status != nil {
		if !ch.Status.IsTerminal() {
			return nil, errors.New("intent status can only change to a terminal one")
		}
		fields["status"] = *ch.Status
	}
	if ch.Armed != nil {
		fields["armed"] = *ch.Armed
	}
	if ch.Dispatched != nil {
		fields["dispatched"] = *ch.Dispatched
	}
	if len(fields) == 0 {
		return nil, errors.New("intent update has no fields")
	}
	var updated []Intent
	res := r.db.GetGormDb().Model(&updated).Clauses(clause.Returning{}).
		Where("operation_id = ? AND status = ?", operationId, IntentWatching).
		Updates(fields)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 || len(updated) == 0 {
		return r.Get(operationId)
	}
	return &updated[0], nil
}
