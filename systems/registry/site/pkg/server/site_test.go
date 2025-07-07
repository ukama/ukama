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

// Test data constants
const (
	testSiteName            = "Test Site"
	testSiteLocation        = "Test Location"
	testUpdatedSiteName     = "Updated Site Name"
	testLongSiteName        = "This is a very long site name that might exceed some limits but should still be processed by the update function"
	testDeactivatedSiteName = "Deactivated Site"
	testExtremeLocation     = "Extreme Location"
	testNewYorkLocation     = "New York"
	testInvalidUUID         = "invalid-uuid-format"
	testInvalidDate         = "invalid-date"
	testValidInstallDate    = "2023-12-01T00:00:00Z"

	// Coordinates
	testLatitude         = 40.7128
	testLongitude        = -74.0060
	testExtremeLatitude  = 90.0  // North Pole
	testExtremeLongitude = 180.0 // International Date Line

	// Site names for list tests
	testSite1Name               = "Site1"
	testSite2Name               = "Site2"
	testDeactivatedSiteListName = "DeactivatedSite"
)

// Test UUIDs - these will be generated fresh for each test
var (
	testSiteId     = uuid.NewV4()
	testNetworkId  = uuid.NewV4()
	testBackhaulId = uuid.NewV4()
	testPowerId    = uuid.NewV4()
	testAccessId   = uuid.NewV4()
	testSwitchId   = uuid.NewV4()
	testSpectrumId = uuid.NewV4()
)

// Helper function to create a mock site with test data
func createMockSite() *db.Site {
	return &db.Site{
		Id:            testSiteId,
		Name:          testSiteName,
		Location:      testSiteLocation,
		NetworkId:     testNetworkId,
		BackhaulId:    testBackhaulId,
		PowerId:       testPowerId,
		AccessId:      testAccessId,
		SwitchId:      testSwitchId,
		SpectrumId:    testSpectrumId,
		IsDeactivated: false,
		Latitude:      testLatitude,
		Longitude:     testLongitude,
		InstallDate:   testValidInstallDate,
	}
}

// Helper function to create a valid AddRequest with test data
func createValidAddRequest() *pb.AddRequest {
	return &pb.AddRequest{
		Name:          testSiteName,
		NetworkId:     testNetworkId.String(),
		BackhaulId:    testBackhaulId.String(),
		PowerId:       testPowerId.String(),
		AccessId:      testAccessId.String(),
		SwitchId:      testSwitchId.String(),
		SpectrumId:    testSpectrumId.String(),
		IsDeactivated: false,
		Latitude:      testLatitude,
		Longitude:     testLongitude,
		Location:      testNewYorkLocation,
		InstallDate:   testValidInstallDate,
	}
}

// Helper function to verify site field mappings
func verifySiteFieldMappings(t *testing.T, expected *db.Site, actual *pb.Site) {
	fieldMappings := []struct {
		name     string
		expected interface{}
		actual   interface{}
	}{
		{"Id", expected.Id.String(), actual.Id},
		{"Name", expected.Name, actual.Name},
		{"Location", expected.Location, actual.Location},
		{"NetworkId", expected.NetworkId.String(), actual.NetworkId},
		{"BackhaulId", expected.BackhaulId.String(), actual.BackhaulId},
		{"PowerId", expected.PowerId.String(), actual.PowerId},
		{"AccessId", expected.AccessId.String(), actual.AccessId},
		{"SwitchId", expected.SwitchId.String(), actual.SwitchId},
		{"SpectrumId", expected.SpectrumId.String(), actual.SpectrumId},
		{"IsDeactivated", expected.IsDeactivated, actual.IsDeactivated},
		{"Latitude", expected.Latitude, actual.Latitude},
		{"Longitude", expected.Longitude, actual.Longitude},
		{"InstallDate", expected.InstallDate, actual.InstallDate},
	}

	for _, field := range fieldMappings {
		assert.Equal(t, field.expected, field.actual, "Field %s mismatch", field.name)
	}
}

// Helper function to verify specific site fields (for Update tests)
func verifySiteFields(t *testing.T, expectedId string, expectedName string, actual *pb.Site) {
	assert.Equal(t, expectedId, actual.Id, "Site ID mismatch")
	assert.Equal(t, expectedName, actual.Name, "Site Name mismatch")
}

// Helper function to verify Add request fields
func verifyAddRequestFields(t *testing.T, expected *pb.AddRequest, actual *pb.Site) {
	assert.Equal(t, expected.Name, actual.Name, "Site Name mismatch")
	assert.Equal(t, expected.NetworkId, actual.NetworkId, "Network ID mismatch")
	assert.Equal(t, expected.BackhaulId, actual.BackhaulId, "Backhaul ID mismatch")
	assert.Equal(t, expected.PowerId, actual.PowerId, "Power ID mismatch")
	assert.Equal(t, expected.AccessId, actual.AccessId, "Access ID mismatch")
	assert.Equal(t, expected.SwitchId, actual.SwitchId, "Switch ID mismatch")
	assert.Equal(t, expected.SpectrumId, actual.SpectrumId, "Spectrum ID mismatch")
	assert.Equal(t, expected.IsDeactivated, actual.IsDeactivated, "IsDeactivated mismatch")
	assert.Equal(t, expected.Latitude, actual.Latitude, "Latitude mismatch")
	assert.Equal(t, expected.Longitude, actual.Longitude, "Longitude mismatch")
	assert.Equal(t, expected.Location, actual.Location, "Location mismatch")
}

func TestSiteService_Get(t *testing.T) {
	siteRepo := &mocks.SiteRepo{}
	msgclientRepo := &cmocks.MsgBusServiceClient{}
	netRepo := &mocks.NetworkClientProvider{}

	s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", nil)

	t.Run("SiteFound", func(t *testing.T) {
		mockSite := createMockSite()

		siteRepo.On("Get", testSiteId).Return(mockSite, nil).Once()

		resp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: testSiteId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)

		// Verify all fields are correctly mapped
		verifySiteFieldMappings(t, mockSite, resp.Site)

		siteRepo.AssertExpectations(t)
	})

	t.Run("SiteNotFound", func(t *testing.T) {
		siteRepo.On("Get", testSiteId).Return(nil, gorm.ErrRecordNotFound).Once()

		resp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: testSiteId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "record not found")
		siteRepo.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		resp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: testInvalidUUID})

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
		dbError := gorm.ErrInvalidDB

		siteRepo.On("Get", testSiteId).Return(nil, dbError).Once()

		resp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: testSiteId.String()})

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

	t.Run("ValidRequestWithMultipleSites", func(t *testing.T) {
		siteRepo.ExpectedCalls = nil

		mockSites := []*db.Site{
			{
				Id:            uuid.NewV4(),
				NetworkId:     testNetworkId,
				Name:          testSite1Name,
				Location:      testSiteLocation,
				BackhaulId:    testBackhaulId,
				PowerId:       testPowerId,
				AccessId:      testAccessId,
				SwitchId:      testSwitchId,
				SpectrumId:    testSpectrumId,
				IsDeactivated: false,
				Latitude:      testLatitude,
				Longitude:     testLongitude,
				InstallDate:   testValidInstallDate,
			},
			{
				Id:            uuid.NewV4(),
				NetworkId:     testNetworkId,
				Name:          testSite2Name,
				Location:      testNewYorkLocation,
				BackhaulId:    testBackhaulId,
				PowerId:       testPowerId,
				AccessId:      testAccessId,
				SwitchId:      testSwitchId,
				SpectrumId:    testSpectrumId,
				IsDeactivated: false,
				Latitude:      testExtremeLatitude,
				Longitude:     testExtremeLongitude,
				InstallDate:   testValidInstallDate,
			},
		}

		var mockSitesConverted []db.Site
		for _, site := range mockSites {
			mockSitesConverted = append(mockSitesConverted, *site)
		}

		siteRepo.On("List", &testNetworkId, false).Return(mockSitesConverted, nil).Once()

		req := &pb.ListRequest{
			NetworkId:     testNetworkId.String(),
			IsDeactivated: false,
		}

		resp, err := s.List(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, len(mockSites), len(resp.Sites))

		// Verify all fields are correctly mapped for each site
		for i, site := range mockSites {
			verifySiteFieldMappings(t, site, resp.Sites[i])
		}
		siteRepo.AssertExpectations(t)
	})

	t.Run("EmptyNetworkId", func(t *testing.T) {
		// Reset mock for this test
		siteRepo.ExpectedCalls = nil

		mockSites := []db.Site{
			{
				Id:            uuid.NewV4(),
				Name:          testSite1Name,
				Location:      testSiteLocation,
				NetworkId:     testNetworkId,
				BackhaulId:    testBackhaulId,
				PowerId:       testPowerId,
				AccessId:      testAccessId,
				SwitchId:      testSwitchId,
				SpectrumId:    testSpectrumId,
				IsDeactivated: false,
				Latitude:      testLatitude,
				Longitude:     testLongitude,
				InstallDate:   testValidInstallDate,
			},
		}
		siteRepo.On("List", (*uuid.UUID)(nil), false).Return(mockSites, nil).Once()

		req := &pb.ListRequest{NetworkId: "", IsDeactivated: false}
		resp, err := s.List(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, len(resp.Sites))
		verifySiteFieldMappings(t, &mockSites[0], resp.Sites[0])
		siteRepo.AssertExpectations(t)
	})

	t.Run("InvalidNetworkId", func(t *testing.T) {
		req := &pb.ListRequest{NetworkId: testInvalidUUID, IsDeactivated: false}
		resp, err := s.List(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid network ID")
	})

	t.Run("IsDeactivatedTrue", func(t *testing.T) {
		siteRepo.ExpectedCalls = nil

		mockSites := []db.Site{
			{
				Id:            uuid.NewV4(),
				NetworkId:     testNetworkId,
				Name:          testDeactivatedSiteListName,
				Location:      testSiteLocation,
				BackhaulId:    testBackhaulId,
				PowerId:       testPowerId,
				AccessId:      testAccessId,
				SwitchId:      testSwitchId,
				SpectrumId:    testSpectrumId,
				IsDeactivated: true,
				Latitude:      testLatitude,
				Longitude:     testLongitude,
				InstallDate:   testValidInstallDate,
			},
		}
		siteRepo.On("List", &testNetworkId, true).Return(mockSites, nil).Once()

		req := &pb.ListRequest{NetworkId: testNetworkId.String(), IsDeactivated: true}
		resp, err := s.List(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, len(resp.Sites))
		verifySiteFieldMappings(t, &mockSites[0], resp.Sites[0])
		siteRepo.AssertExpectations(t)
	})

	t.Run("EmptyResultSet", func(t *testing.T) {
		siteRepo.ExpectedCalls = nil

		mockSites := []db.Site{}
		siteRepo.On("List", &testNetworkId, false).Return(mockSites, nil).Once()

		req := &pb.ListRequest{NetworkId: testNetworkId.String(), IsDeactivated: false}
		resp, err := s.List(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 0, len(resp.Sites))
		assert.Empty(t, resp.Sites)
		siteRepo.AssertExpectations(t)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		siteRepo.ExpectedCalls = nil

		dbError := gorm.ErrInvalidDB
		siteRepo.On("List", &testNetworkId, false).Return(nil, dbError).Once()

		req := &pb.ListRequest{NetworkId: testNetworkId.String(), IsDeactivated: false}
		resp, err := s.List(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "invalid db")
		siteRepo.AssertExpectations(t)
	})
}

func TestSiteService_Update(t *testing.T) {
	siteRepo := &mocks.SiteRepo{}
	msgclientRepo := &cmocks.MsgBusServiceClient{}
	netRepo := &mocks.NetworkClientProvider{}

	s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", nil)

	t.Run("Success", func(t *testing.T) {
		// Mock the site repository update
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(nil).Once()

		// Mock message bus publish
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		req := &pb.UpdateRequest{
			SiteId: testSiteId.String(),
			Name:   testUpdatedSiteName,
		}

		resp, err := s.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifySiteFields(t, testSiteId.String(), testUpdatedSiteName, resp.Site)

		// Verify the mock was called with correct parameters
		siteRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("InvalidSiteId", func(t *testing.T) {
		req := &pb.UpdateRequest{
			SiteId: testInvalidUUID,
			Name:   testSiteName,
		}

		resp, err := s.Update(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("EmptySiteId", func(t *testing.T) {
		req := &pb.UpdateRequest{
			SiteId: "",
			Name:   testSiteName,
		}

		resp, err := s.Update(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("EmptyName", func(t *testing.T) {
		// Mock the site repository update
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(nil).Once()

		// Mock message bus publish
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		req := &pb.UpdateRequest{
			SiteId: testSiteId.String(),
			Name:   "",
		}

		resp, err := s.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifySiteFields(t, testSiteId.String(), "", resp.Site)

		siteRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Mock the site repository update to return an error
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(gorm.ErrRecordNotFound).Once()

		req := &pb.UpdateRequest{
			SiteId: testSiteId.String(),
			Name:   testSiteName,
		}

		resp, err := s.Update(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "record not found")

		siteRepo.AssertExpectations(t)
	})

	t.Run("MessageBusError", func(t *testing.T) {
		// Mock the site repository update
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(nil).Once()

		// Mock message bus publish to return an error
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(gorm.ErrInvalidDB).Once()

		req := &pb.UpdateRequest{
			SiteId: testSiteId.String(),
			Name:   testUpdatedSiteName,
		}

		resp, err := s.Update(context.Background(), req)

		// The update should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifySiteFields(t, testSiteId.String(), testUpdatedSiteName, resp.Site)

		siteRepo.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("NilMessageBus", func(t *testing.T) {
		// Create server without message bus
		sNoMsgBus := NewSiteServer(OrgName, siteRepo, nil, netRepo, "", nil)

		// Mock the site repository update
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(nil).Once()

		req := &pb.UpdateRequest{
			SiteId: testSiteId.String(),
			Name:   testUpdatedSiteName,
		}

		resp, err := sNoMsgBus.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifySiteFields(t, testSiteId.String(), testUpdatedSiteName, resp.Site)

		siteRepo.AssertExpectations(t)
	})

	t.Run("LongName", func(t *testing.T) {
		// Mock the site repository update
		siteRepo.On("Update", mock.AnythingOfType("*db.Site")).Return(nil).Once()

		// Mock message bus publish
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		req := &pb.UpdateRequest{
			SiteId: testSiteId.String(),
			Name:   testLongSiteName,
		}

		resp, err := s.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifySiteFields(t, testSiteId.String(), testLongSiteName, resp.Site)

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

	validRequest := createValidAddRequest()

	t.Run("Success", func(t *testing.T) {
		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: testNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			testBackhaulId.String(),
			testPowerId.String(),
			testAccessId.String(),
			testSwitchId.String(),
			testSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", testNetworkId).Return(int64(1), nil)

		// Mock message bus
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		resp, err := s.Add(context.Background(), validRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifyAddRequestFields(t, validRequest, resp.Site)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("InvalidNetworkId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       testSiteName,
			NetworkId:  testInvalidUUID,
			BackhaulId: testBackhaulId.String(),
			PowerId:    testPowerId.String(),
			AccessId:   testAccessId.String(),
			SwitchId:   testSwitchId.String(),
			SpectrumId: testSpectrumId.String(),
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidBackhaulId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       testSiteName,
			NetworkId:  testNetworkId.String(),
			BackhaulId: testInvalidUUID,
			PowerId:    testPowerId.String(),
			AccessId:   testAccessId.String(),
			SwitchId:   testSwitchId.String(),
			SpectrumId: testSpectrumId.String(),
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidPowerId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       testSiteName,
			NetworkId:  testNetworkId.String(),
			BackhaulId: testBackhaulId.String(),
			PowerId:    testInvalidUUID,
			AccessId:   testAccessId.String(),
			SwitchId:   testSwitchId.String(),
			SpectrumId: testSpectrumId.String(),
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidAccessId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       testSiteName,
			NetworkId:  testNetworkId.String(),
			BackhaulId: testBackhaulId.String(),
			PowerId:    testPowerId.String(),
			AccessId:   testInvalidUUID,
			SwitchId:   testSwitchId.String(),
			SpectrumId: testSpectrumId.String(),
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidSwitchId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       testSiteName,
			NetworkId:  testNetworkId.String(),
			BackhaulId: testBackhaulId.String(),
			PowerId:    testPowerId.String(),
			AccessId:   testAccessId.String(),
			SwitchId:   testInvalidUUID,
			SpectrumId: testSpectrumId.String(),
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidSpectrumId", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:       testSiteName,
			NetworkId:  testNetworkId.String(),
			BackhaulId: testBackhaulId.String(),
			PowerId:    testPowerId.String(),
			AccessId:   testAccessId.String(),
			SwitchId:   testSwitchId.String(),
			SpectrumId: testInvalidUUID,
		}

		resp, err := s.Add(context.Background(), invalidRequest)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "Error parsing UUID")
	})

	t.Run("InvalidInstallDate", func(t *testing.T) {
		invalidRequest := &pb.AddRequest{
			Name:        testSiteName,
			NetworkId:   testNetworkId.String(),
			BackhaulId:  testBackhaulId.String(),
			PowerId:     testPowerId.String(),
			AccessId:    testAccessId.String(),
			SwitchId:    testSwitchId.String(),
			SpectrumId:  testSpectrumId.String(),
			InstallDate: testInvalidDate,
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
			testBackhaulId.String(),
			testPowerId.String(),
			testAccessId.String(),
			testSwitchId.String(),
			testSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock network client to return error
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: testNetworkId.String(),
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
			testBackhaulId.String(),
			testPowerId.String(),
			testAccessId.String(),
			testSwitchId.String(),
			testSpectrumId.String(),
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
			testBackhaulId.String(),
			testPowerId.String(),
			testAccessId.String(),
			testSwitchId.String(),
			testSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: testNetworkId.String(),
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
			NetworkId: testNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			testBackhaulId.String(),
			testPowerId.String(),
			testAccessId.String(),
			testSwitchId.String(),
			testSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", testNetworkId).Return(int64(1), nil)

		// Mock message bus to return error
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(gorm.ErrInvalidDB)

		resp, err := s.Add(context.Background(), validRequest)

		// Should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifyAddRequestFields(t, validRequest, resp.Site)

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
			NetworkId: testNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			testBackhaulId.String(),
			testPowerId.String(),
			testAccessId.String(),
			testSwitchId.String(),
			testSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount to return error
		siteRepo.On("GetSiteCount", testNetworkId).Return(int64(0), gorm.ErrInvalidDB)

		// Mock message bus
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		resp, err := s.Add(context.Background(), validRequest)

		// Should still succeed even if GetSiteCount fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifyAddRequestFields(t, validRequest, resp.Site)

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
			NetworkId: testNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			testBackhaulId.String(),
			testPowerId.String(),
			testAccessId.String(),
			testSwitchId.String(),
			testSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", testNetworkId).Return(int64(1), nil)

		resp, err := sNoMsgBus.Add(context.Background(), validRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifyAddRequestFields(t, validRequest, resp.Site)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
	})

	t.Run("EmptyName", func(t *testing.T) {
		requestWithEmptyName := &pb.AddRequest{
			Name:          "",
			NetworkId:     testNetworkId.String(),
			BackhaulId:    testBackhaulId.String(),
			PowerId:       testPowerId.String(),
			AccessId:      testAccessId.String(),
			SwitchId:      testSwitchId.String(),
			SpectrumId:    testSpectrumId.String(),
			IsDeactivated: false,
			Latitude:      testLatitude,
			Longitude:     testLongitude,
			Location:      testNewYorkLocation,
			InstallDate:   testValidInstallDate,
		}

		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: testNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			testBackhaulId.String(),
			testPowerId.String(),
			testAccessId.String(),
			testSwitchId.String(),
			testSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", testNetworkId).Return(int64(1), nil)

		// Mock message bus
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		resp, err := s.Add(context.Background(), requestWithEmptyName)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifyAddRequestFields(t, requestWithEmptyName, resp.Site)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("ExtremeCoordinates", func(t *testing.T) {
		requestWithExtremeCoords := &pb.AddRequest{
			Name:          testSiteName,
			NetworkId:     testNetworkId.String(),
			BackhaulId:    testBackhaulId.String(),
			PowerId:       testPowerId.String(),
			AccessId:      testAccessId.String(),
			SwitchId:      testSwitchId.String(),
			SpectrumId:    testSpectrumId.String(),
			IsDeactivated: false,
			Latitude:      testExtremeLatitude,
			Longitude:     testExtremeLongitude,
			Location:      testExtremeLocation,
			InstallDate:   testValidInstallDate,
		}

		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: testNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			testBackhaulId.String(),
			testPowerId.String(),
			testAccessId.String(),
			testSwitchId.String(),
			testSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", testNetworkId).Return(int64(1), nil)

		// Mock message bus
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		resp, err := s.Add(context.Background(), requestWithExtremeCoords)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifyAddRequestFields(t, requestWithExtremeCoords, resp.Site)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})

	t.Run("DeactivatedSite", func(t *testing.T) {
		deactivatedRequest := &pb.AddRequest{
			Name:          testDeactivatedSiteName,
			NetworkId:     testNetworkId.String(),
			BackhaulId:    testBackhaulId.String(),
			PowerId:       testPowerId.String(),
			AccessId:      testAccessId.String(),
			SwitchId:      testSwitchId.String(),
			SpectrumId:    testSpectrumId.String(),
			IsDeactivated: true,
			Latitude:      testLatitude,
			Longitude:     testLongitude,
			Location:      testNewYorkLocation,
			InstallDate:   testValidInstallDate,
		}

		// Mock network client
		mockNetworkClient := &netmocks.NetworkServiceClient{}
		netRepo.On("GetClient").Return(mockNetworkClient, nil)
		mockNetworkClient.On("Get", mock.Anything, &npb.GetRequest{
			NetworkId: testNetworkId.String(),
		}).Return(&npb.GetResponse{}, nil)

		// Mock inventory client calls for all components
		for _, componentId := range []string{
			testBackhaulId.String(),
			testPowerId.String(),
			testAccessId.String(),
			testSwitchId.String(),
			testSpectrumId.String(),
		} {
			inventoryClient.On("Get", componentId).Return(&inventory.ComponentInfo{}, nil)
		}

		// Mock site repository Add
		siteRepo.On("Add", mock.AnythingOfType("*db.Site"), mock.AnythingOfType("func(*db.Site, *gorm.DB) error")).Return(nil)

		// Mock GetSiteCount for pushSiteCount
		siteRepo.On("GetSiteCount", testNetworkId).Return(int64(1), nil)

		// Mock message bus
		msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		resp, err := s.Add(context.Background(), deactivatedRequest)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Site)
		verifyAddRequestFields(t, deactivatedRequest, resp.Site)

		siteRepo.AssertExpectations(t)
		netRepo.AssertExpectations(t)
		inventoryClient.AssertExpectations(t)
		msgclientRepo.AssertExpectations(t)
	})
}
