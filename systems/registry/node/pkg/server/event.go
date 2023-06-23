package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type NodeEventServer struct {
	s *NodeServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewNodeEventServer(s *NodeServer) *NodeEventServer {
	return &NodeEventServer{
		s: s,
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

// TODO: I am not sure I fully understand what's happening here in order for me to update
// so, commenting for compiling.
func (n *NodeEventServer) handleNodeOnlineEvent(key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	node, err := n.s.GetNode(context.Background(), &pb.GetNodeRequest{NodeId: msg.GetNodeId()})
	if err != nil {
		log.Error("error getting the node" + err.Error())
		return err
	}

	/* Add node if you can't find a node */
	if node == nil {
		req := &pb.AddNodeRequest{
			NodeId: node.Node.Id,
			OrgId:  n.s.org.String(),
		}

		_, err = n.s.AddNode(context.Background(), req)
		if err != nil {
			return err
		}
	} else {
		/* Update node status */
		_, err = n.s.UpdateNodeStatus(context.Background(), &pb.UpdateNodeStateRequest{
			NodeId:       node.Node.Id,
			Connectivity: db.Online.String(),
		})
		if err != nil {
			return err
		}
	}

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

	node, err := n.s.GetNode(context.Background(), &pb.GetNodeRequest{NodeId: msg.GetNodeId()})
	if err != nil {
		log.Error("error getting the node" + err.Error())
		return err
	}
	if node != nil {
		/* Update node status */
		_, err = n.s.UpdateNodeStatus(context.Background(), &pb.UpdateNodeStateRequest{
			NodeId:       node.Node.Id,
			Connectivity: db.Offline.String(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
