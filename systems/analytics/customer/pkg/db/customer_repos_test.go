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

	cust_db "github.com/ukama/ukama/systems/analytics/customer/pkg/db"
)

var cFrom = time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
var cTo = time.Date(2026, time.June, 8, 0, 0, 0, 0, time.UTC)

func TestCustomerRepo_Search(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := cust_db.NewCustomerRepo(&UkamaDbMock{GormDb: gdb})

	mock.ExpectQuery(`customer_snapshots`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`customer_snapshots`).
		WillReturnRows(sqlmock.NewRows([]string{"customer_id", "name"}).AddRow(uuid.NewV4(), "Jane"))

	out, count, err := r.Search("jane", uuid.Nil, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, int64(1), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCustomerRepo_PackageIntervals(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := cust_db.NewCustomerRepo(&UkamaDbMock{GormDb: gdb})

	mock.ExpectQuery(`customer_package_interval`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "customer_id", "state"}).AddRow(1, uuid.NewV4(), "active"))

	out, err := r.PackageIntervals(uuid.NewV4())

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCustomerRepo_UsageBetween(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := cust_db.NewCustomerRepo(&UkamaDbMock{GormDb: gdb})

	mock.ExpectQuery(`customer_usage_rollup`).WillReturnRows(sqlmock.NewRows([]string{"sum"}).AddRow(150.5))

	out, err := r.UsageBetween(uuid.NewV4(), cFrom, cTo)

	assert.NoError(t, err)
	assert.Equal(t, 150.5, out)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCustomerRepo_SiteNames(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := cust_db.NewCustomerRepo(&UkamaDbMock{GormDb: gdb})

	id := uuid.NewV4()
	mock.ExpectQuery(`site_snapshots`).
		WillReturnRows(sqlmock.NewRows([]string{"site_id", "name"}).AddRow(id, "Site One"))

	out, err := r.SiteNames([]uuid.UUID{id})

	assert.NoError(t, err)
	assert.Equal(t, "Site One", out[id])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCustomerRepo_SiteNames_Empty(t *testing.T) {
	_, gdb := setupMockDb(t)
	r := cust_db.NewCustomerRepo(&UkamaDbMock{GormDb: gdb})

	out, err := r.SiteNames(nil)

	assert.NoError(t, err)
	assert.Len(t, out, 0)
}

func TestSimRepo_List(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := cust_db.NewSimRepo(&UkamaDbMock{GormDb: gdb})

	mock.ExpectQuery(`sim_snapshots`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`sim_snapshots`).
		WillReturnRows(sqlmock.NewRows([]string{"sim_id", "iccid", "status"}).AddRow("sim-1", "8910", "active"))

	out, count, err := r.List(uuid.Nil, "active", 1, 10)

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.Equal(t, int64(1), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSimRepo_PoolCounts(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := cust_db.NewSimRepo(&UkamaDbMock{GormDb: gdb})

	mock.ExpectQuery(`sim_snapshots`).
		WillReturnRows(sqlmock.NewRows([]string{"status", "count"}).
			AddRow("available", 10).
			AddRow("active", 150).
			AddRow("assigned", 30))

	total, available, active, assigned, _, _, err := r.PoolCounts()

	assert.NoError(t, err)
	assert.Equal(t, uint32(190), total)
	assert.Equal(t, uint32(10), available)
	assert.Equal(t, uint32(150), active)
	assert.Equal(t, uint32(30), assigned)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSimRepo_Batches(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := cust_db.NewSimRepo(&UkamaDbMock{GormDb: gdb})

	mock.ExpectQuery(`sim_batch_snapshots`).
		WillReturnRows(sqlmock.NewRows([]string{"batch_id", "quantity", "assigned"}).AddRow("b1", 100, 60))

	out, err := r.Batches()

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupportRepo_RecentActivityFor(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := cust_db.NewSupportRepo(&UkamaDbMock{GormDb: gdb})

	mock.ExpectQuery(`event_logs`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "routing_key", "occurred_at"}).AddRow(1, "event.x", time.Now()))

	out, err := r.RecentActivityFor(uuid.NewV4(), 10)

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupportRepo_SiteHealth(t *testing.T) {
	mock, gdb := setupMockDb(t)
	r := cust_db.NewSupportRepo(&UkamaDbMock{GormDb: gdb})

	mock.ExpectQuery(`site_health_rollup`).
		WillReturnRows(sqlmock.NewRows([]string{"site_id", "uptime_percent"}).AddRow(uuid.NewV4(), 99.0))

	out, err := r.SiteHealth(uuid.NewV4())

	assert.NoError(t, err)
	assert.NotNil(t, out)
	assert.NoError(t, mock.ExpectationsWereMet())
}
