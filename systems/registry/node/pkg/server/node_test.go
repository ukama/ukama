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

var testNodeId = ukama.NewVirtualNodeId("HomeNode")

func TestRegistryServer_GetNode(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}

	nodeRepo.On("Get", testNodeId).Return(&db.Node{NodeID: testNodeId.String(),
		State: db.Pending, Type: db.NodeTypeHome,
	}, nil).Once()

	s := NewNodeServer(nodeRepo, "")

	node, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{NodeId: testNodeId.String()})

	assert.NoError(t, err)
	assert.Equal(t, pb.NodeState_PENDING, node.Node.State)
	assert.Equal(t, pb.NodeType_HOME, node.Node.Type)
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_UpdateNodeState(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}

	nodeRepo.On("Update", testNodeId, mock.MatchedBy(func(ns *db.NodeState) bool {
		return *ns == db.Onboarded
	}), (*string)(nil)).Return(nil).Once()
	nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()

	s := NewNodeServer(nodeRepo, "")

	_, err := s.UpdateNodeState(context.TODO(), &pb.UpdateNodeStateRequest{
		NodeId: testNodeId.String(),
		State:  pb.NodeState_ONBOARDED,
	})

	// Assert
	assert.NoError(t, err)
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_AddNode(t *testing.T) {
	// Arrange
	nodeId := testNodeId.String()
	nodeRepo := &mocks.NodeRepo{}

	nodeRepo.On("Add", mock.MatchedBy(func(n *db.Node) bool {
		return n.State == db.Pending && n.NodeID == nodeId
	})).Return(nil).Once()
	nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()

	s := NewNodeServer(nodeRepo, "")

	// Act
	actNode, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
		Node: &pb.Node{
			NodeId: nodeId,
			State:  pb.NodeState_PENDING,
		},
	})

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, actNode.Node.Name)
	nodeRepo.AssertExpectations(t)
}

func Test_toDbNodeType(t *testing.T) {
	tests := []struct {
		nodeId ukama.NodeID
		want   db.NodeType
	}{
		{
			nodeId: ukama.NewVirtualHomeNodeId(),
			want:   db.NodeTypeHome,
		},
		{
			nodeId: ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE),
			want:   db.NodeTypeTower,
		},
		{
			nodeId: ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_AMPNODE),
			want:   db.NodeTypeAmplifier,
		},
		{
			nodeId: ukama.NewVirtualNodeId("unknown"),
			want:   db.NodeTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.nodeId.String(), func(t *testing.T) {

			got := toDbNodeType(tt.nodeId.GetNodeType())
			assert.Equal(t, tt.want, got)
		})
	}
}
