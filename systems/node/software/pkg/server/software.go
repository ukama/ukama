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
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	healthpb "github.com/ukama/ukama/systems/node/health/pb/gen"
	opmonpb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg"
	swclient "github.com/ukama/ukama/systems/node/software/pkg/client"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
	"github.com/ukama/ukama/systems/node/software/providers"
	opmgrpb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DefaultUpdateWatchInterval = 32 * time.Second
const DefaultUpdateWatchExpiry = 10 * time.Minute

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
	nodeGwIPs            []string
	opManager            swclient.OperationManager
	opMonitor            swclient.OperationMonitor
	opLeaseSecs          uint32
	opDeadlineSecs       uint32
}

func NewSoftwareServer(orgName string, sRepo db.SoftwareRepo, appRepo db.AppRepo, nodeRepo db.NodeRepo, healthClient providers.HealthClientProvider, msgBus mb.MsgBusServiceClient, debug bool, nodeGwIP []string, opMgr swclient.OperationManager, opMon swclient.OperationMonitor, leaseSecs, deadlineSecs uint32) *SoftwareServer {
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
		opManager:            opMgr,
		opMonitor:            opMon,
		opLeaseSecs:          leaseSecs,
		opDeadlineSecs:       deadlineSecs,
	}
}

func (s *SoftwareServer) acquireAndRegister(actionType, resourceKey string) (*opmgrpb.Operation, error) {
	if s.opManager == nil || s.opMonitor == nil {
		log.Warnf("%s running without operation manager/monitor for %s", actionType, resourceKey)
		return &opmgrpb.Operation{
			Id:          "",
			ResourceKey: resourceKey,
		}, nil
	}

	startResp, err := s.opManager.Start(&opmgrpb.StartOperationRequest{
		Type:         actionType,
		System:       "node",
		ResourceKey:  resourceKey,
		RequestedBy:  pkg.ServiceName,
		LeaseSeconds: s.opLeaseSecs,
	})
	if err != nil {
		log.Warnf("%s lock acquire for %s rejected: %v", actionType, resourceKey, err)
		return nil, err
	}
	op := startResp.Operation
	if _, err := s.opMonitor.Register(&opmonpb.RegisterIntentRequest{
		OperationId:     op.Id,
		ResourceKey:     resourceKey,
		ActionType:      actionType,
		FencingToken:    op.FencingToken,
		DeadlineSeconds: s.opDeadlineSecs,
	}); err != nil {
		log.Errorf("%s register intent for op %s failed: %v", actionType, op.Id, err)
		s.failOperation(op, actionType, fmt.Sprintf("register intent failed: %v", err))
		return nil, status.Errorf(codes.Internal, "register intent: %v", err)
	}
	log.Infof("%s acquired lock op=%s token=%d for %s", actionType, op.Id, op.FencingToken, resourceKey)
	return op, nil
}

func (s *SoftwareServer) markRunning(op *opmgrpb.Operation, actionType string) error {
	if s.opManager == nil || op == nil || op.Id == "" {
		return nil
	}

	if _, err := s.opManager.MarkRunning(op.Id, op.FencingToken); err != nil {
		log.Warnf("%s mark running failed for op %s: %v", actionType, op.Id, err)
		return err
	}
	return nil
}

func (s *SoftwareServer) failOperation(op *opmgrpb.Operation, actionType, reason string) {
	if s.opManager == nil || op == nil || op.Id == "" {
		return
	}
	if _, err := s.opManager.FailOperation(op.Id, pkg.ServiceName, reason); err != nil {
		log.Errorf("%s failed to mark operation %s failed: %v", actionType, op.Id, err)
	}
}

func (s *SoftwareServer) CreateApp(ctx context.Context, req *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	log.Infof("Creating app with name: %s, space: %s, notes: %s, metricsKeys: %v", req.Name, req.Space, req.Notes, req.MetricsKeys)
	app := db.App{
		Name:        req.Name,
		Space:       req.Space,
		Notes:       req.Notes,
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

	sw := list[0]

	if validation.IsVersionMismatch(sw.DesiredVersion, req.Tag) {
		log.Infof("Requested tag %s does not match desired version %s", req.Tag, sw.DesiredVersion)
		return &pb.UpdateSoftwareResponse{Message: "Invalid software version provided"}, nil
	}

	if !validation.IsVersionMismatch(sw.CurrentVersion, req.Tag) {
		log.Infof("Software %s already at or above version %s for node %s", req.Name, req.Tag, nId.String())
		return &pb.UpdateSoftwareResponse{Message: "Software is already up to date"}, nil
	}

	log.Infof("Node gw ips: %v", s.nodeGwIPs)
	if len(s.nodeGwIPs) == 0 {
		return nil, status.Errorf(codes.Internal, "failed to get node gw ip: no node gw ip found")
	}
	hosts := make([]string, 0, len(s.nodeGwIPs))
	for _, ip := range s.nodeGwIPs {
		hosts = append(hosts, fmt.Sprintf("http://%s:8080", ip))
	}

	target := fmt.Sprintf("%s...%s", s.orgName, nId.String())
	path := "/starter/v1/update"

	log.Infof("Publishing update for software %s to version %s on node %s using hub %s",
		req.Name, req.Tag, nId.String(), hosts)

	jsonBody := struct {
		Name string   `json:"name"`
		Tag  string   `json:"tag"`
		Hub  []string `json:"hub"`
	}{
		Name: req.Name,
		Tag:  req.Tag,
		Hub:  hosts,
	}

	data, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal update request: %v", err)
	}

	op, err := s.acquireAndRegister("UpdateSoftware", "node:"+nId.String())
	if err != nil {
		return nil, err
	}
	if err := s.markRunning(op, "UpdateSoftware"); err != nil {
		s.failOperation(op, "UpdateSoftware", fmt.Sprintf("mark running failed: %v", err))
		return nil, status.Errorf(codes.Internal, "mark running: %v", err)
	}
	if err := s.publishMessage(target, "POST", path, nId.String(), data); err != nil {
		log.Errorf("Failed to publish update message: %v", err)
		s.failOperation(op, "UpdateSoftware", fmt.Sprintf("publish failed: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to publish update message: %v", err)
	}

	sw.ChangeLogs = append(sw.ChangeLogs, "Updating app "+req.Name+" to version "+req.Tag)
	sw.Status = ukama.SoftwareStatusType(ukama.UpdateInProgress)

	// TODO(item 9/10): keep this legacy synchronous software status update for now.
	// Replace with async completion once operation-monitor verifies target version.
	// sw.CurrentVersion = req.Tag
	// sw.ChangeLogs = append(sw.ChangeLogs, "Software updated to version "+req.Tag)
	// sw.Status = ukama.SoftwareStatusType(ukama.UpToDate)

	if err := s.sRepo.Update(sw); err != nil {
		log.Errorf("Failed to persist software update: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update software: %v", err)
	}

	log.Infof("Software %s updated to %s for node %s", req.Name, req.Tag, nId.String())

	expiry := time.Now().Add(DefaultUpdateWatchExpiry)
	go s.watchSoftwareUpdate(sw.Id, nId.String(), req.Name, req.Tag, expiry, DefaultUpdateWatchInterval)

	return &pb.UpdateSoftwareResponse{Message: "Software updated dipatched successfully"}, nil

	// return &pb.UpdateSoftwareResponse{Message: "Software updated successfully", OperationId: op.Id, ResourceKey: op.ResourceKey, Status: opmgrpb.OperationStatus_RUNNING.String()}, nil
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

func (s *SoftwareServer) watchSoftwareUpdate(recordID uuid.UUID, nodeID, appName, desiredVersion string, expiry time.Time, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Infof("watchSoftwareUpdate: started for record=%s node=%s app=%s desiredVersion=%s expiry=%s",
		recordID, nodeID, appName, desiredVersion, expiry.Format(time.RFC3339))

	for {
		<-ticker.C

		healthClient, err := s.healthClient.GetClient()
		if err != nil {
			log.Errorf("watchSoftwareUpdate: failed to get health client for node=%s: %v", nodeID, err)
		} else {
			resp, err := healthClient.ListApps(context.Background(), &healthpb.ListAppsRequest{
				NodeId:  nodeID,
				AppName: appName,
			})
			if err != nil {
				log.Errorf("watchSoftwareUpdate: ListApps failed for node=%s app=%s: %v", nodeID, appName, err)
			} else {
				for _, app := range resp.GetApps() {
					log.Infof("watchSoftwareUpdate: app=%s, version=%s, desiredVersion=%s, isMismatch=%v", app.GetName(), app.GetVersion(), desiredVersion, validation.IsVersionMismatch(app.GetVersion(), desiredVersion))
					if app.GetName() == appName && !validation.IsVersionMismatch(app.GetVersion(), desiredVersion) {
						log.Infof("watchSoftwareUpdate: record=%s node=%s app=%s reached desired version %s, marking up-to-date",
							recordID, nodeID, appName, desiredVersion)
						s.persistSoftwareStatus(recordID, nodeID, appName, ukama.UpToDate,
							fmt.Sprintf("Software successfully updated to version %s", desiredVersion))
						return
					}
				}
			}
		}

		// Version not confirmed yet — fail if the deadline has now passed.
		if time.Now().After(expiry) {
			log.Warnf("watchSoftwareUpdate: deadline reached for record=%s node=%s app=%s, marking update failed",
				recordID, nodeID, appName)
			s.persistSoftwareStatus(recordID, nodeID, appName, ukama.UpdateFailed,
				fmt.Sprintf("Update timed out waiting for version %s", desiredVersion))
			return
		}

		log.Debugf("watchSoftwareUpdate: node=%s app=%s not yet at version %s, waiting %s",
			nodeID, appName, desiredVersion, interval)
	}
}

// persistSoftwareStatus fetches the software record by its primary key and writes the new
// status and a changelog entry directly, avoiding an extra List query.
func (s *SoftwareServer) persistSoftwareStatus(recordID uuid.UUID, nodeID, appName string, newStatus ukama.SoftwareStatusType, changeLog string) {
	sw, err := s.sRepo.Get(recordID)
	if err != nil {
		log.Errorf("persistSoftwareStatus: failed to get record %s for node=%s app=%s: %v",
			recordID, nodeID, appName, err)
		return
	}

	sw.Status = newStatus
	sw.ChangeLogs = append(sw.ChangeLogs, changeLog)
	if newStatus == ukama.UpToDate {
		sw.CurrentVersion = sw.DesiredVersion
	}

	log.Infof("persistSoftwareStatus: node %s, current version %s, desired version %s, status %s",
		nodeID, sw.CurrentVersion, sw.DesiredVersion, newStatus)
	if err := s.sRepo.Update(&sw); err != nil {
		log.Errorf("persistSoftwareStatus: failed to update status to %s for node=%s app=%s: %v",
			newStatus, nodeID, appName, err)
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
