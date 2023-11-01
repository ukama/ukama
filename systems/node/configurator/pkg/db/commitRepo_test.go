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
	"regexp"
	"testing"

	int_db "github.com/ukama/ukama/systems/node/configurator/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func TestCommitRepo_Get(t *testing.T) {

	t.Run("CommitExist", func(t *testing.T) {
		// Arrange
		const hash = "6b0a48e3d06ae7708b2257321d17b36bd930f670"

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"hash"}).
			AddRow(hash)

		mock.ExpectQuery(`^SELECT.*commits.*`).
			WithArgs(hash).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCommitRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		commit, err := r.Get(hash)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, commit) {
			assert.Equal(t, hash, commit.Hash)
		}
	})

	t.Run("Commit Doesn't Exist", func(t *testing.T) {
		// Arrange
		const hash = "6b0a48e3d06ae7708b2257321d17b36bd930f670"

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*commits.*`).
			WithArgs(hash).
			WillReturnError(gorm.ErrRecordNotFound)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCommitRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.Get(hash)

		// Assert
		if assert.Error(t, err) {
			assert.Equal(t, true, errors.Is(gorm.ErrRecordNotFound, err))
		}

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

	t.Run("GetAll", func(t *testing.T) {
		// Arrange
		const hash0 = "6b0a48e3d06ae7708b2257321d17b36bd930f670"
		const hash1 = "6b0a48e3d06ae7708b2257321d17b36bd930f671"

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"hash"}).
			AddRow(hash0).AddRow(hash1)

		mock.ExpectQuery(`^SELECT.*commits.*`).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCommitRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		commit, err := r.GetAll()

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, commit) {
			assert.Equal(t, hash0, commit[0].Hash)
			assert.Equal(t, hash1, commit[1].Hash)
		}
	})
	t.Run("GetLatest", func(t *testing.T) {
		// Arrange
		const hash0 = "6b0a48e3d06ae7708b2257321d17b36bd930f670"

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"hash"}).
			AddRow(hash0)

		mock.ExpectQuery(`^SELECT.*commits.*`).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewCommitRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		commit, err := r.GetLatest()

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, commit) {
			assert.Equal(t, hash0, commit.Hash)
		}
	})
}

func TestCommitRepo_Add(t *testing.T) {

	t.Run("AddCommit", func(t *testing.T) {
		// Arrange

		const hash0 = "6b0a48e3d06ae7708b2257321d17b36bd930f670"

		commit := int_db.Commit{
			Hash: hash0,
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), commit.Hash).
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

		r := int_db.NewCommitRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(commit.Hash)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}
