package db_test

import (
	extsql "database/sql"
	"errors"
	"log"
	"testing"

	int_db "github.com/ukama/ukama/systems/node/software/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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


func TestSoftwareRepo_Get(t *testing.T) {
	
	t.Run("Sw  Doesn't Exist in software", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error
	
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)
	
			mock.ExpectQuery(`^SELECT.*softwares.*`).
			WillReturnError(gorm.ErrRecordNotFound)
	
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
	
		r := int_db.NewSoftwareRepo(&UkamaDbMock{
			GormDb: gdb,
		})
	
		assert.NoError(t, err)
	
		// Act
		_, err = r.GetLatestSoftwareUpdate()
	
		// Assert
		if assert.Error(t, err) {
			assert.Equal(t, true, errors.Is(gorm.ErrRecordNotFound, err))
		}
	
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	

}

// func TestControllerRepo_Add(t *testing.T) {

// 	t.Run("Add", func(t *testing.T) {
// 		// Arrange

// 		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

// 		var db *extsql.DB
// 		var err error

// 		db, mock, err := sqlmock.New() // mock sql.DB
// 		assert.NoError(t, err)

// 		mock.ExpectBegin()

// 		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
// 			WithArgs(  sqlmock.AnyArg(), sqlmock.AnyArg(),sqlmock.AnyArg(),nid.String()).
// 			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

// 		mock.ExpectCommit()

// 		dialector := postgres.New(postgres.Config{
// 			DSN:                  "sqlmock_db_0",
// 			DriverName:           "postgres",
// 			Conn:                 db,
// 			PreferSimpleProtocol: true,
// 		})
// 		gdb, err := gorm.Open(dialector, &gorm.Config{})
// 		assert.NoError(t, err)

// 		r := int_db.NewSoftwareRepo(&UkamaDbMock{
// 			GormDb: gdb,
// 		})

// 		assert.NoError(t, err)

// 		// Act
// 		err = r.CreateSoftwareUpdate(nid.String())

// 		// Assert
// 		assert.NoError(t, err)

// 		err = mock.ExpectationsWereMet()
// 		assert.NoError(t, err)
// 	})

// }





