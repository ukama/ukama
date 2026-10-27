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
	"sort"
	"strings"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
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
}

type siteRepo struct {
	Db sql.Db
}

func NewSiteRepo(db sql.Db) SiteRepo {
	return &siteRepo{
		Db: db,
	}
}

// AssignNodes assigns the nodes to the site in one transaction and returns only the new
// assignments. A node already on this site is skipped, and a node on another site fails
// the whole call.
func (s *siteRepo) AssignNodes(siteId, networkId uuid.UUID,
	nodeIds []string) ([]Site, error) {
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
				return fmt.Errorf("node %s: %w", id, store.ErrNotFound)
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

				return fmt.Errorf("node %s already belongs to site %s: %w", id,
					current.SiteId, store.ErrAlreadyExists)
			}

			site := Site{NodeId: id, SiteId: siteId, NetworkId: networkId}

			// Site.BeforeSave adds ON CONFLICT DO NOTHING, so a conflict shows up as zero rows.
			result = tx.Create(&site)
			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return fmt.Errorf("node %s already has a site assignment: %w", id,
					store.ErrAlreadyExists)
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

	result := withNodeAssociations(s.Db.GetGormDb()).
		Select("nodes.*, node_statuses.connectivity, node_statuses.state, sites.site_id,"+
			" sites.network_id").
		Joins("INNER JOIN node_statuses ON nodes.id = node_statuses.node_id").
		Joins("JOIN sites ON sites.node_id = nodes.id AND sites.deleted_at IS NULL").
		Where("sites.site_id=? AND sites.deleted_at IS NULL AND node_statuses.deleted_at IS"+
			" NULL",
			siteId.String()).Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("nodes for site %s: %w", siteId, store.ErrNotFound)
	}

	return nodes, nil
}

func (s *siteRepo) GetByNetwork(networkId uuid.UUID) ([]Node, error) {
	var nodes []Node

	result := withNodeAssociations(s.Db.GetGormDb()).
		Select("nodes.*, node_statuses.connectivity, node_statuses.state, sites.site_id,"+
			" sites.network_id").
		Joins("INNER JOIN node_statuses ON nodes.id = node_statuses.node_id").
		Joins("JOIN sites ON sites.node_id = nodes.id AND sites.deleted_at IS NULL").
		Where("sites.network_id=? AND sites.deleted_at IS NULL AND node_statuses.deleted_at"+
			" IS NULL",
			networkId.String()).Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("nodes for network %s: %w", networkId, store.ErrNotFound)
	}

	return nodes, nil
}

func (s *siteRepo) GetFreeNodes() ([]Node, error) {
	var nodes []Node

	result := withNodeAssociations(s.Db.GetGormDb()).Where("id NOT IN (?)",
		s.Db.GetGormDb().Table("sites").Select("node_id").Where("deleted_at IS NULL")).
		Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func (s *siteRepo) RemoveNode(nodeId ukama.NodeID) (*Site, error) {
	id := nodeId.StringLowercase()

	var removed Site

	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		lockErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").
			First(&Node{}, "id = ?", id).Error
		if errors.Is(lockErr, gorm.ErrRecordNotFound) {
			return fmt.Errorf("node %s: %w", id, store.ErrNotFound)
		}
		if lockErr != nil {
			return lockErr
		}

		siteErr := tx.Where("node_id = ?", id).First(&removed).Error
		if errors.Is(siteErr, gorm.ErrRecordNotFound) {
			return fmt.Errorf("node %s is not yet assigned to site/network: %w", id,
				store.ErrPrecondition)
		}
		if siteErr != nil {
			return siteErr
		}

		var grouped int64
		if countErr := tx.Model(&Node{}).
			Where("parent_node_id = ? OR (id = ? AND parent_node_id IS NOT NULL)", id, id).
			Count(&grouped).Error; countErr != nil {
			return countErr
		}

		if grouped > 0 {
			return fmt.Errorf("node %s is still grouped with other nodes: %w", id,
				store.ErrPrecondition)
		}

		return tx.Where("node_id = ?", id).Delete(&Site{}).Error
	})
	if err != nil {
		return nil, err
	}

	return &removed, nil
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
