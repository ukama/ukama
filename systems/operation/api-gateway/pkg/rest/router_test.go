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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tj/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest"

	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

type fakeManager struct {
	lastStartReq          *pb.StartOperationRequest
	lastGetId             string
	lastGetByResourceKey  string
	lastMarkRunningId     string
	lastMarkRunningToken  uint64
	lastMarkRunningActor  string
	lastForceUnlockId     string
	lastForceUnlockActor  string
	lastForceUnlockReason string
	lastTerminateToken    uint64
	lastTerminateActor    string
	lastTerminateReason   string
	terminateErr          error

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
	lastDeleteId      string
	lastDeleteActor   string
	deleteErr         error
}

func (f *fakeManager) Delete(id, actor string) (*pb.DeleteOperationResponse, error) {
	f.lastDeleteId = id
	f.lastDeleteActor = actor
	if f.deleteErr != nil {
		return nil, f.deleteErr
	}
	return &pb.DeleteOperationResponse{}, nil
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

func (f *fakeManager) MarkRunning(id string, fencingToken uint64,
	actor string) (*pb.MarkRunningResponse, error) {
	f.lastMarkRunningId = id
	f.lastMarkRunningToken = fencingToken
	f.lastMarkRunningActor = actor
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

func (f *fakeManager) Complete(id string, fencingToken uint64, actor,
	reason string) (*pb.CompleteOperationResponse, error) {
	f.lastTerminateToken = fencingToken
	f.lastTerminateActor = actor
	f.lastTerminateReason = reason
	if f.terminateErr != nil {
		return nil, f.terminateErr
	}
	return &pb.CompleteOperationResponse{Operation: &pb.Operation{Id: id}}, nil
}

func (f *fakeManager) Fail(id string, fencingToken uint64, actor,
	reason string) (*pb.FailOperationResponse, error) {
	f.lastTerminateToken = fencingToken
	f.lastTerminateActor = actor
	f.lastTerminateReason = reason
	if f.terminateErr != nil {
		return nil, f.terminateErr
	}
	return &pb.FailOperationResponse{Operation: &pb.Operation{Id: id}}, nil
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
		ResourceKeys:   []string{"node:uk-sa2450-tnode-v0-4e86"},
		MustBeFree:     []string{"node:uk-sa2450-cnode-v0-4e86"},
		RequestedBy:    "user-1",
		IdempotencyKey: "idem-1",
		LeaseSeconds:   300,
	}

	t.Run("Success", func(t *testing.T) {
		op := sampleProtoOperation()
		member := sampleProtoOperation()
		member.GroupId = op.Id
		mgr := &fakeManager{
			startResp: &pb.StartOperationResponse{
				Operations: []*pb.Operation{op, member},
			},
		}
		r := newTestRouter(mgr)

		resp, err := r.postStartHandler(&gin.Context{}, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, req.Type, mgr.lastStartReq.Type)
		assert.Equal(t, req.System, mgr.lastStartReq.System)
		assert.Equal(t, req.ResourceKeys, mgr.lastStartReq.ResourceKeys)
		assert.Equal(t, req.MustBeFree, mgr.lastStartReq.MustBeFree)
		assert.Equal(t, req.RequestedBy, mgr.lastStartReq.RequestedBy)
		assert.Equal(t, req.IdempotencyKey, mgr.lastStartReq.IdempotencyKey)
		assert.Equal(t, req.LeaseSeconds, mgr.lastStartReq.LeaseSeconds)
		assert.Len(t, resp.Operations, 2)
		assertOperationMapped(t, op, resp.Operations[0])
		assert.Empty(t, resp.Operations[0].GroupId)
		assert.Equal(t, op.Id, resp.Operations[1].GroupId)
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

func TestPatchOperationHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := "8e13fa4b-a8a7-40aa-8c61-2891cd16dc7f"
	request := func(status string) *UpdateOperationRequest {
		return &UpdateOperationRequest{Id: uid, Status: status, FencingToken: 7,
			Actor: "controller", Reason: "dispatched"}
	}

	t.Run("RunningMarksRunning", func(t *testing.T) {
		op := sampleProtoOperation()
		mgr := &fakeManager{markRunningResp: &pb.MarkRunningResponse{Operation: op}}

		resp, err := newTestRouter(mgr).patchOperationHandler(&gin.Context{},
			request("RUNNING"))

		assert.NoError(t, err)
		assert.Equal(t, uid, mgr.lastMarkRunningId)
		assert.Equal(t, uint64(7), mgr.lastMarkRunningToken)
		assert.Equal(t, "controller", mgr.lastMarkRunningActor)
		assertOperationMapped(t, op, resp.Operation)
	})

	for _, st := range []string{"SUCCESS", "FAILED"} {
		t.Run(st+"Terminates", func(t *testing.T) {
			mgr := &fakeManager{}

			resp, err := newTestRouter(mgr).patchOperationHandler(&gin.Context{},
				request(st))

			assert.NoError(t, err)
			assert.Equal(t, uid, resp.Operation.Id)
			assert.Equal(t, uint64(7), mgr.lastTerminateToken)
			assert.Equal(t, "controller", mgr.lastTerminateActor)
			assert.Equal(t, "dispatched", mgr.lastTerminateReason)
		})
	}

	t.Run("OtherStatusIsBadRequest", func(t *testing.T) {
		_, err := newTestRouter(&fakeManager{}).patchOperationHandler(&gin.Context{},
			request("CANCELLED"))

		assert.Equal(t, rest.HttpError{HttpCode: http.StatusBadRequest,
			Message: `status "CANCELLED" cannot be set`}, err)
	})

	t.Run("PreconditionIsConflict", func(t *testing.T) {
		mgr := &fakeManager{markRunningErr: status.Error(codes.FailedPrecondition,
			"lease expired")}

		_, err := newTestRouter(mgr).patchOperationHandler(&gin.Context{},
			request("RUNNING"))

		assert.Equal(t, rest.HttpError{HttpCode: http.StatusConflict,
			Message: "lease expired"}, err)
	})

	t.Run("OtherErrorIsPassedThrough", func(t *testing.T) {
		mgr := &fakeManager{terminateErr: status.Error(codes.NotFound, "unknown")}

		_, err := newTestRouter(mgr).patchOperationHandler(&gin.Context{},
			request("FAILED"))

		assert.Equal(t, codes.NotFound, status.Code(err))
	})
}

func TestDeleteLockHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := "8e13fa4b-a8a7-40aa-8c61-2891cd16dc7f"
	req := &ReleaseLockRequest{Id: uid, Actor: "admin", Reason: "stuck"}

	t.Run("Success", func(t *testing.T) {
		op := sampleProtoOperation()
		mgr := &fakeManager{forceUnlockResp: &pb.ForceUnlockResponse{Operation: op}}

		resp, err := newTestRouter(mgr).deleteLockHandler(&gin.Context{}, req)

		assert.NoError(t, err)
		assert.Equal(t, uid, mgr.lastForceUnlockId)
		assert.Equal(t, "admin", mgr.lastForceUnlockActor)
		assert.Equal(t, "stuck", mgr.lastForceUnlockReason)
		assertOperationMapped(t, op, resp.Operation)
	})

	t.Run("FinishedOperationIsConflict", func(t *testing.T) {
		mgr := &fakeManager{forceUnlockErr: status.Error(codes.FailedPrecondition,
			"already finished")}

		_, err := newTestRouter(mgr).deleteLockHandler(&gin.Context{}, req)

		assert.Equal(t, rest.HttpError{HttpCode: http.StatusConflict,
			Message: "already finished"}, err)
	})
}

func TestDeleteOperationHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := "8e13fa4b-a8a7-40aa-8c61-2891cd16dc7f"
	req := &DeleteOperationRequest{Id: uid, Actor: "admin"}

	t.Run("Success", func(t *testing.T) {
		mgr := &fakeManager{}

		err := newTestRouter(mgr).deleteOperationHandler(&gin.Context{}, req)

		assert.NoError(t, err)
		assert.Equal(t, uid, mgr.lastDeleteId)
		assert.Equal(t, "admin", mgr.lastDeleteActor)
	})

	t.Run("ActiveOperationIsConflict", func(t *testing.T) {
		mgr := &fakeManager{deleteErr: status.Error(codes.FailedPrecondition,
			"operation is RUNNING")}

		err := newTestRouter(mgr).deleteOperationHandler(&gin.Context{}, req)

		assert.Equal(t, rest.HttpError{HttpCode: http.StatusConflict,
			Message: "operation is RUNNING"}, err)
	})

	t.Run("UnknownIsPassedThrough", func(t *testing.T) {
		mgr := &fakeManager{deleteErr: status.Error(codes.NotFound, "unknown")}

		err := newTestRouter(mgr).deleteOperationHandler(&gin.Context{}, req)

		assert.Equal(t, codes.NotFound, status.Code(err))
	})
}

func TestRouterOperationRoutes(t *testing.T) {
	id := "8e13fa4b-a8a7-40aa-8c61-2891cd16dc7f"
	busy := status.Error(codes.FailedPrecondition, "operation is RUNNING")

	for _, tc := range []struct {
		name   string
		method string
		path   string
		body   string
		mgr    *fakeManager
		code   int
	}{
		{"PatchRunningWithToken", http.MethodPatch, "",
			`{"status":"RUNNING","fencing_token":7,"actor":"controller"}`,
			&fakeManager{}, http.StatusOK},
		{"PatchSuccessWithToken", http.MethodPatch, "",
			`{"status":"SUCCESS","fencing_token":7,"actor":"controller","reason":"ok"}`,
			&fakeManager{}, http.StatusOK},
		{"PatchWithoutToken", http.MethodPatch, "",
			`{"status":"FAILED","actor":"controller","reason":"failed"}`,
			&fakeManager{}, http.StatusBadRequest},
		{"PatchUnknownStatus", http.MethodPatch, "",
			`{"status":"PAUSED","fencing_token":7,"actor":"controller"}`,
			&fakeManager{}, http.StatusBadRequest},
		{"PatchPreconditionIsConflict", http.MethodPatch, "",
			`{"status":"RUNNING","fencing_token":7,"actor":"controller"}`,
			&fakeManager{markRunningErr: busy}, http.StatusConflict},
		{"ReleaseLock", http.MethodDelete, "/lock?actor=admin&reason=stuck", "",
			&fakeManager{}, http.StatusOK},
		{"ReleaseLockWithoutReason", http.MethodDelete, "/lock?actor=admin", "",
			&fakeManager{}, http.StatusBadRequest},
		{"DeleteFinished", http.MethodDelete, "?actor=admin", "",
			&fakeManager{}, http.StatusOK},
		{"DeleteActiveIsConflict", http.MethodDelete, "?actor=admin", "",
			&fakeManager{deleteErr: busy}, http.StatusConflict},
		{"OldRunRouteIsGone", http.MethodPost, "/run",
			`{"fencing_token":7,"actor":"controller"}`,
			&fakeManager{}, http.StatusNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			serverConf := rest.DefaultHTTPConfig()

			r := NewRouter(&Clients{Manager: tc.mgr}, &RouterConfig{
				serverConf: &serverConf,
				auth:       &config.Auth{BypassAuthMode: true},
			}, nil)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, "/v1/operations/"+id+tc.path,
				strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			r.f.Engine().ServeHTTP(w, req)

			assert.Equal(t, tc.code, w.Code, w.Body.String())
		})
	}
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
