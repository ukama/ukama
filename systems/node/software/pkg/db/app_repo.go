/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"gorm.io/gorm"
)

type AppRepo interface {
	Create(app App) error
	GetAll() ([]App, error)
}

type appRepo struct {
	db *gorm.DB
}

func NewAppRepo(db *gorm.DB) AppRepo {
	return &appRepo{db: db}
}

func (r *appRepo) Create(app App) error {
	return r.db.Create(&app).Error
}

func (r *appRepo) GetAll() ([]App, error) {
	var apps []App
	err := r.db.Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}