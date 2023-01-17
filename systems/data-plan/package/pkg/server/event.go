package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/ukama/events/pb/gen"
	pb "github.com/ukama/ukama/systems/data-plan/package/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/package/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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

func (p *PackageEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
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
		break
	case "event.cloud.data-plan.package.delete":
		msg, err := unmarshalMsg(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
		break
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalMsg(msg *anypb.Any) (*pb.DeletePackageResponse, error) {
	p := &pb.DeletePackageResponse{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal REQ message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func handleEvent(key string, msg *pb.DeletePackageResponse) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	return nil
}
