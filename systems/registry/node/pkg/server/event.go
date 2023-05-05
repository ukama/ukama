package server

import (
	"context"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type NodeEventServer struct {
	nodeRepo db.NodeRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewNodeEventServer(nodeRepo db.NodeRepo) *NodeEventServer {
	return &NodeEventServer{
		nodeRepo: nodeRepo,
	}
}

func (n *NodeEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case "event.cloud.node.node.add":
		msg, err := n.unmarshalNodeOnlineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOnlineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *NodeEventServer) unmarshalNodeOnlineEvent(msg *anypb.Any) (*epb.NodeOnlineEvent, error) {
	p := &epb.NodeOnlineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *NodeEventServer) handleNodeOnlineEvent(key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	nodeID, err := ukama.ValidateNodeId(msg.GetNodeId())
	if err != nil {
		logrus.Error("error getting the NodeId from request" + err.Error())
		return err
	}

	node, err := n.nodeRepo.Get(nodeID)
	if err != nil {
		logrus.Error("error getting the node" + err.Error())
		return err
	}

	if node == nil {
		/* Add new node */
		node.NodeID = nodeID.StringLowercase()
		node.Allocation = false
		node.Type = nodeID.GetNodeType()

		err = AddNodeToOrg(n.nodeRepo, node)
		if err != nil {
			return err
		}
	} else {
		state := db.Online
		n.nodeRepo.Update(nodeID, &state, nil)
	}

	return nil
}

func (n *NodeEventServer) unmarshalNodeOfflineEvent(msg *anypb.Any) (*epb.NodeOfflineEvent, error) {
	p := &epb.NodeOfflineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *NodeEventServer) handleNodeOfflineEvent(key string, msg *epb.NodeOfflineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	nodeID, err := ukama.ValidateNodeId(msg.GetNodeId())
	if err != nil {
		logrus.Error("error getting the NodeId from request" + err.Error())
		return err
	}

	node, err := n.nodeRepo.Get(nodeID)
	if err != nil {
		logrus.Error("error getting the node" + err.Error())
		return err
	}

	if node != nil {
		state := db.Offline
		n.nodeRepo.Update(nodeID, &state, nil)
	}

	return nil
}
