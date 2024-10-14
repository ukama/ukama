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
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	stm "github.com/ukama/ukama/systems/common/stateMachine"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

	"github.com/ukama/ukama/systems/node/state/pkg"
)
 
 type StateEventServer struct {
	 s               *StateServer
	 orgName         string
	 msgbus          mb.MsgBusServiceClient
	 stateRoutingKey msgbus.RoutingKeyBuilder
	 epb.UnimplementedEventNotificationServiceServer
	 stateMachine    *stm.StateMachine
	 instances       map[string]*stm.StateMachineInstance
 }
 
 func NewStateEventServer(orgName string, s *StateServer, msgBus mb.MsgBusServiceClient, configFile string) (*StateEventServer, error) {
	 server := &StateEventServer{
		 s:               s,
		 orgName:         orgName,
		 msgbus:          msgBus,
		 stateRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		 instances:       make(map[string]*stm.StateMachineInstance),
	 }
 
	 server.stateMachine = stm.NewStateMachine(server.handleTransition)
 
	 return server, nil
 }
 
 func (n *StateEventServer) handleTransition(event stm.Event) {
	 log.Infof("Transition occurred: %+v", event)
	 // Here you can add any additional logic to handle the transition,
	 // such as updating a database or sending notifications
 }
 
 func (n *StateEventServer) getOrCreateInstance(nodeId, initialState string) (*stm.StateMachineInstance, error) {
	instance, exists := n.instances[nodeId]
	if !exists {
		newInstance, err := n.stateMachine.NewInstance("pkg/nodeState.json", nodeId, initialState)
		if err != nil {
			return nil, fmt.Errorf("failed to create new instance: %v", err)
		}
		n.instances[nodeId] = newInstance
		instance = newInstance
	}
	return instance, nil
}

 func (n *StateEventServer) EventNodeState(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	 log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, e.Msg)
	 
	 var nodeId string
	 var eventName string
 
	 switch e.RoutingKey {
	 case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventCreate]):
		 c := evt.NodeEventToEventConfig[evt.NodeStateEventCreate]
		 msg, err := epb.UnmarshalNodeCreatedEvent(e.Msg, c.Name)
		 if err != nil {
			 return nil, err
		 }
		 nodeId = msg.NodeId
		 eventName = c.Name
 
	 case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventAssign]):
		 c := evt.NodeEventToEventConfig[evt.NodeStateEventAssign]
		 msg, err := epb.UnmarshalEventRegistryNodeAssign(e.Msg, c.Name)
		 if err != nil {
			 return nil, err
		 }
		 nodeId = msg.NodeId
		 eventName = c.Name
 
	 case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOffline]):
		 c := evt.NodeEventToEventConfig[evt.NodeStateEventOffline]
		 msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		 if err != nil {
			 return nil, err
		 }
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
	currentState, err := n.s.GetLatestState(ctx, &pb.GetLatestStateRequest{NodeId: nodeId})
	if err != nil {
		log.Errorf("Failed to get current state for node %s: %v", nodeId, err)
		return err
	}

	var initialState string
	var initialSubstate string
	if currentState == nil {
		initialState = "unknown"
		initialSubstate = ""
	} else {
		initialState = currentState.State.GetCurrentState()
		initialSubstate = currentState.State.GetSubState()
	}

	instance, err := n.getOrCreateInstance(nodeId, initialState)
	if err != nil {
		log.Errorf("Failed to get or create state machine instance for node %s: %v", nodeId, err)
		return err
	}

	// Set the initial substate if it exists
	if initialSubstate != "" {
		instance.CurrentSubstate = initialSubstate
	}

	oldState := instance.CurrentState
	oldSubstate := instance.CurrentSubstate

	instance.Transition(eventName)

	// Check if the state or substate has changed
	stateChanged := oldState != instance.CurrentState
	substateChanged := oldSubstate != instance.CurrentSubstate

	if stateChanged || substateChanged {
		_, err = n.s.AddNodeState(ctx, &pb.AddStateRequest{
			NodeId:       nodeId,
			CurrentState: instance.CurrentState,
			SubState: instance.CurrentSubstate,
			Events:       []string{eventName},
		})
		if err != nil {
			log.Errorf("Error updating node state: %v", err)
			return err
		}

		if stateChanged {
			log.Infof("State transition: %s -> %s", oldState, instance.CurrentState)
		}
		if substateChanged {
			log.Infof("Substate transition: %s -> %s", oldSubstate, instance.CurrentSubstate)
		}
	} else {
		log.Infof("No state or substate change for event %s", eventName)
	}

	log.Infof("Successfully processed event %s for node %s. New state: %s, New substate: %s", 
		eventName, nodeId, instance.CurrentState, instance.CurrentSubstate)
	return nil
}
