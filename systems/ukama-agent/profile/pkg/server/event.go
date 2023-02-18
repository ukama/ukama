package server

import (
	"context"

	db "github.com/ukama/ukama/systems/ukama-agent/profile/pkg/db"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type ProfileEventServer struct {
	profileRepo db.ProfileRepo

	epb.UnimplementedEventNotificationServiceServer
}

func NewProfileEventServer(pRepo db.ProfileRepo) *ProfileEventServer {
	return &ProfileEventServer{
		profileRepo: pRepo,
	}
}

func (l *ProfileEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}
