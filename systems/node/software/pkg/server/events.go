package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/software/pkg/db"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type SoftwareUpdateEventServer struct {
	s       *SoftwareServer
	orgName string
	epb.UnimplementedEventNotificationServiceServer
}

func NewSoftwareEventServer(orgName string, s *SoftwareServer) *SoftwareUpdateEventServer {
	return &SoftwareUpdateEventServer{
		s:       s,
		orgName: orgName,
	}
}
func (n *SoftwareUpdateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.hub.distributor.capp"):
		msg, err := n.unmarshalSoftwareHubEvent(e.Msg)
		if err != nil {
			return nil, err
		}
		err = n.s.sRepo.CreateSoftwareUpdate(&db.Software{
			Id:          uuid.NewV4(),
			Name:        msg.Name,
			Tag:         msg.Version,
			ReleaseDate: time.Now(),
		}, nil)
		if err != nil {
			return nil, err

		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *SoftwareUpdateEventServer) unmarshalSoftwareHubEvent(msg *anypb.Any) (*epb.CappCreatedEvent, error) {
	p := &epb.CappCreatedEvent{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal node health  message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}
