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

		rows := sqlmock.NewRows([]string{"id", "iccid", "msisdn", "is_allocated", "sim_type", "sm_dp_address", "activation_code", "is_physical"})
		rows.AddRow(1,
			"10101010",
			"01010101",
			false,
			"inter_ukama_all",
			"123456789",
			"0000",
			true,
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs("inter_ukama_all").
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

		sp, err := r.GetSimsByType("inter_ukama_all")
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, sp)
	})
}

func Test_GetByIccid(t *testing.T) {
	t.Run("GetByIccid", func(t *testing.T) {

		var db *extsql.DB
		iccid := "10101010"
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "iccid", "msisdn", "is_allocated", "sim_type", "sm_dp_address", "activation_code", "is_physical"})
		rows.AddRow(1,
			iccid,
			"01010101",
			false,
			"inter_ukama_all",
			"123456789",
			"0000",
			true,
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

		rows := sqlmock.NewRows([]string{"id", "iccid", "msisdn", "is_allocated", "sim_type", "sm_dp_address", "activation_code", "is_physical"})
		rows.AddRow(1,
			"10101010",
			"01010101",
			false,
			"inter_ukama_all",
			"123456789",
			"0000",
			true,
		)

		mock.ExpectQuery(`^SELECT.*sims.*`).
			WithArgs(false, true, "inter_ukama_all").
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

		sp, err := r.Get(true, "inter_ukama_all")
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
				Iccid:          "10101010",
				Msisdn:         "01010101",
				IsAllocated:   false,
				SimType:       "inter_ukama_all",
				SmDpAddress:    "123456789",
				ActivationCode: "0000",
				IsPhysical:    true,
				QrCode:         "123456789",
			},
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sim[0].Iccid, sim[0].Msisdn, sim[0].Is_allocated, sim[0].Sim_type, sim[0].SmDpAddress, sim[0].ActivationCode, sim[0].QrCode, sim[0].Is_physical).
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
