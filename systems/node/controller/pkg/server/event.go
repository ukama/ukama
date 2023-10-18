package server

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ControllerEventServer struct {
	s       *ControllerServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
}

func NewControllerEventServer(orgName string, s *ControllerServer) *ControllerEventServer {
	return &ControllerEventServer{
		s:       s,
		orgName: orgName,
	}
}

func (n *ControllerEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
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

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *ControllerEventServer) unmarshalRegistryNodeAddEvent(msg *anypb.Any) (*epb.NodeCreatedEvent, error) {
	p := &epb.NodeCreatedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOnline  message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *ControllerEventServer) handleRegistryNodeAddEvent(key string, msg *epb.NodeCreatedEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	err := n.s.sRepo.Add(strings.ToLower(msg.NodeId))
	if err != nil {
		log.Errorf("Error adding node %s to controller/nodeLog repo.Error: %+v", msg.NodeId, err)
		return err
	}
	return nil
}