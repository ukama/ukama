/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const MaxAttachedNodes = 2

type NodeRepo interface {
	Add(*Node, func(*Node, *gorm.DB) error) error
	Get(ukama.NodeID) (*Node, error)
	GetAll() ([]Node, error)
	GetNodesByState(connectivity, status uint8) ([]Node, error)
	List(nodeId, siteId, networkId, ntype string, connectivity, state *uint8) ([]Node, error)
	Delete(ukama.NodeID, func(ukama.NodeID, *gorm.DB) error) error
	Update(*Node, func(*Node, *gorm.DB) error) error
	AttachNodes(nodeId ukama.NodeID, attachedNodeId []string) error
	DetachNode(detachNodeId ukama.NodeID) error
	GetNodeCount() (int64, int64, int64, error)
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) NodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (n *nodeRepo) Add(node *Node, nestedFunc func(node *Node, tx *gorm.DB) error) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(node, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(node)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (n *nodeRepo) Get(id ukama.NodeID) (*Node, error) {
	var node Node

	result := n.Db.GetGormDb().Preload(clause.Associations).Preload("Attached.Site").
		First(&node, "id=?", id.StringLowercase())

	if result.Error != nil {
		return nil, result.Error
	}

	return &node, nil
}

func (n *nodeRepo) GetNodesByState(connectivity, status uint8) ([]Node, error) {
	var nodes []Node
	result := n.Db.GetGormDb().
		Preload(clause.Associations).
		Preload("Attached.Site").
		Joins("INNER JOIN node_statuses ON nodes.id = node_statuses.node_id").
		Where("node_statuses.connectivity = ? AND node_statuses.state = ? AND"+
			" node_statuses.deleted_at IS NULL",
			connectivity, status).
		Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func (n *nodeRepo) GetAll() ([]Node, error) {
	var nodes []Node

	result := n.Db.GetGormDb().
		Preload(clause.Associations).Preload("Attached.Site").Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func (n *nodeRepo) Delete(nodeId ukama.NodeID,
	nestedFunc func(ukama.NodeID, *gorm.DB) error) error {
	node, err := n.Get(nodeId)
	if err != nil {
		return fmt.Errorf("fail to get node: %w", err)
	}

	if len(node.Attached) > 0 {
		return fmt.Errorf("node %s still have child nodes: %w", node.Id,
			store.ErrPrecondition)
	}

	if node.ParentNodeId != nil {
		return fmt.Errorf("node %s is still attached to a parent node: %w", node.Id,
			store.ErrPrecondition)
	}

	err = n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("node_id = ?", nodeId.StringLowercase()).Delete(&Site{})
		if result.Error != nil {
			return result.Error
		}

		result = tx.Where("node_id = ?", nodeId.StringLowercase()).Delete(&NodeStatus{})
		if result.Error != nil {
			return result.Error
		}

		result = tx.Delete(&Node{Id: nodeId.StringLowercase()})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("node %s: %w", nodeId.StringLowercase(), store.ErrNotFound)
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(nodeId, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}

func (n *nodeRepo) List(nodeId, siteId, networkId, ntype string, connectivity,
	state *uint8) ([]Node, error) {
	var nodes []Node

	query := n.Db.GetGormDb().
		Preload(clause.Associations).
		Preload("Attached.Site").
		Select("nodes.*").
		Joins("INNER JOIN node_statuses ON nodes.id = node_statuses.node_id").
		// Sites are soft deleted on release/reassign: without this predicate a
		// released node still matches its old site_id/network_id filter.
		Joins("LEFT JOIN sites ON nodes.id = sites.node_id AND sites.deleted_at IS NULL").
		Where("node_statuses.deleted_at IS NULL")

	if nodeId != "" {
		query = query.Where("nodes.id = ?", nodeId)
	}

	if siteId != "" {
		query = query.Where("sites.site_id = ?", siteId)
	}

	if networkId != "" {
		query = query.Where("sites.network_id = ?", networkId)
	}

	if connectivity != nil {
		query = query.Where("node_statuses.connectivity = ?", connectivity)
	}

	if state != nil {
		query = query.Where("node_statuses.state = ?", state)
	}

	if ntype != "" {
		query = query.Where("nodes.type = ?", ntype)
	}

	result := query.Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

// Update updated node with `id`. Only non-zero fields on `node` are written
// (GORM's struct-based Updates semantics), eg an empty Name or Latitude is
// left untouched rather than clearing the existing value.
func (n *nodeRepo) Update(node *Node, nestedFunc func(*Node, *gorm.DB) error) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Returning{}).Updates(node)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("node %s: %w", node.Id, store.ErrNotFound)
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(node, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}

func (n *nodeRepo) AttachNodes(nodeId ukama.NodeID, attachedNodeIds []string) error {
	if len(attachedNodeIds) == 0 || len(attachedNodeIds) > MaxAttachedNodes {
		return fmt.Errorf("number of nodes (%d) to attach is not valid: %w",
			len(attachedNodeIds), store.ErrInvalid)
	}

	parentNode, err := n.Get(nodeId)
	if err != nil {
		return fmt.Errorf("fail to get parent node: %w", err)
	}

	if parentNode.Type != ukama.NODE_ID_TYPE_TOWERNODE {
		return fmt.Errorf("parent node (%v) type must be a towernode: %w", parentNode.Id,
			store.ErrInvalid)
	}

	if parentNode.Site.SiteId == uuid.Nil {
		return fmt.Errorf("parent node (%v) must belong to a site: %w", parentNode.Id,
			store.ErrPrecondition)
	}

	attachedNodes, err := n.batchGet(attachedNodeIds)
	if err != nil {
		return fmt.Errorf("fail to get list of nodes to attach: %w", err)
	}

	if len(attachedNodes) == 0 || len(attachedNodes) != len(attachedNodeIds) {
		return fmt.Errorf("some of the nodes from %v were not found: %w", attachedNodeIds,
			store.ErrInvalid)
	}

	err = n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if lockErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&Node{}, "id = ?", parentNode.Id).Error; lockErr != nil {
			return lockErr
		}

		var currentlyAttached int64
		if countErr := tx.Model(&Node{}).
			Where("parent_node_id = ?", parentNode.Id).
			Count(&currentlyAttached).Error; countErr != nil {
			return countErr
		}

		if currentlyAttached+int64(len(attachedNodes)) > MaxAttachedNodes {
			return fmt.Errorf("max number of attached nodes should not be more than %d: %w",
				MaxAttachedNodes, store.ErrPrecondition)
		}

		for _, id := range attachedNodeIds {
			var an Node

			if getErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Preload("Site").First(&an, "id = ?", id).Error; getErr != nil {
				return getErr
			}

			if an.ParentNodeId != nil {
				return fmt.Errorf("node %v is already attached to a parent: %w", an.Id,
					store.ErrPrecondition)
			}

			if an.Type != ukama.NODE_ID_TYPE_AMPNODE {
				return fmt.Errorf("cannot attach non amplifier node %v: %w", an.Id,
					store.ErrInvalid)
			}

			if parentNode.Site.SiteId != an.Site.SiteId {
				return fmt.Errorf("cannot attach nodes from different sites: %w",
					store.ErrInvalid)
			}

			an.ParentNodeId = &parentNode.Id

			d := tx.Save(&an)
			if d.Error != nil {
				return fmt.Errorf("failed to attach node %s: %w", an.Id, d.Error)
			}
		}

		return nil
	})

	return err
}

func (n *nodeRepo) DetachNode(detachNodeId ukama.NodeID) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		var node Node

		if getErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload(clause.Associations).Preload("Attached.Site").
			First(&node, "id = ?", detachNodeId.StringLowercase()).Error; getErr != nil {
			return fmt.Errorf("fail to get node: %w", getErr)
		}

		if node.ParentNodeId == nil {
			return fmt.Errorf("node %s is not attached to a parent node: %w", node.Id,
				store.ErrPrecondition)
		}

		if len(node.Attached) > 0 {
			return fmt.Errorf("node %s still have child nodes: %w", node.Id,
				store.ErrPrecondition)
		}

		node.ParentNodeId = nil

		d := tx.Save(&node)
		if d.Error != nil {
			return fmt.Errorf("failed to detach node %s: %w", node.Id, d.Error)
		}

		return nil
	})

	return err
}

func (n *nodeRepo) GetNodeCount() (nodeCount, onlineCount, offlineCount int64,
	err error) {
	db := n.Db.GetGormDb()

	if err := db.Model(&Node{}).Count(&nodeCount).Error; err != nil {
		return 0, 0, 0, err
	}

	// Use GORM Joins to match the same relationship that GetNode() uses via Preload
	// This ensures consistent behavior between single node lookups and node counts
	if err := db.Model(&Node{}).
		Joins("JOIN node_statuses ON LOWER(nodes.id) = LOWER(node_statuses.node_id) AND"+
			" node_statuses.deleted_at IS NULL").
		Where("node_statuses.connectivity = ?", ukama.NodeConnectivityOnline).
		Count(&onlineCount).Error; err != nil {
		return 0, 0, 0, err
	}

	if err := db.Model(&Node{}).
		Joins("JOIN node_statuses ON LOWER(nodes.id) = LOWER(node_statuses.node_id) AND"+
			" node_statuses.deleted_at IS NULL").
		Where("node_statuses.connectivity = ?", ukama.NodeConnectivityOffline).
		Count(&offlineCount).Error; err != nil {
		return 0, 0, 0, err
	}

	return nodeCount, onlineCount, offlineCount, nil
}

func (n *nodeRepo) batchGet(nodeIds []string) ([]Node, error) {
	var nodes []Node

	result := n.Db.GetGormDb().Preload(clause.Associations).Preload("Attached.Site").
		Where("id IN ?", nodeIds).Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}
