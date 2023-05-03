package client

import (
	"context"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"

	"google.golang.org/grpc"
)

type Billing struct {
	conn          *grpc.ClientConn
	invoiceClient pb.InvoiceServiceClient
	pdfClient     PdfClient
	timeout       time.Duration
	invoiceHost   string
	fileHost      string
}

func NewBilling(invoiceHost, fileHost string, timeout time.Duration) *Billing {
	// using same context for three connections
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, invoiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	invoiceClient := pb.NewInvoiceServiceClient(conn)

	pdfClient, err := NewPdfClient(fileHost, false)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &Billing{
		conn:          conn,
		invoiceClient: invoiceClient,
		pdfClient:     pdfClient,
		timeout:       timeout,
		invoiceHost:   invoiceHost,
		fileHost:      fileHost,
	}
}

func NewBillingFromClient(invoiceClient pb.InvoiceServiceClient) *Billing {
	return &Billing{
		invoiceHost:   "localhost",
		timeout:       1 * time.Second,
		conn:          nil,
		invoiceClient: invoiceClient,
	}
}

func (r *Billing) Close() {
	r.conn.Close()
}

func (r *Billing) AddInvoice(subscriberId string, rawInvoice string) (*pb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.invoiceClient.Add(ctx, &pb.AddRequest{
		SubscriberId: subscriberId,
		RawInvoice:   rawInvoice})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b *Billing) GetInvoice(invoiceId string, AsPDF bool) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	res, err := b.invoiceClient.Get(ctx, &pb.GetRequest{
		InvoiceId: invoiceId,
		AsPdf:     AsPDF})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Billing) GetInvoices(subscriberId string, AsPDF bool) (*pb.GetBySubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	res, err := r.invoiceClient.GetBySubscriber(ctx,
		&pb.GetBySubscriberRequest{SubscriberId: subscriberId,
			AsPdf: AsPDF})

	if err != nil {
		return nil, err
	}

	if res.Invoices == nil {
		return &pb.GetBySubscriberResponse{Invoices: []*pb.Invoice{}, SubscriberId: subscriberId}, nil
	}

	return res, nil
}

func (r *Billing) RemoveInvoice(invoiceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.invoiceClient.Delete(ctx, &pb.DeleteRequest{InvoiceId: invoiceId})

	return err
}

func (b *Billing) GetInvoicePDF(invoiceId string) ([]byte, error) {
	res, err := b.pdfClient.GetPdf(invoiceId)

	if err != nil {
		return nil, err
	}

	return res, nil
}
