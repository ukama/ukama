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
	"testing"
	"time"

	"github.com/ukama/ukama/systems/common/ukama"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	uuid "github.com/ukama/ukama/systems/common/uuid"
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
	log.Fatal("implement me")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB, nestedFuncs ...func(tx *gorm.DB) error) (err error) {
	log.Fatal("implement me")
	return nil
}

// Test utilities
type TestSetup struct {
	Mock   sqlmock.Sqlmock
	GormDB *gorm.DB
	Repo   *packageRepo
}

func setupTestDB(t *testing.T) *TestSetup {
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

	repo := NewPackageRepo(&UkamaDbMock{GormDb: gdb})

	return &TestSetup{
		Mock:   mock,
		GormDB: gdb,
		Repo:   repo,
	}
}

// Test data constants
const (
	TestPackageUUID = "51fbba62-c79f-11eb-b8bc-0242ac130003"
	TestPackageName = "Silver Plan"
	TestCountry     = "USA"
	TestProvider    = "ukama"
	TestCurrency    = "Dollar"
	TestAPN         = "uakam.tel"
)

// Test data values
var (
	TestDuration    uint64  = 30
	TestSmsVolume   uint64  = 1000
	TestDataVolume  uint64  = 5000000
	TestVoiceVolume uint64  = 500
	TestAmount      float64 = 100
	TestSmsMo       float64 = 0.001
	TestSmsMt       float64 = 0.001
	TestData        float64 = 0.010
	TestMarkup      float64 = 20
	TestDlbr        uint64  = 1024000
	TestUlbr        uint64  = 102400
)

func createTestPackage(packID uuid.UUID, name string) *Package {
	return &Package{
		Uuid:         packID,
		Name:         name,
		SimType:      ukama.SimTypeUkamaData,
		OwnerId:      uuid.NewV4(),
		Active:       true,
		Duration:     TestDuration,
		SmsVolume:    TestSmsVolume,
		DataVolume:   TestDataVolume,
		VoiceVolume:  TestVoiceVolume,
		Type:         ukama.PackageTypePostpaid,
		DataUnits:    ukama.DataUnitTypeMB,
		VoiceUnits:   ukama.CallUnitTypeSec,
		MessageUnits: ukama.MessageUnitTypeInt,
		Flatrate:     false,
		Currency:     TestCurrency,
		From:         time.Now(),
		To:           time.Now().Add(time.Hour * 24 * 30),
		Country:      TestCountry,
		Provider:     TestProvider,
	}
}

func createTestPackageRate(packID uuid.UUID) *PackageRate {
	return &PackageRate{
		PackageID: packID,
		Amount:    TestAmount,
		SmsMo:     TestSmsMo,
		SmsMt:     TestSmsMt,
		Data:      TestData,
	}
}

func createTestPackageDetails(packID uuid.UUID) *PackageDetails {
	return &PackageDetails{
		PackageID: packID,
		Dlbr:      TestDlbr,
		Ulbr:      TestUlbr,
		Apn:       TestAPN,
	}
}

func createTestPackageMarkup(packID uuid.UUID) *PackageMarkup {
	return &PackageMarkup{
		PackageID:  packID,
		BaseRateId: uuid.NewV4(),
		Markup:     TestMarkup,
	}
}

func createPackageRows(pack *Package) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"uuid", "owner_id", "name", "sim_type", "active", "duration",
		"sms_volume", "data_volume", "voice_volume", "type", "data_units",
		"voice_units", "message_units", "flat_rate", "currency", "from",
		"to", "country", "provider",
	}).AddRow(
		pack.Uuid, pack.OwnerId, pack.Name, pack.SimType, pack.Active,
		pack.Duration, pack.SmsVolume, pack.DataVolume, pack.VoiceVolume,
		pack.Type, pack.DataUnits, pack.VoiceUnits, pack.MessageUnits,
		pack.Flatrate, pack.Currency, pack.From, pack.To, pack.Country, pack.Provider,
	)
}

func createPackageRateRows(packID uuid.UUID, rate *PackageRate) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"package_id", "amount", "sms_mo", "sms_mt", "data"}).
		AddRow(packID, rate.Amount, rate.SmsMo, rate.SmsMt, rate.Data)
}

func createPackageDetailsRows(packID uuid.UUID, details *PackageDetails) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"package_id", "dlbr", "ulbr", "apn"}).
		AddRow(packID, details.Dlbr, details.Ulbr, details.Apn)
}

func createPackageMarkupRows(packID uuid.UUID, markup *PackageMarkup) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"package_id", "base_rate_id", "markup"}).
		AddRow(packID, markup.BaseRateId, markup.Markup)
}

func expectPackageQuery(mock sqlmock.Sqlmock, packID uuid.UUID, pack *Package) {
	rows := createPackageRows(pack)
	mock.ExpectQuery(`^SELECT.*packages.*`).
		WithArgs(packID, sqlmock.AnyArg()).
		WillReturnRows(rows)
}

func expectPackageRateQuery(mock sqlmock.Sqlmock, packID uuid.UUID, rate *PackageRate) {
	rows := createPackageRateRows(packID, rate)
	mock.ExpectQuery(`^SELECT.*package_rates.*`).
		WithArgs(packID).
		WillReturnRows(rows)
}

func expectPackageDetailsQuery(mock sqlmock.Sqlmock, packID uuid.UUID, details *PackageDetails) {
	rows := createPackageDetailsRows(packID, details)
	mock.ExpectQuery(`^SELECT.*package_details.*`).
		WithArgs(packID).
		WillReturnRows(rows)
}

func expectPackageMarkupQuery(mock sqlmock.Sqlmock, packID uuid.UUID, markup *PackageMarkup) {
	rows := createPackageMarkupRows(packID, markup)
	mock.ExpectQuery(`^SELECT.*package_markups.*`).
		WithArgs(packID).
		WillReturnRows(rows)
}

func Test_Package_Get(t *testing.T) {

	t.Run("PackageExistGet", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)
		packID, _ := uuid.FromString(TestPackageUUID)
		pack := createTestPackage(packID, TestPackageName)
		rate := createTestPackageRate(packID)

		expectPackageQuery(setup.Mock, packID, pack)
		expectPackageRateQuery(setup.Mock, packID, rate)

		// Act
		result, err := setup.Repo.Get(packID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		err = setup.Mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("PackageExistGetDetails", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)
		packID, _ := uuid.FromString(TestPackageUUID)
		pack := createTestPackage(packID, TestPackageName)
		rate := createTestPackageRate(packID)
		details := createTestPackageDetails(packID)
		markup := createTestPackageMarkup(packID)

		expectPackageQuery(setup.Mock, packID, pack)
		expectPackageDetailsQuery(setup.Mock, packID, details)
		expectPackageMarkupQuery(setup.Mock, packID, markup)
		expectPackageRateQuery(setup.Mock, packID, rate)

		// Act
		result, err := setup.Repo.GetDetails(packID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		err = setup.Mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_Package_GetAll(t *testing.T) {

	t.Run("PackageExist", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)
		packID, _ := uuid.FromString(TestPackageUUID)
		pack := createTestPackage(packID, TestPackageName)
		rate := createTestPackageRate(packID)

		rows := createPackageRows(pack)
		setup.Mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnRows(rows)

		expectPackageRateQuery(setup.Mock, packID, rate)

		// Act
		result, err := setup.Repo.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		err = setup.Mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("MultiplePackages", func(t *testing.T) {
		// Arrange
		setup := setupTestDB(t)
		packID1, _ := uuid.FromString("51fbba62-c79f-11eb-b8bc-0242ac130003")
		packID2, _ := uuid.FromString("51fbba62-c79f-11eb-b8bc-0242ac130004")
		ownerID1 := uuid.NewV4()
		ownerID2 := uuid.NewV4()

		// Create test packages with different data
		pack1 := createTestPackage(packID1, "Silver Plan")
		pack1.OwnerId = ownerID1
		pack2 := createTestPackage(packID2, "Gold Plan")
		pack2.OwnerId = ownerID2

		rate1 := createTestPackageRate(packID1)
		rate2 := createTestPackageRate(packID2)

		// Mock multiple packages
		rows := sqlmock.NewRows([]string{
			"uuid", "owner_id", "name", "sim_type", "active", "duration",
			"sms_volume", "data_volume", "voice_volume", "type", "data_units",
			"voice_units", "message_units", "flat_rate", "currency", "from",
			"to", "country", "provider",
		}).
			AddRow(packID1, ownerID1, pack1.Name, pack1.SimType, pack1.Active,
				pack1.Duration, pack1.SmsVolume, pack1.DataVolume, pack1.VoiceVolume,
				pack1.Type, pack1.DataUnits, pack1.VoiceUnits, pack1.MessageUnits,
				pack1.Flatrate, pack1.Currency, pack1.From, pack1.To, pack1.Country, pack1.Provider).
			AddRow(packID2, ownerID2, pack2.Name, pack2.SimType, pack2.Active,
				pack2.Duration, pack2.SmsVolume, pack2.DataVolume, pack2.VoiceVolume,
				pack2.Type, pack2.DataUnits, pack2.VoiceUnits, pack2.MessageUnits,
				pack2.Flatrate, pack2.Currency, pack2.From, pack2.To, pack2.Country, pack2.Provider)

		rrows := sqlmock.NewRows([]string{"package_id", "amount", "sms_mo", "sms_mt", "data"}).
			AddRow(packID1, rate1.Amount, rate1.SmsMo, rate1.SmsMt, rate1.Data).
			AddRow(packID2, rate2.Amount, rate2.SmsMo, rate2.SmsMt, rate2.Data)

		setup.Mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnRows(rows)

		setup.Mock.ExpectQuery(`^SELECT.*package_rates.*`).
			WithArgs(packID1, packID2).
			WillReturnRows(rrows)

		// Act
		packages, err := setup.Repo.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, packages, 2)
		assert.Equal(t, packID1, packages[0].Uuid)
		assert.Equal(t, "Silver Plan", packages[0].Name)
		assert.Equal(t, packID2, packages[1].Uuid)
		assert.Equal(t, "Gold Plan", packages[1].Name)

		err = setup.Mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAll_EmptyResult", func(t *testing.T) {
		setup := setupTestDB(t)

		// Return empty result
		rows := sqlmock.NewRows([]string{})
		setup.Mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnRows(rows)

		packages, err := setup.Repo.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, packages)
	})

	t.Run("GetAll_DatabaseError", func(t *testing.T) {
		setup := setupTestDB(t)

		// Mock database error
		setup.Mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnError(errors.New("connection timeout"))

		packages, err := setup.Repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, packages)
		assert.Contains(t, err.Error(), "connection timeout")

		err = setup.Mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAll_PackageRatePreloadError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID, _ := uuid.FromString("51fbba62-c79f-11eb-b8bc-0242ac130003")
		ownerID := uuid.NewV4()

		pack := createTestPackage(packID, "Silver Plan")
		pack.OwnerId = ownerID

		// Mock successful package query
		rows := createPackageRows(pack)
		setup.Mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnRows(rows)

		// Mock package rate query error
		setup.Mock.ExpectQuery(`^SELECT.*package_rates.*`).
			WithArgs(packID).
			WillReturnError(errors.New("package rate query failed"))

		packages, err := setup.Repo.GetAll()
		assert.Error(t, err)
		assert.Nil(t, packages)
		assert.Contains(t, err.Error(), "package rate query failed")

		err = setup.Mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetAll_PackageWithoutRate", func(t *testing.T) {
		setup := setupTestDB(t)
		packID, _ := uuid.FromString("51fbba62-c79f-11eb-b8bc-0242ac130003")
		ownerID := uuid.NewV4()

		pack := createTestPackage(packID, "No Rate Plan")
		pack.OwnerId = ownerID

		rows := createPackageRows(pack)
		emptyRateRows := sqlmock.NewRows([]string{"package_id", "amount", "sms_mo", "sms_mt", "data"})

		setup.Mock.ExpectQuery(`^SELECT.*packages.*`).
			WithArgs().
			WillReturnRows(rows)

		setup.Mock.ExpectQuery(`^SELECT.*package_rates.*`).
			WithArgs(packID).
			WillReturnRows(emptyRateRows)

		// Act
		packages, err := setup.Repo.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.Len(t, packages, 1)
		assert.Equal(t, packID, packages[0].Uuid)
		assert.Equal(t, "No Rate Plan", packages[0].Name)
		assert.Equal(t, uuid.Nil, packages[0].PackageRate.PackageID)

		err = setup.Mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_Package_Get_Error(t *testing.T) {
	t.Run("Get_NotFound", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnError(gorm.ErrRecordNotFound)

		package_, err := setup.Repo.Get(packID)
		assert.Error(t, err)
		assert.Nil(t, package_)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Get_DBError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnError(errors.New("database error"))

		package_, err := setup.Repo.Get(packID)
		assert.Error(t, err)
		assert.Nil(t, package_)
		assert.Contains(t, err.Error(), "database error")
	})
}

func Test_Package_GetDetails_Error(t *testing.T) {
	t.Run("GetDetails_NotFound", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnError(gorm.ErrRecordNotFound)

		package_, err := setup.Repo.GetDetails(packID)
		assert.Error(t, err)
		assert.Nil(t, package_)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("GetDetails_DBError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectQuery(`^SELECT.*packages.*`).WillReturnError(errors.New("database error"))

		package_, err := setup.Repo.GetDetails(packID)
		assert.Error(t, err)
		assert.Nil(t, package_)
		assert.Contains(t, err.Error(), "database error")
	})
}

func Test_Package_Add(t *testing.T) {
	t.Run("Add_Success", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		package_ := &Package{
			Uuid:    packID,
			Name:    "Test Package",
			Active:  true,
			Country: TestCountry,
		}

		packageRate := &PackageRate{
			PackageID: packID,
			Amount:    TestAmount,
		}

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectCommit()

		err := setup.Repo.Add(package_, packageRate)
		assert.NoError(t, err)
	})

	t.Run("Add_TransactionBeginError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		package_ := &Package{Uuid: packID, Name: "Test Package"}
		packageRate := &PackageRate{PackageID: packID, Amount: TestAmount}

		setup.Mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		err := setup.Repo.Add(package_, packageRate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")
	})

	t.Run("Add_PackageCreationError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		package_ := &Package{Uuid: packID, Name: "Test Package"}
		packageRate := &PackageRate{PackageID: packID, Amount: TestAmount}

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnError(errors.New("package creation error"))
		setup.Mock.ExpectRollback()

		err := setup.Repo.Add(package_, packageRate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "package creation error")
	})

	t.Run("Add_PackageRateCreationError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		package_ := &Package{Uuid: packID, Name: "Test Package"}
		packageRate := &PackageRate{PackageID: packID, Amount: TestAmount}

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnError(errors.New("package rate creation error"))
		setup.Mock.ExpectRollback()

		err := setup.Repo.Add(package_, packageRate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "package rate creation error")
	})

	t.Run("Add_CommitError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		package_ := &Package{Uuid: packID, Name: "Test Package"}
		packageRate := &PackageRate{PackageID: packID, Amount: TestAmount}

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectCommit().WillReturnError(errors.New("commit error"))

		err := setup.Repo.Add(package_, packageRate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "commit error")
	})

	t.Run("Add_WithFullPackageData", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()
		ownerID := uuid.NewV4()

		package_ := createTestPackage(packID, TestPackageName)
		package_.OwnerId = ownerID

		packageRate := createTestPackageRate(packID)

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectCommit()

		err := setup.Repo.Add(package_, packageRate)
		assert.NoError(t, err)
	})

	t.Run("Add_WithMinimalPackageData", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		package_ := &Package{
			Uuid:    packID,
			Name:    "Basic Plan",
			Country: TestCountry,
		}

		packageRate := &PackageRate{
			PackageID: packID,
		}

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectCommit()

		err := setup.Repo.Add(package_, packageRate)
		assert.NoError(t, err)
	})

	t.Run("Add_WithZeroValues", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		package_ := createTestPackage(packID, "Zero Plan")
		package_.Active = false
		package_.Duration = 0
		package_.SmsVolume = 0
		package_.DataVolume = 0

		packageRate := createTestPackageRate(packID)
		packageRate.Amount = 0
		packageRate.SmsMo = 0
		packageRate.SmsMt = 0
		packageRate.Data = 0

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectCommit()

		err := setup.Repo.Add(package_, packageRate)
		assert.NoError(t, err)
	})

	t.Run("Add_WithNegativeValues", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		package_ := createTestPackage(packID, "Negative Plan")

		packageRate := createTestPackageRate(packID)
		packageRate.Amount = -10.0
		packageRate.SmsMo = -0.01
		packageRate.SmsMt = -0.01
		packageRate.Data = -0.05

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^INSERT INTO "packages"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectQuery(`^INSERT INTO "package_rates"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		setup.Mock.ExpectCommit()

		err := setup.Repo.Add(package_, packageRate)
		assert.NoError(t, err)
	})
}

func Test_Package_Delete(t *testing.T) {
	t.Run("Delete_Success", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectExec(`^UPDATE "packages" SET "deleted_at"`).WillReturnResult(sqlmock.NewResult(1, 1))
		setup.Mock.ExpectCommit()

		err := setup.Repo.Delete(packID)
		assert.NoError(t, err)
	})

	t.Run("Delete_DBError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectExec(`^UPDATE "packages" SET "deleted_at"`).WillReturnError(errors.New("delete failed"))
		setup.Mock.ExpectRollback()

		err := setup.Repo.Delete(packID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}

func Test_Package_Update(t *testing.T) {
	t.Run("Update_Success", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^UPDATE "packages" SET.*WHERE.*uuid.*RETURNING`).
			WithArgs(sqlmock.AnyArg(), "Updated Name", packID).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "name"}).AddRow(packID, "Updated Name"))
		setup.Mock.ExpectCommit()

		pkg := &Package{Name: "Updated Name"}
		err := setup.Repo.Update(packID, pkg)
		assert.NoError(t, err)

		err = setup.Mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Update_TransactionBeginError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))

		pkg := &Package{Name: "Updated Name"}
		err := setup.Repo.Update(packID, pkg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction begin error")
	})

	t.Run("Update_NoRowsAffected", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^UPDATE "packages" SET.*WHERE.*uuid.*RETURNING`).
			WithArgs(sqlmock.AnyArg(), "Updated Name", packID).
			WillReturnRows(sqlmock.NewRows([]string{}))
		setup.Mock.ExpectRollback()

		pkg := &Package{Name: "Updated Name"}
		err := setup.Repo.Update(packID, pkg)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Update_DatabaseError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^UPDATE "packages" SET.*WHERE.*uuid.*RETURNING`).
			WithArgs(sqlmock.AnyArg(), "Updated Name", packID).
			WillReturnError(errors.New("database error"))
		setup.Mock.ExpectRollback()

		pkg := &Package{Name: "Updated Name"}
		err := setup.Repo.Update(packID, pkg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("Update_CommitError", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^UPDATE "packages" SET.*WHERE.*uuid.*RETURNING`).
			WithArgs(sqlmock.AnyArg(), "Updated Name", packID).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "name"}).AddRow(packID, "Updated Name"))
		setup.Mock.ExpectCommit().WillReturnError(errors.New("commit error"))

		pkg := &Package{Name: "Updated Name"}
		err := setup.Repo.Update(packID, pkg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "commit error")
	})

	t.Run("Update_WithMultipleFields", func(t *testing.T) {
		setup := setupTestDB(t)
		packID := uuid.NewV4()

		setup.Mock.ExpectBegin()
		setup.Mock.ExpectQuery(`^UPDATE "packages" SET.*WHERE.*uuid.*RETURNING`).
			WithArgs(sqlmock.AnyArg(), "Premium Plan", true, 60, packID).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "name", "active", "duration"}).AddRow(packID, "Premium Plan", true, 60))
		setup.Mock.ExpectCommit()

		pkg := &Package{
			Name:     "Premium Plan",
			Active:   true,
			Duration: 60,
		}
		err := setup.Repo.Update(packID, pkg)
		assert.NoError(t, err)

		err = setup.Mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
