/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/uuid"

	net_db "github.com/ukama/ukama/systems/analytics/network/pkg/db"
)

var nFrom = time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
var nTo = time.Date(2026, time.June, 8, 0, 0, 0, 0, time.UTC)

func rows(cols ...string) *sqlmock.Rows { return sqlmock.NewRows(cols) }

/* ---------- HealthRepo ---------- */

func TestHealthRepo_NetworkHealthLatest(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewHealthRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`network_health_rollup`).WillReturnRows(rows("hour").AddRow(time.Now()))

	out, err := r.NetworkHealthLatest("net-1")
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealthRepo_NetworkHealthSeries(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewHealthRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`network_health_rollup`).WillReturnRows(rows("hour").AddRow(time.Now()))

	out, err := r.NetworkHealthSeries("net-1", nFrom, nTo)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/* ---------- EventRepo ---------- */

func TestEventRepo_Recent(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewEventRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`event_logs`).WillReturnRows(rows("count").AddRow(1))
	mock.ExpectQuery(`event_logs`).WillReturnRows(rows("id", "routing_key", "occurred_at").AddRow(1, "event.x", time.Now()))

	out, total, err := r.Recent("net-1", "", "", nFrom, nTo, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, int64(1), total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/* ---------- AlarmRepo ---------- */

func TestAlarmRepo_List(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewAlarmRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`alarm_events`).WillReturnRows(rows("count").AddRow(1))
	mock.ExpectQuery(`alarm_events`).WillReturnRows(rows("alarm_id", "severity", "state").AddRow("a1", "critical", "open"))

	out, total, err := r.List(net_db.AlarmFilter{NetworkId: "net-1", Page: 1, PageSize: 10})
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, int64(1), total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAlarmRepo_Counts(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewAlarmRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`alarm_events`).WillReturnRows(rows("severity", "cnt").AddRow("critical", 1).AddRow("warning", 2))

	out, err := r.Counts("net-1", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), out.Open)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAlarmRepo_ForResource(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewAlarmRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`alarm_events`).WillReturnRows(rows("alarm_id", "severity").AddRow("a1", "warning"))

	out, err := r.ForResource("radio", "", 20)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAlarmRepo_OpenImpact(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewAlarmRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`alarm_events`).WillReturnRows(rows("customers", "revenue").AddRow(4, 120.0))

	customers, revenue, err := r.OpenImpact("net-1")
	assert.NoError(t, err)
	assert.Equal(t, int64(4), customers)
	assert.Equal(t, 120.0, revenue)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/* ---------- NodeRepo ---------- */

func TestNodeRepo_List(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewNodeRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`node_snapshots`).WillReturnRows(rows("count").AddRow(1))
	mock.ExpectQuery(`node_snapshots`).WillReturnRows(rows("node_id", "name", "status").AddRow("n1", "Node", "online"))

	out, total, err := r.List("net-1", "", "", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, int64(1), total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNodeRepo_ListAll(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewNodeRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`node_snapshots`).WillReturnRows(rows("node_id", "name").AddRow("n1", "Node"))

	out, err := r.ListAll("net-1")
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNodeRepo_StatusCounts(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewNodeRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`node_snapshots`).WillReturnRows(rows("status", "cnt").AddRow("online", 2).AddRow("offline", 1))

	out, err := r.StatusCounts("net-1", "")
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNodeRepo_Get(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewNodeRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`node_snapshots`).WillReturnRows(rows("node_id", "name").AddRow("n1", "Node"))

	out, err := r.Get("n1")
	assert.NoError(t, err)
	assert.Equal(t, "n1", out.NodeId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNodeRepo_UptimeBetween(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewNodeRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`node_state_interval`).WillReturnRows(rows("coalesce").AddRow(3600.0))

	out, err := r.UptimeBetween("n1", nFrom, nTo)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, out, 0.0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNodeRepo_PoolCounts(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewNodeRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`inventory_snapshots`).WillReturnRows(rows("state", "cnt").AddRow("available", 2).AddRow("deployed", 5))
	mock.ExpectQuery(`node_snapshots`).WillReturnRows(rows("count").AddRow(1))

	out, err := r.PoolCounts("net-1")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), out.AvailableToInstall)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNodeRepo_ConfiguringDuration(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewNodeRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`node_state_interval`).WillReturnRows(rows("id", "start_at").AddRow(1, time.Now().Add(-time.Hour)))

	_, err := r.ConfiguringDuration("n1")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNodeRepo_Search(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewNodeRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`node_snapshots`).WillReturnRows(rows("node_id", "name").AddRow("n1", "Node"))

	out, err := r.Search("node", "net-1", 10)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/* ---------- MetricRepo ---------- */

func TestMetricRepo_Rollups(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewMetricRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`metric_rollup`).WillReturnRows(rows("hour", "avg").AddRow(time.Now(), 1.0))

	out, err := r.Rollups("x", "", "res-1", nFrom, nTo)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMetricRepo_LatestSamples(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewMetricRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`metric_samples`).WillReturnRows(rows("metric", "value").AddRow("x", 1.0))

	out, err := r.LatestSamples("res-1", 10)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMetricRepo_MetricNames(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewMetricRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`metric_samples`).WillReturnRows(rows("metric", "unit", "last_sample_at").AddRow("x", "u", time.Now()))

	out, err := r.MetricNames("res-1")
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMetricRepo_RadioRollups(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewMetricRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`radio_rollup`).WillReturnRows(rows("hour", "active_ues").AddRow(time.Now(), 10))

	out, err := r.RadioRollups("n1", nFrom, nTo)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMetricRepo_LatestRadioRollup(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewMetricRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`radio_rollup`).WillReturnRows(rows("hour", "active_ues").AddRow(time.Now(), 10))

	out, err := r.LatestRadioRollup("n1")
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMetricRepo_RadioRollupSums(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewMetricRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`radio_rollup`).WillReturnRows(rows("active_ues", "attach_failures").AddRow(10, 1))

	ues, fails, err := r.RadioRollupSums("net-1", time.Now())
	assert.NoError(t, err)
	assert.Equal(t, int64(10), ues)
	assert.Equal(t, int64(1), fails)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMetricRepo_BackhaulRollups(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewMetricRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`backhaul_rollup`).WillReturnRows(rows("hour", "latency_ms").AddRow(time.Now(), 20.0))

	out, err := r.BackhaulRollups("s1", nFrom, nTo)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMetricRepo_LatestBackhaulRollup(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewMetricRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`backhaul_rollup`).WillReturnRows(rows("hour", "latency_ms").AddRow(time.Now(), 20.0))

	out, err := r.LatestBackhaulRollup("s1")
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMetricRepo_PowerRollups(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewMetricRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`power_rollup`).WillReturnRows(rows("hour", "battery_percent").AddRow(time.Now(), 80.0))

	out, err := r.PowerRollups("s1", nFrom, nTo)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMetricRepo_LatestPowerRollup(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewMetricRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`power_rollup`).WillReturnRows(rows("hour", "battery_percent").AddRow(time.Now(), 80.0))

	out, err := r.LatestPowerRollup("s1")
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/* ---------- SiteRepo (remaining) ---------- */

func TestSiteRepo_List(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewSiteRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`site_snapshots`).WillReturnRows(rows("count").AddRow(1))
	mock.ExpectQuery(`site_snapshots`).WillReturnRows(rows("site_id", "name", "status").AddRow(uuid.NewV4(), "Site", "online"))

	out, total, err := r.List("net-1", "", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, int64(1), total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepo_StatusCounts(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewSiteRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`site_snapshots`).WillReturnRows(rows("status", "cnt").AddRow("online", 4).AddRow("offline", 1))

	out, err := r.StatusCounts("net-1")
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepo_UptimeBetween(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewSiteRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`site_state_interval`).WillReturnRows(rows("seconds", "episodes").AddRow(3600.0, 1))
	mock.ExpectQuery(`site_state_interval`).WillReturnRows(rows("count").AddRow(1))

	out, err := r.UptimeBetween("s1", nFrom, nTo)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, out, 0.0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepo_LatestSiteHealth(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewSiteRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`site_health_rollup`).WillReturnRows(rows("hour", "uptime_percent").AddRow(time.Now(), 99.0))

	out, err := r.LatestSiteHealth("s1")
	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepo_Search(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewSiteRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`site_snapshots`).WillReturnRows(rows("site_id", "name").AddRow(uuid.NewV4(), "Site"))

	out, err := r.Search("site", "net-1", 10)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepo_OfflineDuration(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewSiteRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`site_state_interval`).WillReturnRows(rows("id", "start_at").AddRow(1, time.Now().Add(-time.Hour)))

	_, err := r.OfflineDuration("s1")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepo_SiteHealthSeries(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := net_db.NewSiteRepo(&UkamaDbMock{GormDb: gdb})
	mock.ExpectQuery(`site_health_rollup`).WillReturnRows(rows("hour", "uptime_percent").AddRow(time.Now(), 99.0))

	out, err := r.SiteHealthSeries("s1", nFrom, nTo)
	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}
