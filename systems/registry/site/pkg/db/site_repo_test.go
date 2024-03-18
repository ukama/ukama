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
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"database/sql"
	extsql "database/sql"

	site_db "github.com/ukama/ukama/systems/registry/site/pkg/db"
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
 func Test_SiteRepo_Get(t *testing.T) {
	 t.Run("SiteExist", func(t *testing.T) {
		 // Arrange
		 const siteName = "site1"
		 var siteId = uuid.NewV4()
		 var netId = uuid.NewV4()
 
		 var db *extsql.DB
 
		 db, mock, err := sqlmock.New() // mock sql.DB
		 assert.NoError(t, err)
 
		 rows := sqlmock.NewRows([]string{"id", "name", "network_id"}).
			 AddRow(siteId, siteName, netId)
 
		 mock.ExpectQuery(`^SELECT.*sites.*`).
			 WithArgs(siteId).
			 WillReturnRows(rows)
 
		 dialector := postgres.New(postgres.Config{
			 DSN:                  "sqlmock_db_0",
			 DriverName:           "postgres",
			 Conn:                 db,
			 PreferSimpleProtocol: true,
		 })
 
		 gdb, err := gorm.Open(dialector, &gorm.Config{})
		 assert.NoError(t, err)
 
		 r := site_db.NewSiteRepo(&UkamaDbMock{
			 GormDb: gdb,
		 })
 
		 assert.NoError(t, err)
 
		 // Act
		 site, err := r.Get(siteId,netId)
 
		 // Assert
		 assert.NoError(t, err)
 
		 err = mock.ExpectationsWereMet()
		 assert.NoError(t, err)
		 assert.NotNil(t, site)
	 })
 
	 t.Run("SiteNotFound", func(t *testing.T) {
		 // Arrange
		 var siteId = uuid.NewV4()
		 var netId = uuid.NewV4()

		 var db *extsql.DB
 
		 db, mock, err := sqlmock.New() // mock sql.DB
		 assert.NoError(t, err)
 
		 mock.ExpectQuery(`^SELECT.*sites.*`).
			 WithArgs(siteId).
			 WillReturnError(sql.ErrNoRows)
 
		 dialector := postgres.New(postgres.Config{
			 DSN:                  "sqlmock_db_0",
			 DriverName:           "postgres",
			 Conn:                 db,
			 PreferSimpleProtocol: true,
		 })
 
		 gdb, err := gorm.Open(dialector, &gorm.Config{})
		 assert.NoError(t, err)
 
		 r := site_db.NewSiteRepo(&UkamaDbMock{
			 GormDb: gdb,
		 })
 
		 assert.NoError(t, err)
 
		 // Act
		 site, err := r.Get(siteId,netId)
 
		 // Assert
		 assert.Error(t, err)
 
		 err = mock.ExpectationsWereMet()
		 assert.NoError(t, err)
		 assert.Nil(t, site)
	 })
 }
 

 func Test_SiteRepo_Add(t *testing.T) {
	t.Run("AddSite", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		site := site_db.Site{
			ID:            uuid.NewV4(),
			Name:          "Site A",
			NetworkID:     uuid.NewV4(),
			BackhaulID:    uuid.NewV4(),
			PowerID:       uuid.NewV4(),
			AccessID:      uuid.NewV4(),
			SwitchID:      uuid.NewV4(),
			IsDeactivated: false,
			Latitude:      40.7128,   // Dummy latitude
			Longitude:     -74.0060,  // Dummy longitude
			InstallDate:   time.Now(), // Current time as install date
		}

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(site.ID, site.Name, site.NetworkID,site.BackhaulID,site.PowerID,site.AccessID,site.SwitchID,site.IsDeactivated,site.Latitude,site.Longitude,site.InstallDate, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := site_db.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&site, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
