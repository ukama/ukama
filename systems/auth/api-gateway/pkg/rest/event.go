package rest

import (
	"context"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type AuthEventServer struct {
	epb.UnimplementedEventNotificationServiceServer
}

func NewAuthEventServer() *AuthEventServer {
	return &AuthEventServer{}
}

func (l *AuthEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case "event.cloud.registry.member.add":
		log.Info(e.Msg)
		msg, err := unmarshalAddMemeberReq(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudAddMember(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalAddMemeberReq(msg *anypb.Any) (*pb.AddMemberRequest, error) {
	p := &pb.AddMemberRequest{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddMemberRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func handleEventCloudAddMember(key string, msg *pb.AddMemberRequest) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	// Call Kratos to update user role
	return nil
}
