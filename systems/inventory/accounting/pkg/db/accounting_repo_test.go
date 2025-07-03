/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
	account_db "github.com/ukama/ukama/systems/inventory/accounting/pkg/db"
)

// Test constants for common data
const (
	testItem          = "tower node"
	testInventory     = "5"
	testEffectiveDate = "12/12/2024"
	testOpexFee       = "1.00"
	testVat           = "1.00"
	testDescription   = "Some description"
	testTowerDesc     = "Tower node descp"
)

// Test data structures
type testData struct {
	accountID uuid.UUID
	userID    uuid.UUID
	account   *account_db.Accounting
}

// Helper functions
func createTestData() *testData {
	accountID := uuid.NewV4()
	userID := uuid.NewV4()

	return &testData{
		accountID: accountID,
		userID:    userID,
		account: &account_db.Accounting{
			Id:            accountID,
			Item:          testItem,
			UserId:        userID,
			Inventory:     testInventory,
			EffectiveDate: testEffectiveDate,
			OpexFee:       testOpexFee,
			Vat:           testVat,
			Description:   testDescription,
		},
	}
}

func createTowerNodeTestData() *testData {
	accountID := uuid.NewV4()
	userID := uuid.NewV4()

	return &testData{
		accountID: accountID,
		userID:    userID,
		account: &account_db.Accounting{
			Id:            accountID,
			Item:          testItem,
			UserId:        userID,
			Inventory:     testInventory,
			EffectiveDate: testEffectiveDate,
			OpexFee:       testOpexFee,
			Vat:           testVat,
			Description:   testTowerDesc,
		},
	}
}

func setupMockDB() (sqlmock.Sqlmock, *gorm.DB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return mock, gdb, nil
}

func createMockRows(account *account_db.Accounting) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "item", "user_id", "inventory", "effective_date", "opex_fee", "vat", "description",
	}).AddRow(
		account.Id, account.Item, account.UserId, account.Inventory,
		account.EffectiveDate, account.OpexFee, account.Vat, account.Description,
	)
}

func createAccountingRepo(gdb *gorm.DB) account_db.AccountingRepo {
	return account_db.NewAccountingRepo(&UkamaDbMock{
		GormDb: gdb,
	})
}

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

func Test_AccountRepo_Get(t *testing.T) {
	t.Run("AccountExist", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb, err := setupMockDB()
		assert.NoError(t, err)

		rows := createMockRows(testData.account)
		mock.ExpectQuery(`^SELECT.*accounting.*`).
			WithArgs(testData.accountID, 1).
			WillReturnRows(rows)

		r := createAccountingRepo(gdb)

		// Act
		comp, err := r.Get(testData.accountID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, comp)
		assert.Equal(t, comp.Id, testData.accountID)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_AccountRepo_GetByUser(t *testing.T) {
	t.Run("AccountExist", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb, err := setupMockDB()
		assert.NoError(t, err)

		rows := createMockRows(testData.account)
		mock.ExpectQuery(`^SELECT.*accounting.*`).
			WithArgs(testData.userID.String()).
			WillReturnRows(rows)

		r := createAccountingRepo(gdb)

		// Act
		comps, err := r.GetByUser(testData.userID.String())

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, comps)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Arrange
		testData := createTestData()
		mock, gdb, err := setupMockDB()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*accounting.*`).
			WithArgs(testData.userID.String()).
			WillReturnError(fmt.Errorf("database connection error"))

		r := createAccountingRepo(gdb)

		// Act
		comps, err := r.GetByUser(testData.userID.String())

		// Assert
		assert.Error(t, err)
		assert.Nil(t, comps)
		assert.Contains(t, err.Error(), "database connection error")
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_accounttRepo_Add(t *testing.T) {
	t.Run("AddComponent", func(t *testing.T) {
		// Arrange
		testData := createTowerNodeTestData()
		accounts := []*account_db.Accounting{testData.account}

		mock, gdb, err := setupMockDB()
		assert.NoError(t, err)

		mock.ExpectBegin()
		for _, account := range accounts {
			mock.ExpectExec(`INSERT INTO "accountings" \("id","item","user_id","inventory","effective_date","opex_fee","vat","description"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\) ON CONFLICT \("id"\) DO NOTHING`).
				WithArgs(account.Id, account.Item, account.UserId, account.Inventory, account.EffectiveDate, account.OpexFee, account.Vat, account.Description).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()

		rows := createMockRows(accounts[0])
		mock.ExpectQuery(`^SELECT.*accountings.*`).
			WithArgs(accounts[0].UserId.String()).WillReturnRows(rows)

		r := createAccountingRepo(gdb)

		// Act
		err = r.Add(accounts)
		assert.NoError(t, err)

		res, err := r.GetByUser(testData.userID.String())

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_AccountRepo_Delete(t *testing.T) {
	t.Run("DeleteAccount", func(t *testing.T) {
		// Arrange
		testData := createTowerNodeTestData()
		accounts := []*account_db.Accounting{testData.account}

		mock, gdb, err := setupMockDB()
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "accountings" ("id","item","user_id","inventory","effective_date","opex_fee","vat","description") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT ("id") DO NOTHING`)).
			WithArgs(accounts[0].Id, accounts[0].Item, accounts[0].UserId, accounts[0].Inventory, accounts[0].EffectiveDate, accounts[0].OpexFee, accounts[0].Vat, accounts[0].Description).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM accountings`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery(`^SELECT.*accounting.*`).
			WithArgs(accounts[0].Id, 1).WillReturnRows(sqlmock.NewRows([]string{}))

		r := createAccountingRepo(gdb)

		// Act
		err = r.Add(accounts)
		assert.NoError(t, err)

		err = r.Delete()
		assert.NoError(t, err)

		res, err := r.Get(testData.accountID)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, res)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
