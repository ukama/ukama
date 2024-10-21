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
	"fmt"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	stm "github.com/ukama/ukama/systems/common/stateMachine"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type StateEventServer struct {
	orgName      string
	orgId        string
	stateMachine *stm.StateMachine
	instances    map[string]*stm.StateMachineInstance
	instancesMu  sync.RWMutex
	s            *StateServer
	configPath   string
	epb.UnimplementedEventNotificationServiceServer
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	eventBuffer    map[string][]string
	bufferMu       sync.RWMutex
}

func NewStateEventServer(orgName, orgId string, s *StateServer, configPath string, msgBus mb.MsgBusServiceClient) *StateEventServer {
	server := &StateEventServer{
		orgName:        orgName,
		orgId:          orgId,
		instances:      make(map[string]*stm.StateMachineInstance),
		s:              s,
		configPath:     configPath,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		eventBuffer:    make(map[string][]string),
	}

	server.stateMachine = stm.NewStateMachine(server.handleTransition)

	return server
}

func (n *StateEventServer) handleTransition(event stm.Event) {
	log.Infof("Transition event received: %+v ", event)

	var state, substate string
	if event.IsSubstate {
		state = event.NewState
		substate = event.NewSubstate
	} else {
		state = event.NewState
		substate = event.NewSubstate
	}

	n.publishStateChangeEvent(state, substate, event.InstanceID)
}
func (n *StateEventServer) publishStateChangeEvent(state, substate, nodeID string) {
	if n.msgbus == nil {
		return
	}

	route := n.baseRoutingKey.SetActionCreate().SetObject("node").MustBuild()

	evt := &epb.NodeStateChangeEvent{
		NodeId:   nodeID,
		State:    state,
		Substate: substate,
		Events:   n.getEventsForNode(nodeID),
	}

	err := n.msgbus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
	}
}
func (n *StateEventServer) getEventsForNode(nodeID string) []string {
	n.bufferMu.RLock()
	defer n.bufferMu.RUnlock()
	return n.eventBuffer[nodeID]
}

func (n *StateEventServer) addEventToBuffer(nodeID, eventName string) {
	n.bufferMu.Lock()
	defer n.bufferMu.Unlock()
	n.eventBuffer[nodeID] = append(n.eventBuffer[nodeID], eventName)
}

func (n *StateEventServer) clearEventsForNode(nodeID string) {
	n.bufferMu.Lock()
	defer n.bufferMu.Unlock()
	delete(n.eventBuffer, nodeID)
}

func (n *StateEventServer) getOrCreateInstance(nodeID, initialState string) (*stm.StateMachineInstance, error) {
	n.instancesMu.Lock()
	defer n.instancesMu.Unlock()

	instance, exists := n.instances[nodeID]
	if !exists {
		newInstance, err := n.stateMachine.NewInstance(n.configPath, nodeID, initialState)
		if err != nil {
			return nil, fmt.Errorf("failed to create new instance: %w", err)
		}
		n.instances[nodeID] = newInstance
		instance = newInstance
	}
	return instance, nil
}

func (n *StateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {

	log.Infof("Received event %s, routing key %s", e.Msg, e.RoutingKey)
	var nodeId string
	var eventName string
	var body interface{}
	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventAssign]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventAssign]
		msg, err := epb.UnmarshalEventRegistryNodeAssign(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeAssign event: %w", err)
		}
		body = msg
		nodeId = msg.NodeId
		eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOffline]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventOffline]
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeOffline event: %w", err)
		}
		body = msg
		nodeId = msg.NodeId
		eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventOnline]
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeOnline event: %w", err)
		}
		body = msg
		nodeId = msg.NodeId
		eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventRelease]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventRelease]
		msg, err := epb.UnmarshalEventRegistryNodeRelease(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeRelease event: %w", err)
		}
		body = msg
		nodeId = msg.NodeId
		eventName = c.Name
		//To be added once node ready is implemented
	// case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventReady]):
	// 	c := evt.NodeEventToEventConfig[evt.NodeStateEventReady]
	// 	msg, err := epb.UnmarshalNodeReadyEvent(e.Msg, c.Name)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to unmarshal NodeReady event: %w", err)
	// 	}
	// 	body = msg
	// 	nodeID = msg.NodeId
	// 	eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventConfig]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventConfig]
		msg, err := n.unmarshalConfigNodeEvent(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeConfig event: %w", err)
		}
		body = msg
		nodeId = strings.Split(msg.Target, ".")[3]
		eventName = c.Name
	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventUpdate]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventUpdate]
		msg, err := n.unmarshalNodeConfigUpdateEvent(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeUpdate event: %w", err)
		}
		body = msg
		nodeId = msg.NodeId
		eventName = c.Name

	default:
		log.Infof("Received event %s, routing key %s", e.Msg, e.RoutingKey)
		return &epb.EventResponse{}, nil
	}

	if err := n.ProcessEvent(ctx, eventName, nodeId, body); err != nil {
		log.WithError(err).Error("Error processing event")
		return nil, err
	}

	return &epb.EventResponse{}, nil
}

func (n *StateEventServer) ProcessEvent(ctx context.Context, eventName, nodeId string, msg interface{}) error {
	log.Infof("Processing event %s for node %s", eventName, nodeId)

	n.addEventToBuffer(nodeId, eventName)

	latestState, err := n.s.GetLatestState(ctx, &pb.GetLatestStateRequest{NodeId: nodeId})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return fmt.Errorf("invalid node ID format: %w", err)
			case codes.Internal:
				return fmt.Errorf("internal error while checking node state: %w", err)
			}
		}
		return fmt.Errorf("error getting latest state: %w", err)
	}

	currentState := "unknown"
	currentSubstate := ""
	if latestState != nil && latestState.State != nil {
		currentState = latestState.State.CurrentState
		log.Infof("Node already exists with current state %s and substate %s", currentState, currentSubstate)
	} else {
		log.Infof("Creating initial state entry for node %s", nodeId)
		if err := n.createInitialNodeState(ctx, nodeId, eventName, msg); err != nil {
			return err
		}
	}

	instance, err := n.getOrCreateInstance(nodeId, currentState)
	if err != nil {
		return fmt.Errorf("failed to create state machine instance for node %s: %w", nodeId, err)
	}

	prevState := instance.CurrentState
	prevSubstate := instance.CurrentSubstate

	if err := instance.Transition(eventName); err != nil {
		return fmt.Errorf("failed to transition state for node %s: %w", nodeId, err)
	}

	// Check if it's a main state transition
	if instance.CurrentState != prevState {
		// Main state transition
		newSubstate := currentSubstate
		if newSubstate != "" {
			newSubstate += "," + instance.CurrentSubstate
		} else {
			newSubstate = instance.CurrentSubstate
		}
		_, err = n.s.AddNodeState(ctx, &pb.AddStateRequest{
			NodeId:       nodeId,
			CurrentState: instance.CurrentState,
			SubState:     []string{newSubstate},
			Events:       n.getEventsForNode(nodeId),
		})
		if err != nil {
			return fmt.Errorf("failed to add new state for node %s: %w", nodeId, err)
		}
		log.Infof("Added new state for node %s: state=%s, substate=%s", nodeId, instance.CurrentState, newSubstate)
	} else if instance.CurrentSubstate != prevSubstate {
		// Substate transition only
		newSubstate := currentSubstate
		if newSubstate != "" {
			newSubstate += "," + instance.CurrentSubstate
		} else {
			newSubstate = instance.CurrentSubstate
		}
		_, err = n.s.UpdateState(ctx, &pb.UpdateStateRequest{
			NodeId:   nodeId,
			SubState: []string{newSubstate},
			Events:   n.getEventsForNode(nodeId),
		})
		if err != nil {
			return fmt.Errorf("failed to update substate for node %s: %w", nodeId, err)
		}
		log.Infof("Updated substate for node %s: substate=%s", nodeId, newSubstate)
	} else {
		log.Infof("No state or substate change for node %s", nodeId)
	}

	log.Infof("Events for node %s: %v", nodeId, n.getEventsForNode(nodeId))

	n.clearEventsForNode(nodeId)

	return nil
}
func (n *StateEventServer) createInitialNodeState(ctx context.Context, nodeId, eventName string, msg interface{}) error {
	// Assume the initial event will always be online
	onlineEvent, ok := msg.(*epb.NodeOnlineEvent)
	if !ok {
		return fmt.Errorf("expected *NodeOnlineEvent, got %T", msg)
	}

	addStateRequest := &pb.AddStateRequest{
		NodeId:       nodeId,
		CurrentState: "unknown",
		SubState:     []string{"on"},
		Events:       []string{eventName},
		NodeIp:       onlineEvent.NodeIp,
		NodePort:     int32(onlineEvent.NodePort),
		MeshIp:       onlineEvent.MeshIp,
		MeshPort:     int32(onlineEvent.MeshPort),
		MeshHostName: onlineEvent.MeshHostName,
	}

	_, err := n.s.AddNodeState(ctx, addStateRequest)
	if err != nil {
		return fmt.Errorf("failed to create initial state entry for node %s: %w", nodeId, err)
	}
	return nil
}
func (n *StateEventServer) unmarshalConfigNodeEvent(msg *anypb.Any, emsg string) (*npb.NodeFeederMessage, error) {
	p := &npb.NodeFeederMessage{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *StateEventServer) unmarshalNodeConfigUpdateEvent(msg *anypb.Any, emsg string) (*npb.NodeConfigUpdateEvent, error) {
	p := &npb.NodeConfigUpdateEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeConfigUpdate message with : %+v. Error %s. Message %s", msg, err.Error(), emsg)
		return nil, err
	}
	return p, nil
}
