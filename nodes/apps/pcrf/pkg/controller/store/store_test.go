/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/systems/common/uuid"
)

// NewStore registers a global sql driver name on every call, which panics if
// called more than once per process - so tests in this file that need a full
// Store must share a single instance.
var (
	sharedStoreOnce sync.Once
	sharedTestStore *Store
)

func newSharedTestStore(t *testing.T) *Store {
	t.Helper()

	sharedStoreOnce.Do(func() {
		dir, err := os.MkdirTemp("", "pcrf-store-test")
		assert.NoError(t, err)

		s, err := NewStore(filepath.Join(dir, "pcrf_test.db"))
		assert.NoError(t, err)
		sharedTestStore = s
	})

	return sharedTestStore
}

func TestStore_ServiceOnAndFlowState_RoundTrip(t *testing.T) {
	s := newSharedTestStore(t)

	policy := &api.Policy{
		Uuid:      uuid.NewV4(),
		Ulbr:      1000,
		Dlbr:      2000,
		Data:      1_000_000,
		Burst:     100,
		StartTime: time.Now().Add(-time.Hour).Unix(),
		EndTime:   time.Now().Add(time.Hour).Unix(),
	}
	ip := "192.168.9.200"

	sub, err := s.CreateSubscriber("999991000000001", policy, &ip, nil)
	assert.NoError(t, err)
	assert.True(t, sub.ServiceOn, "new subscribers must default to service-on")

	assert.NoError(t, s.UpdateSubscriberServiceStatus(sub, false))

	got, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.False(t, got.ServiceOn)

	gotByID, err := s.GetSubscriberByID(sub.ID)
	assert.NoError(t, err)
	assert.False(t, gotByID.ServiceOn)

	ns, _, _, err := s.CreateSession(got, "10.10.10.20", "node1")
	assert.NoError(t, err)
	assert.Equal(t, FlowsPaused, ns.FlowState, "session created for an off subscriber must be born paused")

	sess, err := s.GetSessionByID(ns.ID)
	assert.NoError(t, err)
	assert.Equal(t, FlowsPaused, sess.FlowState)

	assert.NoError(t, s.UpdateSessionFlowState(ns.ID, FlowsActive))

	sess, err = s.GetSessionByID(ns.ID)
	assert.NoError(t, err)
	assert.Equal(t, FlowsActive, sess.FlowState)
}

func TestStore_Migration_IdempotentOnRestart(t *testing.T) {
	s := newSharedTestStore(t)

	assert.NoError(t, s.CreateTables())
	assert.NoError(t, s.CreateTables())
}

func TestStore_Migration_UpgradesPreExistingSchema(t *testing.T) {
	dir, err := os.MkdirTemp("", "pcrf-store-migration-test")
	assert.NoError(t, err)

	dbPath := filepath.Join(dir, "old_schema.db")

	rawDB, err := sql.Open("sqlite3", dbPath)
	assert.NoError(t, err)

	_, err = rawDB.Exec(`
		CREATE TABLE IF NOT EXISTS subscribers (
			id INTEGER PRIMARY KEY,
			imsi TEXT UNIQUE,
			policy_id BLOB CHECK(length(policy_id) = 16),
			reroute_id INTEGER
		);
	`)
	assert.NoError(t, err)

	_, err = rawDB.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			node_id TEXT,
			subscriber_id INTEGER,
			policy_id BLOB CHECK(length(policy_id) = 16),
			apnname TEXT,
			ueipaddr TEXT,
			starttime INTEGER,
			endtime INTEGER,
			txbytes INTEGER,
			rxbytes INTEGER,
			totalbytes INTEGER,
			txmeter_id INTEGER,
			rxmeter_id INTEGER,
			state INTEGER,
			sync INTEGER,
			updatedat INTEGER
		);
	`)
	assert.NoError(t, err)
	assert.NoError(t, rawDB.Close())

	db2, err := sql.Open("sqlite3", dbPath)
	assert.NoError(t, err)

	upgraded := &Store{db: db2}
	assert.NoError(t, upgraded.CreateTables())

	// The new columns must now be usable.
	_, err = db2.Exec(`UPDATE subscribers SET service_on = 1 WHERE id = -1`)
	assert.NoError(t, err)

	_, err = db2.Exec(`UPDATE sessions SET flowstate = 1 WHERE id = -1`)
	assert.NoError(t, err)

	assert.NoError(t, upgraded.CreateTables())
}

func TestStore_UpdateSubscriber_RejectsOlderPolicyRollover(t *testing.T) {
	s := newSharedTestStore(t)

	now := time.Now()
	current := &api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      1_000_000,
		Burst:     100,
		StartTime: now.Unix(),
		EndTime:   now.Add(time.Hour).Unix(),
	}

	sub, err := s.CreateSubscriber("999991000000002", current, nil, nil)
	assert.NoError(t, err)

	newer := &api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      2_000_000,
		Burst:     100,
		StartTime: now.Add(time.Minute).Unix(),
		EndTime:   now.Add(2 * time.Hour).Unix(),
	}
	_, err = s.UpdateSubscriber(sub.Imsi, newer)
	assert.NoError(t, err)

	got, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.Equal(t, newer.Uuid, got.PolicyID.ID, "rollover to a genuinely newer policy must apply")

	// A delayed push for the original (older) rollover arrives after the newer
	// one already landed - it must not clobber the newer policy.
	stale := &api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      500_000,
		Burst:     100,
		StartTime: now.Add(-time.Minute).Unix(),
		EndTime:   now.Add(30 * time.Minute).Unix(),
	}
	_, err = s.UpdateSubscriber(sub.Imsi, stale)
	assert.NoError(t, err)

	got, err = s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.Equal(t, newer.Uuid, got.PolicyID.ID, "a stale/delayed policy push must not overwrite a newer policy")
}

func TestStore_UpdateSubscriber_RejectsOlderConsumedForSamePolicy(t *testing.T) {
	s := newSharedTestStore(t)

	now := time.Now()
	policy := &api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      1_000_000,
		Burst:     100,
		StartTime: now.Unix(),
		EndTime:   now.Add(time.Hour).Unix(),
	}

	sub, err := s.CreateSubscriber("999991000000003", policy, nil, nil)
	assert.NoError(t, err)

	ahead := &api.Policy{Uuid: policy.Uuid, Data: policy.Data, Consumed: 500_000, StartTime: policy.StartTime, EndTime: policy.EndTime}
	_, err = s.UpdateSubscriber(sub.Imsi, ahead)
	assert.NoError(t, err)

	got, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.Equal(t, uint64(500_000), got.PolicyID.Consumed)

	// A redelivered/retried update for the same policy carrying an older
	// consumed value must not regress consumed bytes backwards.
	behind := &api.Policy{Uuid: policy.Uuid, Data: policy.Data, Consumed: 100_000, StartTime: policy.StartTime, EndTime: policy.EndTime}
	_, err = s.UpdateSubscriber(sub.Imsi, behind)
	assert.NoError(t, err)

	got, err = s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.Equal(t, uint64(500_000), got.PolicyID.Consumed, "an out-of-order consumed-data update must not regress consumed bytes")
}

func TestStore_EndSession_DoesNotCountUsageFromSupersededPolicy(t *testing.T) {
	s := newSharedTestStore(t)

	now := time.Now()
	policyA := &api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      1_000_000,
		Burst:     100,
		StartTime: now.Unix(),
		EndTime:   now.Add(time.Hour).Unix(),
	}
	ip := "192.168.9.201"

	sub, err := s.CreateSubscriber("999991000000004", policyA, &ip, nil)
	assert.NoError(t, err)

	session, _, _, err := s.CreateSession(sub, "10.10.10.30", "node1")
	assert.NoError(t, err)

	policyB := &api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      500_000,
		Burst:     100,
		StartTime: now.Add(time.Minute).Unix(),
		EndTime:   now.Add(2 * time.Hour).Unix(),
	}
	_, err = s.UpdateSubscriber(sub.Imsi, policyB)
	assert.NoError(t, err)

	session.TxBytes = 300_000
	session.RxBytes = 300_000
	assert.NoError(t, s.EndSession(session))

	usage, err := s.GetUsageByImsi(sub.Imsi)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), usage.Data, "a session ending under a policy the subscriber has since rolled over from must not count against the new policy's usage")

	rolled, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)

	fresh, _, _, err := s.CreateSession(rolled, "10.10.10.31", "node1")
	assert.NoError(t, err)

	fresh.TxBytes = 50_000
	fresh.RxBytes = 50_000
	assert.NoError(t, s.EndSession(fresh))

	usage, err = s.GetUsageByImsi(sub.Imsi)
	assert.NoError(t, err)
	assert.Equal(t, uint64(100_000), usage.Data, "a session ending under the subscriber's current policy must still count normally")
}

func TestStore_CreateSession_UsesFreshServiceStatus(t *testing.T) {
	s := newSharedTestStore(t)

	policy := &api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      1_000_000,
		Burst:     100,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Add(time.Hour).Unix(),
	}
	ip := "192.168.9.202"

	sub, err := s.CreateSubscriber("999991000000005", policy, &ip, nil)
	assert.NoError(t, err)

	stale, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)
	assert.True(t, stale.ServiceOn)

	assert.NoError(t, s.UpdateSubscriberServiceStatus(sub, false))

	ns, _, _, err := s.CreateSession(stale, "10.10.10.40", "node1")
	assert.NoError(t, err)
	assert.Equal(t, FlowsPaused, ns.FlowState,
		"CreateSession must use the subscriber's current service state, not a caller-held stale snapshot")
}

func TestStore_CreateSession_UsesFreshPolicyForCapCheck(t *testing.T) {
	s := newSharedTestStore(t)

	now := time.Now()
	policyA := &api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      1_000_000,
		Burst:     100,
		StartTime: now.Unix(),
		EndTime:   now.Add(time.Hour).Unix(),
	}
	ip := "192.168.9.203"

	sub, err := s.CreateSubscriber("999991000000006", policyA, &ip, nil)
	assert.NoError(t, err)

	stale, err := s.GetSubscriber(sub.Imsi)
	assert.NoError(t, err)

	policyB := &api.Policy{
		Uuid:      uuid.NewV4(),
		Data:      500_000,
		Burst:     100,
		StartTime: now.Add(time.Minute).Unix(),
		EndTime:   now.Add(2 * time.Hour).Unix(),
	}
	_, err = s.UpdateSubscriber(sub.Imsi, policyB)
	assert.NoError(t, err)

	ns, _, _, err := s.CreateSession(stale, "10.10.10.41", "node1")
	assert.NoError(t, err)
	assert.Equal(t, policyB.Uuid, ns.PolicyID.ID,
		"session must be created under the subscriber's current policy, not the caller's stale snapshot")
	assert.Equal(t, FlowsActive, ns.FlowState, "fresh policy with zero usage must not be born paused")
}
