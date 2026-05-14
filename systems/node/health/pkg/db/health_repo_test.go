/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"log"
	"testing"
	"time"

	"github.com/ukama/ukama/systems/common/ukama"
	int_db "github.com/ukama/ukama/systems/node/health/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func TestHealthRepoList(t *testing.T) {

	t.Run("all reports for node", func(t *testing.T) {
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"nodeId"}).
			AddRow(nid.String())

		mock.ExpectQuery(`SELECT.*"health_reports"`).
			WithArgs(nid.String()).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewHealthRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		reports, err := r.List("", nid.String(), nil, ukama.FilterTimeframesTypeAll)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		if assert.Len(t, reports, 1) {
			assert.Equal(t, nid.String(), reports[0].NodeID)
		}
	})

	t.Run("no reports for node", func(t *testing.T) {
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`SELECT.*"health_reports"`).
			WithArgs(nid.String()).
			WillReturnRows(sqlmock.NewRows([]string{"nodeId"}))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewHealthRepo(&UkamaDbMock{GormDb: gdb})
		_, err = r.List("", nid.String(), nil, ukama.FilterTimeframesTypeAll)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("latest row exists", func(t *testing.T) {
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rid := "a0000000-0000-4000-8000-000000000001"
		rows := sqlmock.NewRows([]string{"nodeId", "reportId"}).
			AddRow(nid.String(), rid)

		mock.ExpectQuery(`SELECT.*"node_latest_healths"`).
			WithArgs(nid.String()).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewHealthRepo(&UkamaDbMock{GormDb: gdb})
		reports, err := r.List("", nid.String(), nil, ukama.FilterTimeframesTypeLatest)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		if assert.Len(t, reports, 1) {
			assert.Equal(t, nid.String(), reports[0].NodeID)
		}
	})

	t.Run("latest empty", func(t *testing.T) {
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`SELECT.*"node_latest_healths"`).
			WithArgs(nid.String()).
			WillReturnRows(sqlmock.NewRows([]string{"nodeId"}))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewHealthRepo(&UkamaDbMock{GormDb: gdb})
		reports, err := r.List("", nid.String(), nil, ukama.FilterTimeframesTypeLatest)

		assert.NoError(t, err)
		assert.Empty(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestHealthRepoStoreHealthReport(t *testing.T) {
	t.Run("nil report", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
		r := int_db.NewHealthRepo(&UkamaDbMock{GormDb: gdb})
		err = r.StoreHealthReport(nil, time.Now().UTC())
		assert.Error(t, err)
	})
}
