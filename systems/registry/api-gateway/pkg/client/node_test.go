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

	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	nodemocks "github.com/ukama/ukama/systems/registry/node/pb/gen/mocks"
)

func TestNewNode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// This test is limited since NewNode creates a real gRPC connection
		// In a real scenario, you might want to use a test server or mock the connection
		nodeHost := "localhost:9090"
		timeout := 5 * time.Second

		// Note: This will fail if there's no actual node service running
		// In practice, you might want to use a test server or skip this test
		node := NewNode(nodeHost, timeout)

		assert.NotNil(t, node)
		assert.Equal(t, nodeHost, node.host)
		assert.Equal(t, timeout, node.timeout)
		assert.NotNil(t, node.client)
		assert.NotNil(t, node.conn)

		// Clean up
		node.Close()
	})
}

func TestNewNodeFromClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		assert.NotNil(t, node)
		assert.Equal(t, "localhost", node.host)
		assert.Equal(t, 1*time.Second, node.timeout)
		assert.Nil(t, node.conn)
		assert.Equal(t, mockClient, node.client)
	})
}

func TestNode_Close(t *testing.T) {
	t.Run("WithoutConnection", func(t *testing.T) {
		node := &Node{
			conn: nil,
		}

		// Should not panic
		assert.NotPanics(t, func() {
			node.Close()
		})
	})
}

func TestNode_AddNode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.AddNodeResponse{
			Node: &pb.Node{
				Id:        "test-node-id",
				Name:      "Test Node",
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
		}

		mockClient.On("AddNode", mock.Anything, &pb.AddNodeRequest{
			NodeId:    "test-node-id",
			Name:      "Test Node",
			Latitude:  40.7128,
			Longitude: -74.0060,
		}).Return(expectedResponse, nil)

		response, err := node.AddNode("test-node-id", "Test Node", "active", 40.7128, -74.0060)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.AlreadyExists, "node already exists")
		mockClient.On("AddNode", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := node.AddNode("existing-node-id", "Existing Node", "active", 40.7128, -74.0060)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_GetNode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.GetNodeResponse{
			Node: &pb.Node{
				Id:        "test-node-id",
				Name:      "Test Node",
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
		}

		mockClient.On("GetNode", mock.Anything, &pb.GetNodeRequest{NodeId: "test-node-id"}).
			Return(expectedResponse, nil)

		response, err := node.GetNode("test-node-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "node not found")
		mockClient.On("GetNode", mock.Anything, &pb.GetNodeRequest{NodeId: "non-existent-id"}).
			Return(nil, expectedError)

		response, err := node.GetNode("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_GetNetworkNodes(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.GetByNetworkResponse{
			Nodes: []*pb.Node{
				{
					Id:   "node-1",
					Name: "Node 1",
				},
				{
					Id:   "node-2",
					Name: "Node 2",
				},
			},
		}

		mockClient.On("GetNodesForNetwork", mock.Anything, &pb.GetByNetworkRequest{NetworkId: "network-1"}).
			Return(expectedResponse, nil)

		response, err := node.GetNetworkNodes("network-1")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "network not found")
		mockClient.On("GetNodesForNetwork", mock.Anything, &pb.GetByNetworkRequest{NetworkId: "non-existent-network"}).
			Return(nil, expectedError)

		response, err := node.GetNetworkNodes("non-existent-network")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_GetSiteNodes(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.GetBySiteResponse{
			Nodes: []*pb.Node{
				{
					Id:   "node-1",
					Name: "Node 1",
				},
				{
					Id:   "node-2",
					Name: "Node 2",
				},
			},
		}

		mockClient.On("GetNodesForSite", mock.Anything, &pb.GetBySiteRequest{SiteId: "site-1"}).
			Return(expectedResponse, nil)

		response, err := node.GetSiteNodes("site-1")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "site not found")
		mockClient.On("GetNodesForSite", mock.Anything, &pb.GetBySiteRequest{SiteId: "non-existent-site"}).
			Return(nil, expectedError)

		response, err := node.GetSiteNodes("non-existent-site")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_GetNodes(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.GetNodesResponse{
			Nodes: []*pb.Node{
				{
					Id:   "node-1",
					Name: "Node 1",
				},
				{
					Id:   "node-2",
					Name: "Node 2",
				},
			},
		}

		mockClient.On("GetNodes", mock.Anything, &pb.GetNodesRequest{}).
			Return(expectedResponse, nil)

		response, err := node.GetNodes()

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.Internal, "internal server error")
		mockClient.On("GetNodes", mock.Anything, &pb.GetNodesRequest{}).
			Return(nil, expectedError)

		response, err := node.GetNodes()

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		request := &pb.ListRequest{
			State:        cpb.NodeState_Operational,
			Connectivity: cpb.NodeConnectivity_Online,
		}

		expectedResponse := &pb.ListResponse{
			Nodes: []*pb.Node{
				{
					Id:   "node-1",
					Name: "Node 1",
				},
			},
		}

		mockClient.On("List", mock.Anything, request).Return(expectedResponse, nil)

		response, err := node.List(request)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		request := &pb.ListRequest{
			State: cpb.NodeState_Unknown,
		}

		expectedError := status.Error(codes.InvalidArgument, "invalid state")
		mockClient.On("List", mock.Anything, request).Return(nil, expectedError)

		response, err := node.List(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_GetNodesByState(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.GetNodesResponse{
			Nodes: []*pb.Node{
				{
					Id:   "node-1",
					Name: "Node 1",
				},
			},
		}

		mockClient.On("GetNodesByState", mock.Anything, mock.Anything).Return(expectedResponse, nil)

		response, err := node.GetNodesByState("online", "active")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.InvalidArgument, "invalid connectivity or state")
		mockClient.On("GetNodesByState", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := node.GetNodesByState("invalid-connectivity", "invalid-state")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_UpdateNodeState(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.UpdateNodeResponse{
			Node: &pb.Node{
				Id: "test-node-id",
			},
		}

		mockClient.On("UpdateNodeState", mock.Anything, &pb.UpdateNodeStateRequest{
			NodeId: "test-node-id",
			State:  "inactive",
		}).Return(expectedResponse, nil)

		response, err := node.UpdateNodeState("test-node-id", "inactive")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "node not found")
		mockClient.On("UpdateNodeState", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := node.UpdateNodeState("non-existent-id", "inactive")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_UpdateNode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.UpdateNodeResponse{
			Node: &pb.Node{
				Id:        "test-node-id",
				Name:      "Updated Node",
				Latitude:  41.7128,
				Longitude: -75.0060,
			},
		}

		mockClient.On("UpdateNode", mock.Anything, &pb.UpdateNodeRequest{
			NodeId:    "test-node-id",
			Name:      "Updated Node",
			Latitude:  41.7128,
			Longitude: -75.0060,
		}).Return(expectedResponse, nil)

		response, err := node.UpdateNode("test-node-id", "Updated Node", 41.7128, -75.0060)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "node not found")
		mockClient.On("UpdateNode", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := node.UpdateNode("non-existent-id", "Updated Node", 41.7128, -75.0060)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_DeleteNode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.DeleteNodeResponse{}

		mockClient.On("DeleteNode", mock.Anything, &pb.DeleteNodeRequest{NodeId: "test-node-id"}).
			Return(expectedResponse, nil)

		response, err := node.DeleteNode("test-node-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "node not found")
		mockClient.On("DeleteNode", mock.Anything, &pb.DeleteNodeRequest{NodeId: "non-existent-id"}).
			Return(nil, expectedError)

		response, err := node.DeleteNode("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_AttachNodes(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.AttachNodesResponse{}

		mockClient.On("AttachNodes", mock.Anything, mock.Anything).Return(expectedResponse, nil)

		response, err := node.AttachNodes("test-node", "left-node", "right-node")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("SuccessWithEmptyNodes", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.AttachNodesResponse{}

		mockClient.On("AttachNodes", mock.Anything, mock.Anything).Return(expectedResponse, nil)

		response, err := node.AttachNodes("test-node", "", "")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "node not found")
		mockClient.On("AttachNodes", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := node.AttachNodes("non-existent-node", "left-node", "right-node")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_DetachNode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.DetachNodeResponse{}

		mockClient.On("DetachNode", mock.Anything, &pb.DetachNodeRequest{NodeId: "test-node-id"}).
			Return(expectedResponse, nil)

		response, err := node.DetachNode("test-node-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "node not found")
		mockClient.On("DetachNode", mock.Anything, &pb.DetachNodeRequest{NodeId: "non-existent-id"}).
			Return(nil, expectedError)

		response, err := node.DetachNode("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_AddNodeToSite(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.AddNodeToSiteResponse{}

		mockClient.On("AddNodeToSite", mock.Anything, &pb.AddNodeToSiteRequest{
			NodeId:    "test-node-id",
			NetworkId: "test-network-id",
			SiteId:    "test-site-id",
		}).Return(expectedResponse, nil)

		response, err := node.AddNodeToSite("test-node-id", "test-network-id", "test-site-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "node or site not found")
		mockClient.On("AddNodeToSite", mock.Anything, mock.Anything).Return(nil, expectedError)

		response, err := node.AddNodeToSite("non-existent-node", "non-existent-network", "non-existent-site")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}

func TestNode_ReleaseNodeFromSite(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedResponse := &pb.ReleaseNodeFromSiteResponse{}

		mockClient.On("ReleaseNodeFromSite", mock.Anything, &pb.ReleaseNodeFromSiteRequest{NodeId: "test-node-id"}).
			Return(expectedResponse, nil)

		response, err := node.ReleaseNodeFromSite("test-node-id")

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := &nodemocks.NodeServiceClient{}
		node := NewNodeFromClient(mockClient)

		expectedError := status.Error(codes.NotFound, "node not found")
		mockClient.On("ReleaseNodeFromSite", mock.Anything, &pb.ReleaseNodeFromSiteRequest{NodeId: "non-existent-id"}).
			Return(nil, expectedError)

		response, err := node.ReleaseNodeFromSite("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, expectedError, err)
		mockClient.AssertExpectations(t)
	})
}
