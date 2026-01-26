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
	"gorm.io/gorm/clause"
)

type NodeRepo interface {
	GetNode(nodeId string) (*Node, error)
	CreateNode(node *Node) error
	DeleteNode(nodeId string) error
	UpdateNode(node *Node) error
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) *nodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (n *nodeRepo) GetNode(nodeId string) (*Node, error) {
	var node Node
	result := n.Db.GetGormDb().Where("node_id = ?", nodeId).First(&node)
	if result.Error != nil {
		return nil, result.Error
	}
	return &node, nil
}

func (n *nodeRepo) CreateNode(node *Node) error {
	return n.Db.GetGormDb().Create(node).Error
}

func (n *nodeRepo) DeleteNode(nodeId string) error {
	return n.Db.GetGormDb().Where("node_id = ?", nodeId).Delete(&Node{}).Error
}

func (n *nodeRepo) UpdateNode(node *Node) error {
	tx := n.Db.GetGormDb().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	result := tx.Clauses(clause.Returning{}).Where("id = ?", node.Id.String()).Updates(node)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	return tx.Commit().Error
}