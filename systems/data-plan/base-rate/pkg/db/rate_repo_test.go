package db

import (
	extsql "database/sql"
	"log"
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

		// Act
		rate, err := r.GetBaseRate(rateId)

		// Assert
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

		// Act
		rates, err := r.GetBaseRates("Tycho crater", "", "", "inter_mno_data")

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, rates)
	})
}
