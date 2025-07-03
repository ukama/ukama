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
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/component/mocks"
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
	"github.com/ukama/ukama/systems/inventory/component/pkg/db"
	"gorm.io/gorm"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
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
	TestWarranty2      = 12
	TestSpecification2 = "Test specification"

	// Error messages
	ErrInvalidUUID    = "invalid format of component uuid"
	ErrInvalidDB      = "invalid db"
	ErrGitCloneFailed = "failed to clone git repo"
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
}
