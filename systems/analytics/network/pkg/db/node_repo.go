/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"errors"
	"time"

	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

// NodeStatusCounts holds node counts grouped by status.
type NodeStatusCounts struct {
	Total          int64
	Online         int64
	Offline        int64
	Configuring    int64
	NeedsAttention int64
}

// NodePoolCounts holds inventory pool counts.
type NodePoolCounts struct {
	AvailableToInstall int64
	Deployed           int64
	InInventory        int64
	Rma                int64
}

type NodeRepo interface {
	List(networkId, siteId, status string, page, pageSize uint32) ([]NodeSnapshot, int64, error)
	ListAll(networkId string) ([]NodeSnapshot, error)
	StatusCounts(networkId, siteId string) (*NodeStatusCounts, error)
	Get(nodeId string) (*NodeSnapshot, error)
	UptimeBetween(nodeId string, from, to time.Time) (float64, error)
	PoolCounts(networkId string) (*NodePoolCounts, error)
	ConfiguringDuration(nodeId string) (float64, error)
	Search(query, networkId string, limit int) ([]NodeSnapshot, error)
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) NodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (r *nodeRepo) List(networkId, siteId, status string, page, pageSize uint32) ([]NodeSnapshot, int64, error) {
	var nodes []NodeSnapshot
	var count int64

	q := r.Db.GetGormDb().Model(&NodeSnapshot{})
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}
	if siteId != "" {
		q = q.Where("site_id = ?", siteId)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if pageSize > 0 {
		if page < 1 {
			page = 1
		}
		q = q.Offset(int((page - 1) * pageSize)).Limit(int(pageSize))
	}

	if err := q.Order("name asc").Find(&nodes).Error; err != nil {
		return nil, 0, err
	}

	return nodes, count, nil
}

func (r *nodeRepo) ListAll(networkId string) ([]NodeSnapshot, error) {
	var nodes []NodeSnapshot

	q := r.Db.GetGormDb().Model(&NodeSnapshot{})
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}

	if err := q.Order("name asc").Find(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *nodeRepo) StatusCounts(networkId, siteId string) (*NodeStatusCounts, error) {
	type row struct {
		Status string
		Cnt    int64
	}

	var rows []row

	q := r.Db.GetGormDb().Model(&NodeSnapshot{}).
		Select("status, count(*) as cnt").
		Group("status")
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}
	if siteId != "" {
		q = q.Where("site_id = ?", siteId)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}

	counts := &NodeStatusCounts{}
	for _, rw := range rows {
		counts.Total += rw.Cnt

		switch rw.Status {
		case "online":
			counts.Online = rw.Cnt
		case "offline":
			counts.Offline = rw.Cnt
		case "configuring":
			counts.Configuring = rw.Cnt
		case "needs_attention":
			counts.NeedsAttention = rw.Cnt
		}
	}

	return counts, nil
}

func (r *nodeRepo) Get(nodeId string) (*NodeSnapshot, error) {
	var node NodeSnapshot

	result := r.Db.GetGormDb().Where("node_id = ?", nodeId).First(&node)
	if result.Error != nil {
		return nil, result.Error
	}

	return &node, nil
}

// UptimeBetween computes uptime percent for a node over [from, to] from
// node state intervals (online seconds / total window seconds).
func (r *nodeRepo) UptimeBetween(nodeId string, from, to time.Time) (float64, error) {
	total := to.Sub(from).Seconds()
	if total <= 0 {
		return 0, nil
	}

	var seconds float64

	err := r.Db.GetGormDb().Model(&NodeStateInterval{}).
		Select("coalesce(sum(extract(epoch from (least(coalesce(end_at, ?), ?) - greatest(start_at, ?)))), 0)", to, to, from).
		Where("node_id = ? AND state = ? AND start_at < ? AND (end_at IS NULL OR end_at > ?)",
			nodeId, "online", to, from).
		Scan(&seconds).Error
	if err != nil {
		return 0, err
	}

	uptime := seconds / total * 100
	if uptime > 100 {
		uptime = 100
	}

	return uptime, nil
}

// PoolCounts derives the node pool KPIs from inventory snapshots
// (available/deployed/rma) plus node snapshots not yet assigned to a site
// (in inventory).
func (r *nodeRepo) PoolCounts(networkId string) (*NodePoolCounts, error) {
	type row struct {
		State string
		Cnt   int64
	}

	var rows []row

	err := r.Db.GetGormDb().Model(&InventorySnapshot{}).
		Select("state, count(*) as cnt").
		Group("state").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := &NodePoolCounts{}
	for _, rw := range rows {
		switch rw.State {
		case "available":
			counts.AvailableToInstall = rw.Cnt
		case "deployed":
			counts.Deployed = rw.Cnt
		case "rma":
			counts.Rma = rw.Cnt
		}
	}

	// nodes known to the network but not attached to any site are considered
	// still in inventory.
	q := r.Db.GetGormDb().Model(&NodeSnapshot{}).
		Where("site_id = ?", "00000000-0000-0000-0000-000000000000")
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}

	if err := q.Count(&counts.InInventory).Error; err != nil {
		return nil, err
	}

	return counts, nil
}

// Search finds nodes whose name or id matches the query (case-insensitive).
func (r *nodeRepo) Search(query, networkId string, limit int) ([]NodeSnapshot, error) {
	var nodes []NodeSnapshot

	pattern := "%" + query + "%"

	q := r.Db.GetGormDb().Model(&NodeSnapshot{}).
		Where("(name ILIKE ? OR node_id ILIKE ?)", pattern, pattern)
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}

	if err := q.Order("name asc").Find(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

// ConfiguringDuration returns how long (in seconds) a node has been in its
// currently open "configuring" interval, or 0 if none is open.
func (r *nodeRepo) ConfiguringDuration(nodeId string) (float64, error) {
	var iv NodeStateInterval

	result := r.Db.GetGormDb().
		Where("node_id = ? AND state = ? AND end_at IS NULL", nodeId, "configuring").
		Order("start_at desc").First(&iv)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return 0, result.Error
	}

	return time.Since(iv.StartAt).Seconds(), nil
}
