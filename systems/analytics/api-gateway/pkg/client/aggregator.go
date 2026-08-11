/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/analytics/aggregator/pb/gen"
)

type Aggregator struct {
	conn    *grpc.ClientConn
	client  pb.AggregatorServiceClient
	timeout time.Duration
	host    string
}

func NewAggregator(host string, timeout time.Duration) *Aggregator {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Aggregator server: %v", err)
	}

	return &Aggregator{
		conn:    conn,
		client:  pb.NewAggregatorServiceClient(conn),
		timeout: timeout,
		host:    host,
	}
}

func (a *Aggregator) Close() {
	if err := a.conn.Close(); err != nil {
		log.Warnf("Failed to close Aggregator connection: %v", err)
	}
}

func (a *Aggregator) Query(req *pb.QueryRequest) (*pb.QueryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.Query(ctx, req)
}

func (a *Aggregator) ListKpis() (*pb.ListKpisResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.ListKpis(ctx, &pb.ListKpisRequest{})
}

func (a *Aggregator) GetKpis(keys []string, span, op string, scope map[string]string, groupBy []string) (*pb.GetKpisResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.GetKpis(ctx, &pb.GetKpisRequest{
		Keys:    keys,
		Span:    span,
		Op:      op,
		Scope:   scope,
		GroupBy: groupBy,
	})
}

func (a *Aggregator) GetKpiTimeSeries(keys []string, span, op, from, to string, scope map[string]string, groupBy []string) (*pb.GetKpiTimeSeriesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.GetKpiTimeSeries(ctx, &pb.GetKpiTimeSeriesRequest{
		Keys:    keys,
		Span:    span,
		Op:      op,
		From:    from,
		To:      to,
		Scope:   scope,
		GroupBy: groupBy,
	})
}

func (a *Aggregator) ListReports() (*pb.ListReportsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.ListReports(ctx, &pb.ListReportsRequest{})
}

func (a *Aggregator) GetPerformanceReport(report, span string, scope map[string]string, top int32) (*pb.GetPerformanceReportResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.GetPerformanceReport(ctx, &pb.GetPerformanceReportRequest{
		Report: report,
		Span:   span,
		Scope:  scope,
		Top:    top,
	})
}

func (a *Aggregator) GetKpiBreakdown(key, span, op, by string, top int32, scope map[string]string) (*pb.GetKpiBreakdownResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.GetKpiBreakdown(ctx, &pb.GetKpiBreakdownRequest{
		Key:   key,
		Span:  span,
		Op:    op,
		By:    by,
		Top:   top,
		Scope: scope,
	})
}
