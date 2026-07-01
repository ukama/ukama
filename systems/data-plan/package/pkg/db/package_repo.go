/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	uuid "github.com/ukama/ukama/systems/common/uuid"
)

type PackageRepo interface {
	Add(dataPackage *Package, packageRate *PackageRate) error
	Get(uuid uuid.UUID) (*Package, error)
	GetDetails(uuid.UUID) (*Package, error)
	GetByName(name string) (*Package, error)
	Delete(uuid uuid.UUID) error
	GetAll() ([]Package, error)
	Update(uuid uuid.UUID, pkg *Package) error
}

type packageRepo struct {
	Db sql.Db
}

func NewPackageRepo(db sql.Db) *packageRepo {
	return &packageRepo{
		Db: db,
	}
}

func (r *packageRepo) Add(dataPackage *Package, packageRate *PackageRate) error {
	tx := r.Db.GetGormDb().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	result := tx.Create(dataPackage)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Create(packageRate)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

func (p *packageRepo) Get(uuid uuid.UUID) (*Package, error) {
	var _package Package

	result := p.Db.GetGormDb().Preload("PackageRate").Where("uuid = ?", uuid).First(&_package)

	if result.Error != nil {
		return nil, result.Error
	}

	return &_package, nil
}

func (p *packageRepo) GetDetails(uuid uuid.UUID) (*Package, error) {
	var _package Package

	result := p.Db.GetGormDb().Preload(clause.Associations).Where("uuid = ?", uuid).First(&_package)

	if result.Error != nil {
		return nil, result.Error
	}

	return &_package, nil
}

// GetByName returns the package matching the given name, case-insensitively
// (so "Plan", "plan" and "PLAN" are treated as the same). Soft-deleted packages
// are excluded by GORM, mirroring the partial unique index on LOWER(name) WHERE
// deleted_at IS NULL. Returns gorm.ErrRecordNotFound when the name is free.
func (p *packageRepo) GetByName(name string) (*Package, error) {
	var _package Package

	result := p.Db.GetGormDb().Where("LOWER(name) = LOWER(?)", name).First(&_package)
	if result.Error != nil {
		return nil, result.Error
	}

	return &_package, nil
}

func (p *packageRepo) GetAll() ([]Package, error) {
	var packages []Package
	result := p.Db.GetGormDb().Preload("PackageRate").Find(&packages)

	if result.Error != nil {
		return nil, result.Error
	}
	return packages, nil
}

func (r *packageRepo) Delete(uuid uuid.UUID) error {
	return r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("package_id = ?", uuid).Delete(&PackageRate{}).Error; err != nil {
			return err
		}
		if err := tx.Where("package_id = ?", uuid).Delete(&PackageDetails{}).Error; err != nil {
			return err
		}
		if err := tx.Where("package_id = ?", uuid).Delete(&PackageMarkup{}).Error; err != nil {
			return err
		}
		result := tx.Where("uuid = ?", uuid).Delete(&Package{})
		return result.Error
	})
}

func (b *packageRepo) Update(uuid uuid.UUID, pkg *Package) error {
	tx := b.Db.GetGormDb().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	result := tx.Clauses(clause.Returning{}).Where("uuid = ?", uuid).Updates(pkg)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	// TODO: Update is not updating the associations
	// https://stackoverflow.com/questions/65683156/updates-doesnt-seem-to-update-the-associations

	return tx.Commit().Error
}
