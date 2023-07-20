package db_test

import (
	"database/sql"
	extsql "database/sql"
	"regexp"
	"testing"

	"github.com/ukama/ukama/systems/common/uuid"
	org_db "github.com/ukama/ukama/systems/registry/org/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_UserRepo_Add(t *testing.T) {
	var db *extsql.DB

	var testUUID = uuid.NewV4()

	user := org_db.User{
		Uuid: testUUID,
	}

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

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

	t.Run("AddUser", func(t *testing.T) {
		// Arrange

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(user.Uuid, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		// Act
		err = r.Add(&user, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUserRepo_Get(t *testing.T) {
	var db *extsql.DB

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

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

	t.Run("UserFound", func(t *testing.T) {
		// Arrange

		var userUUID = uuid.NewV4()

		rows := sqlmock.NewRows([]string{"uuid"}).
			AddRow(userUUID)

		mock.ExpectQuery(`^SELECT.*users.*`).
			WithArgs(userUUID).
			WillReturnRows(rows)

		// Act
		org, err := r.Get(userUUID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, org)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Arrange

		var userUUID = uuid.NewV4()

		mock.ExpectQuery(`^SELECT.*users.*`).
			WithArgs(userUUID).
			WillReturnError(sql.ErrNoRows)

		// Act
		org, err := r.Get(userUUID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, org)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUserRepo_Delete(t *testing.T) {
	var db *extsql.DB

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

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

	t.Run("UserFound", func(t *testing.T) {
		var userUUID = uuid.NewV4()

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
			WithArgs(sqlmock.AnyArg(), userUUID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err = r.Delete(userUUID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		var userUUID = uuid.NewV4()

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
			WithArgs(sqlmock.AnyArg(), userUUID).
			WillReturnError(sql.ErrNoRows)

		// Act
		err = r.Delete(userUUID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUserRepo_GetUserCount(t *testing.T) {
	var db *extsql.DB

	db, mock, err := sqlmock.New() // mock sql.DB
	assert.NoError(t, err)

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

	t.Run("UserFound", func(t *testing.T) {
		// Arrange

		rowsCount1 := sqlmock.NewRows([]string{"count"}).
			AddRow(2)

		rowsCount2 := sqlmock.NewRows([]string{"count"}).
			AddRow(1)

		mock.ExpectQuery(`^SELECT count(\\*).*users.*`).
			WillReturnRows(rowsCount1)

		mock.ExpectQuery(`^SELECT count(\\*).*users.*WHERE.*`).
			WillReturnRows(rowsCount2)

		// Act
		activeUsr, inactiveUsr, err := r.GetUserCount()
		assert.NoError(t, err)

		// Assert
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		assert.Equal(t, int64(2), activeUsr)
		assert.Equal(t, int64(1), inactiveUsr)
	})
}
