package db

import (
	extsql "database/sql"
	"log"
	"testing"

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

func Test_Rate_Get(t *testing.T) {

	t.Run("RateExist", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		var db *extsql.DB
		var err error
		expectedRate := &Rate{
			Uuid:        ratID,
			Country:     "India",
			Network:     "Airtel",
			Vpmn:        "123",
			Imsi:        "456",
			SmsMo:       "0.05",
			SmsMt:       "0.06",
			Data:        "0.07",
			X2g:         "0.08",
			X3g:         "0.09",
			X5g:         "0.1",
			Lte:         "0.11",
			LteM:        "0.12",
			Apn:         "apn123",
			EffectiveAt: "2022-01-01",
			EndAt:       "2022-12-31",
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
		rate, err := r.GetBaseRate(ratID)
		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.Equal(t, rate, expectedRate)
		assert.NotNil(t, rate)
	})

}
