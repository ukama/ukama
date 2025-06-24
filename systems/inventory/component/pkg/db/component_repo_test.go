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

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	component_db "github.com/ukama/ukama/systems/inventory/component/pkg/db"
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

func Test_ComponentRepo_Get(t *testing.T) {
	t.Run("ComponentExist", func(t *testing.T) {

		var db *extsql.DB
		var componentID = uuid.NewV4()
		var uID = uuid.NewV4()

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "inventory", "user_id", "category", "type", "description", "datasheet_url", "image_url", "part_number", "manufacturer", "managed", "warranty", "specification"}).
			AddRow(componentID, "5", uID.String(), int64(1), "tower node", "Tower node descp", "http://datasheepurl", "http://imageurl", "123", "ukama", "ukama", 1, "metainfo")

		mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(componentID.String(), 1).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := component_db.NewComponentRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		comp, err := r.Get(componentID)

		assert.NoError(t, err)
		assert.NotNil(t, comp)
		assert.Equal(t, comp.Id, componentID)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_ComponentRepo_GetByUser(t *testing.T) {
	t.Run("ComponentExist", func(t *testing.T) {

		var db *extsql.DB

		var componentID = uuid.NewV4()
		var category = int32(1)
		var uID = uuid.NewV4()

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "inventory", "user_id", "category"}).
			AddRow(componentID, "5", uID.String(), category)

		mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(uID.String(), category).
			WillReturnRows(rows)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := component_db.NewComponentRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)

		comps, err := r.GetByUser(uID.String(), category)

		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		assert.NotNil(t, comps)
		assert.NoError(t, err)
	})

	t.Run("NoComponentsFound", func(t *testing.T) {
		var db *extsql.DB

		var category = int32(1)
		var uID = uuid.NewV4()

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(uID.String(), category).
			WillReturnRows(sqlmock.NewRows([]string{"id", "inventory", "user_id", "category", "type", "description", "datasheet_url", "images_url", "part_number", "manufacturer", "managed", "warranty", "specification"}))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := component_db.NewComponentRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		comps, err := r.GetByUser(uID.String(), category)

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, comps)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		var db *extsql.DB

		var category = int32(1)
		var uID = uuid.NewV4()

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(uID.String(), category).
			WillReturnError(fmt.Errorf("database connection error"))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := component_db.NewComponentRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		comps, err := r.GetByUser(uID.String(), category)

		assert.Error(t, err)
		assert.Nil(t, comps)
		assert.Contains(t, err.Error(), "database connection error")
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_ComponentRepo_Add(t *testing.T) {
	t.Run("AddComponent", func(t *testing.T) {
		var db *extsql.DB

		uId := uuid.NewV4()
		components := []*component_db.Component{
			{
				Id:            uuid.NewV4(),
				Inventory:     "5",
				UserId:        uId,
				Category:      ukama.ACCESS,
				Type:          "tower node",
				Description:   "Tower node descp",
				DatasheetURL:  "http://datasheepurl",
				ImagesURL:     "http://imageurl",
				PartNumber:    "123",
				Manufacturer:  "ukama",
				Managed:       "ukama",
				Warranty:      1,
				Specification: "metainfo",
			},
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		mock.ExpectBegin()
		for _, component := range components {
			mock.ExpectExec(`INSERT INTO "components" \("id","inventory","user_id","category","type","description","datasheet_url","images_url","part_number","manufacturer","managed","warranty","specification"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11,\$12,\$13\) ON CONFLICT \("id"\) DO NOTHING`).
				WithArgs(component.Id, component.Inventory, component.UserId, component.Category, component.Type, component.Description, component.DatasheetURL, component.ImagesURL, component.PartNumber, component.Manufacturer, component.Managed, component.Warranty, component.Specification).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		mock.ExpectCommit()

		mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(components[0].UserId.String(), int32(1)).WillReturnRows(sqlmock.NewRows([]string{
			"id", "inventory", "user_id", "category", "type", "description", "datasheet_url", "images_url", "part_number", "manufacturer", "managed", "warranty", "specification",
		}).AddRow(
			components[0].Id, components[0].Inventory, components[0].UserId, components[0].Category, components[0].Type, components[0].Description, components[0].DatasheetURL, components[0].ImagesURL, components[0].PartNumber, components[0].Manufacturer, components[0].Managed, components[0].Warranty, components[0].Specification,
		))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := component_db.NewComponentRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		assert.NoError(t, err)
		err = r.Add(components)
		assert.NoError(t, err)

		res, err := r.GetByUser(uId.String(), int32(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_ComponentRepo_Delete(t *testing.T) {
	t.Run("DeleteComponent", func(t *testing.T) {
		var db *extsql.DB

		cId := uuid.NewV4()
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		components := []*component_db.Component{
			{
				Id:            cId,
				Inventory:     "5",
				UserId:        uuid.NewV4(),
				Category:      ukama.ACCESS,
				Type:          "tower node",
				Description:   "Tower node descp",
				DatasheetURL:  "http://datasheepurl",
				ImagesURL:     "http://imageurl",
				PartNumber:    "123",
				Manufacturer:  "ukama",
				Managed:       "ukama",
				Warranty:      1,
				Specification: "metainfo",
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "components" ("id","inventory","user_id","category","type","description","datasheet_url","images_url","part_number","manufacturer","managed","warranty","specification") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) ON CONFLICT ("id") DO NOTHING`)).
			WithArgs(components[0].Id, components[0].Inventory, components[0].UserId, components[0].Category, components[0].Type, components[0].Description, components[0].DatasheetURL, components[0].ImagesURL, components[0].PartNumber, components[0].Manufacturer, components[0].Managed, components[0].Warranty, components[0].Specification).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM components`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(components[0].Id, 1).WillReturnRows(sqlmock.NewRows([]string{}))

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 db,
			PreferSimpleProtocol: true,
		})

		gdb, err := gorm.Open(dialector, &gorm.Config{})
		assert.NoError(t, err)

		r := component_db.NewComponentRepo(&UkamaDbMock{
			GormDb: gdb,
		})

		err = r.Add(components)
		assert.NoError(t, err)

		err = r.Delete()
		assert.NoError(t, err)

		res, err := r.Get(cId)
		assert.Error(t, err)
		assert.Empty(t, res)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
