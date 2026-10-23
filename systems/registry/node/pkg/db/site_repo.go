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
	"sort"
	"strings"

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
	AssignNodes(siteId, networkId uuid.UUID, nodeIds []string) ([]Site, error)
	RemoveNode(ukama.NodeID) (*Site, error)
	ReleaseSite(siteId uuid.UUID) ([]Site, error)
	GetFreeNodes() ([]Node, error)
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

// AssignNodes assigns the nodes to the site in one transaction and returns only the new assignments.
// A node already on this site is skipped, and a node on another site fails the whole call.
func (s *siteRepo) AssignNodes(siteId, networkId uuid.UUID, nodeIds []string) ([]Site, error) {
	ids := make([]string, 0, len(nodeIds))
	for _, id := range nodeIds {
		ids = append(ids, strings.ToLower(id))
	}

	// A fixed lock order keeps two concurrent assignments from deadlocking.
	sort.Strings(ids)

	var assigned []Site

	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			var node Node

			result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").
				First(&node, "id = ?", id)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.NotFound, "node %s not found", id)
			}
			if result.Error != nil {
				return result.Error
			}

			var current Site

			// Match the unique index, which is on lower(node_id).
			result = tx.Where("lower(node_id) = ?", id).Limit(1).Find(&current)
			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected > 0 {
				if current.SiteId == siteId {
					continue
				}

				return status.Errorf(codes.AlreadyExists,
					"node %s already belongs to site %s", id, current.SiteId)
			}

			site := Site{NodeId: id, SiteId: siteId, NetworkId: networkId}

			// Site.BeforeSave adds ON CONFLICT DO NOTHING, so a conflict shows up as zero rows.
			result = tx.Create(&site)
			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return status.Errorf(codes.AlreadyExists, "node %s already has a site assignment", id)
			}

			assigned = append(assigned, site)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return assigned, nil
}

func (s *siteRepo) GetNodes(siteId uuid.UUID) ([]Node, error) {
	var nodes []Node

	result := s.Db.GetGormDb().
		Preload(clause.Associations).Preload("Attached.Site").
		Select("nodes.*, node_statuses.connectivity, node_statuses.state, sites.site_id, sites.network_id").
		Joins("INNER JOIN node_statuses ON nodes.id = node_statuses.node_id").
		Joins("JOIN sites ON sites.node_id = nodes.id AND sites.deleted_at IS NULL").
		Where("sites.site_id=? AND sites.deleted_at IS NULL AND node_statuses.deleted_at IS NULL",
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

	result := s.Db.GetGormDb().
		Preload(clause.Associations).Preload("Attached.Site").
		Select("nodes.*, node_statuses.connectivity, node_statuses.state, sites.site_id, sites.network_id").
		Joins("INNER JOIN node_statuses ON nodes.id = node_statuses.node_id").
		Joins("JOIN sites ON sites.node_id = nodes.id AND sites.deleted_at IS NULL").
		Where("sites.network_id=? AND sites.deleted_at IS NULL AND node_statuses.deleted_at IS NULL",
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

// ReleaseSite releases every node assigned to the site and returns their assignments.
// A site with no assigned nodes is not an error, so callers can retry.
func (s *siteRepo) ReleaseSite(siteId uuid.UUID) ([]Site, error) {
	var released []Site

	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("site_id = ?", siteId).Find(&released)
		if result.Error != nil {
			return result.Error
		}

		if len(released) == 0 {
			return nil
		}

		return tx.Where("site_id = ?", siteId).Delete(&Site{}).Error
	})
	if err != nil {
		return nil, err
	}

	return released, nil
}

func (s *siteRepo) IsAllocated(nodeId ukama.NodeID) (bool, *Site) {
	return isAllocated(s.Db.GetGormDb(), nodeId)
}

func isAllocated(db *gorm.DB, nodeId ukama.NodeID) (bool, *Site) {
	var nd Site

	result := db.Where(&Site{NodeId: nodeId.StringLowercase()}).First(&nd)
	return result.Error == nil, &nd
}
