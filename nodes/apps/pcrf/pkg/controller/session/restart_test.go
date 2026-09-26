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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
	"github.com/ukama/ukama/systems/common/uuid"
)

func TestZeroBandwidth_CreatePauseResumeAndEnd(t *testing.T) {
	for i, rates := range [][2]uint64{{0, 0}, {0, 2000}, {1000, 0}} {
		t.Run(fmt.Sprint(rates), func(t *testing.T) {
			s := newTestStore(t)
			d := &fakeDataPath{}
			imsi := fmt.Sprintf("9999900000002%02d", i)
			route := fmt.Sprintf("192.168.9.%d", 120+i)
			p := &api.Policy{
				Uuid: uuid.NewV4(), Data: 1000000, Ulbr: rates[0], Dlbr: rates[1],
				StartTime: time.Now().Add(-time.Minute).Unix(), EndTime: time.Now().Add(time.Hour).Unix(),
			}
			sub, err := s.CreateSubscriber(imsi, p, &route, nil)
			require.NoError(t, err)
			ns, rx, tx, err := s.CreateSession(sub, "192.168.8.2", "tower")
			require.NoError(t, err)
			rxID, txID := uint32(ns.RxMeterID.ID), uint32(ns.TxMeterID.ID)
			if rates[0] == 0 {
				rxID = 0
			}
			if rates[1] == 0 {
				txID = 0
			}
			sm := &sessionManager{store: s, d: d, period: time.Hour, idle: time.Hour, cache: make(map[string]*sessionCache)}
			d.On("AddNewDataPath", ns.UeIpAddr, rxID, txID, uint32(rates[0]), uint32(rates[1]), uint32(0), rx.Cookie, tx.Cookie).Return(nil).Once()
			require.NoError(t, sm.CreateSesssion(context.Background(), sub, ns, rx, tx))
			t.Cleanup(func() { sm.cacheCleanupForTest(sub.Imsi) })
			d.On("DataPathStats", rx.Cookie, tx.Cookie).Return(uint64(20), uint64(1), uint64(80), uint64(1), nil).Once()
			d.On("DeleteFlowOnly", ns.UeIpAddr).Return(nil).Once()
			require.NoError(t, sm.PauseSession(context.Background(), sub))
			d.On("AddFlowOnly", ns.UeIpAddr, rxID, txID, rx.Cookie, tx.Cookie).Return(nil).Once()
			require.NoError(t, sm.ResumeSession(context.Background(), sub))
			d.On("DataPathStats", rx.Cookie, tx.Cookie).Return(uint64(10), uint64(1), uint64(40), uint64(1), nil).Once()
			d.On("DeleteDataPath", ns.UeIpAddr, rxID, txID).Return(nil).Once()
			require.NoError(t, sm.EndSession(context.Background(), sub))
			u, err := s.GetUsageByImsi(imsi)
			require.NoError(t, err)
			require.Equal(t, uint64(150), u.Data)
			d.AssertExpectations(t)
		})
	}
}

func (s *sessionManager) cacheCleanupForTest(imsi string) {
	if sc, ok := s.cache[imsi]; ok && sc.cancel != nil {
		sc.cancel()
	}
}

func TestEndSession_FinalStatsUnavailablePreservesRecordedCounters(t *testing.T) {
	s := newTestStore(t)
	d := &fakeDataPath{}
	sm, sub := newSubscriberWithSession(t, s, d, "999990000000220", true)
	sc := sm.cache[sub.Imsi]
	sc.s.TxBytes, sc.s.RxBytes = 2194938, 45925
	d.On("DataPathStats", sc.rxCookie, sc.txCookie).Return(uint64(0), uint64(0), uint64(0), uint64(0), fmt.Errorf("bridge unavailable")).Once()
	d.On("DeleteDataPath", sc.s.UeIpAddr, uint32(sc.s.RxMeterID.ID), uint32(sc.s.TxMeterID.ID)).Return(nil).Once()
	require.NoError(t, sm.EndSession(context.Background(), sub))
	u, err := s.GetUsageByImsi(sub.Imsi)
	require.NoError(t, err)
	require.Equal(t, uint64(2240863), u.Data)
	_, err = s.GetActiveSessionByImsi(sub.Imsi)
	require.ErrorIs(t, err, store.ErrSessionNotFound)
	d.AssertExpectations(t)
}
