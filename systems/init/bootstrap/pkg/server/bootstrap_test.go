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
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/rest/client/factory"
	"github.com/ukama/ukama/systems/init/bootstrap/mocks"
	pb "github.com/ukama/ukama/systems/init/bootstrap/pb/gen"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg"
	"github.com/ukama/ukama/systems/init/bootstrap/pkg/mesh"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	testOrgName    = "test-org"
	testNodeID123  = "test-node-123"
	testNodeID456  = "test-node-456"
	testNodeID789  = "test-node-789"
	testNodeID303  = "test-node-303"
	testNodeID404  = "test-node-404"
	unknownOrgName = "unknown-org"
	testMeshNs     = testOrgName + "-messaging"
)

var testConfig = &pkg.Config{
	OrgName:       testOrgName,
	MeshNamespace: "messaging",
	DNSMap:        []pkg.OrgDNS{{OrgName: testOrgName, DNS: "localhost"}},
}

func meshTemplate() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "ukama-mesh-svc", Namespace: testMeshNs, Labels: map[string]string{mesh.NameLabel: mesh.MeshName}},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "mesh", Image: "mesh:test"}}}},
		},
	}
}

func readyMesh(nodeId string) []runtime.Object {
	name := testOrgName + "-" + mesh.Component + "-" + nodeId
	labels := map[string]string{mesh.NameLabel: mesh.MeshName, mesh.ComponentLabel: mesh.Component, mesh.InstanceIdLabel: nodeId}
	return []runtime.Object{
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: name + "-w76pt", Namespace: testMeshNs, Labels: labels},
			Status: corev1.PodStatus{
				Phase:      corev1.PodRunning,
				PodIP:      "10.0.110.252",
				Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}},
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testMeshNs, Labels: labels},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{mesh.ComponentLabel: mesh.Component, mesh.InstanceIdLabel: nodeId},
				Ports: []corev1.ServicePort{
					{Name: "http", Protocol: corev1.ProtocolTCP, Port: mesh.ServicePort, TargetPort: intstr.FromInt32(mesh.ServicePort)},
					{Name: "websocket", Protocol: corev1.ProtocolTCP, Port: mesh.WebsocketPort, TargetPort: intstr.FromInt32(mesh.WebsocketPort)},
				},
			},
		},
	}
}

func TestGetNodeCredentials(t *testing.T) {
	tests := []struct {
		name           string
		nodeID         string
		meshObjects    []runtime.Object
		setupMocks     func(*mbmocks.NodeFactoryClient, *mocks.LookupClientProvider, *mbmocks.MsgBusServiceClient)
		expectedResult *pb.GetNodeCredentialsResponse
		expectedError  error
	}{
		{
			name:        "Success - Node with org name, DNS and a ready mesh",
			nodeID:      testNodeID123,
			meshObjects: readyMesh(testNodeID123),
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
				factoryMock.On("Get", testNodeID123).Return(&factory.Node{
					Node: factory.NodeFactoryInfo{
						Id:      testNodeID123,
						OrgName: testOrgName,
					},
				}, nil)
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
			name:   "Error - Missing mesh template returns Internal",
			nodeID: testNodeID404,
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
				factoryMock.On("Get", testNodeID404).Return(&factory.Node{
					Node: factory.NodeFactoryInfo{
						Id:      testNodeID404,
						OrgName: testOrgName,
					},
				}, nil)
			},
			expectedResult: nil,
			expectedError:  errors.New("rpc error: code = Internal desc = no mesh deployment found in namespace " + testMeshNs),
		},
		{
			name:   "Error - Node without org name",
			nodeID: testNodeID456,
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
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
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
				factoryMock.On("Get", testNodeID789).Return(nil, errors.New("factory service unavailable"))
			},
			expectedResult: nil,
			expectedError:  errors.New("factory service unavailable"),
		},
		{
			name:   "Error - DNS not found in map",
			nodeID: testNodeID303,
			setupMocks: func(factoryMock *mbmocks.NodeFactoryClient, lookupMock *mocks.LookupClientProvider, msgBusMock *mbmocks.MsgBusServiceClient) {
				factoryMock.On("Get", testNodeID303).Return(&factory.Node{
					Node: factory.NodeFactoryInfo{
						Id:      testNodeID303,
						OrgName: unknownOrgName,
					},
				}, nil)
			},
			expectedResult: nil,
			expectedError:  errors.New("rpc error: code = NotFound desc = DNS is not found for org " + unknownOrgName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factoryMock := mbmocks.NewNodeFactoryClient(t)
			lookupMock := mocks.NewLookupClientProvider(t)
			msgBusMock := mbmocks.NewMsgBusServiceClient(t)

			tt.setupMocks(factoryMock, lookupMock, msgBusMock)

			dnsMap := map[string]string{testOrgName: "localhost"}
			serverConfig := &BootstrapServerConfig{Config: testConfig, MessagingCert: "test-certificate-data"}
			testDeps := &BootstrapTestDeps{ClientSet: fake.NewSimpleClientset(tt.meshObjects...), DNSMap: dnsMap}
			server := NewBootstrapServerWithDeps(msgBusMock, false, lookupMock, factoryMock, serverConfig, testDeps)

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
		})
	}
}

func TestGetNodeCredentials_OutsideKubernetesSkipsMesh(t *testing.T) {
	factoryMock := mbmocks.NewNodeFactoryClient(t)
	factoryMock.On("Get", testNodeID123).Return(&factory.Node{
		Node: factory.NodeFactoryInfo{Id: testNodeID123, OrgName: testOrgName},
	}, nil)

	serverConfig := &BootstrapServerConfig{Config: testConfig, MessagingCert: "test-certificate-data"}
	server := NewBootstrapServerWithDeps(mbmocks.NewMsgBusServiceClient(t), false, mocks.NewLookupClientProvider(t), factoryMock, serverConfig, nil)

	result, err := server.GetNodeCredentials(context.Background(), &pb.GetNodeCredentialsRequest{Id: testNodeID123})

	require.NoError(t, err)
	assert.Equal(t, testNodeID123, result.Id)
	assert.Equal(t, "127.0.0.1", result.Ip)
}

func TestGetNodeCredentials_NewNodeIsUnavailableUntilMeshIsReady(t *testing.T) {
	factoryMock := mbmocks.NewNodeFactoryClient(t)
	factoryMock.On("Get", testNodeID123).Return(&factory.Node{
		Node: factory.NodeFactoryInfo{Id: testNodeID123, OrgName: testOrgName},
	}, nil)

	cs := fake.NewSimpleClientset(meshTemplate())
	serverConfig := &BootstrapServerConfig{Config: testConfig, MessagingCert: "test-certificate-data"}
	testDeps := &BootstrapTestDeps{ClientSet: cs, DNSMap: map[string]string{testOrgName: "localhost"}}
	server := NewBootstrapServerWithDeps(mbmocks.NewMsgBusServiceClient(t), false, mocks.NewLookupClientProvider(t), factoryMock, serverConfig, testDeps)

	result, err := server.GetNodeCredentials(context.Background(), &pb.GetNodeCredentialsRequest{Id: testNodeID123})

	assert.Nil(t, result)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unavailable, st.Code())
	assert.Contains(t, st.Message(), "is not ready (phase Pending)")

	pods, err := cs.CoreV1().Pods(testMeshNs).List(context.Background(), metav1.ListOptions{})
	require.NoError(t, err)
	assert.Len(t, pods.Items, 1, "the first request creates the mesh pod")
}
