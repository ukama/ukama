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

	server.stateMachine = stm.NewStateMachine(server.handleTransition)

	return server
}

func (n *StateEventServer) handleTransition(event stm.Event) {
	n.publishStateChangeEvent(event.NewState, event.NewSubstate, event.InstanceID)
}

func (n *StateEventServer) publishStateChangeEvent(state, substate, nodeID string) {
	if n.msgbus == nil {
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
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
	}
}

func (n *StateEventServer) getEventsForNode(nodeID string) []string {
	n.bufferMu.RLock()
	defer n.bufferMu.RUnlock()
	return n.eventBuffer[nodeID]
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
	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline]):
		msg, err := epb.UnmarshalNodeOnlineEvent(e.Msg, e.RoutingKey)
		if err != nil {
			return nil, err
		}
		eventName := evt.NodeEventToEventConfig[evt.NodeStateEventOnline].Name
		err = n.ProcessEvent(ctx, eventName, msg.NodeId, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOffline]):
		msg, err := epb.UnmarshalNodeOfflineEvent(e.Msg, e.RoutingKey)
		if err != nil {
			return nil, err
		}
		eventName := evt.NodeEventToEventConfig[evt.NodeStateEventOffline].Name
		err = n.ProcessEvent(ctx, eventName, msg.NodeId, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventAssign]):
		msg, err := epb.UnmarshalEventRegistryNodeAssign(e.Msg, e.RoutingKey)
		if err != nil {
			return nil, err
		}
		eventName := evt.NodeEventToEventConfig[evt.NodeStateEventAssign].Name
		err = n.ProcessEvent(ctx, eventName, msg.NodeId, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(n.orgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventRelease]):
		msg, err := epb.UnmarshalEventRegistryNodeRelease(e.Msg, e.RoutingKey)
		if err != nil {
			return nil, err
		}
		eventName := evt.NodeEventToEventConfig[evt.NodeStateEventRelease].Name
		err = n.ProcessEvent(ctx, eventName, msg.NodeId, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.node.state.node.force"):
		msg, err := n.UnmarshalTransitionEvent(e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.handleEnforceTransitionEvent(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
		
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.node.notify.notification.store"):
		msg, err := n.unmarshalNotifyEvent(e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.handleNotifyEvent(ctx, e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler for routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *StateEventServer) unmarshalNotifyEvent(msg *anypb.Any) (*epb.Notification, error) {
	p := &epb.Notification{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal node notify message with: %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *StateEventServer) UnmarshalTransitionEvent(msg *anypb.Any) (*epb.EnforceNodeStateEvent, error) {
	p := &epb.EnforceNodeStateEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal node transition message with: %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *StateEventServer) handleEnforceTransitionEvent(ctx context.Context, _ string, msg *epb.EnforceNodeStateEvent) error {
	if err := n.ProcessEvent(ctx, msg.Event, msg.NodeId, msg); err != nil {
		log.WithError(err).Error("Error processing event")
		return err
	}
	return nil
}

func (n *StateEventServer) handleNotifyEvent(ctx context.Context, _ string, msg *epb.Notification) error {
	var details map[string]interface{}
	if err := json.Unmarshal(msg.Details, &details); err != nil {
		log.WithError(err).Error("Failed to unmarshal details")
		return err
	}

	value, exists := details["value"]
	if !exists {
		log.Warn("Value key not found in details")
		return fmt.Errorf("value key not found in details")
	}

	valueStr, ok := value.(string)
	if !ok {
		log.Error("Value is not a string type")
		return fmt.Errorf("value is not a string type")
	}
	
	if valueStr == "Node Online" || valueStr == "Node added"  {
		return nil
	}

	if err := n.ProcessEvent(ctx, valueStr, msg.NodeId, msg); err != nil {
		log.WithError(err).Error("Error processing event")
		return err
	}
	return nil
}

func (n *StateEventServer) ProcessEvent(ctx context.Context, eventName, nodeId string, msg interface{}) error {
	mutexValue, _ := n.processingMutex.LoadOrStore(nodeId, &sync.Mutex{})
	mutex := mutexValue.(*sync.Mutex)

	mutex.Lock()
	defer mutex.Unlock()

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

	var currentState npb.NodeState
	if latestState != nil && latestState.State != nil {
		currentState = latestState.State.CurrentState
	} else {
		if err := n.createInitialNodeState(ctx, nodeId, eventName, msg); err != nil {
			return err
		}
	}

	instance, err := n.getOrCreateInstance(nodeId, currentState.String())
	if err != nil {
		return fmt.Errorf("failed to create state machine instance for node %s: %w", nodeId, err)
	}

	prevState := instance.CurrentState

	if err := instance.Transition(eventName); err != nil {
		return fmt.Errorf("failed to transition state for node %s: %w", nodeId, err)
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
		stateValue := npb.NodeState_value[instance.CurrentState]
		
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
	instance, err := n.getOrCreateInstance(nodeId, "Unknown")
	if err != nil {
		return fmt.Errorf("failed to create state machine instance: %w", err)
	}
	
	if err := instance.Transition(eventName); err != nil {
		if eventName == "online" && instance.CurrentSubstate == "" {
			instance.CurrentSubstate = "on"
		}
	}
	
	initialSubstate := instance.CurrentSubstate
	if initialSubstate == "" {
		initialSubstate = "on" 
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
	
	return nil
}