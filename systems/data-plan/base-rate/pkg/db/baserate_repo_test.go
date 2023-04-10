package db

import (
	extsql "database/sql"
	"log"
	"regexp"
	"testing"
	"time"

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

func TestBaseRateRepo_dbTest(t *testing.T) {

	t.Run("BaseRateById", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		var db *extsql.DB
		var err error
		eat, err := time.Parse(time.RFC3339, "2023-10-12T07:20:50.52Z")
		assert.NoError(t, err)
		expectedRate := &BaseRate{
			Uuid:        ratID,
			Country:     "India",
			Network:     "Airtel",
			Vpmn:        "123",
			Imsi:        2,
			SmsMo:       0.05,
			SmsMt:       0.06,
			Data:        0.07,
			X2g:         false,
			X3g:         false,
			X5g:         true,
			Lte:         true,
			LteM:        true,
			Apn:         "apn123",
			EffectiveAt: time.Now(),
			EndAt:       eat,
			SimType:     SimTypeUkamaData,
		}
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"uuid", "country", "network", "vpmn", "imsi", "sms_mo", "sms_mt", "data", "x2g", "x3g", "x5g", "lte", "lte_m", "apn", "effective_at", "end_at", "sim_type", "created_at", "updated_at", "deleted_at"}).
			AddRow(expectedRate.Uuid, expectedRate.Country, expectedRate.Network, expectedRate.Vpmn, expectedRate.Imsi, expectedRate.SmsMo, expectedRate.SmsMt, expectedRate.Data, expectedRate.X2g, expectedRate.X3g, expectedRate.X5g, expectedRate.Lte, expectedRate.LteM, expectedRate.Apn, expectedRate.EffectiveAt, expectedRate.EndAt, expectedRate.SimType, expectedRate.CreatedAt, expectedRate.UpdatedAt, expectedRate.DeletedAt)

		mock.ExpectQuery(`^SELECT.*rate.*`).
			WithArgs(ratID.String()).
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
		rate, err := r.GetBaseRateById(ratID)
		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Equal(t, rate, expectedRate)
		assert.NotNil(t, rate)
	})

	t.Run("BaseRateByCountry", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		var db *extsql.DB
		var err error
		eat, err := time.Parse(time.RFC3339, "2023-10-12T07:20:50.52Z")
		assert.NoError(t, err)
		expectedRate := &BaseRate{
			Uuid:        ratID,
			Country:     "ABC",
			Network:     "XYZ",
			Vpmn:        "123",
			Imsi:        2,
			SmsMo:       0.05,
			SmsMt:       0.06,
			Data:        0.07,
			X2g:         false,
			X3g:         false,
			X5g:         true,
			Lte:         true,
			LteM:        true,
			Apn:         "apn123",
			EffectiveAt: time.Now(),
			EndAt:       eat,
			SimType:     SimTypeUkamaData,
		}
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"uuid", "country", "network", "vpmn", "imsi", "sms_mo", "sms_mt", "data", "x2g", "x3g", "x5g", "lte", "lte_m", "apn", "effective_at", "end_at", "sim_type", "created_at", "updated_at", "deleted_at"}).
			AddRow(expectedRate.Uuid, expectedRate.Country, expectedRate.Network, expectedRate.Vpmn, expectedRate.Imsi, expectedRate.SmsMo, expectedRate.SmsMt, expectedRate.Data, expectedRate.X2g, expectedRate.X3g, expectedRate.X5g, expectedRate.Lte, expectedRate.LteM, expectedRate.Apn, expectedRate.EffectiveAt, expectedRate.EndAt, expectedRate.SimType, expectedRate.CreatedAt, expectedRate.UpdatedAt, expectedRate.DeletedAt)

		mock.ExpectQuery(`^SELECT.*rate.*`).
			WithArgs(expectedRate.Country, expectedRate.Network, expectedRate.SimType, sqlmock.AnyArg()).
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
		rate, err := r.GetBaseRatesByCountry(expectedRate.Country, expectedRate.Network, expectedRate.SimType)
		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Equal(t, &rate[0], expectedRate)
		assert.NotNil(t, rate)
	})

	t.Run("BaseRateHistoryByCountry", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		var db *extsql.DB
		var err error
		sat, err := time.Parse(time.RFC3339, "2021-10-12T07:20:50.52Z")
		assert.NoError(t, err)

		eat, err := time.Parse(time.RFC3339, "2023-10-12T07:20:50.52Z")
		assert.NoError(t, err)

		expectedRate := []BaseRate{
			{
				Uuid:        ratID,
				Country:     "ABC",
				Network:     "XYZ",
				Vpmn:        "123",
				Imsi:        2,
				SmsMo:       0.05,
				SmsMt:       0.06,
				Data:        0.07,
				X2g:         false,
				X3g:         false,
				X5g:         true,
				Lte:         true,
				LteM:        true,
				Apn:         "apn123",
				EffectiveAt: sat,
				EndAt:       eat,
				SimType:     SimTypeUkamaData,
			},
			{
				Uuid:        uuid.NewV4(),
				Country:     "ABCDE",
				Network:     "XYZXX",
				Vpmn:        "123",
				Imsi:        2,
				SmsMo:       0.05,
				SmsMt:       0.06,
				Data:        0.07,
				X2g:         false,
				X3g:         false,
				X5g:         true,
				Lte:         true,
				LteM:        true,
				Apn:         "apn123",
				EffectiveAt: time.Now(),
				EndAt:       eat,
				SimType:     SimTypeUkamaData,
			},
		}

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"uuid", "country", "network", "vpmn", "imsi", "sms_mo", "sms_mt", "data", "x2g", "x3g", "x5g", "lte", "lte_m", "apn", "effective_at", "end_at", "sim_type", "created_at", "updated_at", "deleted_at"}).
			AddRow(expectedRate[0].Uuid, expectedRate[0].Country, expectedRate[0].Network, expectedRate[0].Vpmn, expectedRate[0].Imsi, expectedRate[0].SmsMo, expectedRate[0].SmsMt, expectedRate[0].Data, expectedRate[0].X2g, expectedRate[0].X3g, expectedRate[0].X5g, expectedRate[0].Lte, expectedRate[0].LteM, expectedRate[0].Apn, expectedRate[0].EffectiveAt, expectedRate[0].EndAt, expectedRate[0].SimType, expectedRate[0].CreatedAt, expectedRate[0].UpdatedAt, expectedRate[0].DeletedAt).
			AddRow(expectedRate[1].Uuid, expectedRate[1].Country, expectedRate[1].Network, expectedRate[1].Vpmn, expectedRate[1].Imsi, expectedRate[1].SmsMo, expectedRate[1].SmsMt, expectedRate[1].Data, expectedRate[1].X2g, expectedRate[1].X3g, expectedRate[1].X5g, expectedRate[1].Lte, expectedRate[1].LteM, expectedRate[1].Apn, expectedRate[1].EffectiveAt, expectedRate[1].EndAt, expectedRate[1].SimType, expectedRate[1].CreatedAt, expectedRate[1].UpdatedAt, expectedRate[1].DeletedAt)

		mock.ExpectQuery(`^SELECT.*rate.*`).
			WithArgs(expectedRate[0].Country, expectedRate[0].Network, expectedRate[0].SimType).
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
		rate, err := r.GetBaseRatesHistoryByCountry(expectedRate[0].Country, expectedRate[0].Network, expectedRate[0].SimType)
		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Equal(t, len(expectedRate), len(rate))
		assert.Equal(t, rate, expectedRate)
		assert.NotNil(t, rate)
	})

	t.Run("BaseRateForPeriod", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		var db *extsql.DB
		var err error
		from, err := time.Parse(time.RFC3339, "2022-10-12T07:20:50.52Z")
		assert.NoError(t, err)
		to, err := time.Parse(time.RFC3339, "2023-10-11T07:20:50.52Z")
		assert.NoError(t, err)
		eat, err := time.Parse(time.RFC3339, "2023-10-12T07:20:50.52Z")
		assert.NoError(t, err)
		expectedRate := &BaseRate{
			Uuid:        ratID,
			Country:     "ABC",
			Network:     "XYZ",
			Vpmn:        "123",
			Imsi:        2,
			SmsMo:       0.05,
			SmsMt:       0.06,
			Data:        0.07,
			X2g:         false,
			X3g:         false,
			X5g:         true,
			Lte:         true,
			LteM:        true,
			Apn:         "apn123",
			EffectiveAt: time.Now(),
			EndAt:       eat,
			SimType:     SimTypeUkamaData,
		}
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"uuid", "country", "network", "vpmn", "imsi", "sms_mo", "sms_mt", "data", "x2g", "x3g", "x5g", "lte", "lte_m", "apn", "effective_at", "end_at", "sim_type", "created_at", "updated_at", "deleted_at"}).
			AddRow(expectedRate.Uuid, expectedRate.Country, expectedRate.Network, expectedRate.Vpmn, expectedRate.Imsi, expectedRate.SmsMo, expectedRate.SmsMt, expectedRate.Data, expectedRate.X2g, expectedRate.X3g, expectedRate.X5g, expectedRate.Lte, expectedRate.LteM, expectedRate.Apn, expectedRate.EffectiveAt, expectedRate.EndAt, expectedRate.SimType, expectedRate.CreatedAt, expectedRate.UpdatedAt, expectedRate.DeletedAt)

		mock.ExpectQuery(`^SELECT.*rate.*`).
			WithArgs(expectedRate.Country, expectedRate.Network, expectedRate.SimType, from, to).
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
		rate, err := r.GetBaseRatesForPeriod(expectedRate.Country, expectedRate.Network, from, to, expectedRate.SimType)
		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Equal(t, &rate[0], expectedRate)
		assert.NotNil(t, rate)
	})

	t.Run("UploadBaseRates", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		var db *extsql.DB
		var err error
		sat, err := time.Parse(time.RFC3339, "2021-10-12T07:20:50.52Z")
		assert.NoError(t, err)

		eat, err := time.Parse(time.RFC3339, "2023-10-12T07:20:50.52Z")
		assert.NoError(t, err)

		expectedRate := BaseRate{

			Uuid:        ratID,
			Country:     "ABC",
			Network:     "XYZ",
			Vpmn:        "123",
			Imsi:        2,
			SmsMo:       0.05,
			SmsMt:       0.06,
			Data:        0.07,
			X2g:         false,
			X3g:         false,
			X5g:         true,
			Lte:         true,
			LteM:        true,
			Apn:         "apn123",
			EffectiveAt: sat,
			EndAt:       eat,
			SimType:     SimTypeUkamaData,
			Currency:    "Dollar",
		}

		upRates := []BaseRate{expectedRate}
		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).
			WithArgs(sqlmock.AnyArg(), expectedRate.Country, expectedRate.Network, expectedRate.SimType, expectedRate.EffectiveAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), expectedRate.Uuid, expectedRate.Country, expectedRate.Network, expectedRate.Vpmn, expectedRate.Imsi, expectedRate.SmsMo, expectedRate.SmsMt, expectedRate.Data, expectedRate.X2g, expectedRate.X3g, expectedRate.X5g, expectedRate.Lte, expectedRate.LteM, expectedRate.Apn, expectedRate.EffectiveAt, expectedRate.EndAt, expectedRate.SimType, expectedRate.Currency).
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

		assert.NoError(t, err)

		// Act
		err = r.UploadBaseRates(upRates)
		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
