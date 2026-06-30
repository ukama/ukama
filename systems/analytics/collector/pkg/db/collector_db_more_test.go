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

	col_db "github.com/ukama/ukama/systems/analytics/collector/pkg/db"
)

func dbMock(t *testing.T) (sqlmock.Sqlmock, *UkamaDbMock) {
	mock, gdb := setupMockDB(t)
	return mock, &UkamaDbMock{GormDb: gdb}
}

// execWrite expects a transactional non-returning write (string-PK upsert,
// update or delete).
func execWrite(mock sqlmock.Sqlmock, re string) {
	mock.ExpectBegin()
	mock.ExpectExec(re).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

// queryWrite expects a transactional INSERT ... RETURNING id (auto-increment).
func queryWrite(mock sqlmock.Sqlmock, re string) {
	mock.ExpectBegin()
	mock.ExpectQuery(re).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
}

/* ---------- StateRepo ---------- */

func TestStateRepo(t *testing.T) {
	t.Run("UpsertRefreshState", func(t *testing.T) {
		mock, udb := dbMock(t)
		execWrite(mock, `INSERT INTO`)
		assert.NoError(t, col_db.NewStateRepo(udb).UpsertRefreshState(&col_db.RefreshState{Source: "registry"}))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("GetRefreshStates", func(t *testing.T) {
		mock, udb := dbMock(t)
		mock.ExpectQuery(`refresh_states`).WillReturnRows(sqlmock.NewRows([]string{"source"}).AddRow("registry"))
		out, err := col_db.NewStateRepo(udb).GetRefreshStates()
		assert.NoError(t, err)
		assert.Len(t, out, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("MarkRollupDirty", func(t *testing.T) {
		mock, udb := dbMock(t)
		execWrite(mock, `INSERT INTO`)
		assert.NoError(t, col_db.NewStateRepo(udb).MarkRollupDirty("business_sales_daily"))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("SetRollupWatermark", func(t *testing.T) {
		mock, udb := dbMock(t)
		execWrite(mock, `INSERT INTO`)
		assert.NoError(t, col_db.NewStateRepo(udb).SetRollupWatermark("business_sales_daily", time.Now()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("GetRollupStates", func(t *testing.T) {
		mock, udb := dbMock(t)
		mock.ExpectQuery(`rollup_states`).WillReturnRows(sqlmock.NewRows([]string{"rollup"}).AddRow("business_sales_daily"))
		out, err := col_db.NewStateRepo(udb).GetRollupStates()
		assert.NoError(t, err)
		assert.Len(t, out, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

/* ---------- SnapshotRepo ---------- */

func TestSnapshotRepo_Upserts(t *testing.T) {
	r := func(udb *UkamaDbMock) col_db.SnapshotRepo { return col_db.NewSnapshotRepo(udb) }

	cases := []struct {
		name string
		fn   func(col_db.SnapshotRepo) error
	}{
		{"Network", func(s col_db.SnapshotRepo) error { return s.UpsertNetwork(&col_db.NetworkSnapshot{NetworkId: uuid.NewV4()}) }},
		{"Site", func(s col_db.SnapshotRepo) error { return s.UpsertSite(&col_db.SiteSnapshot{SiteId: uuid.NewV4()}) }},
		{"Node", func(s col_db.SnapshotRepo) error { return s.UpsertNode(&col_db.NodeSnapshot{NodeId: "n1"}) }},
		{"Customer", func(s col_db.SnapshotRepo) error { return s.UpsertCustomer(&col_db.CustomerSnapshot{CustomerId: uuid.NewV4()}) }},
		{"Sim", func(s col_db.SnapshotRepo) error { return s.UpsertSim(&col_db.SimSnapshot{SimId: "s1"}) }},
		{"SimBatch", func(s col_db.SnapshotRepo) error { return s.UpsertSimBatch(&col_db.SimBatchSnapshot{BatchId: "b1"}) }},
		{"Package", func(s col_db.SnapshotRepo) error { return s.UpsertPackage(&col_db.PackageSnapshot{PackageId: uuid.NewV4()}) }},
		{"Inventory", func(s col_db.SnapshotRepo) error { return s.UpsertInventory(&col_db.InventorySnapshot{}) }},
		{"HealthReport", func(s col_db.SnapshotRepo) error { return s.UpsertHealthReport(&col_db.HealthReportSnapshot{}) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock, udb := dbMock(t)
			execWrite(mock, `INSERT INTO`)
			assert.NoError(t, tc.fn(r(udb)))
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	// BillingSnapshot has a numeric Id PK, so gorm appends RETURNING.
	t.Run("Billing", func(t *testing.T) {
		mock, udb := dbMock(t)
		queryWrite(mock, `INSERT INTO`)
		assert.NoError(t, r(udb).UpsertBilling(&col_db.BillingSnapshot{Id: 1}))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSnapshotRepo_UpdatesAndDeletes(t *testing.T) {
	t.Run("UpdateNodeStatus", func(t *testing.T) {
		mock, udb := dbMock(t)
		execWrite(mock, `UPDATE`)
		assert.NoError(t, col_db.NewSnapshotRepo(udb).UpdateNodeStatus("n1", "online", time.Now()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("DeleteCustomer", func(t *testing.T) {
		mock, udb := dbMock(t)
		execWrite(mock, `customer_snapshots`)
		assert.NoError(t, col_db.NewSnapshotRepo(udb).DeleteCustomer(uuid.NewV4().String()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("DeleteSim", func(t *testing.T) {
		mock, udb := dbMock(t)
		execWrite(mock, `sim_snapshots`)
		assert.NoError(t, col_db.NewSnapshotRepo(udb).DeleteSim("s1"))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("DeletePackage", func(t *testing.T) {
		mock, udb := dbMock(t)
		execWrite(mock, `package_snapshots`)
		assert.NoError(t, col_db.NewSnapshotRepo(udb).DeletePackage(uuid.NewV4().String()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

/* ---------- FactRepo ---------- */

func TestFactRepo_Adds(t *testing.T) {
	cases := []struct {
		name string
		fn   func(col_db.FactRepo) error
	}{
		{"Payment", func(f col_db.FactRepo) error { return f.AddPaymentEvent(&col_db.PaymentEvent{ExternalId: "p1"}) }},
		{"Usage", func(f col_db.FactRepo) error { return f.AddUsageEvent(&col_db.UsageEvent{}) }},
		{"Metric", func(f col_db.FactRepo) error { return f.AddMetricSample(&col_db.MetricSample{}) }},
		{"Alarm", func(f col_db.FactRepo) error { return f.AddAlarmEvent(&col_db.AlarmEvent{AlarmId: "a1"}) }},
		{"NodeState", func(f col_db.FactRepo) error { return f.AddNodeStateEvent(&col_db.NodeStateEvent{}) }},
		{"SiteState", func(f col_db.FactRepo) error { return f.AddSiteStateEvent(&col_db.SiteStateEvent{}) }},
		{"Customer", func(f col_db.FactRepo) error { return f.AddCustomerEvent(&col_db.CustomerEvent{}) }},
		{"Sim", func(f col_db.FactRepo) error { return f.AddSimEvent(&col_db.SimEvent{}) }},
		{"Package", func(f col_db.FactRepo) error { return f.AddPackageEvent(&col_db.PackageEvent{}) }},
		{"Inventory", func(f col_db.FactRepo) error { return f.AddInventoryEvent(&col_db.InventoryEvent{}) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock, udb := dbMock(t)
			queryWrite(mock, `INSERT INTO`)
			assert.NoError(t, tc.fn(col_db.NewFactRepo(udb)))
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFactRepo_Transitions(t *testing.T) {
	t.Run("TransitionNodeState", func(t *testing.T) {
		mock, udb := dbMock(t)
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE`).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery(`INSERT INTO`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
		assert.NoError(t, col_db.NewFactRepo(udb).TransitionNodeState("n1", "online", time.Now()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("TransitionSiteState", func(t *testing.T) {
		mock, udb := dbMock(t)
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE`).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery(`INSERT INTO`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
		assert.NoError(t, col_db.NewFactRepo(udb).TransitionSiteState(uuid.NewV4(), "online", time.Now()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("TransitionSimState", func(t *testing.T) {
		mock, udb := dbMock(t)
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE`).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery(`INSERT INTO`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
		assert.NoError(t, col_db.NewFactRepo(udb).TransitionSimState("sim-1", "active", time.Now()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("OpenCustomerPackageInterval", func(t *testing.T) {
		mock, udb := dbMock(t)
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE`).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery(`INSERT INTO`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
		assert.NoError(t, col_db.NewFactRepo(udb).OpenCustomerPackageInterval(uuid.NewV4(), uuid.NewV4(), "active", time.Now()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("CloseCustomerPackageInterval", func(t *testing.T) {
		mock, udb := dbMock(t)
		execWrite(mock, `UPDATE`)
		assert.NoError(t, col_db.NewFactRepo(udb).CloseCustomerPackageInterval(uuid.NewV4(), time.Now()))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

/* ---------- WithGorm constructors + InTransaction ---------- */

type txRunner interface {
	InTransaction(func(col_db.EventRepo, col_db.StateRepo, col_db.SnapshotRepo, col_db.FactRepo) error) error
}

func TestWithGormConstructors(t *testing.T) {
	_, gdb := setupMockDB(t)

	assert.NotNil(t, col_db.NewEventRepoWithGorm(gdb))
	assert.NotNil(t, col_db.NewFactRepoWithGorm(gdb))
	assert.NotNil(t, col_db.NewSnapshotRepoWithGorm(gdb))
	assert.NotNil(t, col_db.NewStateRepoWithGorm(gdb))
	assert.NotNil(t, col_db.NewRollupRepoWithGorm(gdb))
}

func TestEventRepo_InTransaction(t *testing.T) {
	mock, gdb := setupMockDB(t)
	repo := col_db.NewEventRepoWithGorm(gdb)

	mock.ExpectBegin()
	mock.ExpectCommit()

	runner, ok := repo.(txRunner)
	assert.True(t, ok)

	err := runner.InTransaction(func(col_db.EventRepo, col_db.StateRepo, col_db.SnapshotRepo, col_db.FactRepo) error {
		return nil
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

/* ---------- RollupRepo Upserts + Rebuilds ---------- */

func TestRollupRepo_Upserts(t *testing.T) {
	cases := []struct {
		name string
		fn   func(col_db.RollupRepo) error
	}{
		{"BusinessPackageDaily", func(r col_db.RollupRepo) error { return r.UpsertBusinessPackageDaily(&col_db.BusinessPackageRollupDaily{}) }},
		{"BusinessSiteDaily", func(r col_db.RollupRepo) error { return r.UpsertBusinessSiteDaily(&col_db.BusinessSiteRollupDaily{}) }},
		{"BusinessInventoryDaily", func(r col_db.RollupRepo) error { return r.UpsertBusinessInventoryDaily(&col_db.BusinessInventoryRollupDaily{}) }},
		{"BusinessBillingDaily", func(r col_db.RollupRepo) error { return r.UpsertBusinessBillingDaily(&col_db.BusinessBillingRollupDaily{}) }},
		{"CustomerUsageDaily", func(r col_db.RollupRepo) error { return r.UpsertCustomerUsageDaily(&col_db.CustomerUsageRollupDaily{}) }},
		{"CustomerStateDaily", func(r col_db.RollupRepo) error { return r.UpsertCustomerStateDaily(&col_db.CustomerStateRollupDaily{}) }},
		{"NetworkHealthHourly", func(r col_db.RollupRepo) error { return r.UpsertNetworkHealthHourly(&col_db.NetworkHealthRollupHourly{}) }},
		{"SiteHealthHourly", func(r col_db.RollupRepo) error { return r.UpsertSiteHealthHourly(&col_db.SiteHealthRollupHourly{}) }},
		{"NodeHealthHourly", func(r col_db.RollupRepo) error { return r.UpsertNodeHealthHourly(&col_db.NodeHealthRollupHourly{}) }},
		{"MetricHourly", func(r col_db.RollupRepo) error { return r.UpsertMetricHourly(&col_db.MetricRollupHourly{}) }},
		{"AlarmDaily", func(r col_db.RollupRepo) error { return r.UpsertAlarmDaily(&col_db.AlarmRollupDaily{}) }},
		{"RadioHourly", func(r col_db.RollupRepo) error { return r.UpsertRadioHourly(&col_db.RadioRollupHourly{}) }},
		{"BackhaulHourly", func(r col_db.RollupRepo) error { return r.UpsertBackhaulHourly(&col_db.BackhaulRollupHourly{}) }},
		{"PowerHourly", func(r col_db.RollupRepo) error { return r.UpsertPowerHourly(&col_db.PowerRollupHourly{}) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock, udb := dbMock(t)
			queryWrite(mock, `INSERT INTO`)
			assert.NoError(t, tc.fn(col_db.NewRollupRepo(udb)))
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRollupRepo_Rebuilds(t *testing.T) {
	now := time.Now()
	from := now.AddDate(0, 0, -7)

	cases := []struct {
		name string
		fn   func(col_db.RollupRepo) error
	}{
		{"Package", func(r col_db.RollupRepo) error { return r.RebuildPackageDaily(from, now) }},
		{"Billing", func(r col_db.RollupRepo) error { return r.RebuildBillingDaily(from, now) }},
		{"CustomerUsage", func(r col_db.RollupRepo) error { return r.RebuildCustomerUsageDaily(from, now) }},
		{"CustomerState", func(r col_db.RollupRepo) error { return r.RebuildCustomerStateDaily(from, now) }},
		{"Alarm", func(r col_db.RollupRepo) error { return r.RebuildAlarmDaily(from, now) }},
		{"Metric", func(r col_db.RollupRepo) error { return r.RebuildMetricHourly(from, now) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock, udb := dbMock(t)
			mock.ExpectExec(`INSERT INTO`).WillReturnResult(sqlmock.NewResult(0, 3))
			assert.NoError(t, tc.fn(col_db.NewRollupRepo(udb)))
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
