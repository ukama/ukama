/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package client

import (
	"context"
	"time"

	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConfigurationState struct {
	conn    *grpc.ClientConn
	client  pb.StateServiceClient
	timeout time.Duration
}

func NewConfigurationState(host string, timeout time.Duration) (*ConfigurationState, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ConfigurationState{conn: conn, client: pb.NewStateServiceClient(conn), timeout: timeout}, nil
}

func (s *ConfigurationState) Close() error { return s.conn.Close() }

func (s *ConfigurationState) RecordConfiguration(ctx context.Context, req *pb.RecordConfigurationRequest) (*pb.ConfigurationStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.client.RecordConfiguration(ctx, req)
}
