/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"context"
	extsql "database/sql"
	"regexp"
	"testing"

	"github.com/jackc/pgtype"

	"github.com/ukama/ukama/systems/common/uuid"
	int_db "github.com/ukama/ukama/systems/init/lookup/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_orgRepo_Get(t *testing.T) {

	t.Run("OrgExist", func(t *testing.T) {
		// Arrange
		const orgName = "ukama"
		const orgCert = "ukamacert"
		const ip = "0.0.0.0"

		var dIp pgtype.Inet
		err := dIp.Set(ip)
		assert.NoError(t, err)

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"name", "certificate", "ip"}).
			AddRow(orgName, orgCert, dIp)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(orgName, sqlmock.AnyArg()).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewOrgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		org, err := r.GetByName(context.Background(), orgName)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, org)
	})

}

func Test_orgRepo_GetAll(t *testing.T) {

	t.Run("OrgsExist", func(t *testing.T) {
		// Arrange
		const orgName = "ukama"
		const orgCert = "ukamacert"
		const ip = "0.0.0.0"

		var dIp pgtype.Inet
		err := dIp.Set(ip)
		assert.NoError(t, err)

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"name", "certificate", "ip"}).
			AddRow(orgName, orgCert, dIp).
			AddRow("other-org", orgCert, dIp)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewOrgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		orgs, err := r.GetAll(context.Background())

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Len(t, orgs, 2)
	})

}

func Test_orgRepo_GetById(t *testing.T) {

	t.Run("OrgExist", func(t *testing.T) {
		// Arrange
		const orgName = "ukama"
		const orgCert = "ukamacert"
		const ip = "0.0.0.0"
		orgId := uuid.NewV4()

		var dIp pgtype.Inet
		err := dIp.Set(ip)
		assert.NoError(t, err)

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"name", "certificate", "ip", "org_id"}).
			AddRow(orgName, orgCert, dIp, orgId)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(orgId, sqlmock.AnyArg()).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewOrgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		org, err := r.GetById(context.Background(), orgId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, org)
	})

}

func Test_orgRepo_Add(t *testing.T) {

	t.Run("AddOrg", func(t *testing.T) {
		// Arrange
		const ip = "0.0.0.0"

		var dIp pgtype.Inet
		err := dIp.Set(ip)
		assert.NoError(t, err)

		org := int_db.Org{
			Name:        "ukama",
			Certificate: "ukama_certs",
			Ip:          dIp,
			OrgId:       uuid.NewV4(),
		}

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), org.Name, org.OrgId, org.Certificate, org.Ip).
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

		r := int_db.NewOrgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(context.Background(), &org)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}

func Test_orgRepo_Update(t *testing.T) {

	t.Run("UpdateOrg", func(t *testing.T) {
		// Arrange
		const ip = "0.0.0.0"

		var dIp pgtype.Inet
		err := dIp.Set(ip)
		assert.NoError(t, err)

		org := int_db.Org{
			Name:        "ukama",
			Certificate: "ukama_certs",
			Ip:          dIp,
		}

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), org.Name, org.Certificate, org.Ip, org.Name).
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

		r := int_db.NewOrgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Update(context.Background(), &org)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}
