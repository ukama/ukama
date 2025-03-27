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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type AsrEventServer struct {
	asrRepo  db.AsrRecordRepo
	gutiRepo db.GutiRepo
	s        *AsrRecordServer
	orgName  string
	epb.UnimplementedEventNotificationServiceServer
}

func NewAsrEventServer(asrRepo db.AsrRecordRepo, s *AsrRecordServer, gutiRepo db.GutiRepo, org string) *AsrEventServer {
	return &AsrEventServer{
		asrRepo:  asrRepo,
		gutiRepo: gutiRepo,
		orgName:  org,
		s:        s,
	}
}

func (l *AsrEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)
	switch e.RoutingKey {
	case msgbus.PrepareRoute(l.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.cdr.cdr.create"):
		msg, err := l.unmarshalCDRCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = l.handleEventCDRCreate(e.RoutingKey, msg)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func (l *AsrEventServer) unmarshalCDRCreate(msg *anypb.Any) (*epb.CDRReported, error) {
	p := &epb.CDRReported{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddSystemRequest message with : %+v. Error %s.", msg, err.Error())
		return nil, err
	}
	return p, nil
}

func (l *AsrEventServer) handleEventCDRCreate(key string, msg *epb.CDRReported) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)
	err := l.s.UpdateandSyncAsrProfile(msg.GetImsi())
	if err != nil {
		log.Errorf("Failed to update the active subscriber %+s.Error: %+v", msg.Imsi, err)
		return err
	}
	return nil
}
