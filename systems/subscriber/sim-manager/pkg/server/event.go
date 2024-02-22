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
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"

	log "github.com/sirupsen/logrus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

const (
	handlerTimeoutFactor = 3
)

type SimManagerEventServer struct {
	orgName string
	s       *SimManagerServer
	epb.UnimplementedEventNotificationServiceServer
}

func NewSimManagerEventServer(orgName string, s *SimManagerServer) *SimManagerEventServer {
	return &SimManagerEventServer{
		orgName: orgName,
		s:       s,
	}
}

func (es *SimManagerEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate"):
		msg, err := unmarshalSimManagerSimAllocate(e.Msg)
		if err != nil {
			return nil, err
		}

		simType := ukama.ParseSimType(msg.Sim.Type)

		if simType == ukama.SimTypeOperatorData {
			err = handleEventCloudSimManagerOperatorSimAllocate(e.RoutingKey, msg, es.s)
			if err != nil {
				return nil, err
			}
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalSimManagerSimAllocate(msg *anypb.Any) (*pb.AllocateSimResponse, error) {
	p := &pb.AllocateSimResponse{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func handleEventCloudSimManagerOperatorSimAllocate(key string, msg *pb.AllocateSimResponse, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	_, err := s.activateSim(ctx, msg.Sim.Iccid)

	return err
}
