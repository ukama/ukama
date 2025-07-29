/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"testing"
	"time"

	"github.com/tj/assert"
	int_db "github.com/ukama/ukama/systems/notification/mailer/pkg/db"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/utils"

	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"

	log "github.com/sirupsen/logrus"
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
		Status:        ukama.MailStatusPending,
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
		Status:        ukama.MailStatusPending,
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

func Test_UpdateEmailStatus(t *testing.T) {
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

	mailId := uuid.NewV4()
	nextRetryTime := time.Now().Add(5 * time.Minute)

	mailing := int_db.Mailing{
		MailId:        mailId,
		Email:         "test@ukama.com",
		TemplateName:  "test_template",
		SentAt:        nil,
		Status:        ukama.Failed,
		RetryCount:    2,
		NextRetryTime: &nextRetryTime,
		Values: utils.JSONMap{
			"Name": "test",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "mailings" SET "next_retry_time"=\$1,"retry_count"=\$2,"status"=\$3,"updated_at"=\$4 WHERE mail_id = \$5`).
		WithArgs(
			mailing.NextRetryTime,
			mailing.RetryCount,
			mailing.Status,
			sqlmock.AnyArg(),
			mailing.MailId,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.UpdateEmailStatus(&mailing)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_UpdateEmailStatus_DatabaseError(t *testing.T) {
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

	mailId := uuid.NewV4()
	nextRetryTime := time.Now().Add(5 * time.Minute)

	mailing := int_db.Mailing{
		MailId:        mailId,
		Email:         "test@ukama.com",
		TemplateName:  "test_template",
		SentAt:        nil,
		Status:        ukama.Failed,
		RetryCount:    2,
		NextRetryTime: &nextRetryTime,
		Values: utils.JSONMap{
			"Name": "test",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	// Mock database error
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "mailings" SET "next_retry_time"=\$1,"retry_count"=\$2,"status"=\$3,"updated_at"=\$4 WHERE mail_id = \$5`).
		WithArgs(
			mailing.NextRetryTime,
			mailing.RetryCount,
			mailing.Status,
			sqlmock.AnyArg(),
			mailing.MailId,
		).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	err = repo.UpdateEmailStatus(&mailing)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidDB, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_UpdateEmailStatus_RecordNotFound(t *testing.T) {
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

	mailId := uuid.NewV4()
	nextRetryTime := time.Now().Add(5 * time.Minute)

	mailing := int_db.Mailing{
		MailId:        mailId,
		Email:         "test@ukama.com",
		TemplateName:  "test_template",
		SentAt:        nil,
		Status:        ukama.Failed,
		RetryCount:    2,
		NextRetryTime: &nextRetryTime,
		Values: utils.JSONMap{
			"Name": "test",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "mailings" SET "next_retry_time"=\$1,"retry_count"=\$2,"status"=\$3,"updated_at"=\$4 WHERE mail_id = \$5`).
		WithArgs(
			mailing.NextRetryTime,
			mailing.RetryCount,
			mailing.Status,
			sqlmock.AnyArg(),
			mailing.MailId,
		).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
	mock.ExpectCommit()

	err = repo.UpdateEmailStatus(&mailing)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_GetFailedEmails(t *testing.T) {
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

	// Create test data
	mailId1 := uuid.NewV4()
	mailId2 := uuid.NewV4()
	nextRetryTime1 := time.Now().Add(-5 * time.Minute) // Past time
	nextRetryTime2 := time.Now().Add(5 * time.Minute)  // Future time

	email1 := int_db.Mailing{
		MailId:        mailId1,
		Email:         "failed@ukama.com",
		TemplateName:  "test_template",
		SentAt:        nil,
		Status:        ukama.Failed,
		RetryCount:    1,
		NextRetryTime: &nextRetryTime1,
		Values: utils.JSONMap{
			"Name": "failed_user",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	email2 := int_db.Mailing{
		MailId:        mailId2,
		Email:         "retry@ukama.com",
		TemplateName:  "test_template",
		SentAt:        nil,
		Status:        ukama.Retry,
		RetryCount:    2,
		NextRetryTime: &nextRetryTime2,
		Values: utils.JSONMap{
			"Name": "retry_user",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	// Create rows for the mock query result
	rows := sqlmock.NewRows([]string{
		"mail_id", "email", "template_name", "sent_at", "status", "retry_count",
		"next_retry_time", "values", "created_at", "updated_at", "deleted_at",
	}).
		AddRow(
			email1.MailId, email1.Email, email1.TemplateName, email1.SentAt,
			email1.Status, email1.RetryCount, email1.NextRetryTime, email1.Values,
			email1.CreatedAt, email1.UpdatedAt, email1.DeletedAt,
		).
		AddRow(
			email2.MailId, email2.Email, email2.TemplateName, email2.SentAt,
			email2.Status, email2.RetryCount, email2.NextRetryTime, email2.Values,
			email2.CreatedAt, email2.UpdatedAt, email2.DeletedAt,
		)

	// Mock the SELECT query with the complex WHERE clause
	mock.ExpectQuery(`SELECT \* FROM "mailings" WHERE status IN \(\$1, \$2\) AND retry_count < \$3 AND \(next_retry_time <= \$4 OR next_retry_time IS NULL\)`).
		WithArgs(ukama.Failed, ukama.Retry, ukama.MaxRetryCount, sqlmock.AnyArg()).
		WillReturnRows(rows)

	result, err := repo.GetFailedEmails()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	// Verify the returned emails
	assert.Equal(t, email1.MailId, result[0].MailId)
	assert.Equal(t, email1.Status, result[0].Status)
	assert.Equal(t, email1.RetryCount, result[0].RetryCount)

	assert.Equal(t, email2.MailId, result[1].MailId)
	assert.Equal(t, email2.Status, result[1].Status)
	assert.Equal(t, email2.RetryCount, result[1].RetryCount)

	assert.NoError(t, mock.ExpectationsWereMet())
}
