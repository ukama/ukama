package server

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/node/software/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	case msgbus.PrepareRoute(n.orgName, "event.cloud.local.{{ .Org}}.node.health.capps.store"):
		msg, err := n.unmarshalSoftwareUpdateEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = n.handleRegistryNodeAddEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *SoftwareUpdateEventServer) unmarshalSoftwareUpdateEvent(msg *anypb.Any) (*epb.NodeCreatedEvent, error) {
	p := &epb.NodeCreatedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal NodeOnline  message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func unmarshalSoftwareUpdate(msg *anypb.Any) (*gen.SoftwareUpdate, error) {
	p := &gen.SoftwareUpdate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal SoftwareUpdateEvent message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil

}
func (n *SoftwareUpdateEventServer) handleRegistryNodeAddEvent(key string, msg *epb.NodeCreatedEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	anyMsg, err := anypb.New(msg)
	if err != nil {
		return  err
	}
	
	   err = n.publishMessage(n.s.orgName + "." + "." + "." + msg.NodeId, anyMsg,msg.NodeId)
	   if err != nil {
		   log.Errorf("Failed to publish message. Errors %s", err.Error())
		   return  status.Errorf(codes.Internal, "Failed to publish message: %s", err.Error())
   
	   }
	fmt.Println("Received from health service:", msg)
	return nil
}

func (e *SoftwareUpdateEventServer) publishMessage(target string , anyMsg *anypb.Any ,nodeId string) error {
	route := "request.cloud.local" + "." + e.s.orgName + "." + pkg.SystemName + "." + pkg.ServiceName + "." + "nodefeeder" + "." + "publish"
	msg := &cpb.NodeFeederMessage{
		Target:     target,
		HTTPMethod: "POST",
		Path:       "/v1/update/" + nodeId,
		Msg:        anyMsg,
	}
	log.Infof("Published controller %s on route %s on target %s ",anyMsg,route,target)

	err := e.s.msgbus.PublishRequest(route, msg)
	return err
}