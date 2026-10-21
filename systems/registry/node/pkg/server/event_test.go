/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/mocks"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"github.com/ukama/ukama/systems/registry/node/pkg/server"

	evt "github.com/ukama/ukama/systems/common/events"
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
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
		nodeStatusRepo.On("Update", mock.Anything, "connectivity", "state").Return(nil).Once()

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
		nodeStatusRepo := &mocks.NodeStatusRepo{}
		nodeStatusRepo.On("Update", mock.Anything, "state").Return(gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, orgId,
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
		nodeStatusRepo.AssertExpectations(t)
	})

	t.Run("UnknownNodeIgnored", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeStatusRepo.On("Update", mock.Anything, "connectivity", "state").
			Return(gorm.ErrRecordNotFound).Once()

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
		nodeStatusRepo.AssertExpectations(t)
		nodeRepo.AssertNotCalled(t, "Get", mock.Anything)
		nodeRepo.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})

	t.Run("NothingToWriteIsRejected", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.NodeStateChangeEvent{
				NodeId:   nodeId.StringLowercase(),
				Substate: "unknown",
			}),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeStatusRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("NoSubstateWritesStateOnly", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
		nodeStatusRepo.On("Update", mock.MatchedBy(func(ns *db.NodeStatus) bool {
			return ns.State == ukama.NodeStateConfiguring
		}), "state").Return(nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
			nil, orgId, nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.NodeStateChangeEvent{
				NodeId: nodeId.StringLowercase(),
				State:  "Configuring",
			}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
		nodeStatusRepo.AssertExpectations(t)
	})

	t.Run("SubstateOff", func(t *testing.T) {
		nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).
			Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
		nodeStatusRepo.On("Update", mock.Anything, "connectivity", "state").Return(nil).Once()

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

		nodeStatusRepo.On("Update", mock.Anything, "connectivity", "state").Return(gorm.ErrInvalidDB).Once()

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

	t.Run("AlreadyRegisteredIsSuccess", func(t *testing.T) {
		nodeId := ukama.NewVirtualHomeNodeId()

		nodeRepo := &mocks.NodeRepo{}
		nodeRepo.On("Add", mock.Anything, mock.Anything).
			Return(&pgconn.PgError{Code: "23505"}).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nil, "", nil, nil, orgId,
			nil, nil)
		eventServer := server.NewNodeEventServer("testorg", s, nil)

		resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: routingKey,
			Msg: mustAny(t, &epb.EventInventoryNodeComponentAdd{
				PartNumber: nodeId.StringLowercase(),
			}),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		nodeRepo.AssertExpectations(t)
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

		siteRepo := &mocks.SiteRepo{}

		ids := []string{
			towerId.StringLowercase(), ampId.StringLowercase(), ctrlId.StringLowercase(),
		}
		siteRepo.On("AssignNodes", siteId, networkId, ids).Return([]db.Site{}, nil).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil,
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
		siteRepo.AssertExpectations(t)
	})

	t.Run("AccessNodeNotTower", func(t *testing.T) {
		invClient := &cmocks.ComponentClient{}
		invClient.On("Get", "access-1").
			Return(&cinvent.ComponentInfo{
				PartNumber: ukama.NewVirtualHomeNodeId().StringLowercase(),
			}, nil).Once()

		siteRepo := &mocks.SiteRepo{}

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil,
			orgId, nil, nil)
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
		siteRepo.AssertNotCalled(t, "AssignNodes", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("AssignFailsAsAWhole", func(t *testing.T) {
		towerId := ukama.NewVirtualTowerNodeId()

		invClient := &cmocks.ComponentClient{}
		invClient.On("Get", "access-1").
			Return(&cinvent.ComponentInfo{PartNumber: towerId.StringLowercase()},
				nil).Once()

		siteRepo := &mocks.SiteRepo{}
		ampId, err := ukama.GetANodeIdFromTNodeId(towerId.StringLowercase())
		assert.NoError(t, err)

		siteRepo.On("AssignNodes", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, status.Errorf(codes.NotFound, "node %s not found", ampId)).Once()

		s := server.NewNodeServer(OrgName, nil, siteRepo, nil, "", nil, nil,
			orgId, nil, nil)
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
		assert.Contains(t, err.Error(), ampId.StringLowercase())
		assert.Nil(t, resp)
		siteRepo.AssertExpectations(t)
	})
}

func TestNodeEventServer_MeshConnectivity(t *testing.T) {
	cases := []struct {
		name         string
		event        evt.NodeStateEventId
		msg          proto.Message
		connectivity ukama.NodeConnectivity
	}{
		{"Online", evt.NodeStateEventOnline, &epb.NodeOnlineEvent{}, ukama.NodeConnectivityOnline},
		{"Offline", evt.NodeStateEventOffline, &epb.NodeOfflineEvent{}, ukama.NodeConnectivityOffline},
	}

	for _, tc := range cases {
		routingKey := msgbus.PrepareRoute("testorg", evt.NodeStateEventRoutingKey[tc.event])

		t.Run(tc.name+"WritesConnectivityOnly", func(t *testing.T) {
			nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

			nodeRepo := &mocks.NodeRepo{}
			nodeStatusRepo := &mocks.NodeStatusRepo{}

			nodeRepo.On("Get", nodeId).
				Return(&db.Node{Id: nodeId.StringLowercase()}, nil).Once()
			nodeStatusRepo.On("Update", mock.MatchedBy(func(ns *db.NodeStatus) bool {
				return ns.Connectivity == tc.connectivity
			}), "connectivity").Return(nil).Once()

			s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
				nil, orgId, nil, nil)
			eventServer := server.NewNodeEventServer("testorg", s, nil)

			resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
				RoutingKey: routingKey,
				Msg:        mustAny(t, withNodeId(tc.msg, nodeId.StringLowercase())),
			})

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			nodeRepo.AssertExpectations(t)
			nodeStatusRepo.AssertExpectations(t)
		})

		t.Run(tc.name+"UnknownNodeIgnored", func(t *testing.T) {
			nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

			nodeRepo := &mocks.NodeRepo{}
			nodeStatusRepo := &mocks.NodeStatusRepo{}

			nodeStatusRepo.On("Update", mock.Anything, "connectivity").
				Return(gorm.ErrRecordNotFound).Once()

			s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil,
				nil, orgId, nil, nil)
			eventServer := server.NewNodeEventServer("testorg", s, nil)

			resp, err := eventServer.EventNotification(context.TODO(), &epb.Event{
				RoutingKey: routingKey,
				Msg:        mustAny(t, withNodeId(tc.msg, nodeId.StringLowercase())),
			})

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			nodeStatusRepo.AssertExpectations(t)
			nodeRepo.AssertNotCalled(t, "Get", mock.Anything)
		})
	}
}

func withNodeId(msg proto.Message, nodeId string) proto.Message {
	switch m := proto.Clone(msg).(type) {
	case *epb.NodeOnlineEvent:
		m.NodeId = nodeId
		return m
	case *epb.NodeOfflineEvent:
		m.NodeId = nodeId
		return m
	}

	return msg
}
