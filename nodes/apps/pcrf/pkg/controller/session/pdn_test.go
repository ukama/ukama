/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package session

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
)

func TestPauseFailureCleanupDoesNotDoubleCount(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}
	sm, sub := newSubscriberWithSession(t, s, d, "999990000000301", true)
	sc := sm.cache[sub.Imsi]
	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(
		uint64(20), uint64(0), uint64(80), uint64(0), nil).Once()
	d.On("DeleteFlowOnly", sc.s.UeIpAddr).
		Return(errors.New("temporary failure")).Once()
	require.Error(t, sm.PauseSession(context.Background(), sub))
	require.Contains(t, sm.cache, sub.Imsi)
	require.False(t, sm.HasActiveSession(sub.Imsi))
	require.Error(t, sm.ResumeSession(context.Background(), sub))
	require.Zero(t, sc.baseRxBytes)
	require.Zero(t, sc.baseTxBytes)
	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(
		uint64(30), uint64(0), uint64(120), uint64(0), nil).Once()
	d.On("DeleteDataPath", sc.s.UeIpAddr, mock.Anything,
		mock.Anything).Return(nil).Once()
	require.NoError(t, sm.EndSession(context.Background(), sub))
	u, err := s.GetUsageByImsi(sub.Imsi)
	require.NoError(t, err)
	require.Equal(t, uint64(150), u.Data)
	require.NotContains(t, sm.cache, sub.Imsi)
	d.AssertExpectations(t)
}

func TestPartialResumeRemainsUnknownUntilCleanup(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}
	sm, sub := newSubscriberWithSession(t, s, d, "999990000000302", false)
	sc := sm.cache[sub.Imsi]
	d.On("AddFlowOnly", sc.s.UeIpAddr, mock.Anything, mock.Anything,
		sc.rxCookie, sc.txCookie).Return(errors.New("RX install failed")).Once()
	require.Error(t, sm.ResumeSession(context.Background(), sub))
	got, err := s.GetActiveSessionByImsi(sub.Imsi)
	require.NoError(t, err)
	require.Equal(t, store.FlowStateUnknown, got.FlowState)
	require.Error(t, sm.ResumeSession(context.Background(), sub))
	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(
		uint64(0), uint64(0), uint64(0), uint64(0),
		errors.New("incomplete path")).Once()
	d.On("DeleteDataPath", sc.s.UeIpAddr, mock.Anything,
		mock.Anything).Return(nil).Once()
	require.NoError(t, sm.EndSession(context.Background(), sub))
	require.False(t, sm.HasActiveSession(sub.Imsi))
	d.AssertExpectations(t)
}

func TestRecoveryCleansPermissionsBeforeReleasingSession(t *testing.T) {
	s := newTestStore(t)
	// Prior fixtures have canceled their monitors. Close their DB records
	// so this test models one interrupted session on a restarted process.
	require.NoError(t, s.RecoverActiveSessions())
	d := &fakeDataPath{}
	sm, sub := newSubscriberWithSession(t, s, d, "999990000000303", true)
	sc := sm.cache[sub.Imsi]
	sc.s.TxBytes = 100
	require.NoError(t, s.UpdateSessionUsage(sc.s))
	sc.cancel()
	d.On("DeleteDataPath", sc.s.UeIpAddr, mock.Anything,
		mock.Anything).Return(errors.New("switch unavailable")).Once()
	require.Error(t, sm.recoverSessions())
	_, err := s.GetActiveSessionByImsi(sub.Imsi)
	require.NoError(t, err)
	u, err := s.GetUsageByImsi(sub.Imsi)
	require.NoError(t, err)
	require.Zero(t, u.Data)
	d.On("DeleteDataPath", sc.s.UeIpAddr, mock.Anything,
		mock.Anything).Return(nil).Once()
	require.NoError(t, sm.recoverSessions())
	require.NoError(t, sm.recoverSessions())
	_, err = s.GetActiveSessionByImsi(sub.Imsi)
	require.ErrorIs(t, err, store.ErrSessionNotFound)
	u, err = s.GetUsageByImsi(sub.Imsi)
	require.NoError(t, err)
	require.Equal(t, uint64(100), u.Data)
	d.AssertExpectations(t)
}

func TestResetCounterPreservesRecordedUsage(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}
	sm, sub := newSubscriberWithSession(t, s, d, "999990000000304", true)
	sc := sm.cache[sub.Imsi]
	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(
		uint64(20), uint64(0), uint64(80), uint64(0), nil).Once()
	require.NoError(t, sm.storeStats(sub.Imsi, false))
	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(
		uint64(0), uint64(0), uint64(0), uint64(0), nil).Once()
	require.Error(t, sm.storeStats(sub.Imsi, false))
	se, err := s.GetSessionByID(sc.s.ID)
	require.NoError(t, err)
	require.Equal(t, uint64(100), se.TotalBytes)
	d.AssertExpectations(t)
}
