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
)

// ActivityRepo is a read-only repository over the collector event log.
type ActivityRepo interface {
	Recent(limit int) ([]EventLog, error)
}

type activityRepo struct {
	Db sql.Db
}

func NewActivityRepo(db sql.Db) ActivityRepo {
	return &activityRepo{
		Db: db,
	}
}

func (a activityRepo) Recent(limit int) ([]EventLog, error) {
	var logs []EventLog

	if limit <= 0 {
		limit = 10
	}

	result := a.Db.GetGormDb().Model(&EventLog{}).
		Order("occurred_at DESC").
		Limit(limit).
		Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}

	return logs, nil
}
