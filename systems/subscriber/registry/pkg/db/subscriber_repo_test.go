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
	"log"
	"testing"
	"time"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
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

func TestSubscriber_Add(t *testing.T) {
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
    dateStr := "07-03-2023"

    subscriber := int_db.Subscriber{
        SubscriberId:          uuid.NewV4(),
        Name:                  "John",
        NetworkId:             uuid.NewV4(),
        Email:                 "john@example.com",  
        PhoneNumber:           "555-555-5555",
        Gender:                "Male",
        DOB:                   dateStr,
        ProofOfIdentification: "Driver's License",
        SubscriberStatus:      ukama.SubscriberStatusActive, 
        DeletionRetryCount:    0,                   
        DeletionLastAttempt:   nil,                 
        IdSerial:              "ABC123",
        Address:               "123 Main St.",
        CreatedAt:             time.Now(),
        UpdatedAt:             time.Now(),
        DeletedAt:             nil,
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
        subscriber.SubscriberStatus, 
        subscriber.DeletionRetryCount,    
        subscriber.DeletionLastAttempt,   
        subscriber.IdSerial,
        subscriber.Address,
        subscriber.CreatedAt,
        subscriber.UpdatedAt,
        subscriber.DeletedAt).WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    err = repo.Add(&subscriber, nil)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
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

		subRow := sqlmock.NewRows([]string{
			"subscriber_id", "name", "network_id", "email", "phone_number", 
			"gender", "dob", "proof_of_identification", "subscriber_status", 
			"deletion_retry_count", "deletion_last_attempt", 
			"id_serial", "address", "created_at", "updated_at", "deleted_at"}).
			AddRow(subID, "John", uuid.NewV4(), "john@example.com", "555-555-5555", 
				"Male", "1990-01-01", "Passport", ukama.SubscriberStatusActive, 
				0, nil, 
				"ABC123", "123 Main St", time.Now(), time.Now(), nil)

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

	t.Run("NetworkFound", func(t *testing.T) {
		var networkId = uuid.NewV4()

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

		subRow := sqlmock.NewRows([]string{"network_id"}).
			AddRow(networkId)

		mock.ExpectQuery(`^SELECT.*subscribers.*`).
			WithArgs(networkId, sqlmock.AnyArg()).
			WillReturnRows(subRow)

		assert.NoError(t, err)
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		// Act
		sub, err := repo.Get(networkId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sub)
	})
	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		var networkId = uuid.NewV4()

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
			WithArgs(networkId, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		assert.NoError(t, err)

		// Act
		sub, err := repo.Get(networkId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, sub)
	})
}

func TestSubscriber_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
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

	var subscriberId = uuid.NewV4()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"subscribers\"").WithArgs(subscriberId).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.Delete(subscriberId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
