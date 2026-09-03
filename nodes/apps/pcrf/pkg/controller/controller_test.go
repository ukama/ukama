/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/nodes/apps/pcrf/mocks"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/client"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"
)

// store.NewStore registers a global sql driver name on every call, which
// panics if called more than once per process - so tests in this file share
// a single Store instance (distinct imsi/reroute IPs per test keep them from
// colliding in the same database).
var (
	sharedStoreOnce sync.Once
	sharedStore     *store.Store
)

func newTestStore(t *testing.T) *store.Store {
	t.Helper()

	sharedStoreOnce.Do(func() {
		dir, err := os.MkdirTemp("", "pcrf-controller-test")
		assert.NoError(t, err)

		s, err := store.NewStore(filepath.Join(dir, "pcrf_test.db"))
		assert.NoError(t, err)
		sharedStore = s
	})

	return sharedStore
}

func newTestSubscriber(t *testing.T, s *store.Store, imsi string, serviceOn bool) *store.Subscriber {
	t.Helper()

	policy := &api.Policy{
		Uuid:      uuid.NewV4(),
		Ulbr:      1000,
		Dlbr:      2000,
		Data:      1_000_000,
		Burst:     100,
		StartTime: time.Now().Add(-time.Hour).Unix(),
		EndTime:   time.Now().Add(time.Hour).Unix(),
	}
	ip := "192.168.9." + imsi[len(imsi)-2:]

	sub, err := s.CreateSubscriber(imsi, policy, &ip, nil)
	assert.NoError(t, err)

	if !serviceOn {
		assert.NoError(t, s.UpdateSubscriberServiceStatus(sub, false))
		sub, err = s.GetSubscriber(imsi)
		assert.NoError(t, err)
	}

	return sub
}

// toAPIPolicy rebuilds the api.Policy an UpdateSubscriber request would carry
// for an already-persisted policy, so tests can round-trip it unchanged.
func toAPIPolicy(p store.Policy) api.Policy {
	return api.Policy{
		Uuid:      p.ID,
		Ulbr:      p.Ulbr,
		Dlbr:      p.Dlbr,
		Data:      p.Data,
		Consumed:  p.Consumed,
		Burst:     p.Burst,
		StartTime: p.StartTime,
		EndTime:   p.EndTime,
	}
}

func TestController_UpdateSubscriber_FlipToOffWithActiveSession_Pauses(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)

	sub := newTestSubscriber(t, s, "999992000000001", true)

	sm.On("HasActiveSession", sub.Imsi).Return(true).Once()
	sm.On("PauseSession", mock.Anything, mock.MatchedBy(func(x *store.Subscriber) bool {
		return x.Imsi == sub.Imsi
	})).Return(nil).Once()

	c := &Controller{store: s, sm: sm}

	off := false
	err := c.UpdateSubscriber(&gin.Context{}, &api.UpdateSubscriber{
		Imsi:        sub.Imsi,
		Policy:      toAPIPolicy(sub.PolicyID),
		IsServiceOn: &off,
	})
	assert.NoError(t, err)

	got, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.False(t, got.ServiceOn)
}

func TestController_UpdateSubscriber_PauseSessionFails_DoesNotPersist(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)

	sub := newTestSubscriber(t, s, "999992000000009", true)

	sm.On("HasActiveSession", sub.Imsi).Return(true).Once()
	sm.On("PauseSession", mock.Anything, mock.MatchedBy(func(x *store.Subscriber) bool {
		return x.Imsi == sub.Imsi
	})).Return(fmt.Errorf("ovs error")).Once()

	c := &Controller{store: s, sm: sm}

	off := false
	err := c.UpdateSubscriber(&gin.Context{}, &api.UpdateSubscriber{
		Imsi:        sub.Imsi,
		Policy:      toAPIPolicy(sub.PolicyID),
		IsServiceOn: &off,
	})
	assert.Error(t, err)

	got, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.True(t, got.ServiceOn)
}

func TestController_UpdateSubscriber_FlipToOnWithActiveSession_Resumes(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)

	sub := newTestSubscriber(t, s, "999992000000002", false)

	sm.On("HasActiveSession", sub.Imsi).Return(true).Once()
	sm.On("ResumeSession", mock.Anything, mock.MatchedBy(func(x *store.Subscriber) bool {
		return x.Imsi == sub.Imsi
	})).Return(nil).Once()

	c := &Controller{store: s, sm: sm}

	on := true
	err := c.UpdateSubscriber(&gin.Context{}, &api.UpdateSubscriber{
		Imsi:        sub.Imsi,
		Policy:      toAPIPolicy(sub.PolicyID),
		IsServiceOn: &on,
	})
	assert.NoError(t, err)

	got, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.True(t, got.ServiceOn)
}

func TestController_UpdateSubscriber_FlipWithNoActiveSession_PersistsOnly(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)

	sub := newTestSubscriber(t, s, "999992000000003", true)

	sm.On("HasActiveSession", sub.Imsi).Return(false).Once()
	// No PauseSession/ResumeSession expectation registered at all - if the
	// controller called either, mocks.NewSessionManager(t)'s built-in
	// AssertExpectations-on-cleanup would fail the test.

	c := &Controller{store: s, sm: sm}

	off := false
	err := c.UpdateSubscriber(&gin.Context{}, &api.UpdateSubscriber{
		Imsi:        sub.Imsi,
		Policy:      toAPIPolicy(sub.PolicyID),
		IsServiceOn: &off,
	})
	assert.NoError(t, err)

	got, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.False(t, got.ServiceOn, "status must still be persisted for the next attach to respect")
}

func TestController_UpdateSubscriber_UnchangedStatus_IsNoOp(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)

	sub := newTestSubscriber(t, s, "999992000000004", true)

	// No HasActiveSession/PauseSession/ResumeSession expectations at all:
	// nil IsServiceOn and an explicit-but-unchanged true must both be
	// complete no-ops on the session manager.

	c := &Controller{store: s, sm: sm}

	err := c.UpdateSubscriber(&gin.Context{}, &api.UpdateSubscriber{
		Imsi:   sub.Imsi,
		Policy: toAPIPolicy(sub.PolicyID),
	})
	assert.NoError(t, err)

	unchanged := true
	err = c.UpdateSubscriber(&gin.Context{}, &api.UpdateSubscriber{
		Imsi:        sub.Imsi,
		Policy:      toAPIPolicy(sub.PolicyID),
		IsServiceOn: &unchanged,
	})
	assert.NoError(t, err)
}

func TestController_UpdateSubscriber_UnknownSubscriber_AutoCreates(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)

	imsi := "999992000000006"
	policy := api.Policy{
		Uuid:      uuid.NewV4(),
		Ulbr:      1000,
		Dlbr:      2000,
		Data:      1_000_000,
		Burst:     100,
		StartTime: time.Now().Add(-time.Hour).Unix(),
		EndTime:   time.Now().Add(time.Hour).Unix(),
	}

	sm.On("HasActiveSession", imsi).Return(false).Once()

	c := &Controller{store: s, sm: sm}

	off := false
	err := c.UpdateSubscriber(&gin.Context{}, &api.UpdateSubscriber{
		Imsi:        imsi,
		Policy:      policy,
		IsServiceOn: &off,
	})
	assert.NoError(t, err, "UpdateSubscriber for an unknown imsi must auto-create it, not fail")

	got, err := s.GetSubscriber(imsi)
	assert.NoError(t, err)
	assert.False(t, got.ServiceOn, "the newly auto-created subscriber must reflect the requested status")
}

func TestController_AddSubscriber_FlipToOffWithActiveSession_Pauses(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)

	sub := newTestSubscriber(t, s, "999992000000007", true)

	sm.On("HasActiveSession", sub.Imsi).Return(true).Once()
	sm.On("PauseSession", mock.Anything, mock.MatchedBy(func(x *store.Subscriber) bool {
		return x.Imsi == sub.Imsi
	})).Return(nil).Once()

	c := &Controller{store: s, sm: sm}

	off := false
	err := c.AddSubscriber(&gin.Context{}, &api.CreateSubscriber{
		Imsi: sub.Imsi,
		Policy: api.Policy{
			Uuid:      uuid.NewV4(),
			Ulbr:      1000,
			Dlbr:      2000,
			Data:      1_000_000,
			Burst:     100,
			StartTime: time.Now().Add(-time.Hour).Unix(),
			EndTime:   time.Now().Add(time.Hour).Unix(),
		},
		ReRoute:     "192.168.9.98",
		IsServiceOn: &off,
	})
	assert.NoError(t, err)

	got, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.False(t, got.ServiceOn)
}

func TestController_AddSubscriber_PauseSessionFails_DoesNotPersist(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)

	sub := newTestSubscriber(t, s, "999992000000010", true)

	sm.On("HasActiveSession", sub.Imsi).Return(true).Once()
	sm.On("PauseSession", mock.Anything, mock.MatchedBy(func(x *store.Subscriber) bool {
		return x.Imsi == sub.Imsi
	})).Return(fmt.Errorf("ovs error")).Once()

	c := &Controller{store: s, sm: sm}

	off := false
	err := c.AddSubscriber(&gin.Context{}, &api.CreateSubscriber{
		Imsi: sub.Imsi,
		Policy: api.Policy{
			Uuid:      uuid.NewV4(),
			Ulbr:      1000,
			Dlbr:      2000,
			Data:      1_000_000,
			Burst:     100,
			StartTime: time.Now().Add(-time.Hour).Unix(),
			EndTime:   time.Now().Add(time.Hour).Unix(),
		},
		ReRoute:     "192.168.9.96",
		IsServiceOn: &off,
	})
	assert.Error(t, err)

	got, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.True(t, got.ServiceOn)
}

func TestController_AddSubscriber_FlipWithNoActiveSession_PersistsOnly(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)

	sub := newTestSubscriber(t, s, "999992000000008", true)

	sm.On("HasActiveSession", sub.Imsi).Return(false).Once()

	c := &Controller{store: s, sm: sm}

	off := false
	err := c.AddSubscriber(&gin.Context{}, &api.CreateSubscriber{
		Imsi: sub.Imsi,
		Policy: api.Policy{
			Uuid:      uuid.NewV4(),
			Ulbr:      1000,
			Dlbr:      2000,
			Data:      1_000_000,
			Burst:     100,
			StartTime: time.Now().Add(-time.Hour).Unix(),
			EndTime:   time.Now().Add(time.Hour).Unix(),
		},
		ReRoute:     "192.168.9.97",
		IsServiceOn: &off,
	})
	assert.NoError(t, err)

	got, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.False(t, got.ServiceOn)
}

func TestController_CreateSession_IgnoresServiceStatus(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)

	sub := newTestSubscriber(t, s, "999992000000005", false)

	sm.On("IfSessionExist", mock.Anything, sub.Imsi, "10.10.10.30").Return(false).Once()
	sm.On("CreateSesssion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	c := &Controller{store: s, sm: sm, serviceOn: true}

	req := &api.CreateSession{
		ImsiStr: sub.Imsi,
		IpStr:   "10.10.10.30",
	}

	// A subscriber whose service is off must still be able to attach -
	// CreateSession must not gate on ServiceOn anywhere in its path.
	err := c.CreateSession(&gin.Context{}, req)
	assert.NoError(t, err)
}

func TestController_CreateSession_UnknownSubscriber_RemoteNotFound_Returns404(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)
	rc := mocks.NewRemoteController(t)

	imsi := "999992000000011"
	rc.On("GetSubscriberProfile", imsi).Return(nil, client.ErrRemoteSubscriberNotFound).Once()

	c := &Controller{store: s, sm: sm, rc: rc, serviceOn: true}

	err := c.CreateSession(&gin.Context{}, &api.CreateSession{ImsiStr: imsi, IpStr: "10.10.10.51"})
	assert.Error(t, err)

	httpErr, ok := err.(rest.HttpError)
	assert.True(t, ok, "expected rest.HttpError, got %T: %v", err, err)
	assert.Equal(t, http.StatusNotFound, httpErr.HttpCode)
}

func TestController_CreateSession_UnknownSubscriber_RemoteLookupFails_Returns503(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)
	rc := mocks.NewRemoteController(t)

	imsi := "999992000000012"
	rc.On("GetSubscriberProfile", imsi).Return(nil, fmt.Errorf("connection refused")).Once()

	c := &Controller{store: s, sm: sm, rc: rc, serviceOn: true}

	err := c.CreateSession(&gin.Context{}, &api.CreateSession{ImsiStr: imsi, IpStr: "10.10.10.52"})
	assert.Error(t, err)

	httpErr, ok := err.(rest.HttpError)
	assert.True(t, ok, "expected rest.HttpError, got %T: %v", err, err)
	assert.Equal(t, http.StatusServiceUnavailable, httpErr.HttpCode)
}

func TestController_CreateSession_StillInvalidAfterReverseLookup_Returns403(t *testing.T) {
	s := newTestStore(t)
	sm := mocks.NewSessionManager(t)
	rc := mocks.NewRemoteController(t)

	imsi := "999992000000013"
	now := time.Now()

	expiredPolicy := &api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      1_000_000,
		Burst:     100,
		StartTime: now.Add(-time.Hour).Unix(),
		EndTime:   now.Add(-time.Minute).Unix(),
	}
	ip := "192.168.9.13"
	_, err := s.CreateSubscriber(imsi, expiredPolicy, &ip, nil)
	assert.NoError(t, err)

	stillExpiredRemotePolicy := api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      1_000_000,
		Burst:     100,
		StartTime: now.Add(-time.Second).Unix(),
		EndTime:   now.Add(-time.Second).Unix(),
	}
	rc.On("GetSubscriberProfile", imsi).Return(&api.Spr{
		Imsi:    imsi,
		Policy:  stillExpiredRemotePolicy,
		ReRoute: ip,
	}, nil).Once()

	c := &Controller{store: s, sm: sm, rc: rc, serviceOn: true}

	err = c.CreateSession(&gin.Context{}, &api.CreateSession{ImsiStr: imsi, IpStr: "10.10.10.53"})
	assert.Error(t, err)

	httpErr, ok := err.(rest.HttpError)
	assert.True(t, ok, "expected rest.HttpError, got %T: %v", err, err)
	assert.Equal(t, http.StatusForbidden, httpErr.HttpCode)
}
