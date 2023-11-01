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
)

type SoftwareRepo interface {
	CreateSoftwareUpdate(Software *Software, nestedFunc func(string, string) error) error
	GetLatestSoftwareUpdate() (*Software, error)
}

type softwareRepo struct {
	Db sql.Db
}

func NewSoftwareRepo(db sql.Db) SoftwareRepo {
	return &softwareRepo{
		Db: db,
	}
}
func (r *softwareRepo) CreateSoftwareUpdate(Software *Software, nestedFunc func(string, string) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc("", "")
			if nestErr != nil {
				return nestErr
			}
		}
		if err := tx.Create(Software).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r *softwareRepo) GetLatestSoftwareUpdate() (*Software, error) {
	var Software Software
	err := r.Db.GetGormDb().Order("release_date desc").First(&Software).Error
	return &Software, err
}
