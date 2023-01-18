package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
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
	case "event.cloud.simManager.sim.allocation":
		// msg, err := unmarshalLookupOrganizationCreate(e.Msg)
		// if err != nil {
		// 	return nil, err
		// }

		// err = handleEventCloudLookupOrgCreate(e.RoutingKey, msg)
		// if err != nil {
		// 	return nil, err
		// }
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

// func unmarshalLookupOrganizationCreate(msg *anypb.Any) (*pb.AddOrgRequest, error) {
// 	p := &pb.AddOrgRequest{}
// 	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
// 	if err != nil {
// 		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
// 		return nil, err
// 	}
// 	return p, nil
// }

// func handleEventCloudLookupOrgCreate(key string, msg *pb.AddOrgRequest) error {
// 	log.Infof("Keys %s and Proto is: %+v", key, msg)
// 	return nil
// }
