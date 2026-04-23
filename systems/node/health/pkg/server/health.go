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

	"github.com/cloudflare/cfssl/log"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/health/pb/gen"
	"github.com/ukama/ukama/systems/node/health/pkg"
	"github.com/ukama/ukama/systems/node/health/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthServer struct {
	pb.UnimplementedHealhtServiceServer
	sRepo            db.HealthRepo
	healthRoutingKey msgbus.RoutingKeyBuilder
	msgbus           mb.MsgBusServiceClient
	debug            bool
	orgName          string
}

func NewHealthServer(orgName string, sRepo db.HealthRepo, msgBus mb.MsgBusServiceClient, debug bool) *HealthServer {
	return &HealthServer{
		sRepo:            sRepo,
		orgName:          orgName,
		healthRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
		msgbus:           msgBus,
		debug:            debug,
	}
}

func (h *HealthServer) StoreRunningAppsInfo(ctx context.Context, req *pb.StoreRunningAppsInfoRequest) (*pb.StoreRunningAppsInfoResponse, error) {
	log.Infof("StoreRunningAppsInfo: %v", req)
	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	healthID := uuid.NewV4()
	cappID := uuid.NewV4()

	// Create a Health instance
	health := db.Health{
		Id:        healthID,
		NodeId:    nId.StringLowercase(),
		TimeStamp: req.GetTimestamp(),
	}

	// Populate the System array from the request
	for _, sys := range req.GetSystem() {
		health.System = append(health.System, db.System{
			Id:       uuid.NewV4(),
			HealthID: healthID,
			Name:     sys.GetName(),
			Value:    sys.GetValue(),
		})
	}

	for _, capp := range req.GetCapps() {
		health.Capps = append(health.Capps, db.Capp{
			Id:       cappID,
			HealthID: healthID,
			Space:    capp.GetSpace(),
			Name:     capp.GetName(),
			Tag:      capp.GetTag(),
			Status:   db.Status(capp.GetStatus()),
		})

		for _, resource := range capp.GetResources() {
			health.Capps[len(health.Capps)-1].Resources = append(health.Capps[len(health.Capps)-1].Resources, db.Resource{
				Id:     uuid.NewV4(),
				CappID: cappID,
				Name:   resource.GetName(),
				Value:  resource.GetValue(),
			})
		}
	}

	err = h.sRepo.StoreRunningAppsInfo(&health, nil)
	if err != nil {
		return nil, err
	}

	msg := &epb.StoreRunningAppsInfoEvent{
		NodeId:    req.NodeId,
		Timestamp: req.Timestamp,
		System:    []*epb.System{},
		Capps:     []*epb.Capps{},
	}
	for _, sys := range health.System {
		msg.System = append(msg.System, &epb.System{
			Id:       sys.Id.String(),
			HealthId: health.Id.String(),
			Name:     sys.Name,
			Value:    sys.Value,
		})
	}

	for _, capp := range health.Capps {
		capps := &epb.Capps{
			Id:        capp.Id.String(),
			Space:     capp.Space,
			Name:      capp.Name,
			Tag:       capp.Tag,
			Status:    epb.Status(capp.Status), 
		}
		for _, resource := range capp.Resources {
			resource := &epb.Resource{
				Id:     resource.Id.String(),
				Name:   resource.Name,
				Value:  resource.Value,
			}
			capps.Resources = append(capps.Resources, resource)
		}
		msg.Capps = append(msg.Capps, capps)
	}	

	route := h.healthRoutingKey.SetAction("store").SetObject("capps").MustBuild()
	err = h.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.StoreRunningAppsInfoResponse{}, nil
}

func (h *HealthServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Infof("List: %v", req)
	if req.Id == "" && req.NodeId == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			"either provide id or node id")
	}

	timeframe := ukama.ParseFilterTimeframesType(strings.ToLower(req.Timeframe.String()))
	healths, err := h.sRepo.List(req.Id, req.NodeId, req.Timestamp, timeframe)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "health")
	}

	healthsPb := convertToPbHealths(healths)

	return &pb.ListResponse{
		Healths: healthsPb,
	}, nil
}

func convertToPbHealths(healths []*db.Health) []*pb.Health {
	healthsPb := make([]*pb.Health, len(healths))
	for i, h := range healths {
		healthsPb[i] = convertToPbHealth(h)
	}
	return healthsPb
}

func convertToPbHealth(health *db.Health) *pb.Health {
	systems := make([]*pb.System, len(health.System))
	for i, s := range health.System {
		systems[i] = convertToPbSystem(&s)
	}
	capps := make([]*pb.Capps, len(health.Capps))
	for i, c := range health.Capps {
		capps[i] = convertToPbCapp(&c)
	}
	return &pb.Health{
		Id:        health.Id.String(),
		NodeId:    health.NodeId,
		Timestamp: health.TimeStamp,
		System:    systems,
		Capps:     capps,
	}
}

func convertToPbSystem(system *db.System) *pb.System {
	return &pb.System{
		Id:       system.Id.String(),
		HealthId: system.HealthID.String(),
		Name:     system.Name,
		Value:    system.Value,
	}
}

func convertToPbResource(resources *db.Resource) *pb.Resource {
	return &pb.Resource{
		Id:     resources.Id.String(),
		Name:   resources.Name,
		Value:  resources.Value,
		CappId: resources.CappID.String(),
	}
}

func convertToPbCapp(capp *db.Capp) *pb.Capps {
	resources := make([]*pb.Resource, len(capp.Resources))
	for i, r := range capp.Resources {
		resources[i] = convertToPbResource(&r)
	}
	return &pb.Capps{
		Id:        capp.Id.String(),
		Space:     capp.Space,
		Name:      capp.Name,
		Tag:       capp.Tag,
		Status:    pb.Status(capp.Status),
		Resources: resources,
	}
}