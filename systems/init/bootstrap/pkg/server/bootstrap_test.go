/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/rest/client/factory"
	"github.com/ukama/ukama/systems/init/bootstrap/mocks"
	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	lpb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	lmocks "github.com/ukama/ukama/systems/init/lookup/pb/gen/mocks"
)

func TestGetNodeCredentials(t *testing.T) {
	tests := []struct {
		name           string
		nodeID         string
		setupMocks     func(*mocks.NodeFactoryClient, *mocks.LookupClientProvider, *mbmocks.MsgBusServiceClient)
		expectedResult *pb.GetNodeCredentialsResponse
		expectedError  error
	}{
		{
			name:   "Success - Node with org name and messaging system",
			nodeID: "test-node-123",
			setupMocks: func(factoryMock *mocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
				// Setup factory mock to return a node with org name
				factoryMock.On("Get", "test-node-123").Return(&factory.NodeFactoryInfo{
					Id:      "test-node-123",
					OrgName: "test-org",
				}, nil)

				// Setup lookup client provider mock
				lookupClientMock := &lmocks.LookupServiceClient{}
				lookupMock.On("GetClient").Return(lookupClientMock, nil)

				// Setup lookup service mock to return messaging system
				lookupClientMock.On("GetSystemForOrg", mock.Anything, &lpb.GetSystemRequest{
					OrgName:    "test-org",
					SystemName: "messaging",
				}).Return(&lpb.GetSystemResponse{
					Ip:          "0.0.0.0",
					Certificate: "test-certificate-data",
				}, nil)
			},
			expectedResult: &pb.GetNodeCredentialsResponse{
				Id:          "test-node-123",
				OrgName:     "test-org",
				Ip:          "0.0.0.0",
				Certificate: "test-certificate-data",
			},
			expectedError: nil,
		},
		{
			name:   "Success - Node without org name",
			nodeID: "test-node-456",
			setupMocks: func(factoryMock *mocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
				// Setup factory mock to return a node without org name
				factoryMock.On("Get", "test-node-456").Return(&factory.NodeFactoryInfo{
					Id:      "test-node-456",
					OrgName: "",
				}, nil)

				// Setup lookup client provider mock (it's always called)
				lookupClientMock := &lmocks.LookupServiceClient{}
				lookupMock.On("GetClient").Return(lookupClientMock, nil)
			},
			expectedResult: &pb.GetNodeCredentialsResponse{
				Id:          "test-node-456",
				OrgName:     "",
				Ip:          "",
				Certificate: "",
			},
			expectedError: nil,
		},
		{
			name:   "Error - Factory client fails",
			nodeID: "test-node-789",
			setupMocks: func(factoryMock *mocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
				// Setup factory mock to return an error
				factoryMock.On("Get", "test-node-789").Return(nil, errors.New("factory service unavailable"))
			},
			expectedResult: nil,
			expectedError:  errors.New("factory service unavailable"),
		},
		{
			name:   "Error - Lookup client provider fails",
			nodeID: "test-node-101",
			setupMocks: func(factoryMock *mocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
				// Setup factory mock to return a node with org name
				factoryMock.On("Get", "test-node-101").Return(&factory.NodeFactoryInfo{
					Id:      "test-node-101",
					OrgName: "test-org",
				}, nil)

				// Setup lookup client provider mock to return an error
				lookupMock.On("GetClient").Return(nil, errors.New("lookup service connection failed"))
			},
			expectedResult: nil,
			expectedError:  errors.New("lookup service connection failed"),
		},
		{
			name:   "Success - Lookup service fails but response still returned",
			nodeID: "test-node-202",
			setupMocks: func(factoryMock *mocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
				// Setup factory mock to return a node with org name
				factoryMock.On("Get", "test-node-202").Return(&factory.NodeFactoryInfo{
					Id:      "test-node-202",
					OrgName: "test-org",
				}, nil)

				// Setup lookup client provider mock
				lookupClientMock := &lmocks.LookupServiceClient{}
				lookupMock.On("GetClient").Return(lookupClientMock, nil)

				// Setup lookup service mock to return an error
				lookupClientMock.On("GetSystemForOrg", mock.Anything, &lpb.GetSystemRequest{
					OrgName:    "test-org",
					SystemName: "messaging",
				}).Return(nil, errors.New("messaging system not found"))
			},
			expectedResult: &pb.GetNodeCredentialsResponse{
				Id:          "test-node-202",
				OrgName:     "test-org",
				Ip:          "",
				Certificate: "",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			factoryMock := mocks.NewNodeFactoryClient(t)
			lookupMock := mocks.NewLookupClientProvider(t)
			msgBusMock := mbmocks.NewMsgBusServiceClient(t)

			// Setup mocks
			tt.setupMocks(factoryMock, lookupMock, msgBusMock)

			// Create server instance
			server := NewBootstrapServer("test-org", msgBusMock, false, lookupMock, factoryMock)

			// Create request
			req := &pb.GetNodeCredentialsRequest{
				Id: tt.nodeID,
			}

			// Call the method
			result, err := server.GetNodeCredentials(context.Background(), req)

			// Assertions
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Id, result.Id)
				assert.Equal(t, tt.expectedResult.OrgName, result.OrgName)
				assert.Equal(t, tt.expectedResult.Ip, result.Ip)
				assert.Equal(t, tt.expectedResult.Certificate, result.Certificate)
			}

			// Verify all mocks were called as expected
			factoryMock.AssertExpectations(t)
			lookupMock.AssertExpectations(t)
			msgBusMock.AssertExpectations(t)
		})
	}
}

func TestGetNodeCredentials_EdgeCases(t *testing.T) {
	t.Run("Empty node ID", func(t *testing.T) {
		factoryMock := mocks.NewNodeFactoryClient(t)
		lookupMock := mocks.NewLookupClientProvider(t)
		msgBusMock := mbmocks.NewMsgBusServiceClient(t)

		// Setup factory mock to return an error for empty ID
		factoryMock.On("Get", "").Return(nil, errors.New("invalid node ID"))

		server := NewBootstrapServer("test-org", msgBusMock, false, lookupMock, factoryMock)
		req := &pb.GetNodeCredentialsRequest{Id: ""}

		result, err := server.GetNodeCredentials(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid node ID", err.Error())
	})

	t.Run("Nil request", func(t *testing.T) {
		factoryMock := mocks.NewNodeFactoryClient(t)
		lookupMock := mocks.NewLookupClientProvider(t)
		msgBusMock := mbmocks.NewMsgBusServiceClient(t)

		server := NewBootstrapServer("test-org", msgBusMock, false, lookupMock, factoryMock)

		// This should panic or return an error due to nil request
		assert.Panics(t, func() {
			server.GetNodeCredentials(context.Background(), nil)
		})
	})

	t.Run("Context cancellation", func(t *testing.T) {
		factoryMock := mocks.NewNodeFactoryClient(t)
		lookupMock := mocks.NewLookupClientProvider(t)
		msgBusMock := mbmocks.NewMsgBusServiceClient(t)

		// Setup factory mock to simulate a delay
		factoryMock.On("Get", "test-node").Return(&factory.NodeFactoryInfo{
			Id:      "test-node",
			OrgName: "test-org",
		}, nil)

		// Setup lookup client provider mock
		lookupClientMock := &lmocks.LookupServiceClient{}
		lookupMock.On("GetClient").Return(lookupClientMock, nil)

		// Setup lookup service mock to return messaging system
		lookupClientMock.On("GetSystemForOrg", mock.Anything, &lpb.GetSystemRequest{
			OrgName:    "test-org",
			SystemName: "messaging",
		}).Return(&lpb.GetSystemResponse{
			Ip:          "0.0.0.0",
			Certificate: "test-certificate-data",
		}, nil)

		server := NewBootstrapServer("test-org", msgBusMock, false, lookupMock, factoryMock)
		req := &pb.GetNodeCredentialsRequest{Id: "test-node"}

		// Create a cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		result, err := server.GetNodeCredentials(ctx, req)

		// The method should still work as it doesn't use the context for cancellation
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestNewBootstrapServer(t *testing.T) {
	factoryMock := mocks.NewNodeFactoryClient(t)
	lookupMock := mocks.NewLookupClientProvider(t)
	msgBusMock := mbmocks.NewMsgBusServiceClient(t)

	server := NewBootstrapServer("test-org", msgBusMock, true, lookupMock, factoryMock)

	assert.NotNil(t, server)
	assert.Equal(t, "test-org", server.orgName)
	assert.Equal(t, msgBusMock, server.msgbus)
	assert.True(t, server.debug)
	assert.Equal(t, lookupMock, server.lookupClient)
	assert.Equal(t, factoryMock, server.factoryClient)
	assert.NotNil(t, server.bootstrapRoutingKey)
}
