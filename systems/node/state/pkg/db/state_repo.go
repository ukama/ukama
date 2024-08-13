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
	"gorm.io/gorm"
)

type StateRepo interface {
	Create(state *State, nestedFunc func(*State, *gorm.DB) error) error
	GetByNodeId(nodeId ukama.NodeID) (*State, error)
	Update(state *State) error
	Delete(nodeId ukama.NodeID) error
	ListAll() ([]State, error)
	UpdateConnectivity(nodeId ukama.NodeID, connectivity Connectivity) error
	UpdateCurrentState(nodeId ukama.NodeID, currentState NodeStateEnum) error
	GetStateHistory(nodeId ukama.NodeID) ([]StateHistory, error)
}

type stateRepo struct {
	Db sql.Db
}

func NewStateRepo(db sql.Db) StateRepo {
	return &stateRepo{
		Db: db,
	}
}

func (r *stateRepo) Create(state *State, nestedFunc func(state *State, tx *gorm.DB) error) error {

	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(state, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(state)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (r *stateRepo) GetByNodeId(nodeId ukama.NodeID) (*State, error) {
	var state State
	err := r.Db.GetGormDb().First(&state, nodeId).Error
	if err != nil {
		return nil, err
	}
	return &state, nil

}

func (r *stateRepo) Update(state *State) error {

	result := r.Db.GetGormDb().Model(state).Updates(state)
    if result.Error != nil {
        return result.Error
    }

    return nil
}


func (r *stateRepo) Delete(nodeId ukama.NodeID) error {
	result := r.Db.GetGormDb().Where("id = ?", nodeId).Delete(&State{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *stateRepo) ListAll() ([]State, error) {
	var states []State
	err := r.Db.GetGormDb().Find(&states).Error
	return states, err
}

func (r *stateRepo) UpdateConnectivity(nodeId ukama.NodeID, connectivity Connectivity) error {
	return r.Db.GetGormDb().Model(&State{}).Where("node_id = ?", nodeId).Update("connectivity", connectivity).Error
}

func (r *stateRepo) UpdateCurrentState(nodeId ukama.NodeID, currentState NodeStateEnum) error {
	return r.Db.GetGormDb().Model(&State{}).Where("node_id = ?", nodeId).Update("current_state", currentState).Error
}

func (r *stateRepo) GetStateHistory(nodeId ukama.NodeID) ([]StateHistory, error) {
	var history []StateHistory
	err := r.Db.GetGormDb().Where("node_state_id = ?", nodeId).Order("timestamp desc").Find(&history).Error
	if err != nil {
		return nil, err
	}
	
	return history, nil
}

