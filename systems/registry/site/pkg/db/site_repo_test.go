/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	extsql "database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ukama/ukama/systems/common/uuid"
	db_site "github.com/ukama/ukama/systems/registry/site/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/tj/assert"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	return nil
}

func (u UkamaDbMock) Connect() error {
	return nil
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	return nil
}

func TestSiteRepo_GetSite(t *testing.T) {
	t.Run("SiteExist", func(t *testing.T) {
		// Arrange
		siteId := uuid.NewV4()
		site := db_site.Site{
			Id:            siteId,
			Name:          "pamoja-net",
			NetworkId:     uuid.NewV4(),
			BackhaulId:    uuid.NewV4(),
			PowerId:       uuid.NewV4(),
			AccessId:      uuid.NewV4(),
			SwitchId:      uuid.NewV4(),
			IsDeactivated: false,
			Latitude:      40.7128,
			Longitude:     -74.0060,
			InstallDate:   "07-03-2023",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		}

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "latitude", "network_id"}).
			AddRow(siteId, site.Name, site.Latitude, site.NetworkId)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(site.Id, 1).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		rm, err := r.Get(site.Id)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, rm)
		assert.Equal(t, siteId, rm.Id)
		assert.Equal(t, site.Name, rm.Name)
		assert.Equal(t, site.NetworkId, rm.NetworkId)
		assert.Equal(t, site.Latitude, rm.Latitude)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SiteNotFound", func(t *testing.T) {
		// Arrange
		siteId := uuid.NewV4()

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(siteId, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		site, err := r.Get(siteId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, site)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		siteId := uuid.NewV4()

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(siteId, 1).
			WillReturnError(fmt.Errorf("database connection error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		site, err := r.Get(siteId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, site)
		assert.Contains(t, err.Error(), "database connection error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSiteRepo_GetSites(t *testing.T) {
	t.Run("SitesExist", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()
		site1 := db_site.Site{
			Id:          uuid.NewV4(),
			Name:        "Site1",
			NetworkId:   netID,
			Latitude:    40.7128,
			Longitude:   -74.0060,
			InstallDate: "07-03-2023",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   gorm.DeletedAt{},
		}
		site2 := db_site.Site{
			Id:          uuid.NewV4(),
			Name:        "Site2",
			NetworkId:   netID,
			Latitude:    40.7128,
			Longitude:   -74.0060,
			InstallDate: "07-03-2023",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   gorm.DeletedAt{},
		}

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "latitude", "network_id"}).
			AddRow(site1.Id, site1.Name, site1.Latitude, site1.NetworkId).
			AddRow(site2.Id, site2.Name, site2.Latitude, site2.NetworkId)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(netID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		sites, err := r.GetSites(netID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, sites)
		assert.Equal(t, 2, len(sites))
		assert.Equal(t, site1.Id, sites[0].Id)
		assert.Equal(t, site2.Id, sites[1].Id)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NoSitesFound", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		// Return empty result set
		rows := sqlmock.NewRows([]string{"id", "name", "latitude", "network_id"})

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(netID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sites, err := r.GetSites(netID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, sites)
		assert.Equal(t, 0, len(sites))

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SingleSiteFound", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()
		site := db_site.Site{
			Id:          uuid.NewV4(),
			Name:        "SingleSite",
			NetworkId:   netID,
			Latitude:    40.7128,
			Longitude:   -74.0060,
			InstallDate: "07-03-2023",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   gorm.DeletedAt{},
		}

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "latitude", "network_id"}).
			AddRow(site.Id, site.Name, site.Latitude, site.NetworkId)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(netID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sites, err := r.GetSites(netID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, sites)
		assert.Equal(t, 1, len(sites))
		assert.Equal(t, site.Id, sites[0].Id)
		assert.Equal(t, site.Name, sites[0].Name)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(netID).
			WillReturnError(fmt.Errorf("database connection error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		sites, err := r.GetSites(netID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, sites)
		assert.Contains(t, err.Error(), "database connection error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSiteRepo_Add(t *testing.T) {
	t.Run("ValidSite", func(t *testing.T) {
		// Arrange
		site := &db_site.Site{
			Id:            uuid.NewV4(),
			Name:          "valid-site",
			Location:      "Test Location",
			NetworkId:     uuid.NewV4(),
			BackhaulId:    uuid.NewV4(),
			SpectrumId:    uuid.NewV4(),
			PowerId:       uuid.NewV4(),
			AccessId:      uuid.NewV4(),
			SwitchId:      uuid.NewV4(),
			IsDeactivated: false,
			Latitude:      40.7128,
			Longitude:     -74.0060,
			InstallDate:   "07-03-2023",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		}

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(`^INSERT INTO "sites"`).
			WithArgs(
				site.Id, site.Name, site.Location, site.NetworkId, site.BackhaulId,
				site.SpectrumId, site.PowerId, site.AccessId, site.SwitchId, site.IsDeactivated,
				site.Latitude, site.Longitude, site.InstallDate,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(site, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	t.Run("InvalidSiteName", func(t *testing.T) {
		// Arrange
		invalidSite := &db_site.Site{
			Id:            uuid.NewV4(),
			Name:          "invalid_site_name!",
			Location:      "Test Location",
			NetworkId:     uuid.NewV4(),
			BackhaulId:    uuid.NewV4(),
			SpectrumId:    uuid.NewV4(),
			PowerId:       uuid.NewV4(),
			AccessId:      uuid.NewV4(),
			SwitchId:      uuid.NewV4(),
			IsDeactivated: false,
			Latitude:      40.7128,
			Longitude:     -74.0060,
			InstallDate:   "07-03-2023",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		}

		// Create a mock database
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

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err = r.Add(invalidSite, nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid name")
	})

	t.Run("DatabaseErrorDuringSiteCreation", func(t *testing.T) {
		// Arrange
		site := &db_site.Site{
			Id:            uuid.NewV4(),
			Name:          "valid-site-db-error",
			Location:      "Test Location",
			NetworkId:     uuid.NewV4(),
			BackhaulId:    uuid.NewV4(),
			SpectrumId:    uuid.NewV4(),
			PowerId:       uuid.NewV4(),
			AccessId:      uuid.NewV4(),
			SwitchId:      uuid.NewV4(),
			IsDeactivated: false,
			Latitude:      40.7128,
			Longitude:     -74.0060,
			InstallDate:   "07-03-2023",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		}

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(`^INSERT INTO "sites"`).
			WithArgs(
				site.Id, site.Name, site.Location, site.NetworkId, site.BackhaulId,
				site.SpectrumId, site.PowerId, site.AccessId, site.SwitchId, site.IsDeactivated,
				site.Latitude, site.Longitude, site.InstallDate,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(fmt.Errorf("database constraint violation"))
		mock.ExpectRollback()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err = r.Add(site, nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database constraint violation")
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseErrorDuringTransaction", func(t *testing.T) {
		// Arrange
		site := &db_site.Site{
			Id:            uuid.NewV4(),
			Name:          "valid-site-transaction-error",
			Location:      "Test Location",
			NetworkId:     uuid.NewV4(),
			BackhaulId:    uuid.NewV4(),
			SpectrumId:    uuid.NewV4(),
			PowerId:       uuid.NewV4(),
			AccessId:      uuid.NewV4(),
			SwitchId:      uuid.NewV4(),
			IsDeactivated: false,
			Latitude:      40.7128,
			Longitude:     -74.0060,
			InstallDate:   "07-03-2023",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		}

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(`^INSERT INTO "sites"`).
			WithArgs(
				site.Id, site.Name, site.Location, site.NetworkId, site.BackhaulId,
				site.SpectrumId, site.PowerId, site.AccessId, site.SwitchId, site.IsDeactivated,
				site.Latitude, site.Longitude, site.InstallDate,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(fmt.Errorf("transaction commit failed"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err = r.Add(site, nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction commit failed")
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSiteRepo_GetSiteCount(t *testing.T) {
	t.Run("ValidNetworkId", func(t *testing.T) {
		// Arrange
		networkId := uuid.NewV4()
		expectedCount := int64(5)

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"count"}).
			AddRow(expectedCount)

		mock.ExpectQuery(`^SELECT count\(.*\) FROM "sites" WHERE network_id = \$1`).
			WithArgs(networkId).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		count, err := r.GetSiteCount(networkId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ZeroSites", func(t *testing.T) {
		// Arrange
		networkId := uuid.NewV4()
		expectedCount := int64(0)

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"count"}).
			AddRow(expectedCount)

		mock.ExpectQuery(`^SELECT count\(.*\) FROM "sites" WHERE network_id = \$1`).
			WithArgs(networkId).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		count, err := r.GetSiteCount(networkId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		networkId := uuid.NewV4()

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT count\(.*\) FROM "sites" WHERE network_id = \$1`).
			WithArgs(networkId).
			WillReturnError(fmt.Errorf("database connection error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		count, err := r.GetSiteCount(networkId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.Contains(t, err.Error(), "database connection error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSiteRepo_Update(t *testing.T) {
	t.Run("ValidUpdate", func(t *testing.T) {
		// Arrange
		site := &db_site.Site{
			Id:            uuid.NewV4(),
			Name:          "updated-site",
			Location:      "Updated Location",
			NetworkId:     uuid.NewV4(),
			BackhaulId:    uuid.NewV4(),
			SpectrumId:    uuid.NewV4(),
			PowerId:       uuid.NewV4(),
			AccessId:      uuid.NewV4(),
			SwitchId:      uuid.NewV4(),
			IsDeactivated: true,
			Latitude:      42.3601,
			Longitude:     -71.0589,
			InstallDate:   "15-06-2023",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		}

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE "sites" SET`).
			WithArgs(
				site.Name, site.Location, site.NetworkId, site.BackhaulId,
				site.SpectrumId, site.PowerId, site.AccessId, site.SwitchId,
				site.IsDeactivated, site.Latitude, site.Longitude, site.InstallDate,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				site.Id,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := db_site.NewSiteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Update(site)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSiteRepo_List(t *testing.T) {
	t.Run("ValidNetworkId", func(t *testing.T) {
		netID := uuid.NewV4()
		site1 := db_site.Site{
			Id:            uuid.NewV4(),
			Name:          "Site1",
			NetworkId:     netID,
			IsDeactivated: false,
		}
		site2 := db_site.Site{
			Id:            uuid.NewV4(),
			Name:          "Site2",
			NetworkId:     netID,
			IsDeactivated: false,
		}

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "network_id", "is_deactivated"}).
			AddRow(site1.Id, site1.Name, site1.NetworkId, site1.IsDeactivated).
			AddRow(site2.Id, site2.Name, site2.NetworkId, site2.IsDeactivated)

		mock.ExpectQuery(`SELECT(.*?)FROM "sites"`).
			WithArgs(netID, false).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
		r := db_site.NewSiteRepo(&UkamaDbMock{GormDb: gdb})

		// Act
		sites, err := r.List(&netID, false)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, sites, 2)
		assert.Equal(t, site1.Id, sites[0].Id)
		assert.Equal(t, site2.Id, sites[1].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NilNetworkId", func(t *testing.T) {
		site := db_site.Site{
			Id:            uuid.NewV4(),
			Name:          "Site1",
			NetworkId:     uuid.NewV4(),
			IsDeactivated: false,
		}
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		rows := sqlmock.NewRows([]string{"id", "name", "network_id", "is_deactivated"}).
			AddRow(site.Id, site.Name, site.NetworkId, site.IsDeactivated)
		mock.ExpectQuery(`SELECT(.*?)FROM "sites"`).
			WithArgs(false).
			WillReturnRows(rows)
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
		r := db_site.NewSiteRepo(&UkamaDbMock{GormDb: gdb})
		sites, err := r.List(nil, false)
		assert.NoError(t, err)
		assert.Len(t, sites, 1)
		assert.Equal(t, site.Id, sites[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NullUUIDNetworkId", func(t *testing.T) {
		nullUUID := uuid.UUID{}
		site := db_site.Site{
			Id:            uuid.NewV4(),
			Name:          "Site1",
			NetworkId:     uuid.NewV4(),
			IsDeactivated: false,
		}
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		rows := sqlmock.NewRows([]string{"id", "name", "network_id", "is_deactivated"}).
			AddRow(site.Id, site.Name, site.NetworkId, site.IsDeactivated)
		mock.ExpectQuery(`SELECT(.*?)FROM "sites"`).
			WithArgs(false).
			WillReturnRows(rows)
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
		r := db_site.NewSiteRepo(&UkamaDbMock{GormDb: gdb})
		sites, err := r.List(&nullUUID, false)
		assert.NoError(t, err)
		assert.Len(t, sites, 1)
		assert.Equal(t, site.Id, sites[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		netID := uuid.NewV4()
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectQuery(`SELECT(.*?)FROM "sites"`).
			WithArgs(netID, false).
			WillReturnError(fmt.Errorf("db error"))
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
		r := db_site.NewSiteRepo(&UkamaDbMock{GormDb: gdb})
		sites, err := r.List(&netID, false)
		assert.Error(t, err)
		assert.Nil(t, sites)
		assert.Contains(t, err.Error(), "db error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
