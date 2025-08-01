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
		err := s.conn.Close()
		if err != nil {
			log.Warnf("Failed to gracefully close connection to Software Service: %v", err)
		}
	}
}

func (s *SoftwareManager) UpdateSoftware(space string, name string, tag string,
	nodeId string) (*pb.UpdateSoftwareResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.client.UpdateSoftware(ctx, &pb.UpdateSoftwareRequest{
		NodeId: nodeId,
		Space:  space,
		Name:   name,
		Tag:    tag,
	})
}
