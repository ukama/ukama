/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/sql"

	"gorm.io/gorm/clause"
)

type EventRepo interface {
	// LogEvent records an incoming event for idempotency and audit. It
	// returns false (and no error) when an event with the same MsgId was
	// already recorded, in which case the caller must skip processing.
	LogEvent(log *EventLog) (bool, error)
	RecordError(e *EventError) error
	GetRecent(limit int) ([]EventLog, error)
}

type eventRepo struct {
	Db sql.Db
}

func NewEventRepo(db sql.Db) EventRepo {
	return &eventRepo{
		Db: db,
	}
}

func (r *eventRepo) LogEvent(log *EventLog) (bool, error) {
	// Idempotency via ON CONFLICT (msg_id) DO NOTHING: RowsAffected == 0
	// means the event was a duplicate delivery.
	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "msg_id"}},
		DoNothing: true,
	}).Create(log)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (r *eventRepo) RecordError(e *EventError) error {
	result := r.Db.GetGormDb().Create(e)

	return result.Error
}

func (r *eventRepo) GetRecent(limit int) ([]EventLog, error) {
	var logs []EventLog

	result := r.Db.GetGormDb().Order("occurred_at desc").Limit(limit).Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}

	return logs, nil
}
