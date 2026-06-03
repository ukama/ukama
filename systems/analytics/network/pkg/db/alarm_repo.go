/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

// AlarmCounts holds open alarm counts grouped by severity.
type AlarmCounts struct {
	Open     int64
	Critical int64
	Warning  int64
}

// AlarmFilter narrows alarm listing.
type AlarmFilter struct {
	NetworkId string
	SiteId    string
	Severity  string /* critical|warning */
	State     string /* open|closed */
	Page      uint32
	PageSize  uint32
}

type AlarmRepo interface {
	List(filter AlarmFilter) ([]AlarmEvent, int64, error)
	Counts(networkId, siteId string) (*AlarmCounts, error)
	ForResource(resourceType, resourceId string, limit int) ([]AlarmEvent, error)
	OpenImpact(networkId string) (customersAffected int64, revenueAtRisk float64, err error)
}

type alarmRepo struct {
	Db sql.Db
}

func NewAlarmRepo(db sql.Db) AlarmRepo {
	return &alarmRepo{
		Db: db,
	}
}

// siteScopedAlarms returns a query over alarms attached either directly to a
// site (resource_type=site, resource_id=siteId) or to one of the network's
// sites/nodes. Network scoping is best effort: alarms only carry resource
// ids, so we join through snapshots when a network filter is provided.
func (r *alarmRepo) scoped(networkId, siteId string) *gorm.DB {
	q := r.Db.GetGormDb().Model(&AlarmEvent{})

	if siteId != "" {
		q = q.Where(
			"(resource_id = ?"+
				" OR resource_id IN (SELECT node_id FROM analytics_node_snapshots WHERE site_id = ?))",
			siteId, siteId)
	} else if networkId != "" {
		q = q.Where(
			"(resource_id IN (SELECT site_id::text FROM analytics_site_snapshots WHERE network_id = ?)"+
				" OR resource_id IN (SELECT node_id FROM analytics_node_snapshots WHERE network_id = ?))",
			networkId, networkId)
	}

	return q
}

func (r *alarmRepo) List(filter AlarmFilter) ([]AlarmEvent, int64, error) {
	var alarms []AlarmEvent
	var count int64

	q := r.scoped(filter.NetworkId, filter.SiteId)

	if filter.Severity != "" {
		q = q.Where("severity = ?", filter.Severity)
	}
	if filter.State != "" {
		q = q.Where("state = ?", filter.State)
	}

	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if filter.PageSize > 0 {
		page := filter.Page
		if page < 1 {
			page = 1
		}
		q = q.Offset(int((page - 1) * filter.PageSize)).Limit(int(filter.PageSize))
	}

	if err := q.Order("opened_at desc").Find(&alarms).Error; err != nil {
		return nil, 0, err
	}

	return alarms, count, nil
}

func (r *alarmRepo) Counts(networkId, siteId string) (*AlarmCounts, error) {
	type row struct {
		Severity string
		Cnt      int64
	}

	var rows []row

	q := r.scoped(networkId, siteId).
		Select("severity, count(*) as cnt").
		Where("state = ?", "open").
		Group("severity")

	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}

	counts := &AlarmCounts{}
	for _, rw := range rows {
		counts.Open += rw.Cnt

		switch rw.Severity {
		case "critical":
			counts.Critical = rw.Cnt
		case "warning":
			counts.Warning = rw.Cnt
		}
	}

	return counts, nil
}

func (r *alarmRepo) ForResource(resourceType, resourceId string, limit int) ([]AlarmEvent, error) {
	var alarms []AlarmEvent

	q := r.Db.GetGormDb().Model(&AlarmEvent{})
	if resourceType != "" {
		q = q.Where("resource_type = ?", resourceType)
	}
	if resourceId != "" {
		q = q.Where("resource_id = ?", resourceId)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}

	if err := q.Order("opened_at desc").Find(&alarms).Error; err != nil {
		return nil, err
	}

	return alarms, nil
}

// OpenImpact sums customers affected and revenue at risk over open alarms.
func (r *alarmRepo) OpenImpact(networkId string) (int64, float64, error) {
	type row struct {
		Customers int64
		Revenue   float64
	}

	var res row

	q := r.scoped(networkId, "").
		Select("coalesce(sum(customers_affected), 0) as customers, coalesce(sum(revenue_at_risk), 0) as revenue").
		Where("state = ?", "open")

	if err := q.Find(&res).Error; err != nil {
		return 0, 0, err
	}

	return res.Customers, res.Revenue, nil
}
