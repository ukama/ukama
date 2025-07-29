/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// ICCID values
	TestIccid1        = "8910300000003540855"
	TestIccid2        = "8910300000003540856"
	TestIccid3        = "8910300000003540859"
	TestIccidNotFound = "notfound-iccid"

	// MSISDN values
	TestMsisdn1 = "01010101"
	TestMsisdn2 = "01010102"

	// SmDpAddress values
	TestSmDpAddress1 = "123456789"
	TestSmDpAddress2 = "123456790"

	// ActivationCode values
	TestActivationCode1 = "0000"
	TestActivationCode2 = "0001"

	// QR Code values
	TestQrCode1 = "123456789"
	TestQrCode2 = "123456790"

	// Test IDs
	TestId1        = uint64(1)
	TestId2        = uint64(2)
	TestId3        = uint64(3)
	TestIdNotFound = uint64(999)

	// Database configuration
	TestDbDSN    = "sqlmock_db_0"
	TestDbDriver = "postgres"
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

func Test_GetSimsByType(t *testing.T) {
	t.Run("GetSimsByType_Success", func(t *testing.T) {
		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})
		rows.AddRow(
			TestIccid1,
			TestMsisdn1,
			false,
			false,
			ukama.SimTypeTest,
			TestSmDpAddress1,
			TestActivationCode1,
			true,
			TestQrCode1,
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(ukama.SimTypeTest).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  TestDbDSN,
			DriverName:           TestDbDriver,
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sp, err := r.GetSimsByType(ukama.SimTypeTest)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sp)
		assert.Len(t, sp, 1)
		assert.Equal(t, TestIccid1, sp[0].Iccid)
	})

	t.Run("GetSimsByType_MultipleResults", func(t *testing.T) {
		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})
		rows.AddRow(
			TestIccid1,
			TestMsisdn1,
			false,
			false,
			ukama.SimTypeOperatorData,
			TestSmDpAddress1,
			TestActivationCode1,
			true,
			TestQrCode1,
		)
		rows.AddRow(
			TestIccid2,
			TestMsisdn2,
			true,
			false,
			ukama.SimTypeOperatorData,
			TestSmDpAddress2,
			TestActivationCode2,
			false,
			TestQrCode2,
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(ukama.SimTypeOperatorData).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  TestDbDSN,
			DriverName:           TestDbDriver,
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sp, err := r.GetSimsByType(ukama.SimTypeOperatorData)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sp)
		assert.Len(t, sp, 2)
		assert.Equal(t, TestIccid1, sp[0].Iccid)
		assert.Equal(t, TestIccid2, sp[1].Iccid)
	})

	t.Run("GetSimsByType_EmptyResult", func(t *testing.T) {
		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(ukama.SimTypeUkamaData).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  TestDbDSN,
			DriverName:           TestDbDriver,
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sp, err := r.GetSimsByType(ukama.SimTypeUkamaData)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sp)
		assert.Len(t, sp, 0)
	})

	t.Run("GetSimsByType_DatabaseError", func(t *testing.T) {
		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(ukama.SimTypeTest).
			WillReturnError(gorm.ErrInvalidDB)

		dialector := postgres.New(postgres.Config{
			DSN:                  TestDbDSN,
			DriverName:           TestDbDriver,
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sp, err := r.GetSimsByType(ukama.SimTypeTest)
		assert.Error(t, err)
		assert.Nil(t, sp)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_GetByIccid(t *testing.T) {
	t.Run("GetByIccid", func(t *testing.T) {

		var db *extsql.DB
		iccid := TestIccid1
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})
		rows.AddRow(
			iccid,
			TestMsisdn1,
			false,
			false,
			ukama.SimTypeTest,
			TestSmDpAddress1,
			TestActivationCode1,
			true,
			TestQrCode1,
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(false, iccid, sqlmock.AnyArg()).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  TestDbDSN,
			DriverName:           TestDbDriver,
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		sp, err := r.GetByIccid(iccid)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sp)
	})

	t.Run("GetByIccid_RecordNotFound", func(t *testing.T) {
		var db *extsql.DB
		iccid := TestIccidNotFound
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(false, iccid, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		dialector := postgres.New(postgres.Config{
			DSN:                  TestDbDSN,
			DriverName:           TestDbDriver,
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sim, err := r.GetByIccid(iccid)
		assert.Error(t, err)
		assert.Nil(t, sim)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})
}

func Test_Get(t *testing.T) {
	t.Run("GetPhysicalSim", func(t *testing.T) {
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})
		rows.AddRow(
			TestIccid1,
			TestMsisdn1,
			false,
			false,
			ukama.SimTypeTest,
			TestSmDpAddress1,
			TestActivationCode1,
			true,
			TestQrCode1,
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(false, true, ukama.SimTypeTest, sqlmock.AnyArg()).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  TestDbDSN,
			DriverName:           TestDbDriver,
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sp, err := r.Get(true, ukama.SimTypeTest)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sp)
		assert.Equal(t, TestIccid1, sp.Iccid)
		assert.True(t, sp.IsPhysical)
	})

	t.Run("GetVirtualSim", func(t *testing.T) {
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})
		rows.AddRow(
			TestIccid2,
			TestMsisdn2,
			false,
			false,
			ukama.SimTypeOperatorData,
			TestSmDpAddress2,
			TestActivationCode2,
			false,
			TestQrCode2,
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(false, false, ukama.SimTypeOperatorData, sqlmock.AnyArg()).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sp, err := r.Get(false, ukama.SimTypeOperatorData)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sp)
		assert.Equal(t, TestIccid2, sp.Iccid)
		assert.False(t, sp.IsPhysical)
	})

	t.Run("GetDatabaseError", func(t *testing.T) {
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(false, true, ukama.SimTypeTest, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrInvalidDB)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sp, err := r.Get(true, ukama.SimTypeTest)
		assert.Error(t, err)
		assert.Nil(t, sp)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetRecordNotFound", func(t *testing.T) {
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(false, false, ukama.SimTypeUkamaData, sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sp, err := r.Get(false, ukama.SimTypeUkamaData)
		assert.Error(t, err)
		assert.Nil(t, sp)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_Add(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		var db *extsql.DB

		sim := []Sim{
			{
				Iccid:          TestIccid1,
				Msisdn:         TestMsisdn1,
				IsAllocated:    false,
				IsFailed:       false,
				SimType:        ukama.SimTypeTest,
				SmDpAddress:    TestSmDpAddress1,
				ActivationCode: TestActivationCode1,
				IsPhysical:     true,
				QrCode:         TestQrCode1,
			},
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sim[0].Iccid, sim[0].Msisdn, sim[0].IsAllocated, sim[0].IsFailed, sim[0].SimType, sim[0].SmDpAddress, sim[0].ActivationCode, sim[0].QrCode, sim[0].IsPhysical).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		err = r.Add(sim)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Add_MultipleSims", func(t *testing.T) {
		var db *extsql.DB

		sims := []Sim{
			{
				Iccid:          TestIccid1,
				Msisdn:         TestMsisdn1,
				IsAllocated:    false,
				IsFailed:       false,
				SimType:        ukama.SimTypeTest,
				SmDpAddress:    TestSmDpAddress1,
				ActivationCode: TestActivationCode1,
				IsPhysical:     true,
				QrCode:         TestQrCode1,
			},
			{
				Iccid:          TestIccid2,
				Msisdn:         TestMsisdn2,
				IsAllocated:    false,
				IsFailed:       false,
				SimType:        ukama.SimTypeOperatorData,
				SmDpAddress:    TestSmDpAddress2,
				ActivationCode: TestActivationCode2,
				IsPhysical:     false,
				QrCode:         TestQrCode2,
			},
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sims[0].Iccid, sims[0].Msisdn, sims[0].IsAllocated, sims[0].IsFailed, sims[0].SimType, sims[0].SmDpAddress, sims[0].ActivationCode, sims[0].QrCode, sims[0].IsPhysical,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sims[1].Iccid, sims[1].Msisdn, sims[1].IsAllocated, sims[1].IsFailed, sims[1].SimType, sims[1].SmDpAddress, sims[1].ActivationCode, sims[1].QrCode, sims[1].IsPhysical).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.Add(sims)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Add_DatabaseError", func(t *testing.T) {
		var db *extsql.DB

		sim := []Sim{
			{
				Iccid:          TestIccid1,
				Msisdn:         TestMsisdn1,
				IsAllocated:    false,
				IsFailed:       false,
				SimType:        ukama.SimTypeTest,
				SmDpAddress:    TestSmDpAddress1,
				ActivationCode: TestActivationCode1,
				IsPhysical:     true,
				QrCode:         TestQrCode1,
			},
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sim[0].Iccid, sim[0].Msisdn, sim[0].IsAllocated, sim[0].IsFailed, sim[0].SimType, sim[0].SmDpAddress, sim[0].ActivationCode, sim[0].QrCode, sim[0].IsPhysical).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.Add(sim)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_Delete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		simId := []uint64{TestId1}
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(sqlmock.AnyArg(), simId[0]).
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

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		err = r.Delete(simId)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Delete_MultipleIds", func(t *testing.T) {
		simIds := []uint64{TestId1, TestId2, TestId3}
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(sqlmock.AnyArg(), simIds[0], simIds[1], simIds[2]).
			WillReturnResult(sqlmock.NewResult(3, 3))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.Delete(simIds)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Delete_NoRecordsFound", func(t *testing.T) {
		simId := []uint64{TestIdNotFound} // Non-existent ID
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(sqlmock.AnyArg(), simId[0]).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.Delete(simId)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Delete_DatabaseError", func(t *testing.T) {
		simId := []uint64{TestId1}
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(sqlmock.AnyArg(), simId[0]).
			WillReturnError(gorm.ErrInvalidDB)

		mock.ExpectRollback()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.Delete(simId)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_GetSims(t *testing.T) {
	t.Run("GetSimsWithSpecificType", func(t *testing.T) {
		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})
		rows.AddRow(
			TestIccid1,
			TestMsisdn1,
			false,
			false,
			ukama.SimTypeTest,
			TestSmDpAddress1,
			TestActivationCode1,
			true,
			TestQrCode1,
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(ukama.SimTypeTest).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sims, err := r.GetSims(ukama.SimTypeTest)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sims)
		assert.Len(t, sims, 1)
		assert.Equal(t, TestIccid1, sims[0].Iccid)
	})

	t.Run("GetSimsWithSpecificTypeError", func(t *testing.T) {
		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(ukama.SimTypeTest).
			WillReturnError(gorm.ErrInvalidDB)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sims, err := r.GetSims(ukama.SimTypeTest)
		assert.Error(t, err)
		assert.Nil(t, sims)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetSimsWithUnknownType", func(t *testing.T) {
		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})
		rows.AddRow(
			TestIccid1,
			TestMsisdn1,
			false,
			false,
			ukama.SimTypeTest,
			TestSmDpAddress1,
			TestActivationCode1,
			true,
			TestQrCode1,
		)
		rows.AddRow(
			TestIccid2,
			TestMsisdn2,
			true,
			false,
			ukama.SimTypeOperatorData,
			TestSmDpAddress2,
			TestActivationCode2,
			false,
			TestQrCode2,
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sims, err := r.GetSims(ukama.SimTypeUnknown)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sims)
		assert.Len(t, sims, 2)
		assert.Equal(t, TestIccid1, sims[0].Iccid)
		assert.Equal(t, TestIccid2, sims[1].Iccid)
	})

	t.Run("GetSimsWithUnknownTypeError", func(t *testing.T) {
		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WillReturnError(gorm.ErrInvalidDB)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		sims, err := r.GetSims(ukama.SimTypeUnknown)
		assert.Error(t, err)
		assert.Nil(t, sims)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}

func Test_UpdateStatus(t *testing.T) {
	t.Run("UpdateStatusSuccess", func(t *testing.T) {
		var db *extsql.DB
		iccid := TestIccid1
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(true, false, sqlmock.AnyArg(), iccid).
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

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.UpdateStatus(iccid, true, false)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateStatusDatabaseError", func(t *testing.T) {
		var db *extsql.DB
		iccid := TestIccid3
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(true, false, sqlmock.AnyArg(), iccid).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.UpdateStatus(iccid, true, false)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UpdateStatusWithEmptyIccid", func(t *testing.T) {
		var db *extsql.DB
		iccid := ""
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(true, false, sqlmock.AnyArg(), iccid).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.UpdateStatus(iccid, true, false)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}
