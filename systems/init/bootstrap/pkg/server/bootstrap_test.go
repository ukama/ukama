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
	messaging "github.com/ukama/ukama/systems/common/rest/client/messaging"
	"github.com/ukama/ukama/systems/init/bootstrap/mocks"
	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	testOrgName    = "test-org"
	testNodeID123  = "test-node-123"
	testNodeID456  = "test-node-456"
	testNodeID789  = "test-node-789"
	testNodeID303  = "test-node-303"
	unknownOrgName = "unknown-org"
)

var testConfig = &pkg.Config{
	OrgName:       testOrgName,
	MeshNamespace: "messaging",
	DNSMap:        []pkg.OrgDNS{{OrgName: testOrgName, DNS: "localhost"}},
}

func TestGetNodeCredentials(t *testing.T) {
	tests := []struct {
		name           string
		nodeID         string
		setupMocks     func(*mbmocks.NodeFactoryClient, *mocks.LookupClientProvider, *mbmocks.MsgBusServiceClient, *mbmocks.NnsClient)
		expectedResult *pb.GetNodeCredentialsResponse
		expectedError  error
	}{
		{
			name:   "Success - Node with org name and DNS resolves",
			nodeID: testNodeID123,
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient, nnsMock *mbmocks.NnsClient) {
				factoryMock.On("Get", testNodeID123).Return(&factory.Node{
					Node: factory.NodeFactoryInfo{
						Id:      testNodeID123,
						OrgName: testOrgName,
					},
				}, nil)
				nnsMock.On("GetMesh", testNodeID123).Return((*messaging.MeshInfo)(nil), nil)
			},
			expectedResult: &pb.GetNodeCredentialsResponse{
				Id:          "test-node-123",
				OrgName:     "test-org",
				Ip:          "127.0.0.1",
				Certificate: "test-certificate-data",
			},
			expectedError: nil,
		},
		{
			name:   "Error - Node without org name",
			nodeID: testNodeID456,
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient, nnsMock *mbmocks.NnsClient) {
				factoryMock.On("Get", testNodeID456).Return(&factory.Node{
					Node: factory.NodeFactoryInfo{
						Id:      testNodeID456,
						OrgName: "",
					},
				}, nil)
			},
			expectedResult: nil,
			expectedError:  errors.New("rpc error: code = FailedPrecondition desc = Node is not provisioned in any org"),
		},
		{
			name:   "Error - Factory client fails",
			nodeID: testNodeID789,
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient, nnsMock *mbmocks.NnsClient) {
				factoryMock.On("Get", testNodeID789).Return(nil, errors.New("factory service unavailable"))
			},
			expectedResult: nil,
			expectedError:  errors.New("factory service unavailable"),
		},
		{
			name:   "Error - DNS not found in map",
			nodeID: testNodeID303,
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient, nnsMock *mbmocks.NnsClient) {
				factoryMock.On("Get", testNodeID303).Return(&factory.Node{
					Node: factory.NodeFactoryInfo{
						Id:      testNodeID303,
						OrgName: unknownOrgName,
					},
				}, nil)
			},
			expectedResult: nil,
			expectedError:  errors.New("rpc error: code = NotFound desc = DNS is not found for org "+unknownOrgName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factoryMock := mbmocks.NewNodeFactoryClient(t)
			lookupMock := mocks.NewLookupClientProvider(t)
			msgBusMock := mbmocks.NewMsgBusServiceClient(t)
			nnsMock := mbmocks.NewNnsClient(t)

			tt.setupMocks(factoryMock, lookupMock, msgBusMock, nnsMock)

			dnsMap := map[string]string{testOrgName: "localhost"}
			serverConfig := &BootstrapServerConfig{Config: testConfig, MessagingCert: "test-certificate-data"}
			testDeps := &BootstrapTestDeps{ClientSet: fake.NewSimpleClientset(), DNSMap: dnsMap}
			server := NewBootstrapServerWithDeps(msgBusMock, false, lookupMock, factoryMock, nnsMock, serverConfig, testDeps)

			req := &pb.GetNodeCredentialsRequest{
				Id: tt.nodeID,
			}

			result, err := server.GetNodeCredentials(context.Background(), req)

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

			factoryMock.AssertExpectations(t)
			nnsMock.AssertExpectations(t)
		})
	}
}