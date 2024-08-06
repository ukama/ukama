/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/billing/invoice/pb/gen"
)

type Invoice interface {
	Add(rawInvoice string) (*pb.AddResponse, error)
	Get(invoiceId string, asPDF bool) (*pb.GetResponse, error)
	List(invoiceeId, invoiceeType, networkId string,
		isPaid bool, count uint32, sort bool) (*pb.ListResponse, error)
	Remove(invoiceId string) error
}

type invoice struct {
	conn        *grpc.ClientConn
	client      pb.InvoiceServiceClient
	timeout     time.Duration
	invoiceHost string
}

func NewInvoiceClient(invoiceHost string, timeout time.Duration) *invoice {
	// using same context for three connections

	conn, err := grpc.NewClient(invoiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to invoice host %q.Error: %v", invoiceHost, err)
	}

	return &invoice{
		conn:        conn,
		client:      pb.NewInvoiceServiceClient(conn),
		timeout:     timeout,
		invoiceHost: invoiceHost,
	}
}

func NewInvoiceFromClient(invoiceClient pb.InvoiceServiceClient) *invoice {
	return &invoice{
		invoiceHost: "localhost",
		timeout:     1 * time.Second,
		conn:        nil,
		client:      invoiceClient,
	}
}

func (i *invoice) Close() {
	i.conn.Close()
}

func (i *invoice) Add(rawInvoice string) (*pb.AddResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), i.timeout)
	defer cancel()

	res, err := i.client.Add(ctx, &pb.AddRequest{
		RawInvoice: rawInvoice})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (i *invoice) Get(invoiceId string, AsPDF bool) (*pb.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), i.timeout)
	defer cancel()

	res, err := i.client.Get(ctx, &pb.GetRequest{
		InvoiceId: invoiceId,
		AsPdf:     AsPDF})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (i *invoice) List(invoiceeId, invoiceeType, networkId string, isPaid bool, count uint32, sort bool) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), i.timeout)
	defer cancel()

	return i.client.List(ctx,
		&pb.ListRequest{
			InvoiceeId:   invoiceeId,
			InvoiceeType: invoiceeType,
			NetworkId:    networkId,
			IsPaid:       isPaid,
			Count:        count,
			Sort:         sort,
		})
}

func (i *invoice) Remove(invoiceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), i.timeout)
	defer cancel()

	_, err := i.client.Delete(ctx, &pb.DeleteRequest{InvoiceId: invoiceId})

	return err
}
