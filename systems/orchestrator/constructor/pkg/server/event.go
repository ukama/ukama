package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
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

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}
