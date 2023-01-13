package db_test

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"

	int_db "github.com/ukama/ukama/systems/init/msgClient/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic("implement me")
}

func (u UkamaDbMock) Connect() error {
	panic("implement me")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func() error) error {
	log.Fatal("implement me")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) (err error) {
	log.Fatal("implement me")
	return nil
}

func Test_routeRepo_Get(t *testing.T) {

	t.Run("RouteExist", func(t *testing.T) {
		// Arrange
		const key = "event.cloud.lookup.organization.create"

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"key"}).
			AddRow(key)

		mock.ExpectQuery(`^SELECT.*routes.*`).
			WithArgs(key).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewRouteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		route, err := r.Get(key)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, route) {
			assert.Equal(t, route.Key, key)
		}

	})

}

func Test_routeRepo_Delete(t *testing.T) {

	t.Run("DeleteRoute", func(t *testing.T) {

		const key = "event.cloud.lookup.organization.create"

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), key).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewRouteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Remove(key)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}

func Test_routeRepo_Add(t *testing.T) {

	t.Run("AddRouteWithExistingKey", func(t *testing.T) {

		const key = "event.cloud.lookup.organization.create"

		// route := int_db.Route{
		// 	Key: key,
		// }

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"key"}).
			AddRow(key)

		mock.ExpectQuery(`^SELECT.*routes.*`).
			WithArgs(key).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewRouteRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		_, err = r.Add(key)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	/*
		t.Run("AddRouteWithNonExistingKey", func(t *testing.T) {

			const key = "event.cloud.lookup.organization.create"

			// route := int_db.Route{
			// 	Key: key,
			// }

			var db *extsql.DB
			var err error

			db, mock, err := sqlmock.New() // mock sql.DB
			assert.NoError(t, err)

			mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
				WithArgs(key).
				WillReturnRows(sqlmock.NewRows([]string{"key"}))

			// mock.ExpectQuery(regexp.QuoteMeta("INSERT")).
			// 	WithArgs(key).
			// 	WillReturnRows(sqlmock.NewRows([]string{"key"}).AddRow(key))

			dialector := postgres.New(postgres.Config{
				DSN:                  "sqlmock_db_0",
				DriverName:           "postgres",
				Conn:                 db,
				PreferSimpleProtocol: true,
			})
			gdb, err := gorm.Open(dialector, &gorm.Config{})
			assert.NoError(t, err)

			r := int_db.NewRouteRepo(&UkamaDbMock{
				GormDb: gdb,
			})

			assert.NoError(t, err)

			// Act
			_, err = r.Add(key)

			// Assert
			assert.NoError(t, err)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	*/
}
