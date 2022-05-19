package db_test

import (
	extsql "database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	intDb "github.com/ukama/ukama/testing/services/network/internal/db"
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
	panic("implement me")
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) (err error) {
	panic("implement me")
}

func Test_vNodeRepo_GetNodeInfo(t *testing.T) {

	t.Run("GetInfoNodeExist", func(t *testing.T) {

		node := intDb.VNode{
			NodeID: "1001",
			Status: "PowerOn",
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"node_id", "status"}).
			AddRow(node.NodeID, node.Status)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(node.NodeID).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewVNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		rNode, err := r.GetInfo(node.NodeID)

		// Assert
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
		assert.NoError(t, err)

		if assert.NotNil(t, rNode) {
			assert.Equal(t, node.Status, rNode.Status)
		}
	})

}

func Test_VNodeRepo_GetInfoFailure(t *testing.T) {
	t.Run("GetInfoNodeDoesNotExist", func(t *testing.T) {

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectQuery(`SELECT`).
			WithArgs("1002").
			WillReturnError(fmt.Errorf("no matching id"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewVNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		nodeData, err := r.GetInfo("1002")

		// Assert
		if assert.Error(t, err) {
			assert.Equal(t, fmt.Errorf("no matching id"), err)
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}

		assert.NoError(t, err)
		assert.Nil(t, nodeData)
	})

}

func Test_VNodeRepo_Delete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE")).WithArgs("1001").
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

		r := intDb.NewVNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		assert.NoError(t, err)

		// Act
		err = r.Delete("1001")

		// Assert
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}

		assert.NoError(t, err)
	})

}

func Test_VNodeRepo_GetList(t *testing.T) {

	t.Run("List", func(t *testing.T) {

		node := intDb.VNode{
			NodeID: "1001",
			Status: "PowerOn",
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"node_id", "status"}).
			AddRow(node.NodeID, node.Status).
			AddRow("1002", node.Status)

		mock.ExpectQuery(`^SELECT`).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := intDb.NewVNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		nodeList, err := r.List()

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		if assert.NotNil(t, nodeList) {
			assert.Equal(t, 2, len(*nodeList))
		}
	})

}

func Test_VNodeRepo_PowerOn(t *testing.T) {
	t.Run("PowerOn", func(t *testing.T) {

		node := intDb.VNode{
			NodeID: "1001",
			Status: "PowerOn",
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), node.Status, node.NodeID).
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

		r := intDb.NewVNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		assert.NoError(t, err)

		// Act
		err = r.PowerOn(node.NodeID)

		// Assert
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}

		assert.NoError(t, err)
	})

}

func Test_VNodeRepo_PowerOff(t *testing.T) {
	t.Run("PowerOn", func(t *testing.T) {

		node := intDb.VNode{
			NodeID: "1001",
			Status: "PowerOff",
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WithArgs(sqlmock.AnyArg(), node.Status, node.NodeID).
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

		r := intDb.NewVNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		assert.NoError(t, err)

		// Act
		err = r.PowerOff(node.NodeID)

		// Assert
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}

		assert.NoError(t, err)
	})

}
