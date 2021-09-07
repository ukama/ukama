package db_test

import (
	extsql "database/sql"
	"testing"
	int_db "ukamaX/bootstrap/lookup/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	uuid "github.com/satori/go.uuid"
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

func Test_nodeRepo_Get(t *testing.T) {

	t.Run("NodeExist", func(t *testing.T) {
		// Arrange
		const uuidStr = "51fbba62-c79f-11eb-b8bc-0242ac130003"
		const orgId = uint(15)

		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		id, err := uuid.FromString(uuidStr)
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"uuid", "orgid"}).
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
