/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package store

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/systems/common/uuid"
)

func TestStoreUsesSingleConnection(t *testing.T) {
	s := newSharedTestStore(t)
	require.Equal(t, 1, s.db.Stats().MaxOpenConnections)
}

func TestSessionListsWithSingleConnection(t *testing.T) {
	s := openRestartStore(t, filepath.Join(t.TempDir(), "pcrf.db"))
	s.db.SetMaxOpenConns(1)

	sessions := make([]*Session, 3)
	for i := range sessions {
		policy := &api.Policy{
			Uuid:      uuid.NewV4(),
			Data:      1024 * 1024,
			StartTime: time.Now().Add(-time.Minute).Unix(),
			EndTime:   time.Now().Add(time.Hour).Unix(),
		}
		sub, err := s.CreateSubscriber(fmt.Sprintf("00101088987315%d", i), policy, nil, nil)
		require.NoError(t, err)
		sessions[i], _, _, err = s.CreateSession(sub, fmt.Sprintf("192.168.8.%d", i+2), "tower")
		require.NoError(t, err)
	}

	sessions[1].TxBytes = 1234
	require.NoError(t, s.EndSession(sessions[1]))
	require.NoError(t, s.UpdateSessionState(sessions[2].ID, SessionTerminated))

	tests := []struct {
		name string
		read func() ([]Session, error)
		want *Session
	}{
		{"by_imsi", func() ([]Session, error) {
			return s.GetSessionsByImsi(sessions[0].SubscriberID.Imsi)
		}, sessions[0]},
		{"unsynced", func() ([]Session, error) {
			return s.GetUnsyncSessionsByImsiAfterTime(sessions[1].SubscriberID.Imsi, 0)
		}, sessions[1]},
		{"active", s.GetAllActiveSessions, sessions[0]},
		{"unpublished", s.GetAllNonPublishedSessions, sessions[1]},
		{"terminated", s.GetAllNonPublishedTerminatedSessions, sessions[2]},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.read()
			require.NoError(t, err)
			require.Len(t, got, 1)
			require.Equal(t, test.want.ID, got[0].ID)
			require.Equal(t, test.want.PolicyID.ID, got[0].PolicyID.ID)
			require.Equal(t, test.want.SubscriberID.Imsi, got[0].SubscriberID.Imsi)
			require.Equal(t, test.want.TxMeterID.ID, got[0].TxMeterID.ID)
			require.Equal(t, test.want.RxMeterID.ID, got[0].RxMeterID.ID)
			require.Equal(t, test.want.TxBytes, got[0].TxBytes)
		})
	}

	// Exercise publisher reads alongside monitor writes and service cleanup.
	operations := []func() error{
		func() error {
			_, err := s.GetAllActiveSessions()
			return err
		},
		func() error {
			_, err := s.GetAllNonPublishedSessions()
			return err
		},
		func() error { return s.UpdateSessionUsage(sessions[0]) },
		func() error { return s.ClosePDNs() },
		func() error {
			_, err := s.GetLivePDNs()
			return err
		},
	}
	results := make(chan error, len(operations))
	for _, operation := range operations {
		go func(run func() error) {
			for i := 0; i < 100; i++ {
				if err := run(); err != nil {
					results <- err
					return
				}
			}
			results <- nil
		}(operation)
	}
	for range operations {
		require.NoError(t, <-results)
	}
}
