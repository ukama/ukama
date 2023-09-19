package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/msgbus"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
)

type PackageEventServer struct {
	orgName     string
	packageRepo db.PackageRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewPackageEventServer(orgName string, packageRepo db.PackageRepo) *PackageEventServer {
	return &PackageEventServer{
		orgName:     orgName,
		packageRepo: packageRepo,
	}
}

func (p *PackageEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(p.orgName, "event.cloud.local.{{ .Org }}.dataplan.baserate.rate.upload"):
		break
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}
