package db

import (
	extsql "database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDefaultMarkupRepo_Create(t *testing.T) {

	t.Run("Create", func(t *testing.T) {
		// Arrange
		var markup float64 = 10

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), markup).
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

		r := NewDefaultMarkupRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.CreateDefaultMarkupRate(markup)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}

func TestDefaultMarkupRepo_Get(t *testing.T) {

	t.Run("Get_Success", func(t *testing.T) {
		// Arrange
		var markup float64 = 10

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		row := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "markup"}).
			AddRow(1, time.Now(), time.Now(), nil, markup)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WillReturnRows(row)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewDefaultMarkupRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		m, err := r.GetDefaultMarkupRate()
		assert.NoError(t, err)

		assert.NotNil(t, m)
		assert.EqualValues(t, markup, m.Markup)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func TestDefaultMarkupRepo_Delete(t *testing.T) {

	t.Run("Delete_Success", func(t *testing.T) {
		// Arrange

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), nil).
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

		r := NewDefaultMarkupRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.DeleteDefaultMarkupRate()
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func TestDefaultMarkupRepo_Update(t *testing.T) {

	t.Run("Update_Success", func(t *testing.T) {
		// Arrange
		var markup float64 = 10

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), markup).
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

		r := NewDefaultMarkupRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.UpdateDefaultMarkupRate(markup)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func TestDefaultMarkupRepo_GetHistory(t *testing.T) {

	t.Run("GetHistory_Success", func(t *testing.T) {
		// Arrange
		var markup float64 = 10

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		row := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "markup"}).
			AddRow(1, time.Now(), time.Now(), nil, markup)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
			WillReturnRows(row)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewDefaultMarkupRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		m, err := r.GetDefaultMarkupRateHistory()
		assert.NoError(t, err)

		assert.NotNil(t, m)
		assert.EqualValues(t, markup, m[0].Markup)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}
