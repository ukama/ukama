package server

import (
	"context"
	mocks "github.com/ukama/ukama/services/cloud/node/mocks"
	pb "github.com/ukama/ukama/services/cloud/node/pb/gen"
	"github.com/ukama/ukama/services/cloud/node/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/services/common/ukama"
)

var testNodeId = ukama.NewVirtualNodeId("HomeNode")
var testDeviceGatewayHost = "1.1.1.1"

const testOrgName = "org-1"
const testNetName = "net-1"
const testNetId = 98

type qPubStub struct {
}

func (q qPubStub) Publish(payload any, routingKey string) error {
	return nil
}
func TestRegistryServer_GetNode(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	pub := &qPubStub{}

	nodeRepo.On("Get", testNodeId).Return(&db.Node{NodeID: testNodeId.String(),
		State: db.Pending, Type: db.NodeTypeHome,
	}, nil).Once()

	s := NewNodeServer(nodeRepo, pub)
	node, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{NodeId: testNodeId.String()})
	assert.NoError(t, err)
	assert.Equal(t, pb.NodeState_PENDING, node.Node.State)
	assert.Equal(t, pb.NodeType_HOME, node.Node.Type)
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_UpdateNodeState(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}

	pub := qPubStub{}

	nodeRepo.On("Update", testNodeId, mock.MatchedBy(func(ns *db.NodeState) bool {
		return *ns == db.Onboarded
	}), (*string)(nil)).Return(nil).Once()
	s := NewNodeServer(nodeRepo, pub)
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
	pub := &qPubStub{}

	nodeRepo.On("Add", mock.MatchedBy(func(n *db.Node) bool {
		return n.State == db.Pending && n.NodeID == nodeId
	})).Return(nil).Once()
	s := NewNodeServer(nodeRepo, pub)

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
