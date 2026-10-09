/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/mocks"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"github.com/ukama/ukama/systems/registry/node/pkg/server"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	sitepb "github.com/ukama/ukama/systems/registry/site/pb/gen"
	sitemocks "github.com/ukama/ukama/systems/registry/site/pb/gen/mocks"
)

var (
	testNode = ukama.NewVirtualNodeId("HomeNode")
	orgId    = uuid.NewV4()
)

const OrgName = "testorg"

func TestNodeServer_Add(t *testing.T) {
	nodeId := testNode.String()

	msgbusClient := &mbmocks.MsgBusServiceClient{}
	nodeRepo := &mocks.NodeRepo{}
	nodeStatusRepo := &mocks.NodeStatusRepo{}
	siteService := &mocks.SiteClientProvider{}

	const nodeName = "node-A"
	const nodeType = ukama.NODE_TYPE_HOMENODE

	s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "",
		msgbusClient, siteService, orgId, nil, nil)

	node := &db.Node{
		Id:   nodeId,
		Name: nodeName,
		Type: ukama.NodeType(testNode.GetNodeType()),
		Status: db.NodeStatus{
			NodeId:       nodeId,
			State:        ukama.NodeStateUnknown,
			Connectivity: ukama.NodeConnectivityUndefined,
		},
	}

	nodeRepo.On("Add", node, mock.Anything).Return(nil).Once()
	nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
		Return(nil).Once()

	t.Run("NodeStateValid", func(t *testing.T) {
		// Act
		res, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
			NodeId: nodeId,
			Name:   nodeName,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, nodeName, res.Node.Name)
		assert.Equal(t, nodeType, res.Node.Type)
		nodeRepo.AssertExpectations(t)
	})

}

func TestNodeServer_Get(t *testing.T) {
	t.Run("NodeFound", func(t *testing.T) {
		const (
			nodeName = "node-A"
			nodeType = ukama.NODE_TYPE_HOMENODE
		)

		var nodeId = ukama.NewVirtualNodeId(nodeType)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).Return(
			&db.Node{Id: nodeId.StringLowercase(),
				Name: nodeName,
				Type: ukama.NODE_TYPE_HOMENODE,
			}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{
			NodeId: nodeId.StringLowercase()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, nodeId.StringLowercase(), resp.GetNode().GetId())
		assert.Equal(t, nodeName, resp.GetNode().Name)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		var nodeId = ukama.NewVirtualAmplifierNodeId()

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{
			NodeId: nodeId.StringLowercase()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("NodeIdInvalid", func(t *testing.T) {
		var nodeId = uuid.NewV4()

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}
		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{
			NodeId: nodeId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeServer_List(t *testing.T) {
	t.Run("ListWithAllFilters", func(t *testing.T) {
		// Arrange
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		siteId := uuid.NewV4()
		networkId := uuid.NewV4()
		ntype := ukama.NODE_ID_TYPE_HOMENODE
		connectivity := ukama.NodeConnectivityOnline
		state := ukama.NodeStateUnknown

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodes := []db.Node{
			{
				Id:   nodeId.StringLowercase(),
				Name: "node-1",
				Type: ukama.NodeType(ntype),
				Status: db.NodeStatus{
					NodeId:       nodeId.StringLowercase(),
					Connectivity: connectivity,
					State:        state,
				},
				Site: db.Site{
					NodeId:    nodeId.StringLowercase(),
					SiteId:    siteId,
					NetworkId: networkId,
				},
			},
		}

		connectivityVal := uint8(connectivity)
		stateVal := uint8(state)
		nodeRepo.On("List", nodeId.StringLowercase(), siteId.String(),
			networkId.String(), ntype, &connectivityVal, &stateVal).
			Return(nodes, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		// Act
		resp, err := s.List(context.TODO(), &pb.ListRequest{
			NodeId:       nodeId.StringLowercase(),
			SiteId:       siteId.String(),
			NetworkId:    networkId.String(),
			Type:         ntype,
			Connectivity: cpb.NodeConnectivity(connectivity),
			State:        cpb.NodeState(state),
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Nodes, 1)
		assert.Equal(t, nodeId.String(), resp.Nodes[0].Id)
		assert.Equal(t, "node-1", resp.Nodes[0].Name)
		assert.Equal(t, ntype, resp.Nodes[0].Type)
		assert.Equal(t, cpb.NodeConnectivity(connectivity),
			resp.Nodes[0].Status.Connectivity)
		assert.Equal(t, cpb.NodeState(state), resp.Nodes[0].Status.State)
		assert.Equal(t, siteId.String(), resp.Nodes[0].Site.SiteId)
		assert.Equal(t, networkId.String(), resp.Nodes[0].Site.NetworkId)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("ListWithNoFilters", func(t *testing.T) {
		// Arrange
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		siteId := uuid.NewV4()
		networkId := uuid.NewV4()
		ntype := ukama.NODE_ID_TYPE_HOMENODE
		connectivity := cpb.NodeConnectivity_Online
		state := cpb.NodeState_Unknown

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodes := []db.Node{
			{
				Id:   nodeId.StringLowercase(),
				Name: "node-1",
				Type: ukama.NodeType(ntype),
				Status: db.NodeStatus{
					NodeId:       nodeId.StringLowercase(),
					Connectivity: ukama.NodeConnectivityOnline,
					State:        ukama.NodeStateUnknown,
				},
				Site: db.Site{
					NodeId:    nodeId.StringLowercase(),
					SiteId:    siteId,
					NetworkId: networkId,
				},
			},
		}

		connectivityVal := uint8(ukama.NodeConnectivityOnline)
		stateVal := uint8(ukama.NodeStateUnknown)
		nodeRepo.On("List", nodeId.StringLowercase(), siteId.String(),
			networkId.String(), ntype, &connectivityVal, &stateVal).
			Return(nodes, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		// Act
		resp, err := s.List(context.TODO(), &pb.ListRequest{
			NodeId:       nodeId.StringLowercase(),
			SiteId:       siteId.String(),
			NetworkId:    networkId.String(),
			Type:         ntype,
			Connectivity: connectivity,
			State:        state,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Nodes, 1)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("ListWithPartialFilters", func(t *testing.T) {
		// Arrange
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		connectivity := cpb.NodeConnectivity_Online
		state := cpb.NodeState_Unknown

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodes := []db.Node{
			{
				Id:   nodeId.StringLowercase(),
				Name: "node-1",
				Type: ukama.NODE_ID_TYPE_HOMENODE,
				Status: db.NodeStatus{
					NodeId:       nodeId.StringLowercase(),
					Connectivity: ukama.NodeConnectivityOnline,
					State:        ukama.NodeStateUnknown,
				},
			},
		}

		connectivityVal := uint8(ukama.NodeConnectivityOnline)
		stateVal := uint8(ukama.NodeStateUnknown)
		nodeRepo.On("List", nodeId.StringLowercase(), "", "", "", &connectivityVal,
			&stateVal).
			Return(nodes, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		// Act
		resp, err := s.List(context.TODO(), &pb.ListRequest{
			NodeId:       nodeId.StringLowercase(),
			Connectivity: connectivity,
			State:        state,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Nodes, 1)
		assert.Equal(t, nodeId.String(), resp.Nodes[0].Id)
		assert.Equal(t, connectivity, resp.Nodes[0].Status.Connectivity)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("ListWithNoResults", func(t *testing.T) {
		// Arrange
		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		// Create pointers to uint8 values to match the actual behavior
		connectivityVal := uint8(0)
		stateVal := uint8(0)
		nodeRepo.On("List", "", "", "", "", &connectivityVal, &stateVal).
			Return([]db.Node{}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		// Act
		resp, err := s.List(context.TODO(), &pb.ListRequest{})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Nodes, 0)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("ListWithError", func(t *testing.T) {
		// Arrange
		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		// Create pointers to uint8 values to match the actual behavior
		connectivityVal := uint8(0)
		stateVal := uint8(0)
		nodeRepo.On("List", "", "", "", "", &connectivityVal, &stateVal).
			Return(nil, gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		// Act
		resp, err := s.List(context.TODO(), &pb.ListRequest{})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeServer_GetNodesForSite(t *testing.T) {
	t.Run("SiteFound", func(t *testing.T) {
		siteId := uuid.NewV4()
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		siteRepo := &mocks.SiteRepo{}
		nodes := []db.Node{
			{Id: nodeId.StringLowercase(), Name: "node-1"},
		}
		siteRepo.On("GetNodes", siteId).Return(nodes, nil).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.GetNodesForSite(context.TODO(),
			&pb.GetBySiteRequest{SiteId: siteId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Nodes, 1)
		siteRepo.AssertExpectations(t)
	})

	t.Run("InvalidSiteId", func(t *testing.T) {
		siteRepo := &mocks.SiteRepo{}
		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.GetNodesForSite(context.TODO(),
			&pb.GetBySiteRequest{SiteId: "not-a-uuid"})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("RepoError", func(t *testing.T) {
		siteId := uuid.NewV4()
		siteRepo := &mocks.SiteRepo{}
		siteRepo.On("GetNodes", siteId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.GetNodesForSite(context.TODO(),
			&pb.GetBySiteRequest{SiteId: siteId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		siteRepo.AssertExpectations(t)
	})
}

func TestNodeServer_GetNodesForNetwork(t *testing.T) {
	t.Run("NetworkFound", func(t *testing.T) {
		networkId := uuid.NewV4()
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		siteRepo := &mocks.SiteRepo{}
		nodes := []db.Node{
			{Id: nodeId.StringLowercase(), Name: "node-1"},
		}
		siteRepo.On("GetByNetwork", networkId).Return(nodes, nil).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.GetNodesForNetwork(context.TODO(),
			&pb.GetByNetworkRequest{NetworkId: networkId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Nodes, 1)
		siteRepo.AssertExpectations(t)
	})

	t.Run("InvalidNetworkId", func(t *testing.T) {
		siteRepo := &mocks.SiteRepo{}
		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.GetNodesForNetwork(context.TODO(),
			&pb.GetByNetworkRequest{NetworkId: "not-a-uuid"})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("RepoError", func(t *testing.T) {
		networkId := uuid.NewV4()
		siteRepo := &mocks.SiteRepo{}
		siteRepo.On("GetByNetwork", networkId).
			Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.GetNodesForNetwork(context.TODO(),
			&pb.GetByNetworkRequest{NetworkId: networkId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		siteRepo.AssertExpectations(t)
	})
}

/** Deprecated: Use List API instead */
func TestNodeServer_GetNodes(t *testing.T) {
	t.Run("NodesFound", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodes := []db.Node{
			{Id: nodeId.StringLowercase(), Name: "node-1"},
		}
		nodeRepo.On("GetAll").Return(nodes, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.GetNodes(context.TODO(), &pb.GetNodesRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Nodes, 1)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("RepoError", func(t *testing.T) {
		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("GetAll").Return(nil, gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.GetNodes(context.TODO(), &pb.GetNodesRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}

/** Deprecated: Use List API instead */
func TestNodeServer_GetNodesByState(t *testing.T) {
	t.Run("NodesFound", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodes := []db.Node{
			{Id: nodeId.StringLowercase(), Name: "node-1"},
		}
		nodeRepo.On("GetNodesByState", uint8(cpb.NodeConnectivity_Online),
			uint8(cpb.NodeState_Unknown)).
			Return(nodes, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.GetNodesByState(context.TODO(), &pb.GetNodesByStateRequest{
			Connectivity: cpb.NodeConnectivity_Online,
			State:        cpb.NodeState_Unknown,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Nodes, 1)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("RepoError", func(t *testing.T) {
		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("GetNodesByState", uint8(cpb.NodeConnectivity_Online),
			uint8(cpb.NodeState_Unknown)).
			Return(nil, gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.GetNodesByState(context.TODO(), &pb.GetNodesByStateRequest{
			Connectivity: cpb.NodeConnectivity_Online,
			State:        cpb.NodeState_Unknown,
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeServer_UpdateNodeStatus(t *testing.T) {
	t.Run("UpdateSuccess", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		nodeStatusRepo.On("Update", mock.Anything).Return(nil).Once()
		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
		nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "",
			msgbusClient, nil, orgId, nil, nil)

		resp, err := s.UpdateNodeStatus(context.TODO(), &pb.UpdateNodeStateRequest{
			NodeId:       nodeId.StringLowercase(),
			Connectivity: ukama.NodeConnectivityOnline.String(),
			State:        ukama.NodeStateConfiguring.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeStatusRepo.AssertExpectations(t)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("InvalidNodeId", func(t *testing.T) {
		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}
		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		resp, err := s.UpdateNodeStatus(context.TODO(), &pb.UpdateNodeStateRequest{
			NodeId: uuid.NewV4().String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("UpdateRepoError", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeStatusRepo.On("Update", mock.Anything).Return(gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		resp, err := s.UpdateNodeStatus(context.TODO(), &pb.UpdateNodeStateRequest{
			NodeId:       nodeId.StringLowercase(),
			Connectivity: ukama.NodeConnectivityOnline.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeStatusRepo.AssertExpectations(t)
	})

	t.Run("GetNodeError", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeStatusRepo.On("Update", mock.Anything).Return(nil).Once()
		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)

		resp, err := s.UpdateNodeStatus(context.TODO(), &pb.UpdateNodeStateRequest{
			NodeId:       nodeId.StringLowercase(),
			Connectivity: ukama.NodeConnectivityOnline.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeStatusRepo.AssertExpectations(t)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeServer_UpdateNode(t *testing.T) {
	t.Run("UpdateSuccess", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		nodeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase(), Name: "node-new"}, nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", msgbusClient, nil,
			orgId, nil, nil)

		resp, err := s.UpdateNode(context.TODO(), &pb.UpdateNodeRequest{
			NodeId: nodeId.StringLowercase(),
			Name:   "node-new",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "node-new", resp.Node.Name)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("InvalidNodeId", func(t *testing.T) {
		nodeRepo := &mocks.NodeRepo{}
		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.UpdateNode(context.TODO(), &pb.UpdateNodeRequest{
			NodeId: uuid.NewV4().String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("UpdateRepoError", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Update", mock.Anything, mock.Anything).
			Return(gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.UpdateNode(context.TODO(), &pb.UpdateNodeRequest{
			NodeId: nodeId.StringLowercase(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	// When the follow-up Get fails, UpdateNode still reports success using
	// the request-derived node rather than failing the whole call.
	t.Run("GetAfterUpdateFails", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.UpdateNode(context.TODO(), &pb.UpdateNodeRequest{
			NodeId: nodeId.StringLowercase(),
			Name:   "node-new",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "node-new", resp.Node.Name)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeServer_DeleteNode(t *testing.T) {
	t.Run("DeleteSuccess", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		nodeRepo.On("Delete", nodeId, mock.Anything).Return(nil).Once()
		nodeRepo.On("GetNodeCount").Return(int64(0), int64(0), int64(0), nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", msgbusClient, nil,
			orgId, nil, nil)

		resp, err := s.DeleteNode(context.TODO(),
			&pb.DeleteNodeRequest{NodeId: nodeId.StringLowercase()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("InvalidNodeId", func(t *testing.T) {
		nodeRepo := &mocks.NodeRepo{}
		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.DeleteNode(context.TODO(),
			&pb.DeleteNodeRequest{NodeId: uuid.NewV4().String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("DeleteRepoError", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Delete", nodeId, mock.Anything).Return(gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.DeleteNode(context.TODO(),
			&pb.DeleteNodeRequest{NodeId: nodeId.StringLowercase()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeServer_AttachNodes(t *testing.T) {
	t.Run("AttachSuccess", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)
		ampId := ukama.NewVirtualAmplifierNodeId()

		nodeRepo := &mocks.NodeRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		nodeRepo.On("AttachNodes", nodeId, []string{ampId.StringLowercase()}).
			Return(nil).Once()
		nodeRepo.On("GetNodeCount").Return(int64(2), int64(0), int64(0), nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", msgbusClient, nil,
			orgId, nil, nil)

		resp, err := s.AttachNodes(context.TODO(), &pb.AttachNodesRequest{
			NodeId:        nodeId.StringLowercase(),
			AttachedNodes: []string{ampId.StringLowercase()},
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("InvalidNodeId", func(t *testing.T) {
		nodeRepo := &mocks.NodeRepo{}
		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.AttachNodes(context.TODO(),
			&pb.AttachNodesRequest{NodeId: uuid.NewV4().String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("AttachRepoError", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)
		ampId := ukama.NewVirtualAmplifierNodeId()

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("AttachNodes", nodeId, []string{ampId.StringLowercase()}).
			Return(gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.AttachNodes(context.TODO(), &pb.AttachNodesRequest{
			NodeId:        nodeId.StringLowercase(),
			AttachedNodes: []string{ampId.StringLowercase()},
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeServer_DetachNode(t *testing.T) {
	t.Run("DetachSuccess", func(t *testing.T) {
		nodeId := ukama.NewVirtualAmplifierNodeId()

		nodeRepo := &mocks.NodeRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
		nodeRepo.On("DetachNode", nodeId).Return(nil).Once()
		nodeRepo.On("GetNodeCount").Return(int64(2), int64(0), int64(0), nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", msgbusClient, nil,
			orgId, nil, nil)

		resp, err := s.DetachNode(context.TODO(),
			&pb.DetachNodeRequest{NodeId: nodeId.StringLowercase()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("InvalidNodeId", func(t *testing.T) {
		nodeRepo := &mocks.NodeRepo{}
		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.DetachNode(context.TODO(),
			&pb.DetachNodeRequest{NodeId: uuid.NewV4().String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("GetNodeError", func(t *testing.T) {
		nodeId := ukama.NewVirtualAmplifierNodeId()

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.DetachNode(context.TODO(),
			&pb.DetachNodeRequest{NodeId: nodeId.StringLowercase()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("DetachRepoError", func(t *testing.T) {
		nodeId := ukama.NewVirtualAmplifierNodeId()

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
		nodeRepo.On("DetachNode", nodeId).Return(gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.DetachNode(context.TODO(),
			&pb.DetachNodeRequest{NodeId: nodeId.StringLowercase()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeServer_AddNodeToSite(t *testing.T) {
	t.Run("AddSuccess", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		siteId := uuid.NewV4()
		networkId := uuid.NewV4()

		siteRepo := &mocks.SiteRepo{}
		siteService := &mocks.SiteClientProvider{}
		siteClient := &sitemocks.SiteServiceClient{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		siteService.On("GetClient").Return(siteClient, nil).Once()
		siteClient.On("Get", mock.Anything,
			&sitepb.GetRequest{SiteId: siteId.String()}, mock.Anything).
			Return(&sitepb.GetResponse{
				Site: &sitepb.Site{NetworkId: networkId.String()},
			}, nil).Once()
		siteRepo.On("AddNode", mock.Anything, mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", msgbusClient,
			siteService, orgId, nil, nil)

		resp, err := s.AddNodeToSite(context.TODO(), &pb.AddNodeToSiteRequest{
			NodeId:    nodeId.StringLowercase(),
			SiteId:    siteId.String(),
			NetworkId: networkId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		siteRepo.AssertExpectations(t)
		siteClient.AssertExpectations(t)
	})

	t.Run("InvalidNodeId", func(t *testing.T) {
		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)

		resp, err := s.AddNodeToSite(context.TODO(),
			&pb.AddNodeToSiteRequest{NodeId: uuid.NewV4().String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("InvalidNetworkId", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)

		resp, err := s.AddNodeToSite(context.TODO(), &pb.AddNodeToSiteRequest{
			NodeId:    nodeId.StringLowercase(),
			NetworkId: "not-a-uuid",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("InvalidSiteId", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		networkId := uuid.NewV4()
		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)

		resp, err := s.AddNodeToSite(context.TODO(), &pb.AddNodeToSiteRequest{
			NodeId:    nodeId.StringLowercase(),
			NetworkId: networkId.String(),
			SiteId:    "not-a-uuid",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("NetworkMismatch", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		siteId := uuid.NewV4()
		networkId := uuid.NewV4()
		otherNetworkId := uuid.NewV4()

		siteService := &mocks.SiteClientProvider{}
		siteClient := &sitemocks.SiteServiceClient{}

		siteService.On("GetClient").Return(siteClient, nil).Once()
		siteClient.On("Get", mock.Anything,
			&sitepb.GetRequest{SiteId: siteId.String()}, mock.Anything).
			Return(&sitepb.GetResponse{
				Site: &sitepb.Site{NetworkId: otherNetworkId.String()},
			}, nil).Once()

		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, siteService, orgId,
			nil, nil)

		resp, err := s.AddNodeToSite(context.TODO(), &pb.AddNodeToSiteRequest{
			NodeId:    nodeId.StringLowercase(),
			SiteId:    siteId.String(),
			NetworkId: networkId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		siteClient.AssertExpectations(t)
	})
}

func TestNodeServer_ReleaseNodeFromSite(t *testing.T) {
	t.Run("ReleaseSuccess", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
		siteId := uuid.NewV4()
		networkId := uuid.NewV4()

		siteRepo := &mocks.SiteRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		siteRepo.On("RemoveNode", nodeId).
			Return(&db.Site{SiteId: siteId, NetworkId: networkId}, nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", msgbusClient, nil,
			orgId, nil, nil)

		resp, err := s.ReleaseNodeFromSite(context.TODO(),
			&pb.ReleaseNodeFromSiteRequest{NodeId: nodeId.StringLowercase()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		siteRepo.AssertExpectations(t)
	})

	t.Run("InvalidNodeId", func(t *testing.T) {
		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)

		resp, err := s.ReleaseNodeFromSite(context.TODO(),
			&pb.ReleaseNodeFromSiteRequest{NodeId: uuid.NewV4().String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("RepoError", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		siteRepo := &mocks.SiteRepo{}
		siteRepo.On("RemoveNode", nodeId).Return(nil, gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.ReleaseNodeFromSite(context.TODO(),
			&pb.ReleaseNodeFromSiteRequest{NodeId: nodeId.StringLowercase()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		siteRepo.AssertExpectations(t)
	})
}

func TestNodeServer_PushMetrics(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	nodeRepo.On("GetNodeCount").
		Return(int64(0), int64(0), int64(0), gorm.ErrInvalidDB).Once()

	s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
		nil, nil)

	// Act: PushMetrics only wraps pushNodeMeterics, exercised here for the
	// GetNodeCount-error branch (the pushgateway call itself isn't under test).
	s.PushMetrics()

	nodeRepo.AssertExpectations(t)
}
