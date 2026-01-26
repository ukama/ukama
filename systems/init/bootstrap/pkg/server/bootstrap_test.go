/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package server

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	mbmocks "github.com/ukama/ukama/systems/common/mocks"
// 	"github.com/ukama/ukama/systems/common/rest/client/factory"
// 	"github.com/ukama/ukama/systems/init/bootstrap/mocks"
// 	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
// 	lpb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
// 	lmocks "github.com/ukama/ukama/systems/init/lookup/pb/gen/mocks"
// )

// const (
// 	testOrgName     = "test-org"
// 	testNodeID123   = "test-node-123"
// 	testNodeID456   = "test-node-456"
// 	testNodeID789   = "test-node-789"
// 	testNodeID101   = "test-node-101"
// 	testNodeID202   = "test-node-202"
// 	testNodeID303   = "test-node-303"
// 	unknownOrgName  = "unknown-org"
// )

// func TestGetNodeCredentials(t *testing.T) {
// 	 tests := []struct {
// 		 name           string
// 		 nodeID         string
// 		 setupMocks     func(*mbmocks.NodeFactoryClient, *mocks.LookupClientProvider, *mbmocks.MsgBusServiceClient)
// 		 expectedResult *pb.GetNodeCredentialsResponse
// 		 expectedError  error
// 	 }{
// 		{
// 			name:   "Success - Node with org name and messaging system",
// 			nodeID: testNodeID123,
// 			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
// 				// Setup factory mock to return a node with org name
// 				factoryMock.On("Get", testNodeID123).Return(&factory.Node{
// 					Node: factory.NodeFactoryInfo{
// 						Id:      testNodeID123,
// 						OrgName: testOrgName,
// 					},
// 				}, nil)

// 				 // Setup lookup client provider mock
// 				 lookupClientMock := &lmocks.LookupServiceClient{}
// 				 lookupMock.On("GetClient").Return(lookupClientMock, nil)

// 				 // Setup lookup service mock to return messaging system
// 				 lookupClientMock.On("GetSystemForOrg", mock.Anything, &lpb.GetSystemRequest{
// 					 OrgName:    "test-org",
// 					 SystemName: "messaging",
// 				 }).Return(&lpb.GetSystemResponse{
// 					 ApiGwIp:     "192.168.1.100",
// 					 Certificate: "test-certificate-data",
// 				 }, nil)
// 			 },
// 			 expectedResult: &pb.GetNodeCredentialsResponse{
// 				 Id:          "test-node-123",
// 				 OrgName:     "test-org",
// 				 Ip:          "127.0.0.1",
// 				 Certificate: "test-certificate-data",
// 			 },
// 			 expectedError: nil,
// 		 },
// 		{
// 			name:   "Error - Node without org name",
// 			nodeID: testNodeID456,
// 			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
// 				// Setup factory mock to return a node without org name
// 				factoryMock.On("Get", testNodeID456).Return(&factory.Node{
// 					Node: factory.NodeFactoryInfo{
// 						Id:      testNodeID456,
// 						OrgName: "",
// 					},
// 				}, nil)
// 			},
// 			expectedResult: nil,
// 			expectedError:  errors.New("rpc error: code = FailedPrecondition desc = Node is not provisioned in any org"),
// 		},
// 		{
// 			name:   "Error - Factory client fails",
// 			nodeID: testNodeID789,
// 			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
// 				// Setup factory mock to return an error
// 				factoryMock.On("Get", testNodeID789).Return(nil, errors.New("factory service unavailable"))
// 			 },
// 			 expectedResult: nil,
// 			 expectedError:  errors.New("factory service unavailable"),
// 		 },
// 		{
// 			name:   "Error - Lookup client provider fails",
// 			nodeID: testNodeID101,
// 			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
// 				// Setup factory mock to return a node with org name
// 				factoryMock.On("Get", testNodeID101).Return(&factory.Node{
// 					Node: factory.NodeFactoryInfo{
// 						Id:      testNodeID101,
// 						OrgName: testOrgName,
// 					},
// 				}, nil)

// 				 // Setup lookup client provider mock to return an error
// 				 lookupMock.On("GetClient").Return(nil, errors.New("lookup service connection failed"))
// 			 },
// 			 expectedResult: nil,
// 			 expectedError:  errors.New("lookup service connection failed"),
// 		 },
// 		{
// 			name:   "Error - Lookup service fails",
// 			nodeID: testNodeID202,
// 			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
// 				// Setup factory mock to return a node with org name
// 				factoryMock.On("Get", testNodeID202).Return(&factory.Node{
// 					Node: factory.NodeFactoryInfo{
// 						Id:      testNodeID202,
// 						OrgName: testOrgName,
// 					},
// 				}, nil)

// 				// Setup lookup client provider mock
// 				lookupClientMock := &lmocks.LookupServiceClient{}
// 				lookupMock.On("GetClient").Return(lookupClientMock, nil)

// 				// Setup lookup service mock to return an error
// 				lookupClientMock.On("GetSystemForOrg", mock.Anything, &lpb.GetSystemRequest{
// 					OrgName:    testOrgName,
// 					SystemName: "messaging",
// 				}).Return(nil, errors.New("messaging system not found"))
// 			},
// 			expectedResult: nil,
// 			expectedError:  errors.New("messaging system not found"),
// 		},
// 		{
// 			name:   "Error - DNS not found in map",
// 			nodeID: testNodeID303,
// 			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
// 				// Setup factory mock to return a node with org name not in DNS map
// 				factoryMock.On("Get", testNodeID303).Return(&factory.Node{
// 					Node: factory.NodeFactoryInfo{
// 						Id:      testNodeID303,
// 						OrgName: unknownOrgName,
// 					},
// 				}, nil)
// 			},
// 			expectedResult: nil,
// 			expectedError:  errors.New("rpc error: code = NotFound desc = DNS is not found for org " + unknownOrgName),
// 		},
// 	 }

// 	 for _, tt := range tests {
// 		 t.Run(tt.name, func(t *testing.T) {
// 			 // Create mocks
// 			 factoryMock := mbmocks.NewNodeFactoryClient(t)
// 			 lookupMock := mocks.NewLookupClientProvider(t)
// 			 msgBusMock := mbmocks.NewMsgBusServiceClient(t)

// 			// Setup mocks
// 			tt.setupMocks(factoryMock, lookupMock, msgBusMock)

// 			// Create server instance with DNS map - use localhost which will always resolve
// 			dnsMap := map[string]string{testOrgName: "localhost"}
// 			server := NewBootstrapServer(msgBusMock, false, lookupMock, factoryMock, dnsMap, config, nil)

// 			// Create request
// 			req := &pb.GetNodeCredentialsRequest{
// 				Id: tt.nodeID,
// 			}

// 			// Call the method
// 			result, err := server.GetNodeCredentials(context.Background(), req)

// 			// Assertions
// 			if tt.expectedError != nil {
// 				assert.Error(t, err)
// 				assert.Equal(t, tt.expectedError.Error(), err.Error())
// 				assert.Nil(t, result)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotNil(t, result)
// 				assert.Equal(t, tt.expectedResult.Id, result.Id)
// 				assert.Equal(t, tt.expectedResult.OrgName, result.OrgName)
// 				assert.Equal(t, tt.expectedResult.Ip, result.Ip)
// 				assert.Equal(t, tt.expectedResult.Certificate, result.Certificate)
// 			}

// 			// Verify all mocks were called as expected
// 			factoryMock.AssertExpectations(t)
// 			lookupMock.AssertExpectations(t)
// 			// msgBusMock is not used in GetNodeCredentials, so we don't assert expectations
// 		 })
// 	 }
//  }