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
	"github.com/ukama/ukama/systems/node/health/pkg/parser"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)
 
type HealthServer struct {
	pb.UnimplementedHealthServiceServer
	sRepo   		 db.HealthRepo
	debug   		 bool
	orgName 		 string
	msgbus           mb.MsgBusServiceClient
	healthRoutingKey msgbus.RoutingKeyBuilder
}

func NewHealthServer(orgName string, sRepo db.HealthRepo, debug bool, msgBus mb.MsgBusServiceClient) *HealthServer {
	return &HealthServer{
		sRepo:   sRepo,
		orgName: orgName,
		debug:   debug,
		msgbus:           msgBus,
		healthRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (h *HealthServer) StoreHealthReport(ctx context.Context, req *pb.StoreHealthReportRequest) (*pb.StoreHealthReportResponse, error) {
	log.Infof("StoreHealthReport: %v", req)

	nID, err := ukama.ValidateNodeId(req.GetNodeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of node id. Error %s", err.Error())
	}

	raw := req.GetPayload()
	if len(raw) == 0 {
		raw = []byte("{}")
	}
	parsed, err := parser.ParseHealthPayload(raw)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "payload is invalid: %v", err)
	}
	if parsed.NodeType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nodeType is required in payload")
	}
	if parsed.ReportedAt.IsZero() {
		return nil, status.Errorf(codes.InvalidArgument, "reportedAt is required in payload")
	}

	nodeType := ukama.NodeType(parsed.NodeType)
	schemaVersion := parsed.SchemaVersion

	report := &db.HealthReport{
		ID:            uuid.NewV4(),
		NodeID:        nID.StringLowercase(),
		NodeType:      nodeType,
		SchemaVersion: schemaVersion,
		ReportedAt:    parsed.ReportedAt,
		Payload:       raw,
	}

	receivedAt := time.Now().UTC()
	if err := h.sRepo.StoreHealthReport(report, receivedAt); err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "health")
	}

	if h.msgbus != nil {
		route := h.healthRoutingKey.SetAction("store").SetObject("capps").MustBuild()

		evt := &epb.HealthReportEvent{
			Id:        		report.ID.String(),
			NodeId:    		report.NodeID,
			Payload: 		report.Payload,
			SchemaVersion:  report.SchemaVersion,
			NodeType:  		report.NodeType.String(),
			ReportedAt: 	report.ReportedAt.Format(time.RFC3339),
		}
		log.Infof("Publishing event %+v with key %+v", evt, route)
		err = h.msgbus.PublishRequest(route, evt)
		if err != nil {
			log.Errorf("Failed to publish message %+v with key %+v. Errors %s", evt, route, err.Error())
		}
	}

	return &pb.StoreHealthReportResponse{ReportId: report.ID.String()}, nil
}

func (h *HealthServer) ListReports(ctx context.Context, req *pb.ListReportsRequest) (*pb.ListReportsResponse, error) {
	log.Infof("ListReports: %v", req)
	if req.GetReportId() == "" && req.GetNodeId() == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			"either provide reportId or nodeId")
	}

	var reportedAt *time.Time
	if req.GetReportedAt() != nil {
		t := req.GetReportedAt().AsTime()
		reportedAt = &t
	}

	timeframe := ukama.ParseFilterTimeframesType(strings.ToLower(req.GetTimeframe().String()))
	reports, err := h.sRepo.List(req.GetReportId(), req.GetNodeId(), reportedAt, timeframe)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "health")
	}

	out := make([]*pb.HealthReport, len(reports))
	for i, r := range reports {
		out[i] = healthReportToPb(r)
	}

	return &pb.ListReportsResponse{Reports: out}, nil
}

func (h *HealthServer) ListApps(ctx context.Context, req *pb.ListAppsRequest) (*pb.ListAppsResponse, error) {
	log.Infof("ListApps: %v", req)

	reports, err := h.sRepo.List(req.GetReportId(), req.GetNodeId(), nil, ukama.FilterTimeframesTypeLatest)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "health")
	}

	if len(reports) == 0 {
		return &pb.ListAppsResponse{Apps: make([]*pb.App, 0)}, nil
	}

	parsed, err := parser.ParseHealthPayload(reports[0].Payload)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "payload is invalid: %v", err)
	}

	apps := make([]*pb.App, len(parsed.Apps))
	for i, a := range parsed.Apps {
		if req.GetAppName() != "" && a.Name == req.GetAppName() {
			apps[i] = parseAppToPb(&a)
			break
		}
		apps[i] = parseAppToPb(&a)
	}

	return &pb.ListAppsResponse{Apps: apps}, nil
}

func (h *HealthServer) ListInterfaces(ctx context.Context, req *pb.ListInterfacesRequest) (*pb.ListInterfacesResponse, error) {
	log.Infof("ListInterfaces: %v", req)

	reports, err := h.sRepo.List(req.GetReportId(), req.GetNodeId(), nil, ukama.FilterTimeframesTypeLatest)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "health")
	}
	if len(reports) == 0 {
		return &pb.ListInterfacesResponse{Interfaces: &pb.Interface{}}, nil
	}
	parsed, err := parser.ParseHealthPayload(reports[0].Payload)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "payload is invalid: %v", err)
	}

	interfaces := parseInterfaceToPb(&parsed.Interfaces)

	return &pb.ListInterfacesResponse{
		Interfaces: interfaces,
	}, nil
}

func parseInterfaceToPb(i *parser.HealthInterfaces) *pb.Interface {
	return &pb.Interface{
		Cellular:   parseCellularInterfaceToPb(i.Cellular),
		Radio:      parseRadioInterfaceToPb(i.Radio),
		Gps:        parseGPSInterfaceToPb(i.GPS),
		Backhaul:   parseBackhaulInterfaceToPb(i.Backhaul),
		Fem:        parseFEMInterfaceToPb(i.FEM),
		Switch:     parseSwitchInterfaceToPb(i.Switch),
		Controller: parseControllerInterfaceToPb(i.Controller),
	}
}

func parseCellularInterfaceToPb(c *parser.CellularInterface) *pb.CellularInterface {
	if c == nil {
		return nil
	}
	return &pb.CellularInterface{
		Available: c.Available,
		Error:     c.Error,
	}
}

func parseRadioInterfaceToPb(r *parser.RadioInterface) *pb.RadioInterface {
	if r == nil {
		return nil
	}
	return &pb.RadioInterface{
		Available: r.Available,
		State:     r.State,
	}
}

func parseGPSInterfaceToPb(g *parser.GPSInterface) *pb.GPSInterface {
	if g == nil {
		return nil
	}
	return &pb.GPSInterface{
		Available:   g.Available,
		Lock:        g.Lock,
		Coordinates: g.Coordinates,
		Time:        timestamppb.New(g.Time),
	}
}

func parseBackhaulInterfaceToPb(b *parser.BackhaulInterface) *pb.BackhaulInterface {
	if b == nil {
		return nil
	}
	return &pb.BackhaulInterface{
		Available:  b.Available,
		State:      b.State,
		LinkGuess: b.LinkGuess,
		Confidence: b.Confidence,
	}
}

func parseFEMInterfaceToPb(f *parser.FEMInterface) *pb.FEMInterface {
	if f == nil {
		return nil
	}
	return &pb.FEMInterface{
		Available: f.Available,
		Fems:      parseFEMsToPb(f.FEMs),
	}
}

func parseSwitchInterfaceToPb(s *parser.SwitchInterface) *pb.SwitchInterface {
	if s == nil {
		return nil
	}
	return &pb.SwitchInterface{
		Available:       s.Available,
		Reachable:       s.Reachable,
		State:           s.State,
		Model:           s.Model,
		SoftwareVersion: s.SoftwareVersion,
		PortCount:       int32(s.PortCount),
		Policy: &pb.SwitchInterfacePolicy{
			State:  s.Policy.State,
			Hash:   s.Policy.Hash,
			Source: s.Policy.Source,
			Error:  s.Policy.Error,
		},
		Ports: parseSwitchPortsToPb(s.Ports),
	}
}

func parseSwitchPortsToPb(p []parser.SwitchPort) []*pb.SwitchPort {
	ports := make([]*pb.SwitchPort, len(p))
	for i, p := range p {
		ports[i] = &pb.SwitchPort{
			Id:             p.ID,
			Name:           p.Name,
			Present:        p.Present,
			AdminState:     p.AdminState,
			LinkState:      p.LinkState,
			PoeState:       p.PoeState,
			PoeOperational: p.PoeOperational,
			SpeedBps:       p.SpeedBps,
			PowerWatts:     p.PowerWatts,
			Fault:          p.Fault,
		}
	}
	return ports
}

func parseFEMsToPb(f []*parser.FEMUnit) []*pb.FEMUnit {
	fems := make([]*pb.FEMUnit, len(f))
	for i, f := range f {
		fems[i] = &pb.FEMUnit{
			Unit: int32(f.Unit),
			Present: f.Present,
		}
	}
	return fems
}

func parseControllerInterfaceToPb(c *parser.NodeControllerInterface) *pb.NodeControllerInterface {
	if c == nil {
		return nil
	}
	return &pb.NodeControllerInterface{
		Available:        c.Available,
		CommOk:           c.CommOk,
		ChargeState:      c.ChargeState,
		ErrorCode:        int32(c.ErrorCode),
		Error:            c.Error,
		ActiveAlarmCount: int32(c.ActiveAlarmCount),
		Solar:            parseControllerSolarMetricsToPb(&c.Solar),
		Battery:          parseControllerBatteryMetricsToPb(&c.Battery),
		Load:             parseControllerLoadMetricsToPb(&c.Load),
	}
}

func parseControllerSolarMetricsToPb(s *parser.ControllerSolarMetrics) *pb.ControllerSolarMetrics {
	return &pb.ControllerSolarMetrics{
		VoltageV: s.VoltageV,
		CurrentA: s.CurrentA,
		PowerW: s.PowerW,
	}
}

func parseControllerBatteryMetricsToPb(b *parser.ControllerBatteryMetrics) *pb.ControllerBatteryMetrics {
	return &pb.ControllerBatteryMetrics{
		VoltageV: b.VoltageV,
		CurrentA: b.CurrentA,
		SocPct: int32(b.SocPct),
	}
}

func parseControllerLoadMetricsToPb(l *parser.ControllerLoadMetrics) *pb.ControllerLoadMetrics {
	return &pb.ControllerLoadMetrics{
		OutputOn: l.OutputOn,
		CurrentA: l.CurrentA,
	}
}

func parseAppToPb(a *parser.HealthApp) *pb.App {
	return &pb.App{
		Name: a.Name,
		Version: a.Version,
		Tag: a.Tag,
		Status: a.State,
		Resource: &pb.AppResource{
			CpuPercent: float32(a.Resources.CPUPercent),
			MemoryRssKb: float32(a.Resources.MemoryRssKb),
			DiskReadBytes: float32(a.Resources.DiskReadBytes),
			DiskWriteBytes: float32(a.Resources.DiskWriteBytes),
		},
	}
}

func healthReportToPb(r *db.HealthReport) *pb.HealthReport {
	if r == nil {
		return nil
	}
	return &pb.HealthReport{
		Id:            r.ID.String(),
		NodeId:        r.NodeID,
		NodeType:      string(r.NodeType),
		SchemaVersion: r.SchemaVersion,
		ReportedAt:    timestamppb.New(r.ReportedAt),
		ReceivedAt:    timestamppb.New(r.ReceivedAt),
		Payload:       r.Payload,
	}
}
