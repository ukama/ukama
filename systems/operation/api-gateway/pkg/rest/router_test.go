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

type fakeManager struct {
	lastStartReq         *pb.StartOperationRequest
	lastGetId            string
	lastGetByResourceKey string
	lastMarkRunningId    string
	lastMarkRunningToken uint64
	lastForceUnlockId    string
	lastForceUnlockActor string
	lastForceUnlockReason string

	startResp         *pb.StartOperationResponse
	startErr          error
	getResp           *pb.GetOperationResponse
	getErr            error
	getByResourceResp *pb.GetByResourceResponse
	getByResourceErr  error
	markRunningResp   *pb.MarkRunningResponse
	markRunningErr    error
	forceUnlockResp   *pb.ForceUnlockResponse
	forceUnlockErr    error
}

func (f *fakeManager) Start(req *pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
	f.lastStartReq = req
	if f.startErr != nil {
		return nil, f.startErr
	}
	if f.startResp != nil {
		return f.startResp, nil
	}
	return &pb.StartOperationResponse{}, nil
}

func (f *fakeManager) Get(id string) (*pb.GetOperationResponse, error) {
	f.lastGetId = id
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.getResp != nil {
		return f.getResp, nil
	}
	return &pb.GetOperationResponse{}, nil
}

func (f *fakeManager) GetByResource(resourceKey string) (*pb.GetByResourceResponse, error) {
	f.lastGetByResourceKey = resourceKey
	if f.getByResourceErr != nil {
		return nil, f.getByResourceErr
	}
	if f.getByResourceResp != nil {
		return f.getByResourceResp, nil
	}
	return &pb.GetByResourceResponse{}, nil
}

func (f *fakeManager) MarkRunning(id string, fencingToken uint64) (*pb.MarkRunningResponse, error) {
	f.lastMarkRunningId = id
	f.lastMarkRunningToken = fencingToken
	if f.markRunningErr != nil {
		return nil, f.markRunningErr
	}
	if f.markRunningResp != nil {
		return f.markRunningResp, nil
	}
	return &pb.MarkRunningResponse{}, nil
}

func (f *fakeManager) ForceUnlock(id, actor, reason string) (*pb.ForceUnlockResponse, error) {
	f.lastForceUnlockId = id
	f.lastForceUnlockActor = actor
	f.lastForceUnlockReason = reason
	if f.forceUnlockErr != nil {
		return nil, f.forceUnlockErr
	}
	if f.forceUnlockResp != nil {
		return f.forceUnlockResp, nil
	}
	return &pb.ForceUnlockResponse{Operation: &pb.Operation{Id: id}}, nil
}

func (f *fakeManager) Complete(id, actor, reason string) (*pb.ForceUnlockResponse, error) {
	return &pb.ForceUnlockResponse{Operation: &pb.Operation{Id: id}}, nil
}

func newTestRouter(mgr *fakeManager) *Router {
	return &Router{
		clients: &Clients{
			Manager: mgr,
		},
	}
}

func sampleProtoOperation() *pb.Operation {
	now := timestamppb.New(time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC))
	return &pb.Operation{
		Id:             "8e13fa4b-a8a7-40aa-8c61-2891cd16dc7f",
		Type:           "RestartNode",
		System:         "node",
		Status:         pb.OperationStatus_RUNNING,
		FencingToken:   42,
		RequestedBy:    "user-1",
		IdempotencyKey: "idem-1",
		ResourceKey:    "node:uk-sa2450-tnode-v0-4e86",
		LeaseExpiresAt: now,
		StartedAt:      now,
		CreatedAt:      now,
	}
}

func assertOperationMapped(t *testing.T, expected *pb.Operation, actual *Operation) {
	t.Helper()
	assert.NotNil(t, actual)
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.System, actual.System)
	assert.Equal(t, expected.Status.String(), actual.Status)
	assert.Equal(t, expected.FencingToken, actual.FencingToken)
	assert.Equal(t, expected.RequestedBy, actual.RequestedBy)
	assert.Equal(t, expected.IdempotencyKey, actual.IdempotencyKey)
	assert.Equal(t, expected.ResourceKey, actual.ResourceKey)
	assert.Equal(t, expected.Error, actual.Error)
	assert.NotNil(t, actual.LeaseExpiresAt)
	assert.Equal(t, expected.LeaseExpiresAt.AsTime(), *actual.LeaseExpiresAt)
	assert.NotNil(t, actual.StartedAt)
	assert.Equal(t, expected.StartedAt.AsTime(), *actual.StartedAt)
	assert.Nil(t, actual.TerminalAt)
	assert.NotNil(t, actual.CreatedAt)
	assert.Equal(t, expected.CreatedAt.AsTime(), *actual.CreatedAt)
}

func TestPostStartHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := &StartOperationRequest{
		Type:           "RestartNode",
		System:         "node",
		ResourceKey:    "node:uk-sa2450-tnode-v0-4e86",
		RequestedBy:    "user-1",
		IdempotencyKey: "idem-1",
		LeaseSeconds:   300,
	}

	t.Run("Success", func(t *testing.T) {
		op := sampleProtoOperation()
		conflict := &pb.Operation{Id: "conflict-op"}
		mgr := &fakeManager{
			startResp: &pb.StartOperationResponse{
				Operation:            op,
				ConflictingOperation: conflict,
			},
		}
		r := newTestRouter(mgr)

		resp, err := r.postStartHandler(&gin.Context{}, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, req.Type, mgr.lastStartReq.Type)
		assert.Equal(t, req.System, mgr.lastStartReq.System)
		assert.Equal(t, req.ResourceKey, mgr.lastStartReq.ResourceKey)
		assert.Equal(t, req.RequestedBy, mgr.lastStartReq.RequestedBy)
		assert.Equal(t, req.IdempotencyKey, mgr.lastStartReq.IdempotencyKey)
		assert.Equal(t, req.LeaseSeconds, mgr.lastStartReq.LeaseSeconds)
		assertOperationMapped(t, op, resp.Operation)
		assert.Equal(t, "conflict-op", resp.ConflictingOperation.Id)
	})

	t.Run("ManagerError", func(t *testing.T) {
		mgr := &fakeManager{startErr: errors.New("manager down")}
		r := newTestRouter(mgr)

		_, err := r.postStartHandler(&gin.Context{}, req)
		assert.Error(t, err)
		assert.NotNil(t, mgr.lastStartReq)
	})
}

func TestGetOperationHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := &GetOperationRequest{Id: "8e13fa4b-a8a7-40aa-8c61-2891cd16dc7f"}

	t.Run("Success", func(t *testing.T) {
		op := sampleProtoOperation()
		mgr := &fakeManager{
			getResp: &pb.GetOperationResponse{Operation: op},
		}
		r := newTestRouter(mgr)

		resp, err := r.getOperationHandler(&gin.Context{}, req)
		assert.NoError(t, err)
		assert.Equal(t, req.Id, mgr.lastGetId)
		assertOperationMapped(t, op, resp.Operation)
	})

	t.Run("ManagerError", func(t *testing.T) {
		mgr := &fakeManager{getErr: errors.New("not found")}
		r := newTestRouter(mgr)

		_, err := r.getOperationHandler(&gin.Context{}, req)
		assert.Error(t, err)
		assert.Equal(t, req.Id, mgr.lastGetId)
	})
}

func TestGetByResourceHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := &GetByResourceRequest{ResourceKey: "node:uk-sa2450-tnode-v0-4e86"}

	t.Run("Locked", func(t *testing.T) {
		op := sampleProtoOperation()
		mgr := &fakeManager{
			getByResourceResp: &pb.GetByResourceResponse{Operation: op},
		}
		r := newTestRouter(mgr)

		resp, err := r.getByResourceHandler(&gin.Context{}, req)
		assert.NoError(t, err)
		assert.Equal(t, req.ResourceKey, mgr.lastGetByResourceKey)
		assert.True(t, resp.Locked)
		assertOperationMapped(t, op, resp.Operation)
	})

	t.Run("Unlocked", func(t *testing.T) {
		mgr := &fakeManager{
			getByResourceResp: &pb.GetByResourceResponse{},
		}
		r := newTestRouter(mgr)

		resp, err := r.getByResourceHandler(&gin.Context{}, req)
		assert.NoError(t, err)
		assert.False(t, resp.Locked)
		assert.Nil(t, resp.Operation)
	})

	t.Run("ManagerError", func(t *testing.T) {
		mgr := &fakeManager{getByResourceErr: errors.New("manager down")}
		r := newTestRouter(mgr)

		_, err := r.getByResourceHandler(&gin.Context{}, req)
		assert.Error(t, err)
		assert.Equal(t, req.ResourceKey, mgr.lastGetByResourceKey)
	})
}

func TestPostMarkRunningHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req := &MarkRunningRequest{
		Id:           "8e13fa4b-a8a7-40aa-8c61-2891cd16dc7f",
		FencingToken: 42,
	}

	t.Run("Success", func(t *testing.T) {
		op := sampleProtoOperation()
		mgr := &fakeManager{
			markRunningResp: &pb.MarkRunningResponse{Operation: op},
		}
		r := newTestRouter(mgr)

		resp, err := r.postMarkRunningHandler(&gin.Context{}, req)
		assert.NoError(t, err)
		assert.Equal(t, req.Id, mgr.lastMarkRunningId)
		assert.Equal(t, req.FencingToken, mgr.lastMarkRunningToken)
		assertOperationMapped(t, op, resp.Operation)
	})

	t.Run("ManagerError", func(t *testing.T) {
		mgr := &fakeManager{markRunningErr: errors.New("stale token")}
		r := newTestRouter(mgr)

		_, err := r.postMarkRunningHandler(&gin.Context{}, req)
		assert.Error(t, err)
		assert.Equal(t, req.Id, mgr.lastMarkRunningId)
		assert.Equal(t, req.FencingToken, mgr.lastMarkRunningToken)
	})
}

func TestPostForceUnlockHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := "8e13fa4b-a8a7-40aa-8c61-2891cd16dc7f"
	req := &ForceUnlockRequest{Id: uid, UserId: uid, Reason: "stuck"}

	t.Run("Success", func(t *testing.T) {
		op := sampleProtoOperation()
		mgr := &fakeManager{
			forceUnlockResp: &pb.ForceUnlockResponse{Operation: op},
		}
		r := newTestRouter(mgr)

		resp, err := r.postForceUnlockHandler(&gin.Context{}, req)
		assert.NoError(t, err)
		assert.Equal(t, req.Id, mgr.lastForceUnlockId)
		assert.Equal(t, uid, mgr.lastForceUnlockActor)
		assert.Equal(t, "stuck", mgr.lastForceUnlockReason)
		assertOperationMapped(t, op, resp.Operation)
	})

	t.Run("ManagerError", func(t *testing.T) {
		mgr := &fakeManager{forceUnlockErr: errors.New("manager down")}
		r := newTestRouter(mgr)

		_, err := r.postForceUnlockHandler(&gin.Context{}, req)
		assert.Error(t, err)
		assert.Equal(t, req.Id, mgr.lastForceUnlockId)
	})
}

func TestOperationFromProto(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, operationFromProto(nil))
	})

	t.Run("MapsFields", func(t *testing.T) {
		op := sampleProtoOperation()
		mapped := operationFromProto(op)
		assertOperationMapped(t, op, mapped)
	})
}

func TestTimestampAsTime(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		assert.Nil(t, timestampAsTime(nil))
	})

	t.Run("Invalid", func(t *testing.T) {
		assert.Nil(t, timestampAsTime(&timestamppb.Timestamp{Nanos: 1_000_000_000}))
	})

	t.Run("Valid", func(t *testing.T) {
		expected := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
		ts := timestamppb.New(expected)
		actual := timestampAsTime(ts)
		assert.NotNil(t, actual)
		assert.True(t, expected.Equal(*actual))
	})
}
