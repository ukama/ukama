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

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg"
)

type NodeStateEventServer struct {
	s               *NodeStateServer
	orgName         string
	msgbus          mb.MsgBusServiceClient
	stateRoutingKey msgbus.RoutingKeyBuilder
	epb.UnimplementedEventNotificationServiceServer
	stateMachine *stm.StateMachine
}

func (s *NodeStateServer) InitializeStateMachine(configPath string) error {
	sm, err := stm.NewStateMachine(configPath)
	if err != nil {
		return err
	}
	s.stateMachine = sm
	log.Infof("Initialized state machine with config from %s", configPath)

	return nil
}
func NewNodeStateEventServer(orgName string, s *NodeStateServer, msgBus mb.MsgBusServiceClient) *NodeStateEventServer {
	return &NodeStateEventServer{
		s:               s,
		orgName:         orgName,
		msgbus:          msgBus,
		stateRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (n *NodeStateEventServer) EventNodeState(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, e.Msg)
	switch e.RoutingKey {

	case msgbus.PrepareRoute(n.orgName, evt.NodeEventRoutingKey[evt.NodeEventCreate]):
		c := evt.NodeEventToEventConfig[evt.NodeEventCreate]
		msg, err := epb.UnmarshalNodeCreatedEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		_ = n.ProcessEvent(&c, msg.NodeId)

	case msgbus.PrepareRoute(n.orgName, evt.NodeEventRoutingKey[evt.NodeEventAssign]):
		c := evt.NodeEventToEventConfig[evt.NodeEventAssign]
		msg, err := epb.UnmarshalEventRegistryNodeAssign(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		_ = n.ProcessEvent(&c, msg.NodeId)
	case msgbus.PrepareRoute(n.orgName, evt.NodeEventRoutingKey[evt.NodeEventOffline]):
		c := evt.NodeEventToEventConfig[evt.NodeEventOffline]
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		_ = n.ProcessEvent(&c, msg.NodeId)
	case msgbus.PrepareRoute(n.orgName, evt.NodeEventRoutingKey[evt.NodeEventOnline]):
		c := evt.NodeEventToEventConfig[evt.NodeEventOnline]
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		_ = n.ProcessEvent(&c, msg.NodeId)

	case msgbus.PrepareRoute(n.orgName, evt.NodeEventRoutingKey[evt.NodeEventRelease]):
		c := evt.NodeEventToEventConfig[evt.NodeEventRelease]
		msg, err := epb.UnmarshalEventRegistryNodeRelease(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		_ = n.ProcessEvent(&c, msg.NodeId)
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *NodeStateEventServer) ProcessEvent(evt *evt.NodeEventConfig, nodeId string) error {
	currentState, err := n.s.GetLatestNodeState(context.Background(), &pb.GetLatestNodeStateRequest{NodeId: nodeId})
	currentStateName := ""
	if err != nil {
		log.Errorf("Error getting current state: %v", err)
		return err
	}
	if currentState == nil {
		initialState := "unknown"
		_, err = n.s.AddNodeState(context.Background(), &pb.AddNodeStateRequest{
			NodeId:       nodeId,
			CurrentState: initialState,
			Events:       []string{},
		})
		currentStateName = initialState
		if err != nil {
			log.Errorf("Error creating initial node state: %v", err)
			return err
		}
	}

	currentStateName = currentState.NodeState.CurrentState

	receivedEvents := []string{evt.Name}

	latestState := currentStateName

	nextState, err := n.s.stateMachine.GetNextState(latestState, receivedEvents)
	if err != nil {
		log.Errorf("Error getting next state: %v", err)
		return fmt.Errorf("failed to determine next state: %v", err)
	}

	if nextState != latestState {
		_, err = n.s.AddNodeState(context.Background(), &pb.AddNodeStateRequest{
			NodeId:       nodeId,
			CurrentState: nextState,
			Events:       receivedEvents,
		})
		if err != nil {
			log.Errorf("Error adding node state: %v", err)
			return err
		}

		log.Infof("Successfully processed events %v, new state: %s", receivedEvents, nextState)
		receivedEvents = nil
	} else {
		log.Infof("Event %s processed, state remains: %s", evt, latestState)
	}

	return nil
}
