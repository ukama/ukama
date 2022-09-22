package db_test

import (
	extsql "database/sql"
	"regexp"
	"testing"

	"github.com/jackc/pgtype"

	int_db "github.com/ukama/ukama/systems/init/lookup/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_orgRepo_Get(t *testing.T) {

	t.Run("OrgExist", func(t *testing.T) {
		// Arrange
		const orgName = "ukama"
		const orgCert = "ukamacert"
		const ip = "0.0.0.0"

		var dIp pgtype.Inet
		err := dIp.Set(ip)
		assert.NoError(t, err)

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"name", "certificate", "ip"}).
			AddRow(orgName, orgCert, dIp)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(orgName).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewOrgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		org, err := r.GetByName(orgName)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, org)
	})

}

func Test_orgRepo_Add(t *testing.T) {

	t.Run("AddOrg", func(t *testing.T) {
		// Arrange
		const ip = "0.0.0.0"

		var dIp pgtype.Inet
		err := dIp.Set(ip)
		assert.NoError(t, err)

		org := int_db.Org{
			Name:        "ukama",
			Certificate: "ukama_certs",
			Ip:          dIp,
		}

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), org.Name, org.Certificate, org.Ip).
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

		r := int_db.NewOrgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&org)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

}
