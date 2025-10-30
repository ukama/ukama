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

	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	netmocks "github.com/ukama/ukama/systems/registry/network/pb/gen/mocks"
)

func TestNewNetworkRegistry(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// This test is limited since NewNetworkRegistry creates a real gRPC connection
		// In a real scenario, you might want to use a test server or mock the connection
		networkHost := "localhost:9090"
		timeout := 5 * time.Second

		// Note: This will fail if there's no actual network service running
		// In practice, you might want to use a test server or skip this test
		registry := NewNetworkRegistry(networkHost, timeout)

		assert.NotNil(t, registry)
		assert.Equal(t, networkHost, registry.host)
		assert.Equal(t, timeout, registry.timeout)
		assert.NotNil(t, registry.client)
		assert.NotNil(t, registry.conn)

		// Clean up
		registry.Close()
	})
}

func TestNewNetworkRegistryFromClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		assert.NotNil(t, registry)
		assert.Equal(t, "localhost", registry.host)
		assert.Equal(t, 1*time.Second, registry.timeout)
		assert.Nil(t, registry.conn)
		assert.Equal(t, mockClient, registry.client)
	})
}

func TestNetworkRegistry_Close(t *testing.T) {
	t.Run("WithoutConnection", func(t *testing.T) {
		registry := &NetworkRegistry{
			conn: nil,
		}

		// Should not panic
		assert.NotPanics(t, func() {
			registry.Close()
		})
	})
}

func TestNetworkRegistry_AddNetwork(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		expectedResponse := &netpb.AddResponse{
			Network: &netpb.Network{
				Id:   "test-network-id",
				Name: "test-network",
			},
		}

		mockClient.On("Add", mock.Anything, mock.MatchedBy(func(req *netpb.AddRequest) bool {
			return req.Name == "test-network" &&
				req.Budget == 100.0 &&
				req.Overdraft == 50.0 &&
				req.TrafficPolicy == 1 &&
				req.PaymentLinks == true
		})).Return(expectedResponse, nil)

		response, err := registry.AddNetwork(
			"test-network",
			[]string{"US", "CA"},
			[]string{"network1", "network2"},
			100.0,
			50.0,
			1,
			true,
		)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		expectedError := status.Error(codes.Internal, "internal error")
		mockClient.On("Add", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := registry.AddNetwork(
			"test-network",
			[]string{"US"},
			[]string{"network1"},
			100.0,
			50.0,
			1,
			false,
		)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNetworkRegistry_GetNetwork(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		expectedResponse := &netpb.GetResponse{
			Network: &netpb.Network{
				Id:   "test-network-id",
				Name: "test-network",
			},
		}

		mockClient.On("Get", mock.Anything, &netpb.GetRequest{NetworkId: "test-network-id"}).
			Return(expectedResponse, nil)

		response, err := registry.GetNetwork("test-network-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "network not found")
		mockClient.On("Get", mock.Anything, &netpb.GetRequest{NetworkId: "non-existent-id"}).
			Return(nil, expectedError)

		response, err := registry.GetNetwork("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNetworkRegistry_SetNetworkDefault(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		expectedResponse := &netpb.SetDefaultResponse{}

		mockClient.On("SetDefault", mock.Anything, &netpb.SetDefaultRequest{NetworkId: "test-network-id"}).
			Return(expectedResponse, nil)

		response, err := registry.SetNetworkDefault("test-network-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "network not found")
		mockClient.On("SetDefault", mock.Anything, &netpb.SetDefaultRequest{NetworkId: "non-existent-id"}).
			Return(nil, expectedError)

		response, err := registry.SetNetworkDefault("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNetworkRegistry_GetDefault(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		expectedResponse := &netpb.GetDefaultResponse{
			Network: &netpb.Network{
				Id:   "default-network-id",
				Name: "default-network",
			},
		}

		mockClient.On("GetDefault", mock.Anything, &netpb.GetDefaultRequest{}).
			Return(expectedResponse, nil)

		response, err := registry.GetDefault()

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "no default network set")
		mockClient.On("GetDefault", mock.Anything, &netpb.GetDefaultRequest{}).
			Return(nil, expectedError)

		response, err := registry.GetDefault()

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNetworkRegistry_GetNetworks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		expectedResponse := &netpb.GetNetworksResponse{
			Networks: []*netpb.Network{
				{
					Id:   "network-1",
					Name: "Network 1",
				},
				{
					Id:   "network-2",
					Name: "Network 2",
				},
			},
		}

		mockClient.On("GetAll", mock.Anything, &netpb.GetNetworksRequest{}).
			Return(expectedResponse, nil)

		response, err := registry.GetNetworks()

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("SuccessWithNilNetworks", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		// Simulate response with nil networks
		responseWithNilNetworks := &netpb.GetNetworksResponse{
			Networks: nil,
		}

		mockClient.On("GetAll", mock.Anything, &netpb.GetNetworksRequest{}).
			Return(responseWithNilNetworks, nil)

		response, err := registry.GetNetworks()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotNil(t, response.Networks)
		assert.Empty(t, response.Networks)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &netmocks.NetworkServiceClient{}
		registry := NewNetworkRegistryFromClient(mockClient)

		expectedError := status.Error(codes.Internal, "internal server error")
		mockClient.On("GetAll", mock.Anything, &netpb.GetNetworksRequest{}).
			Return(nil, expectedError)

		response, err := registry.GetNetworks()

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}
