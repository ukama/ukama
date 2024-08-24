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
	"fmt"
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

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.operator.cdr.cdr.create"):
		msg, err := unmarshalOperatorCdrCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudOperatorCdrCreate(e.RoutingKey, msg, es.s)
		if err != nil {
			return nil, err
		}

	case msgbus.PrepareRoute(es.orgName, "event.cloud.local.{{ .Org}}.ukamaagent.cdr.cdr.create"):
		msg, err := unmarshalUkamaAgentCdrCreate(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudUkamaAgentCdrCreate(e.RoutingKey, msg, es.s)
		if err != nil {
			return nil, err
		}
	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func handleEventCloudSimManagerOperatorSimAllocate(key string, msg *pb.AllocateSimResponse, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*handlerTimeoutFactor)
	defer cancel()

	_, err := s.activateSim(ctx, msg.Sim.Iccid)

	return err
}

func handleEventCloudOperatorCdrCreate(key string, cdr *epb.EventOperatorCdrReport, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	sim, err := s.simRepo.GetByIccid(cdr.Iccid)
	if err != nil {
		return fmt.Errorf("no corresponding sim found for given iccid %q: %v",
			cdr.Iccid, err)
	}

	usageMsg := &epb.EventSimUsage{
		SimId:        sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		NetworkId:    sim.NetworkId.String(),
		Type:         cdr.Type,
		BytesUsed:    cdr.Duration,
		StartTime:    cdr.ConnectTime,
		EndTime:      cdr.CloseTime,
		Id:           cdr.Id,
		// OrgId:        s.OrgId.String(),
		// SessionId: msg.InventoryId,
	}

	route := s.baseRoutingKey.SetAction("usage").SetObject("sim").MustBuild()

	err = s.msgbus.PublishRequest(route, usageMsg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			usageMsg, route, err.Error())
	}

	return err
}

func handleEventCloudUkamaAgentCdrCreate(key string, cdr *epb.CDRReported, s *SimManagerServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, cdr)

	// TODO: implement simRepo.List
	// sim, err := s.simRepo.GetByImsi(cdr.Imsi)
	sim, err := s.simRepo.GetByIccid(cdr.Imsi)
	if err != nil {
		return fmt.Errorf("no corresponding sim found for given iccid %q: %v",
			cdr.Imsi, err)
	}

	usageMsg := &epb.EventSimUsage{
		SimId:        sim.Id.String(),
		SubscriberId: sim.SubscriberId.String(),
		NetworkId:    sim.NetworkId.String(),
		Type:         ukama.CdrTypeData.String(),
		BytesUsed:    cdr.TotalBytes,
		StartTime:    cdr.StartTime,
		EndTime:      cdr.EndTime,
		// Id:           cdr.Id,
		// OrgId:        s.OrgId.String(),
		// SessionId:    cdr.Session,
	}

	route := s.baseRoutingKey.SetAction("usage").SetObject("sim").MustBuild()

	err = s.msgbus.PublishRequest(route, usageMsg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			usageMsg, route, err.Error())
	}

	return err
}

func unmarshalSimManagerSimAllocate(msg *anypb.Any) (*pb.AllocateSimResponse, error) {
	p := &pb.AllocateSimResponse{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AllocateSim message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalOperatorCdrCreate(msg *anypb.Any) (*epb.EventOperatorCdrReport, error) {
	p := &epb.EventOperatorCdrReport{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal EventOperatorCdrReport message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalUkamaAgentCdrCreate(msg *anypb.Any) (*epb.CDRReported, error) {
	p := &epb.CDRReported{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal UkamaAgent CDRReprted message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}
