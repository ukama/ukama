package server

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/node/software/pb/gen"
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

// func (l *SoftwareUpdateEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
// 	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
// 	//add another case for the other event
// 	switch e.RoutingKey {
// 	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.hub.distributor.capp"):
// 		msg, err := unmarshalSoftwareUpdate(e.Msg)
// 		if err != nil {
// 			return nil, err
// 		}
// 		err = l.s.sRepo.CreateSoftwareUpdate(&db.Software{
// 			Id:          uuid.NewV4(),
// 			Name:        msg.Name,
// 			Tag:         msg.Tag,
// 			ReleaseDate: time.Now(),
// 		}, nil)
// 		if err != nil {
// 			return nil, err

// 		}

// 	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.node.health.capps.store"):
// 		msg, err := unmarshalSoftwareUpdate(e.Msg)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println("Received from health service:", msg)

// 	default:
// 		log.Errorf("Handler not registered for %s", e.RoutingKey)
// 	}

// 	return &epb.EventResponse{}, nil
// }
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

	fmt.Println("Received from registry service:", msg)
	return nil
}