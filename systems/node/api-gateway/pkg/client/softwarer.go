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
	ukamapb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"
)

type SoftwareManager struct {
	conn    *grpc.ClientConn
	client  pb.SoftwareServiceClient
	timeout time.Duration
	host    string
}

func NewSoftwareManager(softwareManagerHost string, timeout time.Duration) *SoftwareManager {
	conn, err := grpc.NewClient(softwareManagerHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Software Service: %v", err)
	}
	client := pb.NewSoftwareServiceClient(conn)

	return &SoftwareManager{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    softwareManagerHost,
	}
}

func NewSoftwareManagerFromClient(mClient pb.SoftwareServiceClient) *SoftwareManager {
	return &SoftwareManager{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  mClient,
	}
}

func (s *SoftwareManager) Close() {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close connection to Software Service: %v", err)
		}
	}
}

func (s *SoftwareManager) UpdateSoftware(nodeId string, name string, tag string) (*pb.UpdateSoftwareResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.client.UpdateSoftware(ctx, &pb.UpdateSoftwareRequest{
		NodeId: nodeId,
		Name:   name,
		Tag:    tag,
	})
}

func (s *SoftwareManager) ListApps() (*pb.GetAppListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.client.GetAppList(ctx, &pb.GetAppListRequest{})
}

func (s *SoftwareManager) ListSoftware(nodeId string, status string, appName string) (*pb.GetSoftwareListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.client.GetSoftwareList(ctx, &pb.GetSoftwareListRequest{
		NodeId: nodeId, Status: ukamapb.SoftwareStatus(ukamapb.SoftwareStatus_value[status]), AppName: appName})
}

