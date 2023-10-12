package server

import (
	"context"

	"github.com/jackc/pgtype"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/init/lookup/internal/db"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type LookupEventServer struct {
	orgName    string
	systemRepo db.SystemRepo
	orgRepo    db.OrgRepo
	nodeRepo   db.NodeRepo
	epb.UnimplementedEventNotificationServiceServer
}

func NewLookupEventServer(orgName string, nodeRepo db.NodeRepo, orgRepo db.OrgRepo, systemRepo db.SystemRepo) *LookupEventServer {
	return &LookupEventServer{
		orgName:    orgName,
		nodeRepo:   nodeRepo,
		orgRepo:    orgRepo,
		systemRepo: systemRepo,
	}
}

func (l *LookupEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.init.lookup.organization.create"):
		msg, err := l.unmarshalLookupOrganizationCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleEventCloudLookupOrgCreate(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.messaging.mesh.ip.update"):
		msg, err := l.unmarshalOrgIpUpdateEvent(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleEventOrgIPUpdateEvent(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (l *LookupEventServer) unmarshalLookupOrganizationCreate(msg *anypb.Any) (*pb.AddOrgRequest, error) {
	p := &pb.AddOrgRequest{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *LookupEventServer) handleEventCloudLookupOrgCreate(key string, msg *pb.AddOrgRequest) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	return nil
}

func (l *LookupEventServer) unmarshalOrgIpUpdateEvent(msg *anypb.Any) (*epb.OrgIPUpdateEvent, error) {
	p := &epb.OrgIPUpdateEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *LookupEventServer) handleEventOrgIPUpdateEvent(key string, msg *epb.OrgIPUpdateEvent) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	/* Read the Org first */
	org, err := l.orgRepo.GetByName(msg.OrgName)
	if err != nil {
		log.Errorf("Org %s not found", msg.OrgName)
		return status.Error(codes.NotFound, "Org not found")
	}

	var orgIp pgtype.Inet

	err = orgIp.Set(msg.Ip)
	if err != nil {
		log.Errorf("Invalid ip %s for Org %s. Error %s", msg.Ip, msg.OrgName, err.Error())
		return err
	}

	org.Ip = orgIp

	/* Update Ip */
	err = l.orgRepo.Update(org)
	if err != nil {
		log.Errorf("Error updating org %s", msg.OrgName)
		return err
	}
	return nil
}
