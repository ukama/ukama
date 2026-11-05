/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
)

type fakeMonitorClient struct {
	pb.OperationMonitorServiceClient
	err      error
	deadline bool
	register *pb.RegisterIntentRequest
	cancel   *pb.CancelIntentRequest
}

func (f *fakeMonitorClient) RegisterIntent(ctx context.Context,
	req *pb.RegisterIntentRequest, _ ...grpc.CallOption) (*pb.RegisterIntentResponse,
	error) {
	_, f.deadline = ctx.Deadline()
	f.register = req
	if f.err != nil {
		return nil, f.err
	}
	return &pb.RegisterIntentResponse{}, nil
}

func (f *fakeMonitorClient) CancelIntent(ctx context.Context,
	req *pb.CancelIntentRequest, _ ...grpc.CallOption) (*pb.CancelIntentResponse,
	error) {
	_, f.deadline = ctx.Deadline()
	f.cancel = req
	if f.err != nil {
		return nil, f.err
	}
	return &pb.CancelIntentResponse{}, nil
}

func TestOperationMonitor(t *testing.T) {
	t.Run("RegisterForwardsWithDeadline", func(t *testing.T) {
		fake := &fakeMonitorClient{}
		m := &operationMonitor{client: fake, timeout: time.Second}

		_, err := m.Register(&pb.RegisterIntentRequest{OperationId: "op"})

		assert.NoError(t, err)
		assert.Equal(t, "op", fake.register.OperationId)
		assert.True(t, fake.deadline)
	})

	t.Run("CancelForwardsWithDeadline", func(t *testing.T) {
		fake := &fakeMonitorClient{}
		m := &operationMonitor{client: fake, timeout: time.Second}

		_, err := m.Cancel(&pb.CancelIntentRequest{OperationId: "op"})

		assert.NoError(t, err)
		assert.Equal(t, "op", fake.cancel.OperationId)
		assert.True(t, fake.deadline)
	})

	t.Run("CancelReturnsError", func(t *testing.T) {
		fake := &fakeMonitorClient{err: status.Error(codes.NotFound, "no intent")}
		m := &operationMonitor{client: fake, timeout: time.Second}

		_, err := m.Cancel(&pb.CancelIntentRequest{OperationId: "op"})

		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("CloseWithoutConnection", func(t *testing.T) {
		m := &operationMonitor{}

		assert.NotPanics(t, m.Close)
	})
}
