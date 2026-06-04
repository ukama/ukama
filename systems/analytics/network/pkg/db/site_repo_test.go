/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package db_test

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"

	net_db "github.com/ukama/ukama/systems/analytics/network/pkg/db"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic("implement me: Init()")
}

func (u UkamaDbMock) Connect() error {
	panic("implement me: Connect()")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	log.Fatal("implement me: ExecuteInTransaction()")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	log.Fatal("implement me: ExecuteInTransaction2()")
	return nil
}

func setupMockDb(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	t.Helper()

	var db *extsql.DB

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return mock, gdb
}

func TestSiteRepo_Get(t *testing.T) {
	t.Run("SiteFound", func(t *testing.T) {
		mock, gdb := setupMockDb(t)

		siteId := uuid.NewV4()
		networkId := uuid.NewV4()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"site_id", "network_id", "name", "status",
			"latitude", "longitude", "node_count", "updated_at"}).
			AddRow(siteId, networkId, "site-a", "online", 1.5, 2.5, 3, now)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "analytics_site_snapshots" WHERE site_id = $1`)).
			WithArgs(siteId.String(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		r := net_db.NewSiteRepo(&UkamaDbMock{GormDb: gdb})

		site, err := r.Get(siteId.String())
		assert.NoError(t, err)
		assert.NotNil(t, site)
		assert.Equal(t, siteId.String(), site.SiteId.String())
		assert.Equal(t, "site-a", site.Name)
		assert.Equal(t, "online", site.Status)
		assert.Equal(t, uint32(3), site.NodeCount)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SiteNotFound", func(t *testing.T) {
		mock, gdb := setupMockDb(t)

		siteId := uuid.NewV4()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "analytics_site_snapshots" WHERE site_id = $1`)).
			WithArgs(siteId.String(), sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		r := net_db.NewSiteRepo(&UkamaDbMock{GormDb: gdb})

		site, err := r.Get(siteId.String())
		assert.Error(t, err)
		assert.Nil(t, site)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSiteRepo_CustomerCount(t *testing.T) {
	mock, gdb := setupMockDb(t)

	siteId := uuid.NewV4()

	mock.ExpectQuery(`SELECT count\(\*\) FROM "analytics_customer_snapshots"`).
		WithArgs(siteId.String()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(7))

	r := net_db.NewSiteRepo(&UkamaDbMock{GormDb: gdb})

	count, err := r.CustomerCount(siteId.String())
	assert.NoError(t, err)
	assert.Equal(t, int64(7), count)

	assert.NoError(t, mock.ExpectationsWereMet())
}
