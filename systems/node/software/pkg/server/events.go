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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"

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
	case msgbus.PrepareRoute(n.orgName, "event.cloud.global.{{ .Org}}.hub.distributor.app.chunkready"):
		msg, err := epb.UnmarshalEventArtifactChunkReady(e.Msg, "Failed to unmarshal chunk ready event")
		if err != nil {
			return nil, err
		}
		_, err = n.s.CreateSoftwareUpdate(ctx, &pb.CreateSoftwareUpdateRequest{
			Name: msg.Name,
			Tag:  msg.Version,
			Space: "event",
		})
		if err != nil {
			return nil, err

		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (n *SoftwareUpdateEventServer) unmarshalSoftwareHubEvent(msg *anypb.Any) (*epb.EventArtifactChunkReady, error) {
	p := &epb.EventArtifactChunkReady{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal node health  message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}
