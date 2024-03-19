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
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type AccountRepo interface {
	Get(id uuid.UUID) (*Account, error)
	GetByCompany(company string, ctype string) ([]*Account, error)

	Add(account *Account, nestedFunc func(*Account, *gorm.DB) error) error
}

type accountRepo struct {
	Db sql.Db
}

func NewAccountRepo(db sql.Db) AccountRepo {
	return &accountRepo{
		Db: db,
	}
}

func (c *accountRepo) Get(id uuid.UUID) (*Account, error) {
	var account Account
	err := c.Db.GetGormDb().First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (c *accountRepo) GetByCompany(company string, ctype string) ([]*Account, error) {
	var accounts []*Account
	err := c.Db.GetGormDb().Where("company = ?", company).Where("type", ctype).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (c *accountRepo) Add(account *Account, nestedFunc func(*Account, *gorm.DB) error) error {
	err := c.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(account, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		if err := tx.Create(account).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
