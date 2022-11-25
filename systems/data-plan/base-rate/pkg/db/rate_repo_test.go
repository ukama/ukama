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

func Test_Rate_Get(t *testing.T) {
	t.Run("GetRate", func(t *testing.T) {
		const rateId = 1

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "country", "network", "vpmn", "imsi", "sms_mo", "sms_mt", "data", "lte", "lte_m", "apn", "end_at", "x2g", "x3g", "x5g", "sim_type", "effective_at"}).
			AddRow(1,
				"Tycho crater",
				"Multi Tel",
				"TTC",
				"1",
				"$0.1",
				"$0.1",
				"$0.4",
				"LTE",
				"",
				"Manual entry required",
				"",
				"2G",
				"3G",
				"",
				"inter_mno_data",
				"2023-10-10",
			)

		mock.ExpectQuery(`^SELECT.*rates.*`).
			WithArgs(rateId).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewBaseRateRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		rate, err := r.GetBaseRate(rateId)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, rate)
	})
}

func Test_Rates_Get(t *testing.T) {
	t.Run("GetRates", func(t *testing.T) {
		var db *extsql.DB
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "country", "network", "vpmn", "imsi", "sms_mo", "sms_mt", "data", "lte", "lte_m", "apn", "end_at", "x2g", "x3g", "x5g", "sim_type", "effective_at"})
		for i := 1; i <= 3; i++ {
			rows.AddRow(i,
				"Tycho crater",
				"Multi Tel",
				"TTC",
				"1",
				"$0.1",
				"$0.1",
				"$0.4",
				"LTE",
				"",
				"Manual entry required",
				"",
				"2G",
				"3G",
				"",
				"inter_mno_data",
				"2023-10-10",
			)
		}

		mock.ExpectQuery(`^SELECT.*rates.*`).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := NewBaseRateRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		rates, err := r.GetBaseRates("Tycho crater", "", "", "inter_mno_data")
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, rates)
	})
}

func Test_Rate_Upload(t *testing.T) {
	t.Run("UploadRates", func(t *testing.T) {
		var db *extsql.DB

		rates := []Rate{{
			Country:      "Tycho crater",
			Data:         "$0.4",
			Effective_at: "2023-10-10",
			Network:      "Multi Tel",
			Sim_type:     "inter_mno_data",
			X2g:          "",
			X3g:          "",
			Apn:          "",
			Imsi:         "",
			Lte:          "",
			Sms_mo:       "",
			Sms_mt:       "",
			Vpmn:         "",
			End_at:       "",
			Lte_m:        "",
			X5g:          "",
		}}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), rates[0].Country, rates[0].Network,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), rates[0].Data,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), rates[0].Effective_at, sqlmock.AnyArg(), rates[0].Sim_type).
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

		r := NewBaseRateRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.UploadBaseRates(rates)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
