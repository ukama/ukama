package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/msgbus"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type CDREventServer struct {
	s       *CDRServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
}

func NewCDREventServer(s *CDRServer, org string) *CDREventServer {
	return &CDREventServer{
		s:       s,
		orgName: org,
	}
}

func (n *CDREventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch msgbus.UpdateToAcceptFromAllOrg(e.RoutingKey) {
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.*.ukamaagent.asr.activesubscriber.create"):
		msg, err := n.unmarshalActiveSubscriberCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleEventActiveSubscriberCreate(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *CDREventServer) unmarshalActiveSubscriberCreate(msg *anypb.Any) (*epb.AsrUpdated, error) {
	p := &epb.AsrUpdated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddSystemRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (n *CDREventServer) handleEventActiveSubscriberCreate(key string, msg *epb.AsrUpdated) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	err := n.s.ResetPackageUsage(msg.Subscriber.Imsi)
	if err != nil {
		log.Errorf("Failed to update the active subscriber %+s.Error: %+v", msg.Subscriber.Imsi, err)
		return err
	}

	return nil
}
