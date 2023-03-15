package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
)

type PackageEventServer struct {
	packageRepo db.PackageRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewPackageEventServer(packageRepo db.PackageRepo) *PackageEventServer {
	return &PackageEventServer{
		packageRepo: packageRepo,
	}
}

func (p *PackageEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case "event.cloud.dataplan.rate.upload":
		break
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}
