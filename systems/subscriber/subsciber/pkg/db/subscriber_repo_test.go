package db

import (
	extsql "database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
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

func Test_Subscriber_Get(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		subscriberID := "7333911a-effb-4da9-949c-9f79fac688dd"

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"subscriber_id", "first_name", "last_name", "network_id", "email", "gender", "phone_number",
			"address", "id_serial"}).
			AddRow(subscriberID, "John", "Doe", "dd109466-8100-4450-87b1-472abb60e319", "john.doe@example.com", "male", "123-456-7890", "123 Main St.", "123456789")
		mock.ExpectQuery(`^SELECT.*subscriber.*`).
			WithArgs(subscriberID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		sub, err := r.Get(subscriberID)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sub)
	})
}
func Test_Package_Delete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		subscriberID := "7333911a-effb-4da9-949c-9f79fac688dd"

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Expect a BEGIN statement to start the transaction
		mock.ExpectBegin()

		// Expect a DELETE statement with the correct subscriber ID
		mock.ExpectExec("DELETE FROM \"subscribers\" WHERE \"subscriber_id\" = ?").
		WithArgs(subscriberID).
		WillReturnResult(sqlmock.NewResult(1, 1))
		// Expect a COMMIT statement to end the transaction
		mock.ExpectCommit()

		// Set up the database connection
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})	
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
		r := NewSubscriberRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.Delete(subscriberID)
		assert.NoError(t, err)

		// Check that all expected queries were executed
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}