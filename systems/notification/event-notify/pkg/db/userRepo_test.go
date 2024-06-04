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

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/roles"
	"github.com/ukama/ukama/systems/common/uuid"
	int_db "github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var user = int_db.Users{
	Id:           uuid.NewV4(),
	OrgId:        testOrgId,
	UserId:       testUserId,
	SubscriberId: uuid.NewV4().String(),
	Role:         roles.TYPE_OWNER,
}

func TestUserRepo_Add(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		em := GetEventMsg()
		assert.NotNil(t, em)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(user.Id, user.OrgId, user.NetworkId, user.SubscriberId, user.UserId, user.Role, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

		r := int_db.NewUserRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&user)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUserRepo_GetUser(t *testing.T) {
	t.Run("GetUser", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "org_id", "network_id", "subscriber_id", "user_id", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(user.Id, user.OrgId, user.NetworkId, user.SubscriberId, user.UserId, user.Role, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(user.UserId).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewUserRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		ru, err := r.GetUser(user.UserId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, ru)
		assert.Equal(t, user.UserId, ru.UserId)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAllUser", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "org_id", "network_id", "subscriber_id", "user_id", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(user.Id, user.OrgId, user.NetworkId, user.SubscriberId, user.UserId, user.Role, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(user.OrgId).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewUserRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		ruL, err := r.GetAllUsers(user.OrgId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, ruL)
		assert.Equal(t, user.UserId, ruL[0].UserId)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetUserWithRoles", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "org_id", "network_id", "subscriber_id", "user_id", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(user.Id, user.OrgId, user.NetworkId, user.SubscriberId, user.UserId, user.Role, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(user.OrgId, user.Role).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewUserRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		ruL, err := r.GetUserWithRoles(user.OrgId, []roles.RoleType{user.Role})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, ruL)
		assert.Equal(t, user.UserId, ruL[0].UserId)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetUsers", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "org_id", "network_id", "subscriber_id", "user_id", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(user.Id, user.OrgId, user.NetworkId, user.SubscriberId, user.UserId, user.Role, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(user.OrgId, user.SubscriberId, user.UserId).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewUserRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		ruL, err := r.GetUsers(user.OrgId, user.NetworkId, user.SubscriberId, user.UserId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, ruL)
		assert.Equal(t, user.UserId, ruL[0].UserId)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUserRepo_GetSubscriber(t *testing.T) {
	t.Run("GetSubscriber", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "org_id", "network_id", "subscriber_id", "user_id", "role", "created_at", "updated_at", "deleted_at"}).
			AddRow(user.Id, user.OrgId, user.NetworkId, user.SubscriberId, user.UserId, user.Role, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WithArgs(user.SubscriberId).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewUserRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		ru, err := r.GetSubscriber(user.SubscriberId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, ru)
		assert.Equal(t, user.UserId, ru.UserId)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
