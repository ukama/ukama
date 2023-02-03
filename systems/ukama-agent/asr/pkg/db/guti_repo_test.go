package db_test

import (
	extsql "database/sql"
	"regexp"
	"testing"
	"time"

	int_db "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

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
var Guti = int_db.Guti{
	Imsi:            Imsi,
	PlmnId:          "00101",
	Mmegi:           101,
	Mmec:            101,
	MTmsi:           101,
	DeviceUpdatedAt: time.Unix(int64(1639144056), 0),
}

func TestGutiRepo_Update(t *testing.T) {

	t.Run("UpdateGuti", func(t *testing.T) {
		// Arrange
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"created_at", "device_updated_at", "imsi", "plmn_id", "mmegi", "mmec", "m_tmsi"})

		mock.ExpectBegin()
		mock.ExpectQuery(`^SELECT.*gutis.*`).
			WithArgs(Imsi, sqlmock.AnyArg()).
			WillReturnRows(rows)

		mock.ExpectExec(regexp.QuoteMeta("DELETE")).
			WithArgs(Imsi, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), Guti.Imsi, Guti.PlmnId, Guti.Mmegi, Guti.Mmec, Guti.MTmsi).
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

		r := int_db.NewGutiRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		err = r.Update(&Guti)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

	})

}

func TestGutiRepo_GetImsi(t *testing.T) {

	t.Run("GetImsi", func(t *testing.T) {
		var db *extsql.DB
		var err error

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)
		now := time.Now()
		rows := sqlmock.NewRows([]string{"created_at", "device_updated_at", "imsi", "plmn_id", "mmegi", "mmec", "m_tmsi"}).
			AddRow(now, now, Guti.Imsi, Guti.PlmnId, Guti.Mmegi, Guti.Mmec, Guti.MTmsi)

		mock.ExpectQuery(`^SELECT.*gutis.*`).
			WithArgs(Imsi).
			WillReturnRows(rows)
		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})
		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := int_db.NewGutiRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		// Act
		i, err := r.GetImsi(Imsi)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)

		if assert.NotNil(t, i) {
			assert.EqualValues(t, i, Imsi)
		}

	})

}
