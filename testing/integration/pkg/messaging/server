package messaging

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type EventServer struct {
	store map[string]interface{}
	epb.UnimplementedEventNotificationServiceServer
}

func NewEventServer() *EventServer {
	emap := make(map[string]interface{})
	return &EventServer{
		store: emap,
	}
}

func (s *EventServer) GetEvent(key string) (interface{}, bool) {
	if m, ok := s.store[key]; ok {
		return m, ok
	}
	return nil, false
}

func (s *EventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	s.store[e.RoutingKey] = e.Msg
	switch e.RoutingKey {
	case "event.cloud.lookup.organization.create":

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}
