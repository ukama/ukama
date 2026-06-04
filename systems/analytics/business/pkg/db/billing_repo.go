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
)

// BillingRepo is a read-only repository over billing snapshots and rollups.
type BillingRepo interface {
	GetBillingSnapshot() (*BillingSnapshot, error)
	InvoiceRollups(from, to time.Time) ([]BusinessBillingRollupDaily, error)
}

type billingRepo struct {
	Db sql.Db
}

func NewBillingRepo(db sql.Db) BillingRepo {
	return &billingRepo{
		Db: db,
	}
}

func (b billingRepo) GetBillingSnapshot() (*BillingSnapshot, error) {
	var snap BillingSnapshot

	// Org-level billing snapshot is a single row with Id = 1.
	result := b.Db.GetGormDb().Where("id = ?", 1).First(&snap)
	if result.Error != nil {
		return nil, result.Error
	}

	return &snap, nil
}

func (b billingRepo) InvoiceRollups(from, to time.Time) ([]BusinessBillingRollupDaily, error) {
	var rollups []BusinessBillingRollupDaily

	result := b.Db.GetGormDb().Model(&BusinessBillingRollupDaily{}).
		Where("day >= ? AND day < ?", from, to).
		Order("day DESC").
		Find(&rollups)
	if result.Error != nil {
		return nil, result.Error
	}

	return rollups, nil
}
