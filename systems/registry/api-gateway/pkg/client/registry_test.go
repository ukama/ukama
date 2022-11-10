package client

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/registry/network/pb/gen"
	"github.com/ukama/ukama/systems/registry/network/pb/gen/mocks"
	pbnode "github.com/ukama/ukama/systems/registry/node/pb/gen"
	ndmock "github.com/ukama/ukama/systems/registry/node/pb/gen/mocks"
)

const TEST_NODE_ID = "uk-aa00001-hnode-a1-0001"
const testName = "testName"

func TestRegistry_UpdateNode(t *testing.T) {
	// mock GetNode
	rc := &mocks.NetworkServiceClient{}
	nodeC := &ndmock.NodeServiceClient{}
	nodeType := pbnode.NodeType_HOME

	setMocks := func() {
		rc = &mocks.NetworkServiceClient{}
		nodeC = &ndmock.NodeServiceClient{}
		nodeC.On("GetNode", mock.Anything, mock.Anything).Return(&pbnode.GetNodeResponse{
			Node: &pbnode.Node{
				NodeId: TEST_NODE_ID,
				Name:   testName,
				Type:   nodeType,
			},
		}, nil)

		nodeC.On("UpdateNode", mock.Anything, mock.MatchedBy(func(r *pbnode.UpdateNodeRequest) bool {
			return r.GetName() == testName && r.NodeId == TEST_NODE_ID
		})).Return(&pbnode.UpdateNodeResponse{
			Node: &pbnode.Node{
				Type: nodeType,
			},
		}, nil)
	}

	t.Run("updateNode", func(t *testing.T) {
		setMocks()
		r := Registry{
			client:     rc,
			nodeClient: nodeC,
		}

		resp, err := r.UpdateNode("org", TEST_NODE_ID, testName)
		if assert.NoError(t, err) {
			rc.AssertExpectations(t)
			assert.Equal(t, testName, resp.Name)
		}
	})

	t.Run("updateNodeAttachFailOnNodeType", func(t *testing.T) {
		setMocks()
		r := Registry{
			client:     rc,
			nodeClient: nodeC,
		}

		_, err := r.UpdateNode("org", TEST_NODE_ID, testName, "test")
		assert.ErrorContains(t, err, "node type")
	})

	t.Run("updateNodeAttachingNodes", func(t *testing.T) {
		nodeType = pbnode.NodeType_TOWER
		setMocks()
		r := Registry{
			client:     rc,
			nodeClient: nodeC,
		}

		toAttach := []string{ukama.NewVirtualAmplifierNodeId().String(), ukama.NewVirtualAmplifierNodeId().String()}
		nodeC.On("AttachNodes", mock.Anything, mock.MatchedBy(func(r *pbnode.AttachNodesRequest) bool {
			return r.ParentNodeId == TEST_NODE_ID && strings.Join(r.GetAttachedNodeIds(), "") == strings.Join(toAttach, "")
		})).Return(&pbnode.AttachNodesResponse{}, nil)

		resp, err := r.UpdateNode("org", TEST_NODE_ID, testName, toAttach...)
		if assert.NoError(t, err) {
			rc.AssertExpectations(t)
			assert.Equal(t, testName, resp.Name)
		}
	})
}

func TestRegistry_AddNode(t *testing.T) {
	var rc *mocks.NetworkServiceClient
	var nodeC *ndmock.NodeServiceClient
	nodeId := "uk-aa00001-hnode-a1-0001"
	nodeType := pbnode.NodeType_HOME

	setMocks := func() {
		rc = &mocks.NetworkServiceClient{}
		nodeC = &ndmock.NodeServiceClient{}
		nodeC.On("GetNode", mock.Anything, mock.Anything).Return(
			&pbnode.GetNodeResponse{
				Node: &pbnode.Node{
					Name: testName,
				},
			}, nil)

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
				Type:   nodeType,
			},
		}, nil)
	}

	t.Run("AddNode", func(t *testing.T) {
		setMocks()
		r := Registry{
			client:     rc,
			nodeClient: nodeC,
		}

		resp, err := r.Add("org", nodeId, testName)
		if assert.NoError(t, err) {
			rc.AssertExpectations(t)
			assert.Equal(t, testName, resp.Name)
		}
	})

	t.Run("AddNodeWithAttachmentsFailedDueType", func(t *testing.T) {
		setMocks()
		r := Registry{
			client:     rc,
			nodeClient: nodeC,
		}

		toAttach := []string{ukama.NewVirtualAmplifierNodeId().String(), ukama.NewVirtualAmplifierNodeId().String()}
		nodeC.On("AttachNodes", mock.Anything, mock.MatchedBy(func(r *pbnode.AttachNodesRequest) bool {
			return r.ParentNodeId == nodeId && strings.Join(r.GetAttachedNodeIds(), "") == strings.Join(toAttach, "")
		})).Return(&pbnode.AttachNodesResponse{}, nil)

		_, err := r.Add("org", nodeId, testName, toAttach...)
		assert.ErrorContains(t, err, "node type")
	})

	t.Run("AddNodeWithAttachments", func(t *testing.T) {
		nodeId = ukama.NewVirtualTowerNodeId().String()
		nodeType = pbnode.NodeType_TOWER
		setMocks()
		r := Registry{
			client:     rc,
			nodeClient: nodeC,
		}

		toAttach := []string{ukama.NewVirtualAmplifierNodeId().String(), ukama.NewVirtualAmplifierNodeId().String()}
		nodeC.On("AttachNodes", mock.Anything, mock.MatchedBy(func(r *pbnode.AttachNodesRequest) bool {
			return r.ParentNodeId == nodeId && strings.Join(r.GetAttachedNodeIds(), "") == strings.Join(toAttach, "")
		})).Return(&pbnode.AttachNodesResponse{}, nil)

		_, err := r.Add("org", nodeId, testName, toAttach...)
		if assert.NoError(t, err) {
			rc.AssertExpectations(t)
			nodeC.AssertExpectations(t)
		}
	})

}
