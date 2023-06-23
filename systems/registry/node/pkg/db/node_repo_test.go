package db_test

import (
	extsql "database/sql"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	nodedb "github.com/ukama/ukama/systems/registry/node/pkg/db"
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

func TestNodeRepo_Add(t *testing.T) {
	var nodeID = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

	var db *extsql.DB

	node := nodedb.Node{
		Id:    nodeID.String(),
		Name:  "node-1",
		OrgId: uuid.NewV4(),
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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("AddNode", func(t *testing.T) {
		// Arrange
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(node.Name, sqlmock.AnyArg(), sqlmock.AnyArg(), node.OrgId,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), node.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectCommit()

		// Act
		err = r.Add(&node, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestNodeRepo_Get(t *testing.T) {
	var nodeID = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
	var name = "node-1"

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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("NodeFound", func(t *testing.T) {
		// Arrange
		row := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(nodeID, name)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeID).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*attached_nodes.*`).
			WithArgs(nodeID).
			WillReturnRows(row)

		// Act
		node, err := r.Get(nodeID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, node)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(nodeID).
			WillReturnError(extsql.ErrNoRows)

		// Act
		node, err := r.Get(nodeID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, node)
	})
}

func TestNodeRepo_GetForOrg(t *testing.T) {
	var orgId = uuid.NewV4()
	var nodeId = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
	var name = "node-1"

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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("OrgFound", func(t *testing.T) {
		// Arrange
		row := sqlmock.NewRows([]string{"id", "name", "org_id"}).
			AddRow(nodeId, name, orgId)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(orgId).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*attached_nodes.*`).
			WithArgs(nodeId).
			WillReturnRows(row)

		// Act
		node, err := r.GetForOrg(orgId)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, node)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("OrgNotFound", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WithArgs(orgId).
			WillReturnError(extsql.ErrNoRows)

		// Act
		node, err := r.GetForOrg(orgId)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, node)
	})
}

func TestNodeRepo_GetAll(t *testing.T) {
	var nodeID = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
	var name = "node-1"

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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("NodeFound", func(t *testing.T) {
		// Arrange
		row := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(nodeID, name)

		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WillReturnRows(row)

		mock.ExpectQuery(`^SELECT.*attached_nodes.*`).
			WithArgs(nodeID).
			WillReturnRows(row)

		// Act
		node, err := r.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, node)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("UnknownError", func(t *testing.T) {
		mock.ExpectQuery(`^SELECT.*nodes.*`).
			WillReturnError(errors.New("internal"))

		// Act
		node, err := r.GetAll()

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, node)
	})
}

func TestNodeRepo_Delete(t *testing.T) {
	var db *extsql.DB

	// Arrange
	var nodeID = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

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

	r := nodedb.NewNodeRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	assert.NoError(t, err)

	t.Run("NodeFound", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), nodeID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Act
		err = r.Delete(nodeID, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
