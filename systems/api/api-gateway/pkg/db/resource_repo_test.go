package db_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/api/api-gateway/pkg/db"
	"github.com/ukama/ukama/systems/common/uuid"
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

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	panic("implement me")
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) (err error) {
	panic("implement me")
}

func TestResourceRepo_Add(t *testing.T) {
	var err error

	nt := db.Resource{
		Id:     uuid.NewV4(),
		Status: db.ParseStatus("pending"),
	}

	mock, gdb := prepare_db(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
		WithArgs(nt.Id, nt.Status, sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r := db.NewResourceRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	assert.NoError(t, err)

	// Act
	err = r.Add(&nt)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestResourceRepo_Get(t *testing.T) {
	t.Run("ResourceFound", func(t *testing.T) {
		// Arrange
		var resourceId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		rows := sqlmock.NewRows([]string{"id", "status"}).
			AddRow(resourceId, db.ParseStatus("completed"))

		mock.ExpectQuery(`^SELECT.*resources.*`).
			WithArgs(resourceId).
			WillReturnRows(rows)

		r := db.NewResourceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		resource, err := r.Get(resourceId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resource)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ResourceNotFound", func(t *testing.T) {
		// Arrange
		var resourceId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectQuery(`^SELECT.*resources.*`).
			WithArgs(resourceId).
			WillReturnError(sql.ErrNoRows)

		r := db.NewResourceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		resource, err := r.Get(resourceId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resource)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
func TestResourceRepo_Update(t *testing.T) {
	var err error

	nt := db.Resource{
		Id:     uuid.NewV4(),
		Status: db.ParseStatus("pending"),
	}

	mock, gdb := prepare_db(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
		WithArgs(nt.Status, sqlmock.AnyArg(), nt.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r := db.NewResourceRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	assert.NoError(t, err)

	// Act
	err = r.Update(&nt)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestResourceRepo_Delete(t *testing.T) {
	t.Run("ResourceFound", func(t *testing.T) {
		// Arrange
		var resourceId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "resources" SET`)).
			WithArgs(sqlmock.AnyArg(), resourceId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		r := db.NewResourceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(resourceId)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ResourceNotFound", func(t *testing.T) {
		// Arrange
		var resourceId = uuid.NewV4()

		mock, gdb := prepare_db(t)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "resources" SET`)).
			WithArgs(sqlmock.AnyArg(), resourceId).
			WillReturnError(sql.ErrNoRows)

		r := db.NewResourceRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err := r.Delete(resourceId)

		// Assert
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func prepare_db(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	var db *sql.DB
	var err error

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

	return mock, gdb
}
