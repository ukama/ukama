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

// SiteStatusCounts holds site counts grouped by status.
type SiteStatusCounts struct {
	Total    int64
	Online   int64
	Degraded int64
	Offline  int64
}

type SiteRepo interface {
	List(networkId string, status string, page, pageSize uint32) ([]SiteSnapshot, int64, error)
	StatusCounts(networkId string) (*SiteStatusCounts, error)
	Get(siteId string) (*SiteSnapshot, error)
	CustomerCount(siteId string) (int64, error)
	UptimeBetween(siteId string, from, to time.Time) (float64, error)
	LatestSiteHealth(siteId string) (*SiteHealthRollupHourly, error)
	SiteHealthSeries(siteId string, from, to time.Time) ([]SiteHealthRollupHourly, error)
	Search(query, networkId string, limit int) ([]SiteSnapshot, error)
	OfflineDuration(siteId string) (float64, error)
}

type siteRepo struct {
	Db sql.Db
}

func NewSiteRepo(db sql.Db) SiteRepo {
	return &siteRepo{
		Db: db,
	}
}

func (r *siteRepo) List(networkId string, status string, page, pageSize uint32) ([]SiteSnapshot, int64, error) {
	var sites []SiteSnapshot
	var count int64

	q := r.Db.GetGormDb().Model(&SiteSnapshot{})
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
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

	if err := q.Order("name asc").Find(&sites).Error; err != nil {
		return nil, 0, err
	}

	return sites, count, nil
}

func (r *siteRepo) StatusCounts(networkId string) (*SiteStatusCounts, error) {
	type row struct {
		Status string
		Cnt    int64
	}

	var rows []row

	q := r.Db.GetGormDb().Model(&SiteSnapshot{}).
		Select("status, count(*) as cnt").
		Group("status")
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}

	counts := &SiteStatusCounts{}
	for _, rw := range rows {
		counts.Total += rw.Cnt

		switch rw.Status {
		case "online":
			counts.Online = rw.Cnt
		case "degraded":
			counts.Degraded = rw.Cnt
		case "offline":
			counts.Offline = rw.Cnt
		}
	}

	return counts, nil
}

func (r *siteRepo) Get(siteId string) (*SiteSnapshot, error) {
	var site SiteSnapshot

	result := r.Db.GetGormDb().Where("site_id = ?", siteId).First(&site)
	if result.Error != nil {
		return nil, result.Error
	}

	return &site, nil
}

func (r *siteRepo) CustomerCount(siteId string) (int64, error) {
	var count int64

	result := r.Db.GetGormDb().Model(&CustomerSnapshot{}).
		Where("site_id = ?", siteId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// UptimeBetween computes uptime percent for a site over [from, to] from
// site state intervals (online seconds / total window seconds). If no
// intervals exist for the window, it falls back to the average of the site
// health rollups over the same window.
func (r *siteRepo) UptimeBetween(siteId string, from, to time.Time) (float64, error) {
	total := to.Sub(from).Seconds()
	if total <= 0 {
		return 0, nil
	}

	type res struct {
		Seconds  float64
		Episodes int64
	}

	var online res

	// clamp intervals to the window; open intervals (end_at IS NULL) extend to `to`.
	err := r.Db.GetGormDb().Model(&SiteStateInterval{}).
		Select("coalesce(sum(extract(epoch from (least(coalesce(end_at, ?), ?) - greatest(start_at, ?)))), 0) as seconds, count(*) as episodes", to, to, from).
		Where("site_id = ? AND state = ? AND start_at < ? AND (end_at IS NULL OR end_at > ?)",
			siteId, "online", to, from).
		Find(&online).Error
	if err != nil {
		return 0, err
	}

	var overlapping int64

	err = r.Db.GetGormDb().Model(&SiteStateInterval{}).
		Where("site_id = ? AND start_at < ? AND (end_at IS NULL OR end_at > ?)", siteId, to, from).
		Count(&overlapping).Error
	if err != nil {
		return 0, err
	}

	if overlapping > 0 {
		uptime := online.Seconds / total * 100
		if uptime > 100 {
			uptime = 100
		}

		return uptime, nil
	}

	// fallback: average of site health rollups over the window.
	var avg *float64

	err = r.Db.GetGormDb().Model(&SiteHealthRollupHourly{}).
		Select("avg(uptime_percent)").
		Where("site_id = ? AND hour >= ? AND hour < ?", siteId, from, to).
		Scan(&avg).Error
	if err != nil {
		return 0, err
	}

	if avg == nil {
		return 0, nil
	}

	return *avg, nil
}

func (r *siteRepo) LatestSiteHealth(siteId string) (*SiteHealthRollupHourly, error) {
	var h SiteHealthRollupHourly

	result := r.Db.GetGormDb().Where("site_id = ?", siteId).
		Order("hour desc").First(&h)
	if result.Error != nil {
		return nil, result.Error
	}

	return &h, nil
}

// Search finds sites whose name or id matches the query (case-insensitive).
func (r *siteRepo) Search(query, networkId string, limit int) ([]SiteSnapshot, error) {
	var sites []SiteSnapshot

	pattern := "%" + query + "%"

	q := r.Db.GetGormDb().Model(&SiteSnapshot{}).
		Where("(name ILIKE ? OR site_id::text ILIKE ?)", pattern, pattern)
	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}

	if err := q.Order("name asc").Find(&sites).Error; err != nil {
		return nil, err
	}

	return sites, nil
}

// OfflineDuration returns how long (in seconds) a site has been in its
// currently open "offline" interval, or 0 if none is open.
func (r *siteRepo) OfflineDuration(siteId string) (float64, error) {
	var iv SiteStateInterval

	result := r.Db.GetGormDb().
		Where("site_id = ? AND state = ? AND end_at IS NULL", siteId, "offline").
		Order("start_at desc").First(&iv)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return 0, result.Error
	}

	return time.Since(iv.StartAt).Seconds(), nil
}

func (r *siteRepo) SiteHealthSeries(siteId string, from, to time.Time) ([]SiteHealthRollupHourly, error) {
	var rows []SiteHealthRollupHourly

	result := r.Db.GetGormDb().
		Where("site_id = ? AND hour >= ? AND hour < ?", siteId, from, to).
		Order("hour asc").Find(&rows)
	if result.Error != nil {
		return nil, result.Error
	}

	return rows, nil
}
