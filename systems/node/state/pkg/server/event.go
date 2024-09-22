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
	"fmt"

	log "github.com/sirupsen/logrus"
	evt "github.com/ukama/ukama/systems/common/events"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	stm "github.com/ukama/ukama/systems/common/stateMachine"

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
 func NewNodeStateEventServer(orgName string,s *NodeStateServer, msgBus mb.MsgBusServiceClient) *NodeStateEventServer {
	 return &NodeStateEventServer{
		 s:               s,
		 orgName:         orgName,
		 msgbus:          msgBus,
		 stateRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	 }
 }
 

 func (es *NodeStateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v.", e.RoutingKey, e.Msg)
	switch e.RoutingKey {

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventOrgAdd]):
		c := evt.EventToEventConfig[evt.EventOrgAdd]
		msg, err := epb.UnmarshalEventOrgCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}
		// Handle Org Add event
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, msg.Id, jmsg)

	case msgbus.PrepareRoute(es.orgName, evt.EventRoutingKey[evt.EventUserAdd]):
		c := evt.EventToEventConfig[evt.EventUserAdd]
		msg, err := epb.UnmarshalEventUserCreate(e.Msg, c.Name)
		if err != nil {
			return nil, err
		}

		// Handle Org Add event
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to store raw message for %s to db. Error %+v", c.Name, err)
		}

		_ = es.ProcessEvent(&c, es.orgName, jmsg)
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (es *NodeStateEventServer) ProcessEvent(c *evt.EventConfig, nodeId string, rawMsg json.RawMessage) error {
    currentState, err := es.s.GetLatestNodeState(context.Background(), &pb.GetLatestNodeStateRequest{NodeId: nodeId})
    if err != nil {
        log.Errorf("Error getting current state: %v", err)
        return err
    }

    currentStateName := currentState.NodeState.CurrentState

    receivedEvents := []string{c.Name}

    latestState := currentStateName

    nextState, err := es.s.stateMachine.GetNextState(latestState, receivedEvents)
    if err != nil {
        log.Errorf("Error getting next state: %v", err)
        return fmt.Errorf("failed to determine next state: %v", err)
    }

    if nextState != latestState {
        _, err = es.s.AddNodeState(context.Background(), &pb.AddNodeStateRequest{
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
        log.Infof("Event %s processed, state remains: %s", c.Name, latestState)
    }

    return nil
}