/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	apb "github.com/ukama/ukama/systems/hub/artifactmanager/pb/gen"
)

type ArtifactManager struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  apb.ArtifactServiceClient
	host    string
}

func NewArtifactManager(host string, maxMsgSize int, timeout time.Duration) *ArtifactManager {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(maxMsgSize),
		grpc.MaxCallSendMsgSize(maxMsgSize)))
	if err != nil {
		log.Fatalf("Failed to connect to ArtifactManager Service: %v", err)
	}
	client := apb.NewArtifactServiceClient(conn)

	return &ArtifactManager{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewArtifactManagerFromClient(c apb.ArtifactServiceClient) *ArtifactManager {
	return &ArtifactManager{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  c,
	}
}

func (a *ArtifactManager) Close() {
	if a.conn != nil {
		if err := a.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close ArtifactManager Service connection: %v", err)
		}
	}
}

func (a *ArtifactManager) StoreArtifact(in *apb.StoreArtifactRequest) (*apb.StoreArtifactResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.StoreArtifact(ctx, in)
}

func (a *ArtifactManager) GetArtifactLocation(in *apb.GetArtifactLocationRequest) (*apb.GetArtifactLocationResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.GetArtifactLocation(ctx, in)
}

func (a *ArtifactManager) GetArtifact(in *apb.GetArtifactRequest) (*apb.GetArtifactResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.GetArtifact(ctx, in)
}

func (a *ArtifactManager) GetArtifactVersionList(in *apb.GetArtifactVersionListRequest) (*apb.GetArtifactVersionListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.GetArtifactVersionList(ctx, in)
}

func (a *ArtifactManager) ListArtifacts(in *apb.ListArtifactRequest) (*apb.ListArtifactResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	return a.client.ListArtifacts(ctx, in)
}
