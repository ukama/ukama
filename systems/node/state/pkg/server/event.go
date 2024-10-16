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
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	npb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	stm "github.com/ukama/ukama/systems/common/stateMachine"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"

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
}

func NewStateEventServer(orgName, orgId string, s *StateServer, configPath string) *StateEventServer {
	server := &StateEventServer{
		orgName:    orgName,
		orgId:      orgId,
		instances:  make(map[string]*stm.StateMachineInstance),
		s:          s,
		configPath: configPath,
	}

	server.stateMachine = stm.NewStateMachine(server.handleTransition)

	return server
}

func (n *StateEventServer) handleTransition(event stm.Event) {
	log.WithFields(log.Fields{
		"event":      event.Name,
		"oldState":   event.OldState,
		"newState":   event.NewState,
		"isSubstate": event.IsSubstate,
		"nodeId":     event.InstanceID,
	}).Info("Transition occurred")

	var subState string
	var currentState string
	var eventHistory []string

	if event.IsSubstate {
		subState = event.NewState
		currentState = event.OldState
	} else {
		currentState = event.NewState
	}
	eventHistory = append(eventHistory, event.Name)

	if !event.IsSubstate {
		ctx := context.Background()

		_, err := n.s.AddNodeState(ctx, &pb.AddStateRequest{
			NodeId:       event.InstanceID,
			CurrentState: currentState,
			SubState:     subState,
			Events:       eventHistory,
		})

		if err != nil {
			log.WithError(err).Error("Error adding node state")
		}
	}
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
	log.WithFields(log.Fields{
		"routingKey": e.RoutingKey,
		"message":    e.Msg,
	}).Info("Received event")

	var nodeID string
	var eventName string
	var body interface{}
	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventCreate]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventCreate]
		msg, err := epb.UnmarshalEventRegistryNodeCreate(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeCreate event: %w", err)
		}
		body = msg
		nodeID = msg.NodeId
		eventName = c.Name
	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventAssign]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventAssign]
		msg, err := epb.UnmarshalEventRegistryNodeAssign(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeAssign event: %w", err)
		}
		body = msg
		nodeID = msg.NodeId
		eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOffline]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventOffline]
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeOffline event: %w", err)
		}
		body = msg
		nodeID = msg.NodeId
		eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventOnline]
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeOnline event: %w", err)
		}
		body = msg
		nodeID = msg.NodeId
		eventName = c.Name

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventRelease]):
		c := evt.NodeEventToEventConfig[evt.NodeStateEventRelease]
		msg, err := epb.UnmarshalEventRegistryNodeRelease(e.Msg, c.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal NodeRelease event: %w", err)
		}
		body = msg
		nodeID = msg.NodeId
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
		nodeID = strings.Split(msg.Target, ".")[3]
		eventName = c.Name
	default:
		log.WithField("routingKey", e.RoutingKey).Error("No handler for routing key")
		return &epb.EventResponse{}, nil
	}

	if err := n.ProcessEvent(ctx, eventName, nodeID, body); err != nil {
		log.WithError(err).Error("Error processing event")
		return nil, err
	}

	return &epb.EventResponse{}, nil
}

func (n *StateEventServer) ProcessEvent(ctx context.Context, eventName, nodeID string, msg interface{}) error {
	log.WithFields(log.Fields{
		"event":  eventName,
		"nodeID": nodeID,
	}).Info("Processing event")

	latestState, err := n.s.GetLatestState(ctx, &pb.GetLatestStateRequest{NodeId: nodeID})
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
	if latestState != nil && latestState.State != nil {
		currentState = latestState.State.CurrentState
		log.WithFields(log.Fields{
			"nodeID": nodeID,
			"state":  currentState,
		}).Info("Node already exists. Attempting transition")
	} else {
		log.WithField("nodeID", nodeID).Info("Node does not exist. Creating new node")

		_, err := n.s.AddNodeState(ctx, &pb.AddStateRequest{
			NodeId:       nodeID,
			CurrentState: currentState,
			SubState:     "on",
			Events:       []string{eventName},
			NodeType:     msg.(epb.EventRegistryNodeCreate).Type,
			NodeIp:       msg.(epb.NodeOnlineEvent).NodeIp,
			NodePort:     int32(msg.(epb.NodeOnlineEvent).NodePort),
			MeshIp:       msg.(epb.NodeOnlineEvent).MeshIp,
			MeshPort:     int32(msg.(epb.NodeOnlineEvent).MeshPort),
			MeshHostName: msg.(epb.NodeOnlineEvent).MeshHostName,
		})
		if err != nil {
			return fmt.Errorf("failed to create initial state entry for node %s: %w", nodeID, err)
		}
	}

	instance, err := n.getOrCreateInstance(nodeID, currentState)
	if err != nil {
		return fmt.Errorf("failed to create state machine instance for node %s: %w", nodeID, err)
	}

	if err := instance.Transition(eventName); err != nil {
		return fmt.Errorf("failed to transition state for node %s: %w", nodeID, err)
	}

	log.WithFields(log.Fields{
		"nodeID": nodeID,
		"state":  instance.CurrentState,
	}).Info("Successfully processed event")

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
