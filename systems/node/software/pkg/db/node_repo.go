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
 
 type NodeRepo interface {
	Create(node Node) error
}

type nodeRepo struct {
	 db *gorm.DB
 }
 
func NewNodeRepo(db sql.Db) NodeRepo {
	return &nodeRepo{db: db.GetGormDb()}
}

func (r *nodeRepo) Create(node Node) error {
return r.db.Create(&node).Error
}