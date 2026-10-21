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
	"errors"
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/gorm/clause"
)

var ErrNodeMissing = errors.New("node missing")

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

// Atomic upsert on lower(node_id); a node re-registering under a different
// org has its org_id moved instead of failing on the unique node_id index.
func (r *nodeRepo) AddOrUpdate(ctx context.Context, node *Node) error {
	result := r.Db.GetGormDb().WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "lower(node_id)", Raw: true},
		},
		DoUpdates: clause.AssignmentColumns([]string{"org_id", "updated_at"}),
	}).Create(node)

	return result.Error
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

	return fmt.Errorf("%s %w", nodeId.String(), ErrNodeMissing)
}
