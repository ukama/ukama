package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	"github.com/ukama/ukama/systems/messaging/nns/pkg/client"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type NnsEventServer struct {
	Nns      *NnsServer
	Registry client.NodeRegistryClient
	Org      string
	epb.UnimplementedEventNotificationServiceServer
}

func NewNnsEventServer(c client.NodeRegistryClient, s *NnsServer, o string) *NnsEventServer {

	return &NnsEventServer{
		Registry: c,
		Nns:      s,
		Org:      o,
	}
}

func (l *NnsEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case "event.cloud.mesh.node.online":
		msg, err := l.unmarshalNodeOnlineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleNodeOnlineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case "event.cloud.mesh.node.offline":
		msg, err := l.unmarshalNodeOfflineEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleNodeOfflineEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	case "event.cloud.registry.node.assigned":
		msg, err := l.unmarshalNodeAssignedEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleNodeAssignedEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case "event.cloud.registry.node.release":
		msg, err := l.unmarshalNodeReleaseEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleNodeReleaseEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (l *NnsEventServer) unmarshalNodeOnlineEvent(msg *anypb.Any) (*epb.NodeOnlineEvent, error) {
	p := &epb.NodeOnlineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleNodeOnlineEvent(key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	log.Infof("Getting org and network for %s", msg.GetNodeId())

	nodeInfo := &client.NodeInfo{}
	var err error
	nodeInfo, err = l.Registry.GetNode(msg.GetNodeId())
	if err != nil {
		log.Errorf("Failed to get org and network. Error: %+v", err)
		log.Warningf("Node id %s won't have org and network info", msg.GetNodeId())
		nodeInfo = &client.NodeInfo{
			Id:      msg.GetNodeId(),
			Network: "",
			Site:    "",
		}
	}

	_, err = l.Nns.Set(context.Background(), &pb.SetNodeIPRequest{
		NodeId:   msg.GetNodeId(),
		NodeIp:   msg.GetMeshIp(),
		MeshIp:   msg.GetMeshIp(),
		NodePort: msg.GetNodePort(),
		MeshPort: msg.GetMeshPort(),
		Org:      l.Org,
		Network:  nodeInfo.Network,
		Site:     nodeInfo.Site,
	})

	if err != nil {
		log.Errorf("Failed to set node IP. Error: %+v", err)
		return err
	}
	log.Infof("Node %s IP set to %s", msg.GetNodeId(), msg.GetMeshIp())

	return nil
}

func (l *NnsEventServer) unmarshalNodeOfflineEvent(msg *anypb.Any) (*epb.NodeOfflineEvent, error) {
	p := &epb.NodeOfflineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleNodeOfflineEvent(key string, msg *epb.NodeOfflineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	return nil
}

func (l *NnsEventServer) unmarshalNodeAssignedEvent(msg *anypb.Any) (*epb.NodeAssignedEvent, error) {
	p := &epb.NodeAssignedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleNodeAssignedEvent(key string, msg *epb.NodeAssignedEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	orgNet, err := l.Nns.nodeOrgMapping.Get(context.Background(), msg.GetNodeId())
	if err != nil {
		log.Errorf("node %s doesn't exist. Error %v", msg.GetNodeId(), err)
		return err
	}

	err = l.Nns.nodeOrgMapping.Add(context.Background(), msg.GetNodeId(), l.Org, orgNet.Network, orgNet.Site, orgNet.NodeIp, orgNet.NodePort, orgNet.MeshPort)
	if err != nil {
		log.Errorf("failed to update labels for %s. Error %v", msg.GetNodeId(), err)
		return err
	}

	return nil
}

func (l *NnsEventServer) unmarshalNodeReleaseEvent(msg *anypb.Any) (*epb.NodeReleaseEvent, error) {
	p := &epb.NodeReleaseEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *NnsEventServer) handleNodeReleaseEvent(key string, msg *epb.NodeReleaseEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	orgNet, err := l.Nns.nodeOrgMapping.Get(context.Background(), msg.GetNodeId())
	if err != nil {
		log.Errorf("node %s doesn't exist. Error %v", msg.GetNodeId(), err)
		return err
	}

	err = l.Nns.nodeOrgMapping.Add(context.Background(), msg.GetNodeId(), l.Org, "", "", orgNet.NodeIp, orgNet.NodePort, orgNet.MeshPort)
	if err != nil {
		log.Errorf("failed to update labels for %s. Error %v", msg.GetNodeId(), err)
		return err
	}

	return nil
}
