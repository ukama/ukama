/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type AppRepo interface {
	Create(app App) error
	GetAll() ([]App, error)
	Get(name string) (App, error)
}

type appRepo struct {
	db *gorm.DB
}

func NewAppRepo(db sql.Db) AppRepo {
	return &appRepo{db: db.GetGormDb()}
}

func (r *appRepo) Create(app App) error {
	return r.db.Create(&app).Error
}

func (r *appRepo) GetAll() ([]App, error) {
	var apps []App
	err := r.db.Find(&apps).Error
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	return apps, nil
}

func (r *appRepo) Get(name string) (App, error) {
	var app App
	err := r.db.Where("name = ?", name).First(&app).Error
	if err != nil {
		return App{}, gorm.ErrRecordNotFound
	}
	return app, nil
}