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
		FirstName:             "John",
		LastName:              "Doe",
		NetworkId:             uuid.NewV4(),
		Email:                 "johndoe@example.com",
		PhoneNumber:           "555-555-5555",
		Gender:                "Male",
		DOB:                   dateStr,
		ProofOfIdentification: "Driver's License",
		IdSerial:              "ABC123",
		Address:               "123 Main St.",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"subscribers\"").WithArgs(
		subscriber.SubscriberId,
		subscriber.FirstName,
		subscriber.LastName,
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
func TestSubscriber_GetByEmail(t *testing.T) {
    t.Run("SubscriberFoundByEmail", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        assert.NoError(t, err)
        defer db.Close()
        gdb, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

        email := "hello@ukama.com"
        subRow := sqlmock.NewRows([]string{"subscriber_id", "email"}).
            AddRow(uuid.NewV4(), email)

        mock.ExpectQuery(`^SELECT \* FROM "subscribers" WHERE email = \$1 ORDER BY "subscribers"."subscriber_id" LIMIT \$2`).
            WithArgs(email, 1).
            WillReturnRows(subRow)

        repo := int_db.NewSubscriberRepo(&UkamaDbMock{GormDb: gdb})
        sub, err := repo.GetByEmail(email)

        assert.NoError(t, err)
        assert.NotNil(t, sub)
        assert.Equal(t, email, sub.Email)
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("SubscriberNotFoundByEmail", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        assert.NoError(t, err)
        defer db.Close()
        gdb, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

        email := "nonexistent@example.com"

        mock.ExpectQuery(`^SELECT \* FROM "subscribers" WHERE email = \$1 ORDER BY "subscribers"."subscriber_id" LIMIT \$2`).
            WithArgs(email, 1).
            WillReturnError(gorm.ErrRecordNotFound)

        repo := int_db.NewSubscriberRepo(&UkamaDbMock{GormDb: gdb})
        sub, err := repo.GetByEmail(email)

        assert.Error(t, err)
        assert.Nil(t, sub)
        assert.Equal(t, gorm.ErrRecordNotFound, err)
        assert.NoError(t, mock.ExpectationsWereMet())
    })
}


func TestSubscriber_GetByNetwork(t *testing.T) {
    t.Run("SubscribersFoundByNetwork", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        assert.NoError(t, err)
        defer db.Close()
        gdb, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

        networkID := uuid.NewV4()
        subRows := sqlmock.NewRows([]string{"subscriber_id", "network_id"}).
            AddRow(uuid.NewV4(), networkID).
            AddRow(uuid.NewV4(), networkID)

        mock.ExpectQuery(`^SELECT.*subscribers.*WHERE.*network_id.*`).
            WithArgs(networkID).
            WillReturnRows(subRows)

        repo := int_db.NewSubscriberRepo(&UkamaDbMock{GormDb: gdb})
        subs, err := repo.GetByNetwork(networkID)

        assert.NoError(t, err)
        assert.Len(t, subs, 2)
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("NoSubscribersFoundByNetwork", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        assert.NoError(t, err)
        defer db.Close()
        gdb, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

        networkID := uuid.NewV4()

        mock.ExpectQuery(`^SELECT.*subscribers.*WHERE.*network_id.*`).
            WithArgs(networkID).
            WillReturnRows(sqlmock.NewRows([]string{}))

        repo := int_db.NewSubscriberRepo(&UkamaDbMock{GormDb: gdb})
        subs, err := repo.GetByNetwork(networkID)

        assert.NoError(t, err)
        assert.Empty(t, subs)
        assert.NoError(t, mock.ExpectationsWereMet())
    })
}
func TestSubscriber_ListSubscribers(t *testing.T) {
    t.Run("ListAllSubscribers", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        assert.NoError(t, err)
        defer db.Close()
        gdb, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

        subRows := sqlmock.NewRows([]string{"subscriber_id"}).
            AddRow(uuid.NewV4()).
            AddRow(uuid.NewV4()).
            AddRow(uuid.NewV4())

        mock.ExpectQuery(`^SELECT.*subscribers.*`).
            WillReturnRows(subRows)

        repo := int_db.NewSubscriberRepo(&UkamaDbMock{GormDb: gdb})
        subs, err := repo.ListSubscribers()

        assert.NoError(t, err)
        assert.Len(t, subs, 3)
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("EmptySubscriberList", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        assert.NoError(t, err)
        defer db.Close()
        gdb, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

        mock.ExpectQuery(`^SELECT.*subscribers.*`).
            WillReturnRows(sqlmock.NewRows([]string{}))

        repo := int_db.NewSubscriberRepo(&UkamaDbMock{GormDb: gdb})
        subs, err := repo.ListSubscribers()

        assert.NoError(t, err)
        assert.Empty(t, subs)
        assert.NoError(t, mock.ExpectationsWereMet())
    })
}
func TestSubscriber_Update(t *testing.T) {
    t.Run("SuccessfulUpdate", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        assert.NoError(t, err)
        defer db.Close()
        gdb, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

        subscriberId := uuid.NewV4()
        updatedSub := int_db.Subscriber{
            FirstName: "UpdatedName",
            LastName:  "UpdatedLastName",
            Email:     "updated@example.com",
        }

        mock.ExpectBegin()
        mock.ExpectExec(`UPDATE "subscribers" SET`).
            WithArgs(
                updatedSub.FirstName,
                updatedSub.LastName,
                updatedSub.Email,
                sqlmock.AnyArg(), // for updated_at
                subscriberId,
            ).
            WillReturnResult(sqlmock.NewResult(0, 1))
        mock.ExpectCommit()

        // Expect the SELECT query after the update, including the LIMIT clause
        mock.ExpectQuery(`SELECT \* FROM "subscribers" WHERE subscriber_id = \$1 ORDER BY "subscribers"."subscriber_id" LIMIT \$2`).
            WithArgs(subscriberId, 1).
            WillReturnRows(sqlmock.NewRows([]string{"subscriber_id", "first_name", "last_name", "email"}).
                AddRow(subscriberId, updatedSub.FirstName, updatedSub.LastName, updatedSub.Email))

        repo := int_db.NewSubscriberRepo(&UkamaDbMock{GormDb: gdb})
        err = repo.Update(subscriberId, updatedSub)

        assert.NoError(t, err)
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("UpdateNotFound", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        assert.NoError(t, err)
        defer db.Close()
        gdb, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

        subscriberId := uuid.NewV4()
        updatedSub := int_db.Subscriber{
            FirstName: "UpdatedName",
            LastName:  "UpdatedLastName",
            Email:     "updated@example.com",
        }

        mock.ExpectBegin()
        mock.ExpectExec(`UPDATE "subscribers" SET`).
            WithArgs(
                updatedSub.FirstName,
                updatedSub.LastName,
                updatedSub.Email,
                sqlmock.AnyArg(), // for updated_at
                subscriberId,
            ).
            WillReturnResult(sqlmock.NewResult(0, 0))
        mock.ExpectCommit()

        // Expect the SELECT query after the update, including the LIMIT clause, but return no rows
        mock.ExpectQuery(`SELECT \* FROM "subscribers" WHERE subscriber_id = \$1 ORDER BY "subscribers"."subscriber_id" LIMIT \$2`).
            WithArgs(subscriberId, 1).
            WillReturnError(gorm.ErrRecordNotFound)

        repo := int_db.NewSubscriberRepo(&UkamaDbMock{GormDb: gdb})
        err = repo.Update(subscriberId, updatedSub)

        assert.Error(t, err)
        assert.Equal(t, gorm.ErrRecordNotFound, err)
        assert.NoError(t, mock.ExpectationsWereMet())
    })

    t.Run("UpdateError", func(t *testing.T) {
        db, mock, err := sqlmock.New()
        assert.NoError(t, err)
        defer db.Close()
        gdb, _ := gorm.Open(postgres.New(postgres.Config{
            Conn: db,
        }), &gorm.Config{})

        subscriberId := uuid.NewV4()
        updatedSub := int_db.Subscriber{
            FirstName: "UpdatedName",
            LastName:  "UpdatedLastName",
            Email:     "updated@example.com",
        }

        mock.ExpectBegin()
        mock.ExpectExec(`UPDATE "subscribers" SET`).
            WithArgs(
                updatedSub.FirstName,
                updatedSub.LastName,
                updatedSub.Email,
                sqlmock.AnyArg(), // for updated_at
                subscriberId,
            ).
            WillReturnError(sql.ErrConnDone)
        mock.ExpectRollback()

        repo := int_db.NewSubscriberRepo(&UkamaDbMock{GormDb: gdb})
        err = repo.Update(subscriberId, updatedSub)

        assert.Error(t, err)
        assert.Equal(t, sql.ErrConnDone, err)
        assert.NoError(t, mock.ExpectationsWereMet())
    })
}