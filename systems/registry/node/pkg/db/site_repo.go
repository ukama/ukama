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

type SiteRepo interface {
	GetNodes(uuid.UUID) ([]Node, error)
	GetByNetwork(uuid.UUID) ([]Node, error)
	AddNode(*Site, func(*Site, *gorm.DB) error) error
	RemoveNode(ukama.NodeID) (*Site, error)
	GetFreeNodes() ([]Node, error)
	GetFreeNodesForOrg(uuid.UUID) ([]Node, error)
	IsAllocated(ukama.NodeID) (bool, *Site)
}

type siteRepo struct {
	Db sql.Db
}

func NewSiteRepo(db sql.Db) SiteRepo {
	return &siteRepo{
		Db: db,
	}
}

func (s *siteRepo) AddNode(node *Site, nestedFunc func(node *Site, tx *gorm.DB) error) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Create(node)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("node already belongs to site")
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

func (s *siteRepo) GetNodes(siteId uuid.UUID) ([]Node, error) {
	var nodes []Node

	result := s.Db.GetGormDb().Joins("JOIN sites on sites.node_id=nodes.id").
		Preload(clause.Associations).Preload("Attached.Site").
		Where("sites.site_id=? AND sites.deleted_at IS NULL",
			siteId.String()).Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return nodes, nil
}

func (s *siteRepo) GetByNetwork(networkId uuid.UUID) ([]Node, error) {
	var nodes []Node

	result := s.Db.GetGormDb().Joins("JOIN sites on sites.node_id=nodes.id").
		Preload(clause.Associations).Preload("Attached.Site").
		Where("sites.network_id=? AND sites.deleted_at IS NULL",
			networkId.String()).Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return nodes, nil
}

func (s *siteRepo) GetFreeNodes() ([]Node, error) {
	var nodes []Node

	result := s.Db.GetGormDb().
		Preload(clause.Associations).Preload("Attached.Site").Where("id NOT IN (?)",
		s.Db.GetGormDb().Table("sites").Select("node_id").Where("deleted_at IS NULL")).
		Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func (s *siteRepo) GetFreeNodesForOrg(orgId uuid.UUID) ([]Node, error) {
	var nodes []Node

	result := s.Db.GetGormDb().
		Preload(clause.Associations).Preload("Attached.Site").Where("id NOT IN (?) AND org_id= ?",
		s.Db.GetGormDb().Table("sites").Select("node_id").Where("deleted_at IS NULL"), orgId).Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func (s *siteRepo) RemoveNode(nodeId ukama.NodeID) (*Site, error) {
	ok, nd := s.IsAllocated(nodeId)
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition,
			"node is not yet assigned to site/network")
	}

	res := s.Db.GetGormDb().
		Exec("select * from nodes where parent_node_id= ?  OR (id= ? AND parent_node_id is NOT NULL)",
			nodeId.StringLowercase(), nodeId.StringLowercase())

	if res.Error != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to get node grouping result. error %s", res.Error.Error())
	}

	if res.RowsAffected > 0 {
		return nil, status.Errorf(codes.FailedPrecondition,
			"node is still grouped with other nodes.")
	}

	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("node_id= ?", nodeId.StringLowercase()).Delete(&Site{})
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return nd, err
}

func (s *siteRepo) IsAllocated(nodeId ukama.NodeID) (bool, *Site) {
	return isAllocated(s.Db.GetGormDb(), nodeId)
}

func isAllocated(db *gorm.DB, nodeId ukama.NodeID) (bool, *Site) {
	var nd Site

	result := db.Where(&Site{NodeId: nodeId.StringLowercase()}).First(&nd)
	return result.Error == nil, &nd
}
