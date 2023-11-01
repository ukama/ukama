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

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
	"github.com/ukama/ukama/systems/billing/invoice/pkg"

	"github.com/ukama/ukama/systems/billing/invoice/pkg/db"
	"github.com/ukama/ukama/systems/billing/invoice/pkg/pdf"
	"github.com/ukama/ukama/systems/billing/invoice/pkg/util"
	"github.com/ukama/ukama/systems/common/grpc"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const defaultTemplate = "templates/invoice.html.tmpl"
const pdfFolder = "/srv/static/"

type InvoiceServer struct {
	invoiceRepo    db.InvoiceRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedInvoiceServiceServer
}

func NewInvoiceServer(orgName string, invoiceRepo db.InvoiceRepo, msgBus mb.MsgBusServiceClient) *InvoiceServer {
	return &InvoiceServer{
		invoiceRepo:    invoiceRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (i *InvoiceServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Infof("Unmarshalling raw invoice from webhook: %v", req.RawInvoice)
	rwInvoceStruct := &util.RawInvoice{}

	err := json.Unmarshal([]byte(req.RawInvoice), rwInvoceStruct)
	if err != nil {
		log.Errorf("Failed to unmarshal RawInvoice JSON to RawInvoice struct %v", err)

		return nil, status.Errorf(codes.InvalidArgument,
			"failed to unmarshal RawInvoice JSON paylod from webhook. Error %s", err)
	}

	if rwInvoceStruct == nil || rwInvoceStruct.Customer == nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid JSON format of RawInvoice. Error %s", req.RawInvoice)
	}

	subscriberId, err := uuid.FromString(rwInvoceStruct.Customer.ExternalID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	invoice := &db.Invoice{
		SubscriberId: subscriberId,
	}

	log.Infof("Adding invoice for subscriber: %s", subscriberId)
	err = i.invoiceRepo.Add(invoice, func(*db.Invoice, *gorm.DB) error {
		invoice.Id = uuid.NewV4()
		rwInvoceStruct.FileURL = fmt.Sprintf("http://{API_ENDPOINT}/pdf/%s.pdf", invoice.Id.String())

		rwInvoiceBytes, err := json.Marshal(rwInvoceStruct)
		if err != nil {
			log.Errorf("Failed to marshal RawInvoice struct to RawInvoice JSON %v", err)

			return fmt.Errorf("failed to marshal RawInvoice struct to RawInvoice JSON: %w", err)
		}

		invoice.Period = time.Now().UTC()
		invoice.RawInvoice = datatypes.JSON(rwInvoiceBytes)

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invoice")
	}

	log.Infof("starting PDF generation")
	go func() {
		err = generateInvoicePDF(rwInvoceStruct, defaultTemplate, filepath.Join(pdfFolder, invoice.Id.String()+".pdf"))
		if err != nil {
			log.Errorf("PDF generation failure: failed to generate invoice PDF: %v", err)
		}

		log.Infof("finishing PDF generation")
	}()

	resp := &pb.AddResponse{
		Invoice: dbInvoiceToPbInvoice(invoice),
	}

	route := i.baseRoutingKey.SetAction("generate").SetObject("invoice").MustBuild()

	err = i.msgbus.PublishRequest(route, resp.Invoice)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			req, route, err.Error())
	}

	return resp, nil
}

func (i *InvoiceServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	invoiceId, err := uuid.FromString(req.InvoiceId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invoice uuid. Error %s", err.Error())
	}

	invoice, err := i.invoiceRepo.Get(invoiceId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invoice")
	}

	return &pb.GetResponse{
		Invoice: dbInvoiceToPbInvoice(invoice),
	}, nil
}

func (i *InvoiceServer) GetBySubscriber(ctx context.Context, req *pb.GetBySubscriberRequest) (*pb.GetBySubscriberResponse, error) {
	subscriberId, err := uuid.FromString(req.SubscriberId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	invoices, err := i.invoiceRepo.GetBySubscriber(subscriberId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invoices")
	}

	resp := &pb.GetBySubscriberResponse{
		SubscriberId: req.SubscriberId,
		Invoices:     dbInvoicesToPbInvoices(invoices),
	}

	return resp, nil
}

func (i *InvoiceServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Infof("Deleting invoice %s", req.InvoiceId)

	invoiceId, err := uuid.FromString(req.InvoiceId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invoice uuid. Error %s", err.Error())
	}

	err = i.invoiceRepo.Delete(invoiceId, nil)
	if err != nil {
		log.Error(err)

		return nil, grpc.SqlErrorToGrpc(err, "invoice")
	}

	route := i.baseRoutingKey.SetAction("delete").SetObject("invoice").MustBuild()

	err = i.msgbus.PublishRequest(route, req)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s",
			req, route, err.Error())
	}

	return &pb.DeleteResponse{}, nil
}

func generateInvoicePDF(data any, templatePath, outputPath string) error {
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

func dbInvoiceToPbInvoice(invoice *db.Invoice) *pb.Invoice {
	inv := &pb.Invoice{
		Id:           invoice.Id.String(),
		SubscriberId: invoice.SubscriberId.String(),
		Period:       timestamppb.New(invoice.Period),
		IsPaid:       invoice.IsPaid,
		CreatedAt:    timestamppb.New(invoice.CreatedAt),
	}

	val := &pb.RawInvoice{}

	m := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	err := m.Unmarshal([]byte(invoice.RawInvoice.String()), val)
	//TODO: error handling

	if err == nil {
		inv.RawInvoice = val
	}

	return inv
}

func dbInvoicesToPbInvoices(invoices []db.Invoice) []*pb.Invoice {
	res := []*pb.Invoice{}

	for _, n := range invoices {
		res = append(res, dbInvoiceToPbInvoice(&n))
	}

	return res
}
