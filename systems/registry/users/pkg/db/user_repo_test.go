package db_test

import (
	"regexp"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/uuid"

	"database/sql"
	extsql "database/sql"

	log "github.com/sirupsen/logrus"
	userdb "github.com/ukama/ukama/systems/registry/users/pkg/db"
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

func TestUserRepo_Add(t *testing.T) {
	// Arrange
	var db *extsql.DB

	user := userdb.User{
		Id:     uuid.NewV4(),
		Name:   "John Doe",
		Email:  "johndoe@example.com",
		Phone:  "00100000000",
		AuthId: uuid.NewV4(),
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

	r := userdb.NewUserRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("AddUser", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(user.Id, user.Name, user.Email, user.Phone, sqlmock.AnyArg(),
				user.AuthId, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

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

	r := userdb.NewUserRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("UserFound", func(t *testing.T) {
		// Arrange
		const name = "John Doe"
		const email = "johndoe@example.com"
		const phone = "00100000000"
		var userId = uuid.NewV4()
		var authId = uuid.NewV4()

		rows := sqlmock.NewRows([]string{"id", "name", "email", "phone", "auth_id"}).
			AddRow(userId, name, email, phone, authId)

		mock.ExpectQuery(`^SELECT.*users.*`).
			WithArgs(userId).
			WillReturnRows(rows)

		// Act
		usr, err := r.Get(userId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, usr)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		assert.Equal(t, name, usr.Name)
		assert.Equal(t, email, usr.Email)
		assert.Equal(t, phone, usr.Phone)
		assert.Equal(t, authId, usr.AuthId)
	})

	t.Run("userNotFound", func(t *testing.T) {
		// Arrange
		var userId = uuid.NewV4()

		mock.ExpectQuery(`^SELECT.*users.*`).
			WithArgs(userId).
			WillReturnError(sql.ErrNoRows)

		// Act
		usr, err := r.Get(userId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, usr)
	})
}

func TestUserRepo_GetByAuthId(t *testing.T) {
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

	r := userdb.NewUserRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("UserFound", func(t *testing.T) {
		// Arrange
		const name = "John Doe"
		const email = "johndoe@example.com"
		const phone = "00100000000"
		var userId = uuid.NewV4()
		var authId = uuid.NewV4()

		rows := sqlmock.NewRows([]string{"id", "name", "email", "phone", "auth_id"}).
			AddRow(userId, name, email, phone, authId)

		mock.ExpectQuery(`^SELECT.*users.*`).
			WithArgs(authId).
			WillReturnRows(rows)

		// Act
		usr, err := r.GetByAuthId(authId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, usr)

		assert.NotNil(t, name, usr.Name)
		assert.NotNil(t, email, usr.Email)
		assert.NotNil(t, phone, usr.Phone)
		assert.NotNil(t, authId, usr.AuthId)
	})

	t.Run("userNotFound", func(t *testing.T) {
		// Arrange
		var authId = uuid.NewV4()

		mock.ExpectQuery(`^SELECT.*users.*`).
			WithArgs(authId).
			WillReturnError(sql.ErrNoRows)

		// Act
		usr, err := r.GetByAuthId(authId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, usr)
	})
}

// func TestUserRepo_Update(t *testing.T) {
// t.Run("UserFound", func(t *testing.T) {
// var db *extsql.DB

// const name = "John Doe"
// const email = "johndoe@example.com"
// const phone = "00100000000"
// var authId = uuid.NewV4()

// var userId = uuid.NewV4()

// usr := &userdb.User{
// Id:     userId,
// Name:   "Fox Doe",
// Email:  "foxdoe@example.com",
// Phone:  "00200000000",
// AuthId: uuid.NewV4(),
// }

// db, mock, err := sqlmock.New() // mock sql.DB
// assert.NoError(t, err)

// rows := sqlmock.NewRows([]string{"id", "name", "email", "phone", "auth_id"}).
// AddRow(userId, name, email, phone, authId)

// mock.ExpectBegin()

// mock.ExpectQuery(`^SELECT.*users.*`).
// WithArgs(userId).
// WillReturnRows(rows)

// mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
// WithArgs(usr.Name, usr.Email, usr.Phone, usr.AuthId,
// sqlmock.AnyArg(), usr.Id, usr.Id).
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

// r := userdb.NewUserRepo(&UkamaDbMock{
// GormDb: gdb,
// })

// assert.NoError(t, err)

// // Act
// err = r.Update(usr, nil)

// // Assert
// assert.NoError(t, err)

// err = mock.ExpectationsWereMet()
// assert.NoError(t, err)
// })

// // t.Run("UserNotFound", func(t *testing.T) {
// // var db *extsql.DB

// // var userId = uuid.NewV4()

// // db, mock, err := sqlmock.New() // mock sql.DB
// // assert.NoError(t, err)

// // mock.ExpectBegin()

// // mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
// // WithArgs(sqlmock.AnyArg(), userId).
// // WillReturnError(sql.ErrNoRows)

// // dialector := postgres.New(postgres.Config{
// // DSN:                  "sqlmock_db_0",
// // DriverName:           "postgres",
// // Conn:                 db,
// // PreferSimpleProtocol: true,
// // })

// // gdb, err := gorm.Open(dialector, &gorm.Config{})
// // assert.NoError(t, err)

// // r := userdb.NewUserRepo(&UkamaDbMock{
// // GormDb: gdb,
// // })

// // assert.NoError(t, err)

// // // Act
// // err = r.Delete(userId, nil)

// // // Assert
// // assert.Error(t, err)

// // err = mock.ExpectationsWereMet()
// // assert.NoError(t, err)
// // })
// }

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

	r := userdb.NewUserRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("UserFound", func(t *testing.T) {
		var userId = uuid.NewV4()

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
			WithArgs(sqlmock.AnyArg(), userId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err = r.Delete(userId, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		var userId = uuid.NewV4()

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
			WithArgs(sqlmock.AnyArg(), userId).
			WillReturnError(sql.ErrNoRows)

		// Act
		err = r.Delete(userId, nil)

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

	r := userdb.NewUserRepo(&UkamaDbMock{
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
