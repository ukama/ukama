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
	List() ([]Node, error)
	Delete(nodeId string) error
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

func (r *nodeRepo) List() ([]Node, error) {
	var nodes []Node
	err := r.db.Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// Delete hard-deletes the node and its software rows. Soft-deleted rows would
// still be read by SoftwareRepo.List and rewritten by every promotion.
func (r *nodeRepo) Delete(nodeId string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("node_id = ?", nodeId).Delete(&Software{}).Error; err != nil {
			return err
		}
		return tx.Where("node_id = ?", nodeId).Delete(&Node{}).Error
	})
}
