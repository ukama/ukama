/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package session

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/datapath"
	"github.com/ukama/ukama/systems/common/uuid"
)

// store.NewStore registers a global sql driver name on every call, which
// panics if called more than once per process - so every test in this file
// must share a single Store instance (distinct imsi/reroute IPs per test
// keep them from colliding in the same database).
var (
	sharedStoreOnce sync.Once
	sharedStore     *store.Store
)

type fakeDataPath struct {
	mock.Mock
}

func (f *fakeDataPath) AddNewDataPath(ip string, rxMeter, txMeter, rxRate, txRate, burstSize uint32, rxCookie, txCookie uint64) error {
	args := f.Called(ip, rxMeter, txMeter, rxRate, txRate, burstSize, rxCookie, txCookie)
	return args.Error(0)
}

func (f *fakeDataPath) DeleteDataPath(ip string, rxMeter, txMeter uint32) error {
	args := f.Called(ip, rxMeter, txMeter)
	return args.Error(0)
}

func (f *fakeDataPath) AddFlowOnly(ip string, rxMeter, txMeter uint32, rxCookie, txCookie uint64) error {
	args := f.Called(ip, rxMeter, txMeter, rxCookie, txCookie)
	return args.Error(0)
}

func (f *fakeDataPath) DeleteFlowOnly(ip string) error {
	args := f.Called(ip)
	return args.Error(0)
}

func (f *fakeDataPath) AddMetersOnly(rxMeter, txMeter, rxRate, txRate, burstSize uint32) error {
	args := f.Called(rxMeter, txMeter, rxRate, txRate, burstSize)
	return args.Error(0)
}

func (f *fakeDataPath) DeleteMetersOnly(rxMeter, txMeter uint32) error {
	args := f.Called(rxMeter, txMeter)
	return args.Error(0)
}

func (f *fakeDataPath) DataPathCount() uint32 {
	args := f.Called()
	return args.Get(0).(uint32)
}

func (f *fakeDataPath) DataPathStats(rxCookieID, txCookieID uint64) (uint64, uint64, uint64, uint64, error) {
	args := f.Called(rxCookieID, txCookieID)
	return args.Get(0).(uint64), args.Get(1).(uint64), args.Get(2).(uint64), args.Get(3).(uint64), args.Error(4)
}

func (f *fakeDataPath) Status() datapath.Status {
	args := f.Called()
	return args.Get(0).(datapath.Status)
}

func newTestStore(t *testing.T) *store.Store {
	t.Helper()

	sharedStoreOnce.Do(func() {
		dir, err := os.MkdirTemp("", "pcrf-session-test")
		assert.NoError(t, err)

		s, err := store.NewStore(filepath.Join(dir, "pcrf_test.db"))
		assert.NoError(t, err)
		sharedStore = s
	})

	return sharedStore
}

func newSubscriberWithSession(t *testing.T, s *store.Store, d *fakeDataPath, imsi string, serviceOn bool) (*sessionManager, *store.Subscriber) {
	t.Helper()

	policy := &api.Policy{
		Uuid:      uuid.NewV4(),
		Ulbr:      1000,
		Dlbr:      2000,
		Data:      1_000_000,
		Consumed:  0,
		Burst:     100,
		StartTime: time.Now().Add(-time.Hour).Unix(),
		EndTime:   time.Now().Add(time.Hour).Unix(),
	}

	rerouteIP := "192.168.9." + imsi[len(imsi)-2:]

	sub, err := s.CreateSubscriber(imsi, policy, &rerouteIP, nil)
	assert.NoError(t, err)

	if !serviceOn {
		assert.NoError(t, s.UpdateSubscriberServiceStatus(sub, false))
		sub, err = s.GetSubscriber(imsi)
		assert.NoError(t, err)
	}

	sm := &sessionManager{
		period: time.Hour,
		idle:   time.Hour,
		store:  s,
		d:      d,
		cache:  make(map[string]*sessionCache),
	}

	ns, rxf, txf, err := s.CreateSession(sub, "10.10.10.11", "node1")
	assert.NoError(t, err)

	if serviceOn {
		d.On("AddNewDataPath", ns.UeIpAddr, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	} else {
		d.On("AddMetersOnly", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	}

	err = sm.CreateSesssion(context.Background(), sub, ns, rxf, txf)
	assert.NoError(t, err)

	t.Cleanup(func() {
		if sc, ok := sm.cache[sub.Imsi]; ok && sc.cancel != nil {
			sc.cancel()
		}
	})

	return sm, sub
}

func TestSessionManager_CreateSesssion_BornPaused(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	sm, sub := newSubscriberWithSession(t, s, d, "999990000000001", false)

	sc, ok := sm.cache[sub.Imsi]
	assert.True(t, ok)
	assert.True(t, sc.paused)
	assert.Equal(t, store.FlowsPaused, sc.s.FlowState)

	d.AssertNotCalled(t, "AddNewDataPath", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	d.AssertExpectations(t)
}

func TestSessionManager_PauseSession(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	sm, sub := newSubscriberWithSession(t, s, d, "999990000000002", true)

	sc := sm.cache[sub.Imsi]
	sc.s.RxBytes = 500
	sc.s.TxBytes = 700

	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(500), uint64(0), uint64(700), uint64(0), nil).Once()
	d.On("DeleteFlowOnly", sc.s.UeIpAddr).Return(nil).Once()

	err := sm.PauseSession(context.Background(), sub)
	assert.NoError(t, err)

	assert.True(t, sc.paused)
	assert.Equal(t, uint64(500), sc.baseRxBytes)
	assert.Equal(t, uint64(700), sc.baseTxBytes)
	assert.Equal(t, store.FlowsPaused, sc.s.FlowState)

	d.AssertNotCalled(t, "DeleteDataPath", mock.Anything, mock.Anything, mock.Anything)
	d.AssertExpectations(t)

	// Idempotent: pausing an already-paused session is a no-op.
	err = sm.PauseSession(context.Background(), sub)
	assert.NoError(t, err)
	d.AssertExpectations(t)
}

func TestSessionManager_ResumeSession(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	sm, sub := newSubscriberWithSession(t, s, d, "999990000000003", true)

	sc := sm.cache[sub.Imsi]
	sc.s.RxBytes = 500
	sc.s.TxBytes = 700

	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(500), uint64(0), uint64(700), uint64(0), nil).Once()
	d.On("DeleteFlowOnly", sc.s.UeIpAddr).Return(nil).Once()
	assert.NoError(t, sm.PauseSession(context.Background(), sub))

	staleUpdatedAt := sc.s.UpdatedAt

	d.On("AddFlowOnly", sc.s.UeIpAddr, uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID), sc.rxCookie, sc.txCookie).Return(nil).Once()

	err := sm.ResumeSession(context.Background(), sub)
	assert.NoError(t, err)

	assert.False(t, sc.paused)
	assert.Equal(t, store.FlowsActive, sc.s.FlowState)
	assert.GreaterOrEqual(t, sc.s.UpdatedAt, staleUpdatedAt)

	d.AssertExpectations(t)

	// Idempotent: resuming an already-active session is a no-op.
	err = sm.ResumeSession(context.Background(), sub)
	assert.NoError(t, err)
	d.AssertExpectations(t)
}

func TestSessionManager_StoreStats_SkipsDataPathStatsWhilePaused(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	sm, sub := newSubscriberWithSession(t, s, d, "999990000000004", true)

	sc := sm.cache[sub.Imsi]
	sc.s.RxBytes = 111
	sc.s.TxBytes = 222

	// PauseSession itself does one legitimate stats flush before removing the
	// flow - that's the only DataPathStats call this test expects, ever.
	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(111), uint64(0), uint64(222), uint64(0), nil).Once()
	d.On("DeleteFlowOnly", sc.s.UeIpAddr).Return(nil).Once()
	assert.NoError(t, sm.PauseSession(context.Background(), sub))

	err := sm.storeStats(sub.Imsi, false)
	assert.NoError(t, err)

	d.AssertExpectations(t)
	assert.Equal(t, uint64(111), sc.s.RxBytes)
	assert.Equal(t, uint64(222), sc.s.TxBytes)
}

func TestSessionManager_StoreStats_AddsBaselineAfterResume(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	sm, sub := newSubscriberWithSession(t, s, d, "999990000000005", true)

	sc := sm.cache[sub.Imsi]
	sc.s.RxBytes = 1000
	sc.s.TxBytes = 2000

	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(1000), uint64(0), uint64(2000), uint64(0), nil).Once()
	d.On("DeleteFlowOnly", sc.s.UeIpAddr).Return(nil).Once()
	assert.NoError(t, sm.PauseSession(context.Background(), sub))
	assert.Equal(t, uint64(1000), sc.baseRxBytes)
	assert.Equal(t, uint64(2000), sc.baseTxBytes)

	d.On("AddFlowOnly", sc.s.UeIpAddr, mock.Anything, mock.Anything, sc.rxCookie, sc.txCookie).Return(nil).Once()
	assert.NoError(t, sm.ResumeSession(context.Background(), sub))

	// Switch resets its own counters to zero on a fresh flow install; storeStats
	// must add them on top of the frozen baseline, not replace it.
	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(50), uint64(0), uint64(75), uint64(0), nil).Once()

	err := sm.storeStats(sub.Imsi, false)
	assert.NoError(t, err)

	assert.Equal(t, uint64(1050), sc.s.RxBytes)
	assert.Equal(t, uint64(2075), sc.s.TxBytes)
}

func TestSessionManager_StoreStats_IdleTimeoutEndsActiveSession(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	sm, sub := newSubscriberWithSession(t, s, d, "999990000000007", true)

	sc := sm.cache[sub.Imsi]

	sc.s.UpdatedAt = uint64(time.Now().Add(-2 * sm.idle).Unix())

	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(0), uint64(0), uint64(0), uint64(0), nil).Twice()
	d.On("DeleteDataPath", sc.s.UeIpAddr, uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID)).Return(nil).Once()

	err := sm.storeStats(sub.Imsi, false)
	assert.Error(t, err, "an idle timeout must be surfaced as an error so the caller knows the session is gone")

	_, ok := sm.cache[sub.Imsi]
	assert.False(t, ok, "an idle-timed-out session must be removed from the cache")

	d.AssertExpectations(t)
}

func TestSessionManager_StoreStats_IdleTimeoutEndsPausedSession(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	sm, sub := newSubscriberWithSession(t, s, d, "999990000000008", true)

	sc := sm.cache[sub.Imsi]

	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(0), uint64(0), uint64(0), uint64(0), nil).Once()
	d.On("DeleteFlowOnly", sc.s.UeIpAddr).Return(nil).Once()
	assert.NoError(t, sm.PauseSession(context.Background(), sub))

	sc.s.UpdatedAt = uint64(time.Now().Add(-2 * sm.idle).Unix())

	d.On("DeleteMetersOnly", uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID)).Return(nil).Once()

	err := sm.storeStats(sub.Imsi, false)
	assert.Error(t, err, "an idle timeout must be surfaced as an error so the caller knows the session is gone")

	_, ok := sm.cache[sub.Imsi]
	assert.False(t, ok, "an idle-timed-out paused session must be removed from the cache")

	d.AssertNotCalled(t, "DeleteDataPath", mock.Anything, mock.Anything, mock.Anything)
	d.AssertExpectations(t)
}

func TestSessionManager_PauseSession_SessionEndedDuringFlush(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	sm, sub := newSubscriberWithSession(t, s, d, "999990000000009", true)

	sc := sm.cache[sub.Imsi]

	sc.s.UpdatedAt = uint64(time.Now().Add(-2 * sm.idle).Unix())

	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(0), uint64(0), uint64(0), uint64(0), nil).Twice()
	d.On("DeleteDataPath", sc.s.UeIpAddr, uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID)).Return(nil).Once()

	err := sm.PauseSession(context.Background(), sub)
	assert.NoError(t, err, "PauseSession must not error just because the session ended during its own flush")

	_, ok := sm.cache[sub.Imsi]
	assert.False(t, ok, "the session ended during flush must stay gone, not be resurrected")

	// Must not have tried to pause a session that's already gone.
	d.AssertNotCalled(t, "DeleteFlowOnly", mock.Anything)
	d.AssertExpectations(t)
}

func TestSessionManager_EndSession_PausedDeletesMetersOnly(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	sm, sub := newSubscriberWithSession(t, s, d, "999990000000006", true)

	sc := sm.cache[sub.Imsi]

	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(0), uint64(0), uint64(0), uint64(0), nil).Once()
	d.On("DeleteFlowOnly", sc.s.UeIpAddr).Return(nil).Once()
	assert.NoError(t, sm.PauseSession(context.Background(), sub))

	d.On("DeleteMetersOnly", uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID)).Return(nil).Once()

	err := sm.EndSession(context.Background(), sub)
	assert.NoError(t, err)

	d.AssertNotCalled(t, "DeleteDataPath", mock.Anything, mock.Anything, mock.Anything)
	d.AssertExpectations(t)

	_, ok := sm.cache[sub.Imsi]
	assert.False(t, ok)
}

func TestSessionManager_EndSession_ClearsCacheEvenWhenDatapathCleanupFails(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	sm, sub := newSubscriberWithSession(t, s, d, "999990000000010", true)

	sc := sm.cache[sub.Imsi]

	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(0), uint64(0), uint64(0), uint64(0), nil).Once()
	d.On("DeleteDataPath", sc.s.UeIpAddr, uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID)).Return(assert.AnError).Once()

	err := sm.EndSession(context.Background(), sub)
	assert.Error(t, err)

	_, ok := sm.cache[sub.Imsi]
	assert.False(t, ok, "the cache entry must be cleared even when datapath cleanup fails, or HasActiveSession keeps reporting a session that no longer exists in the store")
}

func TestSessionManager_ConcurrentTickerAndRequests_NoRace(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}

	policy := &api.Policy{
		Uuid:      uuid.NewV4(),
		Ulbr:      1000,
		Dlbr:      2000,
		Data:      1_000_000,
		Consumed:  0,
		Burst:     100,
		StartTime: time.Now().Add(-time.Hour).Unix(),
		EndTime:   time.Now().Add(time.Hour).Unix(),
	}

	imsi := "999990000000099"
	rerouteIP := "192.168.9.99"

	sub, err := s.CreateSubscriber(imsi, policy, &rerouteIP, nil)
	assert.NoError(t, err)

	sm := &sessionManager{
		period: 2 * time.Millisecond,
		idle:   time.Hour,
		store:  s,
		d:      d,
		cache:  make(map[string]*sessionCache),
	}

	ns, rxf, txf, err := s.CreateSession(sub, "10.10.10.55", "node1")
	assert.NoError(t, err)

	d.On("AddNewDataPath", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	d.On("DataPathStats", mock.Anything, mock.Anything).Return(uint64(0), uint64(0), uint64(0), uint64(0), nil)
	d.On("DeleteFlowOnly", mock.Anything).Return(nil)
	d.On("AddFlowOnly", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err = sm.CreateSesssion(context.Background(), sub, ns, rxf, txf)
	assert.NoError(t, err)

	t.Cleanup(func() {
		if sc, ok := sm.cache[sub.Imsi]; ok && sc.cancel != nil {
			sc.cancel()
		}
	})

	var wg sync.WaitGroup
	stop := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			sm.HasActiveSession(imsi)
			_ = sm.PauseSession(context.Background(), sub)
			_ = sm.ResumeSession(context.Background(), sub)
		}
	}()

	time.Sleep(150 * time.Millisecond)
	close(stop)
	wg.Wait()

	assert.True(t, sm.HasActiveSession(imsi))
}
