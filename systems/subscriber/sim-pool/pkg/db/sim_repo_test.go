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

func Test_GetSimsByType(t *testing.T) {
	t.Run("GetStats", func(t *testing.T) {

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})
		rows.AddRow(
			"8910300000003540855",
			"01010101",
			false,
			false,
			SimTypeTest,
			"123456789",
			"0000",
			true,
			"123456789",
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(SimTypeTest).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		sp, err := r.GetSimsByType(SimTypeTest)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sp)
	})
}

func Test_GetByIccid(t *testing.T) {
	t.Run("GetByIccid", func(t *testing.T) {

		var db *extsql.DB
		iccid := "8910300000003540855"
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})
		rows.AddRow(
			iccid,
			"01010101",
			false,
			false,
			SimTypeTest,
			"123456789",
			"0000",
			true,
			"123456789",
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(false, iccid).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		sp, err := r.GetByIccid(iccid)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sp)
	})
}

func Test_Get(t *testing.T) {
	t.Run("Get", func(t *testing.T) {

		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"iccid", "msisdn", "is_allocated", "is_failed", "sim_type", "sm_dp_address", "activation_code", "is_physical", "qr_code"})
		rows.AddRow(
			"8910300000003540855",
			"01010101",
			false,
			false,
			SimTypeTest,
			"123456789",
			"0000",
			true,
			"123456789",
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(false, true, SimTypeTest).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		sp, err := r.Get(true, SimTypeTest)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sp)
	})
}

func Test_Add(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		var db *extsql.DB

		sim := []Sim{
			{
				Iccid:          "8910300000003540855",
				Msisdn:         "01010101",
				IsAllocated:    false,
				IsFailed:       false,
				SimType:        SimTypeTest,
				SmDpAddress:    "123456789",
				ActivationCode: "0000",
				IsPhysical:     true,
				QrCode:         "123456789",
			},
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sim[0].Iccid, sim[0].Msisdn, sim[0].IsAllocated, sim[0].IsFailed, sim[0].SimType, sim[0].SmDpAddress, sim[0].ActivationCode, sim[0].QrCode, sim[0].IsPhysical).
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

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		err = r.Add(sim)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_Delete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		simId := []uint64{1}
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "sims" SET`)).
			WithArgs(sqlmock.AnyArg(), simId[0]).
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

		r := NewSimRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		err = r.Delete(simId)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
