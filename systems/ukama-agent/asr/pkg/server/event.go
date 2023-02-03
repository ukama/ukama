package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
)

type AsrEventServer struct {
	asrRepo  db.AsrRecordRepo
	gutiRepo db.GutiRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewAsrEventServer(asrRepo db.AsrRecordRepo, gutiRepo db.GutiRepo) *AsrEventServer {
	return &AsrEventServer{
		asrRepo:  asrRepo,
		gutiRepo: gutiRepo,
	}
}

func (l *AsrEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}
