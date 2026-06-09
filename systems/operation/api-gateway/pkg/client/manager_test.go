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

	"github.com/tj/assert"
	"google.golang.org/grpc"

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

// fakeGrpc implements pb.OperationManagerServiceClient for wrapper tests.
type fakeGrpc struct {
	pb.OperationManagerServiceClient
}

func (f *fakeGrpc) StartOperation(_ context.Context, in *pb.StartOperationRequest, _ ...grpc.CallOption) (*pb.StartOperationResponse, error) {
	return &pb.StartOperationResponse{Operation: &pb.Operation{ResourceKey: in.ResourceKey}}, nil
}
func (f *fakeGrpc) GetOperation(_ context.Context, in *pb.GetOperationRequest, _ ...grpc.CallOption) (*pb.GetOperationResponse, error) {
	return &pb.GetOperationResponse{Operation: &pb.Operation{Id: in.Id}}, nil
}
func (f *fakeGrpc) GetByResource(_ context.Context, in *pb.GetByResourceRequest, _ ...grpc.CallOption) (*pb.GetByResourceResponse, error) {
	return &pb.GetByResourceResponse{Operation: &pb.Operation{ResourceKey: in.ResourceKey}}, nil
}
func (f *fakeGrpc) MarkRunning(_ context.Context, in *pb.MarkRunningRequest, _ ...grpc.CallOption) (*pb.MarkRunningResponse, error) {
	return &pb.MarkRunningResponse{Operation: &pb.Operation{Id: in.Id, FencingToken: in.FencingToken}}, nil
}
func (f *fakeGrpc) ForceUnlock(_ context.Context, in *pb.ForceUnlockRequest, _ ...grpc.CallOption) (*pb.ForceUnlockResponse, error) {
	return &pb.ForceUnlockResponse{Operation: &pb.Operation{Id: in.Id}}, nil
}

func newTestManager() *Manager {
	return NewManagerFromClient(&fakeGrpc{})
}

func TestManager_Start(t *testing.T) {
	resp, err := newTestManager().Start(&pb.StartOperationRequest{ResourceKey: "node:abc"})
	assert.NoError(t, err)
	assert.Equal(t, "node:abc", resp.Operation.ResourceKey)
}

func TestManager_Get(t *testing.T) {
	resp, err := newTestManager().Get("op-1")
	assert.NoError(t, err)
	assert.Equal(t, "op-1", resp.Operation.Id)
}

func TestManager_GetByResource(t *testing.T) {
	resp, err := newTestManager().GetByResource("node:abc")
	assert.NoError(t, err)
	assert.Equal(t, "node:abc", resp.Operation.ResourceKey)
}

func TestManager_MarkRunning(t *testing.T) {
	resp, err := newTestManager().MarkRunning("op-1", 5)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5), resp.Operation.FencingToken)
}

func TestManager_ForceUnlock(t *testing.T) {
	resp, err := newTestManager().ForceUnlock("op-1", "owner", "stuck")
	assert.NoError(t, err)
	assert.Equal(t, "op-1", resp.Operation.Id)
}

func TestManager_CloseNilConn(t *testing.T) {
	// NewManagerFromClient leaves conn nil; Close must not panic.
	newTestManager().Close()
}
