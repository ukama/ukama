/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
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

// Common test data
var (
	testSiteId1            = uuid.NewV4()
	testSiteId2            = uuid.NewV4()
	testNetworkId          = uuid.NewV4()
	testBackhaulId         = uuid.NewV4()
	testSpectrumId         = uuid.NewV4()
	testPowerId            = uuid.NewV4()
	testAccessId           = uuid.NewV4()
	testSwitchId           = uuid.NewV4()
	testInstallDate        = "07-03-2023"
	testLocation           = "Test Location"
	testLatitude           = 40.7128
	testLongitude          = -74.0060
	testUpdatedLatitude    = 42.3601
	testUpdatedLongitude   = -71.0589
	testUpdatedInstallDate = "15-06-2023"
	testUpdatedLocation    = "Updated Location"
)

// Common test site data
var testSite = db_site.Site{
	Id:            testSiteId1,
	Name:          "pamoja-net",
	Location:      testLocation,
	NetworkId:     testNetworkId,
	BackhaulId:    testBackhaulId,
	SpectrumId:    testSpectrumId,
	PowerId:       testPowerId,
	AccessId:      testAccessId,
	SwitchId:      testSwitchId,
	IsDeactivated: false,
	Latitude:      testLatitude,
	Longitude:     testLongitude,
	InstallDate:   testInstallDate,
	CreatedAt:     time.Now(),
	UpdatedAt:     time.Now(),
	DeletedAt:     gorm.DeletedAt{},
}

var testSite2 = db_site.Site{
	Id:            testSiteId2,
	Name:          "Site2",
	Location:      testLocation,
	NetworkId:     testNetworkId,
	BackhaulId:    testBackhaulId,
	SpectrumId:    testSpectrumId,
	PowerId:       testPowerId,
	AccessId:      testAccessId,
	SwitchId:      testSwitchId,
	IsDeactivated: false,
	Latitude:      testLatitude,
	Longitude:     testLongitude,
	InstallDate:   testInstallDate,
	CreatedAt:     time.Now(),
	UpdatedAt:     time.Now(),
	DeletedAt:     gorm.DeletedAt{},
}

var updatedTestSite = db_site.Site{
	Id:            testSiteId1,
	Name:          "updated-site",
	Location:      testUpdatedLocation,
	NetworkId:     testNetworkId,
	BackhaulId:    testBackhaulId,
	SpectrumId:    testSpectrumId,
	PowerId:       testPowerId,
	AccessId:      testAccessId,
	SwitchId:      testSwitchId,
	IsDeactivated: true,
	Latitude:      testUpdatedLatitude,
	Longitude:     testUpdatedLongitude,
	InstallDate:   testUpdatedInstallDate,
	CreatedAt:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	UpdatedAt:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	DeletedAt:     gorm.DeletedAt{},
}

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

// Helper function to create mock database and repository
func createMockDBAndRepo(t *testing.T) (sqlmock.Sqlmock, db_site.SiteRepo, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	r := db_site.NewSiteRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	return mock, r, nil
}

func TestSiteRepo_GetSite(t *testing.T) {
	t.Run("SiteExist", func(t *testing.T) {
		// Arrange
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "latitude", "network_id"}).
			AddRow(testSite.Id, testSite.Name, testSite.Latitude, testSite.NetworkId)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(testSite.Id, 1).
			WillReturnRows(rows)

		// Act
		rm, err := r.Get(testSite.Id)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, rm)
		assert.Equal(t, testSiteId1, rm.Id)
		assert.Equal(t, testSite.Name, rm.Name)
		assert.Equal(t, testSite.NetworkId, rm.NetworkId)
		assert.Equal(t, testSite.Latitude, rm.Latitude)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SiteNotFound", func(t *testing.T) {
		// Arrange
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(testSiteId1, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Act
		site, err := r.Get(testSiteId1)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, site)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(testSiteId1, 1).
			WillReturnError(fmt.Errorf("database connection error"))

		// Act
		site, err := r.Get(testSiteId1)

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
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "latitude", "network_id"}).
			AddRow(testSite.Id, testSite.Name, testSite.Latitude, testSite.NetworkId).
			AddRow(testSite2.Id, testSite2.Name, testSite2.Latitude, testSite2.NetworkId)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(testNetworkId).
			WillReturnRows(rows)

		// Act
		sites, err := r.GetSites(testNetworkId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, sites)
		assert.Equal(t, 2, len(sites))
		assert.Equal(t, testSite.Id, sites[0].Id)
		assert.Equal(t, testSite2.Id, sites[1].Id)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NoSitesFound", func(t *testing.T) {
		// Arrange
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		// Return empty result set
		rows := sqlmock.NewRows([]string{"id", "name", "latitude", "network_id"})

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(testNetworkId).
			WillReturnRows(rows)

		// Act
		sites, err := r.GetSites(testNetworkId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, sites)
		assert.Equal(t, 0, len(sites))

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SingleSiteFound", func(t *testing.T) {
		// Arrange
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "latitude", "network_id"}).
			AddRow(testSite.Id, testSite.Name, testSite.Latitude, testSite.NetworkId)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(testNetworkId).
			WillReturnRows(rows)

		// Act
		sites, err := r.GetSites(testNetworkId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, sites)
		assert.Equal(t, 1, len(sites))
		assert.Equal(t, testSite.Id, sites[0].Id)
		assert.Equal(t, testSite.Name, sites[0].Name)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(testNetworkId).
			WillReturnError(fmt.Errorf("database connection error"))

		// Act
		sites, err := r.GetSites(testNetworkId)

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
		site := &testSite
		mock, r, err := createMockDBAndRepo(t)
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
			Location:      testLocation,
			NetworkId:     testNetworkId,
			BackhaulId:    testBackhaulId,
			SpectrumId:    testSpectrumId,
			PowerId:       testPowerId,
			AccessId:      testAccessId,
			SwitchId:      testSwitchId,
			IsDeactivated: false,
			Latitude:      testLatitude,
			Longitude:     testLongitude,
			InstallDate:   testInstallDate,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		}

		_, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

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
			Location:      testLocation,
			NetworkId:     testNetworkId,
			BackhaulId:    testBackhaulId,
			SpectrumId:    testSpectrumId,
			PowerId:       testPowerId,
			AccessId:      testAccessId,
			SwitchId:      testSwitchId,
			IsDeactivated: false,
			Latitude:      testLatitude,
			Longitude:     testLongitude,
			InstallDate:   testInstallDate,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		}

		mock, r, err := createMockDBAndRepo(t)
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
			Location:      testLocation,
			NetworkId:     testNetworkId,
			BackhaulId:    testBackhaulId,
			SpectrumId:    testSpectrumId,
			PowerId:       testPowerId,
			AccessId:      testAccessId,
			SwitchId:      testSwitchId,
			IsDeactivated: false,
			Latitude:      testLatitude,
			Longitude:     testLongitude,
			InstallDate:   testInstallDate,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		}

		mock, r, err := createMockDBAndRepo(t)
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

		// Act
		err = r.Add(site, nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction commit failed")
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NestedFunctionError", func(t *testing.T) {
		// Arrange
		site := &db_site.Site{
			Id:            uuid.NewV4(),
			Name:          "valid-site-nested-error",
			Location:      testLocation,
			NetworkId:     testNetworkId,
			BackhaulId:    testBackhaulId,
			SpectrumId:    testSpectrumId,
			PowerId:       testPowerId,
			AccessId:      testAccessId,
			SwitchId:      testSwitchId,
			IsDeactivated: false,
			Latitude:      testLatitude,
			Longitude:     testLongitude,
			InstallDate:   testInstallDate,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			DeletedAt:     gorm.DeletedAt{},
		}

		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		// Define a nested function that returns an error
		nestedFunc := func(site *db_site.Site, tx *gorm.DB) error {
			return fmt.Errorf("nested function validation failed")
		}

		mock.ExpectBegin()
		mock.ExpectRollback()

		// Act
		err = r.Add(site, nestedFunc)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nested function validation failed")
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSiteRepo_GetSiteCount(t *testing.T) {
	t.Run("ValidNetworkId", func(t *testing.T) {
		// Arrange
		expectedCount := int64(5)
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"count"}).
			AddRow(expectedCount)

		mock.ExpectQuery(`^SELECT count\(.*\) FROM "sites" WHERE network_id = \$1`).
			WithArgs(testNetworkId).
			WillReturnRows(rows)

		// Act
		count, err := r.GetSiteCount(testNetworkId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ZeroSites", func(t *testing.T) {
		// Arrange
		expectedCount := int64(0)
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"count"}).
			AddRow(expectedCount)

		mock.ExpectQuery(`^SELECT count\(.*\) FROM "sites" WHERE network_id = \$1`).
			WithArgs(testNetworkId).
			WillReturnRows(rows)

		// Act
		count, err := r.GetSiteCount(testNetworkId)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT count\(.*\) FROM "sites" WHERE network_id = \$1`).
			WithArgs(testNetworkId).
			WillReturnError(fmt.Errorf("database connection error"))

		// Act
		count, err := r.GetSiteCount(testNetworkId)

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
		site := &updatedTestSite
		mock, r, err := createMockDBAndRepo(t)
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

		// Act
		err = r.Update(site)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseUpdateError", func(t *testing.T) {
		// Arrange
		site := &updatedTestSite
		mock, r, err := createMockDBAndRepo(t)
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
			WillReturnError(fmt.Errorf("database constraint violation: site not found"))
		mock.ExpectRollback()

		// Act
		err = r.Update(site)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database constraint violation: site not found")
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestSiteRepo_List(t *testing.T) {
	t.Run("ValidNetworkId", func(t *testing.T) {
		// Arrange
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "network_id", "is_deactivated"}).
			AddRow(testSite.Id, testSite.Name, testSite.NetworkId, testSite.IsDeactivated).
			AddRow(testSite2.Id, testSite2.Name, testSite2.NetworkId, testSite2.IsDeactivated)

		mock.ExpectQuery(`SELECT(.*?)FROM "sites"`).
			WithArgs(testNetworkId, false).
			WillReturnRows(rows)

		// Act
		sites, err := r.List(&testNetworkId, false)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, sites, 2)
		assert.Equal(t, testSite.Id, sites[0].Id)
		assert.Equal(t, testSite2.Id, sites[1].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NilNetworkId", func(t *testing.T) {
		// Arrange
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "network_id", "is_deactivated"}).
			AddRow(testSite.Id, testSite.Name, testSite.NetworkId, testSite.IsDeactivated)

		mock.ExpectQuery(`SELECT(.*?)FROM "sites"`).
			WithArgs(false).
			WillReturnRows(rows)

		// Act
		sites, err := r.List(nil, false)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, sites, 1)
		assert.Equal(t, testSite.Id, sites[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NullUUIDNetworkId", func(t *testing.T) {
		// Arrange
		nullUUID := uuid.UUID{}
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "network_id", "is_deactivated"}).
			AddRow(testSite.Id, testSite.Name, testSite.NetworkId, testSite.IsDeactivated)

		mock.ExpectQuery(`SELECT(.*?)FROM "sites"`).
			WithArgs(false).
			WillReturnRows(rows)

		// Act
		sites, err := r.List(&nullUUID, false)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, sites, 1)
		assert.Equal(t, testSite.Id, sites[0].Id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mock, r, err := createMockDBAndRepo(t)
		assert.NoError(t, err)

		mock.ExpectQuery(`SELECT(.*?)FROM "sites"`).
			WithArgs(testNetworkId, false).
			WillReturnError(fmt.Errorf("db error"))

		// Act
		sites, err := r.List(&testNetworkId, false)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, sites)
		assert.Contains(t, err.Error(), "db error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
