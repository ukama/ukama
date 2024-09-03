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
	"encoding/json"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	eCfgPb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	utils "github.com/ukama/ukama/systems/node/state/pkg/utils"

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
		"event.cloud.local.{{ .Org}}.registry.node.node.create":      n.handleNodeCreateEvent,
		"event.cloud.local.{{ .Org}}.messaging.mesh.node.online":     n.handleNodeOnlineEvent,
		"event.cloud.local.{{ .Org}}.registry.node.node.assign":      n.handleOnboardingEvent,
		"event.cloud.local.{{ .Org}}.messaging.mesh.node.offline":    n.handleNodeOfflineEvent,
		"event.cloud.local.{{ .Org}}.node.notify.notification.store": n.handleNodeHealthSeverityEvent,
		"event.cloud.local.{{ .Org}}.registry.node.node.release":     n.handleNodeDeassignEvent,
		"event.cloud.local.{{ .Org}}node.configurator.config.add":    n.handleNodeConfigEvent,
	}

	if handler, ok := routingHandlers[e.RoutingKey]; ok {
		if err := handler(e.Msg); err != nil {
			log.Errorf("Error handling event %s: %v", e.RoutingKey, err)
			return nil, err
		}
	} else {
		log.Errorf("No handler for routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *NodeStateEventServer) handleNodeHealthSeverityEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNotification(msg)
	if err != nil {
		return err
	}

	if evt.Type != string(utils.Event) {
		log.Infof("Ignoring message with low severity for node %s", evt.NodeId)
		return nil
	}

	var details map[string]interface{}
	if err = json.Unmarshal(evt.Details, &details); err != nil {
		log.Errorf("Failed to unmarshal details: %v", err)
		return err
	}

	nId, err := ukama.ValidateNodeId(evt.NodeId)
	if err != nil {
		log.Errorf("Error converting NodeId %s to ukama.NodeID. Error: %+v", evt.NodeId, err)
		return err
	}

	currentState, err := n.s.sRepo.GetByNodeId(nId)
	if err != nil {
		log.Errorf("Error getting current state for node %s. Error: %+v", nId, err)
		return err
	}

	var newState ukama.NodeStateEnum
	severity := utils.ToSeverityType(evt.Severity)

	if reboot, ok := details["reboot"].(bool); ok && reboot {
		log.Infof("Received reboot event for node %s. Maintaining current state: %s", nId, currentState.State)
		newState = currentState.State
	} else {
		if severity == utils.Critical {
			newState = ukama.StateFaulty
		} else if configKey, ok := details["config"].(string); ok {
			switch configKey {
			case "ready":
				newState = ukama.StateOperational
			case "update":
				newState = ukama.StateConfigure
			default:
				newState = utils.GetNodeStateBySeverity(severity)
			}
		} else {
			newState = utils.GetNodeStateBySeverity(severity)
		}
	}

	log.Infof("Node %s: Severity: %s, Details: %v, New State: %s", evt.NodeId, evt.Severity, details, newState)

	if err = n.updateNodeState(evt.NodeId, newState, severity, utils.Event); err != nil {
		log.Errorf("Failed to update node state: %v", err)
		return err
	}
	return nil
}

func (n *NodeStateEventServer) handleNodeCreateEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNodeCreateEvent(msg)
	if err != nil {
		return err
	}
	return n.updateNodeState(evt.NodeId, ukama.StateUnknown, utils.Medium, utils.Event)
}

func (n *NodeStateEventServer) handleNodeConfigEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNodeConfigCreateEvent(msg)
	if err != nil {
		return err
	}
	return n.updateNodeState(evt.NodeId, ukama.StateConfigure, utils.Medium, utils.Event)
}

func (n *NodeStateEventServer) handleNodeOnlineEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNodeOnlineEvent(msg)
	if err != nil {
		return err
	}
	return n.updateNodeStateWithConnectivity(evt.NodeId, ukama.Online, utils.Event)
}

func (n *NodeStateEventServer) handleOnboardingEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalOnboardingEvent(msg)
	if err != nil {
		return err
	}
	return n.updateNodeState(evt.NodeId, ukama.StateConfigure, utils.Medium, utils.Event)
}

func (n *NodeStateEventServer) handleNodeOfflineEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNodeOfflineEvent(msg)
	if err != nil {
		return err
	}
	return n.updateNodeStateWithConnectivity(evt.NodeId, ukama.Offline, utils.Event)
}

func (n *NodeStateEventServer) handleNodeDeassignEvent(msg *anypb.Any) error {
	evt, err := n.unmarshalNodeDeassignEvent(msg)
	if err != nil {
		return err
	}
	return n.updateNodeState(evt.NodeId, ukama.StateUnknown, utils.Medium, utils.Event)
}

func (n *NodeStateEventServer) updateNodeState(nodeId string, state ukama.NodeStateEnum, severity utils.SeverityType, eventType utils.NotificationType, connectivity ...*ukama.Connectivity) error {
	nId, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		log.Errorf("Error converting NodeId %s to ukama.NodeID. Error: %+v", nodeId, err)
		return err
	}

	_, err = n.s.sRepo.GetByNodeId(nId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("Error getting latest state for node %s. Error: %+v", nId, err)
		return err
	}

	var connState ukama.Connectivity
	if len(connectivity) > 0 && connectivity[0] != nil {
		connState = *connectivity[0]
	} else {
		connState = ukama.Unknown
	}

	now := time.Now()
	newState := &db.State{
		Id:              uuid.NewV4(),
		NodeId:          nId.String(),
		State:           state,
		Type:            string(eventType),
		LastStateChange: now,
		LastHeartbeat:   now,
		CreatedAt:       now,
		UpdatedAt:       now,
		Connectivity:    connState,
	}

	if _, err = n.s.Create(context.Background(), &pb.CreateStateRequest{State: convertStateToProto(newState)}); err != nil {
		return err
	}

	n.publishNodeStateUpdate(newState, severity, eventType)
	log.Infof("Updated node %s state to %s", nId, state)
	return nil
}

func (n *NodeStateEventServer) updateNodeStateWithConnectivity(nodeId string, connState ukama.Connectivity, eventType utils.NotificationType) error {
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

	return n.updateNodeState(nodeId, currentState.State, utils.Medium, eventType, &connState)
}

func (n *NodeStateEventServer) publishNodeStateUpdate(newState *db.State, severity utils.SeverityType, eventType utils.NotificationType) {
	if n.s.msgbus != nil {
		route := n.s.stateRoutingKey.SetAction("state").SetObject("node").MustBuild()
		evt := &epb.NodeStateHealthEvent{
			NodeId:          newState.NodeId,
			CurrentState:    newState.State.String(),
			LastStateChange: newState.LastStateChange.String(),
			Severity:        string(severity),
			Type:            string(eventType),
		}

		if err := n.s.msgbus.PublishRequest(route, evt); err != nil {
			log.Errorf("Failed to publish node state update event: %+v with key %+v. Error: %s", evt, route, err.Error())
		} else {
			log.Infof("Published node state update event for node %s", newState.NodeId)
		}
	}
}

func (n *NodeStateEventServer) unmarshalNotification(msg *anypb.Any) (*epb.Notification, error) {
	evt := &epb.Notification{}
	if err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}); err != nil {
		log.Errorf("Failed to Unmarshal Node message: %+v. Error: %s.", msg, err.Error())
		return nil, err
	}
	return evt, nil
}

func (n *NodeStateEventServer) unmarshalNodeEvent(msg *anypb.Any, evt proto.Message) error {
	if err := anypb.UnmarshalTo(msg, evt, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}); err != nil {
		log.Errorf("Failed to Unmarshal message with : %+v. Error %s.", msg, err.Error())
		return err
	}
	return nil
}

func (n *NodeStateEventServer) unmarshalNodeOnlineEvent(msg *anypb.Any) (*epb.NodeOnlineEvent, error) {
	evt := &epb.NodeOnlineEvent{}
	return evt, n.unmarshalNodeEvent(msg, evt)
}

func (n *NodeStateEventServer) unmarshalNodeOfflineEvent(msg *anypb.Any) (*epb.NodeOfflineEvent, error) {
	evt := &epb.NodeOfflineEvent{}
	return evt, n.unmarshalNodeEvent(msg, evt)
}

func (n *NodeStateEventServer) unmarshalNodeConfigCreateEvent(msg *anypb.Any) (*eCfgPb.NodeConfigUpdateEvent, error) {
	evt := &eCfgPb.NodeConfigUpdateEvent{}
	return evt, n.unmarshalNodeEvent(msg, evt)
}

func (n *NodeStateEventServer) unmarshalNodeDeassignEvent(msg *anypb.Any) (*epb.NodeReleasedEvent, error) {
	evt := &epb.NodeReleasedEvent{}
	return evt, n.unmarshalNodeEvent(msg, evt)
}

func (n *NodeStateEventServer) unmarshalOnboardingEvent(msg *anypb.Any) (*epb.EventRegistryNodeAssign, error) {
	evt := &epb.EventRegistryNodeAssign{}
	return evt, n.unmarshalNodeEvent(msg, evt)
}

func (n *NodeStateEventServer) unmarshalNodeCreateEvent(msg *anypb.Any) (*epb.EventRegistryNodeCreate, error) {
	evt := &epb.EventRegistryNodeCreate{}
	return evt, n.unmarshalNodeEvent(msg, evt)
}
