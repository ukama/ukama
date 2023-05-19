package db_test

import (
	extsql "database/sql"
	"regexp"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/uuid"
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
		const orgName = "ukama"
		const orgCert = "ukamacert"

		var orgId = uuid.NewV4()
		var orgOwner = uuid.NewV4()

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

// func Test_OrgRepo_Add(t *testing.T) {
// t.Run("AddOrg", func(t *testing.T) {
// // Arrange
// var db *extsql.DB

// org := org_db.Org{
// Name:        "ukama",
// Owner:       uuid.NewV4(),
// Certificate: "ukama_certs",
// }

// db, mock, err := sqlmock.New() // mock sql.DB
// assert.NoError(t, err)

// mock.ExpectBegin()

// mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
// WithArgs(org.Name, org.Owner, org.Certificate,
// sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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
// err = r.Add(&org, nil)

// // Assert
// assert.NoError(t, err)

// err = mock.ExpectationsWereMet()
// assert.NoError(t, err)
// })
// }

func Test_OrgRepo_AddMember(t *testing.T) {
	t.Run("AddMember", func(t *testing.T) {
		// Arrange
		member := org_db.OrgUser{
			OrgId:       uuid.NewV4(),
			UserId:      1,
			Uuid:        uuid.NewV4(),
			Role:        org_db.Member,
			Deactivated: false,
			CreatedAt:   time.Now(),
		}

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(member.OrgId, member.UserId, member.Uuid, member.Deactivated, member.CreatedAt, sqlmock.AnyArg(), org_db.RoleType(org_db.Member)).
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

		r := org_db.NewOrgRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		// Act
		err = r.AddMember(&member)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_OrgRepo_GetMember(t *testing.T) {
	t.Run("MemberExist", func(t *testing.T) {
		// Arrange

		orgID := uuid.NewV4()
		userUUID := uuid.NewV4()

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"org_id", "uuid"}).
			AddRow(orgID, userUUID)

		mock.ExpectQuery(`^SELECT.*org_users.*`).
			WithArgs(orgID, userUUID).
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
		member, err := r.GetMember(orgID, userUUID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, member)
	})
}

func Test_OrgRepo_GetMembers(t *testing.T) {
	t.Run("MembersOfAnOrg", func(t *testing.T) {
		// Arrange

		orgID := uuid.NewV4()

		var db *extsql.DB

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"org_id"}).
			AddRow(orgID)

		mock.ExpectQuery(`^SELECT.*org_users.*`).
			WithArgs(orgID).
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
		members, err := r.GetMembers(orgID)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, members)
	})
}
