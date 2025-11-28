/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	nnmocks "github.com/ukama/ukama/systems/messaging/nns/pb/gen/mocks"
)

// Test data constants
const (
	testNodeId       = "uk-sa3333-uk-0001-0001"
	testNodeIp       = "127.0.0.100"
	testMeshIp       = "127.0.0.200"
	testNodePort     = int32(8080)
	testMeshPort     = int32(9090)
	testNetwork      = "test-network"
	testSite         = "test-site"
	testMeshHostName = "mesh-hostname"
)

func TestNewNnsFromClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		assert.NotNil(t, nns)
		assert.Equal(t, "localhost", nns.host)
		assert.Equal(t, 1*time.Second, nns.timeout)
		assert.Nil(t, nns.conn)
		assert.Equal(t, mockClient, nns.client)
	})
}

func TestNns_Close(t *testing.T) {
	t.Run("WithoutConnection", func(t *testing.T) {
		nns := &Nns{
			conn: nil,
		}

		// Should not panic
		assert.NotPanics(t, func() {
			nns.Close()
		})
	})
}

func TestNns_GetNodeRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.GetNodeRequest{
			NodeId: testNodeId,
		}

		expectedResponse := &pb.GetNodeResponse{
			NodeId:   testNodeId,
			NodeIp:   testNodeIp,
			NodePort: testNodePort,
		}

		mockClient.On("GetNode", mock.Anything, req, mock.Anything).
			Return(expectedResponse, nil)

		response, err := nns.GetNodeRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.GetNodeRequest{
			NodeId: testNodeId,
		}

		expectedError := status.Error(codes.NotFound, "node not found")
		mockClient.On("GetNode", mock.Anything, req, mock.Anything).
			Return(nil, expectedError)

		response, err := nns.GetNodeRequest(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNns_GetMeshRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.GetMeshRequest{}

		expectedResponse := &pb.GetMeshResponse{
			MeshIp:   testMeshIp,
			MeshPort: testMeshPort,
		}

		mockClient.On("GetMesh", mock.Anything, req, mock.Anything).
			Return(expectedResponse, nil)

		response, err := nns.GetMeshRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.GetMeshRequest{}

		expectedError := status.Error(codes.Internal, "internal server error")
		mockClient.On("GetMesh", mock.Anything, req, mock.Anything).
			Return(nil, expectedError)

		response, err := nns.GetMeshRequest(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNns_SetRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.SetRequest{
			NodeId:       testNodeId,
			NodeIp:       testNodeIp,
			NodePort:     testNodePort,
			MeshIp:       testMeshIp,
			MeshPort:     testMeshPort,
			Network:      testNetwork,
			Site:         testSite,
			MeshHostName: testMeshHostName,
		}

		expectedResponse := &pb.SetResponse{}

		mockClient.On("Set", mock.Anything, req, mock.Anything).
			Return(expectedResponse, nil)

		response, err := nns.SetRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.SetRequest{
			NodeId:       testNodeId,
			NodeIp:       testNodeIp,
			NodePort:     testNodePort,
			MeshIp:       testMeshIp,
			MeshPort:     testMeshPort,
			Network:      testNetwork,
			Site:         testSite,
			MeshHostName: testMeshHostName,
		}

		expectedError := status.Error(codes.AlreadyExists, "node already exists")
		mockClient.On("Set", mock.Anything, req, mock.Anything).
			Return(nil, expectedError)

		response, err := nns.SetRequest(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNns_UpdateMeshRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.UpdateMeshRequest{
			MeshIp:   testMeshIp,
			MeshPort: testMeshPort,
		}

		expectedResponse := &pb.UpdateMeshResponse{}

		mockClient.On("UpdateMesh", mock.Anything, req, mock.Anything).
			Return(expectedResponse, nil)

		response, err := nns.UpdateMeshRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.UpdateMeshRequest{
			MeshIp:   testMeshIp,
			MeshPort: testMeshPort,
		}

		expectedError := status.Error(codes.Internal, "failed to update mesh")
		mockClient.On("UpdateMesh", mock.Anything, req, mock.Anything).
			Return(nil, expectedError)

		response, err := nns.UpdateMeshRequest(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNns_UpdateNodeRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.UpdateNodeRequest{
			NodeId:   testNodeId,
			NodeIp:   testNodeIp,
			NodePort: testNodePort,
		}

		expectedResponse := &pb.UpdateNodeResponse{}

		mockClient.On("UpdateNode", mock.Anything, req, mock.Anything).
			Return(expectedResponse, nil)

		response, err := nns.UpdateNodeRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.UpdateNodeRequest{
			NodeId:   testNodeId,
			NodeIp:   testNodeIp,
			NodePort: testNodePort,
		}

		expectedError := status.Error(codes.NotFound, "node not found")
		mockClient.On("UpdateNode", mock.Anything, req, mock.Anything).
			Return(nil, expectedError)

		response, err := nns.UpdateNodeRequest(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNns_DeleteRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.DeleteRequest{
			NodeId: testNodeId,
		}

		expectedResponse := &pb.DeleteResponse{}

		mockClient.On("Delete", mock.Anything, req, mock.Anything).
			Return(expectedResponse, nil)

		response, err := nns.DeleteRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.DeleteRequest{
			NodeId: testNodeId,
		}

		expectedError := status.Error(codes.NotFound, "node not found")
		mockClient.On("Delete", mock.Anything, req, mock.Anything).
			Return(nil, expectedError)

		response, err := nns.DeleteRequest(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNns_ListRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.ListRequest{}

		expectedResponse := &pb.ListResponse{
			List: []*pb.OrgMap{
				{
					NodeId:       testNodeId,
					NodeIp:       testNodeIp,
					NodePort:     testNodePort,
					MeshIp:       testMeshIp,
					MeshPort:     testMeshPort,
					Network:      testNetwork,
					Site:         testSite,
					MeshHostName: testMeshHostName,
				},
			},
		}

		mockClient.On("List", mock.Anything, req, mock.Anything).
			Return(expectedResponse, nil)

		response, err := nns.ListRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		assert.Len(t, response.List, 1)
		assert.Equal(t, testNodeId, response.List[0].NodeId)
		mockClient.AssertExpectations(t)
	})

	t.Run("SuccessWithEmptyList", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.ListRequest{}

		expectedResponse := &pb.ListResponse{
			List: []*pb.OrgMap{},
		}

		mockClient.On("List", mock.Anything, req, mock.Anything).
			Return(expectedResponse, nil)

		response, err := nns.ListRequest(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		assert.Len(t, response.List, 0)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nnmocks.NnsClient{}
		nns := NewNnsFromClient(mockClient)

		req := &pb.ListRequest{}

		expectedError := status.Error(codes.Internal, "internal server error")
		mockClient.On("List", mock.Anything, req, mock.Anything).
			Return(nil, expectedError)

		response, err := nns.ListRequest(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}
