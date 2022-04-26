package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"github.com/ukama/ukamaX/cloud/registry/pb/gen/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

const nodeId = "uk-aa00001-hnode-a1-0001"
const testName = "testName"

func TestRegistry_UpdateNode(t *testing.T) {
	// mock GetNode
	rc := mocks.RegistryServiceClient{}

	rc.On("GetNode", mock.Anything, mock.Anything).Return(&gen.GetNodeResponse{
		Node: &gen.Node{
			NodeId: nodeId,
		},
	}, nil)

	rc.On("UpdateNode", mock.Anything, mock.MatchedBy(func(r *gen.UpdateNodeRequest) bool {
		return r.Name == testName && r.NodeId == nodeId
	})).Return(&gen.UpdateNodeResponse{}, nil)

	r := Registry{
		client: &rc,
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
	rc := mocks.RegistryServiceClient{}

	rc.On("GetNode", mock.Anything, mock.Anything).Return(
		nil, status.Error(codes.NotFound, ""))

	rc.On("AddNode", mock.Anything, mock.MatchedBy(func(r *gen.AddNodeRequest) bool {
		return r.Node.Name == testName && r.Node.NodeId == nodeId
	})).Return(&gen.AddNodeResponse{
		Node: &gen.Node{
			NodeId: nodeId,
			Name:   testName,
		},
	}, nil)

	r := Registry{
		client: &rc,
	}

	resp, isCreated, err := r.AddOrUpdate("org", nodeId, testName)
	if assert.NoError(t, err) {
		rc.AssertExpectations(t)
		assert.Equal(t, testName, resp.Name)
		assert.True(t, isCreated)
	}
}
