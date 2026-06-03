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

	"gorm.io/gorm/clause"
)

type RollupRepo interface {
	UpsertBusinessSalesDaily(r *BusinessSalesRollupDaily) error
	UpsertBusinessPackageDaily(r *BusinessPackageRollupDaily) error
	UpsertBusinessSiteDaily(r *BusinessSiteRollupDaily) error
	UpsertBusinessInventoryDaily(r *BusinessInventoryRollupDaily) error
	UpsertBusinessBillingDaily(r *BusinessBillingRollupDaily) error
	UpsertCustomerUsageDaily(r *CustomerUsageRollupDaily) error
	UpsertCustomerStateDaily(r *CustomerStateRollupDaily) error
	UpsertNetworkHealthHourly(r *NetworkHealthRollupHourly) error
	UpsertSiteHealthHourly(r *SiteHealthRollupHourly) error
	UpsertNodeHealthHourly(r *NodeHealthRollupHourly) error
	UpsertMetricHourly(r *MetricRollupHourly) error
	UpsertAlarmDaily(r *AlarmRollupDaily) error
	UpsertRadioHourly(r *RadioRollupHourly) error
	UpsertBackhaulHourly(r *BackhaulRollupHourly) error
	UpsertPowerHourly(r *PowerRollupHourly) error

	/* Rebuilds from facts via SQL aggregates. */
	RebuildSalesDaily(from, to time.Time) error
	RebuildPackageDaily(from, to time.Time) error
	RebuildBillingDaily(from, to time.Time) error
	RebuildCustomerUsageDaily(from, to time.Time) error
	RebuildCustomerStateDaily(from, to time.Time) error
	RebuildAlarmDaily(from, to time.Time) error
	RebuildMetricHourly(from, to time.Time) error
}

type rollupRepo struct {
	Db sql.Db
}

func NewRollupRepo(db sql.Db) RollupRepo {
	return &rollupRepo{
		Db: db,
	}
}

func (r *rollupRepo) upsert(columns []string, value interface{}) error {
	cols := make([]clause.Column, 0, len(columns))
	for _, c := range columns {
		cols = append(cols, clause.Column{Name: c})
	}

	result := r.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   cols,
		UpdateAll: true,
	}).Create(value)

	return result.Error
}

func (r *rollupRepo) UpsertBusinessSalesDaily(v *BusinessSalesRollupDaily) error {
	return r.upsert([]string{"day", "network_id", "site_id"}, v)
}

func (r *rollupRepo) UpsertBusinessPackageDaily(v *BusinessPackageRollupDaily) error {
	return r.upsert([]string{"day", "package_id"}, v)
}

func (r *rollupRepo) UpsertBusinessSiteDaily(v *BusinessSiteRollupDaily) error {
	return r.upsert([]string{"day", "site_id"}, v)
}

func (r *rollupRepo) UpsertBusinessInventoryDaily(v *BusinessInventoryRollupDaily) error {
	return r.upsert([]string{"day"}, v)
}

func (r *rollupRepo) UpsertBusinessBillingDaily(v *BusinessBillingRollupDaily) error {
	return r.upsert([]string{"day"}, v)
}

func (r *rollupRepo) UpsertCustomerUsageDaily(v *CustomerUsageRollupDaily) error {
	return r.upsert([]string{"day", "customer_id"}, v)
}

func (r *rollupRepo) UpsertCustomerStateDaily(v *CustomerStateRollupDaily) error {
	return r.upsert([]string{"day", "network_id"}, v)
}

func (r *rollupRepo) UpsertNetworkHealthHourly(v *NetworkHealthRollupHourly) error {
	return r.upsert([]string{"hour", "network_id"}, v)
}

func (r *rollupRepo) UpsertSiteHealthHourly(v *SiteHealthRollupHourly) error {
	return r.upsert([]string{"hour", "site_id"}, v)
}

func (r *rollupRepo) UpsertNodeHealthHourly(v *NodeHealthRollupHourly) error {
	return r.upsert([]string{"hour", "node_id"}, v)
}

func (r *rollupRepo) UpsertMetricHourly(v *MetricRollupHourly) error {
	return r.upsert([]string{"hour", "metric", "resource_id"}, v)
}

func (r *rollupRepo) UpsertAlarmDaily(v *AlarmRollupDaily) error {
	return r.upsert([]string{"day"}, v)
}

func (r *rollupRepo) UpsertRadioHourly(v *RadioRollupHourly) error {
	return r.upsert([]string{"hour", "node_id"}, v)
}

func (r *rollupRepo) UpsertBackhaulHourly(v *BackhaulRollupHourly) error {
	return r.upsert([]string{"hour", "site_id"}, v)
}

func (r *rollupRepo) UpsertPowerHourly(v *PowerRollupHourly) error {
	return r.upsert([]string{"hour", "site_id"}, v)
}

// RebuildSalesDaily recomputes the business sales rollup from payment events
// for the given window.
func (r *rollupRepo) RebuildSalesDaily(from, to time.Time) error {
	return r.Db.GetGormDb().Exec(`
		INSERT INTO analytics_business_sales_rollup_daily
			(day, network_id, site_id, revenue, purchases, paid_customers, data_sold_mb)
		SELECT date_trunc('day', paid_at) AS day,
			network_id,
			site_id,
			COALESCE(SUM(amount), 0) AS revenue,
			COUNT(*) AS purchases,
			COUNT(DISTINCT customer_id) AS paid_customers,
			0 AS data_sold_mb
		FROM analytics_payment_events
		WHERE status = 'success' AND paid_at >= ? AND paid_at < ?
		GROUP BY 1, 2, 3
		ON CONFLICT (day, network_id, site_id) DO UPDATE SET
			revenue = EXCLUDED.revenue,
			purchases = EXCLUDED.purchases,
			paid_customers = EXCLUDED.paid_customers`,
		from, to).Error
}

// RebuildPackageDaily recomputes the business package rollup from payment
// events for the given window.
func (r *rollupRepo) RebuildPackageDaily(from, to time.Time) error {
	return r.Db.GetGormDb().Exec(`
		INSERT INTO analytics_business_package_rollup_daily
			(day, package_id, sold_count, revenue, data_used_mb)
		SELECT date_trunc('day', paid_at) AS day,
			package_id,
			COUNT(*) AS sold_count,
			COALESCE(SUM(amount), 0) AS revenue,
			0 AS data_used_mb
		FROM analytics_payment_events
		WHERE status = 'success' AND paid_at >= ? AND paid_at < ?
		GROUP BY 1, 2
		ON CONFLICT (day, package_id) DO UPDATE SET
			sold_count = EXCLUDED.sold_count,
			revenue = EXCLUDED.revenue`,
		from, to).Error
}

// RebuildBillingDaily recomputes the billing rollup from successful payment
// events for the given window.
func (r *rollupRepo) RebuildBillingDaily(from, to time.Time) error {
	return r.Db.GetGormDb().Exec(`
		INSERT INTO analytics_business_billing_rollup_daily
			(day, invoiced_amount, invoice_count)
		SELECT date_trunc('day', paid_at) AS day,
			COALESCE(SUM(amount), 0) AS invoiced_amount,
			COUNT(*) AS invoice_count
		FROM analytics_payment_events
		WHERE status = 'success' AND paid_at >= ? AND paid_at < ?
		GROUP BY 1
		ON CONFLICT (day) DO UPDATE SET
			invoiced_amount = EXCLUDED.invoiced_amount,
			invoice_count = EXCLUDED.invoice_count`,
		from, to).Error
}

// RebuildCustomerUsageDaily recomputes the customer usage rollup from usage
// events for the given window.
func (r *rollupRepo) RebuildCustomerUsageDaily(from, to time.Time) error {
	return r.Db.GetGormDb().Exec(`
		INSERT INTO analytics_customer_usage_rollup_daily
			(day, customer_id, data_used_mb)
		SELECT date_trunc('day', start_at) AS day,
			customer_id,
			COALESCE(SUM(bytes_used), 0) / 1048576.0 AS data_used_mb
		FROM analytics_usage_events
		WHERE start_at >= ? AND start_at < ?
		GROUP BY 1, 2
		ON CONFLICT (day, customer_id) DO UPDATE SET
			data_used_mb = EXCLUDED.data_used_mb`,
		from, to).Error
}

// RebuildCustomerStateDaily recomputes per-network customer counters from
// customer events and customer snapshots for the given window.
func (r *rollupRepo) RebuildCustomerStateDaily(from, to time.Time) error {
	return r.Db.GetGormDb().Exec(`
		INSERT INTO analytics_customer_state_rollup_daily
			(day, network_id, total, active, new, expired, failed_activations)
		SELECT d.day,
			s.network_id,
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE s.status = 'active') AS active,
			COUNT(*) FILTER (WHERE s.source_created_at >= d.day
				AND s.source_created_at < d.day + interval '1 day') AS new,
			COUNT(*) FILTER (WHERE s.status = 'expired') AS expired,
			(SELECT COUNT(*) FROM analytics_customer_events ce
				WHERE ce.kind = 'activation_failed'
				AND ce.occurred_at >= d.day
				AND ce.occurred_at < d.day + interval '1 day') AS failed_activations
		FROM generate_series(date_trunc('day', ?::timestamptz),
			date_trunc('day', ?::timestamptz), interval '1 day') AS d(day)
		CROSS JOIN analytics_customer_snapshots s
		GROUP BY 1, 2
		ON CONFLICT (day, network_id) DO UPDATE SET
			total = EXCLUDED.total,
			active = EXCLUDED.active,
			new = EXCLUDED.new,
			expired = EXCLUDED.expired,
			failed_activations = EXCLUDED.failed_activations`,
		from, to).Error
}

// RebuildAlarmDaily recomputes the alarm rollup from alarm events for the
// given window.
func (r *rollupRepo) RebuildAlarmDaily(from, to time.Time) error {
	return r.Db.GetGormDb().Exec(`
		INSERT INTO analytics_alarm_rollup_daily
			(day, opened, closed, critical, warning)
		SELECT date_trunc('day', opened_at) AS day,
			COUNT(*) AS opened,
			COUNT(closed_at) AS closed,
			COUNT(*) FILTER (WHERE severity = 'critical') AS critical,
			COUNT(*) FILTER (WHERE severity = 'warning') AS warning
		FROM analytics_alarm_events
		WHERE opened_at >= ? AND opened_at < ?
		GROUP BY 1
		ON CONFLICT (day) DO UPDATE SET
			opened = EXCLUDED.opened,
			closed = EXCLUDED.closed,
			critical = EXCLUDED.critical,
			warning = EXCLUDED.warning`,
		from, to).Error
}

// RebuildMetricHourly recomputes the hourly metric rollup from raw metric
// samples for the given window.
func (r *rollupRepo) RebuildMetricHourly(from, to time.Time) error {
	return r.Db.GetGormDb().Exec(`
		INSERT INTO analytics_metric_rollup_hourly
			(hour, metric, resource_type, resource_id, avg, min, max, count)
		SELECT date_trunc('hour', sampled_at) AS hour,
			metric,
			MAX(resource_type) AS resource_type,
			resource_id,
			AVG(value) AS avg,
			MIN(value) AS min,
			MAX(value) AS max,
			COUNT(*) AS count
		FROM analytics_metric_samples
		WHERE sampled_at >= ? AND sampled_at < ?
		GROUP BY 1, 2, 4
		ON CONFLICT (hour, metric, resource_id) DO UPDATE SET
			resource_type = EXCLUDED.resource_type,
			avg = EXCLUDED.avg,
			min = EXCLUDED.min,
			max = EXCLUDED.max,
			count = EXCLUDED.count`,
		from, to).Error
}
