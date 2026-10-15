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
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/mocks"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"github.com/ukama/ukama/systems/registry/node/pkg/server"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	cinvent "github.com/ukama/ukama/systems/common/rest/client/inventory"
	node "github.com/ukama/ukama/systems/common/rest/client/node"
)

func mustAny(t *testing.T, msg proto.Message) *anypb.Any {
	a, err := anypb.New(msg)
	assert.NoError(t, err)

	return a
}

func TestNodeEventServer_EventNotification(t *testing.T) {
	orgName := "testorg"

	t.Run("UnknownRoutingKey", func(t *testing.T) {
		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)
		eventServer := server.NewNodeEventServer(orgName, s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: "event.cloud.local.testorg.unknown.event",
			Msg:        mustAny(t, &epb.HealthReportEvent{}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("UnmarshalError", func(t *testing.T) {
		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)
		eventServer := server.NewNodeEventServer(orgName, s, nil)

		routingKey := msgbus.PrepareRoute(orgName,
			"event.cloud.local.{{ .Org}}.node.health.report.store")

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg:        mustAny(t, &epb.EventRegistryNodeCreate{}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("NodeStateTransition_Dispatched", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Twice()
		nodeStatusRepo.On("Update", mock.Anything).Return(nil).Once()
		nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)
		eventServer := server.NewNodeEventServer(orgName, s, nil)

		routingKey := msgbus.PrepareRoute(orgName,
			"event.cloud.local.{{ .Org}}.node.state.node.transition")

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.NodeStateChangeEvent{
				NodeId:   nodeId.StringLowercase(),
				State:    "operational",
				Substate: "on",
			}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
		nodeStatusRepo.AssertExpectations(t)
	})

	t.Run("HealthReport_Dispatched", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(&db.Node{
			Id:   nodeId.StringLowercase(),
			Type: ukama.NODE_ID_TYPE_HOMENODE,
		}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)
		eventServer := server.NewNodeEventServer(orgName, s, nil)

		routingKey := msgbus.PrepareRoute(orgName,
			"event.cloud.local.{{ .Org}}.node.health.report.store")

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.HealthReportEvent{
				NodeId: nodeId.StringLowercase(),
				Id:     "report-1",
			}),
		})

		// Node is not a tower node: handler returns early without error.
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("AddSite_Dispatched", func(t *testing.T) {
		invClient := &cmocks.ComponentClient{}
		invClient.On("Get", "access-1").Return(nil, errors.New("not found")).Once()

		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)
		eventServer := server.NewNodeEventServer(orgName, s, invClient)

		routingKey := msgbus.PrepareRoute(orgName,
			"event.cloud.local.{{ .Org}}.registry.site.site.create")

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.EventAddSite{
				AccessId:  "access-1",
				SiteId:    uuid.NewV4().String(),
				NetworkId: uuid.NewV4().String(),
			}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		invClient.AssertExpectations(t)
	})

	t.Run("InventoryNodeAdd_Dispatched", func(t *testing.T) {
		nodeId := ukama.NewVirtualHomeNodeId()

		nodeRepo := &mocks.NodeRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		nodeRepo.On("Add", mock.Anything, mock.Anything).Return(nil).Once()
		nodeRepo.On("GetNodeCount").Return(int64(1), int64(0), int64(0), nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", msgbusClient, nil,
			orgId, nil, nil)
		eventServer := server.NewNodeEventServer(orgName, s, nil)

		routingKey := msgbus.PrepareRoute(orgName,
			"event.cloud.local.{{ .Org}}.inventory.component.node.added")

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.EventInventoryNodeComponentAdd{
				PartNumber: nodeId.StringLowercase(),
				Type:       "hnode",
			}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeEventServer_handleHealthReportEvent(t *testing.T) {
	routingKey := msgbus.PrepareRoute("testorg",
		"event.cloud.local.{{ .Org}}.node.health.report.store")

	t.Run("UnknownNodeIgnored", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t,
				&epb.HealthReportEvent{NodeId: nodeId.StringLowercase(), Id: "r1"}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("GetInterfacesError", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(&db.Node{
			Id: nodeId.StringLowercase(), Type: ukama.NODE_ID_TYPE_TOWERNODE,
		}, nil).Once()

		healthClient := &cmocks.NodeHealthClient{}
		healthClient.On("GetInterfaces", nodeId.StringLowercase(), "r1").
			Return(node.InterfaceInfo{}, errors.New("unreachable")).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, healthClient)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t,
				&epb.HealthReportEvent{NodeId: nodeId.StringLowercase(), Id: "r1"}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
		healthClient.AssertExpectations(t)
	})

	t.Run("NoGpsLock_Skipped", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(&db.Node{
			Id: nodeId.StringLowercase(), Type: ukama.NODE_ID_TYPE_TOWERNODE,
		}, nil).Once()

		healthClient := &cmocks.NodeHealthClient{}
		healthClient.On("GetInterfaces", nodeId.StringLowercase(), "r1").
			Return(node.InterfaceInfo{
				Gps: &node.GPSInterfaceInfo{
					Lock: false, Available: true, Coordinates: "1,2",
				},
			}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, healthClient)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t,
				&epb.HealthReportEvent{NodeId: nodeId.StringLowercase(), Id: "r1"}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
		healthClient.AssertExpectations(t)
	})

	t.Run("NoGpsAvailable_Skipped", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(&db.Node{
			Id: nodeId.StringLowercase(), Type: ukama.NODE_ID_TYPE_TOWERNODE,
		}, nil).Once()

		healthClient := &cmocks.NodeHealthClient{}
		healthClient.On("GetInterfaces", nodeId.StringLowercase(), "r1").
			Return(node.InterfaceInfo{
				Gps: &node.GPSInterfaceInfo{
					Lock: true, Available: false, Coordinates: "1,2",
				},
			}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, healthClient)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t,
				&epb.HealthReportEvent{NodeId: nodeId.StringLowercase(), Id: "r1"}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
		healthClient.AssertExpectations(t)
	})

	t.Run("SameLocation_Skipped", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(&db.Node{
			Id: nodeId.StringLowercase(), Type: ukama.NODE_ID_TYPE_TOWERNODE,
			Latitude: "1", Longitude: "2",
		}, nil).Once()

		healthClient := &cmocks.NodeHealthClient{}
		healthClient.On("GetInterfaces", nodeId.StringLowercase(), "r1").
			Return(node.InterfaceInfo{
				Gps: &node.GPSInterfaceInfo{
					Lock: true, Available: true, Coordinates: "1,2",
				},
			}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, healthClient)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t,
				&epb.HealthReportEvent{NodeId: nodeId.StringLowercase(), Id: "r1"}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
		healthClient.AssertExpectations(t)
	})

	t.Run("LocationUpdated", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(&db.Node{
			Id: nodeId.StringLowercase(), Type: ukama.NODE_ID_TYPE_TOWERNODE,
		}, nil).Once()
		nodeRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		nodeRepo.On("Get", nodeId).Return(&db.Node{
			Id: nodeId.StringLowercase(), Type: ukama.NODE_ID_TYPE_TOWERNODE,
			Latitude: "1", Longitude: "2",
		}, nil).Once()

		healthClient := &cmocks.NodeHealthClient{}
		healthClient.On("GetInterfaces", nodeId.StringLowercase(), "r1").
			Return(node.InterfaceInfo{
				Gps: &node.GPSInterfaceInfo{
					Lock: true, Available: true, Coordinates: "1,2",
				},
			}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, healthClient)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t,
				&epb.HealthReportEvent{NodeId: nodeId.StringLowercase(), Id: "r1"}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
		healthClient.AssertExpectations(t)
	})

	t.Run("InvalidCoordinates", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(&db.Node{
			Id: nodeId.StringLowercase(), Type: ukama.NODE_ID_TYPE_TOWERNODE,
		}, nil).Once()

		healthClient := &cmocks.NodeHealthClient{}
		healthClient.On("GetInterfaces", nodeId.StringLowercase(), "r1").
			Return(node.InterfaceInfo{
				Gps: &node.GPSInterfaceInfo{
					Lock: true, Available: true, Coordinates: "invalid",
				},
			}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, healthClient)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t,
				&epb.HealthReportEvent{NodeId: nodeId.StringLowercase(), Id: "r1"}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestNodeEventServer_handleNodeStateTransitionEvent(t *testing.T) {
	routingKey := msgbus.PrepareRoute("testorg",
		"event.cloud.local.{{ .Org}}.node.state.node.transition")

	t.Run("UnknownOffboardedIgnored", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.NodeStateChangeEvent{
				NodeId: nodeId.StringLowercase(),
				State:  "Offboarded",
			}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("UnknownNodeCreatedThenUpdated", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Once()
		nodeRepo.On("Add", mock.Anything, mock.Anything).Return(nil).Once()
		nodeRepo.On("GetNodeCount").Return(int64(1), int64(0), int64(0), nil).Once()
		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Twice()
		nodeStatusRepo.On("Update", mock.Anything).Return(nil).Once()
		nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.NodeStateChangeEvent{
				NodeId:   nodeId.StringLowercase(),
				State:    "operational",
				Substate: "on",
			}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
		nodeStatusRepo.AssertExpectations(t)
	})

	t.Run("CreateFailsAndRecheckFails", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Twice()
		nodeRepo.On("Add", mock.Anything, mock.Anything).
			Return(gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.NodeStateChangeEvent{
				NodeId: nodeId.StringLowercase(),
				State:  "operational",
			}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("OtherGetError", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.NodeStateChangeEvent{
				NodeId: nodeId.StringLowercase(),
				State:  "operational",
			}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("SubstateOff", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Twice()
		nodeStatusRepo.On("Update", mock.Anything).Return(nil).Once()
		nodeRepo.On("GetNodeCount").Return(int64(1), int64(0), int64(1), nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.NodeStateChangeEvent{
				NodeId:   nodeId.StringLowercase(),
				State:    "faulty",
				Substate: "off",
			}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
		nodeStatusRepo.AssertExpectations(t)
	})

	t.Run("UpdateStatusError", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
		nodeStatusRepo.On("Update", mock.Anything).Return(gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.NodeStateChangeEvent{
				NodeId:   nodeId.StringLowercase(),
				State:    "faulty",
				Substate: "off",
			}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
		nodeStatusRepo.AssertExpectations(t)
	})
}

func TestNodeEventServer_handleAddNode(t *testing.T) {
	routingKey := msgbus.PrepareRoute("testorg",
		"event.cloud.local.{{ .Org}}.inventory.component.node.added")

	t.Run("InvalidNodeId", func(t *testing.T) {
		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.EventInventoryNodeComponentAdd{
				PartNumber: "not-a-node-id",
			}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("AddNodeError", func(t *testing.T) {
		nodeId := ukama.NewVirtualHomeNodeId()

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Add", mock.Anything, mock.Anything).
			Return(gorm.ErrInvalidDB).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.EventInventoryNodeComponentAdd{
				PartNumber: nodeId.StringLowercase(),
			}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}

func TestNodeEventServer_handleAddNodeToSite(t *testing.T) {
	routingKey := msgbus.PrepareRoute("testorg",
		"event.cloud.local.{{ .Org}}.registry.site.site.create")

	t.Run("ComponentGetError", func(t *testing.T) {
		invClient := &cmocks.ComponentClient{}
		invClient.On("Get", "access-1").Return(nil, errors.New("not found")).Once()

		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)
		eventServer := server.NewNodeEventServer("testorg", s, invClient)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.EventAddSite{
				AccessId:  "access-1",
				SiteId:    uuid.NewV4().String(),
				NetworkId: uuid.NewV4().String(),
			}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		invClient.AssertExpectations(t)
	})

	t.Run("InvalidPartNumber", func(t *testing.T) {
		invClient := &cmocks.ComponentClient{}
		invClient.On("Get", "access-1").
			Return(&cinvent.ComponentInfo{PartNumber: "bad"}, nil).Once()

		s := server.NewNodeServer(OrgName, nil, nil, nil, "", nil, nil, orgId, nil,
			nil)
		eventServer := server.NewNodeEventServer("testorg", s, invClient)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.EventAddSite{
				AccessId:  "access-1",
				SiteId:    uuid.NewV4().String(),
				NetworkId: uuid.NewV4().String(),
			}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		invClient.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		towerId := ukama.NewVirtualTowerNodeId()
		ampId, err := ukama.GetANodeIdFromTNodeId(towerId.StringLowercase())
		assert.NoError(t, err)
		ctrlId, err := ukama.GetCNodeIdFromTNodeId(towerId.StringLowercase())
		assert.NoError(t, err)

		siteId := uuid.NewV4()
		networkId := uuid.NewV4()

		invClient := &cmocks.ComponentClient{}
		invClient.On("Get", "access-1").
			Return(&cinvent.ComponentInfo{PartNumber: towerId.StringLowercase()},
				nil).Once()

		nodeRepo := &mocks.NodeRepo{}
		siteRepo := &mocks.SiteRepo{}

		ids := []string{
			towerId.StringLowercase(), ampId.StringLowercase(), ctrlId.StringLowercase(),
		}
		for _, id := range ids {
			nodeRepo.On("Get", mock.MatchedBy(func(n ukama.NodeID) bool {
				return n.StringLowercase() == id
			})).Return(&db.Node{Id: id}, nil).Once()
		}
		siteRepo.On("AddNode", mock.Anything, mock.Anything).Return(nil).Times(3)

		s := server.NewNodeServer(OrgName, nodeRepo, siteRepo, nil, "", nil, nil,
			orgId, nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, invClient)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.EventAddSite{
				AccessId:  "access-1",
				SiteId:    siteId.String(),
				NetworkId: networkId.String(),
			}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		invClient.AssertExpectations(t)
		nodeRepo.AssertExpectations(t)
		siteRepo.AssertExpectations(t)
	})
}
