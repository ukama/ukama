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
	case "event.cloud.node.node.online":
		msg, err := n.unmarshalNodeOnlineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOnlineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case "event.cloud.node.node.offline":
		msg, err := n.unmarshalNodeOfflineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleNodeOfflineEvent(e.RoutingKey, msg)
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
		log.Errorf("Failed to Unmarshal NodeOnline  message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

//TODO: I am not sure I fully understand what's happening here in order for me to update
// so, commenting for compiling.
func (n *NodeEventServer) handleNodeOnlineEvent(key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	// nodeId, err := ukama.ValidateNodeId(msg.GetNodeId())
	// if err != nil {
	// logrus.Error("error getting the NodeId from request" + err.Error())
	// return err
	// }

	// node, err := n.nodeRepo.Get(nodeId)
	// if err != nil {
	// logrus.Error("error getting the node" + err.Error())
	// return err
	// }

	// if node == nil {
	// [> Add new node <]
	// node.Id = nodeId.StringLowercase()
	// // node.Allocation = false
	// node.Type = nodeId.GetNodeType()

	// err = AddNodeToOrg(n.nodeRepo, node)
	// if err != nil {
	// return err
	// }
	// } else {
	// state := db.Online
	// n.nodeRepo.Update(nodeId, &state, nil)
	// }

	return nil
}

func (n *NodeEventServer) unmarshalNodeOfflineEvent(msg *anypb.Any) (*epb.NodeOfflineEvent, error) {
	p := &epb.NodeOfflineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOffline message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *NodeEventServer) handleNodeOfflineEvent(key string, msg *epb.NodeOfflineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	nodeId, err := ukama.ValidateNodeId(msg.GetNodeId())
	if err != nil {
		logrus.Error("error getting the NodeId from request" + err.Error())
		return err
	}

	node, err := n.nodeRepo.Get(nodeId)
	if err != nil {
		logrus.Error("error getting the node" + err.Error())
		return err
	}

	if node != nil {
		node.State = db.Offline

		err := n.nodeRepo.Update(node, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
