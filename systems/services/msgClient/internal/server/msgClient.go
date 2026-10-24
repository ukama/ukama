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
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/common/pb/gen/msgclient"
)

type Registry interface {
	Register(name, system, instance, uri string, routes []string) (string, error)
	Resume(uuid string) error
	Pause(uuid string) error
	Unregister(uuid string) error
	Lookup(uuid string) (name, instance string, ok bool)
}

type Publisher interface {
	Publish(key string, body []byte, headers rabbitmq.Table) error
}

type MsgClientServer struct {
	sys             string
	refreshInterval time.Duration
	r               Registry
	pub             Publisher
	p               msgbus.MsgBusShovelProvider
	pb.UnimplementedMsgClientServiceServer
}

func NewMsgClientServer(r Registry, pub Publisher, p msgbus.MsgBusShovelProvider, sys string, refreshInterval time.Duration) *MsgClientServer {
	return &MsgClientServer{
		sys:             sys,
		refreshInterval: refreshInterval,
		r:               r,
		pub:             pub,
		p:               p,
	}
}

func (m *MsgClientServer) RegisterService(ctx context.Context, req *pb.RegisterServiceReq) (*pb.RegisterServiceResp, error) {
	if !strings.EqualFold(m.sys, req.SystemName) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid system name %s in request", req.SystemName)
	}

	if _, err := msgbus.ParseRouteList(req.Routes); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid routes for %s: %v", req.ServiceName, err)
	}

	id, err := m.r.Register(req.ServiceName, req.SystemName, req.InstanceId, req.ServiceURI, req.Routes)
	if err != nil {
		log.Errorf("Failed to register service %s. Error %s", req.ServiceName, err.Error())
		return nil, status.Errorf(codes.Unavailable, "failed to register %s: %v", req.ServiceName, err)
	}

	return &pb.RegisterServiceResp{
		State:           pb.REGISTRAION_STATUS_REGISTERED,
		ServiceUuid:     id,
		RefreshInterval: uint32(max(m.refreshInterval/time.Second, 1)),
	}, nil
}

func (m *MsgClientServer) StartMsgBusHandler(ctx context.Context, req *pb.StartMsgBusHandlerReq) (*pb.StartMsgBusHandlerResp, error) {
	if err := m.r.Resume(req.ServiceUuid); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.StartMsgBusHandlerResp{}, nil
}

func (m *MsgClientServer) StopMsgBusHandler(ctx context.Context, req *pb.StopMsgBusHandlerReq) (*pb.StopMsgBusHandlerResp, error) {
	if err := m.r.Pause(req.ServiceUuid); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.StopMsgBusHandlerResp{}, nil
}

func (m *MsgClientServer) UnregisterService(ctx context.Context, req *pb.UnregisterServiceReq) (*pb.UnregisterServiceResp, error) {
	if err := m.r.Unregister(req.ServiceUuid); err != nil {
		return nil, err
	}

	return &pb.UnregisterServiceResp{}, nil
}

func (m *MsgClientServer) PublishMsg(ctx context.Context, req *pb.PublishMsgRequest) (*pb.PublishMsgResponse, error) {
	body, err := proto.Marshal(req.Msg)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid message: %v", err)
	}

	name, instance, ok := m.r.Lookup(req.ServiceUuid)
	if !ok {
		name = "unknown"
	}

	err = m.pub.Publish(req.RoutingKey, body, rabbitmq.Table{
		"source-service": name,
		"instance-id":    instance,
	})
	if err != nil {
		log.Errorf("Failed to publish %s from %s. Error %s", req.RoutingKey, name, err.Error())
		return nil, status.Errorf(codes.Unavailable, "event %s not accepted: %v", req.RoutingKey, err)
	}

	return &pb.PublishMsgResponse{}, nil
}

func (m *MsgClientServer) CreateShovel(ctx context.Context, in *pb.CreateShovelRequest) (*pb.CreateShovelResponse, error) {
	err := m.p.CreateShovel(in.Name, &msgbus.Shovel{
		SrcProtocol:    in.SrcProtocol,
		DestProtocol:   in.DestProtocol,
		SrcUri:         in.SrcUri,
		DestUri:        in.DestUri,
		SrcExchangeKey: in.SrcExchangeKey,
		SrcExchange:    in.SrcExchange,
		DestExchange:   in.DestExchange,
	})

	if err != nil {
		log.Errorf("Failed to create shovel %v. Error %+v.", in, err)
	}
	return &pb.CreateShovelResponse{}, err
}

func (m *MsgClientServer) RemoveShovel(ctx context.Context, in *pb.RemoveShovelRequest) (*pb.RemoveShovelResponse, error) {
	err := m.p.RemoveShovel(in.Name)
	if err != nil {
		log.Errorf("Failed to remove shovel %s. Error %+v.", in.Name, err)
	}
	return &pb.RemoveShovelResponse{}, err
}
