package server

import (
	"context"
	"testing"

	mocks "github.com/ukama/ukama/systems/registry/node/mocks"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/ukama"
)

var testNode = ukama.NewVirtualNodeId("HomeNode")

func TestRegistryServer_GetNode(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}

	nodeRepo.On("Get", testNode).Return(&db.Node{NodeID: testNode.String(),
		State: db.Onboarded, Type: ukama.NODE_ID_TYPE_HOMENODE,
	}, nil).Once()

	s := NewNodeServer(nodeRepo, "")

	node, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{Node: testNode.String()})

	assert.NoError(t, err)
	assert.Equal(t, "onboarded", node.Node.State)
	assert.Equal(t, ukama.NODE_ID_TYPE_HOMENODE, node.Node.Type)
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_UpdateNodeState(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}

	nodeRepo.On("Update", testNode, mock.MatchedBy(func(ns *db.NodeState) bool {
		return *ns == db.Onboarded
	}), (*string)(nil)).Return(nil).Once()
	nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()

	s := NewNodeServer(nodeRepo, "")

	_, err := s.UpdateNodeState(context.TODO(), &pb.UpdateNodeStateRequest{
		Node:  testNode.String(),
		State: "onboarded",
	})

	// Assert
	assert.NoError(t, err)
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_AddNode(t *testing.T) {
	// Arrange
	nodeId := testNode.String()
	nodeRepo := &mocks.NodeRepo{}

	nodeRepo.On("Add", mock.MatchedBy(func(n *db.Node) bool {
		return n.State == db.Onboarded && n.NodeID == nodeId
	})).Return(nil).Once()
	nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()

	s := NewNodeServer(nodeRepo, "")

	// Act
	actNode, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
		Node: &pb.Node{
			Node:  nodeId,
			State: "onboarded",
		},
	})

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, actNode.Node.Name)
	nodeRepo.AssertExpectations(t)
}
