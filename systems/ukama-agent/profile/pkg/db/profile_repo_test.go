package db_test

import (
	extsql "database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/uuid"
	int_db "github.com/ukama/ukama/systems/ukama-agent/profile/pkg/db"
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

	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) (err error) {
	return u.GormDb.Transaction(func(tx *gorm.DB) error {
		d := dbOperation(tx)

		if d.Error != nil {
			return d.Error
		}

		if len(nestedFuncs) > 0 {
			for _, n := range nestedFuncs {
				if n != nil {
					nestErr := n(tx)
					if nestErr != nil {
						return nestErr
					}
				}
			}
		}

		return nil
	})
}

var Imsi = "012345678912345"
var Iccid = "123456789123456789123"
var Network = "db081ef5-a8ae-4a95-bff3-a7041d52bb9b"
var Org = "abdc0cec-5553-46aa-b3a8-1e31b0ef58ad"
var Package = "fab4f98d-2e82-47e8-adb5-e516346880d8"

var profile = int_db.Profile{
	Iccid:                   Iccid,
	Imsi:                    Imsi,
	UeDlBps:                 10000000,
	UeUlBps:                 1000000,
	ApnName:                 "ukama",
	AllowedTimeOfService:    2592000,
	TotalDataBytes:          1024000,
	ConsumedDataBytes:       0,
	NetworkId:               uuid.FromStringOrNil(Network),
	PackageId:               uuid.FromStringOrNil(Package),
	LastStatusChangeReasons: int_db.ACTIVATION,
	LastStatusChangeAt:      time.Now(),
}

var pack = int_db.PackageDetails{
	PackageId:            uuid.FromStringOrNil(Package),
	UeDlBps:              10000000,
	UeUlBps:              1000000,
	ApnName:              "ukama",
	AllowedTimeOfService: time.Second * 2592000,
	TotalDataBytes:       1024000,
	ConsumedDataBytes:    0,
	LastStatusChangeAt:   time.Now(),
}

func TestProfileRepo_Add(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		//row := sqlmock.NewRows([]string{"iccid", "imsi", "op", "amf", "key", "algo_type", "ue_dl_ambr_bps", "ue_ul_ambr_bps", "sqn", "csg_id_prsent", "csg_id", "default_apn_name", "network_id", "package_id"}).

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), profile.Iccid, profile.Imsi, profile.UeDlBps, profile.UeUlBps, profile.ApnName, profile.NetworkId, profile.PackageId, profile.AllowedTimeOfService, profile.TotalDataBytes, profile.ConsumedDataBytes, profile.LastStatusChangeAt, profile.LastStatusChangeReasons).
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

		r := int_db.NewProfileRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Add(&profile)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestProfileRepo_UpdatePackage(t *testing.T) {
	t.Run("UpdatePackage", func(t *testing.T) {
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(int64(pack.AllowedTimeOfService.Seconds()), pack.ApnName, pack.ConsumedDataBytes, pack.LastStatusChangeAt, int_db.PACKAGE_UPDATE, pack.PackageId, pack.TotalDataBytes, pack.UeDlBps, pack.UeUlBps, sqlmock.AnyArg(), Imsi).
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

		r := int_db.NewProfileRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.UpdatePackage(Imsi, pack)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestProfileRepo_UpdateUsage(t *testing.T) {
	t.Run("UpdatePackage", func(t *testing.T) {
		var db *extsql.DB
		var err error
		var usage uint64 = 1000
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
			WithArgs(sqlmock.AnyArg(), usage, Imsi).
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

		r := int_db.NewProfileRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.UpdateUsage(Imsi, usage)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
