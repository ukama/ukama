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
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

const (
	SimStatusAvailable = "available"
	SimStatusAssigned  = "assigned"
	SimStatusActive    = "active"
	SimStatusSuspended = "suspended"
	SimStatusFaulty    = "faulty"
)

// SimRepo is a read-only repository over the analytics database.
type SimRepo interface {
	List(networkId uuid.UUID, status string, page, pageSize uint32) ([]SimSnapshot, int64, error)
	PoolCounts() (total, available, active, assigned, suspended, faulty uint32, err error)
	Batches() ([]SimBatchSnapshot, error)
}

type simRepo struct {
	Db sql.Db
}

func NewSimRepo(db sql.Db) SimRepo {
	return &simRepo{
		Db: db,
	}
}

func (r simRepo) List(networkId uuid.UUID, status string, page, pageSize uint32) ([]SimSnapshot, int64, error) {
	var sims []SimSnapshot
	var count int64

	page, pageSize = normalizePage(page, pageSize)

	db := r.Db.GetGormDb()

	q := func() *gorm.DB {
		q := db.Model(&SimSnapshot{})

		if status != "" {
			q = q.Where("status = ?", status)
		}

		if networkId != uuid.Nil {
			// sim snapshots have no network id; scope through the customer snapshot.
			q = q.Where(
				"customer_id IN (?)",
				db.Model(&CustomerSnapshot{}).Select("customer_id").
					Where("network_id = ?", networkId),
			)
		}

		return q
	}

	if err := q().Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := q().Order("iccid asc").
		Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).
		Find(&sims).Error; err != nil {
		return nil, 0, err
	}

	return sims, count, nil
}

func (r simRepo) PoolCounts() (total, available, active, assigned, suspended, faulty uint32, err error) {
	type statusCount struct {
		Status string
		Count  int64
	}

	var counts []statusCount

	result := r.Db.GetGormDb().Model(&SimSnapshot{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&counts)
	if result.Error != nil {
		return 0, 0, 0, 0, 0, 0, result.Error
	}

	for _, c := range counts {
		total += uint32(c.Count)

		switch c.Status {
		case SimStatusAvailable:
			available = uint32(c.Count)
		case SimStatusActive:
			active = uint32(c.Count)
		case SimStatusAssigned:
			assigned = uint32(c.Count)
		case SimStatusSuspended:
			suspended = uint32(c.Count)
		case SimStatusFaulty:
			faulty = uint32(c.Count)
		}
	}

	return total, available, active, assigned, suspended, faulty, nil
}

func (r simRepo) Batches() ([]SimBatchSnapshot, error) {
	var batches []SimBatchSnapshot

	result := r.Db.GetGormDb().
		Order("uploaded_at desc").Find(&batches)
	if result.Error != nil {
		return nil, result.Error
	}

	return batches, nil
}
