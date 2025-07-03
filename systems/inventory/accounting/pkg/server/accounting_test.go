/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/accounting/mocks"
	pb "github.com/ukama/ukama/systems/inventory/accounting/pb/gen"
	"github.com/ukama/ukama/systems/inventory/accounting/pkg/db"
	"github.com/ukama/ukama/systems/inventory/accounting/pkg/utils"
	"gorm.io/gorm"
)

// TestData holds all test data for the accounting tests
type TestData struct {
	OrgName              string
	Vat                  string
	Item                 string
	Inventory            string
	OpexFee              string
	EffectiveDate        string
	Description          string
	Company              string
	BranchName           string
	Email                string
	ItemName             string
	ItemDesc             string
	InventoryQty         string
	OpexFeeAmount        string
	VatAmount            string
	BackhaulItem         string
	BackhaulDesc         string
	BackhaulInventory    string
	BackhaulOpexFee      string
	BackhaulVat          string
	ConnectivityProvider struct {
		Company string
		POC     string
		Address string
		Phone   string
		Email   string
	}
	Nodes struct {
		Inventory string
		OnOrder   string
	}
	UkamaItems    []UkamaItem
	BackhaulItems []BackhaulItem
}

// UkamaItem represents a ukama item in test data
type UkamaItem struct {
	Item          string
	Description   string
	Inventory     string
	OpexFee       string
	Vat           string
	EffectiveDate string
}

// BackhaulItem represents a backhaul item in test data
type BackhaulItem struct {
	Item          string
	Description   string
	Inventory     string
	OpexFee       string
	Vat           string
	EffectiveDate string
}

// Error messages
const (
	InvalidUUIDError       = "invalid format of accounting uuid"
	InvalidDBError         = "invalid db"
	GitCloneError          = "git clone failed"
	GitCloneErrorMsg       = "failed to clone git repo"
	FileReadError          = "failed to read root file"
	JSONUnmarshalError     = "failed to unmarshal root json file"
	BranchCheckoutError    = "failed to checkout branch"
	ManifestReadError      = "failed to read manifest file"
	ManifestUnmarshalError = "failed to unmarshal manifest json file"
	UUIDParsingError       = "Error parsing UUID"
)

var (
	testUserId uuid.UUID
	testAccId  uuid.UUID
	testData   *TestData
)

func init() {
	testUserId = uuid.NewV4()
	testAccId = uuid.NewV4()

	// Initialize test data
	testData = &TestData{
		OrgName:           "testorg",
		Vat:               "10",
		Item:              "Product-1",
		Inventory:         "1",
		OpexFee:           "100",
		EffectiveDate:     "2023-01-01",
		Description:       "Product-1 description",
		Company:           "TestCompany",
		BranchName:        "test-branch",
		Email:             "test@example.com",
		ItemName:          "TestItem",
		ItemDesc:          "Test Description",
		InventoryQty:      "10",
		OpexFeeAmount:     "100.00",
		VatAmount:         "20.00",
		BackhaulItem:      "BackhaulItem",
		BackhaulDesc:      "Backhaul Description",
		BackhaulInventory: "5",
		BackhaulOpexFee:   "50.00",
		BackhaulVat:       "10.00",
		ConnectivityProvider: struct {
			Company string
			POC     string
			Address string
			Phone   string
			Email   string
		}{
			Company: "TestConnectivityProvider",
			POC:     "John Doe",
			Address: "123 Test Street, Test City",
			Phone:   "+1-555-123-4567",
			Email:   "contact@testprovider.com",
		},
		Nodes: struct {
			Inventory string
			OnOrder   string
		}{
			Inventory: "100",
			OnOrder:   "50",
		},
		UkamaItems: []UkamaItem{
			{
				Item:          "TestItem",
				Description:   "Test Description",
				Inventory:     "10",
				OpexFee:       "100.00",
				Vat:           "20.00",
				EffectiveDate: "2023-01-01",
			},
			{
				Item:          "UkamaItem2",
				Description:   "Second Ukama Item Description",
				Inventory:     "25",
				OpexFee:       "150.00",
				Vat:           "30.00",
				EffectiveDate: "2023-02-01",
			},
		},
		BackhaulItems: []BackhaulItem{
			{
				Item:          "BackhaulItem",
				Description:   "Backhaul Description",
				Inventory:     "5",
				OpexFee:       "50.00",
				Vat:           "10.00",
				EffectiveDate: "2023-01-01",
			},
			{
				Item:          "BackhaulItem2",
				Description:   "Second Backhaul Item Description",
				Inventory:     "15",
				OpexFee:       "75.00",
				Vat:           "15.00",
				EffectiveDate: "2023-03-01",
			},
		},
	}
}

// Helper functions to generate test data
func (td *TestData) generateRootJSON(userId uuid.UUID) string {
	return `{"test":[{"company":"` + td.Company + `","git_branch_name":"` + td.BranchName + `","email":"` + td.Email + `","user_id":"` + userId.String() + `"}]}`
}

func (td *TestData) generateManifestJSON() string {
	ukamaItems := ""
	for i, item := range td.UkamaItems {
		if i > 0 {
			ukamaItems += ","
		}
		ukamaItems += `{
			"item": "` + item.Item + `",
			"description": "` + item.Description + `",
			"inventory": "` + item.Inventory + `",
			"opex_fee": "` + item.OpexFee + `",
			"vat": "` + item.Vat + `",
			"effective_date": "` + item.EffectiveDate + `"
		}`
	}

	backhaulItems := ""
	for i, item := range td.BackhaulItems {
		if i > 0 {
			backhaulItems += ","
		}
		backhaulItems += `{
			"item": "` + item.Item + `",
			"description": "` + item.Description + `",
			"inventory": "` + item.Inventory + `",
			"opex_fee": "` + item.OpexFee + `",
			"vat": "` + item.Vat + `",
			"effective_date": "` + item.EffectiveDate + `"
		}`
	}

	return `{
		"effective_date": "` + td.EffectiveDate + `",
		"connectivityProvider": {
			"company": "` + td.ConnectivityProvider.Company + `",
			"poc": "` + td.ConnectivityProvider.POC + `",
			"address": "` + td.ConnectivityProvider.Address + `",
			"phone": "` + td.ConnectivityProvider.Phone + `",
			"email": "` + td.ConnectivityProvider.Email + `"
		},
		"nodes": {
			"inventory": "` + td.Nodes.Inventory + `",
			"onOrder": "` + td.Nodes.OnOrder + `"
		},
		"ukama": [` + ukamaItems + `],
		"backhaul": [` + backhaulItems + `]
	}`
}

func (td *TestData) generateMinimalManifestJSON() string {
	return `{
		"effective_date": "` + td.EffectiveDate + `",
		"connectivityProvider": {
			"company": "MinimalProvider",
			"poc": "Jane Smith",
			"address": "456 Minimal Ave",
			"phone": "+1-555-987-6543",
			"email": "jane@minimal.com"
		},
		"nodes": {
			"inventory": "10",
			"onOrder": "5"
		},
		"ukama": [
			{
				"item": "MinimalUkamaItem",
				"description": "Minimal Ukama Description",
				"inventory": "1",
				"opex_fee": "50.00",
				"vat": "10.00",
				"effective_date": "` + td.EffectiveDate + `"
			}
		],
		"backhaul": [
			{
				"item": "MinimalBackhaulItem",
				"description": "Minimal Backhaul Description",
				"inventory": "2",
				"opex_fee": "25.00",
				"vat": "5.00",
				"effective_date": "` + td.EffectiveDate + `"
			}
		]
	}`
}

func (td *TestData) generateSimpleManifestJSON() string {
	return `{"ukama":[{"item":"` + td.ItemName + `","description":"` + td.ItemDesc + `","inventory":"` + td.InventoryQty + `","opex_fee":"` + td.OpexFeeAmount + `","vat":"` + td.VatAmount + `","effective_date":"` + td.EffectiveDate + `"}]}`
}

func (td *TestData) generateRootJSONWithInvalidUUID() string {
	return `{"test":[{"company":"` + td.Company + `","git_branch_name":"` + td.BranchName + `","email":"` + td.Email + `","user_id":"invalid-uuid"}]}`
}

func (td *TestData) generateAccountingRecords(userId uuid.UUID) []*db.Accounting {
	var records []*db.Accounting

	// Add ukama items
	for _, item := range td.UkamaItems {
		records = append(records, &db.Accounting{
			Id:            uuid.NewV4(),
			UserId:        userId,
			Item:          item.Item,
			Description:   item.Description,
			Inventory:     item.Inventory,
			OpexFee:       item.OpexFee,
			Vat:           item.Vat,
			EffectiveDate: item.EffectiveDate,
		})
	}

	// Add backhaul items
	for _, item := range td.BackhaulItems {
		records = append(records, &db.Accounting{
			Id:            uuid.NewV4(),
			UserId:        userId,
			Item:          item.Item,
			Description:   item.Description,
			Inventory:     item.Inventory,
			OpexFee:       item.OpexFee,
			Vat:           item.Vat,
			EffectiveDate: item.EffectiveDate,
		})
	}

	return records
}

func (td *TestData) generateMinimalAccountingRecords(userId uuid.UUID) []*db.Accounting {
	return []*db.Accounting{
		{
			Id:            uuid.NewV4(),
			UserId:        userId,
			Item:          "MinimalUkamaItem",
			Description:   "Minimal Ukama Description",
			Inventory:     "1",
			OpexFee:       "50.00",
			Vat:           "10.00",
			EffectiveDate: td.EffectiveDate,
		},
		{
			Id:            uuid.NewV4(),
			UserId:        userId,
			Item:          "MinimalBackhaulItem",
			Description:   "Minimal Backhaul Description",
			Inventory:     "2",
			OpexFee:       "25.00",
			Vat:           "5.00",
			EffectiveDate: td.EffectiveDate,
		},
	}
}

func (td *TestData) generateManifestJSONWithEmptyArrays() string {
	return `{
		"effective_date": "` + td.EffectiveDate + `",
		"connectivityProvider": {
			"company": "` + td.ConnectivityProvider.Company + `",
			"poc": "` + td.ConnectivityProvider.POC + `",
			"address": "` + td.ConnectivityProvider.Address + `",
			"phone": "` + td.ConnectivityProvider.Phone + `",
			"email": "` + td.ConnectivityProvider.Email + `"
		},
		"nodes": {
			"inventory": "` + td.Nodes.Inventory + `",
			"onOrder": "` + td.Nodes.OnOrder + `"
		},
		"ukama": [],
		"backhaul": []
	}`
}

func (td *TestData) generateManifestJSONWithOnlyUkama() string {
	ukamaItems := ""
	for i, item := range td.UkamaItems {
		if i > 0 {
			ukamaItems += ","
		}
		ukamaItems += `{
			"item": "` + item.Item + `",
			"description": "` + item.Description + `",
			"inventory": "` + item.Inventory + `",
			"opex_fee": "` + item.OpexFee + `",
			"vat": "` + item.Vat + `",
			"effective_date": "` + item.EffectiveDate + `"
		}`
	}

	return `{
		"effective_date": "` + td.EffectiveDate + `",
		"connectivityProvider": {
			"company": "` + td.ConnectivityProvider.Company + `",
			"poc": "` + td.ConnectivityProvider.POC + `",
			"address": "` + td.ConnectivityProvider.Address + `",
			"phone": "` + td.ConnectivityProvider.Phone + `",
			"email": "` + td.ConnectivityProvider.Email + `"
		},
		"nodes": {
			"inventory": "` + td.Nodes.Inventory + `",
			"onOrder": "` + td.Nodes.OnOrder + `"
		},
		"ukama": [` + ukamaItems + `],
		"backhaul": []
	}`
}

func (td *TestData) generateManifestJSONWithOnlyBackhaul() string {
	backhaulItems := ""
	for i, item := range td.BackhaulItems {
		if i > 0 {
			backhaulItems += ","
		}
		backhaulItems += `{
			"item": "` + item.Item + `",
			"description": "` + item.Description + `",
			"inventory": "` + item.Inventory + `",
			"opex_fee": "` + item.OpexFee + `",
			"vat": "` + item.Vat + `",
			"effective_date": "` + item.EffectiveDate + `"
		}`
	}

	return `{
		"effective_date": "` + td.EffectiveDate + `",
		"connectivityProvider": {
			"company": "` + td.ConnectivityProvider.Company + `",
			"poc": "` + td.ConnectivityProvider.POC + `",
			"address": "` + td.ConnectivityProvider.Address + `",
			"phone": "` + td.ConnectivityProvider.Phone + `",
			"email": "` + td.ConnectivityProvider.Email + `"
		},
		"nodes": {
			"inventory": "` + td.Nodes.Inventory + `",
			"onOrder": "` + td.Nodes.OnOrder + `"
		},
		"ukama": [],
		"backhaul": [` + backhaulItems + `]
	}`
}

func (td *TestData) generateManifestJSONWithMissingFields() string {
	return `{
		"effective_date": "` + td.EffectiveDate + `",
		"connectivityProvider": {
			"company": "` + td.ConnectivityProvider.Company + `",
			"poc": "` + td.ConnectivityProvider.POC + `",
			"address": "` + td.ConnectivityProvider.Address + `",
			"phone": "` + td.ConnectivityProvider.Phone + `",
			"email": "` + td.ConnectivityProvider.Email + `"
		},
		"nodes": {
			"inventory": "` + td.Nodes.Inventory + `",
			"onOrder": "` + td.Nodes.OnOrder + `"
		},
		"ukama": [
			{
				"item": "TestItem",
				"description": "Test Description",
				"inventory": "10",
				"opex_fee": "100.00",
				"vat": "20.00"
			}
		],
		"backhaul": [
			{
				"item": "BackhaulItem",
				"description": "Backhaul Description",
				"inventory": "5",
				"opex_fee": "50.00",
				"vat": "10.00"
			}
		]
	}`
}

func (td *TestData) generateManifestJSONWithNullValues() string {
	return `{
		"effective_date": "` + td.EffectiveDate + `",
		"connectivityProvider": {
			"company": "` + td.ConnectivityProvider.Company + `",
			"poc": "` + td.ConnectivityProvider.POC + `",
			"address": "` + td.ConnectivityProvider.Address + `",
			"phone": "` + td.ConnectivityProvider.Phone + `",
			"email": "` + td.ConnectivityProvider.Email + `"
		},
		"nodes": {
			"inventory": "` + td.Nodes.Inventory + `",
			"onOrder": "` + td.Nodes.OnOrder + `"
		},
		"ukama": [
			{
				"item": "TestItem",
				"description": null,
				"inventory": "10",
				"opex_fee": "100.00",
				"vat": "20.00",
				"effective_date": "` + td.EffectiveDate + `"
			}
		],
		"backhaul": [
			{
				"item": "BackhaulItem",
				"description": "Backhaul Description",
				"inventory": null,
				"opex_fee": "50.00",
				"vat": "10.00",
				"effective_date": "` + td.EffectiveDate + `"
			}
		]
	}`
}

func (td *TestData) generateManifestJSONWithSpecialCharacters() string {
	return `{
		"effective_date": "` + td.EffectiveDate + `",
		"connectivityProvider": {
			"company": "Test & Co. Ltd.",
			"poc": "John O'Connor",
			"address": "123 Main St., Suite #100, City, State 12345",
			"phone": "+1-555-123-4567",
			"email": "test+tag@example.com"
		},
		"nodes": {
			"inventory": "100",
			"onOrder": "50"
		},
		"ukama": [
			{
				"item": "Test-Item_123",
				"description": "Test Description with special chars: & < > \" '",
				"inventory": "10",
				"opex_fee": "100.50",
				"vat": "20.25",
				"effective_date": "` + td.EffectiveDate + `"
			}
		],
		"backhaul": [
			{
				"item": "Backhaul-Item_456",
				"description": "Backhaul with special chars: & < > \" '",
				"inventory": "5",
				"opex_fee": "50.75",
				"vat": "10.15",
				"effective_date": "` + td.EffectiveDate + `"
			}
		]
	}`
}

func (td *TestData) generateManifestJSONWithLargeNumbers() string {
	return `{
		"effective_date": "` + td.EffectiveDate + `",
		"connectivityProvider": {
			"company": "` + td.ConnectivityProvider.Company + `",
			"poc": "` + td.ConnectivityProvider.POC + `",
			"address": "` + td.ConnectivityProvider.Address + `",
			"phone": "` + td.ConnectivityProvider.Phone + `",
			"email": "` + td.ConnectivityProvider.Email + `"
		},
		"nodes": {
			"inventory": "999999",
			"onOrder": "999999"
		},
		"ukama": [
			{
				"item": "LargeNumberItem",
				"description": "Item with large numbers",
				"inventory": "999999",
				"opex_fee": "999999.99",
				"vat": "999999.99",
				"effective_date": "` + td.EffectiveDate + `"
			}
		],
		"backhaul": [
			{
				"item": "LargeNumberBackhaul",
				"description": "Backhaul with large numbers",
				"inventory": "999999",
				"opex_fee": "999999.99",
				"vat": "999999.99",
				"effective_date": "` + td.EffectiveDate + `"
			}
		]
	}`
}

func TestAccountingServer_Get(t *testing.T) {
	t.Run("Account record found", func(t *testing.T) {
		var aId = uuid.NewV4()

		accRepo := &mocks.AccountingRepo{}

		accRepo.On("Get", aId).Return(
			&db.Accounting{Id: aId}, nil).Once()

		s := NewAccountingServer(testData.OrgName, accRepo, nil, "", nil, "")
		accResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: aId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, accResp)
		assert.Equal(t, aId.String(), accResp.Accounting.Id)

		accRepo.AssertExpectations(t)
	})

	t.Run("Account record not found", func(t *testing.T) {
		var aId = uuid.NewV4()

		accRepo := &mocks.AccountingRepo{}

		accRepo.On("Get", aId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewAccountingServer(testData.OrgName, accRepo, nil, "", nil, "")
		accResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: aId.String()})

		assert.Error(t, err)
		assert.Nil(t, accResp)
		accRepo.AssertExpectations(t)
	})

	t.Run("Invalid UUID format", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}

		s := NewAccountingServer(testData.OrgName, accRepo, nil, "", nil, "")
		accResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: "invalid-uuid-format"})

		assert.Error(t, err)
		assert.Nil(t, accResp)
		assert.Contains(t, err.Error(), InvalidUUIDError)
		accRepo.AssertExpectations(t)
	})
}

func TestAccountingServer_GetByUser(t *testing.T) {
	t.Run("Account records by user found", func(t *testing.T) {
		var uId = uuid.NewV4()

		accRepo := &mocks.AccountingRepo{}

		accRepo.On("GetByUser", uId.String()).Return(
			testData.generateAccountingRecords(uId), nil).Once()

		s := NewAccountingServer(testData.OrgName, accRepo, nil, "", nil, "")

		accResp, err := s.GetByUser(context.TODO(),
			&pb.GetByUserRequest{
				UserId: uId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, accResp)
		assert.Equal(t, uId.String(), accResp.GetAccounting()[0].GetUserId())
		accRepo.AssertExpectations(t)
	})

	t.Run("Database error when getting accountings by user", func(t *testing.T) {
		var uId = uuid.NewV4()

		accRepo := &mocks.AccountingRepo{}

		accRepo.On("GetByUser", uId.String()).Return(nil, gorm.ErrInvalidDB).Once()

		s := NewAccountingServer(testData.OrgName, accRepo, nil, "", nil, "")

		accResp, err := s.GetByUser(context.TODO(),
			&pb.GetByUserRequest{
				UserId: uId.String()})

		assert.Error(t, err)
		assert.Nil(t, accResp)
		assert.Contains(t, err.Error(), InvalidDBError)
		accRepo.AssertExpectations(t)
	})
}

func TestAccountingServer_SyncAccounting(t *testing.T) {
	t.Run("Successfully sync accounting", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()

		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSON()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(testData.generateAccountingRecords(userId), nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Failed to clone git repo", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(errors.New(GitCloneError))

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), GitCloneErrorMsg)

		gitClient.AssertExpectations(t)
	})

	// Negative test cases for comprehensive coverage
	t.Run("Failed to read root.json file", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return(nil, errors.New("file not found"))

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to read root file")

		gitClient.AssertExpectations(t)
	})

	t.Run("Invalid JSON in root.json file", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte("invalid json"), nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to unmarshal root json file")

		gitClient.AssertExpectations(t)
	})

	t.Run("Branch checkout failure", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()
		rootJSON := testData.generateRootJSON(userId)

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(errors.New("branch not found"))

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to checkout branch")

		gitClient.AssertExpectations(t)
	})

	t.Run("Failed to read manifest.json file", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()
		rootJSON := testData.generateRootJSON(userId)

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return(nil, errors.New("manifest file not found"))

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to read manifest file")

		gitClient.AssertExpectations(t)
	})

	t.Run("Invalid JSON in manifest.json file", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()
		rootJSON := testData.generateRootJSON(userId)

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte("invalid json"), nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to unmarshal manifest json file")

		gitClient.AssertExpectations(t)
	})

	t.Run("Invalid user ID in company data", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		rootJSON := testData.generateRootJSONWithInvalidUUID()
		manifestJSON := testData.generateSimpleManifestJSON()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), UUIDParsingError)

		gitClient.AssertExpectations(t)
	})

	t.Run("Database delete failure", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()
		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSON()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(errors.New("database delete failed"))

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database delete failed")

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
	})

	t.Run("Database add failure", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()
		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSON()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(errors.New("database add failed"))

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database add failed")

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
	})

	t.Run("Database GetByUser failure", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()
		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSON()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(nil, errors.New("database get by user failed"))

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database get by user failed")

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
	})

	t.Run("Message bus publish failure", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()
		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSON()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(testData.generateAccountingRecords(userId), nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("publish failed"))

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Test utils.Accounting struct parsing with minimal data", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()

		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateMinimalManifestJSON()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(testData.generateMinimalAccountingRecords(userId), nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Test utils.Accounting struct parsing with empty arrays", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()

		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSONWithEmptyArrays()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return([]*db.Accounting{}, nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Test utils.Accounting struct parsing with only ukama items", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()

		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSONWithOnlyUkama()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(testData.generateAccountingRecords(userId), nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Test utils.Accounting struct parsing with only backhaul items", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()

		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSONWithOnlyBackhaul()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(testData.generateAccountingRecords(userId), nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Test utils.Accounting struct parsing with missing fields", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()

		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSONWithMissingFields()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(testData.generateAccountingRecords(userId), nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Test utils.Accounting struct parsing with null values", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()

		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSONWithNullValues()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(testData.generateAccountingRecords(userId), nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Test utils.Accounting struct parsing with special characters", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()

		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSONWithSpecialCharacters()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(testData.generateAccountingRecords(userId), nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Test utils.Accounting struct parsing with large numbers", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()

		rootJSON := testData.generateRootJSON(userId)
		manifestJSON := testData.generateManifestJSONWithLargeNumbers()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(testData.generateAccountingRecords(userId), nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Test utils.Accounting struct parsing with multiple companies", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId1 := uuid.NewV4()
		userId2 := uuid.NewV4()

		rootJSON := `{"test":[
			{"company":"` + testData.Company + `","git_branch_name":"` + testData.BranchName + `","email":"` + testData.Email + `","user_id":"` + userId1.String() + `"},
			{"company":"` + testData.Company + `2","git_branch_name":"` + testData.BranchName + `2","email":"` + testData.Email + `2","user_id":"` + userId2.String() + `"}
		]}`
		manifestJSON := testData.generateManifestJSON()

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)
		gitClient.On("BranchCheckout", testData.BranchName+"2").Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil).Times(2)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil).Times(2)
		accRepo.On("GetByUser", userId1.String()).Return(testData.generateAccountingRecords(userId1), nil)
		accRepo.On("GetByUser", userId2.String()).Return(testData.generateAccountingRecords(userId2), nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil).Times(2)

		s := NewAccountingServer(testData.OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})
}

// TestUtilsAccountingStructParsing tests the utils.Accounting struct parsing directly
func TestUtilsAccountingStructParsing(t *testing.T) {
	t.Run("Parse complete utils.Accounting struct", func(t *testing.T) {
		manifestJSON := testData.generateManifestJSON()

		var accounting utils.Accounting
		err := json.Unmarshal([]byte(manifestJSON), &accounting)

		assert.NoError(t, err)
		assert.Equal(t, testData.EffectiveDate, accounting.EffectiveDate)
		assert.Equal(t, testData.ConnectivityProvider.Company, accounting.ConnectivityProvider.Company)
		assert.Equal(t, testData.ConnectivityProvider.POC, accounting.ConnectivityProvider.Poc)
		assert.Equal(t, testData.ConnectivityProvider.Address, accounting.ConnectivityProvider.Address)
		assert.Equal(t, testData.ConnectivityProvider.Phone, accounting.ConnectivityProvider.Phone)
		assert.Equal(t, testData.ConnectivityProvider.Email, accounting.ConnectivityProvider.Email)
		assert.Equal(t, testData.Nodes.Inventory, accounting.Nodes.Inventory)
		assert.Equal(t, testData.Nodes.OnOrder, accounting.Nodes.OnOrder)
		assert.Len(t, accounting.Ukama, len(testData.UkamaItems))
		assert.Len(t, accounting.Backhaul, len(testData.BackhaulItems))

		// Verify ukama items
		for i, item := range accounting.Ukama {
			assert.Equal(t, testData.UkamaItems[i].Item, item.Item)
			assert.Equal(t, testData.UkamaItems[i].Description, item.Description)
			assert.Equal(t, testData.UkamaItems[i].Inventory, item.Inventory)
			assert.Equal(t, testData.UkamaItems[i].OpexFee, item.OpexFee)
			assert.Equal(t, testData.UkamaItems[i].Vat, item.Vat)
			assert.Equal(t, testData.UkamaItems[i].EffectiveDate, item.EffectiveDate)
		}

		// Verify backhaul items
		for i, item := range accounting.Backhaul {
			assert.Equal(t, testData.BackhaulItems[i].Item, item.Item)
			assert.Equal(t, testData.BackhaulItems[i].Description, item.Description)
			assert.Equal(t, testData.BackhaulItems[i].Inventory, item.Inventory)
			assert.Equal(t, testData.BackhaulItems[i].OpexFee, item.OpexFee)
			assert.Equal(t, testData.BackhaulItems[i].Vat, item.Vat)
			assert.Equal(t, testData.BackhaulItems[i].EffectiveDate, item.EffectiveDate)
		}
	})

	t.Run("Parse utils.Accounting struct with empty arrays", func(t *testing.T) {
		manifestJSON := testData.generateManifestJSONWithEmptyArrays()

		var accounting utils.Accounting
		err := json.Unmarshal([]byte(manifestJSON), &accounting)

		assert.NoError(t, err)
		assert.Equal(t, testData.EffectiveDate, accounting.EffectiveDate)
		assert.Empty(t, accounting.Ukama)
		assert.Empty(t, accounting.Backhaul)
	})

	t.Run("Parse utils.Accounting struct with missing fields", func(t *testing.T) {
		manifestJSON := testData.generateManifestJSONWithMissingFields()

		var accounting utils.Accounting
		err := json.Unmarshal([]byte(manifestJSON), &accounting)

		assert.NoError(t, err)
		assert.Len(t, accounting.Ukama, 1)
		assert.Len(t, accounting.Backhaul, 1)

		// Check that missing fields are handled gracefully (empty strings)
		assert.Equal(t, "TestItem", accounting.Ukama[0].Item)
		assert.Equal(t, "Test Description", accounting.Ukama[0].Description)
		assert.Equal(t, "10", accounting.Ukama[0].Inventory)
		assert.Equal(t, "100.00", accounting.Ukama[0].OpexFee)
		assert.Equal(t, "20.00", accounting.Ukama[0].Vat)
		assert.Equal(t, "", accounting.Ukama[0].EffectiveDate) // Missing field

		assert.Equal(t, "BackhaulItem", accounting.Backhaul[0].Item)
		assert.Equal(t, "Backhaul Description", accounting.Backhaul[0].Description)
		assert.Equal(t, "5", accounting.Backhaul[0].Inventory)
		assert.Equal(t, "50.00", accounting.Backhaul[0].OpexFee)
		assert.Equal(t, "10.00", accounting.Backhaul[0].Vat)
		assert.Equal(t, "", accounting.Backhaul[0].EffectiveDate) // Missing field
	})

	t.Run("Parse utils.Accounting struct with null values", func(t *testing.T) {
		manifestJSON := testData.generateManifestJSONWithNullValues()

		var accounting utils.Accounting
		err := json.Unmarshal([]byte(manifestJSON), &accounting)

		assert.NoError(t, err)
		assert.Len(t, accounting.Ukama, 1)
		assert.Len(t, accounting.Backhaul, 1)

		// Check that null values are handled gracefully (empty strings)
		assert.Equal(t, "TestItem", accounting.Ukama[0].Item)
		assert.Equal(t, "", accounting.Ukama[0].Description) // null value
		assert.Equal(t, "10", accounting.Ukama[0].Inventory)
		assert.Equal(t, "100.00", accounting.Ukama[0].OpexFee)
		assert.Equal(t, "20.00", accounting.Ukama[0].Vat)
		assert.Equal(t, testData.EffectiveDate, accounting.Ukama[0].EffectiveDate)

		assert.Equal(t, "BackhaulItem", accounting.Backhaul[0].Item)
		assert.Equal(t, "Backhaul Description", accounting.Backhaul[0].Description)
		assert.Equal(t, "", accounting.Backhaul[0].Inventory) // null value
		assert.Equal(t, "50.00", accounting.Backhaul[0].OpexFee)
		assert.Equal(t, "10.00", accounting.Backhaul[0].Vat)
		assert.Equal(t, testData.EffectiveDate, accounting.Backhaul[0].EffectiveDate)
	})

	t.Run("Parse utils.Accounting struct with special characters", func(t *testing.T) {
		manifestJSON := testData.generateManifestJSONWithSpecialCharacters()

		var accounting utils.Accounting
		err := json.Unmarshal([]byte(manifestJSON), &accounting)

		assert.NoError(t, err)
		assert.Equal(t, "Test & Co. Ltd.", accounting.ConnectivityProvider.Company)
		assert.Equal(t, "John O'Connor", accounting.ConnectivityProvider.Poc)
		assert.Equal(t, "123 Main St., Suite #100, City, State 12345", accounting.ConnectivityProvider.Address)
		assert.Equal(t, "+1-555-123-4567", accounting.ConnectivityProvider.Phone)
		assert.Equal(t, "test+tag@example.com", accounting.ConnectivityProvider.Email)

		assert.Len(t, accounting.Ukama, 1)
		assert.Len(t, accounting.Backhaul, 1)

		assert.Equal(t, "Test-Item_123", accounting.Ukama[0].Item)
		assert.Equal(t, "Test Description with special chars: & < > \" '", accounting.Ukama[0].Description)
		assert.Equal(t, "10", accounting.Ukama[0].Inventory)
		assert.Equal(t, "100.50", accounting.Ukama[0].OpexFee)
		assert.Equal(t, "20.25", accounting.Ukama[0].Vat)

		assert.Equal(t, "Backhaul-Item_456", accounting.Backhaul[0].Item)
		assert.Equal(t, "Backhaul with special chars: & < > \" '", accounting.Backhaul[0].Description)
		assert.Equal(t, "5", accounting.Backhaul[0].Inventory)
		assert.Equal(t, "50.75", accounting.Backhaul[0].OpexFee)
		assert.Equal(t, "10.15", accounting.Backhaul[0].Vat)
	})

	t.Run("Parse utils.Accounting struct with large numbers", func(t *testing.T) {
		manifestJSON := testData.generateManifestJSONWithLargeNumbers()

		var accounting utils.Accounting
		err := json.Unmarshal([]byte(manifestJSON), &accounting)

		assert.NoError(t, err)
		assert.Equal(t, "999999", accounting.Nodes.Inventory)
		assert.Equal(t, "999999", accounting.Nodes.OnOrder)

		assert.Len(t, accounting.Ukama, 1)
		assert.Len(t, accounting.Backhaul, 1)

		assert.Equal(t, "LargeNumberItem", accounting.Ukama[0].Item)
		assert.Equal(t, "Item with large numbers", accounting.Ukama[0].Description)
		assert.Equal(t, "999999", accounting.Ukama[0].Inventory)
		assert.Equal(t, "999999.99", accounting.Ukama[0].OpexFee)
		assert.Equal(t, "999999.99", accounting.Ukama[0].Vat)

		assert.Equal(t, "LargeNumberBackhaul", accounting.Backhaul[0].Item)
		assert.Equal(t, "Backhaul with large numbers", accounting.Backhaul[0].Description)
		assert.Equal(t, "999999", accounting.Backhaul[0].Inventory)
		assert.Equal(t, "999999.99", accounting.Backhaul[0].OpexFee)
		assert.Equal(t, "999999.99", accounting.Backhaul[0].Vat)
	})

	t.Run("Parse invalid JSON for utils.Accounting struct", func(t *testing.T) {
		invalidJSON := `{"invalid": "json", "structure":}`

		var accounting utils.Accounting
		err := json.Unmarshal([]byte(invalidJSON), &accounting)

		assert.Error(t, err)
	})

	t.Run("Parse utils.Accounting struct with only ukama items", func(t *testing.T) {
		manifestJSON := testData.generateManifestJSONWithOnlyUkama()

		var accounting utils.Accounting
		err := json.Unmarshal([]byte(manifestJSON), &accounting)

		assert.NoError(t, err)
		assert.Len(t, accounting.Ukama, len(testData.UkamaItems))
		assert.Empty(t, accounting.Backhaul)

		// Verify ukama items
		for i, item := range accounting.Ukama {
			assert.Equal(t, testData.UkamaItems[i].Item, item.Item)
			assert.Equal(t, testData.UkamaItems[i].Description, item.Description)
			assert.Equal(t, testData.UkamaItems[i].Inventory, item.Inventory)
			assert.Equal(t, testData.UkamaItems[i].OpexFee, item.OpexFee)
			assert.Equal(t, testData.UkamaItems[i].Vat, item.Vat)
			assert.Equal(t, testData.UkamaItems[i].EffectiveDate, item.EffectiveDate)
		}
	})

	t.Run("Parse utils.Accounting struct with only backhaul items", func(t *testing.T) {
		manifestJSON := testData.generateManifestJSONWithOnlyBackhaul()

		var accounting utils.Accounting
		err := json.Unmarshal([]byte(manifestJSON), &accounting)

		assert.NoError(t, err)
		assert.Empty(t, accounting.Ukama)
		assert.Len(t, accounting.Backhaul, len(testData.BackhaulItems))

		// Verify backhaul items
		for i, item := range accounting.Backhaul {
			assert.Equal(t, testData.BackhaulItems[i].Item, item.Item)
			assert.Equal(t, testData.BackhaulItems[i].Description, item.Description)
			assert.Equal(t, testData.BackhaulItems[i].Inventory, item.Inventory)
			assert.Equal(t, testData.BackhaulItems[i].OpexFee, item.OpexFee)
			assert.Equal(t, testData.BackhaulItems[i].Vat, item.Vat)
			assert.Equal(t, testData.BackhaulItems[i].EffectiveDate, item.EffectiveDate)
		}
	})
}
