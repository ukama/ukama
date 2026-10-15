/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"context"
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/gorm/clause"
)

type NodeRepo interface {
	AddOrUpdate(ctx context.Context, node *Node) error
	Get(ctx context.Context, nodeId ukama.NodeID) (*Node, error)
	Delete(ctx context.Context, nodeId ukama.NodeID) error
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) *nodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (r *nodeRepo) AddOrUpdate(ctx context.Context, node *Node) error {
	d := r.Db.GetGormDb().WithContext(ctx).Create(node)
	return d.Error
}

func (r *nodeRepo) Get(ctx context.Context, nodeId ukama.NodeID) (*Node, error) {
	var node Node
	result := r.Db.GetGormDb().WithContext(ctx).Preload(clause.Associations).First(&node, "node_id = ?", nodeId.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}
	return &node, nil
}

func (r *nodeRepo) Delete(ctx context.Context, nodeId ukama.NodeID) error {
	var node Node
	result := r.Db.GetGormDb().WithContext(ctx).Unscoped().Delete(&node, "node_id = ?", nodeId.StringLowercase())
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	}

	return fmt.Errorf("%s node missing", nodeId.String())
}
