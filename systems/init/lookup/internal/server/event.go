package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/init/lookup/internal/db"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
)

type LookupEventServer struct {
	systemRepo db.SystemRepo
	orgRepo    db.OrgRepo
	nodeRepo   db.NodeRepo
	pb.UnimplementedEventNotificationServiceServer
}

func NewLookupEventServer(nodeRepo db.NodeRepo, orgRepo db.OrgRepo, systemRepo db.SystemRepo) *LookupEventServer {
	return &LookupEventServer{
		nodeRepo:   nodeRepo,
		orgRepo:    orgRepo,
		systemRepo: systemRepo,
	}
}

func (l *LookupEventServer) EventNotification(ctx context.Context, e *pb.Event) (*pb.EventResponse, error) {
	log.Infof("Recieved a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	return &pb.EventResponse{}, nil
}
