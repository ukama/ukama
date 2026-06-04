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

// MetricName describes one metric available for a resource, with the last
// sample time used for freshness checks.
type MetricName struct {
	Metric       string
	Unit         string
	LastSampleAt time.Time
}

type MetricRepo interface {
	Rollups(metric, resourceType, resourceId string, from, to time.Time) ([]MetricRollupHourly, error)
	LatestSamples(resourceId string, limit int) ([]MetricSample, error)
	MetricNames(resourceId string) ([]MetricName, error)
	RadioRollups(nodeId string, from, to time.Time) ([]RadioRollupHourly, error)
	LatestRadioRollup(nodeId string) (*RadioRollupHourly, error)
	RadioRollupSums(networkId string, at time.Time) (activeUes int64, attachFailures int64, err error)
	BackhaulRollups(siteId string, from, to time.Time) ([]BackhaulRollupHourly, error)
	LatestBackhaulRollup(siteId string) (*BackhaulRollupHourly, error)
	PowerRollups(siteId string, from, to time.Time) ([]PowerRollupHourly, error)
	LatestPowerRollup(siteId string) (*PowerRollupHourly, error)
}

type metricRepo struct {
	Db sql.Db
}

func NewMetricRepo(db sql.Db) MetricRepo {
	return &metricRepo{
		Db: db,
	}
}

func (r *metricRepo) Rollups(metric, resourceType, resourceId string, from, to time.Time) ([]MetricRollupHourly, error) {
	var rows []MetricRollupHourly

	q := r.Db.GetGormDb().Model(&MetricRollupHourly{}).
		Where("hour >= ? AND hour < ?", from, to)
	if metric != "" {
		q = q.Where("metric = ?", metric)
	}
	if resourceType != "" {
		q = q.Where("resource_type = ?", resourceType)
	}
	if resourceId != "" {
		q = q.Where("resource_id = ?", resourceId)
	}

	if err := q.Order("hour asc").Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *metricRepo) LatestSamples(resourceId string, limit int) ([]MetricSample, error) {
	var rows []MetricSample

	q := r.Db.GetGormDb().Model(&MetricSample{}).
		Where("resource_id = ?", resourceId).
		Order("sampled_at desc")
	if limit > 0 {
		q = q.Limit(limit)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *metricRepo) MetricNames(resourceId string) ([]MetricName, error) {
	var rows []MetricName

	q := r.Db.GetGormDb().Model(&MetricSample{}).
		Select("metric, max(unit) as unit, max(sampled_at) as last_sample_at").
		Group("metric").
		Order("metric asc")
	if resourceId != "" {
		q = q.Where("resource_id = ?", resourceId)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *metricRepo) RadioRollups(nodeId string, from, to time.Time) ([]RadioRollupHourly, error) {
	var rows []RadioRollupHourly

	q := r.Db.GetGormDb().Model(&RadioRollupHourly{}).
		Where("hour >= ? AND hour < ?", from, to)
	if nodeId != "" {
		q = q.Where("node_id = ?", nodeId)
	}

	if err := q.Order("hour asc").Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *metricRepo) LatestRadioRollup(nodeId string) (*RadioRollupHourly, error) {
	var row RadioRollupHourly

	result := r.Db.GetGormDb().Where("node_id = ?", nodeId).
		Order("hour desc").First(&row)
	if result.Error != nil {
		return nil, result.Error
	}

	return &row, nil
}

// RadioRollupSums sums the latest hour's radio rollups over all the
// network's nodes (active UEs, attach failures).
func (r *metricRepo) RadioRollupSums(networkId string, at time.Time) (int64, int64, error) {
	type row struct {
		ActiveUes      int64
		AttachFailures int64
	}

	var res row

	q := r.Db.GetGormDb().Model(&RadioRollupHourly{}).
		Select("coalesce(sum(active_ues), 0) as active_ues, coalesce(sum(attach_failures), 0) as attach_failures").
		Where("hour = (SELECT max(hour) FROM analytics_radio_rollup_hourly WHERE hour <= ?)", at)
	if networkId != "" {
		q = q.Where("node_id IN (SELECT node_id FROM analytics_node_snapshots WHERE network_id = ?)", networkId)
	}

	if err := q.Find(&res).Error; err != nil {
		return 0, 0, err
	}

	return res.ActiveUes, res.AttachFailures, nil
}

func (r *metricRepo) BackhaulRollups(siteId string, from, to time.Time) ([]BackhaulRollupHourly, error) {
	var rows []BackhaulRollupHourly

	q := r.Db.GetGormDb().Model(&BackhaulRollupHourly{}).
		Where("hour >= ? AND hour < ?", from, to)
	if siteId != "" {
		q = q.Where("site_id = ?", siteId)
	}

	if err := q.Order("hour asc").Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *metricRepo) LatestBackhaulRollup(siteId string) (*BackhaulRollupHourly, error) {
	var row BackhaulRollupHourly

	result := r.Db.GetGormDb().Where("site_id = ?", siteId).
		Order("hour desc").First(&row)
	if result.Error != nil {
		return nil, result.Error
	}

	return &row, nil
}

func (r *metricRepo) PowerRollups(siteId string, from, to time.Time) ([]PowerRollupHourly, error) {
	var rows []PowerRollupHourly

	q := r.Db.GetGormDb().Model(&PowerRollupHourly{}).
		Where("hour >= ? AND hour < ?", from, to)
	if siteId != "" {
		q = q.Where("site_id = ?", siteId)
	}

	if err := q.Order("hour asc").Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *metricRepo) LatestPowerRollup(siteId string) (*PowerRollupHourly, error) {
	var row PowerRollupHourly

	result := r.Db.GetGormDb().Where("site_id = ?", siteId).
		Order("hour desc").First(&row)
	if result.Error != nil {
		return nil, result.Error
	}

	return &row, nil
}
