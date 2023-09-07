package db_test

import (
	extsql "database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	sidedb "github.com/ukama/ukama/systems/registry/node/pkg/db"
)

func TestSiteRepo_GetNodes(t *testing.T) {
	// Arrange
	var siteID = uuid.NewV4()
	var nodeIDa = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
	var nodeIDb = ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

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

	r := sidedb.NewSiteRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	t.Run("SiteFound", func(t *testing.T) {
		siteRows := sqlmock.NewRows([]string{"node_id", "site_id"}).
			AddRow(nodeIDa, siteID).
			AddRow(nodeIDb, siteID)

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(siteID).
			WillReturnRows(siteRows)

		// Act
		nodes, err := r.GetNodes(siteID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, nodes)
		assert.Equal(t, 2, len(nodes))

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SiteNotFound", func(t *testing.T) {
		// Arrange

		mock.ExpectQuery(`^SELECT.*sites.*`).
			WithArgs(siteID).
			WillReturnError(extsql.ErrNoRows)

		// Act
		nodes, err := r.GetNodes(siteID)

		// Assert
		assert.Error(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Nil(t, nodes)
	})
}
