/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	eCfgPb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/state/pkg"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
)

type NodeStateEventServer struct {
	s               *StateServer
	orgName         string
	msgbus          mb.MsgBusServiceClient
	stateRoutingKey msgbus.RoutingKeyBuilder
	epb.UnimplementedEventNotificationServiceServer
}

func NewControllerEventServer(orgName string, s *StateServer, msgBus mb.MsgBusServiceClient) *NodeStateEventServer {
	return &NodeStateEventServer{
		s:               s,
		orgName:         orgName,
		msgbus:          msgBus,
		stateRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (n *NodeStateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	routingHandlers := map[string]func(*anypb.Any) error{
		"event.cloud.local.{{ .Org}}.registry.node.node.create":         n.handleNodeCreateEvent,
		"event.cloud.local.{{ .Org}}.messaging.mesh.node.online":        n.handleNodeOnlineEvent,
		"event.cloud.local.{{ .Org}}.registry.node.node.assign":         n.handleOnboardingEvent,
		"event.cloud.local.{{ .Org}}.messaging.mesh.node.offline":       n.handleNodeOfflineEvent,
		"event.cloud.local.{{ .Org}}.node.notify.notification.store":    n.handleNodeHealthSeverityHighEvent,
		"event.cloud.local.{{ .Org}}.node.notify.notification.config.ready": n.handleNodeConfigReadyEvent,
		"event.cloud.local.{{ .Org}}.registry.node.node.release":        n.handleNodeDeassignEvent,
	}

	if handler, ok := routingHandlers[e.RoutingKey]; ok {
		err := handler(e.Msg)
		if err != nil {
			log.Errorf("Error handling event %s: %v", e.RoutingKey, err)
			return nil, err
		}
	} else {
		log.Errorf("No handler for routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *NodeStateEventServer) handleNodeHealthSeverityHighEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNotification(msg)
	if err != nil {
		return err
	}
	if evt.Severity != "high" {
		log.Infof("Ignoring message with low severity for node %s", evt.NodeId)
		return nil
	}
	return n.updateNodeState(evt.NodeId, ukama.StateFaulty)
}

func (n *NodeStateEventServer) handleNodeCreateEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNodeCreateEvent(msg)
	if err != nil {
		return err
	}
	return n.updateNodeState(evt.NodeId, ukama.StateUnknown)
}

func (n *NodeStateEventServer) handleNodeOnlineEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNodeOnlineEvent(msg)
	if err != nil {
		return err
	}
	nId, err := ukama.ValidateNodeId(evt.NodeId)
	if err != nil {
		log.Errorf("Error converting NodeId %s to ukama.NodeID. Error: %+v", nId, err)
		return err
	}

	currentState, err := n.s.sRepo.GetByNodeId(nId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("Error getting latest state for node %s. Error: %+v", nId, err)
		return err
	}
	conn := db.Online

	return n.updateNodeState(evt.NodeId, currentState.State,&conn)
}

func (n *NodeStateEventServer) handleOnboardingEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalOnboardingEvent(msg)
	if err != nil {
		return err
	}
	return n.updateNodeState(evt.NodeId, ukama.StateConfigure)
}

func (n *NodeStateEventServer) handleNodeOfflineEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNodeOfflineEvent(msg)
	if err != nil {
		return err
	}
	nId, err := ukama.ValidateNodeId(evt.NodeId)
	if err != nil {
		log.Errorf("Error converting NodeId %s to ukama.NodeID. Error: %+v", nId, err)
		return err
	}

	currentState, err := n.s.sRepo.GetByNodeId(nId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("Error getting latest state for node %s. Error: %+v", nId, err)
		return err
	}
	conn := db.Offline
	return n.updateNodeState(evt.NodeId,currentState.State,&conn)
}

func (n *NodeStateEventServer) handleNodeConfigReadyEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNodeConfigReadyEvent(msg)
	if err != nil {
		return err
	}
	return n.updateNodeState(evt.NodeId, ukama.StateOperational)
}

func (n *NodeStateEventServer) handleNodeDeassignEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNodeDeassignEvent(msg)
	if err != nil {
		return err
	}
	return n.updateNodeState(evt.NodeId, ukama.StateUnknown)
}

func (n *NodeStateEventServer) updateNodeState(nodeId string, state ukama.NodeStateEnum,connectivity ...*db.Connectivity) error {
	// Validate the Node ID
	nId, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		log.Errorf("Error converting NodeId %s to ukama.NodeID. Error: %+v", nodeId, err)
		return err
	}

	currentState, err := n.s.sRepo.GetByNodeId(nId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("Error getting latest state for node %s. Error: %+v", nId, err)
		return err
	}

	var connState db.Connectivity
	if len(connectivity) > 0 && connectivity[0] != nil {
		connState = *connectivity[0]
	} else {
		connState = db.Unknown 
	}

	// Set the current time for state updates
	now := time.Now()
	// Prepare the new state with optional connectivity status
	newState := &db.State{
		Id:              uuid.NewV4(),
		NodeId:          nId.String(),
		State:           state,
		Type:            currentState.Type,
		LastStateChange: now,
		LastHeartbeat:   now,
		CreatedAt:       now,
		UpdatedAt:       now,
		Connectivity:   connState, 
	}

	// Create the new state record in the repository
	if err := n.s.sRepo.Create(newState, nil); err != nil {
		log.Errorf("Error creating new state for node %s in Nodestate repo. Error: %+v", nId, err)
		return err
	}

	// Publish a state update event if the message bus is available
	if n.s.msgbus != nil {
		route := n.s.stateRoutingKey.SetAction("state").SetObject("node").MustBuild()
		evt := &epb.EventNodeStateUpdate{
			NodeId:          newState.NodeId,
			CurrentState:    newState.State.String(),
			LastStateChange: newState.LastStateChange.String(),
		}

		if err := n.s.msgbus.PublishRequest(route, evt); err != nil {
			log.Errorf("Failed to publish node state update event: %+v with key %+v. Error: %s", evt, route, err.Error())
		} else {
			log.Infof("Published node state update event for node %s", nId.String())
		}
	}

	log.Infof("Updated node %s state to %s", nId, state)
	return nil
}


func (n *NodeStateEventServer) unmarshalNotification(msg *anypb.Any) (*epb.Notification, error) {
	evt := &epb.Notification{}
	err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal Node message: %+v. Error: %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}

func (n *NodeStateEventServer) unmarshalNodeOnlineEvent(msg *anypb.Any) (*epb.NodeOnlineEvent, error) {
	evt := &epb.NodeOnlineEvent{}
	err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOnline message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}

func (n *NodeStateEventServer) unmarshalNodeOfflineEvent(msg *anypb.Any) (*epb.NodeOfflineEvent, error) {
	evt := &epb.NodeOfflineEvent{}
	err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOffline message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}
func (n *NodeStateEventServer) unmarshalNodeConfigCreateEvent(msg *anypb.Any) (*eCfgPb.NodeConfigUpdateEvent, error) {
	evt := &eCfgPb.NodeConfigUpdateEvent{}
	err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeConfigCreate message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}
func (n *NodeStateEventServer) unmarshalNodeConfigReadyEvent(msg *anypb.Any) (*epb.NotificationNodeConfigReady, error) {
	evt := &epb.NotificationNodeConfigReady{}
	err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal Node create message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}
func (n *NodeStateEventServer) unmarshalNodeDeassignEvent(msg *anypb.Any) (*epb.NodeReleasedEvent, error) {
	evt := &epb.NodeReleasedEvent{}
	err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal Node release from site message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}

func (n *NodeStateEventServer) unmarshalNodeHealthSeverityHighEvent(msg *anypb.Any) (*epb.Notification, error) {
	evt := &epb.Notification{}
	err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal Node severity high message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}


func (n *NodeStateEventServer) unmarshalOnboardingEvent(msg *anypb.Any) (*epb.EventRegistryNodeAssign, error) {
	evt := &epb.EventRegistryNodeAssign{}
	err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeCreated message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}


func (n *NodeStateEventServer) unmarshalNodeCreateEvent(msg *anypb.Any) (*epb.EventRegistryNodeCreate, error) {
	evt := &epb.EventRegistryNodeCreate{}
	err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal Node create message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}