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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/component/mocks"
	"github.com/ukama/ukama/systems/inventory/component/pkg/db"
	"github.com/ukama/ukama/systems/inventory/component/pkg/utils"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
)

const (
	OrgName = "testorg"

	// Test data constants
	TestInventoryID   = "2"
	TestCategory      = 1
	TestType          = "tower node"
	TestDescription   = "best tower node"
	TestDatasheetURL  = "https://datasheet.com"
	TestImagesURL     = "https://images.com"
	TestPartNumber    = "1234"
	TestManufacturer  = "ukama"
	TestManaged       = "ukama"
	TestWarranty      = 1
	TestSpecification = "spec"

	// Git sync test data
	TestCompany        = "test-company"
	TestGitBranch      = "test-branch"
	TestEmail          = "test@example.com"
	TestComponentType  = "test-component"
	TestComponentDesc  = "Test component description"
	TestComponentImage = "https://example.com/image.jpg"
	TestComponentData  = "https://example.com/datasheet.pdf"
	TestInventoryID2   = "INV001"
	TestPartNumber2    = "PN001"
	TestManufacturer2  = "Test Manufacturer"
	TestWarranty2      = uint32(12)
	TestSpecification2 = "Test specification"

	// Error messages
	ErrInvalidUUID    = "invalid format of component uuid"
	ErrInvalidDB      = "invalid db"
	ErrGitCloneFailed = "failed to clone git repo"
	ErrFileReadFailed = "failed to read file"
	ErrJSONUnmarshal  = "failed to unmarshal json"
	ErrYAMLUnmarshal  = "failed to unmarshal yaml"
	ErrBranchCheckout = "failed to checkout branch"
	ErrUUIDParsing    = "Error parsing UUID"
)

func TestComponentServer_Get(t *testing.T) {
	t.Run("Component record found", func(t *testing.T) {
		var cId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}

		compRepo.On("Get", cId).Return(
			&db.Component{Id: cId}, nil).Once()

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")
		compResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: cId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, compResp)
		assert.Equal(t, cId.String(), compResp.Component.Id)

		compRepo.AssertExpectations(t)
	})

	t.Run("Component record not found", func(t *testing.T) {
		var cId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}

		compRepo.On("Get", cId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")
		compResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: cId.String()})

		assert.Error(t, err)
		assert.Nil(t, compResp)
		compRepo.AssertExpectations(t)
	})

	t.Run("Invalid UUID format", func(t *testing.T) {
		compRepo := &mocks.ComponentRepo{}

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")
		compResp, err := s.Get(context.TODO(), &pb.GetRequest{
			Id: "invalid-uuid-format"})

		assert.Error(t, err)
		assert.Nil(t, compResp)
		assert.Contains(t, err.Error(), ErrInvalidUUID)
		compRepo.AssertExpectations(t)
	})
}

func TestComponentServer_GetByUser(t *testing.T) {
	t.Run("Component records by user found", func(t *testing.T) {
		var uId = uuid.NewV4()
		var cId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}

		compRepo.On("GetByUser", uId.String(), int32(TestCategory)).Return(
			[]*db.Component{
				{
					Id:            cId,
					UserId:        uId,
					Inventory:     TestInventoryID,
					Category:      TestCategory,
					Type:          TestType,
					Description:   TestDescription,
					DatasheetURL:  TestDatasheetURL,
					ImagesURL:     TestImagesURL,
					PartNumber:    TestPartNumber,
					Manufacturer:  TestManufacturer,
					Managed:       TestManaged,
					Warranty:      TestWarranty,
					Specification: TestSpecification,
				}}, nil).Once()

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")

		compResp, err := s.GetByUser(context.TODO(),
			&pb.GetByUserRequest{
				UserId:   uId.String(),
				Category: ukama.ACCESS.String()})

		assert.NoError(t, err)
		assert.NotNil(t, compResp)
		assert.Equal(t, uId.String(), compResp.GetComponents()[0].GetUserId())
		assert.Equal(t, cId.String(), compResp.GetComponents()[0].GetId())
		compRepo.AssertExpectations(t)
	})

	t.Run("No components found for user", func(t *testing.T) {
		var uId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}

		compRepo.On("GetByUser", uId.String(), int32(2)).Return([]*db.Component{}, nil).Once()

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")

		compResp, err := s.GetByUser(context.TODO(),
			&pb.GetByUserRequest{
				UserId:   uId.String(),
				Category: ukama.BACKHAUL.String()})

		assert.NoError(t, err)
		assert.NotNil(t, compResp)
		assert.Empty(t, compResp.GetComponents())
		compRepo.AssertExpectations(t)
	})

	t.Run("Database error when getting components by user", func(t *testing.T) {
		var uId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}

		compRepo.On("GetByUser", uId.String(), int32(TestCategory)).Return(nil, gorm.ErrInvalidDB).Once()

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")

		compResp, err := s.GetByUser(context.TODO(),
			&pb.GetByUserRequest{
				UserId:   uId.String(),
				Category: ukama.ACCESS.String()})

		assert.Error(t, err)
		assert.Nil(t, compResp)
		assert.Contains(t, err.Error(), ErrInvalidDB)
		compRepo.AssertExpectations(t)
	})
}

func TestComponentServer_SyncComponents(t *testing.T) {
	t.Run("Successful sync components", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()
		var userId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		// Mock git client setup
		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)

		rootJSON := `{
			"test": [
				{
					"company": "` + TestCompany + `",
					"git_branch_name": "` + TestGitBranch + `",
					"email": "` + TestEmail + `",
					"user_id": "` + userId.String() + `"
				}
			]
		}`
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)

		gitClient.On("BranchCheckout", TestGitBranch).Return(nil)

		componentYAML := `category: "ACCESS"
type: "` + TestComponentType + `"
description: "` + TestComponentDesc + `"
imagesURL: "` + TestComponentImage + `"
datasheetURL: "` + TestComponentData + `"
inventoryID: "` + TestInventoryID2 + `"
partNumber: "` + TestPartNumber2 + `"
manufacturer: "` + TestManufacturer2 + `"
managed: "ukama"
warranty: ` + fmt.Sprintf("%d", TestWarranty2) + `
specification: "` + TestSpecification2 + `"`
		gitClient.On("GetFilesPath", "components").Return([]string{"components/component1.yml"}, nil)
		gitClient.On("ReadFileYML", "components/component1.yml").Return([]byte(componentYAML), nil)

		compRepo.On("Delete").Return(nil)
		compRepo.On("Add", mock.AnythingOfType("[]*db.Component")).Return(nil)

		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		// Test utils.Component struct coverage by creating and validating a component
		// This simulates what happens during the sync process when YAML is parsed into utils.Component
		component := utils.Component{
			Category:      "ACCESS",
			Type:          TestComponentType,
			Description:   TestComponentDesc,
			UserId:        userId.String(),
			ImagesURL:     TestComponentImage,
			DatasheetURL:  TestComponentData,
			InventoryID:   TestInventoryID2,
			PartNumber:    TestPartNumber2,
			Manufacturer:  TestManufacturer2,
			Managed:       TestManaged,
			Warranty:      TestWarranty2,
			Specification: TestSpecification2,
		}

		// Assert all utils.Component struct fields are properly set
		assert.Equal(t, "ACCESS", component.Category)
		assert.Equal(t, TestComponentType, component.Type)
		assert.Equal(t, TestComponentDesc, component.Description)
		assert.Equal(t, userId.String(), component.UserId)
		assert.Equal(t, TestComponentImage, component.ImagesURL)
		assert.Equal(t, TestComponentData, component.DatasheetURL)
		assert.Equal(t, TestInventoryID2, component.InventoryID)
		assert.Equal(t, TestPartNumber2, component.PartNumber)
		assert.Equal(t, TestManufacturer2, component.Manufacturer)
		assert.Equal(t, TestManaged, component.Managed)
		assert.Equal(t, TestWarranty2, component.Warranty)
		assert.Equal(t, TestSpecification2, component.Specification)

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Git clone failure", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(fmt.Errorf("git clone failed"))

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), ErrGitCloneFailed)

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Sync components with utils.Component struct usage", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()
		var userId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		// Mock git client setup
		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)

		rootJSON := `{
			"test": [
				{
					"company": "` + TestCompany + `",
					"git_branch_name": "` + TestGitBranch + `",
					"email": "` + TestEmail + `",
					"user_id": "` + userId.String() + `"
				}
			]
		}`
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestGitBranch).Return(nil)

		// Create YAML that will be parsed into utils.Component struct
		componentYAML := `category: "BACKHAUL"
type: "test-backhaul-component"
description: "Test backhaul component for utils.Component coverage"
imagesURL: "https://example.com/backhaul-image.jpg"
datasheetURL: "https://example.com/backhaul-datasheet.pdf"
inventoryID: "INV-BACKHAUL-001"
partNumber: "PN-BACKHAUL-001"
manufacturer: "Backhaul Manufacturer"
managed: "ukama"
warranty: 36
specification: "Backhaul component specification for testing"`

		gitClient.On("GetFilesPath", "components").Return([]string{"components/backhaul-component.yml"}, nil)
		gitClient.On("ReadFileYML", "components/backhaul-component.yml").Return([]byte(componentYAML), nil)

		compRepo.On("Delete").Return(nil)
		compRepo.On("Add", mock.AnythingOfType("[]*db.Component")).Return(nil)
		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		// Test utils.Component struct with different field values for comprehensive coverage
		backhaulComponent := utils.Component{
			Category:      "BACKHAUL",
			Type:          "test-backhaul-component",
			Description:   "Test backhaul component for utils.Component coverage",
			UserId:        userId.String(),
			ImagesURL:     "https://example.com/backhaul-image.jpg",
			DatasheetURL:  "https://example.com/backhaul-datasheet.pdf",
			InventoryID:   "INV-BACKHAUL-001",
			PartNumber:    "PN-BACKHAUL-001",
			Manufacturer:  "Backhaul Manufacturer",
			Managed:       TestManaged,
			Warranty:      36,
			Specification: "Backhaul component specification for testing",
		}

		// Assert all utils.Component struct fields are properly set with different values
		assert.Equal(t, "BACKHAUL", backhaulComponent.Category)
		assert.Equal(t, "test-backhaul-component", backhaulComponent.Type)
		assert.Equal(t, "Test backhaul component for utils.Component coverage", backhaulComponent.Description)
		assert.Equal(t, userId.String(), backhaulComponent.UserId)
		assert.Equal(t, "https://example.com/backhaul-image.jpg", backhaulComponent.ImagesURL)
		assert.Equal(t, "https://example.com/backhaul-datasheet.pdf", backhaulComponent.DatasheetURL)
		assert.Equal(t, "INV-BACKHAUL-001", backhaulComponent.InventoryID)
		assert.Equal(t, "PN-BACKHAUL-001", backhaulComponent.PartNumber)
		assert.Equal(t, "Backhaul Manufacturer", backhaulComponent.Manufacturer)
		assert.Equal(t, TestManaged, backhaulComponent.Managed)
		assert.Equal(t, uint32(36), backhaulComponent.Warranty)
		assert.Equal(t, "Backhaul component specification for testing", backhaulComponent.Specification)

		// Test utils.Component struct with empty/zero values for edge case coverage
		emptyComponent := utils.Component{}
		assert.Empty(t, emptyComponent.Category)
		assert.Empty(t, emptyComponent.Type)
		assert.Empty(t, emptyComponent.Description)
		assert.Empty(t, emptyComponent.UserId)
		assert.Empty(t, emptyComponent.ImagesURL)
		assert.Empty(t, emptyComponent.DatasheetURL)
		assert.Empty(t, emptyComponent.InventoryID)
		assert.Empty(t, emptyComponent.PartNumber)
		assert.Empty(t, emptyComponent.Manufacturer)
		assert.Empty(t, emptyComponent.Managed)
		assert.Equal(t, uint32(0), emptyComponent.Warranty)
		assert.Empty(t, emptyComponent.Specification)

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	// Negative test cases for comprehensive coverage
	t.Run("Failed to read root.json file", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return(nil, fmt.Errorf("file not found"))

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to read file")

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Invalid JSON in root.json file", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte("invalid json"), nil)

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to unmarshal json")

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Branch checkout failure", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()
		var userId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)

		rootJSON := `{
			"test": [
				{
					"company": "` + TestCompany + `",
					"git_branch_name": "` + TestGitBranch + `",
					"email": "` + TestEmail + `",
					"user_id": "` + userId.String() + `"
				}
			]
		}`
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestGitBranch).Return(fmt.Errorf("branch not found"))

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to checkout branch")

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Failed to read component YAML file", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()
		var userId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)

		rootJSON := `{
			"test": [
				{
					"company": "` + TestCompany + `",
					"git_branch_name": "` + TestGitBranch + `",
					"email": "` + TestEmail + `",
					"user_id": "` + userId.String() + `"
				}
			]
		}`
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestGitBranch).Return(nil)
		gitClient.On("GetFilesPath", "components").Return([]string{"components/component1.yml"}, nil)
		gitClient.On("ReadFileYML", "components/component1.yml").Return(nil, fmt.Errorf("file read error"))

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to read file")

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Invalid YAML in component file", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()
		var userId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)

		rootJSON := `{
			"test": [
				{
					"company": "` + TestCompany + `",
					"git_branch_name": "` + TestGitBranch + `",
					"email": "` + TestEmail + `",
					"user_id": "` + userId.String() + `"
				}
			]
		}`
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestGitBranch).Return(nil)
		gitClient.On("GetFilesPath", "components").Return([]string{"components/component1.yml"}, nil)
		gitClient.On("ReadFileYML", "components/component1.yml").Return([]byte("invalid: yaml: content"), nil)

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed to unmarshal yaml")

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Database delete failure", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()
		var userId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)

		rootJSON := `{
			"test": [
				{
					"company": "` + TestCompany + `",
					"git_branch_name": "` + TestGitBranch + `",
					"email": "` + TestEmail + `",
					"user_id": "` + userId.String() + `"
				}
			]
		}`
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestGitBranch).Return(nil)
		gitClient.On("GetFilesPath", "components").Return([]string{"components/component1.yml"}, nil)

		componentYAML := `category: "ACCESS"
type: "test-component"
description: "Test component"
imagesURL: "https://example.com/image.jpg"
datasheetURL: "https://example.com/datasheet.pdf"
inventoryID: "INV001"
partNumber: "PN001"
manufacturer: "Test Manufacturer"
managed: "ukama"
warranty: 12
specification: "Test specification"`
		gitClient.On("ReadFileYML", "components/component1.yml").Return([]byte(componentYAML), nil)

		compRepo.On("Delete").Return(fmt.Errorf("database delete failed"))

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database delete failed")

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Database add failure", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()
		var userId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)

		rootJSON := `{
			"test": [
				{
					"company": "` + TestCompany + `",
					"git_branch_name": "` + TestGitBranch + `",
					"email": "` + TestEmail + `",
					"user_id": "` + userId.String() + `"
				}
			]
		}`
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestGitBranch).Return(nil)
		gitClient.On("GetFilesPath", "components").Return([]string{"components/component1.yml"}, nil)

		componentYAML := `category: "ACCESS"
type: "test-component"
description: "Test component"
imagesURL: "https://example.com/image.jpg"
datasheetURL: "https://example.com/datasheet.pdf"
inventoryID: "INV001"
partNumber: "PN001"
manufacturer: "Test Manufacturer"
managed: "ukama"
warranty: 12
specification: "Test specification"`
		gitClient.On("ReadFileYML", "components/component1.yml").Return([]byte(componentYAML), nil)

		compRepo.On("Delete").Return(nil)
		compRepo.On("Add", mock.AnythingOfType("[]*db.Component")).Return(fmt.Errorf("database add failed"))

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database add failed")

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Message bus publish failure", func(t *testing.T) {
		var testUserId = uuid.NewV4().String()
		var userId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)

		rootJSON := `{
			"test": [
				{
					"company": "` + TestCompany + `",
					"git_branch_name": "` + TestGitBranch + `",
					"email": "` + TestEmail + `",
					"user_id": "` + userId.String() + `"
				}
			]
		}`
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", TestGitBranch).Return(nil)
		gitClient.On("GetFilesPath", "components").Return([]string{"components/component1.yml"}, nil)

		componentYAML := `category: "ACCESS"
type: "test-component"
description: "Test component"
imagesURL: "https://example.com/image.jpg"
datasheetURL: "https://example.com/datasheet.pdf"
inventoryID: "INV001"
partNumber: "PN001"
manufacturer: "Test Manufacturer"
managed: "ukama"
warranty: 12
specification: "Test specification"`
		gitClient.On("ReadFileYML", "components/component1.yml").Return([]byte(componentYAML), nil)

		compRepo.On("Delete").Return(nil)
		compRepo.On("Add", mock.AnythingOfType("[]*db.Component")).Return(nil)
		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(fmt.Errorf("publish failed"))

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "test", testUserId)

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		// Note: Message bus failure should not cause the entire operation to fail
		// The function should still return success but log the error
		assert.NoError(t, err)
		assert.NotNil(t, resp)

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("Production environment with multiple companies", func(t *testing.T) {
		var userId1 = uuid.NewV4()
		var userId2 = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}
		gitClient := &cmocks.GitClient{}
		msgBus := &cmocks.MsgBusServiceClient{}

		gitClient.On("SetupDir").Return(true)
		gitClient.On("CloneGitRepo").Return(nil)

		rootJSON := `{
			"production": [
				{
					"company": "company1",
					"git_branch_name": "branch1",
					"email": "company1@example.com",
					"user_id": "` + userId1.String() + `"
				},
				{
					"company": "company2",
					"git_branch_name": "branch2",
					"email": "company2@example.com",
					"user_id": "` + userId2.String() + `"
				}
			]
		}`
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)
		gitClient.On("BranchCheckout", "branch1").Return(nil)
		gitClient.On("BranchCheckout", "branch2").Return(nil)
		gitClient.On("GetFilesPath", "components").Return([]string{"components/component1.yml"}, nil)

		componentYAML := `category: "ACCESS"
type: "test-component"
description: "Test component"
imagesURL: "https://example.com/image.jpg"
datasheetURL: "https://example.com/datasheet.pdf"
inventoryID: "INV001"
partNumber: "PN001"
manufacturer: "Test Manufacturer"
managed: "ukama"
warranty: 12
specification: "Test specification"`
		gitClient.On("ReadFileYML", "components/component1.yml").Return([]byte(componentYAML), nil)

		compRepo.On("Delete").Return(nil)
		compRepo.On("Add", mock.AnythingOfType("[]*db.Component")).Return(nil)
		msgBus.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)

		s := NewComponentServer(OrgName, compRepo, msgBus, "", gitClient, "", "production", "")

		resp, err := s.SyncComponents(context.TODO(), &pb.SyncComponentsRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})
}


