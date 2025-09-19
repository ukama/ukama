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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

// Test fixtures and constants
const (
	testEmail        = "test@ukama.com"
	testTemplateName = "test_template"
	testName         = "test_user"
	failedEmail      = "failed@ukama.com"
	retryEmail       = "retry@ukama.com"
	failedUserName   = "failed_user"
	retryUserName    = "retry_user"
)

// Test data structures
type testSetup struct {
	mock    sqlmock.Sqlmock
	repo    int_db.MailerRepo
	cleanup func()
}

// UkamaDbMock implements the database interface for testing
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

// Helper functions

// setupTestDB creates a mock database connection and repository for testing
func setupTestDB(t *testing.T) *testSetup {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	cleanup := func() {
		if err := sqlDB.Close(); err != nil {
			t.Logf("Error closing database: %v", err)
		}
	}

	repo := int_db.NewMailerRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	return &testSetup{
		mock:    mock,
		repo:    repo,
		cleanup: cleanup,
	}
}

// createTestMailing creates a Mailing struct with default test values
func createTestMailing(opts ...func(*int_db.Mailing)) *int_db.Mailing {
	now := time.Now()
	mailing := &int_db.Mailing{
		MailId:        uuid.NewV4(),
		Email:         testEmail,
		TemplateName:  testTemplateName,
		SentAt:        nil,
		Status:        ukama.MailStatusPending,
		RetryCount:    0,
		NextRetryTime: &time.Time{},
		Values: utils.JSONMap{
			"Name": testName,
		},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	// Apply any custom options
	for _, opt := range opts {
		opt(mailing)
	}

	return mailing
}

// withEmail sets a custom email for the mailing
func withEmail(email string) func(*int_db.Mailing) {
	return func(m *int_db.Mailing) {
		m.Email = email
	}
}

// withStatus sets a custom status for the mailing
func withStatus(status ukama.MailStatus) func(*int_db.Mailing) {
	return func(m *int_db.Mailing) {
		m.Status = status
	}
}

// withRetryCount sets a custom retry count for the mailing
func withRetryCount(count int) func(*int_db.Mailing) {
	return func(m *int_db.Mailing) {
		m.RetryCount = count
	}
}

// withNextRetryTime sets a custom next retry time for the mailing
func withNextRetryTime(nextRetryTime time.Time) func(*int_db.Mailing) {
	return func(m *int_db.Mailing) {
		m.NextRetryTime = &nextRetryTime
	}
}

// withName sets a custom name in the values map
func withName(name string) func(*int_db.Mailing) {
	return func(m *int_db.Mailing) {
		if m.Values == nil {
			m.Values = make(utils.JSONMap)
		}
		m.Values["Name"] = name
	}
}

// createMockRows creates sqlmock.Rows for a mailing
func createMockRows(mailing *int_db.Mailing) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"mail_id", "email", "template_name", "sent_at", "status", "retry_count",
		"next_retry_time", "values", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		mailing.MailId, mailing.Email, mailing.TemplateName, mailing.SentAt,
		mailing.Status, mailing.RetryCount, mailing.NextRetryTime, mailing.Values,
		mailing.CreatedAt, mailing.UpdatedAt, mailing.DeletedAt,
	)
}

// createMultipleMockRows creates sqlmock.Rows for multiple mailings
func createMultipleMockRows(mailings ...*int_db.Mailing) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"mail_id", "email", "template_name", "sent_at", "status", "retry_count",
		"next_retry_time", "values", "created_at", "updated_at", "deleted_at",
	})

	for _, mailing := range mailings {
		rows.AddRow(
			mailing.MailId, mailing.Email, mailing.TemplateName, mailing.SentAt,
			mailing.Status, mailing.RetryCount, mailing.NextRetryTime, mailing.Values,
			mailing.CreatedAt, mailing.UpdatedAt, mailing.DeletedAt,
		)
	}

	return rows
}

// Test functions

func Test_SendEmail(t *testing.T) {
	setup := setupTestDB(t)
	defer setup.cleanup()

	email := createTestMailing()

	setup.mock.ExpectBegin()
	setup.mock.ExpectExec("INSERT INTO \"mailings\"").WithArgs(
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
	setup.mock.ExpectCommit()

	err := setup.repo.CreateEmail(email)
	assert.NoError(t, err)
	assert.NoError(t, setup.mock.ExpectationsWereMet())
}

func Test_GetEmailById(t *testing.T) {
	setup := setupTestDB(t)
	defer setup.cleanup()

	email := createTestMailing()

	rows := createMockRows(email)

	setup.mock.ExpectQuery(`SELECT \* FROM "mailings" WHERE mail_id = \$1`).
		WithArgs(email.MailId, sqlmock.AnyArg()).
		WillReturnRows(rows)

	result, err := setup.repo.GetEmailById(email.MailId)

	assert.NotNil(t, result)
	assert.NoError(t, err)
	assert.NoError(t, setup.mock.ExpectationsWereMet())
}

func Test_UpdateEmailStatus(t *testing.T) {
	setup := setupTestDB(t)
	defer setup.cleanup()

	nextRetryTime := time.Now().Add(5 * time.Minute)
	mailing := createTestMailing(
		withStatus(ukama.MailStatusFailed),
		withRetryCount(2),
		withNextRetryTime(nextRetryTime),
	)

	setup.mock.ExpectBegin()
	setup.mock.ExpectExec(`UPDATE "mailings" SET "next_retry_time"=\$1,"retry_count"=\$2,"status"=\$3,"updated_at"=\$4 WHERE mail_id = \$5`).
		WithArgs(
			mailing.NextRetryTime,
			mailing.RetryCount,
			mailing.Status,
			sqlmock.AnyArg(),
			mailing.MailId,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	setup.mock.ExpectCommit()

	err := setup.repo.UpdateEmailStatus(mailing)
	assert.NoError(t, err)
	assert.NoError(t, setup.mock.ExpectationsWereMet())
}

func Test_UpdateEmailStatus_DatabaseError(t *testing.T) {
	setup := setupTestDB(t)
	defer setup.cleanup()

	nextRetryTime := time.Now().Add(5 * time.Minute)
	mailing := createTestMailing(
		withStatus(ukama.MailStatusFailed),
		withRetryCount(2),
		withNextRetryTime(nextRetryTime),
	)

	// Mock database error
	setup.mock.ExpectBegin()
	setup.mock.ExpectExec(`UPDATE "mailings" SET "next_retry_time"=\$1,"retry_count"=\$2,"status"=\$3,"updated_at"=\$4 WHERE mail_id = \$5`).
		WithArgs(
			mailing.NextRetryTime,
			mailing.RetryCount,
			mailing.Status,
			sqlmock.AnyArg(),
			mailing.MailId,
		).
		WillReturnError(gorm.ErrInvalidDB)
	setup.mock.ExpectRollback()

	err := setup.repo.UpdateEmailStatus(mailing)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidDB, err)
	assert.NoError(t, setup.mock.ExpectationsWereMet())
}

func Test_UpdateEmailStatus_RecordNotFound(t *testing.T) {
	setup := setupTestDB(t)
	defer setup.cleanup()

	nextRetryTime := time.Now().Add(5 * time.Minute)
	mailing := createTestMailing(
		withStatus(ukama.MailStatusFailed),
		withRetryCount(2),
		withNextRetryTime(nextRetryTime),
	)

	setup.mock.ExpectBegin()
	setup.mock.ExpectExec(`UPDATE "mailings" SET "next_retry_time"=\$1,"retry_count"=\$2,"status"=\$3,"updated_at"=\$4 WHERE mail_id = \$5`).
		WithArgs(
			mailing.NextRetryTime,
			mailing.RetryCount,
			mailing.Status,
			sqlmock.AnyArg(),
			mailing.MailId,
		).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
	setup.mock.ExpectCommit()

	err := setup.repo.UpdateEmailStatus(mailing)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, setup.mock.ExpectationsWereMet())
}

func Test_GetFailedEmails(t *testing.T) {
	setup := setupTestDB(t)
	defer setup.cleanup()

	// Create test data with different scenarios
	nextRetryTime1 := time.Now().Add(-5 * time.Minute) // Past time
	nextRetryTime2 := time.Now().Add(5 * time.Minute)  // Future time

	email1 := createTestMailing(
		withEmail(failedEmail),
		withStatus(ukama.MailStatusFailed),
		withRetryCount(1),
		withNextRetryTime(nextRetryTime1),
		withName(failedUserName),
	)

	email2 := createTestMailing(
		withEmail(retryEmail),
		withStatus(ukama.MailStatusRetry),
		withRetryCount(2),
		withNextRetryTime(nextRetryTime2),
		withName(retryUserName),
	)

	rows := createMultipleMockRows(email1, email2)

	// Mock the SELECT query with the complex WHERE clause
	setup.mock.ExpectQuery(`SELECT \* FROM "mailings" WHERE status IN \(\$1, \$2\) AND retry_count < \$3 AND \(next_retry_time <= \$4 OR next_retry_time IS NULL\)`).
		WithArgs(ukama.MailStatusFailed, ukama.MailStatusRetry, ukama.MaxRetryCount, sqlmock.AnyArg()).
		WillReturnRows(rows)

	result, err := setup.repo.GetFailedEmails()

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

	assert.NoError(t, setup.mock.ExpectationsWereMet())
}
