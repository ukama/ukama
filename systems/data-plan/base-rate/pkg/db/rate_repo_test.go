package db

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	uuid "github.com/ukama/ukama/systems/common/uuid"
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
		var rateId = uuid.NewV4()

		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "uuid", "country", "network", "vpmn", "imsi", "sms_mo", "sms_mt", "data", "lte", "lte_m", "apn", "end_at", "x2g", "x3g", "x5g", "sim_type", "effective_at"}).
			AddRow(1,
				rateId,
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
				3,
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

		rows := sqlmock.NewRows([]string{"id", "uuid", "country", "network", "vpmn", "imsi", "sms_mo", "sms_mt", "data", "lte", "lte_m", "apn", "end_at", "x2g", "x3g", "x5g", "sim_type", "effective_at"})
		for i := 1; i <= 3; i++ {
			rows.AddRow(i,
				uuid.NewV4(),
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
				3,
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

		rates, err := r.GetBaseRates("Tycho crater", "", "", SimTypeUkamaData)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, rates)
	})
}

func Test_Rate_Add(t *testing.T) {
	var rate_uuid = uuid.NewV4()
	t.Run("AddRates", func(t *testing.T) {
		var db *extsql.DB

		rates := []Rate{{
			Uuid:        rate_uuid,
			Country:     "Tycho crater",
			Data:        "1024",
			EffectiveAt: "2009-11-10T23:00:00",
			Network:     "Multi Tel",
			SimType:     3,
			X2g:         "",
			X3g:         "",
			Apn:         "",
			Imsi:        "",
			Lte:         "",
			SmsMo:       "",
			SmsMt:       "",
			Vpmn:        "",
			EndAt:       "",
			LteM:        "",
			X5g:         "",
		}}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), rate_uuid, rates[0].Country, rates[0].Network,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), rates[0].Data,
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), rates[0].EffectiveAt, sqlmock.AnyArg(), rates[0].SimType).
			WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(rate_uuid))

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
