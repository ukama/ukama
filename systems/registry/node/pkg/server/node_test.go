package server_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/mocks"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"github.com/ukama/ukama/systems/registry/node/pkg/server"
	"gorm.io/gorm"

	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	opb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	omocks "github.com/ukama/ukama/systems/registry/org/pb/gen/mocks"
)

var testNode = ukama.NewVirtualNodeId("HomeNode")

// func TestRegistryServer_UpdateNodeState(t *testing.T) {
// nodeRepo := &mocks.NodeRepo{}

// nodeRepo.On("Update", testNode, mock.MatchedBy(func(ns *db.NodeState) bool {
// return *ns == db.Onboarded
// }), (*string)(nil)).Return(nil).Once()
// nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()

// s := NewNodeServer(nodeRepo, "")

// _, err := s.UpdateNodeState(context.TODO(), &pb.UpdateNodeStateRequest{
// Node:  testNode.String(),
// State: "onboarded",
// })

// // Assert
// assert.NoError(t, err)
// nodeRepo.AssertExpectations(t)
// }

func TestNodeServer_Add(t *testing.T) {
	nodeId := testNode.String()
	nodeState := "online"

	nodeRepo := &mocks.NodeRepo{}
	orgService := &mocks.OrgClientProvider{}
	networkService := &mocks.NetworkClientProvider{}

	const nodeName = "node-A"
	const nodeType = "hnode"
	var orgId = uuid.NewV4()

	s := server.NewNodeServer(nodeRepo, nil, "", orgService, networkService)

	node := &db.Node{
		Id:    nodeId,
		Name:  nodeName,
		OrgId: orgId,
		Type:  testNode.GetNodeType(),
		State: db.Online,
	}

	nodeRepo.On("Add", node, mock.Anything).Return(nil).Once()

	nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()

	t.Run("NodeStateValid", func(t *testing.T) {
		// Arrange
		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("Get", mock.Anything,
			&opb.GetRequest{Id: orgId.String()}).
			Return(&opb.GetResponse{
				Org: &opb.Organization{
					Id: orgId.String(),
				},
			}, nil).Once()

		// Act
		res, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
			NodeId: nodeId,
			Name:   nodeName,
			OrgId:  orgId.String(),
			State:  nodeState,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, nodeName, res.Node.Name)
		assert.Equal(t, nodeType, res.Node.Type)
		assert.Equal(t, node.State.String(), res.Node.State)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("NodeStateInvalid", func(t *testing.T) {
		// Arrange
		const nState = "unknown"

		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("Get", mock.Anything,
			&opb.GetRequest{Id: orgId.String()}).
			Return(&opb.GetResponse{
				Org: &opb.Organization{
					Id: orgId.String(),
				},
			}, nil).Once()

		// Act
		res, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
			NodeId: nodeId,
			Name:   nodeName,
			OrgId:  orgId.String(),
			State:  nState,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("RemoteOrgIsDeactivated", func(t *testing.T) {
		// Arrange
		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("Get", mock.Anything,
			&opb.GetRequest{Id: orgId.String()}).
			Return(&opb.GetResponse{
				Org: &opb.Organization{
					Id:            orgId.String(),
					IsDeactivated: true,
				},
			}, nil).Once()

		// Act
		res, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
			// NodeId: nodeId,
			Name:  nodeName,
			OrgId: orgId.String(),
			State: nodeState,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("RemoteOrgNotFound", func(t *testing.T) {
		// Arrange
		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("Get", mock.Anything,
			&opb.GetRequest{Id: orgId.String()}).
			Return(nil, gorm.ErrRecordNotFound).Once()

		// Act
		res, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
			// NodeId: nodeId,
			Name:  nodeName,
			OrgId: orgId.String(),
			State: nodeState,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("OrgServiceFailure", func(t *testing.T) {
		// Arrange
		orgService.On("GetClient").
			Return(nil, errors.New("fail to get client")).
			Once()

		// Act
		res, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
			// NodeId: nodeId,
			Name:  nodeName,
			OrgId: orgId.String(),
			State: nodeState,
		})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, res)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeServer_Get(t *testing.T) {
	t.Run("NodeFound", func(t *testing.T) {
		const nodeName = "node-A"
		const nodeType = ukama.NODE_ID_TYPE_HOMENODE
		var nodeId = ukama.NewVirtualNodeId(nodeType)

		nodeRepo := &mocks.NodeRepo{}

		nodeRepo.On("Get", nodeId).Return(
			&db.Node{Id: nodeId.StringLowercase(),
				Name:  nodeName,
				Type:  ukama.NODE_ID_TYPE_HOMENODE,
				State: db.Online,
			}, nil).Once()

		s := server.NewNodeServer(nodeRepo, nil, "", nil, nil)

		resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{
			NodeId: nodeId.StringLowercase()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, nodeId.String(), resp.GetNode().GetId())
		assert.Equal(t, nodeName, resp.GetNode().Name)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		var nodeId = ukama.NewVirtualAmplifierNodeId()

		nodeRepo := &mocks.NodeRepo{}

		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(nodeRepo, nil, "", nil, nil)

		resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{
			NodeId: nodeId.StringLowercase()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("NodeIdInvalid", func(t *testing.T) {
		var nodeId = uuid.NewV4()

		nodeRepo := &mocks.NodeRepo{}

		s := server.NewNodeServer(nodeRepo, nil, "", nil, nil)

		resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{
			NodeId: nodeId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeServer_Delete(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	s := server.NewNodeServer(nodeRepo, nil, "", nil, nil)

	nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()

	t.Run("NodeFound", func(t *testing.T) {
		const nodeType = ukama.NODE_ID_TYPE_HOMENODE
		var nodeId = ukama.NewVirtualNodeId(nodeType)

		nodeRepo.On("Delete", nodeId, mock.Anything).Return(nil).Once()

		resp, err := s.DeleteNode(context.TODO(), &pb.DeleteNodeRequest{
			NodeId: nodeId.StringLowercase()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		var nodeId = ukama.NewVirtualAmplifierNodeId()

		nodeRepo.On("Delete", nodeId, mock.Anything).
			Return(gorm.ErrRecordNotFound).Once()

		resp, err := s.DeleteNode(context.TODO(), &pb.DeleteNodeRequest{
			NodeId: nodeId.StringLowercase()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("NodeIdInvalid", func(t *testing.T) {
		var nodeId = uuid.NewV4()

		resp, err := s.DeleteNode(context.TODO(), &pb.DeleteNodeRequest{
			NodeId: nodeId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}
