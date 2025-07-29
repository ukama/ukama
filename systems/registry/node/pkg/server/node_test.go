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
)

var testNode = ukama.NewVirtualNodeId("HomeNode")
var orgId = uuid.NewV4()

const OrgName = "testorg"

func TestNodeServer_Add(t *testing.T) {
	nodeId := testNode.String()

	msgbusClient := &mbmocks.MsgBusServiceClient{}
	nodeRepo := &mocks.NodeRepo{}
	nodeStatusRepo := &mocks.NodeStatusRepo{}
	siteService := &mocks.SiteClientProvider{}

	const nodeName = "node-A"
	const nodeType = "hnode"

	s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", msgbusClient, siteService, orgId, nil)

	node := &db.Node{
		Id:   nodeId,
		Name: nodeName,
		Type: testNode.GetNodeType(),
		Status: db.NodeStatus{
			NodeId:       nodeId,
			State:        ukama.NodeStateUnknown,
			Connectivity: ukama.NodeConnectivityUndefined,
		},
	}

	nodeRepo.On("Add", node, mock.Anything).Return(nil).Once()
	nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

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
		const nodeName = "node-A"
		const nodeType = ukama.NODE_ID_TYPE_HOMENODE
		var nodeId = ukama.NewVirtualNodeId(nodeType)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).Return(
			&db.Node{Id: nodeId.StringLowercase(),
				Name: nodeName,
				Type: ukama.NODE_ID_TYPE_HOMENODE,
			}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, orgId, nil)

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
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, orgId, nil)

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
		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, orgId, nil)

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
				Type: ntype,
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
		nodeRepo.On("List", nodeId.StringLowercase(), siteId.String(), networkId.String(), ntype, &connectivityVal, &stateVal).
			Return(nodes, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, orgId, nil)

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
		assert.Equal(t, cpb.NodeConnectivity(connectivity), resp.Nodes[0].Status.Connectivity)
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
				Type: ntype,
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
		nodeRepo.On("List", nodeId.StringLowercase(), siteId.String(), networkId.String(), ntype, &connectivityVal, &stateVal).
			Return(nodes, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, orgId, nil)

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
		nodeRepo.On("List", nodeId.StringLowercase(), "", "", "", &connectivityVal, &stateVal).
			Return(nodes, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, orgId, nil)

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

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, orgId, nil)

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

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, orgId, nil)

		// Act
		resp, err := s.List(context.TODO(), &pb.ListRequest{})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}
