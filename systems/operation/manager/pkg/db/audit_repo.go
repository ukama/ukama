/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"context"
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type AuditRepo interface {
	Add(ctx context.Context, audit *OperationAudit) error
	AddInTx(tx *gorm.DB, audit *OperationAudit) error
}

type auditRepo struct {
	db sql.Db
}

func NewAuditRepo(db sql.Db) AuditRepo {
	return &auditRepo{db: db}
}

func (r *auditRepo) Add(ctx context.Context, audit *OperationAudit) error {
	return r.AddInTx(r.db.GetGormDb().WithContext(ctx), audit)
}

func (r *auditRepo) AddInTx(tx *gorm.DB, audit *OperationAudit) error {
	if audit.Id == uuid.Nil {
		audit.Id = uuid.NewV4()
	}
	if audit.At.IsZero() {
		audit.At = time.Now().UTC()
	}
	return tx.Create(audit).Error
}
