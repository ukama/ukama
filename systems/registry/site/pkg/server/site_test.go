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
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/rest/client/inventory"
	"github.com/ukama/ukama/systems/common/uuid"
	npb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	netmocks "github.com/ukama/ukama/systems/registry/network/pb/gen/mocks"
	"github.com/ukama/ukama/systems/registry/site/mocks"
	pb "github.com/ukama/ukama/systems/registry/site/pb/gen"
	"github.com/ukama/ukama/systems/registry/site/pkg/db"
	"gorm.io/gorm"
)

const OrgName = "Ukama"

func TestSiteService_Get(t *testing.T) {
	siteRepo := &mocks.SiteRepo{}
	msgclientRepo := &cmocks.MsgBusServiceClient{}
	netRepo := &mocks.NetworkClientProvider{}

	s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", nil)

	t.Run("SiteFound", func(t *testing.T) {
		siteId := uuid.NewV4()
		networkId := uuid.NewV4()
		backhaulId := uuid.NewV4()
		powerId := uuid.NewV4()
		accessId := uuid.NewV4()
		switchId := uuid.NewV4()
		spectrumId := uuid.NewV4()

		mockSite := &db.Site{
			Id:            siteId,
			Name:          "Test Site",
			Location:      "Test Location",
			NetworkId:     networkId,
			BackhaulId:    backhaulId,
			PowerId:       powerId,
			AccessId:      accessId,
			SwitchId:      switchId,
			SpectrumId:    spectrumId,
			IsDeactivated: false,
			Latitude:      40.7128,
			Longitude:     -74.0060,
			InstallDate:   "2023-12-01T00:00:00Z",
		}

		siteRepo.On("Get", siteId).Return(mockSite, nil).Once()

		resp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: siteId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)

		// Verify all fields are correctly mapped
		assert.Equal(t, siteId.String(), resp.Site.Id)
		assert.Equal(t, mockSite.Name, resp.Site.Name)
		assert.Equal(t, mockSite.Location, resp.Site.Location)
		assert.Equal(t, networkId.String(), resp.Site.NetworkId)
		assert.Equal(t, backhaulId.String(), resp.Site.BackhaulId)
		assert.Equal(t, powerId.String(), resp.Site.PowerId)
		assert.Equal(t, accessId.String(), resp.Site.AccessId)
		assert.Equal(t, switchId.String(), resp.Site.SwitchId)
		assert.Equal(t, spectrumId.String(), resp.Site.SpectrumId)
		assert.Equal(t, mockSite.IsDeactivated, resp.Site.IsDeactivated)
		assert.Equal(t, mockSite.Latitude, resp.Site.Latitude)
		assert.Equal(t, mockSite.Longitude, resp.Site.Longitude)
		assert.Equal(t, mockSite.InstallDate, resp.Site.InstallDate)

		siteRepo.AssertExpectations(t)
	})

	t.Run("SiteNotFound", func(t *testing.T) {
		siteId := uuid.NewV4()

		siteRepo.On("Get", siteId).Return(nil, gorm.ErrRecordNotFound).Once()

		resp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: siteId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "record not found")
		siteRepo.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		invalidUUID := "invalid-uuid-format"

		resp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: invalidUUID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("EmptySiteId", func(t *testing.T) {
		resp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: ""})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("DatabaseError", func(t *testing.T) {
		siteId := uuid.NewV4()
		dbError := gorm.ErrInvalidDB

		siteRepo.On("Get", siteId).Return(nil, dbError).Once()

		resp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: siteId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid db")
		siteRepo.AssertExpectations(t)
	})
}

func TestSiteService_List(t *testing.T) {
	siteRepo := &mocks.SiteRepo{}
	msgclientRepo := &cmocks.MsgBusServiceClient{}
	netRepo := &mocks.NetworkClientProvider{}

	s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", nil)

	t.Run("ValidRequest", func(t *testing.T) {
		netId := uuid.NewV4()

		mockSites := []*db.Site{
			{
				Id:        uuid.NewV4(),
				NetworkId: netId,
				Name:      "Site1",
			},
			{
				Id:        uuid.NewV4(),
				NetworkId: netId,
				Name:      "Site2",
			},
		}

		var mockSitesConverted []db.Site
		for _, site := range mockSites {
			mockSitesConverted = append(mockSitesConverted, *site)
		}

		siteRepo.On("List", &netId, false).Return(mockSitesConverted, nil)

		req := &pb.ListRequest{
			NetworkId:     netId.String(),
			IsDeactivated: false,
		}

		resp, err := s.List(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, len(mockSites), len(resp.Sites))
		for i, site := range mockSites {
			assert.Equal(t, site.Id.String(), resp.Sites[i].Id)
			assert.Equal(t, site.Name, resp.Sites[i].Name)
		}
		siteRepo.AssertExpectations(t)
	})

	t.Run("EmptyNetworkId", func(t *testing.T) {
		mockSites := []db.Site{
			{Id: uuid.NewV4(), Name: "Site1"},
		}
		siteRepo.On("List", (*uuid.UUID)(nil), false).Return(mockSites, nil)

		req := &pb.ListRequest{NetworkId: "", IsDeactivated: false}
		resp, err := s.List(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, len(resp.Sites))
		assert.Equal(t, mockSites[0].Id.String(), resp.Sites[0].Id)
		siteRepo.AssertExpectations(t)
	})

	t.Run("InvalidNetworkId", func(t *testing.T) {
		req := &pb.ListRequest{NetworkId: "invalid-uuid", IsDeactivated: false}
		resp, err := s.List(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid network ID")
	})

	t.Run("IsDeactivatedTrue", func(t *testing.T) {
		netId := uuid.NewV4()
		mockSites := []db.Site{
			{Id: uuid.NewV4(), NetworkId: netId, Name: "DeactivatedSite", IsDeactivated: true},
		}
		siteRepo.On("List", &netId, true).Return(mockSites, nil)

		req := &pb.ListRequest{NetworkId: netId.String(), IsDeactivated: true}
		resp, err := s.List(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, len(resp.Sites))
		assert.True(t, resp.Sites[0].IsDeactivated)
		siteRepo.AssertExpectations(t)
	})

	t.Run("RepoError", func(t *testing.T) {
		netId := uuid.NewV4()
		siteRepo.On("List", &netId, false).Return(nil, gorm.ErrInvalidDB)

		req := &pb.ListRequest{NetworkId: netId.String(), IsDeactivated: false}
		resp, err := s.List(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid db")
		siteRepo.AssertExpectations(t)
	})

	t.Run("EmptyResult", func(t *testing.T) {
		netId := uuid.NewV4()
		siteRepo.On("List", &netId, false).Return([]db.Site{}, nil)

		req := &pb.ListRequest{NetworkId: netId.String(), IsDeactivated: false}
		resp, err := s.List(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 0, len(resp.Sites))
		siteRepo.AssertExpectations(t)
	})
}

func TestSiteService_Update(t *testing.T) {
	siteRepo := &mocks.SiteRepo{}
	msgclientRepo := &cmocks.MsgBusServiceClient{}
	netRepo := &mocks.NetworkClientProvider{}

	s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", nil)

	t.Run("Success", func(t *testing.T) {
		siteId := uuid.NewV4()
		newName := "Updated Site Name"

		// Mock the site repository update
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(nil).Once()

		// Mock message bus publish
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		req := &pb.UpdateRequest{
			SiteId: siteId.String(),
			Name:   newName,
		}

		resp, err := s.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, siteId.String(), resp.Site.Id)
		assert.Equal(t, newName, resp.Site.Name)

		// Verify the mock was called with correct parameters
		siteRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("InvalidSiteId", func(t *testing.T) {
		req := &pb.UpdateRequest{
			SiteId: "invalid-uuid",
			Name:   "Test Site",
		}

		resp, err := s.Update(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("EmptySiteId", func(t *testing.T) {
		req := &pb.UpdateRequest{
			SiteId: "",
			Name:   "Test Site",
		}

		resp, err := s.Update(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("EmptyName", func(t *testing.T) {
		siteId := uuid.NewV4()

		// Mock the site repository update
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(nil).Once()

		// Mock message bus publish
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		req := &pb.UpdateRequest{
			SiteId: siteId.String(),
			Name:   "",
		}

		resp, err := s.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, siteId.String(), resp.Site.Id)
		assert.Equal(t, "", resp.Site.Name)

		siteRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		siteId := uuid.NewV4()

		// Mock the site repository update to return an error
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(gorm.ErrRecordNotFound).Once()

		req := &pb.UpdateRequest{
			SiteId: siteId.String(),
			Name:   "Test Site",
		}

		resp, err := s.Update(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "record not found")

		siteRepo.AssertExpectations(t)
	})

	t.Run("MessageBusError", func(t *testing.T) {
		siteId := uuid.NewV4()
		newName := "Updated Site Name"

		// Mock the site repository update
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(nil).Once()

		// Mock message bus publish to return an error
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(gorm.ErrInvalidDB).Once()

		req := &pb.UpdateRequest{
			SiteId: siteId.String(),
			Name:   newName,
		}

		resp, err := s.Update(context.Background(), req)

		// The update should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, siteId.String(), resp.Site.Id)
		assert.Equal(t, newName, resp.Site.Name)

		siteRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("NilMessageBus", func(t *testing.T) {
		// Create server without message bus
		sNoMsgBus := NewSiteServer(OrgName, siteRepo, nil, netRepo, "", nil)
		siteId := uuid.NewV4()
		newName := "Updated Site Name"

		// Mock the site repository update
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(nil).Once()

		req := &pb.UpdateRequest{
			SiteId: siteId.String(),
			Name:   newName,
		}

		resp, err := sNoMsgBus.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, siteId.String(), resp.Site.Id)
		assert.Equal(t, newName, resp.Site.Name)

		siteRepo.AssertExpectations(t)
	})

	t.Run("LongName", func(t *testing.T) {
		siteId := uuid.NewV4()
		longName := "This is a very long site name that might exceed some limits but should still be processed by the update function"

		// Mock the site repository update
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(nil).Once()

		// Mock message bus publish
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		req := &pb.UpdateRequest{
			SiteId: siteId.String(),
			Name:   longName,
		}

		resp, err := s.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, siteId.String(), resp.Site.Id)
		assert.Equal(t, longName, resp.Site.Name)

		siteRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})
}

func TestSiteService_Add(t *testing.T) {
	siteRepo := &mocks.SiteRepo{}
	msgclientRepo := &cmocks.MsgBusServiceClient{}
	netRepo := &mocks.NetworkClientProvider{}
	inventoryClient := &cmocks.ComponentClient{}

	s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", inventoryClient)

	// Create valid UUIDs for testing
	validNetworkId := uuid.NewV4()
	validBackhaulId := uuid.NewV4()
	validPowerId := uuid.NewV4()
	validAccessId := uuid.NewV4()
	validSwitchId := uuid.NewV4()
	validSpectrumId := uuid.NewV4()
	validInstallDate := "2023-12-01T00:00:00Z"

	validRequest := &pb.AddRequest{
		Name:          "Test Site",
		NetworkId:     validNetworkId.String(),
		BackhaulId:    validBackhaulId.String(),
		PowerId:       validPowerId.String(),
		AccessId:      validAccessId.String(),
		SwitchId:      validSwitchId.String(),
		SpectrumId:    validSpectrumId.String(),
		IsDeactivated: false,
		Latitude:      40.7128,
		Longitude:     -74.0060,
		Location:      "New York",
		InstallDate:   validInstallDate,
	}

	t.Run("Success", func(t *testing.T) {
		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: validNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			validBackhaulId.String(),
			validPowerId.String(),
			validAccessId.String(),
			validSwitchId.String(),
			validSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", validNetworkId).Return(int64(1), nil)

		// Mock message bus
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		resp, err := s.Add(context.Background(), validRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, validRequest.Name, resp.Site.Name)
		assert.Equal(t, validRequest.NetworkId, resp.Site.NetworkId)
		assert.Equal(t, validRequest.BackhaulId, resp.Site.BackhaulId)
		assert.Equal(t, validRequest.PowerId, resp.Site.PowerId)
		assert.Equal(t, validRequest.AccessId, resp.Site.AccessId)
		assert.Equal(t, validRequest.SwitchId, resp.Site.SwitchId)
		assert.Equal(t, validRequest.SpectrumId, resp.Site.SpectrumId)
		assert.Equal(t, validRequest.IsDeactivated, resp.Site.IsDeactivated)
		assert.Equal(t, validRequest.Latitude, resp.Site.Latitude)
		assert.Equal(t, validRequest.Longitude, resp.Site.Longitude)
		assert.Equal(t, validRequest.Location, resp.Site.Location)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("InvalidNetworkId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       "Test Site",
			NetworkId:  "invalid-uuid",
			BackhaulId: validBackhaulId.String(),
			PowerId:    validPowerId.String(),
			AccessId:   validAccessId.String(),
			SwitchId:   validSwitchId.String(),
			SpectrumId: validSpectrumId.String(),
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidBackhaulId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       "Test Site",
			NetworkId:  validNetworkId.String(),
			BackhaulId: "invalid-uuid",
			PowerId:    validPowerId.String(),
			AccessId:   validAccessId.String(),
			SwitchId:   validSwitchId.String(),
			SpectrumId: validSpectrumId.String(),
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidPowerId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       "Test Site",
			NetworkId:  validNetworkId.String(),
			BackhaulId: validBackhaulId.String(),
			PowerId:    "invalid-uuid",
			AccessId:   validAccessId.String(),
			SwitchId:   validSwitchId.String(),
			SpectrumId: validSpectrumId.String(),
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidAccessId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       "Test Site",
			NetworkId:  validNetworkId.String(),
			BackhaulId: validBackhaulId.String(),
			PowerId:    validPowerId.String(),
			AccessId:   "invalid-uuid",
			SwitchId:   validSwitchId.String(),
			SpectrumId: validSpectrumId.String(),
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidSwitchId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       "Test Site",
			NetworkId:  validNetworkId.String(),
			BackhaulId: validBackhaulId.String(),
			PowerId:    validPowerId.String(),
			AccessId:   validAccessId.String(),
			SwitchId:   "invalid-uuid",
			SpectrumId: validSpectrumId.String(),
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidSpectrumId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       "Test Site",
			NetworkId:  validNetworkId.String(),
			BackhaulId: validBackhaulId.String(),
			PowerId:    validPowerId.String(),
			AccessId:   validAccessId.String(),
			SwitchId:   validSwitchId.String(),
			SpectrumId: "invalid-uuid",
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidInstallDate", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:        "Test Site",
			NetworkId:   validNetworkId.String(),
			BackhaulId:  validBackhaulId.String(),
			PowerId:     validPowerId.String(),
			AccessId:    validAccessId.String(),
			SwitchId:    validSwitchId.String(),
			SpectrumId:  validSpectrumId.String(),
			InstallDate: "invalid-date",
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid date format")
	})

	t.Run("NetworkServiceError", func(t *testing.T) {
		// Create fresh mocks for this test
		siteRepo := &mocks.SiteRepo{}
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		netRepo := &mocks.NetworkClientProvider{}
		inventoryClient := &cmocks.ComponentClient{}
		s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", inventoryClient)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			validBackhaulId.String(),
			validPowerId.String(),
			validAccessId.String(),
			validSwitchId.String(),
			validSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock network client to return error
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: validNetworkId.String(),
		}).Return(nil, gorm.ErrRecordNotFound)

		resp, err := s.Add(context.Background(), validRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "record not found")
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
	})

	t.Run("NetworkServiceClientError", func(t *testing.T) {
		// Create fresh mocks for this test
		siteRepo := &mocks.SiteRepo{}
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		netRepo := &mocks.NetworkClientProvider{}
		inventoryClient := &cmocks.ComponentClient{}
		s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", inventoryClient)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			validBackhaulId.String(),
			validPowerId.String(),
			validAccessId.String(),
			validSwitchId.String(),
			validSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock network service provider to return error
		netRepo.On("GetClient").Return(nil, gorm.ErrInvalidDB)

		resp, err := s.Add(context.Background(), validRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid db")
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
	})

	t.Run("SiteRepositoryError", func(t *testing.T) {
		// Create fresh mocks for this test
		siteRepo := &mocks.SiteRepo{}
		msgclientRepo := &cmocks.MsgBusServiceClient{}
		netRepo := &mocks.NetworkClientProvider{}
		inventoryClient := &cmocks.ComponentClient{}
		s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", inventoryClient)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			validBackhaulId.String(),
			validPowerId.String(),
			validAccessId.String(),
			validSwitchId.String(),
			validSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: validNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock site repository Add to return error
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(gorm.ErrInvalidDB)

		resp, err := s.Add(context.Background(), validRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid db")
		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
	})

	t.Run("MessageBusError", func(t *testing.T) {
		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: validNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			validBackhaulId.String(),
			validPowerId.String(),
			validAccessId.String(),
			validSwitchId.String(),
			validSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", validNetworkId).Return(int64(1), nil)

		// Mock message bus to return error
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(gorm.ErrInvalidDB)

		resp, err := s.Add(context.Background(), validRequest)

		// Should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, validRequest.Name, resp.Site.Name)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("GetSiteCountError", func(t *testing.T) {
		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: validNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			validBackhaulId.String(),
			validPowerId.String(),
			validAccessId.String(),
			validSwitchId.String(),
			validSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount to return error
		siteRepo.On("GetSiteCount", validNetworkId).Return(int64(0), gorm.ErrInvalidDB)

		// Mock message bus
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		resp, err := s.Add(context.Background(), validRequest)

		// Should still succeed even if GetSiteCount fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, validRequest.Name, resp.Site.Name)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("NilMessageBus", func(t *testing.T) {
		// Create server without message bus
		sNoMsgBus := NewSiteServer(OrgName, siteRepo, nil, netRepo, "", inventoryClient)

		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: validNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			validBackhaulId.String(),
			validPowerId.String(),
			validAccessId.String(),
			validSwitchId.String(),
			validSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", validNetworkId).Return(int64(1), nil)

		resp, err := sNoMsgBus.Add(context.Background(), validRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, validRequest.Name, resp.Site.Name)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
	})

	t.Run("EmptyName", func(t *testing.T) {
		requestWithEmptyName := &pb.AddRequest{
			Name:          "",
			NetworkId:     validNetworkId.String(),
			BackhaulId:    validBackhaulId.String(),
			PowerId:       validPowerId.String(),
			AccessId:      validAccessId.String(),
			SwitchId:      validSwitchId.String(),
			SpectrumId:    validSpectrumId.String(),
			IsDeactivated: false,
			Latitude:      40.7128,
			Longitude:     -74.0060,
			Location:      "New York",
			InstallDate:   validInstallDate,
		}

		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: validNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			validBackhaulId.String(),
			validPowerId.String(),
			validAccessId.String(),
			validSwitchId.String(),
			validSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", validNetworkId).Return(int64(1), nil)

		// Mock message bus
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		resp, err := s.Add(context.Background(), requestWithEmptyName)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, "", resp.Site.Name)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("ExtremeCoordinates", func(t *testing.T) {
		requestWithExtremeCoords := &pb.AddRequest{
			Name:          "Test Site",
			NetworkId:     validNetworkId.String(),
			BackhaulId:    validBackhaulId.String(),
			PowerId:       validPowerId.String(),
			AccessId:      validAccessId.String(),
			SwitchId:      validSwitchId.String(),
			SpectrumId:    validSpectrumId.String(),
			IsDeactivated: false,
			Latitude:      90.0,  // North Pole
			Longitude:     180.0, // International Date Line
			Location:      "Extreme Location",
			InstallDate:   validInstallDate,
		}

		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: validNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			validBackhaulId.String(),
			validPowerId.String(),
			validAccessId.String(),
			validSwitchId.String(),
			validSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", validNetworkId).Return(int64(1), nil)

		// Mock message bus
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		resp, err := s.Add(context.Background(), requestWithExtremeCoords)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.Equal(t, 90.0, resp.Site.Latitude)
		assert.Equal(t, 180.0, resp.Site.Longitude)
		assert.Equal(t, "Extreme Location", resp.Site.Location)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("DeactivatedSite", func(t *testing.T) {
		deactivatedRequest := &pb.AddRequest{
			Name:          "Deactivated Site",
			NetworkId:     validNetworkId.String(),
			BackhaulId:    validBackhaulId.String(),
			PowerId:       validPowerId.String(),
			AccessId:      validAccessId.String(),
			SwitchId:      validSwitchId.String(),
			SpectrumId:    validSpectrumId.String(),
			IsDeactivated: true,
			Latitude:      40.7128,
			Longitude:     -74.0060,
			Location:      "New York",
			InstallDate:   validInstallDate,
		}

		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: validNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			validBackhaulId.String(),
			validPowerId.String(),
			validAccessId.String(),
			validSwitchId.String(),
			validSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", validNetworkId).Return(int64(1), nil)

		// Mock message bus
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		resp, err := s.Add(context.Background(), deactivatedRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		assert.True(t, resp.Site.IsDeactivated)
		assert.Equal(t, "Deactivated Site", resp.Site.Name)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})
}
