package db_test

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"

	init_db "github.com/ukama/ukama/systems/data-plan/package/pkg/db"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
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

func Test_Package_Get(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		var packageId = uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		gdb, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

		pkgRow := sqlmock.NewRows([]string{"uuid"}).
			AddRow(packageId)

		mock.ExpectQuery(`^SELECT.*package.*`).
			WithArgs(packageId).
			WillReturnRows(pkgRow)

		assert.NoError(t, err)
		repo := init_db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		// Act
		sub, err := repo.Get(packageId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sub)
	})

}

func Test_Package_GetByOrg(t *testing.T) {
	t.Run("GetByOrg", func(t *testing.T) {
		var packageId = uuid.NewV4()
		var orgId = uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		gdb, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

		pkgRow := sqlmock.NewRows([]string{"uuid", "org_id"}).
			AddRow(packageId, orgId)

		mock.ExpectQuery(`^SELECT.*package.*`).
			WithArgs(orgId).
			WillReturnRows(pkgRow)

		assert.NoError(t, err)
		repo := init_db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})
		// Act
		sub, err := repo.GetByOrg(orgId)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sub)
	})
}
func Test_Package_Delete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		packageId := uuid.NewV4()

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "packages" SET`)).
			WithArgs(sqlmock.AnyArg(), packageId).
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
		r := init_db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)
		err = r.Delete(packageId)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_Package_Add(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		var db *extsql.DB

		pkg := init_db.Package{
			Uuid:        uuid.NewV4(),
			Name:        "Monthly",
			SimType:     1,
			Active:      false,
			Duration:    360000,
			SmsVolume:   10,
			DataVolume:  1024,
			VoiceVolume: 10,
			OrgRatesID:  1,
			OrgID:       uuid.NewV4(),
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), pkg.Uuid.String(), pkg.Name, pkg.SimType, pkg.OrgID,
				pkg.Active, pkg.Duration, pkg.SmsVolume, pkg.DataVolume, pkg.VoiceVolume, pkg.OrgRatesID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := init_db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		err = r.Add(&pkg)
		assert.NotNil(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
