/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"errors"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type NodeStateRepo interface {
	GetNodeStateById(id uuid.UUID) (*NodeState, error)
	GetNodeStatesByNodeId(nodeId string) ([]NodeState, error)
	GetNodeStateHistory(nodeId string) ([]NodeState, error)
	AddNodeState(newState *NodeState, previousState *NodeState) error
    GetLatestNodeState(nodeId string) (*NodeState, error)
}

type nodeStateRepo struct {
	Db sql.Db
}

func NewNodeStateRepo(db sql.Db) NodeStateRepo {
	return &nodeStateRepo{
		Db: db,
	}
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
	err := r.Db.GetGormDb().Where("node_id = ?", nodeId).Order("created_at DESC").First(&nodeState).Error
	if err != nil {
		return nil, err
	}
	return &nodeState, nil
}

func (r *nodeStateRepo) GetNodeStatesByNodeId(nodeId string) ([]NodeState, error) {
	var nodeStates []NodeState
	if err := r.Db.GetGormDb().Where("node_id = ?", nodeId).Order("created_at DESC").Find(&nodeStates).Error; err != nil {
		return nil, err
	}
	return nodeStates, nil
}

func (r *nodeStateRepo) GetNodeStateHistory(nodeId string) ([]NodeState, error) {
	var history []NodeState
	currentState, err := r.GetCurrentNodeState(nodeId)
	if err != nil {
		return nil, err
	}

	for currentState != nil {
		history = append(history, *currentState)
		if currentState.PreviousStateId == nil {
			break
		}
		currentState, err = r.GetNodeStateById(*currentState.PreviousStateId)
		if err != nil {
			return nil, err
		}
	}

	return history, nil
}
func (r *nodeStateRepo) GetLatestNodeState(nodeId string) (*NodeState, error) {
    var latestState NodeState
    result := r.Db.GetGormDb().Where("node_id = ?", nodeId).Order("created_at DESC").First(&latestState)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil 
        }
        return nil, result.Error
    }
    return &latestState, nil
}

func (r *nodeStateRepo) AddNodeState(newState *NodeState, previousState *NodeState) error {
    return r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
        if previousState != nil {
            newState.PreviousStateId = &previousState.Id
        }
        
        if err := tx.Create(newState).Error; err != nil {
            return err
        }
        
        return nil
    })
}
