/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConfigRepo interface {
	Add(id string) error
	Get(id string) (*Configuration, error)
	GetAll() ([]Configuration, error)
	Delete(id string) error
	//Update(c Configuration) error
	UpdateCurrentCommit(c Configuration, state *CommitState) error
	UpdateLastCommit(c Configuration, state *CommitState) error
	UpdateLastCommitState(nodeid string, state CommitState) error
	UpdateCommitState(nodeid string, state CommitState) error
}

type configRepo struct {
	Db sql.Db
}

func NewConfigRepo(db sql.Db) ConfigRepo {
	return &configRepo{
		Db: db,
	}
}

func (n *configRepo) Add(node string) error {
	config := Configuration{
		NodeId: node,
		State:  Default,
	}

	r := n.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "node_id"}},
		DoNothing: true,
	}).Create(&config)

	return r.Error
}

func (n *configRepo) Get(id string) (*Configuration, error) {
	var config Configuration

	result := n.Db.GetGormDb().Preload("Commit").First(&config, "node_id=?", strings.ToLower(id))
	if result.Error != nil {
		return nil, result.Error
	}

	return &config, nil
}

func (n *configRepo) GetAll() ([]Configuration, error) {
	var configs []Configuration

	result := n.Db.GetGormDb().Preload("Commit").Find(&configs)

	if result.Error != nil {
		return nil, result.Error
	}

	return configs, nil
}

func (n *configRepo) Delete(id string) error {
	var configs Configuration
	result := n.Db.GetGormDb().Where("node_id=?", strings.ToLower(id)).Delete(&configs)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update updated node with `id`. Only fields that are not nil are updated, eg name and state.
// func (n *configRepo) Update(c Configuration) error {

// 	result := n.Db.GetGormDb().Where("node_id=?", strings.ToLower(c.NodeId)).Updates(&c)
// 	if result.Error != nil {
// 		return result.Error
// 	}

//		return result.Error
//	}
//
// TODO: Check this one.
func (n *configRepo) UpdateLastCommit(c Configuration, state *CommitState) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&c).Association("LastCommit").Replace(&(c.Commit))
		if err != nil {
			log.Errorf("Failed to update last commit.Error: %v", err)
			return err
		}

		if state != nil {
			res := tx.Model(&Configuration{}).Where("node_id = ?", c.NodeId).Update("last_commit_state", *state)
			if res.Error != nil {
				log.Errorf("Failed to update configuration for node %s. Error: %v", c.NodeId, res.Error)
				return res.Error
			}
		}

		return nil

	})

	return err
}

func (n *configRepo) UpdateCurrentCommit(c Configuration, state *CommitState) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&c).Association("Commit").Replace(&(c.Commit))
		if err != nil {
			log.Errorf("Failed to update commit.Error: %v", err)
			return err
		}

		if state != nil {
			res := tx.Model(&Configuration{}).Where("node_id = ?", c.NodeId).Update("state", *state)
			if res.Error != nil {
				log.Errorf("Failed to update state. Error: %v", res.Error)
				return res.Error
			}
		}

		return nil

	})

	return err
}

func (n *configRepo) UpdateLastCommitState(nodeid string, state CommitState) error {

	result := n.Db.GetGormDb().Where("node_id=?", nodeid).Updates(&Configuration{LastCommitState: state})
	if result.Error != nil {
		return result.Error
	}

	return result.Error
}

func (n *configRepo) UpdateCommitState(nodeid string, state CommitState) error {

	result := n.Db.GetGormDb().Where("node_id=?", nodeid).Updates(&Configuration{State: state})
	if result.Error != nil {
		return result.Error
	}

	return result.Error
}
