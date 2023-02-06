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
		// msg, err := unmarshalMsg(e.Msg)
		// if err != nil {
		// 	return nil, err
		// }

		// err = handleEvent(e.RoutingKey, msg)
		// if err != nil {
		// 	return nil, err
		// }
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &pb.EventResponse{}, nil
}

// func unmarshalMsg(msg *anypb.Any) (*pb., error) {
// 	p := &pb.{}
// 	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
// 	if err != nil {
// 		log.Errorf("Failed to Unmarshal REQ message with : %+v. Error %s.", msg, err.Error())
// 		return nil, err
// 	}
// 	return p, nil
// }

// func handleEvent(key string, msg *pb.) error {
// 	log.Infof("Keys %s and Proto is: %+v", key, msg)
// 	return nil
// }
