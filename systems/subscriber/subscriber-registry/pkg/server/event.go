package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/subscriber/subscriber-registry/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type SubcriberRegistryEventServer struct {
	subscriberRepo db.SubscriberRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewSubscriberEventServer(subscriberRepo db.SubscriberRepo) *SubcriberRegistryEventServer {
	return &SubcriberRegistryEventServer{
		subscriberRepo: subscriberRepo,
	}
}

func (l *SubcriberRegistryEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case "event.cloud.simManager.sim.allocation":
	default:
		log.Errorf("handler not registered for %s",e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalAllocateSim(msg *anypb.Any) (*pb.AddSubscriberRequest, error) {
	p := &pb.AddSubscriberRequest{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal UploadSimRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func handleEventAllocateSim(key string, msg *pb.AddSubscriberRequest) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	return nil
}