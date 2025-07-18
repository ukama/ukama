/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
)

type OrgRepo interface {
	/* Orgs */
	Add(org *Org, nestedFunc func(*Org, *gorm.DB) error) error
	Get(id uuid.UUID) (*Org, error)
	GetByName(name string) (*Org, error)
	GetByOwner(uuid uuid.UUID) ([]Org, error)
	GetByMember(id uint) ([]Org, error)
	GetAll() ([]Org, error)
	GetOrgCount() (int64, int64, error)
	AddUser(org *Org, user *User) error
	RemoveUser(org *Org, user *User) error
	// Update(id uint) error
	// Deactivate(id uint) error
	// Delete(id uint) error

}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) OrgRepo {
	return &orgRepo{
		Db: db,
	}
}

func (r *orgRepo) Add(org *Org, nestedFunc func(*Org, *gorm.DB) error) (err error) {
	if !validation.IsValidDnsLabelName(org.Name) {
		return fmt.Errorf("invalid name must be less then 253 " +
			"characters and consist of lowercase characters with a hyphen")
	}

	err = r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(org, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		d := tx.Create(org)
		if d.Error != nil {
			return d.Error
		}

		return nil
	})

	return err
}

func (r *orgRepo) Get(id uuid.UUID) (*Org, error) {
	var org Org

	result := r.Db.GetGormDb().First(&org, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

func (r *orgRepo) GetByName(name string) (*Org, error) {
	var org Org

	result := r.Db.GetGormDb().Where(&Org{Name: name}).First(&org)
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

func (r *orgRepo) GetByOwner(uuid uuid.UUID) ([]Org, error) {
	var orgs []Org

	result := r.Db.GetGormDb().Where(&Org{Owner: uuid}).Find(&orgs)
	if result.Error != nil {
		return nil, result.Error
	}

	return orgs, nil
}

func (r *orgRepo) GetByMember(id uint) ([]Org, error) {
	var membOrgs []Org
	result := r.Db.GetGormDb().Preload("Users", "id IN (?)", id).Find(&membOrgs)
	if result.Error != nil {
		return nil, result.Error
	}

	return membOrgs, nil
}

func (r *orgRepo) GetAll() ([]Org, error) {
	var orgs []Org

	result := r.Db.GetGormDb().Where(&Org{}).Find(&orgs)
	if result.Error != nil {
		return nil, result.Error
	}

	return orgs, nil
}

func (r *orgRepo) AddUser(org *Org, user *User) error {
	err := r.Db.GetGormDb().Model(org).Association("Users").Append(user)
	return err
}

func (r *orgRepo) RemoveUser(org *Org, user *User) error {
	err := r.Db.GetGormDb().Model(org).Association("Users").Append(user)
	return err
}

func (r *orgRepo) GetOrgCount() (int64, int64, error) {
	var activeOrgCount int64
	var deactiveOrgCount int64

	result := r.Db.GetGormDb().Model(&Org{}).
		Where("deactivated = ?", false).Count(&activeOrgCount)

	if result.Error != nil {
		return 0, 0, result.Error
	}

	result = r.Db.GetGormDb().Model(&Org{}).
		Where("deactivated = ?", true).Count(&deactiveOrgCount)

	if result.Error != nil {
		return 0, 0, result.Error
	}

	return activeOrgCount, deactiveOrgCount, nil
}
