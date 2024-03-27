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

type AccountingRepo interface {
	Get(id uuid.UUID) (*Accounting, error)
	GetByCompany(company string) ([]*Accounting, error)
	GetByUser(userId string) ([]*Accounting, error)

	Add(accounts []*Accounting) error
	Delete(ids []string) error
}

type accountingRepo struct {
	Db sql.Db
}

func NewAccountingRepo(db sql.Db) AccountingRepo {
	return &accountingRepo{
		Db: db,
	}
}

func (c *accountingRepo) Get(id uuid.UUID) (*Accounting, error) {
	var account Accounting
	err := c.Db.GetGormDb().First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (c *accountingRepo) GetByCompany(company string) ([]*Accounting, error) {
	var accounts []*Accounting
	err := c.Db.GetGormDb().Where("company = ?", company).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (c *accountingRepo) GetByUser(userId string) ([]*Accounting, error) {
	var accounts []*Accounting
	err := c.Db.GetGormDb().Where("user_id = ?", userId).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (c *accountingRepo) Add(accounts []*Accounting) error {
	db := c.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	})

	result := db.Create(&accounts)

	return result.Error
}

func (c *accountingRepo) Delete(ids []string) error {
	db := c.Db.GetGormDb().Where("inventory IN ?", ids).Delete(&Accounting{})

	return db.Error
}
