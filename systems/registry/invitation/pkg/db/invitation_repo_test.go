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
	 "time"
 
	 "github.com/ukama/ukama/systems/common/roles"
	 "github.com/ukama/ukama/systems/common/uuid"
	 db_inv "github.com/ukama/ukama/systems/registry/invitation/pkg/db"
 
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
 
 func TestInvitationRepo_AddInvitation(t *testing.T) {
	 invitation := db_inv.Invitation{
		 Id:        uuid.NewV4(),
		 Name:      "test",
		 Email:     "test@ukama.com",
		 Role:      roles.TYPE_ADMIN,
		 Status:    db_inv.Pending,
		 UserId:    uuid.NewV4().String(),
		 ExpiresAt: time.Date(2023, 8, 25, 17, 59, 43, 831000000, time.UTC),
		 Link:      "https://ukama.com/invitation/accept/" + uuid.NewV4().String(),
		 CreatedAt: time.Now(),
		 UpdatedAt: time.Now(),
		 DeletedAt: gorm.DeletedAt{},
	 }
 
	 db, mock, err := sqlmock.New()
	 assert.NoError(t, err)
 
	 dialector := postgres.New(postgres.Config{
		 DSN:                  "sqlmock_db_0",
		 DriverName:           "postgres",
		 Conn:                 db,
		 PreferSimpleProtocol: true,
	 })
 
	 gdb, err := gorm.Open(dialector, &gorm.Config{})
	 assert.NoError(t, err)
 
	 r := db_inv.NewInvitationRepo(&UkamaDbMock{
		 GormDb: gdb,
	 })
 
	 t.Run("AddValidInvitation", func(t *testing.T) {
		 mock.ExpectBegin()
 
		 mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			 WithArgs(invitation.Id, invitation.Link, invitation.Email, invitation.Name, invitation.ExpiresAt, invitation.Role, invitation.Status,
				 invitation.UserId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			 WillReturnResult(sqlmock.NewResult(1, 1))
 
		 mock.ExpectCommit()
 
		 // Act
		 err = r.Add(&invitation, nil)
 
		 // Assert
		 assert.NoError(t, err)
 
		 err = mock.ExpectationsWereMet()
		 assert.NoError(t, err)
	 })
 }
 func TestInvitationRepo_Getinvitation(t *testing.T) {
	 t.Run("InvitationExist", func(t *testing.T) {
		 // Arrange
		 invId := uuid.NewV4()
		 invitation := db_inv.Invitation{
			 Id:        invId,
			 Name:      "test",
			 Email:     "test@ukama.com",
			 Role:      roles.TYPE_ADMIN,
			 Status:    db_inv.Pending,
			 UserId:    uuid.NewV4().String(),
			 ExpiresAt: time.Date(2023, 8, 25, 17, 59, 43, 831000000, time.UTC),
			 Link:      "https://ukama.com/invitation/accept/" + uuid.NewV4().String(),
			 CreatedAt: time.Now(),
			 UpdatedAt: time.Now(),
			 DeletedAt: gorm.DeletedAt{},
		 }
 
		 var db *extsql.DB
 
		 db, mock, err := sqlmock.New() // mock sql.DB
		 assert.NoError(t, err)
 
		 rows := sqlmock.NewRows([]string{"id", "user_id", "role"}).
			 AddRow(invId, invitation.UserId, invitation.Role)
 
		 mock.ExpectQuery(`^SELECT.*invitations.*`).
			 WithArgs(invitation.Id, sqlmock.AnyArg()).
			 WillReturnRows(rows)
 
		 dialector := postgres.New(postgres.Config{
			 DSN:                  "sqlmock_db_0",
			 DriverName:           "postgres",
			 Conn:                 db,
			 PreferSimpleProtocol: true,
		 })
 
		 gdb, err := gorm.Open(dialector, &gorm.Config{})
		 assert.NoError(t, err)
 
		 r := db_inv.NewInvitationRepo(&UkamaDbMock{
			 GormDb: gdb,
		 })
 
		 assert.NoError(t, err)
 
		 // Act
		 rm, err := r.Get(invitation.Id)
 
		 // Assert
		 assert.NoError(t, err)
 
		 err = mock.ExpectationsWereMet()
		 assert.NoError(t, err)
		 assert.NotNil(t, rm)
	 })
 }
 
 func TestInvitationRepo_GetByOrg(t *testing.T) {
	 t.Run("InvitationExist", func(t *testing.T) {
		 // Arrange
		 invitation := db_inv.Invitation{
			 Id:        uuid.NewV4(),
			 Name:      "test",
			 Email:     "test@ukama.com",
			 Role:      roles.TYPE_ADMIN,
			 Status:    db_inv.Pending,
			 UserId:    uuid.NewV4().String(),
			 ExpiresAt: time.Date(2023, 8, 25, 17, 59, 43, 831000000, time.UTC),
			 Link:      "https://ukama.com/invitation/accept/" + uuid.NewV4().String(),
			 CreatedAt: time.Now(),
			 UpdatedAt: time.Now(),
			 DeletedAt: gorm.DeletedAt{},
		 }
 
		 var db *extsql.DB
 
		 db, mock, err := sqlmock.New() // mock sql.DB
		 assert.NoError(t, err)
 
		 rows := sqlmock.NewRows([]string{"id", "user_id", "role"}).
			 AddRow(invitation.Id, invitation.UserId, invitation.Role)
 
		 mock.ExpectQuery(`^SELECT.*invitations.*`).
			 WithArgs().
			 WillReturnRows(rows)
 
		 dialector := postgres.New(postgres.Config{
			 DSN:                  "sqlmock_db_0",
			 DriverName:           "postgres",
			 Conn:                 db,
			 PreferSimpleProtocol: true,
		 })
 
		 gdb, err := gorm.Open(dialector, &gorm.Config{})
		 assert.NoError(t, err)
 
		 r := db_inv.NewInvitationRepo(&UkamaDbMock{
			 GormDb: gdb,
		 })
 
		 assert.NoError(t, err)
 
		 // Act
		 rm, err := r.GetAll()
 
		 // Assert
		 assert.NoError(t, err)
 
		 err = mock.ExpectationsWereMet()
		 assert.NoError(t, err)
		 assert.NotNil(t, rm)
	 })
 }
 
 func TestInvitationRepo_Delete(t *testing.T) {
	 t.Run("DeleteInvitation", func(t *testing.T) {
		 // Arrange
		 invitation := db_inv.Invitation{
			 Id:        uuid.NewV4(),
			 Name:      "test",
			 Email:     "test@ukama",
			 Role:      roles.TYPE_ADMIN,
			 Status:    db_inv.Pending,
			 UserId:    uuid.NewV4().String(),
			 ExpiresAt: time.Date(2023, 8, 25, 17, 59, 43, 831000000, time.UTC),
			 Link:      "https://ukama.com/invitation/accept/" + uuid.NewV4().String(),
			 CreatedAt: time.Now(),
			 UpdatedAt: time.Now(),
			 DeletedAt: gorm.DeletedAt{},
		 }
 
		 var db *extsql.DB
 
		 db, mock, err := sqlmock.New() // mock sql.DB
		 assert.NoError(t, err)
 
		 mock.ExpectBegin()
 
		 mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invitations"`)).
			 WithArgs(sqlmock.AnyArg(), invitation.Id).
			 WillReturnResult(sqlmock.NewResult(1, 1))
 
		 mock.ExpectCommit()
 
		 dialector := postgres.New(postgres.Config{
			 DSN:        "sqlmock_db_0",
			 DriverName: "postgres",
			 Conn:       db,
 
			 PreferSimpleProtocol: true,
		 })
 
		 gdb, err := gorm.Open(dialector, &gorm.Config{})
 
		 assert.NoError(t, err)
 
		 r := db_inv.NewInvitationRepo(&UkamaDbMock{
			 GormDb: gdb,
		 })
 
		 assert.NoError(t, err)
 
		 // Act
		 err = r.Delete(invitation.Id, nil)
 
		 // Assert
 
		 assert.NoError(t, err)
 
		 err = mock.ExpectationsWereMet()
		 assert.NoError(t, err)
	 })
 }
 