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

type ComponentRepo interface {
	Get(id uuid.UUID) (*Component, error)
	GetByCompany(company string, ctype string) ([]*Component, error)

	Add(invitation *Component, nestedFunc func(*Component, *gorm.DB) error) error
}

type componentRepo struct {
	Db sql.Db
}

func NewComponentRepo(db sql.Db) ComponentRepo {
	return &componentRepo{
		Db: db,
	}
}

func (c *componentRepo) Get(id uuid.UUID) (*Component, error) {
	var component Component
	err := c.Db.GetGormDb().First(&component, id).Error
	if err != nil {
		return nil, err
	}
	return &component, nil
}

func (c *componentRepo) GetByCompany(company string, ctype string) ([]*Component, error) {
	var components []*Component
	err := c.Db.GetGormDb().Where("company = ?", company).Where("type", ctype).Find(&components).Error
	if err != nil {
		return nil, err
	}
	return components, nil
}

func (c componentRepo) Add(component *Component, nestedFunc func(*Component, *gorm.DB) error) error {
	err := c.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(component, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		if err := tx.Create(component).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
