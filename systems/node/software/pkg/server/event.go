package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/node/health/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg"

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
		anyMsg, err := anypb.New(msg)
		if err != nil {
			return nil, err
		}
		log.Infof("Received from health service: %v", msg)
 
		err = n.publishMessage(n.s.orgName + "." + "." + "." + msg.NodeId, anyMsg, msg.NodeId)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *SoftwareUpdateEventServer) unmarshalSoftwareUpdateEvent(msg *anypb.Any) (*pb.StoreRunningAppsInfoRequest, error) {
	p := &pb.StoreRunningAppsInfoRequest{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal node health  message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}



func (c *SoftwareUpdateEventServer) publishMessage(target string , anyMsg *anypb.Any ,nodeId string) error {
	route := "request.cloud.local" + "." + c.orgName + "." + pkg.SystemName + "." + pkg.ServiceName + "." + "nodefeeder" + "." + "publish"
	msg := &cpb.NodeFeederMessage{
		Target:     target,
		HTTPMethod: "POST",
		Path:       "/v1/update/" + nodeId,
		Msg:        anyMsg,
	}
	log.Infof("Published controller %s on route %s on target %s ",anyMsg,route,target)

	err := c.s.msgbus.PublishRequest(route, msg)
	return err
}