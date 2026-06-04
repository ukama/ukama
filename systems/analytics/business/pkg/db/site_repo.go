/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db

import (
	"time"

	"github.com/ukama/ukama/systems/common/sql"
)

// SiteRepo is a read-only repository over site snapshots, site rollups and
// site health rollups.
type SiteRepo interface {
	ListSites(networkId string, page, pageSize int) ([]SiteSnapshot, int64, error)
	GetSite(siteId string) (*SiteSnapshot, error)
	SiteRollups(siteId string, from, to time.Time) ([]BusinessSiteRollupDaily, error)
	SiteUptime(siteId string, from, to time.Time) (float64, error)
}

type siteRepo struct {
	Db sql.Db
}

func NewSiteRepo(db sql.Db) SiteRepo {
	return &siteRepo{
		Db: db,
	}
}

func (s siteRepo) ListSites(networkId string, page, pageSize int) ([]SiteSnapshot, int64, error) {
	var sites []SiteSnapshot
	var count int64

	countQ := s.Db.GetGormDb().Model(&SiteSnapshot{})
	if networkId != "" {
		countQ = countQ.Where("network_id = ?", networkId)
	}

	if err := countQ.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	q := s.Db.GetGormDb().Model(&SiteSnapshot{}).Order("name ASC")
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}

	if pageSize > 0 {
		if page < 1 {
			page = 1
		}
		q = q.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	result := q.Find(&sites)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return sites, count, nil
}

func (s siteRepo) GetSite(siteId string) (*SiteSnapshot, error) {
	var site SiteSnapshot

	result := s.Db.GetGormDb().Where("site_id = ?", siteId).First(&site)
	if result.Error != nil {
		return nil, result.Error
	}

	return &site, nil
}

func (s siteRepo) SiteRollups(siteId string, from, to time.Time) ([]BusinessSiteRollupDaily, error) {
	var rollups []BusinessSiteRollupDaily

	q := s.Db.GetGormDb().Model(&BusinessSiteRollupDaily{}).
		Where("day >= ? AND day < ?", from, to).
		Order("day ASC")

	if siteId != "" {
		q = q.Where("site_id = ?", siteId)
	}

	result := q.Find(&rollups)
	if result.Error != nil {
		return nil, result.Error
	}

	return rollups, nil
}

func (s siteRepo) SiteUptime(siteId string, from, to time.Time) (float64, error) {
	var uptime float64

	q := s.Db.GetGormDb().Model(&SiteHealthRollupHourly{}).
		Select("COALESCE(AVG(uptime_percent), 0)").
		Where("hour >= ? AND hour < ?", from, to)

	if siteId != "" {
		q = q.Where("site_id = ?", siteId)
	}

	result := q.Scan(&uptime)
	if result.Error != nil {
		return 0, result.Error
	}

	return uptime, nil
}
