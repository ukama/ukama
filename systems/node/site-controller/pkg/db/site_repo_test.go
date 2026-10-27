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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ukamaDbMock struct {
	GormDb *gorm.DB
}

func (u ukamaDbMock) Init(model ...interface{}) error { panic("implement me: Init()") }
func (u ukamaDbMock) Connect() error                  { panic("implement me: Connect()") }
func (u ukamaDbMock) GetGormDb() *gorm.DB             { return u.GormDb }
func (u ukamaDbMock) InitDB() error                   { return nil }
func (u ukamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func() error) error {
	panic("implement me: ExecuteInTransaction()")
}
func (u ukamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) error {
	panic("implement me: ExecuteInTransaction2()")
}

func newSiteRepo(t *testing.T) (db.SiteRepo, sqlmock.Sqlmock) {
	conn, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 conn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return db.NewSiteRepo(ukamaDbMock{GormDb: gdb}), mock
}

func TestSiteRepoDelete_RemovesSiteAndChildRows(t *testing.T) {
	siteID := "11111111-1111-1111-1111-111111111111"
	repo, mock := newSiteRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "site_intent_flights" WHERE site_intent_id IN \(SELECT "id" FROM "site_intents" WHERE site_id = \$1\)`).
		WithArgs(siteID).WillReturnResult(sqlmock.NewResult(0, 2))
	for _, table := range []string{"site_intents", "site_states", "site_components", "site_port_maps", "sites"} {
		mock.ExpectExec(`DELETE FROM "` + table + `" WHERE site_id = \$1`).
			WithArgs(siteID).WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectCommit()

	assert.NoError(t, repo.Delete(siteID))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepoDelete_RollsBackOnFailure(t *testing.T) {
	siteID := "22222222-2222-2222-2222-222222222222"
	repo, mock := newSiteRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "site_intent_flights"`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(`DELETE FROM "site_intents"`).WillReturnError(assert.AnError)
	mock.ExpectRollback()

	assert.Error(t, repo.Delete(siteID))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSiteRepoDelete_RequiresSiteID(t *testing.T) {
	repo, _ := newSiteRepo(t)

	assert.Error(t, repo.Delete(""))
}
