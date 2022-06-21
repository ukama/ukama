package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/services/cloud/network/pb/gen"
	"github.com/ukama/ukama/services/cloud/network/pb/gen/mocks"
	pbnode "github.com/ukama/ukama/services/cloud/node/pb/gen"
	ndmock "github.com/ukama/ukama/services/cloud/node/pb/gen/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const nodeId = "uk-aa00001-hnode-a1-0001"
const testName = "testName"

func TestRegistry_UpdateNode(t *testing.T) {
	// mock GetNode
	rc := mocks.NetworkServiceClient{}
	nodeC := ndmock.NodeServiceClient{}

	nodeC.On("GetNode", mock.Anything, mock.Anything).Return(&pbnode.GetNodeResponse{
		Node: &pbnode.Node{
			NodeId: nodeId,
		},
	}, nil)

	nodeC.On("UpdateNode", mock.Anything, mock.MatchedBy(func(r *pbnode.UpdateNodeRequest) bool {
		return r.GetName() == testName && r.NodeId == nodeId
	})).Return(&pbnode.UpdateNodeResponse{}, nil)

	r := Registry{
		client:     &rc,
		nodeClient: &nodeC,
	}

	resp, isCreated, err := r.AddOrUpdate("org", nodeId, testName)
	if assert.NoError(t, err) {
		rc.AssertExpectations(t)
		assert.Equal(t, testName, resp.Name)
		assert.False(t, isCreated)
	}
}

func TestRegistry_AddNode(t *testing.T) {
	// mock GetNode
	rc := mocks.NetworkServiceClient{}
	nodeC := ndmock.NodeServiceClient{}

	nodeC.On("GetNode", mock.Anything, mock.Anything).Return(
		nil, status.Error(codes.NotFound, ""))

	rc.On("AddNode", mock.Anything, mock.MatchedBy(func(r *gen.AddNodeRequest) bool {
		return r.Node.Name == testName && r.Node.NodeId == nodeId
	})).Return(&gen.AddNodeResponse{
		Node: &gen.Node{
			NodeId: nodeId,
			Name:   testName,
		},
	}, nil)

	nodeC.On("AddNode", mock.Anything, mock.MatchedBy(func(r *pbnode.AddNodeRequest) bool {
		return r.Node.Name == testName && r.Node.NodeId == nodeId
	})).Return(&pbnode.AddNodeResponse{
		Node: &pbnode.Node{
			NodeId: nodeId,
			Name:   testName,
		},
	}, nil)

	r := Registry{
		client:     &rc,
		nodeClient: &nodeC,
	}

	resp, isCreated, err := r.AddOrUpdate("org", nodeId, testName)
	if assert.NoError(t, err) {
		rc.AssertExpectations(t)
		assert.Equal(t, testName, resp.Name)
		assert.True(t, isCreated)
	}
}
