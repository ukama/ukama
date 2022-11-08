package db_test

import (
	extsql "database/sql"
	"regexp"
	"testing"

	"github.com/google/uuid"
	org_db "github.com/ukama/ukama/systems/registry/org/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_UserRepo_Get(t *testing.T) {
	t.Run("GetUser", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		var userUUID = uuid.New()

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"uuid"}).
			AddRow(userUUID)

		mock.ExpectQuery(`^SELECT.*users.*`).
			WithArgs(userUUID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := org_db.NewUserRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		org, err := r.Get(userUUID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, org)
	})
}

func Test_UserRepo_Add(t *testing.T) {
	t.Run("AddUser", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		var testUUID = uuid.New()

		user := org_db.User{
			Uuid: testUUID,
		}

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(user.Uuid, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := org_db.NewUserRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&user)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

// func Test_UserRepo_Delete(t *testing.T) {
// t.Run("DeleteUser", func(t *testing.T) {
// var db *extsql.DB

// var userUUID = uuid.New()

// db, mock, err := sqlmock.New() // mock sql.DB
// assert.NoError(t, err)

// mock.ExpectBegin()

// mock.ExpectExec(regexp.QuoteMeta("DELETE")).WithArgs(userUUID).
// WillReturnResult(sqlmock.NewResult(1, 1))

// mock.ExpectCommit()

// dialector := postgres.New(postgres.Config{
// DSN:                  "sqlmock_db_0",
// DriverName:           "postgres",
// Conn:                 db,
// PreferSimpleProtocol: true,
// })

// gdb, err := gorm.Open(dialector, &gorm.Config{})
// assert.NoError(t, err)

// r := org_db.NewUserRepo(&UkamaDbMock{
// GormDb: gdb,
// })

// assert.NoError(t, err)

// // Act
// err = r.Delete(userUUID)

// // Assert
// assert.NoError(t, err)

// err = mock.ExpectationsWereMet()
// assert.NoError(t, err)
// })
// }
