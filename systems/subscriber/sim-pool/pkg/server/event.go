package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type SimPoolEventServer struct {
	simPoolRepo db.SimRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewSimPoolEventServer(simPoolRepo db.SimRepo) *SimPoolEventServer {
	return &SimPoolEventServer{
		simPoolRepo: simPoolRepo,
	}
}
func (l *SimPoolEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case "event.cloud.simManager.sim.allocate":
		msg, err := unmarshalAllocateSim(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.simPoolRepo.UpdateStatus(msg.Iccid, true, false)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("handler not registered for %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalAllocateSim(msg *anypb.Any) (*pb.Sim, error) {
	p := &pb.Sim{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal UploadSimRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}
