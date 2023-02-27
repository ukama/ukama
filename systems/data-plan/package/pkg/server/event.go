package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
)

type PackageEventServer struct {
	packageRepo db.PackageRepo
	pb.UnimplementedPackagesServiceServer
}

func NewPackageEventServer(packageRepo db.PackageRepo) *PackageEventServer {
	return &PackageEventServer{
		packageRepo: packageRepo,
	}
}

func (p *PackageEventServer) EventNotification(ctx context.Context, e *pb.Event) (*pb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case "event.cloud.data-plan.base-rate.upload":
		break
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &pb.EventResponse{}, nil
}
