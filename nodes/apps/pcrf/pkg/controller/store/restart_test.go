/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package store

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/systems/common/uuid"
)

func openRestartStore(t *testing.T, path string) *Store {
	t.Helper()
	db, err := sql.Open("sqlite3", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	s := &Store{db: db}
	require.NoError(t, s.CreateTables())
	return s
}

func restartSubscriber(t *testing.T, s *Store) *Subscriber {
	t.Helper()
	p := &api.Policy{
		Uuid: uuid.NewV4(), Data: 5 * 1024 * 1024,
		StartTime: time.Now().Add(-time.Minute).Unix(),
		EndTime: time.Now().Add(time.Hour).Unix(),
	}
	route := "10.10.10.11"
	sub, err := s.CreateSubscriber("001010889873152", p, &route, nil)
	require.NoError(t, err)
	return sub
}

func TestRecoverActiveSessions_PreservesAllowanceAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pcrf.db")
	s := openRestartStore(t, path)
	sub := restartSubscriber(t, s)
	before, _, _, err := s.CreateSession(sub, "192.168.8.2", "tower")
	require.NoError(t, err)
	before.TxBytes, before.RxBytes = 2194938, 45925
	require.NoError(t, s.UpdateSessionUsage(before))
	require.NoError(t, s.db.Close())

	s = openRestartStore(t, path)
	require.NoError(t, s.RecoverActiveSessions())
	require.NoError(t, s.RecoverActiveSessions())
	u, err := s.GetUsageByImsi(sub.Imsi)
	require.NoError(t, err)
	require.Equal(t, uint64(2240863), u.Data)
	_, err = s.GetActiveSessionByImsi(sub.Imsi)
	require.ErrorIs(t, err, ErrSessionNotFound)
	closed, err := s.GetSessionByID(before.ID)
	require.NoError(t, err)
	require.Equal(t, SessionCompleted, closed.State)
	require.Equal(t, FlowsPaused, closed.FlowState)
	require.Equal(t, SessionSyncReady, closed.Sync)

	after, _, _, err := s.CreateSession(sub, "192.168.8.2", "tower")
	require.NoError(t, err)
	after.TxBytes = 3 * 1024 * 1024
	require.NoError(t, s.EndSession(after))
	require.NoError(t, s.EndSession(after))
	u, err = s.GetUsageByImsi(sub.Imsi)
	require.NoError(t, err)
	require.Equal(t, uint64(2240863+3*1024*1024), u.Data)
	require.Error(t, s.ValidateDataCapLimits(sub.Imsi, &sub.PolicyID))
}

func TestEndSession_RollsBackStateWhenUsageWriteFails(t *testing.T) {
	s := openRestartStore(t, filepath.Join(t.TempDir(), "pcrf.db"))
	sub := restartSubscriber(t, s)
	session, _, _, err := s.CreateSession(sub, "192.168.8.2", "tower")
	require.NoError(t, err)
	session.TxBytes = 200
	_, err = s.db.Exec(`CREATE TRIGGER fail_usage BEFORE UPDATE ON usages
		BEGIN SELECT RAISE(ABORT, 'injected usage failure'); END;`)
	require.NoError(t, err)
	require.Error(t, s.EndSession(session))
	got, err := s.GetSessionByID(session.ID)
	require.NoError(t, err)
	require.Equal(t, SessionActive, got.State)
	require.Zero(t, got.EndTime)
	u, err := s.GetUsageByImsi(sub.Imsi)
	require.NoError(t, err)
	require.Zero(t, u.Data)

	_, err = s.db.Exec(`DROP TRIGGER fail_usage`)
	require.NoError(t, err)
	require.NoError(t, s.EndSession(session))
	require.NoError(t, s.EndSession(session))
	u, err = s.GetUsageByImsi(sub.Imsi)
	require.NoError(t, err)
	require.Equal(t, uint64(200), u.Data)
}

func TestRecoverActiveSessions_DoesNotChargeSupersededPolicy(t *testing.T) {
	s := openRestartStore(t, filepath.Join(t.TempDir(), "pcrf.db"))
	sub := restartSubscriber(t, s)
	session, _, _, err := s.CreateSession(sub, "192.168.8.2", "tower")
	require.NoError(t, err)
	session.TxBytes = 200
	require.NoError(t, s.UpdateSessionUsage(session))
	p := &api.Policy{
		Uuid: uuid.NewV4(), Data: 1000,
		StartTime: time.Now().Unix(), EndTime: time.Now().Add(time.Hour).Unix(),
	}
	_, err = s.UpdateSubscriber(sub.Imsi, p)
	require.NoError(t, err)
	require.NoError(t, s.RecoverActiveSessions())
	u, err := s.GetUsageByImsi(sub.Imsi)
	require.NoError(t, err)
	require.Zero(t, u.Data)
	closed, err := s.GetSessionByID(session.ID)
	require.NoError(t, err)
	require.Equal(t, uint64(200), closed.TotalBytes)
	require.Equal(t, SessionSyncReady, closed.Sync)
}
