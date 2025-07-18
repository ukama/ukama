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
)

type NodeLogRepo interface {
	Get(nodeId string) (*NodeLog, error)
	Add(nodeLog string) error
}

type nodeLogRepo struct {
	Db sql.Db
}

func NewNodeLogRepo(db sql.Db) NodeLogRepo {
	return &nodeLogRepo{
		Db: db,
	}
}

func (r *nodeLogRepo) Add(nodeId string) error {
	var nodeLog NodeLog
	if err := r.Db.GetGormDb().Where("node_id = ?", nodeId).First(&nodeLog).Error; err != nil {
		if err := r.Db.GetGormDb().Create(&NodeLog{NodeId: nodeId}).Error; err != nil {
			return err
		}
	} else {
		return errors.New("duplicate record: a record with the same nodeId already exists")
	}
	return nil
}

func (r *nodeLogRepo) Get(nodeId string) (*NodeLog, error) {
	var nodeLog NodeLog
	if err := r.Db.GetGormDb().Where("node_id = ?", nodeId).First(&nodeLog).Error; err != nil {
		return nil, err
	}
	return &nodeLog, nil
}
