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
			Latitude:      40.7128,   // Dummy latitude
			Longitude:     -74.0060,  // Dummy longitude
			InstallDate:   "07-03-2023", // Current time as install date
			 CreatedAt: time.Now(),
			 UpdatedAt: time.Now(),
			 DeletedAt: gorm.DeletedAt{},
		 }
 
		 var db *extsql.DB
 
		 db, mock, err := sqlmock.New() // mock sql.DB
		 assert.NoError(t, err)
 
		 rows := sqlmock.NewRows([]string{"id", "name", "latitude","network_id"}).
			 AddRow(siteId, site.Name, site.Latitude,site.NetworkId)


			 mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(site.Id,site.NetworkId, 1). // Add '1' for the LIMIT clause
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
 
		 err = mock.ExpectationsWereMet()
		 assert.NoError(t, err)
		 assert.NotNil(t, rm)
	 })
 }
 
 func TestSiteRepo_GetSites(t *testing.T) {
    t.Run("SitesExist", func(t *testing.T) {
        // Arrange
        netID := uuid.NewV4()
        site1 := db_site.Site{
            Id:            uuid.NewV4(),
            Name:          "Site1",
            NetworkId:     netID,
            Latitude:      40.7128,
            Longitude:     -74.0060,
            InstallDate:   "07-03-2023",
            CreatedAt:     time.Now(),
            UpdatedAt:     time.Now(),
            DeletedAt:     gorm.DeletedAt{},
        }
        site2 := db_site.Site{
            Id:            uuid.NewV4(),
            Name:          "Site2",
            NetworkId:     netID,
            Latitude:      40.7128,
            Longitude:     -74.0060,
            InstallDate:    "07-03-2023",
            CreatedAt:     time.Now(),
            UpdatedAt:     time.Now(),
            DeletedAt:     gorm.DeletedAt{},
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
        assert.Equal(t, 2, len(sites)) // Ensure that both sites are retrieved
        assert.Equal(t, site1.Id, sites[0].Id)
        assert.Equal(t, site2.Id, sites[1].Id)

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
            NetworkId:     uuid.NewV4(),
            BackhaulId:    uuid.NewV4(),
            AccessId:      uuid.NewV4(),
            PowerId:       uuid.NewV4(),
            SwitchId:      uuid.NewV4(),
            IsDeactivated: false,
            Latitude:      40.7128,
            Longitude:     -74.0060,
            InstallDate:    "07-03-2023",
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
                site.Id, site.Name, site.NetworkId, site.BackhaulId,
                site.PowerId, site.AccessId, site.SwitchId, site.IsDeactivated,
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
            Name:          "invalid_site_name!", // Invalid site name with special characters
            NetworkId:     uuid.NewV4(),
            Latitude:      40.7128,
            Longitude:     -74.0060,
            InstallDate:    "07-03-2023",
            CreatedAt:     time.Now(),
            UpdatedAt:     time.Now(),
            DeletedAt:     gorm.DeletedAt{},
        }

        r := db_site.NewSiteRepo(&UkamaDbMock{})

        // Act
        err := r.Add(invalidSite, nil)

        // Assert
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "invalid name")
    })
}


