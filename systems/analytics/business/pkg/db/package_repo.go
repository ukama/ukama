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

// PackageRepo is a read-only repository over package snapshots and rollups.
type PackageRepo interface {
	ListPackages(page, pageSize int) ([]PackageSnapshot, int64, error)
	PackageRollups(from, to time.Time) ([]BusinessPackageRollupDaily, error)
}

type packageRepo struct {
	Db sql.Db
}

func NewPackageRepo(db sql.Db) PackageRepo {
	return &packageRepo{
		Db: db,
	}
}

func (p packageRepo) ListPackages(page, pageSize int) ([]PackageSnapshot, int64, error) {
	var pkgs []PackageSnapshot
	var count int64

	if err := p.Db.GetGormDb().Model(&PackageSnapshot{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	q := p.Db.GetGormDb().Model(&PackageSnapshot{}).Order("name ASC")

	if pageSize > 0 {
		if page < 1 {
			page = 1
		}
		q = q.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	result := q.Find(&pkgs)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return pkgs, count, nil
}

func (p packageRepo) PackageRollups(from, to time.Time) ([]BusinessPackageRollupDaily, error) {
	var rollups []BusinessPackageRollupDaily

	result := p.Db.GetGormDb().Model(&BusinessPackageRollupDaily{}).
		Where("day >= ? AND day < ?", from, to).
		Order("day ASC").
		Find(&rollups)
	if result.Error != nil {
		return nil, result.Error
	}

	return rollups, nil
}
