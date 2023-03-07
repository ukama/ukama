package db_test

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"

	int_db "github.com/ukama/ukama/systems/data-plan/package/pkg/db"

	uuid "github.com/ukama/ukama/systems/common/uuid"

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

func Test_Package_Get(t *testing.T) {

	t.Run("PackageExist", func(t *testing.T) {
		// Arrange
		const uuidStr = "51fbba62-c79f-11eb-b8bc-0242ac130003"
		packID, _ := uuid.FromString(uuidStr)
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		id := uuid.NewV4()
		pack := &int_db.Package{
			Uuid:        packID,
			Name:        "Silver Plan",
			SimType:     int_db.SimTypeTest,
			OrgId:       uuid.NewV4(),
			Active:      true,
			Duration:    30,
			SmsVolume:   1000,
			DataVolume:  5000000,
			VoiceVolume: 500,
			OrgRatesId:  1,
		}

		rows := sqlmock.NewRows([]string{"package_id", "name", "sim_type", "org_id", "active", "duration", "sms_volume", "data_volume", "voice_volume", "org_rate_id"}).
			AddRow(packID, pack.Name, pack.SimType, pack.OrgId, pack.Active, pack.Duration, pack.SmsVolume, pack.DataVolume, pack.VoiceVolume, pack.OrgRatesId)

		mock.ExpectQuery(`^SELECT.*packages.*`).
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

		r := int_db.NewPackageRepo(&UkamaDbMock{
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

func Test_Package_GetByOrg(t *testing.T) {

	t.Run("PackageExist", func(t *testing.T) {
		// Arrange
		const uuidStr = "51fbba62-c79f-11eb-b8bc-0242ac130003"
		packID, _ := uuid.FromString(uuidStr)
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		id := uuid.NewV4()
		pack := &int_db.Package{
			Uuid:        packID,
			Name:        "Silver Plan",
			SimType:     int_db.SimTypeTest,
			OrgId:       uuid.NewV4(),
			Active:      true,
			Duration:    30,
			SmsVolume:   1000,
			DataVolume:  5000000,
			VoiceVolume: 500,
			OrgRatesId:  1,
		}

		rows := sqlmock.NewRows([]string{"package_id", "name", "sim_type", "org_id", "active", "duration", "sms_volume", "data_volume", "voice_volume", "org_rate_id"}).
			AddRow(packID, pack.Name, pack.SimType, pack.OrgId, pack.Active, pack.Duration, pack.SmsVolume, pack.DataVolume, pack.VoiceVolume, pack.OrgRatesId)

		mock.ExpectQuery(`^SELECT.*packages.*`).
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

		r := int_db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		node, err := r.GetByOrg(id)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, node)
	})

}

func Test_Package_Add(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		var db *extsql.DB

		pkg := int_db.Package{
			Uuid:        uuid.NewV4(),
			Name:        "Monthly",
			SimType:     int_db.SimTypeUkamaData,
			Active:      false,
			Duration:    360000,
			SmsVolume:   10,
			DataVolume:  1024,
			VoiceVolume: 10,
			OrgRatesId:  1,
			OrgId:       uuid.NewV4(),
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), pkg.Uuid, pkg.Name, pkg.SimType, pkg.OrgId, pkg.Active, pkg.Duration, pkg.SmsVolume, pkg.DataVolume, pkg.VoiceVolume, pkg.OrgRatesId).
			WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(pkg.Uuid))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewPackageRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		err = r.Add(&pkg)
		assert.NotNil(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
