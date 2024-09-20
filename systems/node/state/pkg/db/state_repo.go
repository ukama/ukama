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

type NodeStateRepo interface {
	AddNodeState(nodeState *NodeState, nestedFunc func(*NodeState, *gorm.DB) error) error
	GetNodeStateById(id uuid.UUID) (*NodeState, error)
	GetCurrentNodeState(nodeId string) (*NodeState, error)
	GetNodeStatesByNodeId(nodeId string) ([]NodeState, error)
}


type nodeStateRepo struct {
	Db sql.Db
}

func NewNodeStateRepo(db sql.Db) NodeStateRepo {
	return &nodeStateRepo{
		Db: db,
	}
}




func (r *nodeStateRepo) AddNodeState(nodeState *NodeState, nestedFunc func(nodeState *NodeState, tx *gorm.DB) error) error {

	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(nodeState, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(nodeState)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (r *nodeStateRepo) GetNodeStateById(id uuid.UUID) (*NodeState, error) {
	var nodeState NodeState

	err := r.Db.GetGormDb().First(&nodeState, id).Error
	if err != nil {
		return nil, err
	}
	return &nodeState, nil
}

func (r *nodeStateRepo) GetCurrentNodeState(nodeId string) (*NodeState, error) {
	var nodeState NodeState

	err := r.Db.GetGormDb().First(&nodeState, nodeId).Error
	if err != nil {
		return nil, err
	}
	return &nodeState, nil
}


func (r *nodeStateRepo) GetNodeStatesByNodeId(nodeId string) ([]NodeState, error) {
	var nodeStates []NodeState
	if err := r.Db.GetGormDb().Where("node_id = ?", nodeId).Find(&nodeStates).Error; err != nil {
		return nil, err
	}
	return nodeStates, nil
}
