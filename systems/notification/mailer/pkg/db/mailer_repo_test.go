/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"log"
	"testing"
	"time"

	"github.com/tj/assert"
	int_db "github.com/ukama/ukama/systems/notification/mailer/pkg/db"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/utils"

	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func Test_SendEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	repo := int_db.NewMailerRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	email := int_db.Mailing{
		MailId:        uuid.NewV4(),
		Email:         "brackley@ukama.com",
		TemplateName:  "test_template",
		SentAt:        nil,
		Status:        ukama.Pending,
		RetryCount:    0,
		NextRetryTime: &time.Time{},
		Values: utils.JSONMap{
			"Name": "joe",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"mailings\"").WithArgs(
		email.MailId,
		email.Email,
		email.TemplateName,
		email.SentAt,
		email.Status,
		email.RetryCount,
		email.NextRetryTime,
		email.Values,
		email.CreatedAt,
		email.UpdatedAt,
		email.DeletedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = repo.CreateEmail(&email)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func Test_GetEmailById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	repo := int_db.NewMailerRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	email := int_db.Mailing{
		MailId:        uuid.NewV4(),
		Email:         "brackley@ukama.com",
		TemplateName:  "test_template",
		SentAt:        nil,
		Status:        ukama.Pending,
		RetryCount:    0,
		NextRetryTime: &time.Time{},
		Values: utils.JSONMap{
			"Name": "joe",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	rows := sqlmock.NewRows([]string{"mail_id", "email", "template_name", "sent_at", "status", "retry_count", "next_retry_time", "values", "created_at", "updated_at", "deleted_at"}).
		AddRow(email.MailId, email.Email, email.TemplateName, email.SentAt, email.Status, email.RetryCount, email.NextRetryTime, email.Values, email.CreatedAt, email.UpdatedAt, email.DeletedAt)

	mock.ExpectQuery(`SELECT \* FROM "mailings" WHERE mail_id = \$1`).WithArgs(email.MailId, sqlmock.AnyArg()).WillReturnRows(rows)

	result, err := repo.GetEmailById(email.MailId)

	assert.NotNil(t, result)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
