/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ukama/ukama/systems/registry/site/pb/gen"
	sitemocks "github.com/ukama/ukama/systems/registry/site/pb/gen/mocks"
)

func TestNewSiteRegistry(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// This test is limited since NewSiteRegistry creates a real gRPC connection
		// In a real scenario, you might want to use a test server or mock the connection
		siteHost := "localhost:9090"
		timeout := 5 * time.Second

		// Note: This will fail if there's no actual site service running
		// In practice, you might want to use a test server or skip this test
		registry := NewSiteRegistry(siteHost, timeout)

		assert.NotNil(t, registry)
		assert.Equal(t, siteHost, registry.host)
		assert.Equal(t, timeout, registry.timeout)
		assert.NotNil(t, registry.client)
		assert.NotNil(t, registry.conn)

		// Clean up
		registry.Close()
	})
}

func TestNewSiteRegistryFromClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &sitemocks.SiteServiceClient{}
		registry := NewSiteRegistryFromClient(mockClient)

		assert.NotNil(t, registry)
		assert.Equal(t, "localhost", registry.host)
		assert.Equal(t, 1*time.Second, registry.timeout)
		assert.Nil(t, registry.conn)
		assert.Equal(t, mockClient, registry.client)
	})
}

func TestSiteRegistry_Close(t *testing.T) {
	t.Run("WithoutConnection", func(t *testing.T) {
		registry := &SiteRegistry{
			conn: nil,
		}

		// Should not panic
		assert.NotPanics(t, func() {
			registry.Close()
		})
	})
}

func TestSiteRegistry_GetSite(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &sitemocks.SiteServiceClient{}
		registry := NewSiteRegistryFromClient(mockClient)

		expectedResponse := &pb.GetResponse{
			Site: &pb.Site{
				Id:            "test-site-id",
				Name:          "Test Site",
				NetworkId:     "test-network-id",
				Location:      "Test Location",
				BackhaulId:    "backhaul-1",
				PowerId:       "power-1",
				AccessId:      "access-1",
				SwitchId:      "switch-1",
				SpectrumId:    "spectrum-1",
				IsDeactivated: false,
				Latitude:      "40.7128",
				Longitude:     "-74.0060",
				InstallDate:   "2023-01-01",
			},
		}

		mockClient.On("Get", mock.Anything, &pb.GetRequest{SiteId: "test-site-id"}).
			Return(expectedResponse, nil)

		response, err := registry.GetSite("test-site-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &sitemocks.SiteServiceClient{}
		registry := NewSiteRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "site not found")
		mockClient.On("Get", mock.Anything, &pb.GetRequest{SiteId: "non-existent-id"}).
			Return(nil, expectedError)

		response, err := registry.GetSite("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestSiteRegistry_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &sitemocks.SiteServiceClient{}
		registry := NewSiteRegistryFromClient(mockClient)

		expectedResponse := &pb.ListResponse{
			Sites: []*pb.Site{
				{
					Id:            "site-1",
					Name:          "Site 1",
					NetworkId:     "network-1",
					Location:      "Location 1",
					IsDeactivated: false,
				},
				{
					Id:            "site-2",
					Name:          "Site 2",
					NetworkId:     "network-1",
					Location:      "Location 2",
					IsDeactivated: true,
				},
			},
		}

		mockClient.On("List", mock.Anything, &pb.ListRequest{NetworkId: "network-1", IsDeactivated: false}).
			Return(expectedResponse, nil)

		response, err := registry.List("network-1", false)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &sitemocks.SiteServiceClient{}
		registry := NewSiteRegistryFromClient(mockClient)

		expectedError := status.Error(codes.Internal, "internal server error")
		mockClient.On("List", mock.Anything, &pb.ListRequest{NetworkId: "network-1", IsDeactivated: false}).
			Return(nil, expectedError)

		response, err := registry.List("network-1", false)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestSiteRegistry_AddSite(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &sitemocks.SiteServiceClient{}
		registry := NewSiteRegistryFromClient(mockClient)

		expectedResponse := &pb.AddResponse{
			Site: &pb.Site{
				Id:            "new-site-id",
				Name:          "New Site",
				NetworkId:     "network-1",
				Location:      "New Location",
				BackhaulId:    "backhaul-1",
				PowerId:       "power-1",
				AccessId:      "access-1",
				SwitchId:      "switch-1",
				SpectrumId:    "spectrum-1",
				IsDeactivated: false,
				Latitude:      "40.7128",
				Longitude:     "-74.0060",
				InstallDate:   "2023-01-01",
			},
		}

		mockClient.On("Add", mock.Anything, mock.MatchedBy(func(req *pb.AddRequest) bool {
			return req.Name == "New Site" &&
				req.NetworkId == "network-1" &&
				req.Location == "New Location" &&
				req.BackhaulId == "backhaul-1" &&
				req.PowerId == "power-1" &&
				req.AccessId == "access-1" &&
				req.SwitchId == "switch-1" &&
				req.SpectrumId == "spectrum-1" &&
				req.IsDeactivated == false &&
				req.Latitude == "40.7128" &&
				req.Longitude == "-74.0060" &&
				req.InstallDate == "2023-01-01"
		})).Return(expectedResponse, nil)

		response, err := registry.AddSite(
			"network-1",
			"New Site",
			"backhaul-1",
			"power-1",
			"access-1",
			"switch-1",
			"New Location",
			"spectrum-1",
			false,
			"40.7128",
			"-74.0060",
			"2023-01-01",
		)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &sitemocks.SiteServiceClient{}
		registry := NewSiteRegistryFromClient(mockClient)

		expectedError := status.Error(codes.InvalidArgument, "invalid site data")
		mockClient.On("Add", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := registry.AddSite(
			"network-1",
			"Invalid Site",
			"backhaul-1",
			"power-1",
			"access-1",
			"switch-1",
			"Invalid Location",
			"spectrum-1",
			false,
			"40.7128",
			"-74.0060",
			"invalid-date",
		)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestSiteRegistry_UpdateSite(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &sitemocks.SiteServiceClient{}
		registry := NewSiteRegistryFromClient(mockClient)

		expectedResponse := &pb.UpdateResponse{
			Site: &pb.Site{
				Id:   "test-site-id",
				Name: "Updated Site Name",
			},
		}

		mockClient.On("Update", mock.Anything, mock.MatchedBy(func(req *pb.UpdateRequest) bool {
			return req.SiteId == "test-site-id" &&
				req.Name == "Updated Site Name"
		})).Return(expectedResponse, nil)

		response, err := registry.UpdateSite("test-site-id", "Updated Site Name")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &sitemocks.SiteServiceClient{}
		registry := NewSiteRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "site not found")
		mockClient.On("Update", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := registry.UpdateSite("non-existent-id", "Updated Name")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}
