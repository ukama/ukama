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
	"errors"
	"log"
	"testing"
	"time"

	"github.com/ukama/ukama/systems/common/ukama"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	uuid "github.com/ukama/ukama/systems/common/uuid"
	int_db "github.com/ukama/ukama/systems/data-plan/package/pkg/db"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic("implement me")
}

func (u UkamaDbMock) Connect() error {
	panic("implement me")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func() error) error {
	log.Fatal("implement me")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) (err error) {
	log.Fatal("implement me")
	return nil
}

func Test_Package_Get(t *testing.T) {

	t.Run("PackageExistGet", func(t *testing.T) {
		// Arrange
		const uuidStr = "51fbba62-c79f-11eb-b8bc-0242ac130003"

		packID, _ := uuid.FromString(uuidStr)
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		pack := &int_db.Package{
			Uuid:        packID,
			Name:        "Silver Plan",
			SimType:     ukama.SimTypeTest,
			OwnerId:     uuid.NewV4(),
			Active:      true,
			Duration:    30,
			SmsVolume:   1000,
			DataVolume:  5000000,
			VoiceVolume: 500,
			Type:        ukama.PackageTypePostpaid,

			DataUnits:    ukama.DataUnitTypeMB,
			VoiceUnits:   ukama.CallUnitTypeSec,
			MessageUnits: ukama.MessageUnitTypeInt,
			Flatrate:     false,
			Currency:     "Dollar",
			From:         time.Now(),
			To:           time.Now().Add(time.Hour * 24 * 30),
			Country:      "USA",
			Provider:     "ukama",
		}

		rows := sqlmock.NewRows([]string{"uuid", "owner_id", "name", "sim_type", "active", "duration", "sms_volume", "data_volume", "voice_volume", "type", "data_units", "voice_units", "message_units", "flat_rate", "currency", "from", "to", "country", "provider"}).
			AddRow(packID, pack.OwnerId, pack.Name, pack.SimType, pack.Active, pack.Duration, pack.SmsVolume, pack.DataVolume, pack.VoiceVolume, pack.Type, pack.DataUnits, pack.VoiceUnits, pack.MessageUnits, pack.Flatrate, pack.Currency, pack.From, pack.To, pack.Country, pack.Provider)

		rrows := sqlmock.NewRows([]string{"package_id", "amount", "sms_mo", "sms_mt", "data"}).
			AddRow(packID, 100, 0.001, 0.001, 0.010)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(packID, sqlmock.AnyArg()).
			WillReturnRows(rows)
		mock.ExpectQuery(`^SELECT.*package_rates.*`).
			WithArgs(packID).
			WillReturnRows(rrows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		node, err := r.Get(packID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, node)
	})

	t.Run("PackageExistGetDetails", func(t *testing.T) {
		// Arrange
		const uuidStr = "51fbba62-c79f-11eb-b8bc-0242ac130003"
		baserate := uuid.NewV4().String()
		packID, _ := uuid.FromString(uuidStr)
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		pack := &int_db.Package{
			Uuid:        packID,
			Name:        "Silver Plan",
			SimType:     ukama.SimTypeTest,
			OwnerId:     uuid.NewV4(),
			Active:      true,
			Duration:    30,
			SmsVolume:   1000,
			DataVolume:  5000000,
			VoiceVolume: 500,
			Type:        ukama.PackageTypePostpaid,

			DataUnits:    ukama.DataUnitTypeMB,
			VoiceUnits:   ukama.CallUnitTypeSec,
			MessageUnits: ukama.MessageUnitTypeInt,
			Flatrate:     false,
			Currency:     "Dollar",
			From:         time.Now(),
			To:           time.Now().Add(time.Hour * 24 * 30),
			Country:      "USA",
			Provider:     "ukama",
		}

		rows := sqlmock.NewRows([]string{"uuid", "owner_id", "name", "sim_type", "active", "duration", "sms_volume", "data_volume", "voice_volume", "type", "data_units", "voice_units", "message_units", "flat_rate", "currency", "from", "to", "country", "provider"}).
			AddRow(packID, pack.OwnerId, pack.Name, pack.SimType, pack.Active, pack.Duration, pack.SmsVolume, pack.DataVolume, pack.VoiceVolume, pack.Type, pack.DataUnits, pack.VoiceUnits, pack.MessageUnits, pack.Flatrate, pack.Currency, pack.From, pack.To, pack.Country, pack.Provider)

		drows := sqlmock.NewRows([]string{"package_id", "dlbr", "ulbr", "apn"}).
			AddRow(packID, 1024000, 102400, "uakam.tel")

		rrows := sqlmock.NewRows([]string{"package_id", "amount", "sms_mo", "sms_mt", "data"}).
			AddRow(packID, 100, 0.001, 0.001, 0.010)

		mrows := sqlmock.NewRows([]string{"package_id", "base_rate_id", "markup"}).
			AddRow(packID, baserate, 20)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(packID, sqlmock.AnyArg()).
			WillReturnRows(rows)

		mock.ExpectQuery(`^SELECT.*package_details.*`).
			WithArgs(packID).
			WillReturnRows(drows)

		mock.ExpectQuery(`^SELECT.*package_markups.*`).
			WithArgs(packID).
			WillReturnRows(mrows)

		mock.ExpectQuery(`^SELECT.*package_rates.*`).
			WithArgs(packID).
			WillReturnRows(rrows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		node, err := r.GetDetails(packID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, node)
	})
}

func Test_Package_GetAll(t *testing.T) {

	t.Run("PackageExist", func(t *testing.T) {
		// Arrange
		const uuidStr = "51fbba62-c79f-11eb-b8bc-0242ac130003"
		packID, _ := uuid.FromString(uuidStr)
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		pack := &int_db.Package{
			Uuid:        packID,
			Name:        "Silver Plan",
			SimType:     ukama.SimTypeTest,
			OwnerId:     uuid.NewV4(),
			Active:      true,
			Duration:    30,
			SmsVolume:   1000,
			DataVolume:  5000000,
			VoiceVolume: 500,
			Type:        ukama.PackageTypePostpaid,

			DataUnits:    ukama.DataUnitTypeMB,
			VoiceUnits:   ukama.CallUnitTypeSec,
			MessageUnits: ukama.MessageUnitTypeInt,
			Flatrate:     false,
			Currency:     "Dollar",
			From:         time.Now(),
			To:           time.Now().Add(time.Hour * 24 * 30),
			Country:      "USA",
			Provider:     "ukama",
		}

		rows := sqlmock.NewRows([]string{"uuid", "owner_id", "name", "sim_type", "active", "duration", "sms_volume", "data_volume", "voice_volume", "type", "data_units", "voice_units", "message_units", "flat_rate", "currency", "from", "to", "country", "provider"}).
			AddRow(packID, pack.OwnerId, pack.Name, pack.SimType, pack.Active, pack.Duration, pack.SmsVolume, pack.DataVolume, pack.VoiceVolume, pack.Type, pack.DataUnits, pack.VoiceUnits, pack.MessageUnits, pack.Flatrate, pack.Currency, pack.From, pack.To, pack.Country, pack.Provider)

		rrows := sqlmock.NewRows([]string{"package_id", "amount", "sms_mo", "sms_mt", "data"}).
			AddRow(packID, 100, 0.001, 0.001, 0.010)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnRows(rows)

		mock.ExpectQuery(`^SELECT.*package_rates.*`).
			WithArgs(packID).
			WillReturnRows(rrows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		node, err := r.GetAll()

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, node)
	})

	t.Run("MultiplePackages", func(t *testing.T) {
		// Arrange
		packID1, _ := uuid.FromString("51fbba62-c79f-11eb-b8bc-0242ac130003")
		packID2, _ := uuid.FromString("51fbba62-c79f-11eb-b8bc-0242ac130004")
		ownerID1 := uuid.NewV4()
		ownerID2 := uuid.NewV4()

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		// Mock multiple packages
		rows := sqlmock.NewRows([]string{"uuid", "owner_id", "name", "sim_type", "active", "duration", "sms_volume", "data_volume", "voice_volume", "type", "data_units", "voice_units", "message_units", "flat_rate", "currency", "from", "to", "country", "provider"}).
			AddRow(packID1, ownerID1, "Silver Plan", ukama.SimTypeTest, true, 30, 1000, 5000000, 500, ukama.PackageTypePostpaid, ukama.DataUnitTypeMB, ukama.CallUnitTypeSec, ukama.MessageUnitTypeInt, false, "Dollar", time.Now(), time.Now().Add(time.Hour*24*30), "USA", "ukama").
			AddRow(packID2, ownerID2, "Gold Plan", ukama.SimTypeTest, true, 60, 2000, 10000000, 1000, ukama.PackageTypePrepaid, ukama.DataUnitTypeGB, ukama.CallUnitTypeMin, ukama.MessageUnitTypeInt, true, "Euro", time.Now(), time.Now().Add(time.Hour*24*60), "UK", "ukama")

		rrows := sqlmock.NewRows([]string{"package_id", "amount", "sms_mo", "sms_mt", "data"}).
			AddRow(packID1, 100, 0.001, 0.001, 0.010).
			AddRow(packID2, 200, 0.002, 0.002, 0.020)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnRows(rows)

		mock.ExpectQuery(`^SELECT.*package_rates.*`).
			WithArgs(packID1, packID2).
			WillReturnRows(rrows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		// Act
		packages, err := r.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, packages, 2)
		assert.Equal(t, packID1, packages[0].Uuid)
		assert.Equal(t, "Silver Plan", packages[0].Name)
		assert.Equal(t, packID2, packages[1].Uuid)
		assert.Equal(t, "Gold Plan", packages[1].Name)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAll_EmptyResult", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		// Return empty result
		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		packages, err := r.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, packages)
	})

	t.Run("GetAll_DatabaseError", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		// Mock database error
		mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnError(errors.New("connection timeout"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		packages, err := r.GetAll()
		assert.Error(t, err)
		assert.Nil(t, packages)
		assert.Contains(t, err.Error(), "connection timeout")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAll_PackageRatePreloadError", func(t *testing.T) {
		packID, _ := uuid.FromString("51fbba62-c79f-11eb-b8bc-0242ac130003")
		ownerID := uuid.NewV4()

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		// Mock successful package query
		rows := sqlmock.NewRows([]string{"uuid", "owner_id", "name", "sim_type", "active", "duration", "sms_volume", "data_volume", "voice_volume", "type", "data_units", "voice_units", "message_units", "flat_rate", "currency", "from", "to", "country", "provider"}).
			AddRow(packID, ownerID, "Silver Plan", ukama.SimTypeTest, true, 30, 1000, 5000000, 500, ukama.PackageTypePostpaid, ukama.DataUnitTypeMB, ukama.CallUnitTypeSec, ukama.MessageUnitTypeInt, false, "Dollar", time.Now(), time.Now().Add(time.Hour*24*30), "USA", "ukama")

		// Mock package rate query error
		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnRows(rows)

		mock.ExpectQuery(`^SELECT.*package_rates.*`).
			WithArgs(packID).
			WillReturnError(errors.New("package rate query failed"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		packages, err := r.GetAll()
		assert.Error(t, err)
		assert.Nil(t, packages)
		assert.Contains(t, err.Error(), "package rate query failed")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAll_PackageWithoutRate", func(t *testing.T) {
		// Arrange
		packID, _ := uuid.FromString("51fbba62-c79f-11eb-b8bc-0242ac130003")
		ownerID := uuid.NewV4()

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"uuid", "owner_id", "name", "sim_type", "active", "duration", "sms_volume", "data_volume", "voice_volume", "type", "data_units", "voice_units", "message_units", "flat_rate", "currency", "from", "to", "country", "provider"}).
			AddRow(packID, ownerID, "No Rate Plan", ukama.SimTypeTest, true, 30, 1000, 5000000, 500, ukama.PackageTypePostpaid, ukama.DataUnitTypeMB, ukama.CallUnitTypeSec, ukama.MessageUnitTypeInt, false, "Dollar", time.Now(), time.Now().Add(time.Hour*24*30), "USA", "ukama")

		emptyRateRows := sqlmock.NewRows([]string{"package_id", "amount", "sms_mo", "sms_mt", "data"})

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnRows(rows)

		mock.ExpectQuery(`^SELECT.*package_rates.*`).
			WithArgs(packID).
			WillReturnRows(emptyRateRows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		// Act
		packages, err := r.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, packages, 1)
		assert.Equal(t, packID, packages[0].Uuid)
		assert.Equal(t, "No Rate Plan", packages[0].Name)
		assert.Equal(t, uuid.Nil, packages[0].PackageRate.PackageID)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_Package_Get_Error(t *testing.T) {
	t.Run("Get_NotFound", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnError(gorm.ErrRecordNotFound)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		package_, err := r.Get(packID)
		assert.Error(t, err)
		assert.Nil(t, package_)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Get_DBError", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnError(errors.New("database error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		package_, err := r.Get(packID)
		assert.Error(t, err)
		assert.Nil(t, package_)
		assert.Contains(t, err.Error(), "database error")
	})
}

func Test_Package_GetDetails_Error(t *testing.T) {
	t.Run("GetDetails_NotFound", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnError(gorm.ErrRecordNotFound)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		package_, err := r.GetDetails(packID)
		assert.Error(t, err)
		assert.Nil(t, package_)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("GetDetails_DBError", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnError(errors.New("database error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		package_, err := r.GetDetails(packID)
		assert.Error(t, err)
		assert.Nil(t, package_)
		assert.Contains(t, err.Error(), "database error")
	})
}

func Test_Package_Add(t *testing.T) {
	t.Run("Add_Success", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		package_ := &int_db.Package{
			Uuid:    packID,
			Name:    "Test Package",
			Active:  true,
			Country: "USA",
		}

		packageRate := &int_db.PackageRate{
			PackageID: packID,
			Amount:    100,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Add(package_, packageRate)
		assert.NoError(t, err)
	})

	t.Run("Add_TransactionBeginError", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		package_ := &int_db.Package{Uuid: packID, Name: "Test Package"}
		packageRate := &int_db.PackageRate{PackageID: packID, Amount: 100}

		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Add(package_, packageRate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")
	})

	t.Run("Add_PackageCreationError", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		package_ := &int_db.Package{Uuid: packID, Name: "Test Package"}
		packageRate := &int_db.PackageRate{PackageID: packID, Amount: 100}

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnError(errors.New("package creation error"))
		mock.ExpectRollback()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Add(package_, packageRate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "package creation error")
	})

	t.Run("Add_PackageRateCreationError", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		package_ := &int_db.Package{Uuid: packID, Name: "Test Package"}
		packageRate := &int_db.PackageRate{PackageID: packID, Amount: 100}

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnError(errors.New("package rate creation error"))
		mock.ExpectRollback()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Add(package_, packageRate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "package rate creation error")
	})

	t.Run("Add_CommitError", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		package_ := &int_db.Package{Uuid: packID, Name: "Test Package"}
		packageRate := &int_db.PackageRate{PackageID: packID, Amount: 100}

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit().WillReturnError(errors.New("commit error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Add(package_, packageRate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "commit error")
	})

	t.Run("Add_WithFullPackageData", func(t *testing.T) {
		packID := uuid.NewV4()
		ownerID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		package_ := &int_db.Package{
			Uuid:          packID,
			OwnerId:       ownerID,
			Name:          "Premium Data Plan",
			SimType:       ukama.SimTypeTest,
			Active:        true,
			Duration:      30,
			SmsVolume:     1000,
			DataVolume:    5000000,
			VoiceVolume:   500,
			Type:          ukama.PackageTypePostpaid,
			DataUnits:     ukama.DataUnitTypeMB,
			VoiceUnits:    ukama.CallUnitTypeSec,
			MessageUnits:  ukama.MessageUnitTypeInt,
			Flatrate:      false,
			Currency:      "USD",
			From:          time.Now(),
			To:            time.Now().AddDate(0, 1, 0),
			Country:       "USA",
			Provider:      "ukama",
			Overdraft:     50.0,
			TrafficPolicy: 1,
			Networks:      []string{"network1", "network2"},
			SyncStatus:    ukama.StatusTypeCompleted,
		}

		packageRate := &int_db.PackageRate{
			PackageID: packID,
			Amount:    99.99,
			SmsMo:     0.01,
			SmsMt:     0.01,
			Data:      0.05,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Add(package_, packageRate)
		assert.NoError(t, err)
	})

	t.Run("Add_WithMinimalPackageData", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		package_ := &int_db.Package{
			Uuid:    packID,
			Name:    "Basic Plan",
			Country: "USA",
		}

		packageRate := &int_db.PackageRate{
			PackageID: packID,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Add(package_, packageRate)
		assert.NoError(t, err)
	})

	t.Run("Add_WithZeroValues", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		package_ := &int_db.Package{
			Uuid:          packID,
			Name:          "Zero Plan",
			Active:        false,
			Duration:      0,
			SmsVolume:     0,
			DataVolume:    0,
			VoiceVolume:   0,
			Flatrate:      false,
			Currency:      "",
			Country:       "",
			Provider:      "",
			Overdraft:     0,
			TrafficPolicy: 0,
		}

		packageRate := &int_db.PackageRate{
			PackageID: packID,
			Amount:    0,
			SmsMo:     0,
			SmsMt:     0,
			Data:      0,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Add(package_, packageRate)
		assert.NoError(t, err)
	})

	t.Run("Add_WithNegativeValues", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		package_ := &int_db.Package{
			Uuid:    packID,
			Name:    "Negative Plan",
			Country: "USA",
		}

		packageRate := &int_db.PackageRate{
			PackageID: packID,
			Amount:    -10.0,
			SmsMo:     -0.01,
			SmsMt:     -0.01,
			Data:      -0.05,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Add(package_, packageRate)
		assert.NoError(t, err)
	})
}

func Test_Package_Delete(t *testing.T) {
	t.Run("Delete_Success", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE "packages" SET "deleted_at"`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Delete(packID)
		assert.NoError(t, err)
	})

	t.Run("Delete_DBError", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE "packages" SET "deleted_at"`).WillReturnError(errors.New("delete failed"))
		mock.ExpectRollback()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		err = r.Delete(packID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}

func Test_Package_Update(t *testing.T) {
	t.Run("Update_Success", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(`^UPDATE "packages" SET.*WHERE.*uuid.*RETURNING`).
			WithArgs(sqlmock.AnyArg(), "Updated Name", packID).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "name"}).AddRow(packID, "Updated Name"))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		pkg := &int_db.Package{Name: "Updated Name"}
		err = r.Update(packID, pkg)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Update_TransactionBeginError", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		pkg := &int_db.Package{Name: "Updated Name"}
		err = r.Update(packID, pkg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")
	})

	t.Run("Update_NoRowsAffected", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(`^UPDATE "packages" SET.*WHERE.*uuid.*RETURNING`).
			WithArgs(sqlmock.AnyArg(), "Updated Name", packID).
			WillReturnRows(sqlmock.NewRows([]string{}))
		mock.ExpectRollback()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		pkg := &int_db.Package{Name: "Updated Name"}
		err = r.Update(packID, pkg)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Update_DatabaseError", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(`^UPDATE "packages" SET.*WHERE.*uuid.*RETURNING`).
			WithArgs(sqlmock.AnyArg(), "Updated Name", packID).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		pkg := &int_db.Package{Name: "Updated Name"}
		err = r.Update(packID, pkg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("Update_CommitError", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(`^UPDATE "packages" SET.*WHERE.*uuid.*RETURNING`).
			WithArgs(sqlmock.AnyArg(), "Updated Name", packID).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "name"}).AddRow(packID, "Updated Name"))
		mock.ExpectCommit().WillReturnError(errors.New("commit error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		pkg := &int_db.Package{Name: "Updated Name"}
		err = r.Update(packID, pkg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "commit error")
	})

	t.Run("Update_WithMultipleFields", func(t *testing.T) {
		packID := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(`^UPDATE "packages" SET.*WHERE.*uuid.*RETURNING`).
			WithArgs(sqlmock.AnyArg(), "Premium Plan", true, 60, packID).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "name", "active", "duration"}).AddRow(packID, "Premium Plan", true, 60))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{GormDb: gdb})

		pkg := &int_db.Package{
			Name:     "Premium Plan",
			Active:   true,
			Duration: 60,
		}
		err = r.Update(packID, pkg)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
