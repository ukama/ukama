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

	"github.com/stretchr/testify/require"
)

func pdnTestStore(t *testing.T, path string) *Store {
	t.Helper()
	db, err := sql.Open("sqlite3", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	s := &Store{db: db}
	require.NoError(t, s.CreateTables())
	return s
}

func TestPDNReservationsAndTombstonesSurviveReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pdn.db")
	s := pdnTestStore(t, path)
	require.NoError(t, s.BeginEPC("first", ""))
	p := &PDN{InstanceID: "first", ID: "pdn-a", Imsi: "00101",
		IP: "192.168.8.2", ApnName: "internet", State: PDNLive}
	require.NoError(t, s.InsertPDN(p))
	conflict := *p
	conflict.ID = "pdn-b"
	require.Error(t, s.InsertPDN(&conflict))
	conflict.Imsi = "00102"
	require.Error(t, s.InsertPDN(&conflict))
	require.NoError(t, s.RetireEPCs())
	require.NoError(t, s.db.Close())
	s = pdnTestStore(t, path)
	got, err := s.GetLivePDN(p.Imsi)
	require.NoError(t, err)
	require.Equal(t, PDNClosing, got.State)
	require.Equal(t, "internet", got.ApnName)
	require.ErrorIs(t, s.BeginEPC("first", ""), ErrInstanceRetired)
	require.Error(t, s.InsertPDN(&conflict))
	require.NoError(t, s.SetPDNState(got, PDNClosed))
	require.NoError(t, s.InsertPDN(&conflict))
	tombstone, err := s.GetPDN("first", "pdn-a")
	require.NoError(t, err)
	require.Equal(t, PDNClosed, tombstone.State)
	require.Error(t, s.InsertPDN(p))
}

func TestEPCReplacementIsConditionalAndRetryable(t *testing.T) {
	s := pdnTestStore(t, filepath.Join(t.TempDir(), "pdn.db"))
	require.NoError(t, s.BeginEPC("first", ""))
	require.NoError(t, s.BeginEPC("first", ""))
	require.ErrorIs(t, s.BeginEPC("second", "wrong"), ErrInstanceConflict)
	require.NoError(t, s.BeginEPC("second", "first"))
	require.NoError(t, s.BeginEPC("second", "first"))
	require.ErrorIs(t, s.BeginEPC("first", "second"), ErrInstanceRetired)
	require.ErrorIs(t, s.BeginEPC("late", "first"), ErrInstanceConflict)
	owner, err := s.CurrentEPC()
	require.NoError(t, err)
	require.Equal(t, "second", owner)
}
