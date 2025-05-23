/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
)

type PackageRepo interface {
	Add(pkg *Package, nestedFunc func(*Package, *gorm.DB) error) error
	Get(packageID uuid.UUID) (*Package, error)
	List(simId, dataPlanId, fromStartDate, toSartDate, fromEndDate, toEndDate string,
		isActive, asExpired bool, count uint32, sort bool) ([]Package, error)

	// Deprecated: Use db.PackageRepo.List with simId as filtering param instead.
	GetBySim(simID uuid.UUID) ([]Package, error)

	GetOverlap(*Package) ([]Package, error)
	Update(pkg *Package, nestedFunc func(*Package, *gorm.DB) error) error
	Delete(packageID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error
}

type packageRepo struct {
	Db sql.Db
}

func NewPackageRepo(db sql.Db) PackageRepo {
	return &packageRepo{
		Db: db,
	}
}

func (p *packageRepo) Add(pkg *Package, nestedFunc func(pkg *Package, tx *gorm.DB) error) error {
	err := p.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(pkg, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		log.Info("Adding package", pkg)
		result := tx.Create(pkg)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (p *packageRepo) Get(packageID uuid.UUID) (*Package, error) {
	pkg := &Package{}

	result := p.Db.GetGormDb().Where("id = ?", packageID).First(pkg)
	if result.Error != nil {
		return nil, result.Error
	}

	return pkg, nil
}

func (p *packageRepo) List(simId, dataPlanId, fromStartDate, toStartDate,
	fromEndDate, toEndDate string, isActive, asExpired bool, count uint32, sort bool) ([]Package, error) {
	packages := []Package{}

	tx := p.Db.GetGormDb().Preload(clause.Associations)

	if simId != "" {
		tx = tx.Where("sim_id = ?", simId)
	}

	if dataPlanId != "" {
		tx = tx.Where("package_id = ?", dataPlanId)
	}

	if fromStartDate != "" {
		tx = tx.Where("start_date >= ?", fromStartDate)
	}

	if toStartDate != "" {
		tx = tx.Where("start_date <= ?", toStartDate)
	}

	if fromEndDate != "" {
		tx = tx.Where("end_date >= ?", fromEndDate)
	}

	if toEndDate != "" {
		tx = tx.Where("end_date <= ?", toEndDate)
	}

	if isActive {
		tx = tx.Where("is_active = ?", true)
	}

	if asExpired {
		tx = tx.Where("as_expired = ?", true)
	}

	if sort {
		tx = tx.Order("start_date ASC")
	}

	if count > 0 {
		tx = tx.Limit(int(count))
	}

	result := tx.Find(&packages)
	if result.Error != nil {
		return nil, result.Error
	}

	return packages, nil
}

// Deprecated: Use db.PackageRepo.List with simId as filtering param instead.
func (p *packageRepo) GetBySim(simID uuid.UUID) ([]Package, error) {
	var packages []Package

	result := p.Db.GetGormDb().Where(&Package{SimId: simID}).Find(&packages)
	if result.Error != nil {
		return nil, result.Error
	}

	return packages, nil
}

func (p *packageRepo) GetOverlap(pkg *Package) ([]Package, error) {
	var packages []Package

	result := p.Db.GetGormDb().Where(&Package{SimId: pkg.SimId}).Find(&packages,
		"end_date >= ? AND start_date <= ?", pkg.StartDate, pkg.EndDate)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}

	return packages, nil
}

func (p *packageRepo) Update(pkg *Package, nestedFunc func(*Package, *gorm.DB) error) error {
	err := p.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(pkg, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Clauses(clause.Returning{}).Updates(pkg)

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (p *packageRepo) Delete(packageID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error {
	err := p.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id=?", packageID).Delete(&Package{})
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(packageID, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}
