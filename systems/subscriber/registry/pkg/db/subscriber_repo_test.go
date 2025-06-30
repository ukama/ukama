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

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, int_db.SubscriberRepo) {
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

	return db, mock, repo
}

func createTestSubscriber(name, email string) int_db.Subscriber {
	return int_db.Subscriber{
		SubscriberId:          uuid.NewV4(),
		Name:                  name,
		NetworkId:             uuid.NewV4(),
		Email:                 email,
		PhoneNumber:           "555-555-5555",
		Gender:                "Male",
		DOB:                   "07-03-2023",
		ProofOfIdentification: "Driver's License",
		IdSerial:              "ABC123",
		Address:               "123 Main St.",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	}
}

func TestSubscriber_Add(t *testing.T) {
	t.Run("SuccessWithoutNestedFunc", func(t *testing.T) {
		// Arrange
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		subscriber := createTestSubscriber("John", "john@example.com")

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
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		subscriber := int_db.Subscriber{
			SubscriberId:          uuid.NewV4(),
			Name:                  "Jane",
			NetworkId:             uuid.NewV4(),
			Email:                 "jane@example.com",
			PhoneNumber:           "555-555-5556",
			Gender:                "Female",
			DOB:                   "15-06-1990",
			ProofOfIdentification: "Passport",
			IdSerial:              "XYZ789",
			Address:               "456 Oak Ave.",
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
		err = repo.Add(&subscriber, nestedFunc)

		// Assert
		assert.NoError(t, err)
		assert.True(t, nestedFuncCalled)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NestedFuncReturnsError", func(t *testing.T) {
		// Arrange
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
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		subscriber := int_db.Subscriber{
			SubscriberId:          uuid.NewV4(),
			Name:                  "Bob",
			NetworkId:             uuid.NewV4(),
			Email:                 "bob@example.com",
			PhoneNumber:           "555-555-5557",
			Gender:                "Male",
			DOB:                   "22-12-1985",
			ProofOfIdentification: "National ID",
			IdSerial:              "DEF456",
			Address:               "789 Pine St.",
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
			DeletedAt:             nil,
		}

		expectedError := errors.New("expected error")
		nestedFunc := func(sub *int_db.Subscriber, tx *gorm.DB) error {
			return expectedError
		}

		mock.ExpectBegin()
		mock.ExpectRollback()

		// Act
		err = repo.Add(&subscriber, nestedFunc)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseCreateFails", func(t *testing.T) {
		// Arrange
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
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		subscriber := int_db.Subscriber{
			SubscriberId:          uuid.NewV4(),
			Name:                  "Alice",
			NetworkId:             uuid.NewV4(),
			Email:                 "alice@example.com",
			PhoneNumber:           "555-555-5558",
			Gender:                "Female",
			DOB:                   "03-09-1992",
			ProofOfIdentification: "Birth Certificate",
			IdSerial:              "GHI789",
			Address:               "321 Elm St.",
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
			DeletedAt:             nil,
		}

		expectedError := errors.New("expected error")
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
		err = repo.Add(&subscriber, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseTransactionFails", func(t *testing.T) {
		// Arrange
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
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		subscriber := int_db.Subscriber{
			SubscriberId:          uuid.NewV4(),
			Name:                  "Charlie",
			NetworkId:             uuid.NewV4(),
			Email:                 "charlie@example.com",
			PhoneNumber:           "555-555-5559",
			Gender:                "Male",
			DOB:                   "18-04-1988",
			ProofOfIdentification: "Social Security Card",
			IdSerial:              "JKL012",
			Address:               "654 Maple Dr.",
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
			DeletedAt:             nil,
		}

		expectedError := errors.New("expected error")
		mock.ExpectBegin().WillReturnError(expectedError)

		// Act
		err = repo.Add(&subscriber, nil)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NestedFuncWithDatabaseError", func(t *testing.T) {
		// Arrange
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
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		subscriber := int_db.Subscriber{
			SubscriberId:          uuid.NewV4(),
			Name:                  "Diana",
			NetworkId:             uuid.NewV4(),
			Email:                 "diana@example.com",
			PhoneNumber:           "555-555-5560",
			Gender:                "Female",
			DOB:                   "11-07-1995",
			ProofOfIdentification: "Voter ID",
			IdSerial:              "MNO345",
			Address:               "987 Cedar Ln.",
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
		err = repo.Add(&subscriber, nestedFunc)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrInvalidDB, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SubscriberWithMinimalFields", func(t *testing.T) {
		// Arrange
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
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		subscriber := int_db.Subscriber{
			SubscriberId: uuid.NewV4(),
			Name:         "Minimal",
			NetworkId:    uuid.NewV4(),
			Email:        "minimal@example.com",
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
		err = repo.Add(&subscriber, nil)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSubscriber_Get(t *testing.T) {

	t.Run("SubscriberFound", func(t *testing.T) {
		var subID = uuid.NewV4()

		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		gdb, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

		subRow := sqlmock.NewRows([]string{"subscriber_id"}).
			AddRow(subID)

		mock.ExpectQuery(`^SELECT.*subscribers.*`).
			WithArgs(subID, sqlmock.AnyArg()).
			WillReturnRows(subRow)

		assert.NoError(t, err)
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})
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
		db, mock, err := sqlmock.New()

		assert.NoError(t, err)

		defer db.Close()
		gdb, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*subscribers.*`).
			WithArgs(subID, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		assert.NoError(t, err)

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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		networkId := uuid.NewV4()
		subscriber1 := createTestSubscriber("John", "john@example.com")
		subscriber1.NetworkId = networkId
		subscriber2 := createTestSubscriber("Jane", "jane@example.com")
		subscriber2.NetworkId = networkId
		subscriber3 := createTestSubscriber("Bob", "bob@example.com")
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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		networkId := uuid.NewV4()
		subscriber := createTestSubscriber("Alice", "alice@example.com")
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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		networkId := uuid.NewV4()
		dbErr := errors.New("database connection error")

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
		db, mock, repo := setupTestDB(t)
		defer db.Close()
		subscriberId := uuid.NewV4()
		mock.ExpectBegin()
		mock.ExpectExec(`DELETE FROM "subscribers"`).WithArgs(subscriberId).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		err := repo.Delete(subscriberId)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		db, mock, repo := setupTestDB(t)
		defer db.Close()
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
		db, mock, repo := setupTestDB(t)
		defer db.Close()
		subscriberId := uuid.NewV4()
		dbErr := errors.New("db error")
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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		subscriber1 := createTestSubscriber("John", "john@example.com")
		subscriber2 := createTestSubscriber("Jane", "jane@example.com")
		subscriber3 := createTestSubscriber("Bob", "bob@example.com")

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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		subscriber := createTestSubscriber("Alice", "alice@example.com")

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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		expectedError := errors.New("database connection failed")
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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		expectedEmail := "john@example.com"
		subscriber := createTestSubscriber("John", expectedEmail)

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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		nonExistentEmail := "nonexistent@example.com"

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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		email := "test@example.com"
		expectedError := errors.New("database connection failed")

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
		db, mock, repo := setupTestDB(t)
		defer db.Close()

		subscriberId := uuid.NewV4()
		originalSubscriber := createTestSubscriber("John", "john@example.com")
		originalSubscriber.SubscriberId = subscriberId

		updatedSubscriber := int_db.Subscriber{
			Name:                  "John Updated",
			NetworkId:             uuid.NewV4(),
			Email:                 "john.updated@example.com",
			PhoneNumber:           "555-555-9999",
			Gender:                "Male",
			DOB:                   "15-08-1990",
			ProofOfIdentification: "Passport",
			IdSerial:              "XYZ789",
			Address:               "456 Updated St.",
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
