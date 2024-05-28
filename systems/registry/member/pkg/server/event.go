/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type MemberEventServer struct {
	orgName string
	m       *MemberServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewPackageEventServer(orgName string, ms *MemberServer) *MemberEventServer {
	return &MemberEventServer{
		m:       ms,
		orgName: orgName,
	}
}

func (p *MemberEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(p.orgName, "event.cloud.local.{{ .Org }}.registry.invitation.invitation.create"):
		msg, err := unmarshalInvitationCreatedEvent(e.Msg)
		if err != nil {
			log.Errorf("Failed to unmarshal InvitationCreatedEvent message with error %s", err.Error())
			return &epb.EventResponse{}, err
		}
		if msg.Status == epb.StatusType_Accepted {
			p.m.AddMember(ctx, &pb.AddMemberRequest{
				UserUuid: msg.UserId,
				Role:     pb.RoleType(msg.Role),
			})
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalInvitationCreatedEvent(msg *anypb.Any) (*epb.EventInvitationCreated, error) {
	p := &epb.EventInvitationCreated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}
