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

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/metrics/exporter/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Exporter struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.ExporterServiceClient
	host    string
}

func NewExporter(host string, timeout time.Duration) *Exporter {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := pb.NewExporterServiceClient(conn)

	return &Exporter{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewExporterFromClient(c pb.ExporterServiceClient) *Exporter {
	return &Exporter{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  c,
	}
}

func (r *Exporter) Close() {
	err := r.conn.Close()
	if err != nil {
		log.Warnf("failed to properly close exporter client. Error: %v", err)
	}
}

func (e *Exporter) Dummy(req *pb.DummyParameter) (*pb.DummyParameter, error) {
	_, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	return &pb.DummyParameter{}, nil
}
