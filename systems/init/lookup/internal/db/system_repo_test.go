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

	"github.com/jackc/pgtype"
	uuid "github.com/satori/go.uuid"
	int_db "github.com/ukama/ukama/systems/init/lookup/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_systemRepo_Get(t *testing.T) {

	t.Run("SystemExist", func(t *testing.T) {
		// Arrange
		const name = "sys"
		const uuidStr = "51fbba62-c79f-11eb-b8bc-0242ac130003"
		const ip = "0.0.0.0"
		const certs = "ukama_certs"
		const port = 101
		const health = 100
		const orgId = uint(15)
		var dIp pgtype.Inet
		err := dIp.Set(ip)
		assert.NoError(t, err)

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"name", "uuid", "certificate", "ip", "port", "health"}).
			AddRow(name, uuidStr, certs, dIp, port, health)

		mock.ExpectQuery(`^SELECT.*systems.*`).
			WithArgs(name, orgId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewSystemRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		node, err := r.GetByName(name, orgId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, node)
	})

}

func Test_systemRepo_Delete(t *testing.T) {

	t.Run("DeleteSystem", func(t *testing.T) {

		const name = "sys"
		const orgId = uint(15)
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("DELETE")).WithArgs(name, orgId).
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

		r := int_db.NewSystemRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Delete(name, orgId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}

func Test_systemRepo_Add(t *testing.T) {
	t.Run("Add_WithExistingSystemNameandOrgID", func(t *testing.T) {
		// Arrange
		const ip = "0.0.0.0"
		const orgId = uint(15)

		var dIp pgtype.Inet
		err := dIp.Set(ip)
		assert.NoError(t, err)

		system := int_db.System{
			Name:        "sys",
			Certificate: "sys_certs",
			Ip:          dIp,
			Port:        100,
			Uuid:        uuid.NewV4().String(),
			OrgID:       orgId,
			Health:      100,
		}

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), system.Name, system.Uuid, system.Certificate, system.Ip, system.Port, system.OrgID, system.Health, system.OrgID, system.Name).
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

		r := int_db.NewSystemRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&system)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}

func Test_systemRepo_Update(t *testing.T) {

	t.Run("Update", func(t *testing.T) {
		// Arrange
		const ip = "0.0.0.0"
		const orgId = uint(15)

		var dIp pgtype.Inet
		err := dIp.Set(ip)
		assert.NoError(t, err)

		system := int_db.System{
			Name:        "sys",
			Certificate: "sys_certs",
			Ip:          dIp,
			Port:        100,
			Uuid:        uuid.NewV4().String(),
			OrgID:       orgId,
			Health:      100,
		}

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), system.Name, system.Uuid, system.Certificate, system.Ip, system.Port, system.OrgID, system.Health, system.Name, system.OrgID).
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

		r := int_db.NewSystemRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Update(&system, orgId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}
