package db_test

import (
	extsql "database/sql"
	"errors"
	"log"
	"testing"

	"github.com/ukama/ukama/systems/common/ukama"
	int_db "github.com/ukama/ukama/systems/node/health/pkg/db"

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


func TestHealthRepo_GetRunningApps(t *testing.T) {

	t.Run("apps exit in health", func(t *testing.T) {
		// Arrange
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
	
		var db *extsql.DB
		var err error
	
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)
	
		rows := sqlmock.NewRows([]string{"node_id"}).
			AddRow(nid.String())
	
		
			mock.ExpectQuery(`^SELECT.*healths.*`).
			WithArgs(nid.String()).
			WillReturnRows(rows)
	
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
	
		r := int_db.NewHealthRepo(&UkamaDbMock{
			GormDb: gdb,
		})
	
		assert.NoError(t, err)
	
		// Act
		c, err := r.GetRunningAppsInfo(nid)
	
		// Assert
		assert.NoError(t, err)
	
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NoError(t, err)
		if assert.NotNil(t, c) {
			assert.Equal(t, nid.String(),c.NodeId)
		}
	})
	
	t.Run("Apps Doesn't Exist in health", func(t *testing.T) {
		// Arrange
		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		var db *extsql.DB
		var err error
	
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)
	
			mock.ExpectQuery(`^SELECT.*healths.*`).
			WithArgs(nid.String()).
			WillReturnError(gorm.ErrRecordNotFound)
	
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)
	
		r := int_db.NewHealthRepo(&UkamaDbMock{
			GormDb: gdb,
		})
	
		assert.NoError(t, err)
	
		// Act
		_, err = r.GetRunningAppsInfo(nid)
	
		// Assert
		if assert.Error(t, err) {
			assert.Equal(t, true, errors.Is(gorm.ErrRecordNotFound, err))
		}
	
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
	

}









func TestHealthRepo_StoreRunningAppsInfo(t *testing.T) {

	t.Run("StoreRunningAppsInfo", func(t *testing.T) {
		// Arrange

		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		expectedQuery := `INSERT INTO "healths" (.+)`
		mock.ExpectExec(expectedQuery).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()   
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewHealthRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)
		nestedFunc := func(string, string) error {
			// Implement the behavior of the nested function here (if needed)
			return nil
		}
		// Act
		err = r.StoreRunningAppsInfo(
			&int_db.Health{
				NodeId: nid.String(),
				TimeStamp: "12-12-2024",
			},
			nestedFunc,)


		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}