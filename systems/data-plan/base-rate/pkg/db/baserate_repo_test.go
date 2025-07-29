/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"errors"
	"log"
	"regexp"
	"testing"
	"time"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/tj/assert"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Test data constants
const (
	testCountry  = "ABC"
	testProvider = "XYZ"
	testVpmn     = "123"
	testApn      = "apn123"
	testCurrency = "Dollar"
	testImsi     = 2
	testSmsMo    = 0.05
	testSmsMt    = 0.06
	testData     = 0.07
)

// Test time constants
var (
	testEndAt    = time.Date(2023, 10, 12, 7, 20, 50, 520000000, time.UTC)
	testStartAt  = time.Date(2021, 10, 12, 7, 20, 50, 520000000, time.UTC)
	testFromDate = time.Date(2022, 10, 12, 7, 20, 50, 520000000, time.UTC)
	testToDate   = time.Date(2023, 10, 11, 7, 20, 50, 520000000, time.UTC)
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

// Test helper functions
func createTestBaseRate(id uuid.UUID, country, provider string, effectiveAt, endAt time.Time) *BaseRate {
	return &BaseRate{
		Uuid:        id,
		Country:     country,
		Provider:    provider,
		Vpmn:        testVpmn,
		Imsi:        testImsi,
		SmsMo:       testSmsMo,
		SmsMt:       testSmsMt,
		Data:        testData,
		X2g:         false,
		X3g:         false,
		X5g:         true,
		Lte:         true,
		LteM:        true,
		Apn:         testApn,
		EffectiveAt: effectiveAt,
		EndAt:       endAt,
		SimType:     ukama.SimTypeUkamaData,
		Currency:    testCurrency,
	}
}

func createTestBaseRateWithoutCurrency(id uuid.UUID, country, provider string, effectiveAt, endAt time.Time) *BaseRate {
	rate := createTestBaseRate(id, country, provider, effectiveAt, endAt)
	rate.Currency = ""
	return rate
}

func setupTestDB(t *testing.T) (*baseRateRepo, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := NewBaseRateRepo(&UkamaDbMock{GormDb: gdb})
	cleanup := func() { db.Close() }

	return repo, mock, cleanup
}

func createMockRows(rate *BaseRate) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"uuid", "country", "provider", "vpmn", "imsi", "sms_mo", "sms_mt", "data",
		"x2g", "x3g", "x5g", "lte", "lte_m", "apn", "effective_at", "end_at",
		"sim_type", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		rate.Uuid, rate.Country, rate.Provider, rate.Vpmn, rate.Imsi,
		rate.SmsMo, rate.SmsMt, rate.Data, rate.X2g, rate.X3g, rate.X5g,
		rate.Lte, rate.LteM, rate.Apn, rate.EffectiveAt, rate.EndAt,
		rate.SimType, rate.CreatedAt, rate.UpdatedAt, rate.DeletedAt,
	)
}

func createMockRowsForMultipleRates(rates []BaseRate) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"uuid", "country", "provider", "vpmn", "imsi", "sms_mo", "sms_mt", "data",
		"x2g", "x3g", "x5g", "lte", "lte_m", "apn", "effective_at", "end_at",
		"sim_type", "created_at", "updated_at", "deleted_at",
	})

	for _, rate := range rates {
		rows.AddRow(
			rate.Uuid, rate.Country, rate.Provider, rate.Vpmn, rate.Imsi,
			rate.SmsMo, rate.SmsMt, rate.Data, rate.X2g, rate.X3g, rate.X5g,
			rate.Lte, rate.LteM, rate.Apn, rate.EffectiveAt, rate.EndAt,
			rate.SimType, rate.CreatedAt, rate.UpdatedAt, rate.DeletedAt,
		)
	}

	return rows
}

func TestBaseRateRepo_dbTest(t *testing.T) {

	t.Run("BaseRateById", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		expectedRate := createTestBaseRateWithoutCurrency(ratID, "India", "Airtel", time.Now(), testEndAt)

		repo, mock, cleanup := setupTestDB(t)
		defer cleanup()

		rows := createMockRows(expectedRate)
		mock.ExpectQuery(`^SELECT.*rate.*`).
			WithArgs(ratID.String(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		rate, err := repo.GetBaseRateById(ratID)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, rate, expectedRate)
		assert.NotNil(t, rate)
	})

	t.Run("BaseRateByCountry", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		expectedRate := createTestBaseRateWithoutCurrency(ratID, testCountry, testProvider, time.Now(), testEndAt)

		repo, mock, cleanup := setupTestDB(t)
		defer cleanup()

		rows := createMockRows(expectedRate)
		mock.ExpectQuery(`^SELECT.*rate.*`).
			WithArgs(expectedRate.Country, expectedRate.Provider, expectedRate.SimType, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		// Act
		rate, err := repo.GetBaseRatesByCountry(expectedRate.Country, expectedRate.Provider, expectedRate.SimType)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, &rate[0], expectedRate)
		assert.NotNil(t, rate)
	})

	t.Run("BaseRateHistoryByCountry", func(t *testing.T) {
		// Arrange
		ratID1 := uuid.NewV4()
		ratID2 := uuid.NewV4()

		expectedRates := []BaseRate{
			*createTestBaseRateWithoutCurrency(ratID1, testCountry, testProvider, testStartAt, testEndAt),
			*createTestBaseRateWithoutCurrency(ratID2, "ABCDE", "XYZXX", time.Now(), testEndAt),
		}

		repo, mock, cleanup := setupTestDB(t)
		defer cleanup()

		rows := createMockRowsForMultipleRates(expectedRates)
		mock.ExpectQuery(`^SELECT.*rate.*`).
			WithArgs(expectedRates[0].Country, expectedRates[0].Provider, expectedRates[0].SimType).
			WillReturnRows(rows)

		// Act
		rate, err := repo.GetBaseRatesHistoryByCountry(expectedRates[0].Country, expectedRates[0].Provider, expectedRates[0].SimType)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, len(expectedRates), len(rate))
		assert.Equal(t, rate, expectedRates)
		assert.NotNil(t, rate)
	})

	t.Run("BaseRateForPeriod", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		expectedRate := createTestBaseRateWithoutCurrency(ratID, testCountry, testProvider, time.Now(), testEndAt)

		repo, mock, cleanup := setupTestDB(t)
		defer cleanup()

		rows := createMockRows(expectedRate)
		mock.ExpectQuery(`^SELECT.*rate.*`).
			WithArgs(expectedRate.Country, expectedRate.Provider, expectedRate.SimType, testFromDate, testToDate).
			WillReturnRows(rows)

		// Act
		rate, err := repo.GetBaseRatesForPeriod(expectedRate.Country, expectedRate.Provider, testFromDate, testToDate, expectedRate.SimType)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, &rate[0], expectedRate)
		assert.NotNil(t, rate)
	})

	t.Run("UploadBaseRates", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		expectedRate := createTestBaseRate(ratID, testCountry, testProvider, testStartAt, testEndAt)
		upRates := []BaseRate{*expectedRate}

		repo, mock, cleanup := setupTestDB(t)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).
			WithArgs(sqlmock.AnyArg(), expectedRate.Country, expectedRate.Provider, expectedRate.SimType, expectedRate.EffectiveAt).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), expectedRate.Uuid, expectedRate.Country, expectedRate.Provider, expectedRate.Vpmn, expectedRate.Imsi, expectedRate.SmsMo, expectedRate.SmsMt, expectedRate.Data, expectedRate.X2g, expectedRate.X3g, expectedRate.X5g, expectedRate.Lte, expectedRate.LteM, expectedRate.Apn, expectedRate.EffectiveAt, expectedRate.EndAt, expectedRate.SimType, expectedRate.Currency).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Act
		err := repo.UploadBaseRates(upRates)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UploadBaseRates create error", func(t *testing.T) {
		// Arrange
		ratID := uuid.NewV4()
		expectedRate := createTestBaseRate(ratID, testCountry, testProvider, testStartAt, testEndAt)
		upRates := []BaseRate{*expectedRate}

		repo, mock, cleanup := setupTestDB(t)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).
			WithArgs(sqlmock.AnyArg(), expectedRate.Country, expectedRate.Provider, expectedRate.SimType, expectedRate.EffectiveAt).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), expectedRate.Uuid, expectedRate.Country, expectedRate.Provider, expectedRate.Vpmn, expectedRate.Imsi, expectedRate.SmsMo, expectedRate.SmsMt, expectedRate.Data, expectedRate.X2g, expectedRate.X3g, expectedRate.X5g, expectedRate.Lte, expectedRate.LteM, expectedRate.Apn, expectedRate.EffectiveAt, expectedRate.EndAt, expectedRate.SimType, expectedRate.Currency).
			WillReturnError(errors.New("create error"))
		mock.ExpectRollback()

		// Act
		err := repo.UploadBaseRates(upRates)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBaseRateRepo_ErrorCases(t *testing.T) {
	newRepo := func() (*baseRateRepo, sqlmock.Sqlmock, func()) {
		return setupTestDB(t)
	}

	t.Run("GetBaseRateById returns error", func(t *testing.T) {
		repo, mock, cleanup := newRepo()
		defer cleanup()
		mock.ExpectQuery("SELECT.*rate.*").WillReturnError(errors.New("db error"))
		_, err := repo.GetBaseRateById(uuid.NewV4())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("GetBaseRatesHistoryByCountry returns error", func(t *testing.T) {
		repo, mock, cleanup := newRepo()
		defer cleanup()
		mock.ExpectQuery("SELECT.*rate.*").WillReturnError(errors.New("db error"))
		_, err := repo.GetBaseRatesHistoryByCountry("c", "p", ukama.SimTypeUkamaData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("GetBaseRatesByCountry returns error", func(t *testing.T) {
		repo, mock, cleanup := newRepo()
		defer cleanup()
		mock.ExpectQuery("SELECT.*rate.*").WillReturnError(errors.New("db error"))
		_, err := repo.GetBaseRatesByCountry("c", "p", ukama.SimTypeUkamaData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("GetBaseRatesForPeriod returns error", func(t *testing.T) {
		repo, mock, cleanup := newRepo()
		defer cleanup()
		mock.ExpectQuery("SELECT.*rate.*").WillReturnError(errors.New("db error"))
		_, err := repo.GetBaseRatesForPeriod("c", "p", time.Now(), time.Now(), ukama.SimTypeUkamaData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("GetBaseRatesForPackage returns error", func(t *testing.T) {
		repo, mock, cleanup := newRepo()
		defer cleanup()
		mock.ExpectQuery("SELECT.*rate.*").WillReturnError(errors.New("db error"))
		_, err := repo.GetBaseRatesForPackage("c", "p", time.Now(), time.Now(), ukama.SimTypeUkamaData)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("UploadBaseRates delete error (not NotFound)", func(t *testing.T) {
		repo, mock, cleanup := newRepo()
		defer cleanup()
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE")).WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()
		rate := BaseRate{Country: "c", Provider: "p", SimType: ukama.SimTypeUkamaData, EffectiveAt: time.Now()}
		err := repo.UploadBaseRates([]BaseRate{rate})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete error")
	})

	t.Run("GetBaseRatesByCountry empty result", func(t *testing.T) {
		repo, mock, cleanup := newRepo()
		defer cleanup()
		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT.*rate.*").WillReturnRows(rows)
		rates, err := repo.GetBaseRatesByCountry("c", "p", ukama.SimTypeUkamaData)
		assert.NoError(t, err)
		assert.Empty(t, rates)
	})

	t.Run("GetBaseRatesForPeriod empty result", func(t *testing.T) {
		repo, mock, cleanup := newRepo()
		defer cleanup()
		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT.*rate.*").WillReturnRows(rows)
		rates, err := repo.GetBaseRatesForPeriod("c", "p", time.Now(), time.Now(), ukama.SimTypeUkamaData)
		assert.NoError(t, err)
		assert.Empty(t, rates)
	})

	t.Run("GetBaseRatesForPackage empty result", func(t *testing.T) {
		repo, mock, cleanup := newRepo()
		defer cleanup()
		rows := sqlmock.NewRows([]string{})
		mock.ExpectQuery("SELECT.*rate.*").WillReturnRows(rows)
		rates, err := repo.GetBaseRatesForPackage("c", "p", time.Now(), time.Now(), ukama.SimTypeUkamaData)
		assert.NoError(t, err)
		assert.Empty(t, rates)
	})

	t.Run("UploadBaseRates transaction begin error", func(t *testing.T) {
		repo, mock, cleanup := newRepo()
		defer cleanup()
		mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))
		rate := BaseRate{Country: "c", Provider: "p", SimType: ukama.SimTypeUkamaData, EffectiveAt: time.Now()}
		err := repo.UploadBaseRates([]BaseRate{rate})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")
	})
}
