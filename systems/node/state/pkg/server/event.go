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
	s       *StateServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
	msgbus          mb.MsgBusServiceClient
	stateRoutingKey msgbus.RoutingKeyBuilder
}

func NewControllerEventServer(orgName string, s *StateServer, msgBus mb.MsgBusServiceClient) *NodeStateEventServer {
	return &NodeStateEventServer{
		s:               s,
		orgName:         orgName,
		stateRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:          msgBus,
	}
}

func (n *NodeStateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.create"):
		msg, err := n.unmarshalNodeCreateEvent(e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.handleNodeCreateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.online"):
		msg, err := n.unmarshalNodeOnlineEvent(e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.handleNodeOnlineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.assign"):
		msg, err := n.unmarshalOnboardingEventEvent(e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.handleOnboardingEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline"):
		msg, err := n.unmarshalNodeOfflineEvent(e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.handleNodeOfflineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.node.notify.notification.store"):
		msg, err := n.unmarshalNodeHealthSeverityHighEvent(e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.handleNodeHealthSeverityHighEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler for routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *NodeStateEventServer) handleNodeHealthSeverityHighEvent(key string, msg *epb.Notification) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	nId, err := ukama.ValidateNodeId(msg.NodeId)
	if err != nil {
		log.Errorf("Error converting NodeId %s to ukama.NodeID. Error: %+v", msg.NodeId, err)
		return err
	}

	currentState, err := n.s.sRepo.GetByNodeId(nId)
	if err != nil {
		log.Errorf("Error getting latest state for node %s from Nodestate repo. Error: %+v", nId, err)
		return err
	}

	now := time.Now()
	newState := &db.State{
		Id:              uuid.NewV4(),
		NodeId:          nId.String(),
		State:           ukama.StateFaulty,
		Type:            currentState.Type,
		LastStateChange: now,
		LastHeartbeat:   currentState.LastHeartbeat,
		Version:         currentState.Version,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err = n.s.sRepo.Create(newState, nil)
	if err != nil {
		log.Errorf("Error creating new state for node %s in Nodestate repo. Error: %+v", nId, err)
		return err
	}

	if n.s.msgbus != nil {
		route := n.s.stateRoutingKey.SetAction("state").SetObject("node").MustBuild()
		evt := &epb.EventNodeStateUpdate{
			NodeId:          newState.NodeId,
			CurrentState:    newState.State.String(),
			LastStateChange: newState.LastStateChange.String(),
		}

		err = n.s.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish node state update event: %+v with key %+v. Error: %s", evt, route, err.Error())
		} else {
			log.Infof("Published node state update event for node %s", nId.String())
		}
	}

	log.Infof("Updated node %s state to Faulty", nId)
	return nil
}

func (n *NodeStateEventServer) unmarshalNodeCreateEvent(msg *anypb.Any) (*epb.EventRegistryNodeCreate, error) {
	p := &epb.EventRegistryNodeCreate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal Node create message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}
func (n *NodeStateEventServer) unmarshalNodeHealthSeverityHighEvent(msg *anypb.Any) (*epb.Notification, error) {
	p := &epb.Notification{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal Node severity high message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *NodeStateEventServer) handleOnboardingEvent(key string, msg *epb.EventRegistryNodeAssign) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	now := time.Now()
	state := &db.State{
		Id:              uuid.NewV4(),
		NodeId:          msg.NodeId,
		State:           ukama.StateConfigure,
		Type:            msg.Type,
		LastStateChange: now,
		LastHeartbeat:   now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err := n.s.sRepo.Create(state, nil)
	if err != nil {
		log.Errorf("Error adding node %s to Nodestate repo. Error: %+v", msg.NodeId, err)
		return err
	}
	if n.s.msgbus != nil {
		route := n.s.stateRoutingKey.SetAction("state").SetObject("node").MustBuild()
		evt := &epb.EventNodeStateUpdate{
			NodeId:          msg.NodeId,
			CurrentState:    state.State.String(),
			LastStateChange: state.LastStateChange.String(),
		}

		err = n.s.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish node state update event: %+v with key %+v. Error: %s", evt, route, err.Error())
		} else {
			log.Infof("Published node state update event for node %s", msg.NodeId)
		}
	}
	return nil
}

func (n *NodeStateEventServer) handleNodeOnlineEvent(key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	nId, err := ukama.ValidateNodeId(msg.NodeId)
	if err != nil {
		log.Errorf("Error converting NodeId %s to ukama.NodeID. Error: %+v", msg.NodeId, err)
		return err
	}

	currentState, err := n.s.sRepo.GetByNodeId(nId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Node doesn't exist, create a new one
			now := time.Now()
			newState := &db.State{
				Id:              uuid.NewV4(),
				NodeId:          nId.String(),
				State:           ukama.StateUnknown,
				LastStateChange: now,
				LastHeartbeat:   now,
				CreatedAt:       now,
				UpdatedAt:       now,
			}
			err = n.s.sRepo.Create(newState, nil)
			if err != nil {
				log.Errorf("Error creating new node %s in Nodestate repo. Error: %+v", nId, err)
				return err
			}
			log.Infof("Created new node %s with Unknown state", nId)
			return nil
		}
		log.Errorf("Error getting latest state for node %s from Nodestate repo. Error: %+v", nId, err)
		return err
	}

	// If node was already online, ignore the event
	if currentState.LastHeartbeat.Add(time.Minute * 5).After(time.Now()) {
		log.Infof("Node %s was already online, ignoring online event", nId)
		return nil
	}

	now := time.Now()
	newState := &db.State{
		Id:              uuid.NewV4(),
		NodeId:          nId.String(),
		State:           currentState.State,
		Type:            currentState.Type,
		LastStateChange: currentState.LastStateChange,
		LastHeartbeat:   now,
		Version:         currentState.Version,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err = n.s.sRepo.Create(newState, nil)
	if err != nil {
		log.Errorf("Error creating new state for node %s in Nodestate repo. Error: %+v", nId, err)
		return err
	}
	if n.s.msgbus != nil {
		route := n.s.stateRoutingKey.SetAction("state").SetObject("node").MustBuild()
		evt := &epb.EventNodeStateUpdate{
			NodeId:          newState.NodeId,
			CurrentState:    newState.State.String(),
			LastStateChange: newState.LastStateChange.String(),
		}

		err = n.s.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish node state update event: %+v with key %+v. Error: %s", evt, route, err.Error())
		} else {
			log.Infof("Published node state update event for node %s", nId.String())
		}
	}
	log.Infof("Updated node %s heartbeat", nId)
	return nil
}

func (n *NodeStateEventServer) handleNodeCreateEvent(key string, msg *epb.EventRegistryNodeCreate) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	nId, err := ukama.ValidateNodeId(msg.NodeId)
	if err != nil {
		log.Errorf("Error converting NodeId %s to ukama.NodeID. Error: %+v", msg.NodeId, err)
		return err
	}

	currentState, err := n.s.sRepo.GetByNodeId(nId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Node doesn't exist, create a new one
			now := time.Now()
			newState := &db.State{
				Id:              uuid.NewV4(),
				NodeId:          nId.String(),
				State:           ukama.StateUnknown,
				LastStateChange: now,
				LastHeartbeat:   now,
				CreatedAt:       now,
				UpdatedAt:       now,
			}
			err = n.s.sRepo.Create(newState, nil)
			if err != nil {
				log.Errorf("Error creating new node %s in Nodestate repo. Error: %+v", nId, err)
				return err
			}
			log.Infof("Created new node %s with Unknown state", nId)
			return nil
		}
		log.Errorf("Error getting latest state for node %s from Nodestate repo. Error: %+v", nId, err)
		return err
	}

	// If node was already online, ignore the event
	if currentState.LastHeartbeat.Add(time.Minute * 5).After(time.Now()) {
		log.Infof("Node %s was already online, ignoring online event", nId)
		return nil
	}

	now := time.Now()
	newState := &db.State{
		Id:              uuid.NewV4(),
		NodeId:          nId.String(),
		State:           currentState.State,
		Type:            currentState.Type,
		LastStateChange: currentState.LastStateChange,
		LastHeartbeat:   now,
		Version:         currentState.Version,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err = n.s.sRepo.Create(newState, nil)
	if err != nil {
		log.Errorf("Error creating new state for node %s in Nodestate repo. Error: %+v", nId, err)
		return err
	}
	if n.s.msgbus != nil {
		route := n.s.stateRoutingKey.SetAction("state").SetObject("node").MustBuild()
		evt := &epb.EventNodeStateUpdate{
			NodeId:          newState.NodeId,
			CurrentState:    newState.State.String(),
			LastStateChange: newState.LastStateChange.String(),
		}

		err = n.s.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish node state update event: %+v with key %+v. Error: %s", evt, route, err.Error())
		} else {
			log.Infof("Published node state update event for node %s", nId.String())
		}
	}

	log.Infof("Updated node %s heartbeat", nId)
	return nil
}

func (n *NodeStateEventServer) handleNodeOfflineEvent(key string, msg *epb.NodeOfflineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	// We're not changing the state for offline events, just logging it
	log.Infof("Node %s is offline", msg.NodeId)
	return nil
}

func (n *NodeStateEventServer) unmarshalOnboardingEventEvent(msg *anypb.Any) (*epb.EventRegistryNodeAssign, error) {
	p := &epb.EventRegistryNodeAssign{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeCreated message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *NodeStateEventServer) unmarshalNodeOnlineEvent(msg *anypb.Any) (*epb.NodeOnlineEvent, error) {
	p := &epb.NodeOnlineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOnline message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *NodeStateEventServer) unmarshalNodeOfflineEvent(msg *anypb.Any) (*epb.NodeOfflineEvent, error) {
	p := &epb.NodeOfflineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOffline message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *NodeStateEventServer) unmarshalNodeConfigCreateEvent(msg *anypb.Any) (*eCfgPb.NodeConfigUpdateEvent, error) {
	p := &eCfgPb.NodeConfigUpdateEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeConfigCreate message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}
