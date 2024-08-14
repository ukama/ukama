package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type NodeStateEventServer struct {
    s       *StateServer
    orgName string
    epb.UnimplementedEventNotificationServiceServer
}

func NewControllerEventServer(orgName string, s *StateServer) *NodeStateEventServer {
    return &NodeStateEventServer{
        s:       s,
        orgName: orgName,
    }
}

func (n *NodeStateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
    log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
    switch e.RoutingKey {
    case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.registry.node.node.create"):
        msg, err := n.unmarshalRegistryNodeAddEvent(e.Msg)
        if err != nil {
            return nil, err
        }
        err = n.handleRegistryNodeAddEvent(e.RoutingKey, msg)
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
    case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline"):
        msg, err := n.unmarshalNodeOfflineEvent(e.Msg)
        if err != nil {
            return nil, err
        }
        err = n.handleNodeOfflineEvent(e.RoutingKey, msg)
        if err != nil {
            return nil, err
        }
    default:
        log.Errorf("No handler for routing key %s", e.RoutingKey)
    }

    return &epb.EventResponse{}, nil
}

func (n *NodeStateEventServer) unmarshalRegistryNodeAddEvent(msg *anypb.Any) (*epb.NodeCreatedEvent, error) {
	p := &epb.NodeCreatedEvent{}
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

func (n *NodeStateEventServer) handleRegistryNodeAddEvent(key string, msg *epb.NodeCreatedEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	state := &db.State{
		NodeId:       msg.NodeId,
		CurrentState: db.StateOnboarded,
		Connectivity: db.Unknown,
		Type:         msg.Type,
	}
	err := n.s.sRepo.Create(state, nil)
	if err != nil {
		log.Errorf("Error adding node %s to Nodestate repo. Error: %+v", msg.NodeId, err)
		return err
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
    state, err := n.s.sRepo.GetByNodeId(nId)
    if err != nil {
        log.Errorf("Error getting node %s from Nodestate repo. Error: %+v", msg.NodeId, err)
        return err
    }

    if state.Connectivity == db.Online {
        log.Infof("Node %s is already online. Ignoring event.", msg.NodeId)
        return nil
    }

    state.Connectivity = db.Online
    state.LastHeartbeat = time.Now()

    if state.CurrentState == db.StateFaulty {
        state.CurrentState = db.StateActive
    }

    err = n.s.sRepo.Update(state)
    if err != nil {
        log.Errorf("Error updating node %s in Nodestate repo. Error: %+v", msg.NodeId, err)
        return err
    }

    return nil
}

func (n *NodeStateEventServer) handleNodeOfflineEvent(key string, msg *epb.NodeOfflineEvent) error {
    log.Infof("Keys %s and Proto is: %+v", key, msg)
    nId, err := ukama.ValidateNodeId(msg.NodeId)
    if err != nil {
        log.Errorf("Error converting NodeId %s to ukama.NodeID. Error: %+v", msg.NodeId, err)
        return err
    }
    state, err := n.s.sRepo.GetByNodeId(nId)
    if err != nil {
        log.Errorf("Error getting node %s from Nodestate repo. Error: %+v", msg.NodeId, err)
        return err
    }

    if state.Connectivity == db.Offline {
        log.Infof("Node %s is already offline. Ignoring event.", msg.NodeId)
        return nil
    }

    state.Connectivity = db.Offline

    if state.CurrentState == db.StateActive {
        state.CurrentState = db.StateFaulty
    }

    err = n.s.sRepo.Update(state)
    if err != nil {
        log.Errorf("Error updating node %s in Nodestate repo. Error: %+v", msg.NodeId, err)
        return err
    }

    return nil
}