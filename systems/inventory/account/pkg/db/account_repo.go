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
	"gorm.io/gorm/clause"
)

type AccountRepo interface {
	Get(id uuid.UUID) (*Account, error)
	GetByCompany(company string) ([]*Account, error)

	Add(accounts []Account) error
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

func (c *accountRepo) GetByCompany(company string) ([]*Account, error) {
	var accounts []*Account
	err := c.Db.GetGormDb().Where("company = ?", company).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (c *accountRepo) Add(accounts []Account) error {
	db := c.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	})

	result := db.Create(&accounts)

	return result.Error
}
