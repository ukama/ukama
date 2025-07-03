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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/accounting/mocks"
	pb "github.com/ukama/ukama/systems/inventory/accounting/pb/gen"
	"github.com/ukama/ukama/systems/inventory/accounting/pkg/db"
	"gorm.io/gorm"
)

const (
	OrgName = "testorg"

	// Test data constants
	TestVat           = "10"
	TestItem          = "Product-1"
	TestInventory     = "1"
	TestOpexFee       = "100"
	TestEffectiveDate = "2023-01-01"
	TestDescription   = "Product-1 description"

	// Sync accounting test data
	TestCompany       = "TestCompany"
	TestBranchName    = "test-branch"
	TestEmail         = "test@example.com"
	TestItemName      = "TestItem"
	TestItemDesc      = "Test Description"
	TestInventoryQty  = "10"
	TestOpexFeeAmount = "100.00"
	TestVatAmount     = "20.00"
	BackhaulItem      = "BackhaulItem"
	BackhaulDesc      = "Backhaul Description"
	BackhaulInventory = "5"
	BackhaulOpexFee   = "50.00"
	BackhaulVat       = "10.00"

	// Error messages
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
)

func init() {
	testUserId = uuid.NewV4()
	testAccId = uuid.NewV4()
}

func TestAccountingServer_Get(t *testing.T) {
	t.Run("Account record found", func(t *testing.T) {
		var aId = uuid.NewV4()

		accRepo := &mocks.AccountingRepo{}

		accRepo.On("Get", aId).Return(
			&db.Accounting{Id: aId}, nil).Once()

		s := NewAccountingServer(OrgName, accRepo, nil, "", nil, "")
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

		s := NewAccountingServer(OrgName, accRepo, nil, "", nil, "")
		accResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: aId.String()})

		assert.Error(t, err)
		assert.Nil(t, accResp)
		accRepo.AssertExpectations(t)
	})

	t.Run("Invalid UUID format", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}

		s := NewAccountingServer(OrgName, accRepo, nil, "", nil, "")
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
		var aId = uuid.NewV4()

		accRepo := &mocks.AccountingRepo{}

		accRepo.On("GetByUser", uId.String()).Return(
			[]*db.Accounting{
				{
					Id:            aId,
					Vat:           TestVat,
					Item:          TestItem,
					UserId:        uId,
					Inventory:     TestInventory,
					OpexFee:       TestOpexFee,
					EffectiveDate: TestEffectiveDate,
					Description:   TestDescription,
				}}, nil).Once()

		s := NewAccountingServer(OrgName, accRepo, nil, "", nil, "")

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

		s := NewAccountingServer(OrgName, accRepo, nil, "", nil, "")

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

		rootJSON := `{"test":[{"company":"` + TestCompany + `","git_branch_name":"` + TestBranchName + `","email":"` + TestEmail + `","user_id":"` + userId.String() + `"}]}`
		manifestJSON := `{"ukama":[{"item":"` + TestItemName + `","description":"` + TestItemDesc + `","inventory":"` + TestInventoryQty + `","opex_fee":"` + TestOpexFeeAmount + `","vat":"` + TestVatAmount + `","effective_date":"` + TestEffectiveDate + `"}],"backhaul":[{"item":"` + BackhaulItem + `","description":"` + BackhaulDesc + `","inventory":"` + BackhaulInventory + `","opex_fee":"` + BackhaulOpexFee + `","vat":"` + BackhaulVat + `","effective_date":"` + TestEffectiveDate + `"}]}`

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestBranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return([]*db.Accounting{
			{
				Id:            uuid.NewV4(),
				UserId:        userId,
				Item:          TestItemName,
				Description:   TestItemDesc,
				Inventory:     TestInventoryQty,
				OpexFee:       TestOpexFeeAmount,
				Vat:           TestVatAmount,
				EffectiveDate: TestEffectiveDate,
			},
		}, nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

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

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

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

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

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

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

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
		rootJSON := `{"test":[{"company":"` + TestCompany + `","git_branch_name":"` + TestBranchName + `","email":"` + TestEmail + `","user_id":"` + userId.String() + `"}]}`

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestBranchName).Return(errors.New("branch not found"))

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

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
		rootJSON := `{"test":[{"company":"` + TestCompany + `","git_branch_name":"` + TestBranchName + `","email":"` + TestEmail + `","user_id":"` + userId.String() + `"}]}`

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestBranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return(nil, errors.New("manifest file not found"))

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

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
		rootJSON := `{"test":[{"company":"` + TestCompany + `","git_branch_name":"` + TestBranchName + `","email":"` + TestEmail + `","user_id":"` + userId.String() + `"}]}`

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestBranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte("invalid json"), nil)

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

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

		rootJSON := `{"test":[{"company":"` + TestCompany + `","git_branch_name":"` + TestBranchName + `","email":"` + TestEmail + `","user_id":"invalid-uuid"}]}`
		manifestJSON := `{"ukama":[{"item":"` + TestItemName + `","description":"` + TestItemDesc + `","inventory":"` + TestInventoryQty + `","opex_fee":"` + TestOpexFeeAmount + `","vat":"` + TestVatAmount + `","effective_date":"` + TestEffectiveDate + `"}]}`

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestBranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")

		gitClient.AssertExpectations(t)
	})

	t.Run("Database delete failure", func(t *testing.T) {
		accRepo := &mocks.AccountingRepo{}
		msgBus := &cmocks.MsgBusServiceClient{}
		gitClient := &cmocks.GitClient{}

		userId := uuid.NewV4()
		rootJSON := `{"test":[{"company":"` + TestCompany + `","git_branch_name":"` + TestBranchName + `","email":"` + TestEmail + `","user_id":"` + userId.String() + `"}]}`
		manifestJSON := `{"ukama":[{"item":"` + TestItemName + `","description":"` + TestItemDesc + `","inventory":"` + TestInventoryQty + `","opex_fee":"` + TestOpexFeeAmount + `","vat":"` + TestVatAmount + `","effective_date":"` + TestEffectiveDate + `"}]}`

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestBranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(errors.New("database delete failed"))

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

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
		rootJSON := `{"test":[{"company":"` + TestCompany + `","git_branch_name":"` + TestBranchName + `","email":"` + TestEmail + `","user_id":"` + userId.String() + `"}]}`
		manifestJSON := `{"ukama":[{"item":"` + TestItemName + `","description":"` + TestItemDesc + `","inventory":"` + TestInventoryQty + `","opex_fee":"` + TestOpexFeeAmount + `","vat":"` + TestVatAmount + `","effective_date":"` + TestEffectiveDate + `"}]}`

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestBranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(errors.New("database add failed"))

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

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
		rootJSON := `{"test":[{"company":"` + TestCompany + `","git_branch_name":"` + TestBranchName + `","email":"` + TestEmail + `","user_id":"` + userId.String() + `"}]}`
		manifestJSON := `{"ukama":[{"item":"` + TestItemName + `","description":"` + TestItemDesc + `","inventory":"` + TestInventoryQty + `","opex_fee":"` + TestOpexFeeAmount + `","vat":"` + TestVatAmount + `","effective_date":"` + TestEffectiveDate + `"}]}`

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestBranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return(nil, errors.New("database get by user failed"))

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

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
		rootJSON := `{"test":[{"company":"` + TestCompany + `","git_branch_name":"` + TestBranchName + `","email":"` + TestEmail + `","user_id":"` + userId.String() + `"}]}`
		manifestJSON := `{"ukama":[{"item":"` + TestItemName + `","description":"` + TestItemDesc + `","inventory":"` + TestInventoryQty + `","opex_fee":"` + TestOpexFeeAmount + `","vat":"` + TestVatAmount + `","effective_date":"` + TestEffectiveDate + `"}]}`

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestBranchName).Return(nil)
		gitClient.On("ReadFileJSON", "/manifest.json").Return([]byte(manifestJSON), nil)

		accRepo.On("Delete").Return(nil)
		accRepo.On("Add", mock.AnythingOfType("[]*db.Accounting")).Return(nil)
		accRepo.On("GetByUser", userId.String()).Return([]*db.Accounting{
			{
				Id:            uuid.NewV4(),
				UserId:        userId,
				Item:          TestItemName,
				Description:   TestItemDesc,
				Inventory:     TestInventoryQty,
				OpexFee:       TestOpexFeeAmount,
				Vat:           TestVatAmount,
				EffectiveDate: TestEffectiveDate,
			},
		}, nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("publish failed"))

		s := NewAccountingServer(OrgName, accRepo, msgBus, "", gitClient, "")

		resp, err := s.SyncAccounting(context.TODO(), &pb.SyncAcountingRequest{})

		// Note: Message bus failure should not cause the entire operation to fail
		// The function should still return success but log the error
		assert.NoError(t, err)
		assert.NotNil(t, resp)

		gitClient.AssertExpectations(t)
		accRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})
}
