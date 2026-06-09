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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"

	biz_db "github.com/ukama/ukama/systems/analytics/business/pkg/db"
)

func newDbMock(t *testing.T) (sqlmock.Sqlmock, *UkamaDbMock) {
	sdb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "sqlmock_db_0", DriverName: "postgres", Conn: sdb, PreferSimpleProtocol: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return mock, &UkamaDbMock{GormDb: gdb}
}

var winFrom = time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
var winTo = time.Date(2026, time.June, 8, 0, 0, 0, 0, time.UTC)

func TestActivityRepo_Recent(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewActivityRepo(udb)

	mock.ExpectQuery(`event_logs`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "routing_key", "occurred_at"}).
			AddRow(1, "event.x", time.Now()))

	out, err := repo.Recent(5)

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActivityRepo_RecentDefaultLimit(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewActivityRepo(udb)

	mock.ExpectQuery(`event_logs`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "routing_key", "occurred_at"}))

	out, err := repo.Recent(0)

	assert.NoError(t, err)
	assert.Len(t, out, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBillingRepo_GetBillingSnapshot(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewBillingRepo(udb)

	mock.ExpectQuery(`billing_snapshots`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "balance_cents"}).AddRow(1, 25000))

	snap, err := repo.GetBillingSnapshot()

	assert.NoError(t, err)
	assert.NotNil(t, snap)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBillingRepo_InvoiceRollups(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewBillingRepo(udb)

	mock.ExpectQuery(`business_billing_rollup`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "day", "invoice_count"}).
			AddRow(1, time.Now(), 3))

	out, err := repo.InvoiceRollups(winFrom, winTo)

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInventoryRepo_SimCounts(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewInventoryRepo(udb)

	mock.ExpectQuery(`sim_snapshots`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(100))
	mock.ExpectQuery(`sim_snapshots`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(40))

	available, active, err := repo.SimCounts()

	assert.NoError(t, err)
	assert.Equal(t, uint32(100), available)
	assert.Equal(t, uint32(40), active)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInventoryRepo_NodeCounts(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewInventoryRepo(udb)

	mock.ExpectQuery(`inventory_snapshots`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))
	mock.ExpectQuery(`inventory_snapshots`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(8))

	available, deployed, err := repo.NodeCounts()

	assert.NoError(t, err)
	assert.Equal(t, uint32(20), available)
	assert.Equal(t, uint32(8), deployed)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPackageRepo_ListPackages(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewPackageRepo(udb)

	mock.ExpectQuery(`count`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`package_snapshots`).
		WillReturnRows(sqlmock.NewRows([]string{"package_id", "name"}).AddRow(uuid.NewV4(), "Starter"))

	pkgs, total, err := repo.ListPackages(1, 10)

	assert.NoError(t, err)
	assert.Len(t, pkgs, 1)
	assert.Equal(t, int64(1), total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPackageRepo_PackageRollups(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewPackageRepo(udb)

	mock.ExpectQuery(`business_package_rollup`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "package_id", "sold_count"}).
			AddRow(1, uuid.NewV4(), 5))

	out, err := repo.PackageRollups(winFrom, winTo)

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepo_ListSites(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewSiteRepo(udb)

	mock.ExpectQuery(`count`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`site_snapshots`).
		WillReturnRows(sqlmock.NewRows([]string{"site_id", "name", "status"}).
			AddRow(uuid.NewV4(), "Site One", "online"))

	sites, total, err := repo.ListSites("net-1", 1, 10)

	assert.NoError(t, err)
	assert.Len(t, sites, 1)
	assert.Equal(t, int64(1), total)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepo_GetSite(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewSiteRepo(udb)

	id := uuid.NewV4()
	mock.ExpectQuery(`site_snapshots`).
		WillReturnRows(sqlmock.NewRows([]string{"site_id", "name", "status"}).AddRow(id, "Site One", "online"))

	site, err := repo.GetSite(id.String())

	assert.NoError(t, err)
	assert.Equal(t, id, site.SiteId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepo_SiteRollups(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewSiteRepo(udb)

	mock.ExpectQuery(`business_site_rollup`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "site_id", "customers"}).
			AddRow(1, uuid.NewV4(), 3))

	out, err := repo.SiteRollups("site-1", winFrom, winTo)

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepo_SiteUptime(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewSiteRepo(udb)

	mock.ExpectQuery(`site_health_rollup`).
		WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(99.5))

	uptime, err := repo.SiteUptime("site-1", winFrom, winTo)

	assert.NoError(t, err)
	assert.Equal(t, 99.5, uptime)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesRepo_PurchasesBetween(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewSalesRepo(udb)

	mock.ExpectQuery(`payment_events`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(8))

	out, err := repo.PurchasesBetween("net-1", winFrom, winTo)

	assert.NoError(t, err)
	assert.Equal(t, uint32(8), out)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesRepo_PaidCustomersBetween(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewSalesRepo(udb)

	mock.ExpectQuery(`payment_events`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(4))

	out, err := repo.PaidCustomersBetween("net-1", winFrom, winTo)

	assert.NoError(t, err)
	assert.Equal(t, uint32(4), out)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesRepo_RevenueTrendDaily(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewSalesRepo(udb)

	mock.ExpectQuery(`business_sales_rollup`).
		WillReturnRows(sqlmock.NewRows([]string{"day", "value"}).AddRow(time.Now(), 120.5))

	out, err := repo.RevenueTrendDaily("net-1", winFrom, winTo)

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesRepo_RevenueBySite(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewSalesRepo(udb)

	mock.ExpectQuery(`business_sales_rollup`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value"}).AddRow("site-1", "Site One", 50.0))

	out, err := repo.RevenueBySite("net-1", winFrom, winTo)

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSalesRepo_RevenueByPackage(t *testing.T) {
	mock, udb := newDbMock(t)
	repo := biz_db.NewSalesRepo(udb)

	mock.ExpectQuery(`business_package_rollup`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value"}).AddRow("pkg-1", "Starter", 30.0))

	out, err := repo.RevenueByPackage("net-1", winFrom, winTo)

	assert.NoError(t, err)
	assert.Len(t, out, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}
