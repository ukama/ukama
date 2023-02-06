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
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	repo := int_db.NewSubscriberRepo(&UkamaDbMock{
		GormDb: gdb,
	})
	subscriber := int_db.Subscriber{
		SubscriberID:          uuid.NewV4(),
		FirstName:             "John",
		LastName:              "Doe",
		NetworkID:             uuid.NewV4(),
		OrgID:                 uuid.NewV4(),
		Email:                 "johndoe@example.com",
		PhoneNumber:           "555-555-5555",
		Gender:                "Male",
		DOB:                   time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
		ProofOfIdentification: "Driver's License",
		IdSerial:              "ABC123",
		Address:               "123 Main St.",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"subscribers\"").WithArgs(
		subscriber.SubscriberID,
		subscriber.FirstName,
		subscriber.LastName,
		subscriber.NetworkID,
		subscriber.OrgID,
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

	err = repo.Add(&subscriber)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestSubscriber_Get(t *testing.T) {

	t.Run("SubscriberFound", func(t *testing.T) {
		var subID = uuid.NewV4()

		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
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
			WithArgs(subID).
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
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		gdb, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*subscribers.*`).
			WithArgs(subID).
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
		var networkID = uuid.NewV4()

		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		gdb, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

		subRow := sqlmock.NewRows([]string{"network_id"}).
			AddRow(networkID)

		mock.ExpectQuery(`^SELECT.*subscribers.*`).
			WithArgs(networkID).
			WillReturnRows(subRow)

		assert.NoError(t, err)
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		// Act
		sub, err := repo.Get(networkID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sub)
	})
	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		var networkID = uuid.NewV4()

		// Arrange
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		gdb, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*subscribers.*`).
			WithArgs(networkID).
			WillReturnError(sql.ErrNoRows)
		repo := int_db.NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		assert.NoError(t, err)

		// Act
		sub, err := repo.Get(networkID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, sub)
	})
}
