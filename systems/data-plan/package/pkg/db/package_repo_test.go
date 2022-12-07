package db

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"

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

func Test_Package_Get(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		const packageId = 1

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "org_id", "active", "duration", "sms_volume",
			"data_volume", "voice_volume", "sim_type", "org_rate_id"}).
			AddRow(packageId, "Monthly Super", 1, "t", 360000, 10, 1024, 10, "inter_ukama_all", 1)

		mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs(packageId).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		pkg, err := r.Get(packageId)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, pkg)
	})
}
func Test_Package_GetByOrg(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		const orgId = 1

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "name", "org_id", "active", "duration", "sms_volume",
			"data_volume", "voice_volume", "sim_type", "org_rate_id"}).
			AddRow(1, "Monthly Super", orgId, "t", 360000, 10, 1024, 10, "inter_ukama_all", 1)

		mock.ExpectQuery(`^SELECT.*packages.*`).
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

		r := NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		pkg, err := r.GetByOrg(orgId)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, pkg)
	})
}
func Test_Package_Delete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		packageId := 1

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "packages" SET`)).
			WithArgs(sqlmock.AnyArg(), uint64(packageId)).
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
		r := NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)
		err = r.Delete(uint64(packageId))
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_Package_Update(t *testing.T) {
	t.Run("Update", func(t *testing.T) {
		packageId := 1

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectBegin()

		mock.ExpectExec("UPDATE").WithArgs("Monthly", "inter_ukama_all", 360000, 10, 1024, 10, 1, packageId).
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
		r := NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		_package := Package{
			Name:         "Monthly",
			Sim_type:     "inter_ukama_all",
			Active:       false,
			Duration:     360000,
			Sms_volume:   10,
			Data_volume:  1024,
			Voice_volume: 10,
			Org_rates_id: 1,
		}

		assert.NoError(t, err)
		_, err = r.Update(uint64(packageId), _package)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_Package_Add(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		var db *extsql.DB

		pkg := Package{
			Name:         "Monthly",
			Sim_type:     "inter_ukama_all",
			Active:       false,
			Duration:     360000,
			Sms_volume:   10,
			Data_volume:  1024,
			Voice_volume: 10,
			Org_rates_id: 1,
			Org_id:       1,
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), pkg.Name, pkg.Sim_type, pkg.Org_id,
				pkg.Active, pkg.Duration, pkg.Sms_volume, pkg.Data_volume, pkg.Voice_volume, pkg.Org_rates_id).
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

		r := NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		err = r.Add(&pkg)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
