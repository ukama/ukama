package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type SimPoolEventServer struct {
	orgName     string
	simPoolRepo db.SimRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewSimPoolEventServer(orgName string, simPoolRepo db.SimRepo) *SimPoolEventServer {
	return &SimPoolEventServer{
		orgName:     orgName,
		simPoolRepo: simPoolRepo,
	}
}
func (l *SimPoolEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocation"):
		msg, err := unmarshalAllocateSim(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventAllocateSim(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("handler not registered for %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalAllocateSim(msg *anypb.Any) (*pb.UploadRequest, error) {
	p := &pb.UploadRequest{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal UploadSimRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func handleEventAllocateSim(key string, msg *pb.UploadRequest) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	return nil
}
