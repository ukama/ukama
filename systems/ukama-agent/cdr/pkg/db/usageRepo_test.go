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
	"regexp"
	"testing"

	int_db "github.com/ukama/ukama/systems/ukama-agent/cdr/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var usage = int_db.Usage{
	Imsi:             "123456789012345678",
	Historical:       0,
	Usage:            0,
	LastSessionUsage: 0,
	LastSessionId:    0,
}

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

	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) (err error) {
	return u.GormDb.Transaction(func(tx *gorm.DB) error {
		d := dbOperation(tx)

		if d.Error != nil {
			return d.Error
		}

		if len(nestedFuncs) > 0 {
			for _, n := range nestedFuncs {
				if n != nil {
					nestErr := n(tx)
					if nestErr != nil {
						return nestErr
					}
				}
			}
		}

		return nil
	})
}
func TestUsageRepo_Add(t *testing.T) {

	t.Run("Add", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), usage.Imsi, usage.Historical, usage.Usage, usage.LastSessionUsage, usage.LastSessionId, usage.LastNodeId, usage.LastCDRUpdatedAt, usage.Policy).
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

		r := int_db.NewUsageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&usage)

		// Assert12
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func TestUsageRepo_Get(t *testing.T) {

	t.Run("ImsiFound", func(t *testing.T) {
		var ID uint = 1
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		urow := sqlmock.NewRows([]string{"ID", "imsi", "usage", "historical", "last_session_id", "last_session_usage", "last_node_id", "last_cdr_updated_at", "policy"}).
			AddRow(ID, usage.Imsi, usage.Historical, usage.Usage, usage.LastSessionUsage, usage.LastSessionId, usage.LastNodeId, usage.LastCDRUpdatedAt, usage.Policy)

		mock.ExpectQuery(`^SELECT.*usages.*`).
			WithArgs(usage.Imsi).
			WillReturnRows(urow)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewUsageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		tu, err := r.Get(usage.Imsi)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, tu) {
			assert.EqualValues(t, ID, tu.ID)
		}

	})

	t.Run("ImsiNotFound", func(t *testing.T) {

		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		// urow := sqlmock.NewRows([]string{"ID", "usage", "historical", "last_session_id", "last_session_usage", "last_node_id", "last_cdr_updated_at", "policy"}).
		// 	AddRow(ID, usage.Imsi, usage.Historical, usage.Usage, usage.LastSessionUsage, usage.LastSessionId, usage.LastNodeId, usage.LastCDRUpdatedAt, usage.Policy)

		mock.ExpectQuery(`^SELECT.*usages.*`).
			WithArgs(usage.Imsi).
			WillReturnError(gorm.ErrRecordNotFound)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewUsageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.Get(usage.Imsi)

		// Assert
		assert.Error(t, err)
		if assert.Error(t, err) {
			assert.Equal(t, gorm.ErrRecordNotFound, err)
		}

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})
}
