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
)

type SoftwareRepo interface {
	Create(Software *Software) error
	GetAll(nodeId string, status string) ([]Software, error)
	Update(Software *Software) error
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

func (r *softwareRepo) GetAll(nodeId string, status string) ([]Software, error) {
	var Software []Software
	err := r.Db.GetGormDb().Where("node_id = ? AND status = ?", nodeId, status).Find(&Software).Error
	if err != nil {
		return nil, err
	}
	return Software, nil
}

func (r *softwareRepo) Update(Software *Software) error {
	return r.Db.GetGormDb().Save(Software).Error
}