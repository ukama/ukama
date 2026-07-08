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

	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/sql"
)

type ErrorRepo interface {
	Insert(orgID, datasetKey string, windowID int64, iteration, errMsg string) error
}

type errorRepo struct {
	db sql.Db
}

func NewErrorRepo(db sql.Db) ErrorRepo {
	return &errorRepo{db: db}
}

func (e *errorRepo) Insert(orgID, datasetKey string, windowID int64, iteration, errMsg string) error {
	return e.db.GetGormDb().Create(&schema.IngestError{
		OrgID:      orgID,
		DatasetKey: datasetKey,
		WindowID:   windowID,
		Iteration:  iteration,
		Error:      errMsg,
		CreatedAt:  time.Now().UTC(),
	}).Error
}
