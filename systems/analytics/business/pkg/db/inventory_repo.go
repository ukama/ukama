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

// InventoryRepo is a read-only repository over sim and inventory snapshots.
type InventoryRepo interface {
	SimCounts() (available, active uint32, err error)
	NodeCounts() (available, deployed uint32, err error)
}

type inventoryRepo struct {
	Db sql.Db
}

func NewInventoryRepo(db sql.Db) InventoryRepo {
	return &inventoryRepo{
		Db: db,
	}
}

func (i inventoryRepo) SimCounts() (uint32, uint32, error) {
	var availableCount, activeCount int64

	if err := i.Db.GetGormDb().Model(&SimSnapshot{}).
		Where("status = ?", "available").
		Count(&availableCount).Error; err != nil {
		return 0, 0, err
	}

	if err := i.Db.GetGormDb().Model(&SimSnapshot{}).
		Where("status = ?", "active").
		Count(&activeCount).Error; err != nil {
		return 0, 0, err
	}

	return uint32(availableCount), uint32(activeCount), nil
}

func (i inventoryRepo) NodeCounts() (uint32, uint32, error) {
	var availableCount, deployedCount int64

	if err := i.Db.GetGormDb().Model(&InventorySnapshot{}).
		Where("state = ?", "available").
		Count(&availableCount).Error; err != nil {
		return 0, 0, err
	}

	if err := i.Db.GetGormDb().Model(&InventorySnapshot{}).
		Where("state = ?", "deployed").
		Count(&deployedCount).Error; err != nil {
		return 0, 0, err
	}

	return uint32(availableCount), uint32(deployedCount), nil
}
