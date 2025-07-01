/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/tj/assert"
	int_db "github.com/ukama/ukama/systems/subscriber/registry/pkg/db"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Test data constants
const (
	// Names
	TEST_NAME_JOHN         = "John"
	TEST_NAME_JANE         = "Jane"
	TEST_NAME_BOB          = "Bob"
	TEST_NAME_ALICE        = "Alice"
	TEST_NAME_CHARLIE      = "Charlie"
	TEST_NAME_DIANA        = "Diana"
	TEST_NAME_MINIMAL      = "Minimal"
	TEST_NAME_JOHN_UPDATED = "John Updated"

	// Emails
	TEST_EMAIL_JOHN         = "john@example.com"
	TEST_EMAIL_JANE         = "jane@example.com"
	TEST_EMAIL_BOB          = "bob@example.com"
	TEST_EMAIL_ALICE        = "alice@example.com"
	TEST_EMAIL_CHARLIE      = "charlie@example.com"
	TEST_EMAIL_DIANA        = "diana@example.com"
	TEST_EMAIL_MINIMAL      = "minimal@example.com"
	TEST_EMAIL_JOHN_UPDATED = "john.updated@example.com"
	TEST_EMAIL_NONEXISTENT  = "nonexistent@example.com"
	TEST_EMAIL_TEST         = "test@example.com"

	// Phone numbers
	TEST_PHONE_DEFAULT = "555-555-5555"
	TEST_PHONE_JANE    = "555-555-5556"
	TEST_PHONE_BOB     = "555-555-5557"
	TEST_PHONE_ALICE   = "555-555-5558"
	TEST_PHONE_CHARLIE = "555-555-5559"
	TEST_PHONE_DIANA   = "555-555-5560"
	TEST_PHONE_UPDATED = "555-555-9999"

	// Gender
	TEST_GENDER_MALE   = "Male"
	TEST_GENDER_FEMALE = "Female"

	// Dates of birth
	TEST_DOB_DEFAULT = "07-03-2023"
	TEST_DOB_JANE    = "15-06-1990"
	TEST_DOB_BOB     = "22-12-1985"
	TEST_DOB_ALICE   = "03-09-1992"
	TEST_DOB_CHARLIE = "18-04-1988"
	TEST_DOB_DIANA   = "11-07-1995"
	TEST_DOB_UPDATED = "15-08-1990"

	// Proof of identification
	TEST_PROOF_DEFAULT = "Driver's License"
	TEST_PROOF_JANE    = "Passport"
	TEST_PROOF_BOB     = "National ID"
	TEST_PROOF_ALICE   = "Birth Certificate"
	TEST_PROOF_CHARLIE = "Social Security Card"
	TEST_PROOF_DIANA   = "Voter ID"

	// ID serials
	TEST_ID_SERIAL_DEFAULT = "ABC123"
	TEST_ID_SERIAL_JANE    = "XYZ789"
	TEST_ID_SERIAL_BOB     = "DEF456"
	TEST_ID_SERIAL_ALICE   = "GHI789"
	TEST_ID_SERIAL_CHARLIE = "JKL012"
	TEST_ID_SERIAL_DIANA   = "MNO345"

	// Addresses
	TEST_ADDRESS_DEFAULT = "123 Main St."
	TEST_ADDRESS_JANE    = "456 Oak Ave."
	TEST_ADDRESS_BOB     = "789 Pine St."
	TEST_ADDRESS_ALICE   = "321 Elm St."
	TEST_ADDRESS_CHARLIE = "654 Maple Dr."
	TEST_ADDRESS_DIANA   = "987 Cedar Ln."
	TEST_ADDRESS_UPDATED = "456 Updated St."

	// Error messages
	TEST_ERROR_EXPECTED      = "expected error"
	TEST_ERROR_DB_CONNECTION = "database connection error"
	TEST_ERROR_DB_FAILED     = "database connection failed"
	TEST_ERROR_DB            = "db error"
	TEST_ERROR_IMPLEMENT_ME  = "implement me"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic(TEST_ERROR_IMPLEMENT_ME + ": Init()")
}

func (u UkamaDbMock) Connect() error {
	panic(TEST_ERROR_IMPLEMENT_ME + ": Connect()")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	log.Fatal(TEST_ERROR_IMPLEMENT_ME + ": ExecuteInTransaction()")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	log.Fatal(TEST_ERROR_IMPLEMENT_ME + ": ExecuteInTransaction2()")
	return nil
}

func setupTestDB(t *testing.T) (sqlmock.Sqlmock, int_db.SubscriberRepo) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	repo := int_db.NewSubscriberRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	return mock, repo
}

func createTestSubscriber(name, email string) int_db.Subscriber {
	return int_db.Subscriber{
		SubscriberId:          uuid.NewV4(),
		Name:                  name,
		NetworkId:             uuid.NewV4(),
		Email:                 email,
		PhoneNumber:           TEST_PHONE_DEFAULT,
		Gender:                TEST_GENDER_MALE,
		DOB:                   TEST_DOB_DEFAULT,
		ProofOfIdentification: TEST_PROOF_DEFAULT,
		IdSerial:              TEST_ID_SERIAL_DEFAULT,
		Address:               TEST_ADDRESS_DEFAULT,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	}
}

func TestSubscriber_Add(t *testing.T) {
	t.Run("SuccessWithoutNestedFunc", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		subscriber := createTestSubscriber(TEST_NAME_JOHN, TEST_EMAIL_JOHN)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO \"subscribers\"").WithArgs(
			subscriber.SubscriberId,
			subscriber.Name,
			subscriber.NetworkId,
			subscriber.Email,
			subscriber.PhoneNumber,
			subscriber.Gender,
			subscriber.DOB,
			subscriber.ProofOfIdentification,
			subscriber.IdSerial,
			subscriber.Address,
			subscriber.CreatedAt,
			subscriber.UpdatedAt,
			subscriber.DeletedAt).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Act
		err := repo.Add(&subscriber, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SuccessWithNestedFunc", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		subscriber := int_db.Subscriber{
			SubscriberId:          uuid.NewV4(),
			Name:                  TEST_NAME_JANE,
			NetworkId:             uuid.NewV4(),
			Email:                 TEST_EMAIL_JANE,
			PhoneNumber:           TEST_PHONE_JANE,
			Gender:                TEST_GENDER_FEMALE,
			DOB:                   TEST_DOB_JANE,
			ProofOfIdentification: TEST_PROOF_JANE,
			IdSerial:              TEST_ID_SERIAL_JANE,
			Address:               TEST_ADDRESS_JANE,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
			DeletedAt:             nil,
		}

		nestedFuncCalled := false
		nestedFunc := func(sub *int_db.Subscriber, tx *gorm.DB) error {
			nestedFuncCalled = true
			assert.Equal(t, subscriber.SubscriberId, sub.SubscriberId)
			assert.NotNil(t, tx)
			return nil
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO \"subscribers\"").WithArgs(
			subscriber.SubscriberId,
			subscriber.Name,
			subscriber.NetworkId,
			subscriber.Email,
			subscriber.PhoneNumber,
			subscriber.Gender,
			subscriber.DOB,
			subscriber.ProofOfIdentification,
			subscriber.IdSerial,
			subscriber.Address,
			subscriber.CreatedAt,
			subscriber.UpdatedAt,
			subscriber.DeletedAt).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Act
		err := repo.Add(&subscriber, nestedFunc)

		// Assert
		assert.NoError(t, err)
		assert.True(t, nestedFuncCalled)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NestedFuncReturnsError", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		subscriber := int_db.Subscriber{
			SubscriberId:          uuid.NewV4(),
			Name:                  TEST_NAME_BOB,
			NetworkId:             uuid.NewV4(),
			Email:                 TEST_EMAIL_BOB,
			PhoneNumber:           TEST_PHONE_BOB,
			Gender:                TEST_GENDER_MALE,
			DOB:                   TEST_DOB_BOB,
			ProofOfIdentification: TEST_PROOF_BOB,
			IdSerial:              TEST_ID_SERIAL_BOB,
			Address:               TEST_ADDRESS_BOB,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
			DeletedAt:             nil,
		}

		expectedError := errors.New(TEST_ERROR_EXPECTED)
		nestedFunc := func(sub *int_db.Subscriber, tx *gorm.DB) error {
			return expectedError
		}

		mock.ExpectBegin()
		mock.ExpectRollback()

		// Act
		err := repo.Add(&subscriber, nestedFunc)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseCreateFails", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		subscriber := int_db.Subscriber{
			SubscriberId:          uuid.NewV4(),
			Name:                  TEST_NAME_ALICE,
			NetworkId:             uuid.NewV4(),
			Email:                 TEST_EMAIL_ALICE,
			PhoneNumber:           TEST_PHONE_ALICE,
			Gender:                TEST_GENDER_FEMALE,
			DOB:                   TEST_DOB_ALICE,
			ProofOfIdentification: TEST_PROOF_ALICE,
			IdSerial:              TEST_ID_SERIAL_ALICE,
			Address:               TEST_ADDRESS_ALICE,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
			DeletedAt:             nil,
		}

		expectedError := errors.New(TEST_ERROR_EXPECTED)
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO \"subscribers\"").WithArgs(
			subscriber.SubscriberId,
			subscriber.Name,
			subscriber.NetworkId,
			subscriber.Email,
			subscriber.PhoneNumber,
			subscriber.Gender,
			subscriber.DOB,
			subscriber.ProofOfIdentification,
			subscriber.IdSerial,
			subscriber.Address,
			subscriber.CreatedAt,
			subscriber.UpdatedAt,
			subscriber.DeletedAt).WillReturnError(expectedError)
		mock.ExpectRollback()

		// Act
		err := repo.Add(&subscriber, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseTransactionFails", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		subscriber := int_db.Subscriber{
			SubscriberId:          uuid.NewV4(),
			Name:                  TEST_NAME_CHARLIE,
			NetworkId:             uuid.NewV4(),
			Email:                 TEST_EMAIL_CHARLIE,
			PhoneNumber:           TEST_PHONE_CHARLIE,
			Gender:                TEST_GENDER_MALE,
			DOB:                   TEST_DOB_CHARLIE,
			ProofOfIdentification: TEST_PROOF_CHARLIE,
			IdSerial:              TEST_ID_SERIAL_CHARLIE,
			Address:               TEST_ADDRESS_CHARLIE,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
			DeletedAt:             nil,
		}

		expectedError := errors.New(TEST_ERROR_EXPECTED)
		mock.ExpectBegin().WillReturnError(expectedError)

		// Act
		err := repo.Add(&subscriber, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NestedFuncWithDatabaseError", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		subscriber := int_db.Subscriber{
			SubscriberId:          uuid.NewV4(),
			Name:                  TEST_NAME_DIANA,
			NetworkId:             uuid.NewV4(),
			Email:                 TEST_EMAIL_DIANA,
			PhoneNumber:           TEST_PHONE_DIANA,
			Gender:                TEST_GENDER_FEMALE,
			DOB:                   TEST_DOB_DIANA,
			ProofOfIdentification: TEST_PROOF_DIANA,
			IdSerial:              TEST_ID_SERIAL_DIANA,
			Address:               TEST_ADDRESS_DIANA,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
			DeletedAt:             nil,
		}

		nestedFunc := func(sub *int_db.Subscriber, tx *gorm.DB) error {
			// Simulate a database operation that fails
			return gorm.ErrInvalidDB
		}

		mock.ExpectBegin()
		mock.ExpectRollback()

		// Act
		err := repo.Add(&subscriber, nestedFunc)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SubscriberWithMinimalFields", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		subscriber := int_db.Subscriber{
			SubscriberId: uuid.NewV4(),
			Name:         TEST_NAME_MINIMAL,
			NetworkId:    uuid.NewV4(),
			Email:        TEST_EMAIL_MINIMAL,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO \"subscribers\"").WithArgs(
			subscriber.SubscriberId,
			subscriber.Name,
			subscriber.NetworkId,
			subscriber.Email,
			subscriber.PhoneNumber,
			subscriber.Gender,
			subscriber.DOB,
			subscriber.ProofOfIdentification,
			subscriber.IdSerial,
			subscriber.Address,
			subscriber.CreatedAt,
			subscriber.UpdatedAt,
			subscriber.DeletedAt).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Act
		err := repo.Add(&subscriber, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriber_Get(t *testing.T) {

	t.Run("SubscriberFound", func(t *testing.T) {
		var subID = uuid.NewV4()

		// Arrange
		mock, repo := setupTestDB(t)

		subRow := sqlmock.NewRows([]string{"subscriber_id"}).
			AddRow(subID)

		mock.ExpectQuery(`^SELECT.*subscribers.*`).
			WithArgs(subID, sqlmock.AnyArg()).
			WillReturnRows(subRow)

		// Act
		sub, err := repo.Get(subID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sub)
	})
	t.Run("SubscriberNotFound", func(t *testing.T) {
		// Arrange
		var subID = uuid.NewV4()

		// Arrange
		mock, repo := setupTestDB(t)

		mock.ExpectQuery(`^SELECT.*subscribers.*`).
			WithArgs(subID, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		// Act
		sub, err := repo.Get(subID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, sub)
	})
}

func TestSubscriber_GetByNetwork(t *testing.T) {
	t.Run("SuccessWithMultipleSubscribers", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		networkId := uuid.NewV4()
		subscriber1 := createTestSubscriber(TEST_NAME_JOHN, TEST_EMAIL_JOHN)
		subscriber1.NetworkId = networkId
		subscriber2 := createTestSubscriber(TEST_NAME_JANE, TEST_EMAIL_JANE)
		subscriber2.NetworkId = networkId
		subscriber3 := createTestSubscriber(TEST_NAME_BOB, TEST_EMAIL_BOB)
		subscriber3.NetworkId = networkId

		rows := sqlmock.NewRows([]string{
			"subscriber_id", "name", "network_id", "email", "phone_number",
			"gender", "dob", "proof_of_identification", "id_serial", "address",
			"created_at", "updated_at", "deleted_at",
		}).AddRow(
			subscriber1.SubscriberId, subscriber1.Name, subscriber1.NetworkId,
			subscriber1.Email, subscriber1.PhoneNumber, subscriber1.Gender,
			subscriber1.DOB, subscriber1.ProofOfIdentification, subscriber1.IdSerial,
			subscriber1.Address, subscriber1.CreatedAt, subscriber1.UpdatedAt, subscriber1.DeletedAt,
		).AddRow(
			subscriber2.SubscriberId, subscriber2.Name, subscriber2.NetworkId,
			subscriber2.Email, subscriber2.PhoneNumber, subscriber2.Gender,
			subscriber2.DOB, subscriber2.ProofOfIdentification, subscriber2.IdSerial,
			subscriber2.Address, subscriber2.CreatedAt, subscriber2.UpdatedAt, subscriber2.DeletedAt,
		).AddRow(
			subscriber3.SubscriberId, subscriber3.Name, subscriber3.NetworkId,
			subscriber3.Email, subscriber3.PhoneNumber, subscriber3.Gender,
			subscriber3.DOB, subscriber3.ProofOfIdentification, subscriber3.IdSerial,
			subscriber3.Address, subscriber3.CreatedAt, subscriber3.UpdatedAt, subscriber3.DeletedAt,
		)

		mock.ExpectQuery(`^SELECT \* FROM "subscribers" WHERE "subscribers"."network_id" = \$1`).
			WithArgs(networkId).
			WillReturnRows(rows)

		// Act
		subscribers, err := repo.GetByNetwork(networkId)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, subscribers, 3)
		assert.Equal(t, subscriber1.SubscriberId, subscribers[0].SubscriberId)
		assert.Equal(t, subscriber1.Name, subscribers[0].Name)
		assert.Equal(t, subscriber1.Email, subscribers[0].Email)
		assert.Equal(t, networkId, subscribers[0].NetworkId)
		assert.Equal(t, subscriber2.SubscriberId, subscribers[1].SubscriberId)
		assert.Equal(t, subscriber2.Name, subscribers[1].Name)
		assert.Equal(t, subscriber2.Email, subscribers[1].Email)
		assert.Equal(t, networkId, subscribers[1].NetworkId)
		assert.Equal(t, subscriber3.SubscriberId, subscribers[2].SubscriberId)
		assert.Equal(t, subscriber3.Name, subscribers[2].Name)
		assert.Equal(t, subscriber3.Email, subscribers[2].Email)
		assert.Equal(t, networkId, subscribers[2].NetworkId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SuccessWithSingleSubscriber", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		networkId := uuid.NewV4()
		subscriber := createTestSubscriber(TEST_NAME_ALICE, TEST_EMAIL_ALICE)
		subscriber.NetworkId = networkId

		rows := sqlmock.NewRows([]string{
			"subscriber_id", "name", "network_id", "email", "phone_number",
			"gender", "dob", "proof_of_identification", "id_serial", "address",
			"created_at", "updated_at", "deleted_at",
		}).AddRow(
			subscriber.SubscriberId, subscriber.Name, subscriber.NetworkId,
			subscriber.Email, subscriber.PhoneNumber, subscriber.Gender,
			subscriber.DOB, subscriber.ProofOfIdentification, subscriber.IdSerial,
			subscriber.Address, subscriber.CreatedAt, subscriber.UpdatedAt, subscriber.DeletedAt,
		)

		mock.ExpectQuery(`^SELECT \* FROM "subscribers" WHERE "subscribers"."network_id" = \$1`).
			WithArgs(networkId).
			WillReturnRows(rows)

		// Act
		subscribers, err := repo.GetByNetwork(networkId)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, subscribers, 1)
		assert.Equal(t, subscriber.SubscriberId, subscribers[0].SubscriberId)
		assert.Equal(t, subscriber.Name, subscribers[0].Name)
		assert.Equal(t, subscriber.Email, subscribers[0].Email)
		assert.Equal(t, networkId, subscribers[0].NetworkId)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SuccessWithNoSubscribers", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		networkId := uuid.NewV4()

		rows := sqlmock.NewRows([]string{
			"subscriber_id", "name", "network_id", "email", "phone_number",
			"gender", "dob", "proof_of_identification", "id_serial", "address",
			"created_at", "updated_at", "deleted_at",
		})

		mock.ExpectQuery(`^SELECT \* FROM "subscribers" WHERE "subscribers"."network_id" = \$1`).
			WithArgs(networkId).
			WillReturnRows(rows)

		// Act
		subscribers, err := repo.GetByNetwork(networkId)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, subscribers, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		networkId := uuid.NewV4()
		dbErr := errors.New(TEST_ERROR_DB_CONNECTION)

		mock.ExpectQuery(`^SELECT \* FROM "subscribers" WHERE "subscribers"."network_id" = \$1`).
			WithArgs(networkId).
			WillReturnError(dbErr)

		// Act
		subscribers, err := repo.GetByNetwork(networkId)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.Nil(t, subscribers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriber_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mock, repo := setupTestDB(t)
		subscriberId := uuid.NewV4()
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "subscribers"`).WithArgs(subscriberId).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		err := repo.Delete(subscriberId)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock, repo := setupTestDB(t)
		subscriberId := uuid.NewV4()
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "subscribers"`).WithArgs(subscriberId).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()
		err := repo.Delete(subscriberId)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock, repo := setupTestDB(t)
		subscriberId := uuid.NewV4()
		dbErr := errors.New(TEST_ERROR_DB)
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "subscribers"`).WithArgs(subscriberId).WillReturnError(dbErr)
		mock.ExpectRollback()
		err := repo.Delete(subscriberId)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriber_ListSubscribers(t *testing.T) {
	t.Run("SuccessWithMultipleSubscribers", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		subscriber1 := createTestSubscriber(TEST_NAME_JOHN, TEST_EMAIL_JOHN)
		subscriber2 := createTestSubscriber(TEST_NAME_JANE, TEST_EMAIL_JANE)
		subscriber3 := createTestSubscriber(TEST_NAME_BOB, TEST_EMAIL_BOB)

		rows := sqlmock.NewRows([]string{
			"subscriber_id", "name", "network_id", "email", "phone_number",
			"gender", "dob", "proof_of_identification", "id_serial", "address",
			"created_at", "updated_at", "deleted_at",
		}).AddRow(
			subscriber1.SubscriberId, subscriber1.Name, subscriber1.NetworkId,
			subscriber1.Email, subscriber1.PhoneNumber, subscriber1.Gender,
			subscriber1.DOB, subscriber1.ProofOfIdentification, subscriber1.IdSerial,
			subscriber1.Address, subscriber1.CreatedAt, subscriber1.UpdatedAt, subscriber1.DeletedAt,
		).AddRow(
			subscriber2.SubscriberId, subscriber2.Name, subscriber2.NetworkId,
			subscriber2.Email, subscriber2.PhoneNumber, subscriber2.Gender,
			subscriber2.DOB, subscriber2.ProofOfIdentification, subscriber2.IdSerial,
			subscriber2.Address, subscriber2.CreatedAt, subscriber2.UpdatedAt, subscriber2.DeletedAt,
		).AddRow(
			subscriber3.SubscriberId, subscriber3.Name, subscriber3.NetworkId,
			subscriber3.Email, subscriber3.PhoneNumber, subscriber3.Gender,
			subscriber3.DOB, subscriber3.ProofOfIdentification, subscriber3.IdSerial,
			subscriber3.Address, subscriber3.CreatedAt, subscriber3.UpdatedAt, subscriber3.DeletedAt,
		)

		mock.ExpectQuery(`^SELECT \* FROM "subscribers"`).WillReturnRows(rows)

		// Act
		subscribers, err := repo.ListSubscribers()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, subscribers, 3)
		assert.Equal(t, subscriber1.SubscriberId, subscribers[0].SubscriberId)
		assert.Equal(t, subscriber1.Name, subscribers[0].Name)
		assert.Equal(t, subscriber1.Email, subscribers[0].Email)
		assert.Equal(t, subscriber2.SubscriberId, subscribers[1].SubscriberId)
		assert.Equal(t, subscriber2.Name, subscribers[1].Name)
		assert.Equal(t, subscriber2.Email, subscribers[1].Email)
		assert.Equal(t, subscriber3.SubscriberId, subscribers[2].SubscriberId)
		assert.Equal(t, subscriber3.Name, subscribers[2].Name)
		assert.Equal(t, subscriber3.Email, subscribers[2].Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SuccessWithSingleSubscriber", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		subscriber := createTestSubscriber(TEST_NAME_ALICE, TEST_EMAIL_ALICE)

		rows := sqlmock.NewRows([]string{
			"subscriber_id", "name", "network_id", "email", "phone_number",
			"gender", "dob", "proof_of_identification", "id_serial", "address",
			"created_at", "updated_at", "deleted_at",
		}).AddRow(
			subscriber.SubscriberId, subscriber.Name, subscriber.NetworkId,
			subscriber.Email, subscriber.PhoneNumber, subscriber.Gender,
			subscriber.DOB, subscriber.ProofOfIdentification, subscriber.IdSerial,
			subscriber.Address, subscriber.CreatedAt, subscriber.UpdatedAt, subscriber.DeletedAt,
		)

		mock.ExpectQuery(`^SELECT \* FROM "subscribers"`).WillReturnRows(rows)

		// Act
		subscribers, err := repo.ListSubscribers()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, subscribers, 1)
		assert.Equal(t, subscriber.SubscriberId, subscribers[0].SubscriberId)
		assert.Equal(t, subscriber.Name, subscribers[0].Name)
		assert.Equal(t, subscriber.Email, subscribers[0].Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SuccessWithNoSubscribers", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		rows := sqlmock.NewRows([]string{
			"subscriber_id", "name", "network_id", "email", "phone_number",
			"gender", "dob", "proof_of_identification", "id_serial", "address",
			"created_at", "updated_at", "deleted_at",
		})

		mock.ExpectQuery(`^SELECT \* FROM "subscribers"`).WillReturnRows(rows)

		// Act
		subscribers, err := repo.ListSubscribers()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, subscribers, 0)
		assert.Empty(t, subscribers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		expectedError := errors.New(TEST_ERROR_DB_FAILED)
		mock.ExpectQuery(`^SELECT \* FROM "subscribers"`).WillReturnError(expectedError)

		// Act
		subscribers, err := repo.ListSubscribers()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, subscribers)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriber_GetByEmail(t *testing.T) {
	t.Run("SubscriberFound", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		expectedEmail := TEST_EMAIL_JOHN
		subscriber := createTestSubscriber(TEST_NAME_JOHN, expectedEmail)

		subRow := sqlmock.NewRows([]string{
			"subscriber_id", "name", "network_id", "email", "phone_number",
			"gender", "dob", "proof_of_identification", "id_serial", "address",
			"created_at", "updated_at", "deleted_at",
		}).AddRow(
			subscriber.SubscriberId, subscriber.Name, subscriber.NetworkId,
			subscriber.Email, subscriber.PhoneNumber, subscriber.Gender,
			subscriber.DOB, subscriber.ProofOfIdentification, subscriber.IdSerial,
			subscriber.Address, subscriber.CreatedAt, subscriber.UpdatedAt, subscriber.DeletedAt,
		)

		mock.ExpectQuery(`^SELECT \* FROM "subscribers" WHERE email = \$1 ORDER BY "subscribers"\."subscriber_id" LIMIT \$2`).
			WithArgs(expectedEmail, 1).
			WillReturnRows(subRow)

		// Act
		result, err := repo.GetByEmail(expectedEmail)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, subscriber.SubscriberId, result.SubscriberId)
		assert.Equal(t, subscriber.Name, result.Name)
		assert.Equal(t, subscriber.Email, result.Email)
		assert.Equal(t, subscriber.NetworkId, result.NetworkId)
		assert.Equal(t, subscriber.PhoneNumber, result.PhoneNumber)
		assert.Equal(t, subscriber.Gender, result.Gender)
		assert.Equal(t, subscriber.DOB, result.DOB)
		assert.Equal(t, subscriber.ProofOfIdentification, result.ProofOfIdentification)
		assert.Equal(t, subscriber.IdSerial, result.IdSerial)
		assert.Equal(t, subscriber.Address, result.Address)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SubscriberNotFound", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		nonExistentEmail := TEST_EMAIL_NONEXISTENT

		mock.ExpectQuery(`^SELECT \* FROM "subscribers" WHERE email = \$1 ORDER BY "subscribers"\."subscriber_id" LIMIT \$2`).
			WithArgs(nonExistentEmail, 1).
			WillReturnError(sql.ErrNoRows)

		// Act
		result, err := repo.GetByEmail(nonExistentEmail)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		email := TEST_EMAIL_TEST
		expectedError := errors.New(TEST_ERROR_DB_FAILED)

		mock.ExpectQuery(`^SELECT \* FROM "subscribers" WHERE email = \$1 ORDER BY "subscribers"\."subscriber_id" LIMIT \$2`).
			WithArgs(email, 1).
			WillReturnError(expectedError)

		// Act
		result, err := repo.GetByEmail(email)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriber_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mock, repo := setupTestDB(t)

		subscriberId := uuid.NewV4()
		originalSubscriber := createTestSubscriber(TEST_NAME_JOHN, TEST_EMAIL_JOHN)
		originalSubscriber.SubscriberId = subscriberId

		updatedSubscriber := int_db.Subscriber{
			Name:                  TEST_NAME_JOHN_UPDATED,
			NetworkId:             uuid.NewV4(),
			Email:                 TEST_EMAIL_JOHN_UPDATED,
			PhoneNumber:           TEST_PHONE_UPDATED,
			Gender:                TEST_GENDER_MALE,
			DOB:                   TEST_DOB_UPDATED,
			ProofOfIdentification: TEST_PROOF_JANE,
			IdSerial:              TEST_ID_SERIAL_JANE,
			Address:               TEST_ADDRESS_UPDATED,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`^UPDATE "subscribers" SET`).
			WithArgs(
				updatedSubscriber.Name,
				updatedSubscriber.NetworkId,
				updatedSubscriber.Email,
				updatedSubscriber.PhoneNumber,
				updatedSubscriber.Gender,
				updatedSubscriber.DOB,
				updatedSubscriber.ProofOfIdentification,
				updatedSubscriber.IdSerial,
				updatedSubscriber.Address,
				sqlmock.AnyArg(),
				subscriberId,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// Act
		err := repo.Update(subscriberId, updatedSubscriber)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
