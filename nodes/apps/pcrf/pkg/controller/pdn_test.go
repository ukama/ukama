/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/session"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/datapath"
	"github.com/ukama/ukama/systems/common/rest"
	"github.com/ukama/ukama/systems/common/uuid"
)

// The controller tests exercise SQLite and stub only datapath session work.
// Actual counter, pause/resume and cleanup operations are tested in session.
type pdnSessions struct {
	store     *store.Store
	active    map[string]int
	created   int
	endErr    error
	resumeErr error
	dp        datapath.Status
}

func (s *pdnSessions) Status() session.Status {
	return session.Status{DataPath: s.dp}
}
func (s *pdnSessions) HasActiveSession(imsi string) bool {
	_, ok := s.active[imsi]
	return ok
}
func (s *pdnSessions) IfSessionExist(ctx context.Context,
	imsi, ip string) bool {
	return s.HasActiveSession(imsi)
}
func (s *pdnSessions) CreateSesssion(ctx context.Context,
	sub *store.Subscriber, se *store.Session, rx, tx *store.Flow) error {
	s.active[sub.Imsi] = se.ID
	s.created++
	return nil
}
func (s *pdnSessions) PauseSession(ctx context.Context,
	sub *store.Subscriber) error {
	return s.store.UpdateSessionFlowState(s.active[sub.Imsi], store.FlowsPaused)
}
func (s *pdnSessions) ResumeSession(ctx context.Context,
	sub *store.Subscriber) error {
	if s.resumeErr != nil {
		return s.resumeErr
	}
	return s.store.UpdateSessionFlowState(s.active[sub.Imsi], store.FlowsActive)
}
func (s *pdnSessions) EndSession(ctx context.Context,
	sub *store.Subscriber) error {
	if s.endErr != nil {
		return s.endErr
	}
	se, err := s.store.GetActiveSessionByImsi(sub.Imsi)
	if errors.Is(err, store.ErrSessionNotFound) {
		delete(s.active, sub.Imsi)
		return nil
	}
	if err != nil {
		return err
	}
	if err := s.store.EndSession(se); err != nil {
		return err
	}
	delete(s.active, sub.Imsi)
	return nil
}
func (s *pdnSessions) EndAllSessions() error {
	for imsi := range s.active {
		if err := s.EndSession(context.Background(),
			&store.Subscriber{Imsi: imsi}); err != nil {
			return err
		}
	}
	return nil
}

var pdnTestNumber int

func pdnFixture(t *testing.T) (*Controller, *pdnSessions, *api.CreateSession) {
	t.Helper()
	s := newTestStore(t)
	sm := &pdnSessions{store: s, active: make(map[string]int),
		dp: datapath.Status{Connected: true, ConnectedCount: 1}}
	c := &Controller{store: s, sm: sm, apn: "internet",
		nodeId: "test-node", serviceOn: true}
	boot := uuid.NewV4().String()
	_, err := c.RegisterEPC(&gin.Context{}, &api.RegisterEPC{InstanceID: boot})
	require.NoError(t, err)
	pdnTestNumber++
	req := &api.CreateSession{InstanceID: boot, PDNID: uuid.NewV4().String(),
		ImsiStr: fmt.Sprintf("999993%09d", pdnTestNumber),
		IpStr:   "192.168.8.2", ApnName: "internet"}
	t.Cleanup(func() {
		if c.pdnCancel != nil {
			c.pdnCancel()
			<-c.pdnDone
		}
		sm.endErr = nil
		require.NoError(t, c.store.RetireEPCs())
		require.NoError(t, c.cleanupPDNsLocked(context.Background()))
	})
	return c, sm, req
}

func pdnPolicy() api.Policy {
	return api.Policy{Uuid: uuid.NewV4(), Data: 1000000,
		StartTime: time.Now().Add(-time.Hour).Unix(),
		EndTime:   time.Now().Add(time.Hour).Unix()}
}
func pushPDNPolicy(t *testing.T, c *Controller, imsi string, p api.Policy) {
	t.Helper()
	require.NoError(t, c.UpdateSubscriber(&gin.Context{},
		&api.UpdateSubscriber{Imsi: imsi, Policy: p}))
}
func checkHTTPError(t *testing.T, err error, code int) {
	t.Helper()
	var httpErr rest.HttpError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, code, httpErr.HttpCode)
}
func pdnDelete(req *api.CreateSession) *api.EndSession {
	return &api.EndSession{InstanceID: req.InstanceID,
		PDNID: req.PDNID, ImsiStr: req.ImsiStr}
}

func TestPDNInvalidAPNHasNoSideEffects(t *testing.T) {
	c, sm, req := pdnFixture(t)
	for _, apn := range []string{"", "wrong", " internet", "internet.",
		"a..b", "a/b", strings.Repeat("a", 100)} {
		req.ApnName = apn
		_, err := c.CreateSession(&gin.Context{}, req)
		require.Error(t, err, apn)
		_, err = c.store.GetLivePDN(req.ImsiStr)
		require.ErrorIs(t, err, store.ErrPDNNotFound)
		require.Zero(t, sm.created)
	}
	req.ApnName = "internet"
	result, err := c.CreateSession(&gin.Context{}, req)
	require.NoError(t, err)
	require.Equal(t, "blocked", result.DataState)
}

func TestPDNPolicyArrivalPreservesAPNAndBinding(t *testing.T) {
	c, sm, req := pdnFixture(t)
	req.ApnName = "INTERNET"
	result, err := c.CreateSession(&gin.Context{}, req)
	require.NoError(t, err)
	require.Equal(t, "policy_missing", result.Reason)
	require.Equal(t, "internet", result.ApnName)
	require.Zero(t, sm.created)
	pushPDNPolicy(t, c, req.ImsiStr, pdnPolicy())
	got, err := c.GetPDN(&gin.Context{}, &api.GetPDN{Imsi: req.ImsiStr})
	require.NoError(t, err)
	require.Equal(t, req.PDNID, got.PDNID)
	require.Equal(t, "allowed", got.DataState)
	se, err := c.store.GetActiveSessionByImsi(req.ImsiStr)
	require.NoError(t, err)
	require.Equal(t, "internet", store.PrepareCDR(se).ApnName)
	flows, err := c.GetFlowsForImsi(&gin.Context{},
		&api.GetFlowsForImsi{Imsi: req.ImsiStr})
	require.NoError(t, err)
	require.Len(t, flows, 2)
}

func TestPDNPolicyBeforeAttachCreatesNoForwarding(t *testing.T) {
	c, sm, req := pdnFixture(t)
	pushPDNPolicy(t, c, req.ImsiStr, pdnPolicy())
	require.Zero(t, sm.created)
	_, err := c.store.GetActiveSessionByImsi(req.ImsiStr)
	require.ErrorIs(t, err, store.ErrSessionNotFound)
	result, err := c.CreateSession(&gin.Context{}, req)
	require.NoError(t, err)
	require.Equal(t, "allowed", result.DataState)
	original := sm.active[req.ImsiStr]
	for i := 0; i < 3; i++ {
		_, err := c.CreateSession(&gin.Context{}, req)
		require.NoError(t, err)
	}
	require.Equal(t, original, sm.active[req.ImsiStr])
	require.Equal(t, 1, sm.created)
	conflict := *req
	conflict.IpStr = "192.168.8.3"
	_, err = c.CreateSession(&gin.Context{}, &conflict)
	checkHTTPError(t, err, http.StatusConflict)
	conflict = *req
	conflict.PDNID = uuid.NewV4().String()
	_, err = c.CreateSession(&gin.Context{}, &conflict)
	checkHTTPError(t, err, http.StatusConflict)
	conflict.ImsiStr = "999993888888888"
	_, err = c.CreateSession(&gin.Context{}, &conflict)
	checkHTTPError(t, err, http.StatusConflict)
	require.Equal(t, original, sm.active[req.ImsiStr])
}

func TestPDNPolicyReasons(t *testing.T) {
	cases := []string{"package_expired", "package_not_started",
		"quota_exhausted", "subscriber_disabled", "service_disabled"}
	for _, reason := range cases {
		t.Run(reason, func(t *testing.T) {
			c, sm, req := pdnFixture(t)
			p := pdnPolicy()
			switch reason {
			case "package_expired":
				p.EndTime = time.Now().Unix()
			case "package_not_started":
				p.StartTime = time.Now().Add(time.Hour).Unix()
			case "quota_exhausted":
				p.Consumed = p.Data + 1
			}
			pushPDNPolicy(t, c, req.ImsiStr, p)
			if reason == "subscriber_disabled" {
				sub, err := c.store.GetSubscriber(req.ImsiStr)
				require.NoError(t, err)
				require.NoError(t,
					c.store.UpdateSubscriberServiceStatus(sub, false))
			}
			if reason == "service_disabled" {
				c.serviceOn = false
			}
			result, err := c.CreateSession(&gin.Context{}, req)
			require.NoError(t, err)
			require.Equal(t, "live", result.State)
			require.Equal(t, "blocked", result.DataState)
			require.Equal(t, reason, result.Reason)
			require.Zero(t, sm.created)
		})
	}
}

func TestPDNServiceToggleKeepsPDNs(t *testing.T) {
	c, sm, req := pdnFixture(t)
	p := pdnPolicy()
	pushPDNPolicy(t, c, req.ImsiStr, p)
	_, err := c.CreateSession(&gin.Context{}, req)
	require.NoError(t, err)
	original := sm.active[req.ImsiStr]
	require.NoError(t, c.SetService(context.Background(),
		&api.ServiceRequest{State: "off"}))
	got, err := c.GetPDN(&gin.Context{}, &api.GetPDN{Imsi: req.ImsiStr})
	require.NoError(t, err)
	require.Equal(t, "blocked", got.DataState)
	require.Equal(t, original, sm.active[req.ImsiStr])
	other := *req
	other.ImsiStr = "999993777777777"
	other.IpStr = "192.168.8.3"
	other.PDNID = uuid.NewV4().String()
	pushPDNPolicy(t, c, other.ImsiStr, pdnPolicy())
	got, err = c.CreateSession(&gin.Context{}, &other)
	require.NoError(t, err)
	require.Equal(t, "blocked", got.DataState)
	require.NoError(t, c.SetService(context.Background(),
		&api.ServiceRequest{State: "on"}))
	require.Equal(t, original, sm.active[req.ImsiStr])
	require.True(t, sm.HasActiveSession(other.ImsiStr))
	off := false
	require.NoError(t, c.UpdateSubscriber(&gin.Context{},
		&api.UpdateSubscriber{Imsi: req.ImsiStr, Policy: p,
			IsServiceOn: &off}))
	got, err = c.GetPDN(&gin.Context{}, &api.GetPDN{Imsi: req.ImsiStr})
	require.NoError(t, err)
	require.Equal(t, "subscriber_disabled", got.Reason)
	require.Equal(t, original, sm.active[req.ImsiStr])
}

func TestPDNExpirationIsCheckedWhileIdle(t *testing.T) {
	c, sm, req := pdnFixture(t)
	p := pdnPolicy()
	p.EndTime = time.Now().Unix() + 2
	pushPDNPolicy(t, c, req.ImsiStr, p)
	_, err := c.CreateSession(&gin.Context{}, req)
	require.NoError(t, err)
	original := sm.active[req.ImsiStr]
	c.startPDNMonitor(10 * time.Millisecond)
	require.Eventually(t, func() bool {
		got, err := c.GetPDN(&gin.Context{}, &api.GetPDN{Imsi: req.ImsiStr})
		return err == nil && got.DataState == "blocked" &&
			got.Reason == "package_expired"
	}, 3*time.Second, 10*time.Millisecond)
	c.pdnMu.Lock()
	require.Equal(t, original, sm.active[req.ImsiStr])
	c.pdnMu.Unlock()
}

func TestPDNQuotaAndNewPackageRetainBinding(t *testing.T) {
	c, sm, req := pdnFixture(t)
	p := pdnPolicy()
	pushPDNPolicy(t, c, req.ImsiStr, p)
	_, err := c.CreateSession(&gin.Context{}, req)
	require.NoError(t, err)
	se, err := c.store.GetActiveSessionByImsi(req.ImsiStr)
	require.NoError(t, err)
	se.TxBytes = p.Data
	require.NoError(t, c.store.UpdateSessionUsage(se))
	require.NoError(t, c.reconcileImsiLocked(context.Background(), req.ImsiStr))
	got, err := c.GetPDN(&gin.Context{}, &api.GetPDN{Imsi: req.ImsiStr})
	require.NoError(t, err)
	require.Equal(t, "quota_exhausted", got.Reason)
	require.Equal(t, se.ID, sm.active[req.ImsiStr])
	next := pdnPolicy()
	next.StartTime = p.StartTime + 1
	pushPDNPolicy(t, c, req.ImsiStr, next)
	got, err = c.GetPDN(&gin.Context{}, &api.GetPDN{Imsi: req.ImsiStr})
	require.NoError(t, err)
	require.Equal(t, "allowed", got.DataState)
	require.Equal(t, req.PDNID, got.PDNID)
	require.NotEqual(t, se.ID, sm.active[req.ImsiStr])
	old, err := c.store.GetSessionByID(se.ID)
	require.NoError(t, err)
	require.Equal(t, store.SessionCompleted, old.State)
	require.Equal(t, p.Uuid, old.PolicyID.ID)
}

func TestPDNFailedDeleteAndLateRetry(t *testing.T) {
	c, sm, req := pdnFixture(t)
	pushPDNPolicy(t, c, req.ImsiStr, pdnPolicy())
	_, err := c.CreateSession(&gin.Context{}, req)
	require.NoError(t, err)
	sm.endErr = errors.New("OVS unavailable")
	checkHTTPError(t, c.EndSession(&gin.Context{}, pdnDelete(req)),
		http.StatusServiceUnavailable)
	p, err := c.store.GetLivePDN(req.ImsiStr)
	require.NoError(t, err)
	require.Equal(t, store.PDNClosing, p.State)
	retry := *req
	retry.PDNID = uuid.NewV4().String()
	_, err = c.CreateSession(&gin.Context{}, &retry)
	checkHTTPError(t, err, http.StatusServiceUnavailable)
	sm.endErr = nil
	require.NoError(t, c.EndSession(&gin.Context{}, pdnDelete(req)))
	require.NoError(t, c.EndSession(&gin.Context{}, pdnDelete(req)))
	_, err = c.CreateSession(&gin.Context{}, req)
	checkHTTPError(t, err, http.StatusGone)
	_, err = c.CreateSession(&gin.Context{}, &retry)
	require.NoError(t, err)
	replacement := sm.active[req.ImsiStr]
	require.NoError(t, c.EndSession(&gin.Context{}, pdnDelete(req)))
	require.Equal(t, replacement, sm.active[req.ImsiStr])
}

func TestPDNDeleteCanOvertakeCreate(t *testing.T) {
	c, _, req := pdnFixture(t)
	require.NoError(t, c.EndSession(&gin.Context{}, pdnDelete(req)))
	_, err := c.CreateSession(&gin.Context{}, req)
	checkHTTPError(t, err, http.StatusGone)
}

func TestPDNEPCRestartRetiresOldIdentity(t *testing.T) {
	c, _, req := pdnFixture(t)
	pushPDNPolicy(t, c, req.ImsiStr, pdnPolicy())
	_, err := c.CreateSession(&gin.Context{}, req)
	require.NoError(t, err)
	newBoot := uuid.NewV4().String()
	_, err = c.RegisterEPC(&gin.Context{},
		&api.RegisterEPC{InstanceID: newBoot,
			ExpectedInstanceID: req.InstanceID})
	require.NoError(t, err)
	_, err = c.RegisterEPC(&gin.Context{},
		&api.RegisterEPC{InstanceID: req.InstanceID})
	checkHTTPError(t, err, http.StatusGone)
	_, err = c.CreateSession(&gin.Context{}, req)
	checkHTTPError(t, err, http.StatusGone)
	old := *req
	req.InstanceID = newBoot
	req.PDNID = uuid.NewV4().String()
	_, err = c.CreateSession(&gin.Context{}, req)
	require.NoError(t, err)
	require.NoError(t, c.EndSession(&gin.Context{}, pdnDelete(&old)))
	got, err := c.GetPDN(&gin.Context{}, &api.GetPDN{Imsi: req.ImsiStr})
	require.NoError(t, err)
	require.Equal(t, newBoot, got.InstanceID)
}

func TestPDNGetDoesNotInstallFlows(t *testing.T) {
	c, sm, req := pdnFixture(t)
	_, err := c.CreateSession(&gin.Context{}, req)
	require.NoError(t, err)
	p := pdnPolicy()
	_, err = c.store.UpdateSubscriber(req.ImsiStr, &p)
	require.NoError(t, err)
	got, err := c.GetPDN(&gin.Context{}, &api.GetPDN{Imsi: req.ImsiStr})
	require.NoError(t, err)
	require.Equal(t, "blocked", got.DataState)
	require.Equal(t, "policy_update_pending", got.Reason)
	require.Zero(t, sm.created)
}

func TestPDNFailedInstallationCannotLaterAuthorize(t *testing.T) {
	c, sm, req := pdnFixture(t)
	pushPDNPolicy(t, c, req.ImsiStr, pdnPolicy())
	sm.resumeErr = errors.New("partial installation")
	sm.endErr = errors.New("cleanup unavailable")
	_, err := c.CreateSession(&gin.Context{}, req)
	checkHTTPError(t, err, http.StatusServiceUnavailable)
	p, err := c.store.GetLivePDN(req.ImsiStr)
	require.NoError(t, err)
	require.Equal(t, store.PDNClosing, p.State)
	sm.endErr = nil
	sm.resumeErr = nil
	require.NoError(t, c.reconcilePDNsLocked(context.Background()))
	_, err = c.store.GetLivePDN(req.ImsiStr)
	require.ErrorIs(t, err, store.ErrPDNNotFound)
	require.False(t, sm.HasActiveSession(req.ImsiStr))
}

func TestPDNConcurrentDuplicateRegistration(t *testing.T) {
	c, sm, req := pdnFixture(t)
	pushPDNPolicy(t, c, req.ImsiStr, pdnPolicy())
	var wg sync.WaitGroup
	errs := make(chan error, 16)
	for i := 0; i < cap(errs); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := c.CreateSession(&gin.Context{}, req)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
	require.Equal(t, 1, sm.created)
}

func TestPDNLateBootCannotReplaceNewerProcess(t *testing.T) {
	c, _, req := pdnFixture(t)
	next := &api.RegisterEPC{InstanceID: uuid.NewV4().String(),
		ExpectedInstanceID: req.InstanceID}
	_, err := c.RegisterEPC(&gin.Context{}, next)
	require.NoError(t, err)
	late := &api.RegisterEPC{InstanceID: uuid.NewV4().String(),
		ExpectedInstanceID: req.InstanceID}
	_, err = c.RegisterEPC(&gin.Context{}, late)
	checkHTTPError(t, err, http.StatusConflict)
	got, err := c.GetEPC(&gin.Context{})
	require.NoError(t, err)
	require.Equal(t, next.InstanceID, got.InstanceID)
}
