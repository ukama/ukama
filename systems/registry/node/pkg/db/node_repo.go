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
	"slices"
	"sort"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/ukama"

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
	id := nodeId.StringLowercase()

	return n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		var node Node

		lockErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id", "parent_node_id").First(&node, "id = ?", id).Error
		if errors.Is(lockErr, gorm.ErrRecordNotFound) {
			return fmt.Errorf("node %s: %w", id, store.ErrNotFound)
		}
		if lockErr != nil {
			return lockErr
		}

		if node.ParentNodeId != nil {
			return fmt.Errorf("node %s is still attached to a parent node: %w", id,
				store.ErrPrecondition)
		}

		var children int64
		if countErr := tx.Model(&Node{}).Where("parent_node_id = ?", id).
			Count(&children).Error; countErr != nil {
			return countErr
		}

		if children > 0 {
			return fmt.Errorf("node %s still have child nodes: %w", id,
				store.ErrPrecondition)
		}

		result := tx.Where("node_id = ?", id).Delete(&Site{})
		if result.Error != nil {
			return result.Error
		}

		result = tx.Where("node_id = ?", id).Delete(&NodeStatus{})
		if result.Error != nil {
			return result.Error
		}

		result = tx.Delete(&Node{Id: id})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("node %s: %w", id, store.ErrNotFound)
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(nodeId, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})
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

	ids := slices.Clone(attachedNodeIds)
	sort.Strings(ids)

	lockIds := append([]string{parentNode.Id}, ids...)
	sort.Strings(lockIds)

	err = n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		for _, id := range lockIds {
			lockErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").
				First(&Node{}, "id = ?", id).Error
			if errors.Is(lockErr, gorm.ErrRecordNotFound) {
				return fmt.Errorf("node %s: %w", id, store.ErrNotFound)
			}
			if lockErr != nil {
				return lockErr
			}
		}

		var parentSite Site
		if siteErr := tx.Where("node_id = ?", parentNode.Id).First(&parentSite).
			Error; siteErr != nil {
			if errors.Is(siteErr, gorm.ErrRecordNotFound) {
				return fmt.Errorf("parent node (%v) must belong to a site: %w",
					parentNode.Id, store.ErrPrecondition)
			}

			return siteErr
		}

		var currentlyAttached int64
		if countErr := tx.Model(&Node{}).
			Where("parent_node_id = ?", parentNode.Id).
			Count(&currentlyAttached).Error; countErr != nil {
			return countErr
		}

		if currentlyAttached+int64(len(ids)) > MaxAttachedNodes {
			return fmt.Errorf("max number of attached nodes should not be more than %d: %w",
				MaxAttachedNodes, store.ErrPrecondition)
		}

		for _, id := range ids {
			var an Node

			if getErr := tx.Preload("Site").First(&an, "id = ?", id).Error; getErr != nil {
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

			if parentSite.SiteId != an.Site.SiteId {
				return fmt.Errorf("cannot attach nodes from different sites: %w",
					store.ErrInvalid)
			}

			updErr := tx.Model(&Node{}).Where("id = ?", an.Id).
				Update("parent_node_id", parentNode.Id).Error
			if updErr != nil {
				return fmt.Errorf("failed to attach node %s: %w", an.Id, updErr)
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
			if errors.Is(getErr, gorm.ErrRecordNotFound) {
				return fmt.Errorf("node %s: %w", detachNodeId.StringLowercase(),
					store.ErrNotFound)
			}

			return getErr
		}

		if node.ParentNodeId == nil {
			return fmt.Errorf("node %s is not attached to a parent node: %w", node.Id,
				store.ErrPrecondition)
		}

		if len(node.Attached) > 0 {
			return fmt.Errorf("node %s still have child nodes: %w", node.Id,
				store.ErrPrecondition)
		}

		updErr := tx.Model(&Node{}).Where("id = ?", node.Id).
			Update("parent_node_id", nil).Error
		if updErr != nil {
			return fmt.Errorf("failed to detach node %s: %w", node.Id, updErr)
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
