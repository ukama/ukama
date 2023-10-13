package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type EventServer struct {
	pb.UnimplementedEventNotificationServiceServer
}

func NewEventServer() *EventServer {
	return &EventServer{}
}

func (e *EventServer) EventNotification(ctx context.Context, evt *pb.Event) (*pb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", evt.RoutingKey, evt.Msg)
	return &pb.EventResponse{}, nil
}
