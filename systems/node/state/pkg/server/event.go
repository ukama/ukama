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

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	stm "github.com/ukama/ukama/systems/common/stateMachine"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StateEventServer struct {
	orgName      string
	orgId        string
	stateMachine *stm.StateMachine
	instances    map[string]*stm.StateMachineInstance
	s            *StateServer
	configPath   string
	epb.UnimplementedEventNotificationServiceServer
}

func NewStateEventServer(orgName string, orgId string, s *StateServer, stateMachine *stm.StateMachine, instances map[string]*stm.StateMachineInstance, configPath string) *StateEventServer {
	server := &StateEventServer{
		orgName:      orgName,
		orgId:        orgId,
		instances:    make(map[string]*stm.StateMachineInstance),
		stateMachine: stateMachine,
		s:            s,
		configPath:   configPath,
	}

	server.stateMachine = stm.NewStateMachine(server.handleTransition)

	return server
}

func (n *StateEventServer) handleTransition(event stm.Event) {
	log.Infof("Transition occurred: %+v", event)

	var subState string
	if event.IsSubstate {
		subState = event.NewState
		log.Infof("Substate Transition: Event: %s, Old Substate: %s, New Substate: %s\n", event.Name, event.OldState, event.NewState)
	}
	ctx := context.Background()
	_, err := n.s.AddNodeState(ctx, &pb.AddStateRequest{
		NodeId:       event.InstanceId,
		CurrentState: event.NewState,
		SubState:     subState,
		Events:       []string{event.Name},
	})
	if err != nil {
		log.Errorf("Error adding node state: %v", err)
	}
}

func (n *StateEventServer) getOrCreateInstance(nodeId, initialState string) (*stm.StateMachineInstance, error) {
	instance, exists := n.instances[nodeId]
	if !exists {
		newInstance, err := n.stateMachine.NewInstance(n.configPath, nodeId, initialState)
		if err != nil {
			return nil, fmt.Errorf("failed to create new instance: %v", err)
		}
		n.instances[nodeId] = newInstance
		instance = newInstance
	}
	return instance, nil
}

func (n *StateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, e.Msg)

	var nodeId string
	var eventName string
	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventCreate]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventCreate]
		msg, err := epb.UnmarshalEventRegistryNodeCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, msg)
		nodeId = msg.NodeId
		eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventAssign]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventAssign]
		msg, err := epb.UnmarshalEventRegistryNodeAssign(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, msg)
		nodeId = msg.NodeId
		eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOffline]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventOffline]
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, msg)
		nodeId = msg.NodeId
		eventName = c.Name

		//  case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventConfig]):
		// 	 c := evt.NodeEventToEventConfig[evt.NodeStateEventConfig]
		// 	 msg, err := epb.UnmarshalNodeConfigUpdateEvent(e.Msg, c.Name)
		// 	 if err != nil {
		// 		 return nil, err
		// 	 }
		// 	 nodeId = msg.NodeId
		// 	 eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventOnline]
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, msg)
		nodeId = msg.NodeId
		eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventRelease]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventRelease]
		msg, err := epb.UnmarshalEventRegistryNodeRelease(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		nodeId = msg.NodeId
		eventName = c.Name

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
		return &epb.EventResponse{}, nil
	}

	err := n.ProcessEvent(eventName, nodeId)
	if err != nil {
		log.Errorf("Error processing event: %v", err)
		return nil, err
	}

	return &epb.EventResponse{}, nil
}

func (n *StateEventServer) ProcessEvent(eventName, nodeId string) error {
	log.Infof("Processing event %s for node %s", eventName, nodeId)

	ctx := context.Background()

	latestState, err := n.s.GetLatestState(ctx, &pb.GetLatestStateRequest{NodeId: nodeId})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				log.Errorf("Invalid node ID format: %v", err)
				return err
			case codes.Internal:
				log.Errorf("Internal error while checking node state: %v", err)
				return err
			}
		}
		return err
	}

	currentState := "unknown"
	if latestState != nil && latestState.State != nil {
		currentState = latestState.State.CurrentState
		log.Infof("Node %s already exists with state %s. Attempting transition.", nodeId, currentState)
	} else {
		log.Infof("Node %s does not exist. Creating new node.", nodeId)
	}

	instance, err := n.getOrCreateInstance(nodeId, currentState)
	if err != nil {
		log.Errorf("Failed to create state machine instance for node %s: %v", nodeId, err)
		return err
	}

	instance.Transition(eventName)

	log.Infof("Successfully processed event for node %s with state '%s'", nodeId, instance.CurrentState)
	return nil
}
