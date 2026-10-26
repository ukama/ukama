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
	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type SoftwareRepo interface {
	Create(Software *Software) error
	Get(id uuid.UUID) (Software, error)
	List(nodeId string, status ukama.SoftwareStatusType, appName string) ([]*Software, error)
	Update(Software *Software) error
	SetDesiredForApp(appName, desired string) (int64, error)
}

type softwareRepo struct {
	Db sql.Db
}

func NewSoftwareRepo(db sql.Db) SoftwareRepo {
	return &softwareRepo{
		Db: db,
	}
}

func (r *softwareRepo) Create(Software *Software) error {
	return r.Db.GetGormDb().Create(Software).Error
}

func (r *softwareRepo) Get(id uuid.UUID) (Software, error) {
	var software Software
	err := r.Db.GetGormDb().Where("id = ?", id).Preload("App").First(&software).Error
	if err != nil {
		return Software{}, err
	}
	return software, nil
}

func (r *softwareRepo) List(nodeId string, status ukama.SoftwareStatusType, appName string) ([]*Software, error) {
	var software []*Software

	tx := r.Db.GetGormDb().Model(&Software{}).Preload("App")
	if appName != "" {
		tx = tx.Where("app_name = ?", appName)
	}
	if nodeId != "" {
		tx = tx.Where("node_id = ?", nodeId)
	}
	if status != ukama.SoftwareStatusType(0) {
		tx = tx.Where("status = ?", status)
	}

	result := tx.Find(&software)
	if result.Error != nil {
		return nil, result.Error
	}

	return software, nil
}

func (r *softwareRepo) Update(Software *Software) error {
	return r.Db.GetGormDb().Save(Software).Error
}

// SetDesiredForApp sets the desired version on every idle row of the app in one
// statement. Rows in UpdateInProgress keep their target for the update watcher.
func (r *softwareRepo) SetDesiredForApp(appName, desired string) (int64, error) {
	res := r.Db.GetGormDb().Model(&Software{}).
		Where("app_name = ? AND status <> ?", appName, uint8(ukama.UpdateInProgress)).
		Updates(map[string]interface{}{
			"desired_version": desired,
			"status": gorm.Expr("CASE WHEN current_version <> ? THEN ? ELSE ? END",
				desired, uint8(ukama.UpdateAvailable), uint8(ukama.UpToDate)),
		})
	return res.RowsAffected, res.Error
}
