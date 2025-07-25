/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"

	log "github.com/sirupsen/logrus"
)

type NodeStatusRepo interface {
	Update(*NodeStatus) error
	Get(ukama.NodeID) (*NodeStatus, error)
	Delete(ukama.NodeID) error
	GetAll() ([]NodeStatus, error)
	GetNodeCount() (onlineNodeCount, offlineNodeCount int64, err error)
}

type nodeStatusRepo struct {
	Db sql.Db
}

func NewNodeStatusRepo(db sql.Db) NodeStatusRepo {
	return &nodeStatusRepo{
		Db: db,
	}
}

func (n *nodeStatusRepo) Update(ns *NodeStatus) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		t := tx.Where("node_id = ?", ns.NodeId).Delete(&NodeStatus{})
		if t.RowsAffected > 0 {
			log.Debugf("Marking old state.")
		}

		result := tx.Create(ns)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (n *nodeStatusRepo) Delete(id ukama.NodeID) error {

	result := n.Db.GetGormDb().Where("node_id=?", id.StringLowercase()).Delete(&NodeStatus{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (n *nodeStatusRepo) Get(id ukama.NodeID) (*NodeStatus, error) {
	var ns NodeStatus

	result := n.Db.GetGormDb().First(&ns, "node_id=?", id.StringLowercase())

	if result.Error != nil {
		return nil, result.Error
	}

	return &ns, nil
}

func (n *nodeStatusRepo) GetAll() ([]NodeStatus, error) {
	var ns []NodeStatus

	result := n.Db.GetGormDb().Find(&ns)

	if result.Error != nil {
		return nil, result.Error
	}

	return ns, nil
}

func (n *nodeStatusRepo) GetNodeCount() (onlineNodeCount, offlineNodeCount int64, err error) {
	db := n.Db.GetGormDb()

	if err := db.Model(&NodeStatus{}).Where("conn = ?", ukama.NodeConnectivityOnline).Count(&onlineNodeCount).Error; err != nil {
		return 0, 0, err
	}

	if err := db.Model(&NodeStatus{}).Where("conn = ?", ukama.NodeConnectivityOffline).Count(&offlineNodeCount).Error; err != nil {
		return 0, 0, err
	}

	return onlineNodeCount, offlineNodeCount, nil
}
