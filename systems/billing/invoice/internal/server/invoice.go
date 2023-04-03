package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/billing/invoice/internal/db"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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

	logrus.Infof("Adding invoice for subscriber: %s", subscriberId)
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

	nt, err := i.invoiceRepo.Get(invoiceId)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "invoice")
	}

	return &pb.GetResponse{
		Invoice: dbInvoiceToPbInvoice(nt),
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
	logrus.Infof("Deleting invoice %s", req.InvoiceId)

	invoiceId, err := uuid.FromString(req.InvoiceId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid format of invoice uuid. Error %s", err.Error())
	}

	err = i.invoiceRepo.Delete(invoiceId, nil)
	if err != nil {
		logrus.Error(err)

		return nil, grpc.SqlErrorToGrpc(err, "invoice")
	}

	return &pb.DeleteResponse{}, nil
}

func dbInvoiceToPbInvoice(invoice *db.Invoice) *pb.Invoice {
	inv := &pb.Invoice{
		Id:           invoice.Id.String(),
		SubscriberId: invoice.SubscriberId.String(),
		Period:       timestamppb.New(invoice.Period),
		IsPaid:       invoice.IsPaid,
		CreatedAt:    timestamppb.New(invoice.CreatedAt),
	}

	val, err := structpb.NewValue(invoice.RawInvoice.String())
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
