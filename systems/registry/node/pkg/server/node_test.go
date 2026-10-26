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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/mocks"
	"github.com/ukama/ukama/systems/registry/node/pkg"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"github.com/ukama/ukama/systems/registry/node/pkg/server"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
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

		nodeStatusRepo.On("Update", mock.Anything, "connectivity", "state").Return(nil).Once()
		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
		// The event must carry both fields when a request changes both (#1784).
		msgbusClient.On("PublishRequest",
			"event.cloud.local.testorg.registry.node.node.state.update",
			mock.MatchedBy(func(e *epb.EventRegistryNodeStatusUpdate) bool {
				return e.Connectivity == ukama.NodeConnectivityOnline.String() &&
					e.State == ukama.NodeStateConfiguring.String()
			})).Return(nil).Once()

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
		msgbusClient.AssertExpectations(t)
	})

	t.Run("ConnectivityOnlyEventLeavesStateEmpty", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		nodeStatusRepo.On("Update", mock.Anything, "connectivity").Return(nil).Once()
		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything,
			mock.MatchedBy(func(e *epb.EventRegistryNodeStatusUpdate) bool {
				return e.Connectivity == ukama.NodeConnectivityOffline.String() && e.State == ""
			})).Return(nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "",
			msgbusClient, nil, orgId, nil, nil)

		_, err := s.UpdateNodeStatus(context.TODO(), &pb.UpdateNodeStateRequest{
			NodeId:       nodeId.StringLowercase(),
			Connectivity: ukama.NodeConnectivityOffline.String(),
		})

		assert.NoError(t, err)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("UnknownStateIsWritten", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeStatusRepo.On("Update", mock.MatchedBy(func(ns *db.NodeStatus) bool {
			return ns.State == ukama.NodeStateUnknown
		}), "state").Return(nil).Once()
		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "",
			nil, nil, orgId, nil, nil)

		_, err := s.UpdateNodeStatus(context.TODO(), &pb.UpdateNodeStateRequest{
			NodeId: nodeId.StringLowercase(),
			State:  ukama.NodeStateUnknown.String(),
		})

		assert.NoError(t, err)
		nodeStatusRepo.AssertExpectations(t)
	})

	invalidCases := []struct {
		name         string
		state        string
		connectivity string
	}{
		{"EmptyUpdate", "", ""},
		{"UnknownState", "bogus", ""},
		{"NumericState", "3", ""},
		{"UnknownConnectivity", "", "bogus"},
		{"NumericConnectivity", "", "1"},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name+"IsRejected", func(t *testing.T) {
			nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

			nodeRepo := &mocks.NodeRepo{}
			nodeStatusRepo := &mocks.NodeStatusRepo{}
			msgbusClient := &mbmocks.MsgBusServiceClient{}

			s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "",
				msgbusClient, nil, orgId, nil, nil)

			resp, err := s.UpdateNodeStatus(context.TODO(), &pb.UpdateNodeStateRequest{
				NodeId:       nodeId.StringLowercase(),
				State:        tc.state,
				Connectivity: tc.connectivity,
			})

			assert.Equal(t, codes.InvalidArgument, status.Code(err))
			assert.Nil(t, resp)
			nodeStatusRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			nodeStatusRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything,
				mock.Anything)
			msgbusClient.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
		})
	}

	t.Run("MixedCaseValuesAreAccepted", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeStatusRepo.On("Update", mock.MatchedBy(func(ns *db.NodeStatus) bool {
			return ns.State == ukama.NodeStateOperational &&
				ns.Connectivity == ukama.NodeConnectivityOffline
		}), "connectivity", "state").Return(nil).Once()
		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "",
			nil, nil, orgId, nil, nil)

		_, err := s.UpdateNodeStatus(context.TODO(), &pb.UpdateNodeStateRequest{
			NodeId:       nodeId.StringLowercase(),
			State:        "OPERATIONAL",
			Connectivity: "Offline",
		})

		assert.NoError(t, err)
		nodeStatusRepo.AssertExpectations(t)
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

		nodeStatusRepo.On("Update", mock.Anything, "connectivity").
			Return(gorm.ErrInvalidDB).Once()

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

		nodeStatusRepo.On("Update", mock.Anything, "connectivity").Return(nil).Once()
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

	t.Run("LocationOnlyPublishesStoredName", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)

		nodeRepo := &mocks.NodeRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		nodeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		nodeRepo.On("Get", nodeId).Return(&db.Node{
			Id: nodeId.StringLowercase(), Name: "tower-1", Latitude: "1.5", Longitude: "2.5",
		}, nil).Once()
		msgbusClient.On("PublishRequest", "event.cloud.local.testorg.registry.node.node.update",
			mock.MatchedBy(func(e *epb.EventRegistryNodeUpdate) bool {
				return e.Name == "tower-1" && e.Latitude == "1.5" && e.Longitude == "2.5"
			})).Return(nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", msgbusClient, nil,
			orgId, nil, nil)

		_, err := s.UpdateNode(context.TODO(), &pb.UpdateNodeRequest{
			NodeId:    nodeId.StringLowercase(),
			Latitude:  "1.5",
			Longitude: "2.5",
		})

		assert.NoError(t, err)
		msgbusClient.AssertExpectations(t)
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
		siteRepo.On("AssignNodes", siteId, networkId, []string{nodeId.StringLowercase()}).
			Return([]db.Site{{NodeId: nodeId.StringLowercase(), SiteId: siteId,
				NetworkId: networkId}}, nil).Once()
		msgbusClient.On("PublishRequest", "event.cloud.local.testorg.registry.node.node.assign",
			mock.MatchedBy(func(e *epb.EventRegistryNodeAssign) bool {
				return e.NodeId == nodeId.StringLowercase() && e.Site == siteId.String()
			})).Return(nil).Once()

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

	assignCases := []struct {
		name     string
		assigned []db.Site
		repoErr  error
		code     codes.Code
	}{
		{"AlreadyOnSameSiteIsSuccess", []db.Site{}, nil, codes.OK},
		{"OnOtherSiteIsAlreadyExists", nil,
			fmt.Errorf("node already belongs to site: %w", store.ErrAlreadyExists),
			codes.AlreadyExists},
		{"MissingNodeIsNotFound", nil, fmt.Errorf("node: %w", store.ErrNotFound),
			codes.NotFound},
	}

	for _, tc := range assignCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
			siteId := uuid.NewV4()
			networkId := uuid.NewV4()

			siteRepo := &mocks.SiteRepo{}
			siteService := &mocks.SiteClientProvider{}
			siteClient := &sitemocks.SiteServiceClient{}
			msgbusClient := &mbmocks.MsgBusServiceClient{}

			siteService.On("GetClient").Return(siteClient, nil).Once()
			siteClient.On("Get", mock.Anything, mock.Anything, mock.Anything).
				Return(&sitepb.GetResponse{
					Site: &sitepb.Site{NetworkId: networkId.String()},
				}, nil).Once()
			siteRepo.On("AssignNodes", siteId, networkId, []string{nodeId.StringLowercase()}).
				Return(tc.assigned, tc.repoErr).Once()

			s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", msgbusClient,
				siteService, orgId, nil, nil)

			_, err := s.AddNodeToSite(context.TODO(), &pb.AddNodeToSiteRequest{
				NodeId:    nodeId.StringLowercase(),
				SiteId:    siteId.String(),
				NetworkId: networkId.String(),
			})

			assert.Equal(t, tc.code, status.Code(err))
			siteRepo.AssertExpectations(t)
			msgbusClient.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
		})
	}

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
		msgbusClient.On("PublishRequest",
			"event.cloud.local.testorg.registry.node.node.release",
			mock.MatchedBy(func(e *epb.NodeReleasedEvent) bool {
				return e.NodeId == nodeId.StringLowercase() && e.Site == siteId.String() &&
					e.Network == networkId.String()
			})).Return(nil).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", msgbusClient, nil,
			orgId, nil, nil)

		resp, err := s.ReleaseNodeFromSite(context.TODO(),
			&pb.ReleaseNodeFromSiteRequest{NodeId: nodeId.StringLowercase()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		siteRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("InvalidNodeId", func(t *testing.T) {
		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)

		resp, err := s.ReleaseNodeFromSite(context.TODO(),
			&pb.ReleaseNodeFromSiteRequest{NodeId: uuid.NewV4().String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("NotAssignedIsFailedPrecondition", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		siteRepo := &mocks.SiteRepo{}
		siteRepo.On("RemoveNode", nodeId).Return(nil,
			fmt.Errorf("node is not yet assigned: %w", store.ErrPrecondition)).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil, orgId,
			nil, nil)

		_, err := s.ReleaseNodeFromSite(context.TODO(),
			&pb.ReleaseNodeFromSiteRequest{NodeId: nodeId.StringLowercase()})

		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		siteRepo.AssertExpectations(t)
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

func TestNodeServer_PushMetricsPushesAllCounts(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	nodeRepo.On("GetNodeCount").Return(int64(5), int64(3), int64(2), nil).Once()

	s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "http://%zz", nil, nil,
		orgId, nil, nil)

	s.PushMetrics()

	pushed := map[string]float64{}
	for _, m := range pkg.NodeMetric {
		pushed[m.Name] = m.Value
	}

	assert.Equal(t, map[string]float64{
		pkg.NumberOfNodes:        5,
		pkg.NumberOfOnlineNodes:  3,
		pkg.NumberOfOfflineNodes: 2,
	}, pushed)
	nodeRepo.AssertExpectations(t)
}

func TestNodeServer_ReleaseSite(t *testing.T) {
	releaseRoute := "event.cloud.local.testorg.registry.node.node.release"

	t.Run("PublishesReleasePerNode", func(t *testing.T) {
		siteId := uuid.NewV4()
		networkId := uuid.NewV4()
		tNode := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE).StringLowercase()
		aNode := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_AMPNODE).StringLowercase()

		siteRepo := &mocks.SiteRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		siteRepo.On("ReleaseSite", siteId).Return([]db.Site{
			{NodeId: tNode, SiteId: siteId, NetworkId: networkId},
			{NodeId: aNode, SiteId: siteId, NetworkId: networkId},
		}, nil).Once()
		for _, id := range []string{tNode, aNode} {
			nodeId := id
			msgbusClient.On("PublishRequest", releaseRoute,
				mock.MatchedBy(func(e *epb.NodeReleasedEvent) bool {
					return e.NodeId == nodeId && e.Site == siteId.String() &&
						e.Network == networkId.String()
				})).Return(nil).Once()
		}

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", msgbusClient, nil,
			orgId, nil, nil)

		resp, err := s.ReleaseSite(context.TODO(),
			&pb.ReleaseSiteRequest{SiteId: siteId.String()})

		assert.NoError(t, err)
		assert.Equal(t, []string{tNode, aNode}, resp.NodeIds)
		siteRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
	})

	t.Run("EmptySiteIsNoop", func(t *testing.T) {
		siteId := uuid.NewV4()

		siteRepo := &mocks.SiteRepo{}
		msgbusClient := &mbmocks.MsgBusServiceClient{}

		siteRepo.On("ReleaseSite", siteId).Return([]db.Site{}, nil).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", msgbusClient, nil,
			orgId, nil, nil)

		resp, err := s.ReleaseSite(context.TODO(),
			&pb.ReleaseSiteRequest{SiteId: siteId.String()})

		assert.NoError(t, err)
		assert.Empty(t, resp.NodeIds)
		msgbusClient.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
	})

	t.Run("InvalidSiteId", func(t *testing.T) {
		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)

		resp, err := s.ReleaseSite(context.TODO(),
			&pb.ReleaseSiteRequest{SiteId: "not-a-uuid"})

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Nil(t, resp)
	})

	t.Run("RepoError", func(t *testing.T) {
		siteId := uuid.NewV4()

		siteRepo := &mocks.SiteRepo{}
		siteRepo.On("ReleaseSite", siteId).Return(nil, gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil, orgId,
			nil, nil)

		resp, err := s.ReleaseSite(context.TODO(),
			&pb.ReleaseSiteRequest{SiteId: siteId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		siteRepo.AssertExpectations(t)
	})
}

func TestNodeServer_StartMetricsPusher(t *testing.T) {
	t.Run("BurstOfWritesPushesOnce", func(t *testing.T) {
		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Add", mock.Anything, mock.Anything).Return(nil).Times(3)

		pushed := make(chan struct{}, 3)
		nodeRepo.On("GetNodeCount").Return(int64(3), int64(0), int64(0), nil).
			Run(func(mock.Arguments) { pushed <- struct{}{} })

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "http://%zz", nil, nil,
			orgId, nil, nil)

		// No pusher is running yet, so these writes must not block on metrics.
		for i := 0; i < 3; i++ {
			_, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
				NodeId: ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String(),
			})
			assert.NoError(t, err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		s.StartMetricsPusher(ctx)

		select {
		case <-pushed:
		case <-time.After(2 * time.Second):
			t.Fatal("pusher did not push node counts")
		}

		select {
		case <-pushed:
			t.Fatal("three writes should coalesce into one push")
		case <-time.After(200 * time.Millisecond):
		}
	})

	t.Run("PushesAfterWrite", func(t *testing.T) {
		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Add", mock.Anything, mock.Anything).Return(nil).Once()

		pushed := make(chan struct{}, 1)
		nodeRepo.On("GetNodeCount").Return(int64(1), int64(0), int64(0), nil).
			Run(func(mock.Arguments) { pushed <- struct{}{} }).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "http://%zz", nil, nil,
			orgId, nil, nil)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		s.StartMetricsPusher(ctx)

		_, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
			NodeId: ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE).String(),
		})
		assert.NoError(t, err)

		select {
		case <-pushed:
		case <-time.After(2 * time.Second):
			t.Fatal("pusher did not push node counts after a write")
		}
	})
}

type publishFixture struct {
	nodeRepo       *mocks.NodeRepo
	nodeStatusRepo *mocks.NodeStatusRepo
	siteRepo       *mocks.SiteRepo
	siteService    *mocks.SiteClientProvider
}

func TestNodeServer_PublishFailureKeepsSuccess(t *testing.T) {
	for _, tc := range []struct {
		name string
		call func(f publishFixture, s *server.NodeServer) error
	}{
		{"AddNode", func(f publishFixture, s *server.NodeServer) error {
			nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
			f.nodeRepo.On("Add", mock.Anything, mock.Anything).Return(nil).Once()

			_, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
				NodeId: nodeId.StringLowercase(),
				Name:   "node-A",
			})
			return err
		}},
		{"UpdateNodeStatus", func(f publishFixture, s *server.NodeServer) error {
			nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
			f.nodeStatusRepo.On("Update", mock.Anything, "connectivity").Return(nil).Once()
			f.nodeRepo.On("Get", nodeId).
				Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()

			_, err := s.UpdateNodeStatus(context.TODO(), &pb.UpdateNodeStateRequest{
				NodeId:       nodeId.StringLowercase(),
				Connectivity: ukama.NodeConnectivityOffline.String(),
			})
			return err
		}},
		{"UpdateNode", func(f publishFixture, s *server.NodeServer) error {
			nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
			f.nodeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			f.nodeRepo.On("Get", nodeId).
				Return(&db.Node{Id: nodeId.StringLowercase(), Name: "node-new"}, nil).Once()

			_, err := s.UpdateNode(context.TODO(), &pb.UpdateNodeRequest{
				NodeId: nodeId.StringLowercase(),
				Name:   "node-new",
			})
			return err
		}},
		{"DeleteNode", func(f publishFixture, s *server.NodeServer) error {
			nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
			f.nodeRepo.On("Delete", nodeId, mock.Anything).Return(nil).Once()

			_, err := s.DeleteNode(context.TODO(),
				&pb.DeleteNodeRequest{NodeId: nodeId.StringLowercase()})
			return err
		}},
		{"AttachNodes", func(f publishFixture, s *server.NodeServer) error {
			nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)
			ampId := ukama.NewVirtualAmplifierNodeId()
			f.nodeRepo.On("AttachNodes", nodeId, []string{ampId.StringLowercase()}).
				Return(nil).Once()

			_, err := s.AttachNodes(context.TODO(), &pb.AttachNodesRequest{
				NodeId:        nodeId.StringLowercase(),
				AttachedNodes: []string{ampId.StringLowercase()},
			})
			return err
		}},
		{"DetachNode", func(f publishFixture, s *server.NodeServer) error {
			nodeId := ukama.NewVirtualAmplifierNodeId()
			f.nodeRepo.On("Get", nodeId).
				Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
			f.nodeRepo.On("DetachNode", nodeId).Return(nil).Once()

			_, err := s.DetachNode(context.TODO(),
				&pb.DetachNodeRequest{NodeId: nodeId.StringLowercase()})
			return err
		}},
		{"AddNodeToSite", func(f publishFixture, s *server.NodeServer) error {
			nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)
			siteId := uuid.NewV4()
			networkId := uuid.NewV4()
			siteClient := &sitemocks.SiteServiceClient{}
			f.siteService.On("GetClient").Return(siteClient, nil).Once()
			siteClient.On("Get", mock.Anything,
				&sitepb.GetRequest{SiteId: siteId.String()}, mock.Anything).
				Return(&sitepb.GetResponse{
					Site: &sitepb.Site{NetworkId: networkId.String()},
				}, nil).Once()
			f.siteRepo.On("AssignNodes", siteId, networkId,
				[]string{nodeId.StringLowercase()}).
				Return([]db.Site{{NodeId: nodeId.StringLowercase(), SiteId: siteId,
					NetworkId: networkId}}, nil).Once()

			_, err := s.AddNodeToSite(context.TODO(), &pb.AddNodeToSiteRequest{
				NodeId:    nodeId.StringLowercase(),
				SiteId:    siteId.String(),
				NetworkId: networkId.String(),
			})
			return err
		}},
		{"ReleaseSite", func(f publishFixture, s *server.NodeServer) error {
			siteId := uuid.NewV4()
			f.siteRepo.On("ReleaseSite", siteId).Return([]db.Site{{
				NodeId:    ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE).StringLowercase(),
				SiteId:    siteId,
				NetworkId: uuid.NewV4(),
			}}, nil).Once()

			_, err := s.ReleaseSite(context.TODO(),
				&pb.ReleaseSiteRequest{SiteId: siteId.String()})
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := publishFixture{
				nodeRepo:       &mocks.NodeRepo{},
				nodeStatusRepo: &mocks.NodeStatusRepo{},
				siteRepo:       &mocks.SiteRepo{},
				siteService:    &mocks.SiteClientProvider{},
			}
			msgbusClient := &mbmocks.MsgBusServiceClient{}
			msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
				Return(errors.New("msgbus unavailable")).Once()

			s := server.NewNodeServer(OrgName, f.nodeRepo, f.siteRepo, f.nodeStatusRepo,
				"", msgbusClient, f.siteService, orgId, nil, nil)

			assert.NoError(t, tc.call(f, s))
			msgbusClient.AssertExpectations(t)
			f.nodeRepo.AssertExpectations(t)
			f.nodeStatusRepo.AssertExpectations(t)
			f.siteRepo.AssertExpectations(t)
			f.siteService.AssertExpectations(t)
		})
	}
}

func TestNodeServer_RepoErrorsMapToGrpcCodes(t *testing.T) {
	towerId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)
	ampId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_AMPNODE).StringLowercase()

	errCases := []struct {
		name string
		err  error
		code codes.Code
	}{
		{"Invalid", fmt.Errorf("not an amplifier: %w", store.ErrInvalid),
			codes.InvalidArgument},
		{"Precondition", fmt.Errorf("not attached: %w", store.ErrPrecondition),
			codes.FailedPrecondition},
		{"NotFound", fmt.Errorf("node: %w", store.ErrNotFound), codes.NotFound},
		{"Other", errors.New("boom"), codes.Internal},
	}

	rpcs := []struct {
		name string
		call func(s *server.NodeServer, r *mocks.NodeRepo, repoErr error) error
	}{
		{"AttachNodes", func(s *server.NodeServer, r *mocks.NodeRepo, repoErr error) error {
			r.On("AttachNodes", towerId, []string{ampId}).Return(repoErr).Once()
			_, err := s.AttachNodes(context.TODO(), &pb.AttachNodesRequest{
				NodeId:        towerId.StringLowercase(),
				AttachedNodes: []string{ampId},
			})

			return err
		}},
		{"DetachNode", func(s *server.NodeServer, r *mocks.NodeRepo, repoErr error) error {
			r.On("Get", towerId).Return(&db.Node{Id: towerId.StringLowercase()}, nil).Once()
			r.On("DetachNode", towerId).Return(repoErr).Once()
			_, err := s.DetachNode(context.TODO(),
				&pb.DetachNodeRequest{NodeId: towerId.StringLowercase()})

			return err
		}},
		{"DeleteNode", func(s *server.NodeServer, r *mocks.NodeRepo, repoErr error) error {
			r.On("Delete", towerId, mock.Anything).Return(repoErr).Once()
			_, err := s.DeleteNode(context.TODO(),
				&pb.DeleteNodeRequest{NodeId: towerId.StringLowercase()})

			return err
		}},
	}

	for _, rpc := range rpcs {
		for _, tc := range errCases {
			t.Run(rpc.name+"/"+tc.name, func(t *testing.T) {
				nodeRepo := &mocks.NodeRepo{}
				s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil,
					orgId, nil, nil)

				err := rpc.call(s, nodeRepo, tc.err)

				assert.Equal(t, tc.code, status.Code(err))
				nodeRepo.AssertExpectations(t)
			})
		}
	}
}
