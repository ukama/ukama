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
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

const (
	DefaultPageSize = 20

	StatusActive  = "active"
	StatusExpired = "expired"

	CustomerEventCreate           = "create"
	CustomerEventActivationFailed = "activation_failed"
)

// CustomerRepo is a read-only repository over the analytics database.
// The collector service owns the schema; this repo never writes.
type CustomerRepo interface {
	Counts(networkId uuid.UUID, from, to time.Time) (total, active, newCount, expired, failed uint32, err error)
	List(networkId, siteId uuid.UUID, status string, page, pageSize uint32) ([]CustomerSnapshot, int64, error)
	Search(query string, networkId uuid.UUID, page, pageSize uint32) ([]CustomerSnapshot, int64, error)
	Get(customerId uuid.UUID) (*CustomerSnapshot, error)
	PackageIntervals(customerId uuid.UUID) ([]CustomerPackageInterval, error)
	UsageBetween(customerId uuid.UUID, from, to time.Time) (float64, error)
	SiteNames(siteIds []uuid.UUID) (map[uuid.UUID]string, error)
}

type customerRepo struct {
	Db sql.Db
}

func NewCustomerRepo(db sql.Db) CustomerRepo {
	return &customerRepo{
		Db: db,
	}
}

func (r customerRepo) Counts(networkId uuid.UUID, from, to time.Time) (total, active, newCount, expired, failed uint32, err error) {
	db := r.Db.GetGormDb()

	snaps := func() *gorm.DB {
		q := db.Model(&CustomerSnapshot{})
		if networkId != uuid.Nil {
			q = q.Where("network_id = ?", networkId)
		}

		return q
	}

	events := func(kind string) *gorm.DB {
		q := db.Model(&CustomerEvent{}).
			Where("kind = ?", kind).
			Where("occurred_at >= ? AND occurred_at < ?", from, to)

		if networkId != uuid.Nil {
			q = q.Where(
				"customer_id IN (?)",
				db.Model(&CustomerSnapshot{}).Select("customer_id").
					Where("network_id = ?", networkId),
			)
		}

		return q
	}

	var totalCount, activeCount, expiredCount, createCount, failedCount int64

	if err = snaps().Count(&totalCount).Error; err != nil {
		return 0, 0, 0, 0, 0, err
	}

	if err = snaps().Where("status = ?", StatusActive).
		Count(&activeCount).Error; err != nil {
		return 0, 0, 0, 0, 0, err
	}

	if err = snaps().Where("status = ?", StatusExpired).
		Count(&expiredCount).Error; err != nil {
		return 0, 0, 0, 0, 0, err
	}

	if err = events(CustomerEventCreate).Count(&createCount).Error; err != nil {
		return 0, 0, 0, 0, 0, err
	}

	if err = events(CustomerEventActivationFailed).Count(&failedCount).Error; err != nil {
		return 0, 0, 0, 0, 0, err
	}

	return uint32(totalCount), uint32(activeCount), uint32(createCount),
		uint32(expiredCount), uint32(failedCount), nil
}

// normalizePage guards pagination inputs (page is 1-based).
func normalizePage(page, pageSize uint32) (uint32, uint32) {
	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = DefaultPageSize
	}

	return page, pageSize
}

func (r customerRepo) List(networkId, siteId uuid.UUID, status string, page, pageSize uint32) ([]CustomerSnapshot, int64, error) {
	var customers []CustomerSnapshot
	var count int64

	page, pageSize = normalizePage(page, pageSize)

	q := func() *gorm.DB {
		q := r.Db.GetGormDb().Model(&CustomerSnapshot{})

		if networkId != uuid.Nil {
			q = q.Where("network_id = ?", networkId)
		}

		if siteId != uuid.Nil {
			q = q.Where("site_id = ?", siteId)
		}

		if status != "" {
			q = q.Where("status = ?", status)
		}

		return q
	}

	if err := q().Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := q().Order("name asc").
		Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).
		Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, count, nil
}

func (r customerRepo) Search(query string, networkId uuid.UUID, page, pageSize uint32) ([]CustomerSnapshot, int64, error) {
	var customers []CustomerSnapshot
	var count int64

	page, pageSize = normalizePage(page, pageSize)

	pattern := "%" + query + "%"

	q := func() *gorm.DB {
		q := r.Db.GetGormDb().Model(&CustomerSnapshot{}).
			Where("name ILIKE ? OR email ILIKE ? OR sim_iccid ILIKE ?",
				pattern, pattern, pattern)

		if networkId != uuid.Nil {
			q = q.Where("network_id = ?", networkId)
		}

		return q
	}

	if err := q().Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := q().Order("name asc").
		Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).
		Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, count, nil
}

func (r customerRepo) Get(customerId uuid.UUID) (*CustomerSnapshot, error) {
	var customer CustomerSnapshot

	result := r.Db.GetGormDb().
		Where("customer_id = ?", customerId).First(&customer)
	if result.Error != nil {
		return nil, result.Error
	}

	return &customer, nil
}

func (r customerRepo) PackageIntervals(customerId uuid.UUID) ([]CustomerPackageInterval, error) {
	var intervals []CustomerPackageInterval

	result := r.Db.GetGormDb().
		Where("customer_id = ?", customerId).
		Order("start_at desc").
		Find(&intervals)
	if result.Error != nil {
		return nil, result.Error
	}

	return intervals, nil
}

func (r customerRepo) UsageBetween(customerId uuid.UUID, from, to time.Time) (float64, error) {
	var usage *float64

	result := r.Db.GetGormDb().Model(&CustomerUsageRollupDaily{}).
		Select("SUM(data_used_mb)").
		Where("customer_id = ? AND day >= ? AND day < ?", customerId, from, to).
		Scan(&usage)
	if result.Error != nil {
		return 0, result.Error
	}

	if usage == nil {
		return 0, nil
	}

	return *usage, nil
}

func (r customerRepo) SiteNames(siteIds []uuid.UUID) (map[uuid.UUID]string, error) {
	names := make(map[uuid.UUID]string)

	if len(siteIds) == 0 {
		return names, nil
	}

	var sites []SiteSnapshot

	result := r.Db.GetGormDb().
		Where("site_id IN ?", siteIds).Find(&sites)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, s := range sites {
		names[s.SiteId] = s.Name
	}

	return names, nil
}
