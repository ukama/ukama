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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/validation"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
	"github.com/ukama/ukama/systems/node/software/providers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SoftwareServer struct {
	pb.UnimplementedSoftwareServiceServer
	sRepo                db.SoftwareRepo
	appRepo              db.AppRepo
	nodeRepo             db.NodeRepo
	nodeFeederRoutingKey msgbus.RoutingKeyBuilder
	msgbus               mb.MsgBusServiceClient
	healthClient         providers.HealthClientProvider
	debug                bool
	orgName              string
	nodeGwIPs             []string
}

func NewSoftwareServer(orgName string, sRepo db.SoftwareRepo, appRepo db.AppRepo, nodeRepo db.NodeRepo, healthClient providers.HealthClientProvider, msgBus mb.MsgBusServiceClient, debug bool, nodeGwIP []string) *SoftwareServer {
	return &SoftwareServer{
		sRepo:                sRepo,
		debug:                debug,
		msgbus:               msgBus,
		appRepo:              appRepo,
		nodeRepo:             nodeRepo,
		healthClient:         healthClient,
		orgName:              orgName,
		nodeGwIPs:            nodeGwIP,
		nodeFeederRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (s *SoftwareServer) CreateApp(ctx context.Context, req *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	log.Infof("Creating app with name: %s, space: %s, notes: %s, metricsKeys: %v", req.Name, req.Space, req.Notes, req.MetricsKeys)
	app := db.App{
		Name: req.Name,
		Space: req.Space,
		Notes: req.Notes,
		MetricsKeys: req.MetricsKeys,
	}
	err := s.appRepo.Create(app)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create app: %v", err)
	}
	return &pb.CreateAppResponse{Message: "App created successfully"}, nil
}

func (s *SoftwareServer) GetAppList(ctx context.Context, req *pb.GetAppListRequest) (*pb.GetAppListResponse, error) {
	log.Infof("Getting apps list")
	apps, err := s.appRepo.GetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get app list: %v", err)
	}
	appsPb := make([]*pb.App, len(apps))
	for i, app := range apps {
		appsPb[i] = dbAppToPbApp(&app)
	}
	return &pb.GetAppListResponse{Apps: appsPb}, nil
}

func (s *SoftwareServer) GetSoftwareList(ctx context.Context, req *pb.GetSoftwareListRequest) (*pb.GetSoftwareListResponse, error) {
	log.Infof("Getting software list with args: %s, %d, %s", req.NodeId, req.Status, req.AppName)
	var nId string
	if req.NodeId != "" {
		ukamaNodeId, err := ukama.ValidateNodeId(req.NodeId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format of node id. Error %s", err.Error())
		}
		nId = ukamaNodeId.String()
	}

	log.Infof("List req for args: %s, %d, %s", nId, ukama.SoftwareStatusType(req.Status), req.AppName)

	software, err := s.sRepo.List(nId, ukama.SoftwareStatusType(req.Status), req.AppName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get software list: %v", err)
	}
	softwarePb := make([]*pb.Software, len(software))
	for i, software := range software {
		softwarePb[i] = dbSoftwareToPbSoftware(software)
	}
	return &pb.GetSoftwareListResponse{Software: softwarePb}, nil
}

func (s *SoftwareServer) UpdateSoftware(ctx context.Context, req *pb.UpdateSoftwareRequest) (*pb.UpdateSoftwareResponse, error) {
	log.Infof("Updating software for node %s app %s to tag %s", req.NodeId, req.Name, req.Tag)

	nId, err := ukama.ValidateNodeId(req.NodeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid node id: %s", err.Error())
	}

	list, err := s.sRepo.List(nId.String(), ukama.UpdateAvailable, req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get software: %v", err)
	}
	if len(list) == 0 {
		return nil, status.Errorf(codes.NotFound, "software not found or already up to date")
	}

	// Unique index on (node_id, app_name) implies at most one record for this request
	sw := list[0]

	if validation.IsVersionMismatch(sw.DesiredVersion, req.Tag) {
		log.Infof("Requested tag %s does not match desired version %s", req.Tag, sw.DesiredVersion)
		return &pb.UpdateSoftwareResponse{Message: "Invalid software version provided"}, nil
	}

	if !validation.IsVersionMismatch(sw.CurrentVersion, req.Tag) {
		log.Infof("Software %s already at or above version %s for node %s", req.Name, req.Tag, nId.String())
		return &pb.UpdateSoftwareResponse{Message: "Software is already up to date"}, nil
	}

	target := fmt.Sprintf("%s...%s", s.orgName, nId.String())
	path := fmt.Sprintf("/starter/v1/update/%s/%s", req.Name, req.Tag)
	log.Infof("Publishing update for software %s to version %s on node %s", req.Name, req.Tag, nId.String())

	nodeGwIP, err := s.getNodeGwIP()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get node gw ip: %v", err)
	}
	log.Infof("Node gw ip: %s", nodeGwIP)
	jsonBody := map[string]string{"host": fmt.Sprintf("http://%s:8080", nodeGwIP)}
	data, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get software: %v", err)
	}
	if len(list) == 0 {
		return nil, status.Errorf(codes.NotFound, "software not found or already up to date")
	}
	
	if err := s.publishMessage(target, "POST", path, nId.String(), data); err != nil {
		log.Errorf("Failed to publish update message: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to publish update message: %v", err)
	}
	sw.CurrentVersion = req.Tag
	sw.ChangeLogs = append(sw.ChangeLogs, "Software updated to version "+req.Tag)
	sw.Status = ukama.SoftwareStatusType(ukama.UpToDate)
	if err := s.sRepo.Update(sw); err != nil {
		log.Errorf("Failed to persist software update: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update software: %v", err)
	}

	log.Infof("Software %s updated to %s for node %s", req.Name, req.Tag, nId.String())
	return &pb.UpdateSoftwareResponse{Message: "Software updated successfully"}, nil
}

func dbSoftwareToPbSoftware(software *db.Software) *pb.Software {
	return &pb.Software{
		Id:             software.Id.String(),
		ReleaseDate:    software.ReleaseDate.Format(time.RFC3339),
		Status:         ukama.SoftwareStatusType(software.Status).String(),
		CurrentVersion: software.CurrentVersion,
		DesiredVersion: software.DesiredVersion,
		Name:           software.App.Name,
		Space:          software.App.Space,
		Notes:          software.App.Notes,
		MetricsKeys:    software.App.MetricsKeys,
		NodeId:         software.NodeId,
		ChangeLog:      software.ChangeLogs,
		CreatedAt:      software.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      software.UpdatedAt.Format(time.RFC3339),
	}
}

func dbAppToPbApp(app *db.App) *pb.App {
	return &pb.App{
		Name:        app.Name,
		Space:       app.Space,
		Notes:       app.Notes,
		MetricsKeys: app.MetricsKeys,
	}
}

func (c *SoftwareServer) publishMessage(target string, method string, path string, nodeId string, data []byte) error {
	route := "request.cloud.local" + "." + c.orgName + "." + pkg.SystemName + "." + pkg.ServiceName + "." + "nodefeeder" + "." + "publish"
	msg := &epb.NodeFeederMessage{
		Target:     target,
		HttpMethod: method,
		Path:       path,
		Msg:        data,
		NodeId:     nodeId,
	}
	log.Infof("Published software update node %s on path %s on target %s ", nodeId, path, target)
	err := c.msgbus.PublishRequest(route, msg)
	return err
}

func (c *SoftwareServer) getNodeGwIP() (string, error) {
	if len(c.nodeGwIPs) == 0 {
		return "", errors.New("no node gw ip found")
	}

	for _, ip := range c.nodeGwIPs {
		log.Infof("validating IP : %s", ip)
		if net.ParseIP(ip) != nil {
			return ip, nil
		}
	}
	
	return "", errors.New("no valid node gw ip found")
}
