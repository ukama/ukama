package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type AsrEventServer struct {
	asrRepo  db.AsrRecordRepo
	gutiRepo db.GutiRepo
	s        *AsrRecordServer
	orgName  string
	epb.UnimplementedEventNotificationServiceServer
}

func NewAsrEventServer(asrRepo db.AsrRecordRepo, s *AsrRecordServer, gutiRepo db.GutiRepo, org string) *AsrEventServer {
	return &AsrEventServer{
		asrRepo:  asrRepo,
		gutiRepo: gutiRepo,
		orgName:  org,
		s:        s,
	}
}

func (l *AsrEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.*.ukamaagent.cdr.cdr.create"):
		msg, err := l.unmarshalCDRCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleEventCDRCreate(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (l *AsrEventServer) unmarshalCDRCreate(msg *anypb.Any) (*epb.CDRReported, error) {
	p := &epb.CDRReported{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddSystemRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *AsrEventServer) handleEventCDRCreate(key string, msg *epb.CDRReported) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	// err := l.s.UpdateUsage()
	// if err != nil {
	// 	log.Errorf("Failed to update the active subscriber %+s.Error: %+v", msg.Subscriber.Imsi, err)
	// 	return err
	// }

	return nil
}
