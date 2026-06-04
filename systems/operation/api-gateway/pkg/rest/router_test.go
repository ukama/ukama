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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tj/assert"

	crest "github.com/ukama/ukama/systems/common/rest"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	pb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

type fakeManager struct {
	forceUnlockCalled bool
}

func (f *fakeManager) Start(*pb.StartOperationRequest) (*pb.StartOperationResponse, error) {
	return &pb.StartOperationResponse{}, nil
}
func (f *fakeManager) Get(string) (*pb.GetOperationResponse, error) {
	return &pb.GetOperationResponse{}, nil
}
func (f *fakeManager) GetByResource(string) (*pb.GetByResourceResponse, error) {
	return &pb.GetByResourceResponse{}, nil
}
func (f *fakeManager) MarkRunning(string, uint64) (*pb.MarkRunningResponse, error) {
	return &pb.MarkRunningResponse{}, nil
}
func (f *fakeManager) ForceUnlock(id, actor, reason string) (*pb.ForceUnlockResponse, error) {
	f.forceUnlockCalled = true
	return &pb.ForceUnlockResponse{Operation: &pb.Operation{Id: id}}, nil
}

type fakeMember struct {
	role string
	err  error
}

func (f *fakeMember) GetByUserId(string) (*creg.MemberInfoResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &creg.MemberInfoResponse{Member: creg.MemberInfo{Role: f.role}}, nil
}

func newTestRouter(role string, memberErr error) (*Router, *fakeManager) {
	mgr := &fakeManager{}
	return &Router{
		clients: &Clients{
			Manager: mgr,
			Member:  &fakeMember{role: role, err: memberErr},
		},
	}, mgr
}

func TestForceUnlockAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := "8e13fa4b-a8a7-40aa-8c61-2891cd16dc7f"
	req := &ForceUnlockRequest{Id: "op-1", UserId: uid, Reason: "stuck"}

	t.Run("OwnerAllowed", func(t *testing.T) {
		r, mgr := newTestRouter("ROLE_OWNER", nil)
		resp, err := r.deleteForceUnlockHandler(&gin.Context{}, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, mgr.forceUnlockCalled)
	})

	t.Run("AdminAllowed", func(t *testing.T) {
		r, mgr := newTestRouter("ROLE_ADMIN", nil)
		_, err := r.deleteForceUnlockHandler(&gin.Context{}, req)
		assert.NoError(t, err)
		assert.True(t, mgr.forceUnlockCalled)
	})

	t.Run("UserDenied", func(t *testing.T) {
		r, mgr := newTestRouter("ROLE_USER", nil)
		_, err := r.deleteForceUnlockHandler(&gin.Context{}, req)
		assert.Error(t, err)
		he, ok := err.(crest.HttpError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, he.HttpCode)
		assert.False(t, mgr.forceUnlockCalled)
	})

	t.Run("MemberLookupFailureDenied", func(t *testing.T) {
		r, mgr := newTestRouter("", errors.New("registry down"))
		_, err := r.deleteForceUnlockHandler(&gin.Context{}, req)
		assert.Error(t, err)
		he, ok := err.(crest.HttpError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, he.HttpCode)
		assert.False(t, mgr.forceUnlockCalled)
	})
}
