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
}
