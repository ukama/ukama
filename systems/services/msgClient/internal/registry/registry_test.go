/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package registry

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fakeTopo struct {
	mu      sync.Mutex
	ensured [][]string
	deleted []string
	owned   []string
}

func (f *fakeTopo) EnsureQueue(queue string, routes []string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensured = append(f.ensured, routes)
	return nil
}

func (f *fakeTopo) DeleteQueue(queue string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.deleted = append(f.deleted, queue)
	return nil
}

func (f *fakeTopo) SetMaxLength(queue string, max int64) error { return nil }

func (f *fakeTopo) OwnedQueues() ([]string, error) { return f.owned, nil }

type fakeListener struct {
	stopped       atomic.Bool
	onUnreachable func()
}

func (l *fakeListener) Stop() { l.stopped.Store(true) }

type harness struct {
	r         *Registry
	topo      *fakeTopo
	listeners []*fakeListener
	clock     time.Time
}

func newHarness(t *testing.T) *harness {
	h := &harness{topo: &fakeTopo{}, clock: time.Unix(1000, 0)}
	h.r = New(h.topo, func(service, uri string, onUnreachable func()) (Listener, error) {
		l := &fakeListener{onUnreachable: onUnreachable}
		h.listeners = append(h.listeners, l)
		return l, nil
	}, Config{RefreshInterval: 2 * time.Second, MissedRefreshes: 3, FlushAfter: time.Hour})
	h.r.now = func() time.Time { return h.clock }
	return h
}

func (h *harness) register(t *testing.T, routes ...string) string {
	t.Helper()
	id, err := h.r.Register("node", "registry", "1", "node:9090", routes)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	return id
}

func (h *harness) active() int {
	n := 0
	for _, l := range h.listeners {
		if !l.stopped.Load() {
			n++
		}
	}
	return n
}

func TestRegisterUnchangedMakesNoBrokerCalls(t *testing.T) {
	h := newHarness(t)
	id1 := h.register(t, "b", "a")
	id2 := h.register(t, "a", "b")

	if id1 != id2 {
		t.Fatalf("uuid changed between registrations")
	}
	if len(h.topo.ensured) != 1 || len(h.listeners) != 1 {
		t.Fatalf("expected 1 queue setup and 1 listener, got %d and %d", len(h.topo.ensured), len(h.listeners))
	}
}

func TestRegisterReconcilesChangedRoutes(t *testing.T) {
	h := newHarness(t)
	h.register(t, "a", "b")
	h.register(t, "a")

	if len(h.topo.ensured) != 2 || len(h.topo.ensured[1]) != 1 || h.topo.ensured[1][0] != "a" {
		t.Fatalf("expected second setup with routes [a], got %v", h.topo.ensured)
	}
	if h.active() != 1 {
		t.Fatalf("expected the listener to keep running, got %d active", h.active())
	}
}

func TestNoRoutesMeansNoQueueOrListener(t *testing.T) {
	h := newHarness(t)
	h.register(t)

	if len(h.topo.ensured) != 0 || len(h.listeners) != 0 {
		t.Fatalf("expected no queue and no listener")
	}
}

func TestLapsedRegistrationPausesUntilNextRegistration(t *testing.T) {
	h := newHarness(t)
	h.register(t, "a")

	h.clock = h.clock.Add(7 * time.Second)
	h.r.sweep()
	waitFor(t, func() bool { return h.active() == 0 })

	h.register(t, "a")
	if h.active() != 1 || len(h.listeners) != 2 {
		t.Fatalf("expected a fresh listener after registering again")
	}
}

func TestUnreachableServicePausesUntilNextRegistration(t *testing.T) {
	h := newHarness(t)
	h.register(t, "a")

	h.listeners[0].onUnreachable()
	if h.active() != 0 {
		t.Fatalf("expected delivery paused")
	}

	h.register(t, "a")
	if h.active() != 1 {
		t.Fatalf("expected delivery resumed")
	}
}

func TestStaleCallbackDoesNotStopNewListener(t *testing.T) {
	h := newHarness(t)
	h.register(t, "a")
	h.listeners[0].onUnreachable()
	h.register(t, "a")

	h.listeners[0].onUnreachable()
	if h.active() != 1 {
		t.Fatalf("old listener's callback stopped the new listener")
	}
}

func TestFlushDeletesQueueOnlyAfterFlushTime(t *testing.T) {
	h := newHarness(t)
	h.register(t, "a")

	h.clock = h.clock.Add(59 * time.Minute)
	h.r.sweep()
	if len(h.topo.deleted) != 0 {
		t.Fatalf("queue flushed too early")
	}

	h.clock = h.clock.Add(2 * time.Minute)
	h.r.sweep()
	if len(h.topo.deleted) != 1 {
		t.Fatalf("expected queue flushed after flush time")
	}

	h.register(t, "a")
	if len(h.topo.ensured) != 2 {
		t.Fatalf("expected queue set up again after flush")
	}
}

func TestQueuesFoundAtStartupFlushFromStartTime(t *testing.T) {
	h := newHarness(t)
	h.topo.owned = []string{"retired"}
	if err := h.r.LoadOwned(); err != nil {
		t.Fatal(err)
	}

	h.clock = h.clock.Add(30 * time.Minute)
	h.r.sweep()
	if len(h.topo.deleted) != 0 {
		t.Fatalf("queue flushed before flush time since start")
	}

	h.clock = h.clock.Add(31 * time.Minute)
	h.r.sweep()
	if len(h.topo.deleted) != 1 || h.topo.deleted[0] != "retired" {
		t.Fatalf("expected retired queue flushed, got %v", h.topo.deleted)
	}
}

func TestUuidIsStableAcrossRestarts(t *testing.T) {
	a := newHarness(t).register(t, "a")
	b := newHarness(t).register(t, "a")
	if a != b {
		t.Fatalf("uuid differs across registries: %s vs %s", a, b)
	}
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("condition not met")
		}
		time.Sleep(time.Millisecond)
	}
}
