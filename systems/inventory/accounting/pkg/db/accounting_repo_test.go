/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	extsql "database/sql"
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

		var db *extsql.DB
		var accountID = uuid.NewV4()
		var uID = uuid.NewV4()

		db, mock, err := sqlmock.New() // mock sql.DB
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "item", "user_id", "inventory", "effective_date", "opex_fee", "vat", "description"}).
			AddRow(accountID, "tower node", uID, "5", "12/12/2024", "1.00", "1.00", "Some description")

		mock.ExpectQuery(`^SELECT.*accounting.*`).
			WithArgs(accountID, 1).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := account_db.NewAccountingRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		comp, err := r.Get(accountID)
		assert.NoError(t, err)
		assert.NotNil(t, comp)
		assert.Equal(t, comp.Id, accountID)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_AccountRepo_GetByUser(t *testing.T) {
	t.Run("AccountExist", func(t *testing.T) {

		var db *extsql.DB

		var accountID = uuid.NewV4()
		var uID = uuid.NewV4()

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "item", "user_id", "inventory", "effective_date", "opex_fee", "vat", "description"}).
			AddRow(accountID, "tower node", uID, "5", "12/12/2024", "1.00", "1.00", "Some description")

		mock.ExpectQuery(`^SELECT.*accounting.*`).
			WithArgs(uID.String()).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := account_db.NewAccountingRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		comps, err := r.GetByUser(uID.String())

		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, comps)
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		var db *extsql.DB

		var uID = uuid.NewV4()

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*accounting.*`).
			WithArgs(uID.String()).
			WillReturnError(fmt.Errorf("database connection error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := account_db.NewAccountingRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		comps, err := r.GetByUser(uID.String())

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
		var db *extsql.DB

		uId := uuid.NewV4()
		accounts := []*account_db.Accounting{
			{
				Id:            uuid.NewV4(),
				Inventory:     "5",
				UserId:        uId,
				Description:   "Tower node descp",
				Item:          "tower node",
				EffectiveDate: "12/12/2024",
				OpexFee:       "1.00",
				Vat:           "1.00",
			},
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		for _, account := range accounts {
			mock.ExpectExec(`INSERT INTO "accountings" \("id","item","user_id","inventory","effective_date","opex_fee","vat","description"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\) ON CONFLICT \("id"\) DO NOTHING`).
				WithArgs(account.Id, account.Item, account.UserId, account.Inventory, account.EffectiveDate, account.OpexFee, account.Vat, account.Description).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()

		mock.ExpectQuery(`^SELECT.*accountings.*`).
			WithArgs(accounts[0].UserId.String()).WillReturnRows(sqlmock.NewRows([]string{
			"id", "item", "user_id", "inventory", "effective_date", "opex_fee", "vat", "description",
		}).AddRow(accounts[0].Id, accounts[0].Item, accounts[0].UserId, accounts[0].Inventory, accounts[0].EffectiveDate, accounts[0].OpexFee, accounts[0].Vat, accounts[0].Description))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := account_db.NewAccountingRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)
		err = r.Add(accounts)
		assert.NoError(t, err)

		res, err := r.GetByUser(uId.String())
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_AccountRepo_Delete(t *testing.T) {
	t.Run("DeleteAccount", func(t *testing.T) {
		var db *extsql.DB

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		uId := uuid.NewV4()
		aId := uuid.NewV4()

		accounts := []*account_db.Accounting{
			{
				Id:            aId,
				Inventory:     "5",
				UserId:        uId,
				Description:   "Tower node descp",
				Item:          "tower node",
				EffectiveDate: "12/12/2024",
				OpexFee:       "1.00",
				Vat:           "1.00",
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "accountings" ("id","item","user_id","inventory","effective_date","opex_fee","vat","description") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT ("id") DO NOTHING`)).
			WithArgs(accounts[0].Id, accounts[0].Item, accounts[0].UserId, accounts[0].Inventory, accounts[0].EffectiveDate, accounts[0].OpexFee, accounts[0].Vat, accounts[0].Description).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM accountings`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery(`^SELECT.*accounting.*`).
			WithArgs(accounts[0].Id, 1).WillReturnRows(sqlmock.NewRows([]string{}))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := account_db.NewAccountingRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.Add(accounts)
		assert.NoError(t, err)

		err = r.Delete()
		assert.NoError(t, err)

		res, err := r.Get(aId)
		assert.Error(t, err)
		assert.Empty(t, res)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
