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

 const (
	 NodeNotifyEventReady  = "ready"
	 NodeNotifyEventReboot = "reboot"
	 NodeNotifyEventOnline = "Node Online"
	 NodeNotifyEventAdded  = "Node added"
	 
	 DefaultSubstate = "on"
     ForceTransitionRoutingKeyTemplate = "event.cloud.local.{{ .Org}}.node.state.node.force"
	 NotifyEventRoutingKeyTemplate     = "event.cloud.local.{{ .Org}}.node.notify.notification.store"
 )
 

 type StateEventServer struct {
	 orgName        string
	 orgId          string
	 stateMachine   *stm.StateMachine
	 instances      map[string]*stm.StateMachineInstance
	 instancesMu    sync.RWMutex
	 s              *StateServer
	 configPath     string
	 epb.UnimplementedEventNotificationServiceServer
	 msgbus         mb.MsgBusServiceClient
	 baseRoutingKey msgbus.RoutingKeyBuilder
	 eventBuffer    map[string][]string
	 bufferMu       sync.RWMutex
	 processingMutex sync.Map
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
		 processingMutex: sync.Map{},
	 }
 
	 if configPath == "" {
		 log.Warn("State machine config path is empty, using default configuration")
	 }
 
	 server.stateMachine = stm.NewStateMachine(server.handleTransition)
 
	 return server
 }
 

 func (n *StateEventServer) handleTransition(event stm.Event) {
	 n.publishStateChangeEvent(event.NewState, event.NewSubstate, event.InstanceID)
 }
 

 func (n *StateEventServer) publishStateChangeEvent(state, substate, nodeID string) {
	 if n.msgbus == nil {
		 log.Warn("Message bus client is nil, skipping state change event publication")
		 return
	 }
 
	 route := n.baseRoutingKey.SetAction("transition").SetObject("node").MustBuild()
	 
	 eventsForNode := n.getEventsForNode(nodeID)
 
	 evt := &epb.NodeStateChangeEvent{
		 NodeId:   nodeID,
		 State:    state,
		 Substate: substate,
		 Events:   eventsForNode,
	 }
 
	 err := n.msgbus.PublishRequest(route, evt)
	 if err != nil {
		 log.Errorf("Failed to publish message %+v with key %+v. Error: %s", evt, route, err.Error())
	 }
	 
	 n.clearEventsForNode(nodeID)
 }
 
 func (n *StateEventServer) getEventsForNode(nodeID string) []string {
	 n.bufferMu.RLock()
	 defer n.bufferMu.RUnlock()
	 
	 events, ok := n.eventBuffer[nodeID]
	 if !ok {
		 return []string{}
	 }
	 return events
 }
 
 func (n *StateEventServer) clearEventsForNode(nodeID string) {
	 n.bufferMu.Lock()
	 defer n.bufferMu.Unlock()
	 delete(n.eventBuffer, nodeID)
 }
 

 func (n *StateEventServer) getOrCreateInstance(nodeID, initialState string) (*stm.StateMachineInstance, error) {
	 if nodeID == "" {
		 return nil, fmt.Errorf("node ID cannot be empty")
	 }
 
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
	 if e == nil {
		 return nil, status.Error(codes.InvalidArgument, "event cannot be nil")
	 }
 
	 log.Infof("Received event with routing key: %s", e.RoutingKey)
 
	 switch e.RoutingKey {
	 case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline]):
		 return n.handleNodeOnlineEvent(ctx, e)
 
	 case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOffline]):
		 return n.handleNodeOfflineEvent(ctx, e)
 
	 case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventAssign]):
		 return n.handleNodeAssignEvent(ctx, e)
 
	 case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventRelease]):
		 return n.handleNodeReleaseEvent(ctx, e)
 
	 case msgbus.PrepareRoute(n.orgName, ForceTransitionRoutingKeyTemplate):
		 return n.handleForceTransitionEvent(ctx, e)
		 
	 case msgbus.PrepareRoute(n.orgName, NotifyEventRoutingKeyTemplate):
		 return n.handleNodeNotifyEvent(ctx, e)
 
	 default:
		 log.Warnf("No handler for routing key %s", e.RoutingKey)
		 return &epb.EventResponse{}, nil
	 }
 }
 
 func (n *StateEventServer) handleNodeOnlineEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	 msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, e.RoutingKey)
	 if err != nil {
		 return nil, fmt.Errorf("failed to unmarshal node online event: %w", err)
	 }
	 eventName := evt.NodeEventToEventConfig[evt.NodeStateEventOnline].Name
	 if err := n.ProcessEvent(ctx, eventName, msg.NodeId, msg); err != nil {
		 return nil, fmt.Errorf("failed to process node online event: %w", err)
	 }
	 return &epb.EventResponse{}, nil
 }
 
 func (n *StateEventServer) handleNodeOfflineEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	 msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, e.RoutingKey)
	 if err != nil {
		 return nil, fmt.Errorf("failed to unmarshal node offline event: %w", err)
	 }
	 eventName := evt.NodeEventToEventConfig[evt.NodeStateEventOffline].Name
	 if err := n.ProcessEvent(ctx, eventName, msg.NodeId, msg); err != nil {
		 return nil, fmt.Errorf("failed to process node offline event: %w", err)
	 }
	 return &epb.EventResponse{}, nil
 }
 
 func (n *StateEventServer) handleNodeAssignEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	 msg, err := epb.UnmarshalEventRegistryNodeAssign(e.Msg, e.RoutingKey)
	 if err != nil {
		 return nil, fmt.Errorf("failed to unmarshal node assign event: %w", err)
	 }
	 eventName := evt.NodeEventToEventConfig[evt.NodeStateEventAssign].Name
	 if err := n.ProcessEvent(ctx, eventName, msg.NodeId, msg); err != nil {
		 return nil, fmt.Errorf("failed to process node assign event: %w", err)
	 }
	 return &epb.EventResponse{}, nil
 }
 
 func (n *StateEventServer) handleNodeReleaseEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	 msg, err := epb.UnmarshalEventRegistryNodeRelease(e.Msg, e.RoutingKey)
	 if err != nil {
		 return nil, fmt.Errorf("failed to unmarshal node release event: %w", err)
	 }
	 eventName := evt.NodeEventToEventConfig[evt.NodeStateEventRelease].Name
	 if err := n.ProcessEvent(ctx, eventName, msg.NodeId, msg); err != nil {
		 return nil, fmt.Errorf("failed to process node release event: %w", err)
	 }
	 return &epb.EventResponse{}, nil
 }
 
 func (n *StateEventServer) handleForceTransitionEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	 msg, err := n.UnmarshalTransitionEvent(e.Msg)
	 if err != nil {
		 return nil, fmt.Errorf("failed to unmarshal force transition event: %w", err)
	 }
	 
	 if err := n.ProcessEvent(ctx, msg.Event, msg.NodeId, msg); err != nil {
		 return nil, fmt.Errorf("failed to process force transition event: %w", err)
	 }
	 
	 return &epb.EventResponse{}, nil
 }
 
 func (n *StateEventServer) handleNodeNotifyEvent(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	 msg, err := n.unmarshalNotifyEvent(e.Msg)
	 if err != nil {
		 return nil, fmt.Errorf("failed to unmarshal notify event: %w", err)
	 }
	 
	 if err := n.handleNotifyEvent(ctx, e.RoutingKey, msg); err != nil {
		 return nil, fmt.Errorf("failed to process notify event: %w", err)
	 }
	 
	 return &epb.EventResponse{}, nil
 }
 
 func (n *StateEventServer) unmarshalNotifyEvent(msg *anypb.Any) (*epb.Notification, error) {
	 if msg == nil {
		 return nil, fmt.Errorf("notification message cannot be nil")
	 }
	 
	 p := &epb.Notification{}
	 err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	 if err != nil {
		 log.Errorf("Failed to unmarshal node notify message with: %+v. Error: %s", msg, err.Error())
		 return nil, fmt.Errorf("failed to unmarshal notification: %w", err)
	 }
	 return p, nil
 }
 
 func (n *StateEventServer) UnmarshalTransitionEvent(msg *anypb.Any) (*epb.EnforceNodeStateEvent, error) {
	 if msg == nil {
		 return nil, fmt.Errorf("transition event message cannot be nil")
	 }
	 
	 p := &epb.EnforceNodeStateEvent{}
	 err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	 if err != nil {
		 log.Errorf("Failed to unmarshal node transition message with: %+v. Error: %s", msg, err.Error())
		 return nil, fmt.Errorf("failed to unmarshal transition event: %w", err)
	 }
	 return p, nil
 }
 
 func (n *StateEventServer) handleNotifyEvent(ctx context.Context, _ string, msg *epb.Notification) error {
	 if msg == nil {
		 return fmt.Errorf("notification message cannot be nil")
	 }
	 
	 if msg.NodeId == "" {
		 return fmt.Errorf("node ID cannot be empty in notification")
	 }
	 
	 var details map[string]interface{}
	 if err := json.Unmarshal(msg.Details, &details); err != nil {
		 return fmt.Errorf("failed to unmarshal notification details: %w", err)
	 }
 
	 value, exists := details["value"]
	 if !exists {
		 return fmt.Errorf("value key not found in notification details")
	 }
 
	 valueStr, ok := value.(string)
	 if !ok {
		 return fmt.Errorf("value is not a string type in notification details")
	 }
	 
	 if valueStr == NodeNotifyEventOnline || valueStr == NodeNotifyEventAdded {
		 log.Infof("Skipping notification event processing for %s event", valueStr)
		 return nil
	 }
	 
	 log.Infof("Processing notification event %s for node %s", valueStr, msg.NodeId)
	 
	 if err := n.ProcessEvent(ctx, valueStr, msg.NodeId, msg); err != nil {
		 return fmt.Errorf("failed to process %s notification event: %w", valueStr, err)
	 }
	 
	 return nil
 }
 

 func (n *StateEventServer) ProcessEvent(ctx context.Context, eventName, nodeId string, msg interface{}) error {
	 if eventName == "" {
		 return fmt.Errorf("event name cannot be empty")
	 }
	 
	 if nodeId == "" {
		 return fmt.Errorf("node ID cannot be empty")
	 }
 
	 mutexValue, _ := n.processingMutex.LoadOrStore(nodeId, &sync.Mutex{})
	 mutex := mutexValue.(*sync.Mutex)
 
	 mutex.Lock()
	 defer mutex.Unlock()
 
	 log.Infof("Processing event %s for node %s", eventName, nodeId)
 
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
		 
		 log.Infof("No state found for node %s, creating initial state", nodeId)
		 return n.createInitialNodeState(ctx, nodeId, eventName, msg)
	 }
 
	 var currentState npb.NodeState
	 if latestState != nil && latestState.State != nil {
		 currentState = latestState.State.CurrentState
	 } else {
		 log.Infof("State information incomplete for node %s, creating initial state", nodeId)
		 return n.createInitialNodeState(ctx, nodeId, eventName, msg)
	 }
 
	 instance, err := n.getOrCreateInstance(nodeId, currentState.String())
	 if err != nil {
		 return fmt.Errorf("failed to create state machine instance for node %s: %w", nodeId, err)
	 }
 
	 prevState := instance.CurrentState
 
	 if err := instance.Transition(eventName); err != nil {
		 return fmt.Errorf("failed to transition state for node %s with event %s: %w", nodeId, eventName, err)
	 }
 
	 _, err = n.s.UpdateState(ctx, &pb.UpdateStateRequest{
		 NodeId:   nodeId,
		 SubState: []string{instance.CurrentSubstate}, 
		 Events:   []string{eventName},
	 })
	 if err != nil {
		 return fmt.Errorf("failed to update state for node %s: %w", nodeId, err)
	 }
 
	 if instance.CurrentState != prevState {
		 log.Infof("Node %s transitioning from state %s to %s", nodeId, prevState, instance.CurrentState)
		 
		 stateValue, ok := npb.NodeState_value[instance.CurrentState]
		 if !ok {
			 log.Warnf("Unknown state %s for node %s, defaulting to Unknown state", instance.CurrentState, nodeId)
			 stateValue = int32(npb.NodeState_Unknown)
		 }
		 
		 _, err = n.s.AddNodeState(ctx, &pb.AddStateRequest{
			 NodeId:       nodeId,
			 CurrentState: npb.NodeState(stateValue),
			 SubState:     []string{instance.CurrentSubstate},
			 Events:       []string{},
		 })
		 if err != nil {
			 return fmt.Errorf("failed to add new state for node %s: %w", nodeId, err)
		 }
	 }
	 
	 return nil
 }
 

 func (n *StateEventServer) createInitialNodeState(ctx context.Context, nodeId, eventName string, msg interface{}) error {
	 if nodeId == "" {
		 return fmt.Errorf("node ID cannot be empty")
	 }
	 
	 log.Infof("Creating initial state for node %s with event %s", nodeId, eventName)
	 
	 instance, err := n.getOrCreateInstance(nodeId, "Unknown")
	 if err != nil {
		 return fmt.Errorf("failed to create state machine instance: %w", err)
	 }
	 
	 if err := instance.Transition(eventName); err != nil {
		 log.Warnf("Initial transition failed for node %s with event %s: %v", nodeId, eventName, err)
		 
		 if eventName == "online" && instance.CurrentSubstate == "" {
			 instance.CurrentSubstate = DefaultSubstate
			 log.Infof("Setting default substate '%s' for node %s", DefaultSubstate, nodeId)
		 }
	 }
	 
	 initialSubstate := instance.CurrentSubstate
	 if initialSubstate == "" {
		 initialSubstate = DefaultSubstate 
		 log.Infof("Using default substate '%s' for node %s", DefaultSubstate, nodeId)
	 }
	 
	 var addStateRequest *pb.AddStateRequest
 
	 switch m := msg.(type) {
	 case *epb.NodeOnlineEvent:
		 addStateRequest = &pb.AddStateRequest{
			 NodeId:       nodeId,
			 CurrentState: npb.NodeState_Unknown,
			 SubState:     []string{initialSubstate},
			 Events:       []string{eventName},
			 NodeIp:       m.NodeIp,
			 NodePort:     int32(m.NodePort),
			 MeshIp:       m.MeshIp,
			 MeshPort:     int32(m.MeshPort),
			 MeshHostName: m.MeshHostName,
		 }
	 default:
		 addStateRequest = &pb.AddStateRequest{
			 NodeId:       nodeId,
			 CurrentState: npb.NodeState_Unknown,
			 SubState:     []string{initialSubstate},
			 Events:       []string{eventName},
		 }
	 }
	 
	 _, err = n.s.AddNodeState(ctx, addStateRequest)
	 if err != nil {
		 return fmt.Errorf("failed to create initial state entry for node %s: %w", nodeId, err)
	 }
	 
	 log.Infof("Initial state created for node %s with state Unknown, substate %s", nodeId, initialSubstate)
	 return nil
 }