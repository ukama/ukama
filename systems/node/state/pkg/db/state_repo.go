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
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type StateRepo interface {
	GetStateById(id uuid.UUID) (*State, error)
	GetStatesByNodeId(nodeId string) ([]State, error)
	GetStateHistory(nodeId string) ([]State, error)
	AddState(newState *State, previousState *State) error
	GetLatestState(nodeId string) (*State, error)
	UpdateState(nodeId string, subStates []string, events []string) (*State, error)
	GetNodeConfig(nodeId string) (*NodeConfig, error)
	GetStateHistoryWithFilter(nodeId string, startTime time.Time, endTime time.Time, pageSize int) ([]*State, error)
}

type stateRepo struct {
	Db sql.Db
}

func NewStateRepo(db sql.Db) StateRepo {
	return &stateRepo{
		Db: db,
	}
}

func (r *stateRepo) GetStateById(id uuid.UUID) (*State, error) {
	var state State
	err := r.Db.GetGormDb().First(&state, id).Error
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *stateRepo) GetCurrentNodeState(nodeId string) (*State, error) {
	var state State
	err := r.Db.GetGormDb().Where("node_id = ?", nodeId).Order("created_at DESC").First(&state).Error
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *stateRepo) GetStatesByNodeId(nodeId string) ([]State, error) {
	var states []State
	if err := r.Db.GetGormDb().Where("node_id = ?", nodeId).Order("created_at DESC").Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

func (r *stateRepo) GetStateHistory(nodeId string) ([]State, error) {
	var tempHistory []State
	currentState, err := r.GetCurrentNodeState(nodeId)
	if err != nil {
		return nil, err
	}

	for currentState != nil {
		tempHistory = append(tempHistory, *currentState)
		if currentState.PreviousStateId == nil {
			break
		}
		currentState, err = r.GetStateById(*currentState.PreviousStateId)
		if err != nil {
			return nil, err
		}
	}

	history := make([]State, len(tempHistory))
	for i, state := range tempHistory {
		history[len(tempHistory)-1-i] = state
	}

	return history, nil
}
func (r *stateRepo) GetLatestState(nodeId string) (*State, error) {
	var latestState State
	result := r.Db.GetGormDb().Where("node_id = ?", nodeId).Order("created_at DESC").First(&latestState)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching latest state: %w", result.Error)
	}

	return &latestState, nil
}
func (r *stateRepo) AddState(newState *State, previousState *State) error {
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

func (r *stateRepo) UpdateState(nodeId string, subStates []string, events []string) (*State, error) {
	var state State

	err := r.Db.GetGormDb().Where("node_id = ?", nodeId).Order("created_at DESC").First(&state).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching current state: %w", err)
	}
	state.SubState = subStates

	if len(events) > 0 {
		state.Events = events
	}

	if err := r.Db.GetGormDb().Save(&state).Error; err != nil {
		return nil, fmt.Errorf("error updating state: %w", err)
	}

	return &state, nil
}

func (r *stateRepo) GetNodeConfig(nodeId string) (*NodeConfig, error) {
	var config NodeConfig
	err := r.Db.GetGormDb().
		Where("node_id = ?", nodeId).
		Order("id DESC").
		First(&config).Error

	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *stateRepo) GetStateHistoryWithFilter(nodeId string, startTime time.Time, endTime time.Time, pageSize int) ([]*State, error) {
	var states []*State

	query := r.Db.GetGormDb().Model(&State{}).Where("node_id = ?", nodeId)

	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}

	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	query = query.Order("created_at DESC")

	if pageSize > 0 {
		query = query.Limit(pageSize)
	}

	err := query.Find(&states).Error
	if err != nil {
		return nil, err
	}

	return states, nil
}
