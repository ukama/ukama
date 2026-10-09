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

// Update updates the node's status row in place. Only non-zero fields on ns
// are written (GORM's struct-based Updates semantics), so a caller that only
// sets Connectivity leaves the existing State untouched, and vice versa.
func (n *nodeStatusRepo) Update(ns *NodeStatus) error {
	result := n.Db.GetGormDb().Model(&NodeStatus{}).
		Where("node_id = ?", ns.NodeId).
		Updates(ns)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
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

func (n *nodeStatusRepo) GetNodeCount() (onlineNodeCount, offlineNodeCount int64,
	err error) {
	db := n.Db.GetGormDb()

	if err := db.Model(&NodeStatus{}).Where("connectivity = ?",
		ukama.NodeConnectivityOnline).Count(&onlineNodeCount).Error; err != nil {
		return 0, 0, err
	}

	if err := db.Model(&NodeStatus{}).Where("connectivity = ?",
		ukama.NodeConnectivityOffline).Count(&offlineNodeCount).Error; err != nil {
		return 0, 0, err
	}

	return onlineNodeCount, offlineNodeCount, nil
}
