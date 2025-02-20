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
	"gorm.io/gorm/clause"
)

type ComponentRepo interface {
	Get(id uuid.UUID) (*Component, error)
	GetByUser(userId string, category int32) ([]*Component, error)
	Add(components []*Component) error
	Delete() error
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

func (c *componentRepo) GetByUser(userId string, category int32) ([]*Component, error) {
	var components []*Component

	tx := c.Db.GetGormDb().Preload(clause.Associations)
	tx = tx.Where("user_id = ?", userId)

	if category != 0 {
		tx = tx.Where("category = ?", category)
	}

	result := tx.Find(&components)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return components, nil
}

func (c *componentRepo) Add(components []*Component) error {
	db := c.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	})

	result := db.Create(&components)

	return result.Error
}

func (c *componentRepo) Delete() error {
	db := c.Db.GetGormDb().Exec("DELETE FROM components")
	return db.Error
}
