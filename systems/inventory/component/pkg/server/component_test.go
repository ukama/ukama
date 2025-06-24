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
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/inventory/component/mocks"
	pb "github.com/ukama/ukama/systems/inventory/component/pb/gen"
	"github.com/ukama/ukama/systems/inventory/component/pkg/db"
	"gorm.io/gorm"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
)

const OrgName = "testorg"

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
		assert.Contains(t, err.Error(), "invalid format of component uuid")
		compRepo.AssertExpectations(t)
	})
}

func TestComponentServer_GetByUser(t *testing.T) {
	t.Run("Component records by user found", func(t *testing.T) {
		var uId = uuid.NewV4()
		var cId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}

		compRepo.On("GetByUser", uId.String(), int32(1)).Return(
			[]*db.Component{
				{
					Id:            cId,
					UserId:        uId,
					Inventory:     "2",
					Category:      1,
					Type:          "tower node",
					Description:   "best tower node",
					DatasheetURL:  "https://datasheet.com",
					ImagesURL:     "https://images.com",
					PartNumber:    "1234",
					Manufacturer:  "ukama",
					Managed:       "ukama",
					Warranty:      1,
					Specification: "spec",
				}}, nil).Once()

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")

		compResp, err := s.GetByUser(context.TODO(),
			&pb.GetByUserRequest{
				UserId:   uId.String(),
				Category: pb.ComponentCategory_ACCESS})

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
				Category: pb.ComponentCategory_BACKHAUL})

		assert.NoError(t, err)
		assert.NotNil(t, compResp)
		assert.Empty(t, compResp.GetComponents())
		compRepo.AssertExpectations(t)
	})

	t.Run("Database error when getting components by user", func(t *testing.T) {
		var uId = uuid.NewV4()

		compRepo := &mocks.ComponentRepo{}

		compRepo.On("GetByUser", uId.String(), int32(1)).Return(nil, gorm.ErrInvalidDB).Once()

		s := NewComponentServer(OrgName, compRepo, nil, "", nil, "", "", "")

		compResp, err := s.GetByUser(context.TODO(),
			&pb.GetByUserRequest{
				UserId:   uId.String(),
				Category: pb.ComponentCategory_ACCESS})

		assert.Error(t, err)
		assert.Nil(t, compResp)
		assert.Contains(t, err.Error(), "invalid db")
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
					"company": "test-company",
					"git_branch_name": "test-branch",
					"email": "test@example.com",
					"user_id": "` + userId.String() + `"
				}
			]
		}`
		gitClient.On("ReadFileJSON", "/root.json").Return([]byte(rootJSON), nil)

		gitClient.On("BranchCheckout", "test-branch").Return(nil)

		componentYAML := `category: "ACCESS"
type: "test-component"
description: "Test component description"
imagesURL: "https://example.com/image.jpg"
datasheetURL: "https://example.com/datasheet.pdf"
inventoryID: "INV001"
partNumber: "PN001"
manufacturer: "Test Manufacturer"
managed: "ukama"
warranty: 12
specification: "Test specification"`
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
		assert.Contains(t, err.Error(), "failed to clone git repo")

		compRepo.AssertExpectations(t)
		gitClient.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})
}
