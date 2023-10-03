package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/software-manager/pb/gen"
	"github.com/ukama/ukama/systems/node/software-manager/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type SoftwareUpdateEventServer struct {
	orgName string
	sRepo   db.SoftwareManagerRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewSoftwareUpdateEventServer(orgName string, sRepo db.SoftwareManagerRepo) *SoftwareUpdateEventServer {
	return &SoftwareUpdateEventServer{
		orgName: orgName,
		sRepo:   sRepo,
	}
}

func (l *SoftwareUpdateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.hub.distributor.capp"):
		msg, err := unmarshalSoftwareUpdate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.sRepo.CreateSoftware(&db.Software{
			Id:          uuid.NewV4(),
			Name:        msg.Name,
			Tag:         msg.Version,
			Description: msg.Description,
			ReleaseDate: time.Now(),
		}, nil)
		if err != nil {
			return nil, err

		}

	default:
		log.Errorf("handler not registered for %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalSoftwareUpdate(msg *anypb.Any) (*gen.SoftwareUpdate, error) {
	p := &gen.SoftwareUpdate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal SoftwareUpdateEvent message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil

}
