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

	"github.com/ukama/ukama/systems/billing/invoice/pkg"
	"github.com/ukama/ukama/systems/billing/invoice/pkg/db"
	"github.com/ukama/ukama/systems/billing/invoice/pkg/pdf"
	"github.com/ukama/ukama/systems/billing/invoice/pkg/util"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

const defaultTemplate = "templates/invoice.html.tmpl"
const pdfFolder = "/srv/static/"

type InvoiceServer struct {
	OrgName          string
	OrgId            uuid.UUID
	invoiceRepo      db.InvoiceRepo
	subscriberClient csub.SubscriberClient
	msgbus           mb.MsgBusServiceClient
	baseRoutingKey   msgbus.RoutingKeyBuilder
	pb.UnimplementedInvoiceServiceServer
}

func NewInvoiceServer(orgName, org string, invoiceRepo db.InvoiceRepo, subscriberClient csub.SubscriberClient, msgBus mb.MsgBusServiceClient) *InvoiceServer {
	orgId, err := uuid.FromString(org)
	if err != nil {
		panic(fmt.Sprintf("invalid format of org uuid: %s", org))
	}

	return &InvoiceServer{
		OrgName:          orgName,
		OrgId:            orgId,
		invoiceRepo:      invoiceRepo,
		subscriberClient: subscriberClient,
		msgbus:           msgBus,
		baseRoutingKey:   msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
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

	invoiceeId, err := uuid.FromString(rwInvoceStruct.Customer.ExternalID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invoicee uuid. Error %s", err.Error())
	}

	invoice := &db.Invoice{
		InvoiceeId:   invoiceeId,
		InvoiceeType: db.InvoiceeTypeOrg,
	}

	if invoiceeId != i.OrgId {
		subscriberInfo, err := i.subscriberClient.Get(invoiceeId.String())
		if err != nil {
			return nil, err
		}

		invoice.NetworkId = subscriberInfo.NetworkId
		invoice.InvoiceeType = db.InvoiceeTypeSubscriber
	}

	log.Infof("Adding invoice for invoicee: %s", invoiceeId)
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

func (i *InvoiceServer) GetByInvoicee(ctx context.Context, req *pb.GetByInvoiceeRequest) (*pb.GetByInvoiceeResponse, error) {
	invoiceeId, err := uuid.FromString(req.InvoiceeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invoicee uuid. Error %s", err.Error())
	}

	invoices, err := i.invoiceRepo.GetByInvoicee(invoiceeId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invoices")
	}

	resp := &pb.GetByInvoiceeResponse{
		InvoiceeId: req.InvoiceeId,
		Invoices:   dbInvoicesToPbInvoices(invoices),
	}

	return resp, nil
}

func (i *InvoiceServer) GetByNetwork(ctx context.Context, req *pb.GetByNetworkRequest) (*pb.GetByNetworkResponse, error) {
	networkId, err := uuid.FromString(req.NetworkId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of network uuid. Error %s", err.Error())
	}

	invoices, err := i.invoiceRepo.GetByNetwork(networkId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invoices")
	}

	resp := &pb.GetByNetworkResponse{
		NetworkId: req.NetworkId,
		Invoices:  dbInvoicesToPbInvoices(invoices),
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
		InvoiceeId:   invoice.InvoiceeId.String(),
		InvoiceeType: invoice.InvoiceeType.String(),
		Period:       invoice.Period.String(),
		IsPaid:       invoice.IsPaid,
		CreatedAt:    invoice.CreatedAt.String(),
	}

	if invoice.NetworkId != uuid.Nil {
		inv.NetworkId = invoice.NetworkId.String()
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
