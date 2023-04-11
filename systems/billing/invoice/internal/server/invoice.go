package server

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"

	"github.com/ukama/ukama/systems/billing/invoice/internal/db"
	"github.com/ukama/ukama/systems/billing/invoice/internal/pdf"
	"github.com/ukama/ukama/systems/billing/invoice/internal/util"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const defaultTemplate = "templates/test.html.tmpl"

type InvoiceServer struct {
	pb.UnimplementedInvoiceServiceServer
	invoiceRepo db.InvoiceRepo
}

func NewInvoiceServer(invoiceRepo db.InvoiceRepo) *InvoiceServer {
	return &InvoiceServer{
		invoiceRepo: invoiceRepo,
	}
}

func (i *InvoiceServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	subscriberId, err := uuid.FromString(req.SubscriberId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of subscriber uuid. Error %s", err.Error())
	}

	invoice := &db.Invoice{
		SubscriberId: subscriberId,
		Period:       req.GetPeriod().AsTime(),
		RawInvoice:   datatypes.JSON([]byte(req.RawInvoice)),
	}

	log.Infof("Adding invoice for subscriber: %s", subscriberId)
	err = i.invoiceRepo.Add(invoice, func(*db.Invoice, *gorm.DB) error {
		invoice.Id = uuid.NewV4()

		return nil
	})

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invoice")
	}

	return &pb.AddResponse{
		Invoice: dbInvoiceToPbInvoice(invoice),
	}, nil
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

	log.Infof("starting PDF operation")

	rw := &util.RawInvoice{}

	err = json.Unmarshal([]byte(invoice.RawInvoice.String()), rw)
	if err != nil {
		log.Errorf("PDF operation failure: %v", err)

		err = nil

		// return nil, status.Errorf(codes.Internal,
		// "failure to unmarshal raw invoice. Error %s", err.Error())
	}

	err = generateInvoicePDF(rw, defaultTemplate, "./test.pdf")
	if err != nil {
		log.Errorf("PDF operation failure: %v", err)

		err = nil

		// return nil, status.Errorf(codes.Internal,
		// "failure to generate PDF file. Error %s", err.Error())
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

	return &pb.DeleteResponse{}, nil
}

func generateInvoicePDF(data any, templatePath, outputPath string) error {
	r := pdf.NewInvoicePDF("")

	templateData := struct {
		Title    string
		FileName string
		Body     string
	}{
		Title:    "Markdown Preview Tool",
		FileName: "templateFile",
		Body:     "hello world",
	}

	// err := r.ParseSlimTemplate(templatePath, data)
	err := r.ParseTemplate(templatePath, templateData)
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
