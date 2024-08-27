/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
)

type ReportRepo interface {
	Add(report *Report, nestedFunc func(*Report, *gorm.DB) error) error
	Get(id uuid.UUID) (*Report, error)
	List(ownerId string, ownerType OwnerType, networkId string, reportType ReportType,
		isPaid bool, count uint32, sort bool) ([]Report, error)

	// Update(orgId uint, network *Network) error
	Delete(reportId uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error
}

type reportRepo struct {
	Db sql.Db
}

func NewReportRepo(db sql.Db) ReportRepo {
	return &reportRepo{
		Db: db,
	}
}

func (r *reportRepo) Add(report *Report, nestedFunc func(report *Report, tx *gorm.DB) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			if nestErr := nestedFunc(report, tx); nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(report)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (i *reportRepo) Get(id uuid.UUID) (*Report, error) {
	var rep Report

	result := i.Db.GetGormDb().First(&rep, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &rep, nil
}

func (r *reportRepo) List(ownerId string, ownerType OwnerType, networkId string,
	reportType ReportType, isPaid bool, count uint32, sort bool) ([]Report, error) {
	reports := []Report{}

	tx := r.Db.GetGormDb().Preload(clause.Associations)

	if ownerId != "" {
		tx = tx.Where("owner_id = ?", ownerId)
	}

	if ownerType != OwnerTypeUnknown {
		tx = tx.Where("owner_type = ?", ownerType)
	}

	if networkId != "" {
		tx = tx.Where("network_id = ?", networkId)
	}

	if reportType != ReportTypeUnknown {
		tx = tx.Where("type = ?", reportType)
	}

	if isPaid {
		tx = tx.Where("is_paid = ?", isPaid)
	}

	if sort {
		tx = tx.Order("time DESC")
	}

	if count > 0 {
		tx = tx.Limit(int(count))
	}

	result := tx.Find(&reports)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return reports, nil
}

func (r *reportRepo) Delete(ownerId uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Report{}, ownerId)
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			if nestErr := nestedFunc(ownerId, tx); nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}
