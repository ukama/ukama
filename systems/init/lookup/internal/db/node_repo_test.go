package db_test

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"

	int_db "github.com/ukama/ukama/systems/init/lookup/internal/db"

	"github.com/ukama/ukama/systems/common/ukama"

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

func Test_nodeRepo_Get(t *testing.T) {

	t.Run("NodeExist", func(t *testing.T) {
		// Arrange
		const uuidStr = "51fbba62-c79f-11eb-b8bc-0242ac130003"
		const orgId = uint(15)

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		id := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		rows := sqlmock.NewRows([]string{"node_id", "orgid"}).
			AddRow(uuidStr, orgId)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(id).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		node, err := r.Get(id)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, node)
	})

}

func Test_nodeRepo_Delete(t *testing.T) {

	t.Run("DeleteNode", func(t *testing.T) {

		id := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("DELETE")).WithArgs(id).
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

		r := int_db.NewNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Delete(id)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}

func Test_nodeRepo_Add(t *testing.T) {

	t.Run("AddNode", func(t *testing.T) {
		// Arrange
		//const uuidStr = "51fbba62-c79f-11eb-b8bc-0242ac130003"
		const orgId = uint(15)

		nid := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		node := int_db.Node{
			NodeID: nid.StringLowercase(),
			OrgID:  orgId,
		}

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), node.NodeID, node.OrgID).
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

		r := int_db.NewNodeRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.AddOrUpdate(&node)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}
