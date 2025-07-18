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
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	result := n.Db.GetGormDb().Preload(clause.Associations).Preload("Attached.Site").Joins("INNER JOIN node_statuses ON nodes.id = node_statuses.node_id").
		Where("node_statuses.connectivity = ? AND node_statuses.state = ? AND node_statuses.deleted_at IS NULL", connectivity, status).
		Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func (n *nodeRepo) GetAll() ([]Node, error) {
	var nodes []Node

	result := n.Db.GetGormDb().Preload(clause.Associations).Preload("Attached.Site").Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func (n *nodeRepo) Delete(nodeId ukama.NodeID, nestedFunc func(ukama.NodeID, *gorm.DB) error) error {
	node, err := n.Get(nodeId)
	if err != nil {
		return fmt.Errorf("fail to get node: %w", err)
	}

	if len(node.Attached) > 0 {
		return fmt.Errorf("node %s still have child nodes", node.Id)
	}

	if node.ParentNodeId != nil {
		return fmt.Errorf("node %s is still attached to a parent node", node.Id)
	}

	if node.Site.SiteId != uuid.Nil {
		return fmt.Errorf("node is still assigned to site/network")
	}

	err = n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("node_id", nodeId.StringLowercase()).Delete(&NodeStatus{})
		if result.Error != nil {
			return result.Error
		}

		result = tx.Delete(&Node{Id: nodeId.StringLowercase()})
		if result.Error != nil {
			return result.Error
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

func (n *nodeRepo) List(nodeId, siteId, networkId, ntype string, connectivity, state *uint8) ([]Node, error) {
	var nodes []Node

	query := n.Db.GetGormDb().
		Preload(clause.Associations).
		Preload("Attached.Site").
		Select("nodes.*, node_statuses.connectivity, node_statuses.state, sites.site_id, sites.network_id").
		Joins("INNER JOIN node_statuses ON nodes.id = node_statuses.node_id").
		Joins("LEFT JOIN sites ON nodes.id = sites.node_id").
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

// Update updated node with `id`. Only fields that are not nil are updated, eg name and state.
func (n *nodeRepo) Update(node *Node, nestedFunc func(*Node, *gorm.DB) error) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Returning{}).Updates(node)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
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
		return fmt.Errorf("number of nodes (%d) to attach is not valid", len(attachedNodeIds))
	}

	parentNode, err := n.Get(nodeId)
	if err != nil {
		return fmt.Errorf("fail to get parent node: %w", err)
	}

	if parentNode.Type != ukama.NODE_ID_TYPE_TOWERNODE {
		return status.Errorf(codes.InvalidArgument,
			"parent node (%v) type must be a towernode", parentNode.Id)
	}

	if parentNode.Site.SiteId == uuid.Nil {
		return status.Errorf(codes.FailedPrecondition,
			"parent node (%v)  must belong to a site", parentNode.Id)
	}

	attachedNodes, err := n.batchGet(attachedNodeIds)
	if err != nil {
		return fmt.Errorf("fail to get list of nodes to attach: %w", err)
	}

	if len(attachedNodes) == 0 || len(attachedNodes) != len(attachedNodeIds) {
		return fmt.Errorf("some of the nodes from %v were not found or were duplicated",
			attachedNodeIds)
	}

	if parentNode.Attached == nil {
		parentNode.Attached = make([]*Node, 0)
	}

	if len(attachedNodes)+len(parentNode.Attached) > MaxAttachedNodes {
		return status.Errorf(codes.InvalidArgument,
			"max number of attached nodes should not be more than %d", MaxAttachedNodes)
	}

	err = n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		for _, an := range attachedNodes {
			if an.ParentNodeId != nil {
				return status.Errorf(codes.InvalidArgument,
					"node %v is already attached to a parent", an.Id)
			}

			if an.Type != ukama.NODE_ID_TYPE_AMPNODE {
				return status.Errorf(codes.InvalidArgument,
					"cannot attach non amplifier node: %v", an.Id)
			}

			if parentNode.Site.SiteId != an.Site.SiteId {
				return status.Errorf(codes.InvalidArgument,
					"cannot attach nodes from different sites")
			}

			an.ParentNodeId = &parentNode.Id

			d := tx.Save(&an)
			if d.Error != nil {
				return status.Errorf(codes.Internal,
					"failed to attach node: %s . Error %s", an.Id, d.Error.Error())
			}
		}

		return nil
	})

	return err
}

func (n *nodeRepo) DetachNode(detachNodeId ukama.NodeID) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		node, err := n.Get(detachNodeId)
		if err != nil {
			return fmt.Errorf("fail to get node: %w", err)
		}

		if node.ParentNodeId == nil {
			return fmt.Errorf("node %s is not attached to a parent node", node.Id)
		}

		if len(node.Attached) > 0 {
			return fmt.Errorf("node %s still have child nodes", node.Id)
		}

		node.ParentNodeId = nil

		d := tx.Save(&node)
		if d.Error != nil {
			return fmt.Errorf("failed to detach node: %s . Error %s", node.Id, d.Error.Error())
		}

		return nil
	})

	return err
}

func (n *nodeRepo) GetNodeCount() (nodeCount, onlineCount, offlineCount int64, err error) {
	db := n.Db.GetGormDb()

	if err := db.Model(&Node{}).Count(&nodeCount).Error; err != nil {
		return 0, 0, 0, err
	}

	res1 := db.Raw("select COUNT(*) from nodes LEFT JOIN node_statuses ON nodes.id = node_statuses.node_id WHERE node_statuses.connectivity = ? AND node_statuses.deleted_at IS NULL",
		ukama.NodeConnectivityOnline).Scan(&onlineCount)
	if res1.Error != nil {
		return 0, 0, 0, err
	}

	res2 := db.Raw("select COUNT(*) from nodes LEFT JOIN node_statuses ON nodes.id = node_statuses.node_id WHERE node_statuses.connectivity = ? AND node_statuses.deleted_at IS NULL",
		ukama.NodeConnectivityOffline).Scan(&offlineCount)
	if res2.Error != nil {
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
