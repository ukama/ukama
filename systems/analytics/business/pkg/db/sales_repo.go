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

// DayValue is a single day bucket of an aggregated value.
type DayValue struct {
	Day   time.Time
	Value float64
}

// NamedAmount is an aggregated value attached to an entity (site/package).
type NamedAmount struct {
	Id    string
	Name  string
	Value float64
}

// SalesRepo is a read-only repository over payment events and sales rollups.
type SalesRepo interface {
	RevenueBetween(networkId string, from, to time.Time) (float64, error)
	PurchasesBetween(networkId string, from, to time.Time) (uint32, error)
	PaidCustomersBetween(networkId string, from, to time.Time) (uint32, error)
	RevenueTrendDaily(networkId string, from, to time.Time) ([]DayValue, error)
	RevenueBySite(networkId string, from, to time.Time) ([]NamedAmount, error)
	RevenueByPackage(networkId string, from, to time.Time) ([]NamedAmount, error)
}

type salesRepo struct {
	Db sql.Db
}

func NewSalesRepo(db sql.Db) SalesRepo {
	return &salesRepo{
		Db: db,
	}
}

func (s salesRepo) RevenueBetween(networkId string, from, to time.Time) (float64, error) {
	var revenue float64

	q := s.Db.GetGormDb().Model(&PaymentEvent{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("status = ?", "success").
		Where("paid_at >= ? AND paid_at < ?", from, to)

	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}

	result := q.Scan(&revenue)
	if result.Error != nil {
		return 0, result.Error
	}

	return revenue, nil
}

func (s salesRepo) PurchasesBetween(networkId string, from, to time.Time) (uint32, error) {
	var count int64

	q := s.Db.GetGormDb().Model(&PaymentEvent{}).
		Where("status = ?", "success").
		Where("paid_at >= ? AND paid_at < ?", from, to)

	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}

	result := q.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return uint32(count), nil
}

func (s salesRepo) PaidCustomersBetween(networkId string, from, to time.Time) (uint32, error) {
	var count int64

	q := s.Db.GetGormDb().Model(&PaymentEvent{}).
		Distinct("customer_id").
		Where("status = ?", "success").
		Where("paid_at >= ? AND paid_at < ?", from, to)

	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}

	result := q.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return uint32(count), nil
}

func (s salesRepo) RevenueTrendDaily(networkId string, from, to time.Time) ([]DayValue, error) {
	var trend []DayValue

	q := s.Db.GetGormDb().Model(&BusinessSalesRollupDaily{}).
		Select("day AS day, COALESCE(SUM(revenue), 0) AS value").
		Where("day >= ? AND day < ?", from, to).
		Group("day").
		Order("day ASC")

	if networkId != "" {
		q = q.Where("network_id = ?", networkId)
	}

	result := q.Scan(&trend)
	if result.Error != nil {
		return nil, result.Error
	}

	return trend, nil
}

func (s salesRepo) RevenueBySite(networkId string, from, to time.Time) ([]NamedAmount, error) {
	var rows []NamedAmount

	q := s.Db.GetGormDb().Model(&BusinessSalesRollupDaily{}).
		Select("analytics_business_sales_rollup_daily.site_id AS id, " +
			"COALESCE(analytics_site_snapshots.name, '') AS name, " +
			"COALESCE(SUM(analytics_business_sales_rollup_daily.revenue), 0) AS value").
		Joins("LEFT JOIN analytics_site_snapshots ON " +
			"analytics_site_snapshots.site_id = analytics_business_sales_rollup_daily.site_id").
		Where("analytics_business_sales_rollup_daily.day >= ? AND "+
			"analytics_business_sales_rollup_daily.day < ?", from, to).
		Group("analytics_business_sales_rollup_daily.site_id, analytics_site_snapshots.name").
		Order("value DESC")

	if networkId != "" {
		q = q.Where("analytics_business_sales_rollup_daily.network_id = ?", networkId)
	}

	result := q.Scan(&rows)
	if result.Error != nil {
		return nil, result.Error
	}

	return rows, nil
}

func (s salesRepo) RevenueByPackage(networkId string, from, to time.Time) ([]NamedAmount, error) {
	var rows []NamedAmount

	// Package rollups are org-level (no network dimension); networkId is
	// accepted for interface symmetry and currently ignored.
	_ = networkId

	q := s.Db.GetGormDb().Model(&BusinessPackageRollupDaily{}).
		Select("analytics_business_package_rollup_daily.package_id AS id, " +
			"COALESCE(analytics_package_snapshots.name, '') AS name, " +
			"COALESCE(SUM(analytics_business_package_rollup_daily.revenue), 0) AS value").
		Joins("LEFT JOIN analytics_package_snapshots ON " +
			"analytics_package_snapshots.package_id = analytics_business_package_rollup_daily.package_id").
		Where("analytics_business_package_rollup_daily.day >= ? AND "+
			"analytics_business_package_rollup_daily.day < ?", from, to).
		Group("analytics_business_package_rollup_daily.package_id, analytics_package_snapshots.name").
		Order("value DESC")

	result := q.Scan(&rows)
	if result.Error != nil {
		return nil, result.Error
	}

	return rows, nil
}
