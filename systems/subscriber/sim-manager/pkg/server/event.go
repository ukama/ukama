package server

import (
	"context"
	"time"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"
)

const (
	simAllocationRoutingKey = "event.cloud.simmanager.sim.allocate"
	handlerTimeoutFactor    = 3
)

type SimManagerEventServer struct {
	s *SimManagerServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewSimManagerEventServer(s *SimManagerServer) *SimManagerEventServer {
	return &SimManagerEventServer{
		s: s,
	}
}

func (es *SimManagerEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case simAllocationRoutingKey:
		msg, err := unmarshalSimManagerSimAllocate(e.Msg)
		if err != nil {
			return nil, err
		}

		if msg.Sim.Type == sims.SimTypeOperatorData.String() {
			err = handleEventCloudSimManagerOperatorSimAllocate(e.RoutingKey, msg, es.s)
			if err != nil {
				return nil, err
			}
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalSimManagerSimAllocate(msg *anypb.Any) (*pb.AllocateSimResponse, error) {
	p := &pb.AllocateSimResponse{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func handleEventCloudSimManagerOperatorSimAllocate(key string, msg *pb.AllocateSimResponse, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	_, err := s.activateSim(ctx, msg.Sim.Iccid)

	return err
}
