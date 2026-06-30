/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package rest

import (
	"errors"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tj/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

// stubManager is a configurable manager fake for exercising every handler
// success and error branch.
type stubManager struct {
	startResp *pb.StartOperationResponse
	getResp   *pb.GetOperationResponse
	byResResp *pb.GetByResourceResponse
	markResp  *pb.MarkRunningResponse
	err       error
}

func (m *stubManager) Start(*pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
	return m.startResp, m.err
}
func (m *stubManager) Get(string) (*pb.GetOperationResponse, error) {
	return m.getResp, m.err
}
func (m *stubManager) GetByResource(string) (*pb.GetByResourceResponse, error) {
	return m.byResResp, m.err
}
func (m *stubManager) MarkRunning(string, uint64) (*pb.MarkRunningResponse, error) {
	return m.markResp, m.err
}
func (m *stubManager) ForceUnlock(id, _, _ string) (*pb.ForceUnlockResponse, error) {
	return &pb.ForceUnlockResponse{Operation: &pb.Operation{Id: id}}, m.err
}

func routerWith(m manager) *Router {
	return &Router{clients: &Clients{Manager: m, Member: &fakeMember{role: "ROLE_OWNER"}}}
}

func TestPostStartHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := &stubManager{startResp: &pb.StartOperationResponse{
			Operation: &pb.Operation{Id: "op-1", ResourceKey: "node:abc"},
		}}
		r := routerWith(m)

		resp, err := r.postStartHandler(&gin.Context{}, &StartOperationRequest{
			Type: "RestartNode", System: "node", ResourceKey: "node:abc",
		})

		assert.NoError(t, err)
		assert.Equal(t, "op-1", resp.Operation.Id)
	})

	t.Run("Conflict", func(t *testing.T) {
		m := &stubManager{startResp: &pb.StartOperationResponse{
			ConflictingOperation: &pb.Operation{Id: "holder", ResourceKey: "node:abc"},
		}}
		r := routerWith(m)

		resp, err := r.postStartHandler(&gin.Context{}, &StartOperationRequest{ResourceKey: "node:abc"})

		assert.NoError(t, err)
		assert.Nil(t, resp.Operation)
		assert.Equal(t, "holder", resp.ConflictingOperation.Id)
	})

	t.Run("Error", func(t *testing.T) {
		r := routerWith(&stubManager{err: errors.New("grpc down")})

		_, err := r.postStartHandler(&gin.Context{}, &StartOperationRequest{ResourceKey: "node:abc"})

		assert.Error(t, err)
	})
}

func TestGetOperationHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := &stubManager{getResp: &pb.GetOperationResponse{
			Operation: &pb.Operation{Id: "op-1", Status: pb.OperationStatus_RUNNING},
		}}
		r := routerWith(m)

		resp, err := r.getOperationHandler(&gin.Context{}, &GetOperationRequest{Id: "op-1"})

		assert.NoError(t, err)
		assert.Equal(t, "op-1", resp.Operation.Id)
	})

	t.Run("Error", func(t *testing.T) {
		r := routerWith(&stubManager{err: errors.New("not found")})

		_, err := r.getOperationHandler(&gin.Context{}, &GetOperationRequest{Id: "x"})

		assert.Error(t, err)
	})
}

func TestGetByResourceHandler(t *testing.T) {
	t.Run("Locked", func(t *testing.T) {
		m := &stubManager{byResResp: &pb.GetByResourceResponse{
			Operation: &pb.Operation{Id: "op-1", ResourceKey: "node:abc"},
		}}
		r := routerWith(m)

		resp, err := r.getByResourceHandler(&gin.Context{}, &GetByResourceRequest{ResourceKey: "node:abc"})

		assert.NoError(t, err)
		assert.True(t, resp.Locked)
		assert.Equal(t, "op-1", resp.Operation.Id)
	})

	t.Run("Free", func(t *testing.T) {
		r := routerWith(&stubManager{byResResp: &pb.GetByResourceResponse{}})

		resp, err := r.getByResourceHandler(&gin.Context{}, &GetByResourceRequest{ResourceKey: "node:free"})

		assert.NoError(t, err)
		assert.False(t, resp.Locked)
		assert.Nil(t, resp.Operation)
	})

	t.Run("Error", func(t *testing.T) {
		r := routerWith(&stubManager{err: errors.New("grpc down")})

		_, err := r.getByResourceHandler(&gin.Context{}, &GetByResourceRequest{ResourceKey: "node:abc"})

		assert.Error(t, err)
	})
}

func TestPostMarkRunningHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		m := &stubManager{markResp: &pb.MarkRunningResponse{
			Operation: &pb.Operation{Id: "op-1", Status: pb.OperationStatus_RUNNING},
		}}
		r := routerWith(m)

		resp, err := r.postMarkRunningHandler(&gin.Context{}, &MarkRunningRequest{Id: "op-1", FencingToken: 3})

		assert.NoError(t, err)
		assert.Equal(t, "op-1", resp.Operation.Id)
	})

	t.Run("Error", func(t *testing.T) {
		r := routerWith(&stubManager{err: errors.New("token mismatch")})

		_, err := r.postMarkRunningHandler(&gin.Context{}, &MarkRunningRequest{Id: "op-1", FencingToken: 3})

		assert.Error(t, err)
	})
}

func TestOperationFromProto(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, operationFromProto(nil))
	})

	t.Run("FullyPopulated", func(t *testing.T) {
		now := timestamppb.New(time.Now().UTC())
		op := &pb.Operation{
			Id:             "op-1",
			Type:           "RestartNode",
			System:         "node",
			Status:         pb.OperationStatus_RUNNING,
			FencingToken:   5,
			RequestedBy:    "user-1",
			IdempotencyKey: "key-1",
			ResourceKey:    "node:abc",
			LeaseExpiresAt: now,
			Error:          "",
			StartedAt:      now,
			TerminalAt:     now,
			CreatedAt:      now,
		}

		out := operationFromProto(op)

		assert.Equal(t, "op-1", out.Id)
		assert.Equal(t, pb.OperationStatus_RUNNING.String(), out.Status)
		assert.Equal(t, uint64(5), out.FencingToken)
		assert.NotNil(t, out.StartedAt)
		assert.NotNil(t, out.TerminalAt)
		assert.NotNil(t, out.LeaseExpiresAt)
	})

	t.Run("NilTimestampsBecomeNil", func(t *testing.T) {
		op := &pb.Operation{Id: "op-2", Status: pb.OperationStatus_PENDING}

		out := operationFromProto(op)

		assert.Nil(t, out.StartedAt)
		assert.Nil(t, out.TerminalAt)
	})
}
