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

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	component_db "github.com/ukama/ukama/systems/inventory/component/pkg/db"
)

// Test constants
const (
	testInventory     = "5"
	testType          = "tower node"
	testDescription   = "Tower node descp"
	testDatasheetURL  = "http://datasheepurl"
	testImageURL      = "http://imageurl"
	testPartNumber    = "123"
	testManufacturer  = "ukama"
	testManaged       = "ukama"
	testWarranty      = 1
	testSpecification = "metainfo"
	testCategory      = ukama.ACCESS
)

// Test data structures
type testComponentData struct {
	ID            uuid.UUID
	Inventory     string
	UserID        uuid.UUID
	Category      ukama.ComponentCategory
	Type          string
	Description   string
	DatasheetURL  string
	ImagesURL     string
	PartNumber    string
	Manufacturer  string
	Managed       string
	Warranty      uint32
	Specification string
}

type testDBSetup struct {
	mock   sqlmock.Sqlmock
	gormDB *gorm.DB
	repo   component_db.ComponentRepo
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

func createTestComponentData() *testComponentData {
	return &testComponentData{
		ID:            uuid.NewV4(),
		Inventory:     testInventory,
		UserID:        uuid.NewV4(),
		Category:      testCategory,
		Type:          testType,
		Description:   testDescription,
		DatasheetURL:  testDatasheetURL,
		ImagesURL:     testImageURL,
		PartNumber:    testPartNumber,
		Manufacturer:  testManufacturer,
		Managed:       testManaged,
		Warranty:      testWarranty,
		Specification: testSpecification,
	}
}

func convertToComponent(data *testComponentData) *component_db.Component {
	return &component_db.Component{
		Id:            data.ID,
		Inventory:     data.Inventory,
		UserId:        data.UserID,
		Category:      data.Category,
		Type:          data.Type,
		Description:   data.Description,
		DatasheetURL:  data.DatasheetURL,
		ImagesURL:     data.ImagesURL,
		PartNumber:    data.PartNumber,
		Manufacturer:  data.Manufacturer,
		Managed:       data.Managed,
		Warranty:      data.Warranty,
		Specification: data.Specification,
	}
}

func setupTestDB(t *testing.T) *testDBSetup {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := component_db.NewComponentRepo(&UkamaDbMock{
		GormDb: gormDB,
	})

	return &testDBSetup{
		mock:   mock,
		gormDB: gormDB,
		repo:   repo,
	}
}

func getComponentColumns() []string {
	return []string{
		"id", "inventory", "user_id", "category", "type", "description",
		"datasheet_url", "images_url", "part_number", "manufacturer",
		"managed", "warranty", "specification",
	}
}

func Test_ComponentRepo_Get(t *testing.T) {
	t.Run("ComponentExist", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		rows := sqlmock.NewRows(getComponentColumns()).
			AddRow(testData.ID, testData.Inventory, testData.UserID.String(), testData.Category,
				testData.Type, testData.Description, testData.DatasheetURL, testData.ImagesURL,
				testData.PartNumber, testData.Manufacturer, testData.Managed, testData.Warranty,
				testData.Specification)

		testSetup.mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(testData.ID.String(), testData.Category).
			WillReturnRows(rows)

		comp, err := testSetup.repo.Get(testData.ID)

		assert.NoError(t, err)
		assert.NotNil(t, comp)
		assert.Equal(t, comp.Id, testData.ID)

		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_ComponentRepo_GetByUser(t *testing.T) {
	t.Run("ComponentExist", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		rows := sqlmock.NewRows(getComponentColumns()).
			AddRow(testData.ID, testData.Inventory, testData.UserID.String(), testData.Category,
				testData.Type, testData.Description, testData.DatasheetURL, testData.ImagesURL,
				testData.PartNumber, testData.Manufacturer, testData.Managed, testData.Warranty,
				testData.Specification)

		testSetup.mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(testData.UserID.String(), int32(testData.Category)).
			WillReturnRows(rows)

		comps, err := testSetup.repo.GetByUser(testData.UserID.String(), int32(testData.Category))

		assert.NoError(t, err)
		assert.NotNil(t, comps)
		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NoComponentsFound", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		testSetup.mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(testData.UserID.String(), int32(testData.Category)).
			WillReturnRows(sqlmock.NewRows(getComponentColumns()))

		comps, err := testSetup.repo.GetByUser(testData.UserID.String(), int32(testData.Category))

		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, comps)
		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		testSetup.mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(testData.UserID.String(), int32(testData.Category)).
			WillReturnError(fmt.Errorf("database connection error"))

		comps, err := testSetup.repo.GetByUser(testData.UserID.String(), int32(testData.Category))

		assert.Error(t, err)
		assert.Nil(t, comps)
		assert.Contains(t, err.Error(), "database connection error")
		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_ComponentRepo_Add(t *testing.T) {
	t.Run("AddComponent", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		components := []*component_db.Component{convertToComponent(testData)}

		testSetup.mock.ExpectBegin()
		for _, component := range components {
			testSetup.mock.ExpectExec(`INSERT INTO "components" \("id","inventory","user_id","category","type","description","datasheet_url","images_url","part_number","manufacturer","managed","warranty","specification"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11,\$12,\$13\) ON CONFLICT \("id"\) DO NOTHING`).
				WithArgs(component.Id, component.Inventory, component.UserId, component.Category, component.Type, component.Description, component.DatasheetURL, component.ImagesURL, component.PartNumber, component.Manufacturer, component.Managed, component.Warranty, component.Specification).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		testSetup.mock.ExpectCommit()

		testSetup.mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(testData.UserID.String(), int32(testData.Category)).WillReturnRows(sqlmock.NewRows(getComponentColumns()).AddRow(
			testData.ID, testData.Inventory, testData.UserID, testData.Category, testData.Type, testData.Description, testData.DatasheetURL, testData.ImagesURL, testData.PartNumber, testData.Manufacturer, testData.Managed, testData.Warranty, testData.Specification,
		))

		err := testSetup.repo.Add(components)
		assert.NoError(t, err)

		res, err := testSetup.repo.GetByUser(testData.UserID.String(), int32(testData.Category))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_ComponentRepo_Delete(t *testing.T) {
	t.Run("DeleteComponent", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		components := []*component_db.Component{convertToComponent(testData)}

		testSetup.mock.ExpectBegin()
		testSetup.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "components" ("id","inventory","user_id","category","type","description","datasheet_url","images_url","part_number","manufacturer","managed","warranty","specification") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) ON CONFLICT ("id") DO NOTHING`)).
			WithArgs(components[0].Id, components[0].Inventory, components[0].UserId, components[0].Category, components[0].Type, components[0].Description, components[0].DatasheetURL, components[0].ImagesURL, components[0].PartNumber, components[0].Manufacturer, components[0].Managed, components[0].Warranty, components[0].Specification).
			WillReturnResult(sqlmock.NewResult(1, 1))
		testSetup.mock.ExpectCommit()

		testSetup.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM components`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		testSetup.mock.ExpectQuery(`^SELECT.*components.*`).
			WithArgs(components[0].Id, testData.Category).WillReturnRows(sqlmock.NewRows([]string{}))

		err := testSetup.repo.Add(components)
		assert.NoError(t, err)

		err = testSetup.repo.Delete()
		assert.NoError(t, err)

		res, err := testSetup.repo.Get(testData.ID)
		assert.Error(t, err)
		assert.Empty(t, res)

		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_ComponentRepo_List(t *testing.T) {
	t.Run("ListWithNoFilters", func(t *testing.T) {
		testSetup := setupTestDB(t)

		// Expect query with no WHERE clauses
		testSetup.mock.ExpectQuery(`^SELECT.*components.*`).
			WillReturnRows(sqlmock.NewRows(getComponentColumns()))

		components, err := testSetup.repo.List("", "", "", 0)

		assert.NoError(t, err)
		assert.NotNil(t, components)
		assert.Len(t, components, 0)

		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ListWithIDFilter", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		// Expect query with ID filter
		testSetup.mock.ExpectQuery(`^SELECT.*components.*WHERE.*id.*`).
			WithArgs(testData.ID.String()).
			WillReturnRows(sqlmock.NewRows(getComponentColumns()).AddRow(
				testData.ID, testData.Inventory, testData.UserID, testData.Category, testData.Type, testData.Description, testData.DatasheetURL, testData.ImagesURL, testData.PartNumber, testData.Manufacturer, testData.Managed, testData.Warranty, testData.Specification,
			))

		components, err := testSetup.repo.List(testData.ID.String(), "", "", 0)

		assert.NoError(t, err)
		assert.NotNil(t, components)
		assert.Len(t, components, 1)
		assert.Equal(t, testData.ID, components[0].Id)

		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ListWithUserIDFilter", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		// Expect query with UserID filter
		testSetup.mock.ExpectQuery(`^SELECT.*components.*WHERE.*user_id.*`).
			WithArgs(testData.UserID.String()).
			WillReturnRows(sqlmock.NewRows(getComponentColumns()).AddRow(
				testData.ID, testData.Inventory, testData.UserID, testData.Category, testData.Type, testData.Description, testData.DatasheetURL, testData.ImagesURL, testData.PartNumber, testData.Manufacturer, testData.Managed, testData.Warranty, testData.Specification,
			))

		components, err := testSetup.repo.List("", testData.UserID.String(), "", 0)

		assert.NoError(t, err)
		assert.NotNil(t, components)
		assert.Len(t, components, 1)
		assert.Equal(t, testData.UserID, components[0].UserId)

		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ListWithPartNumberFilter", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		// Expect query with PartNumber filter
		testSetup.mock.ExpectQuery(`^SELECT.*components.*WHERE.*part_number.*`).
			WithArgs(testData.PartNumber).
			WillReturnRows(sqlmock.NewRows(getComponentColumns()).AddRow(
				testData.ID, testData.Inventory, testData.UserID, testData.Category, testData.Type, testData.Description, testData.DatasheetURL, testData.ImagesURL, testData.PartNumber, testData.Manufacturer, testData.Managed, testData.Warranty, testData.Specification,
			))

		components, err := testSetup.repo.List("", "", testData.PartNumber, 0)

		assert.NoError(t, err)
		assert.NotNil(t, components)
		assert.Len(t, components, 1)
		assert.Equal(t, testData.PartNumber, components[0].PartNumber)

		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ListWithCategoryFilter", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		// Expect query with Category filter
		testSetup.mock.ExpectQuery(`^SELECT.*components.*WHERE.*category.*`).
			WithArgs(int32(testData.Category)).
			WillReturnRows(sqlmock.NewRows(getComponentColumns()).AddRow(
				testData.ID, testData.Inventory, testData.UserID, testData.Category, testData.Type, testData.Description, testData.DatasheetURL, testData.ImagesURL, testData.PartNumber, testData.Manufacturer, testData.Managed, testData.Warranty, testData.Specification,
			))

		components, err := testSetup.repo.List("", "", "", int32(testData.Category))

		assert.NoError(t, err)
		assert.NotNil(t, components)
		assert.Len(t, components, 1)
		assert.Equal(t, testData.Category, components[0].Category)

		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ListWithMultipleFilters", func(t *testing.T) {
		testData := createTestComponentData()
		testSetup := setupTestDB(t)

		// Expect query with multiple filters (ID, UserID, PartNumber, Category)
		testSetup.mock.ExpectQuery(`^SELECT.*components.*WHERE.*id.*AND.*user_id.*AND.*part_number.*AND.*category.*`).
			WithArgs(testData.ID.String(), testData.UserID.String(), testData.PartNumber, int32(testData.Category)).
			WillReturnRows(sqlmock.NewRows(getComponentColumns()).AddRow(
				testData.ID, testData.Inventory, testData.UserID, testData.Category, testData.Type, testData.Description, testData.DatasheetURL, testData.ImagesURL, testData.PartNumber, testData.Manufacturer, testData.Managed, testData.Warranty, testData.Specification,
			))

		components, err := testSetup.repo.List(testData.ID.String(), testData.UserID.String(), testData.PartNumber, int32(testData.Category))

		assert.NoError(t, err)
		assert.NotNil(t, components)
		assert.Len(t, components, 1)
		assert.Equal(t, testData.ID, components[0].Id)
		assert.Equal(t, testData.UserID, components[0].UserId)
		assert.Equal(t, testData.PartNumber, components[0].PartNumber)
		assert.Equal(t, testData.Category, components[0].Category)

		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ListWithDatabaseError", func(t *testing.T) {
		testSetup := setupTestDB(t)

		// Expect query to return error
		testSetup.mock.ExpectQuery(`^SELECT.*components.*`).
			WillReturnError(fmt.Errorf("database connection error"))

		components, err := testSetup.repo.List("", "", "", 0)

		assert.Error(t, err)
		assert.Nil(t, components)
		assert.Contains(t, err.Error(), "database connection error")

		err = testSetup.mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
