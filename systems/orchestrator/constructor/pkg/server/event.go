package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ConstructorEventServer struct {
	s *ConstructorServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewConstructorEventServer(s *ConstructorServer) *ConstructorEventServer {
	return &ConstructorEventServer{
		s: s,
	}
}

func (n *ConstructorEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case "event.cloud.mesh.node.online":
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

func (n *ConstructorEventServer) unmarshalNodeOnlineEvent(msg *anypb.Any) (*epb.NodeOnlineEvent, error) {
	p := &epb.NodeOnlineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOnline  message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

// so, commenting for compiling.
func (n *ConstructorEventServer) handleNodeOnlineEvent(key string, msg *epb.NodeOnlineEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	return nil
}
