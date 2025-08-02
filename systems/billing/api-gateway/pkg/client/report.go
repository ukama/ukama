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
	pb "github.com/ukama/ukama/systems/billing/report/pb/gen"
)

type Report interface {
	Add(rawReport string) (*pb.ReportResponse, error)
	Get(reportId string) (*pb.ReportResponse, error)
	List(ownerId, ownerType, networkId, reportType string,
		isPaid bool, count uint32, sort bool) (*pb.ListResponse, error)

	Update(reportId string, isPaid bool) (*pb.ReportResponse, error)
	Remove(reportId string) error
}

type report struct {
	conn       *grpc.ClientConn
	client     pb.ReportServiceClient
	timeout    time.Duration
	reportHost string
}

func NewReportClient(reportHost string, timeout time.Duration) *report {
	// using same context for three connections
	conn, err := grpc.NewClient(reportHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to report host %q.Error: %v", reportHost, err)
	}

	return &report{
		conn:       conn,
		client:     pb.NewReportServiceClient(conn),
		timeout:    timeout,
		reportHost: reportHost,
	}
}

func NewReportFromClient(reportClient pb.ReportServiceClient) *report {
	return &report{
		reportHost: "localhost",
		timeout:    1 * time.Second,
		conn:       nil,
		client:     reportClient,
	}
}

func (r *report) Close() {
	if r.conn != nil {
		err := r.conn.Close()
		if err != nil {
			log.Warnf("Failed to gracefully close connection to Report Service: %v", err)
		}
	}
}

func (r *report) Add(rawReport string) (*pb.ReportResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Add(ctx, &pb.AddRequest{
		RawReport: rawReport})

}

func (r *report) Get(reportId string) (*pb.ReportResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Get(ctx, &pb.GetRequest{
		ReportId: reportId})
}

func (r *report) List(ownerId, ownerType, networkId, reportType string,
	isPaid bool, count uint32, sort bool) (*pb.ListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.List(ctx,
		&pb.ListRequest{
			OwnerId:    ownerId,
			OwnerType:  ownerType,
			NetworkId:  networkId,
			ReportType: reportType,
			IsPaid:     isPaid,
			Count:      count,
			Sort:       sort,
		})
}

func (r *report) Update(reportId string, isPaid bool) (*pb.ReportResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Update(ctx, &pb.UpdateRequest{
		ReportId: reportId,
		IsPaid:   isPaid,
	})
}

func (r *report) Remove(reportId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.client.Delete(ctx, &pb.DeleteRequest{ReportId: reportId})

	return err
}
