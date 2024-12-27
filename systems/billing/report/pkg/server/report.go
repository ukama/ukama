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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/billing/report/pkg"
	"github.com/ukama/ukama/systems/billing/report/pkg/db"
	"github.com/ukama/ukama/systems/billing/report/pkg/util"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/billing/report/pb/gen"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

type ReportServer struct {
	OrgName          string
	OrgId            uuid.UUID
	reportRepo       db.ReportRepo
	subscriberClient csub.SubscriberClient
	msgBus           mb.MsgBusServiceClient
	baseRoutingKey   msgbus.RoutingKeyBuilder
	pb.UnimplementedReportServiceServer
}

func NewReportServer(orgName, org string, reportRepo db.ReportRepo, subscriberClient csub.SubscriberClient, msgBus mb.MsgBusServiceClient) *ReportServer {
	orgId, err := uuid.FromString(org)
	if err != nil {
		panic(fmt.Sprintf("invalid format of org uuid: %s", org))
	}

	return &ReportServer{
		OrgName:          orgName,
		OrgId:            orgId,
		reportRepo:       reportRepo,
		subscriberClient: subscriberClient,
		msgBus:           msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().
			SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (r *ReportServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.ReportResponse, error) {
	log.Infof("Unmarshalling raw report from webhook: %v", req.RawReport)

	rwInvoceStruct := &util.RawInvoice{}

	err := json.Unmarshal([]byte(req.RawReport), rwInvoceStruct)
	if err != nil {
		log.Errorf("Failed to unmarshal RawReport JSON to RawReport struct %v", err)

		return nil, status.Errorf(codes.InvalidArgument,
			"failed to unmarshal RawReport JSON paylod from webhook. Error %s", err)
	}

	if rwInvoceStruct == nil || rwInvoceStruct.Customer == nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid JSON format of RawReport. Error %s", req.RawReport)
	}

	ownerId, err := uuid.FromString(rwInvoceStruct.Customer.ExternalID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of owner uuid. Error %s", err.Error())
	}

	report := &db.Report{
		OwnerId:   ownerId,
		OwnerType: ukama.OwnerTypeOrg,
		Type:      ukama.ReportTypeInvoice,
	}

	if ownerId != r.OrgId {
		subscriberInfo, err := r.subscriberClient.Get(ownerId.String())
		if err != nil {
			return nil, err
		}

		report.NetworkId = subscriberInfo.NetworkId
		report.OwnerType = ukama.OwnerTypeSubscriber
		report.Type = ukama.ReportTypeConsumption
	}

	rwInvoceStruct.FileURL = fmt.Sprintf("http://{API_ENDPOINT}/pdf/%s.pdf", report.Id.String())

	rwReportBytes, err := json.Marshal(rwInvoceStruct)
	if err != nil {
		log.Errorf("Failed to marshal RawReport struct to RawReport JSON %v", err)

		return nil, fmt.Errorf("failed to marshal RawReport struct to RawReport JSON: %w", err)
	}

	report.RawReport = datatypes.JSON(rwReportBytes)

	log.Infof("Adding report for owner: %s", ownerId)
	err = r.reportRepo.Add(report, func(*db.Report, *gorm.DB) error {
		report.Id = uuid.NewV4()
		report.Period = time.Now().UTC()

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "report")
	}

	pbReport := dbReportToPbReport(report)

	resp := &pb.ReportResponse{
		Report: pbReport,
	}

	route := r.baseRoutingKey.SetAction("generate").SetObject(report.Type.String()).MustBuild()

	val := &epb.RawReport{}
	m := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	err = m.Unmarshal([]byte(report.RawReport.String()), val)
	if err != nil {
		log.Errorf("Failed to unmarshal RawReport JSON to epb.RawReport proto: %v", err)

		return nil, status.Errorf(codes.InvalidArgument,
			"failed to unmarshal RawReport JSON paylod to epb.RawReport. Error %s", err)
	}

	evt := &epb.Report{
		Id:        pbReport.Id,
		OwnerId:   pbReport.OwnerId,
		OwnerType: pbReport.OwnerType,
		NetworkId: pbReport.NetworkId,
		Type:      pbReport.Type,
		Period:    pbReport.Period,
		RawReport: val,
		IsPaid:    pbReport.IsPaid,
		CreatedAt: pbReport.CreatedAt,
	}

	err = r.msgBus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			req, route, err.Error())
	}

	return resp, nil
}

func (r *ReportServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.ReportResponse, error) {
	reportId, err := uuid.FromString(req.ReportId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of report uuid. Error %s", err.Error())
	}

	report, err := r.reportRepo.Get(reportId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "report")
	}

	return &pb.ReportResponse{
		Report: dbReportToPbReport(report),
	}, nil
}

func (r *ReportServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Infof("Getting reports matching: %v", req)

	if req.OwnerId != "" {
		ownerId, err := uuid.FromString(req.GetOwnerId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format for %s owner uuid: %s. Error %v",
				req.OwnerType, req.OwnerId, err)
		}

		req.OwnerId = ownerId.String()
	}

	ownerType := ukama.OwnerTypeUnknown
	if req.OwnerType != "" {
		ownerType = ukama.ParseOwnerType(req.OwnerType)
		if ownerType == ukama.OwnerTypeUnknown {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid value for owner type: %s", req.OwnerType)
		}
	}

	if req.NetworkId != "" {
		networkId, err := uuid.FromString(req.GetNetworkId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid format for network uuid: %s. Error %v", req.NetworkId, err)
		}

		req.NetworkId = networkId.String()
	}

	reportType := ukama.ReportTypeUnknown
	if req.ReportType != "" {
		reportType = ukama.ParseReportType(req.ReportType)
		if reportType == ukama.ReportTypeUnknown {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid value for report type: %s", req.ReportType)
		}
	}

	reports, err := r.reportRepo.List(req.OwnerId, ownerType, req.NetworkId, reportType, req.IsPaid, req.Count, req.Sort)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "reports")
	}

	return &pb.ListResponse{Reports: dbReportsToPbReports(reports)}, nil
}

func (r *ReportServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.ReportResponse, error) {
	report, err := update(req.ReportId, req.IsPaid, r.reportRepo, r.msgBus, r.baseRoutingKey)

	if err != nil {
		return nil, err
	}

	return &pb.ReportResponse{Report: dbReportToPbReport(report)}, nil
}

func (r *ReportServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Infof("Deleting report %s", req.ReportId)

	reportId, err := uuid.FromString(req.ReportId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of report uuid. Error %s", err.Error())
	}

	err = r.reportRepo.Delete(reportId, nil)
	if err != nil {
		log.Error(err)

		return nil, grpc.SqlErrorToGrpc(err, "report")
	}

	route := r.baseRoutingKey.SetAction("delete").SetObject("invoice").MustBuild()

	err = r.msgBus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			req, route, err.Error())
	}

	return &pb.DeleteResponse{}, nil
}

func update(reportId string, isPaid bool, reportRepo db.ReportRepo, msgBus mb.MsgBusServiceClient,
	baseRoutingKey msgbus.RoutingKeyBuilder) (*db.Report, error) {

	log.Infof("Updating report: %v", reportId)

	repId, err := uuid.FromString(reportId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format for payment uuid: %s format for report update. Error %v", reportId, err)
	}

	report, err := reportRepo.Get(repId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "report")
	}

	report.IsPaid = isPaid

	err = reportRepo.Update(report, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "report")
	}

	route := baseRoutingKey.SetAction("update").SetObject(report.Type.String()).MustBuild()

	val := &epb.RawReport{}
	m := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	err = m.Unmarshal([]byte(report.RawReport.String()), val)
	if err != nil {
		log.Errorf("Failed to unmarshal RawReport JSON to epb.RawReport proto: %v", err)

		return nil, status.Errorf(codes.InvalidArgument,
			"failed to unmarshal RawReport JSON paylod to epb.RawReport. Error %s", err)
	}

	evt := &epb.Report{
		Id:        report.Id.String(),
		OwnerId:   report.OwnerId.String(),
		OwnerType: report.OwnerType.String(),
		NetworkId: report.NetworkId.String(),
		Type:      report.Type.String(),
		Period:    report.Period.Format(time.RFC3339),
		RawReport: val,
		IsPaid:    report.IsPaid,
		CreatedAt: report.CreatedAt.Format(time.RFC3339),
	}

	err = msgBus.PublishRequest(route, evt)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			evt, route, err.Error())
	}

	return report, nil
}

func dbReportToPbReport(report *db.Report) *pb.Report {
	inv := &pb.Report{
		Id:        report.Id.String(),
		OwnerId:   report.OwnerId.String(),
		OwnerType: report.OwnerType.String(),
		NetworkId: report.NetworkId.String(),
		Type:      report.Type.String(),
		Period:    report.Period.String(),
		IsPaid:    report.IsPaid,
		CreatedAt: report.CreatedAt.String(),
	}

	if report.NetworkId != uuid.Nil {
		inv.NetworkId = report.NetworkId.String()
	}

	val := &pb.RawReport{}

	m := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	err := m.Unmarshal([]byte(report.RawReport.String()), val)
	//TODO: error handling

	if err == nil {
		inv.RawReport = val
	}

	return inv
}

func dbReportsToPbReports(reports []db.Report) []*pb.Report {
	res := []*pb.Report{}

	for _, r := range reports {
		res = append(res, dbReportToPbReport(&r))
	}

	return res
}
