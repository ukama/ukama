package db_test

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"

	"github.com/google/uuid"
	org_db "github.com/ukama/ukama/systems/registry/org/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
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

func Test_OrgRepo_Get(t *testing.T) {
	t.Run("OrgExist", func(t *testing.T) {
		// Arrange
		const orgId = 1
		const orgName = "ukama"
		const orgCert = "ukamacert"

		var orgOwner = uuid.New()

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "owner", "certificate"}).
			AddRow(orgId, orgName, orgOwner, orgCert)

		mock.ExpectQuery(`^SELECT.*orgs.*`).
			WithArgs(orgId).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := org_db.NewOrgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		org, err := r.Get(orgId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, org)
	})
}

func Test_OrgRepo_Add(t *testing.T) {
	t.Run("AddOrg", func(t *testing.T) {
		// Arrange
		var db *extsql.DB

		org := org_db.Org{
			Name:        "ukama",
			Owner:       uuid.New(),
			Certificate: "ukama_certs",
		}

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(org.Name, org.Owner, org.Certificate, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

		r := org_db.NewOrgRepo(&UkamaDbMock{
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

// func Test_OrgRepo_Add(t *testing.T) {
// t.Run("AddOrg", func(t *testing.T) {
// // Arrange
// org := org_db.Org{
// Name:        "ukama",
// Owner:       uuid.New(),
// Certificate: "ukama_certs",
// }

// var db *extsql.DB

// db, mock, err := sqlmock.New() // mock sql.DB
// assert.NoError(t, err)

// mock.ExpectBegin()

// mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
// WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), org.Name, org.Owner, org.Certificate).
// WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

// mock.ExpectCommit()

// dialector := postgres.New(postgres.Config{
// DSN:                  "sqlmock_db_0",
// DriverName:           "postgres",
// Conn:                 db,
// PreferSimpleProtocol: true,
// })

// gdb, err := gorm.Open(dialector, &gorm.Config{})
// assert.NoError(t, err)

// r := org_db.NewOrgRepo(&UkamaDbMock{
// GormDb: gdb,
// })

// assert.NoError(t, err)

// // Act
// err = r.Add(&org)

// // Assert
// assert.NoError(t, err)

// err = mock.ExpectationsWereMet()
// assert.NoError(t, err)
// })
// }
