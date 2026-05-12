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
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/node/health/pb/gen"
	"github.com/ukama/ukama/systems/node/health/pkg/db"
	"github.com/ukama/ukama/systems/node/health/pkg/parser"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HealthServer struct {
	pb.UnimplementedHealthServiceServer
	sRepo   db.HealthRepo
	debug   bool
	orgName string
}

func NewHealthServer(orgName string, sRepo db.HealthRepo, debug bool) *HealthServer {
	return &HealthServer{
		sRepo:   sRepo,
		orgName: orgName,
		debug:   debug,
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

	return &pb.StoreHealthReportResponse{ReportId: report.ID.String()}, nil
}

func (h *HealthServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Infof("List: %v", req)
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

	return &pb.ListResponse{Reports: out}, nil
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
