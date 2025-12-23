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
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/rest/client/factory"
	initclient "github.com/ukama/ukama/systems/common/rest/client/initclient"
	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
)

func TestGetNodeCredentials(t *testing.T) {
	tests := []struct {
		name           string
		nodeID         string
		setupMocks     func(*mbmocks.NodeFactoryClient, *mbmocks.InitClient, *mbmocks.MsgBusServiceClient)
		expectedResult *pb.GetNodeCredentialsResponse
		expectedError  error
	}{
		{
			name:   "Success - Node with org name and messaging system",
			nodeID: "test-node-123",
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mbmocks.InitClient, msgBusMock *mbmocks.MsgBusServiceClient) {
				// Setup factory mock to return a node with org name
				factoryMock.On("Get", "test-node-123").Return(&factory.Node{
					Node: factory.NodeFactoryInfo{
						Id:      "test-node-123",
						OrgName: "test-org",
					},
				}, nil)

				// Setup lookup client mock to return messaging system
				lookupMock.On("GetSystem", "test-org", "messaging").Return(&initclient.SystemIPInfo{
					Ip:          "192.168.1.100",
					Certificate: "test-certificate-data",
				}, nil)
			},
			expectedResult: &pb.GetNodeCredentialsResponse{
				Id:          "test-node-123",
				OrgName:     "test-org",
				Ip:          "192.168.1.100",
				Certificate: "test-certificate-data",
			},
			expectedError: nil,
		},
		{
			name:   "Error - Node without org name",
			nodeID: "test-node-456",
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mbmocks.InitClient, msgBusMock *mbmocks.MsgBusServiceClient) {
				// Setup factory mock to return a node without org name
				factoryMock.On("Get", "test-node-456").Return(&factory.Node{
					Node: factory.NodeFactoryInfo{
						Id:      "test-node-456",
						OrgName: "",
					},
				}, nil)

				// Setup lookup client mock to return an error when called with empty org
				lookupMock.On("GetSystem", "", "messaging").Return(nil, errors.New("org name is required"))
			},
			expectedResult: nil,
			expectedError:  errors.New("org name is required"),
		},
		{
			name:   "Error - Factory client fails",
			nodeID: "test-node-789",
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mbmocks.InitClient, msgBusMock *mbmocks.MsgBusServiceClient) {
				// Setup factory mock to return an error
				factoryMock.On("Get", "test-node-789").Return(nil, errors.New("factory service unavailable"))
			},
			expectedResult: nil,
			expectedError:  errors.New("factory service unavailable"),
		},
		{
			name:   "Error - Lookup service fails",
			nodeID: "test-node-101",
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mbmocks.InitClient, msgBusMock *mbmocks.MsgBusServiceClient) {
				// Setup factory mock to return a node with org name
				factoryMock.On("Get", "test-node-101").Return(&factory.Node{
					Node: factory.NodeFactoryInfo{
						Id:      "test-node-101",
						OrgName: "test-org",
					},
				}, nil)

				// Setup lookup client mock to return an error
				lookupMock.On("GetSystem", "test-org", "messaging").Return(nil, errors.New("lookup service connection failed"))
			},
			expectedResult: nil,
			expectedError:  errors.New("lookup service connection failed"),
		},
		{
			name:   "Error - Messaging system not found",
			nodeID: "test-node-202",
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mbmocks.InitClient, msgBusMock *mbmocks.MsgBusServiceClient) {
				// Setup factory mock to return a node with org name
				factoryMock.On("Get", "test-node-202").Return(&factory.Node{
					Node: factory.NodeFactoryInfo{
						Id:      "test-node-202",
						OrgName: "test-org",
					},
				}, nil)

				// Setup lookup service mock to return an error
				lookupMock.On("GetSystem", "test-org", "messaging").Return(nil, errors.New("messaging system not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("messaging system not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			factoryMock := mbmocks.NewNodeFactoryClient(t)
			lookupMock := mbmocks.NewInitClient(t)
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
