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
	"path/filepath"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/billing/report/pkg"
	"github.com/ukama/ukama/systems/billing/report/pkg/db"
	"github.com/ukama/ukama/systems/billing/report/pkg/pdf"
	"github.com/ukama/ukama/systems/billing/report/pkg/util"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/billing/report/pb/gen"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

const (
	defaultTemplate = "templates/invoice.html.tmpl"
	pdfFolder       = "/srv/static/"
)

type ReportServer struct {
	OrgName          string
	OrgId            uuid.UUID
	reportRepo       db.ReportRepo
	subscriberClient csub.SubscriberClient
	msgbus           mb.MsgBusServiceClient
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
		msgbus:           msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().
			SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (i *ReportServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
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

	if ownerId != i.OrgId {
		subscriberInfo, err := i.subscriberClient.Get(ownerId.String())
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
	err = i.reportRepo.Add(report, func(*db.Report, *gorm.DB) error {
		report.Id = uuid.NewV4()
		report.Period = time.Now().UTC()

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "report")
	}

	log.Infof("starting PDF generation")
	go func() {
		err = generateReportPDF(rwInvoceStruct, defaultTemplate, filepath.Join(pdfFolder, report.Id.String()+".pdf"))
		if err != nil {
			log.Errorf("PDF generation failure: failed to generate invoice PDF: %v", err)
		}

		log.Infof("finishing PDF generation")
	}()

	resp := &pb.AddResponse{
		Report: dbReportToPbReport(report),
	}

	route := i.baseRoutingKey.SetAction("generate").SetObject(report.Type.String()).MustBuild()

	err = i.msgbus.PublishRequest(route, resp.Report)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			req, route, err.Error())
	}

	return resp, nil
}

func (i *ReportServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	reportId, err := uuid.FromString(req.ReportId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of report uuid. Error %s", err.Error())
	}

	report, err := i.reportRepo.Get(reportId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "report")
	}

	return &pb.GetResponse{
		Report: dbReportToPbReport(report),
	}, nil
}

func (i *ReportServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
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

	reports, err := i.reportRepo.List(req.OwnerId, ownerType, req.NetworkId, reportType, req.IsPaid, req.Count, req.Sort)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "reports")
	}

	return &pb.ListResponse{Reports: dbReportsToPbReports(reports)}, nil
}

func (i *ReportServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Infof("Deleting report %s", req.ReportId)

	reportId, err := uuid.FromString(req.ReportId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of report uuid. Error %s", err.Error())
	}

	err = i.reportRepo.Delete(reportId, nil)
	if err != nil {
		log.Error(err)

		return nil, grpc.SqlErrorToGrpc(err, "report")
	}

	route := i.baseRoutingKey.SetAction("delete").SetObject("invoice").MustBuild()

	err = i.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			req, route, err.Error())
	}

	return &pb.DeleteResponse{}, nil
}

func generateReportPDF(data any, templatePath, outputPath string) error {
	r := pdf.NewInvoicePDF("")

	err := r.ParseTemplate(templatePath, data)
	if err != nil {
		log.Errorf("failed to parse PDF template: %v", err)

		return fmt.Errorf("failed to parse PDF template: %w", err)

	}

	err = r.GeneratePDF(outputPath)
	if err != nil {
		log.Errorf("failed to generate PDF invoice: %v", err)

		return fmt.Errorf("failed to generate PDF invoice: %w", err)
	}

	log.Info("PDF invoice generated successfully")

	return nil
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
