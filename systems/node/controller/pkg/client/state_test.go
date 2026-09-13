/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package client

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type configurationStateServer struct {
	pb.UnimplementedStateServiceServer
}

func (configurationStateServer) RecordConfiguration(ctx context.Context, req *pb.RecordConfigurationRequest) (*pb.ConfigurationStatus, error) {
	if req.RequestId == "reject" {
		return nil, status.Error(codes.FailedPrecondition, "cancelled request")
	}
	if req.RequestId == "wait" {
		<-ctx.Done()
		return nil, status.FromContextError(ctx.Err()).Err()
	}
	return &pb.ConfigurationStatus{RequestId: req.RequestId, Cancelled: req.Cancelled}, nil
}

func TestConfigurationStateRPC(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	server := grpc.NewServer()
	pb.RegisterStateServiceServer(server, configurationStateServer{})
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	t.Cleanup(func() { server.Stop(); require.NoError(t, <-done) })
	client, err := NewConfigurationState(listener.Addr().String(), time.Second)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	result, err := client.RecordConfiguration(context.Background(), &pb.RecordConfigurationRequest{RequestId: "attempt-1", Cancelled: true})
	require.NoError(t, err)
	require.Equal(t, "attempt-1", result.RequestId)
	require.True(t, result.Cancelled)
	_, err = client.RecordConfiguration(context.Background(), &pb.RecordConfigurationRequest{RequestId: "reject"})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = client.RecordConfiguration(ctx, &pb.RecordConfigurationRequest{RequestId: "attempt-1"})
	require.Equal(t, codes.Canceled, status.Code(err))
	client.timeout = 10 * time.Millisecond
	_, err = client.RecordConfiguration(context.Background(), &pb.RecordConfigurationRequest{RequestId: "wait"})
	require.Equal(t, codes.DeadlineExceeded, status.Code(err))
}
