/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package controller

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/api"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller/store"
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

func TestSafeRun_RecoversPanicWithoutPropagating(t *testing.T) {
	assert.NotPanics(t, func() {
		safeRun("boom", func() {
			panic("simulated panic")
		})
	})
}

func TestSafeRun_RunsFnWhenNoPanic(t *testing.T) {
	ran := false

	safeRun("ok", func() {
		ran = true
	})

	assert.True(t, ran)
}

func TestSafeRun_SecondCallStillRunsAfterEarlierPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		safeRun("first", func() {
			panic("first panic")
		})
	})

	ran := false
	safeRun("second", func() {
		ran = true
	})

	assert.True(t, ran, "a panic in one safeRun call must not prevent a later call from running")
}
