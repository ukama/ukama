/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Copyright (c) 2026-present, Ukama Inc.
 */
package server

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/site/mocks"
	"github.com/ukama/ukama/systems/registry/site/pkg/db"
)

type fakeProvisionClient struct {
	call func(context.Context, provisionNode) (provisionResult, error)
}

func (c fakeProvisionClient) Reconcile(ctx context.Context, n provisionNode) (provisionResult, error) {
	return c.call(ctx, n)
}

type memoryProvisions struct {
	mu    sync.Mutex
	op    db.SiteProvision
	begin func(*db.Site, []string) (*db.SiteProvision, error)
}

func (m *memoryProvisions) Create(_ context.Context, site *db.Site, nodes []string) (*db.SiteProvision, error) {
	return m.begin(site, nodes)
}
func (m *memoryProvisions) Get(_ context.Context, _ string) (*db.SiteProvision, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	op := m.op
	return &op, nil
}
func (m *memoryProvisions) Save(_ context.Context, op *db.SiteProvision) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.op = *op
	return nil
}
func (m *memoryProvisions) Pending(context.Context) ([]db.SiteProvision, error) { return nil, nil }
func (m *memoryProvisions) Run(context.Context, string, func(context.Context, *db.SiteProvision) error) error {
	return nil
}

func testProvision() *db.SiteProvision {
	return &db.SiteProvision{ID: uuid.NewV4().String(), Site: db.Site{NetworkId: uuid.NewV4()}, Nodes: []string{"tower", "amplifier", "controller"}, Phase: "configuring", Attempt: 1, Deadline: time.Now().Add(time.Second)}
}

func TestProvisionWaitsForAllThree(t *testing.T) {
	started := make(chan string, 3)
	release := make(chan struct{})
	server := &SiteServer{provisionClient: fakeProvisionClient{call: func(ctx context.Context, n provisionNode) (provisionResult, error) {
		started <- n.NodeID
		select {
		case <-release:
			return provisionResult{Completed: true}, nil
		case <-ctx.Done():
			return provisionResult{}, ctx.Err()
		}
	}}}
	op := testProvision()
	done := make(chan error, 1)
	go func() { done <- server.waitNodes(context.Background(), op, "configure", op.Deadline) }()
	for i := 0; i < 3; i++ {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("nodes did not start concurrently")
		}
	}
	select {
	case <-done:
		t.Fatal("continued before all nodes completed")
	default:
	}
	close(release)
	require.NoError(t, <-done)
}

func TestProvisionSharedDeadline(t *testing.T) {
	server := &SiteServer{provisionClient: fakeProvisionClient{call: func(ctx context.Context, n provisionNode) (provisionResult, error) {
		if n.NodeID != "controller" {
			return provisionResult{Completed: true}, nil
		}
		<-ctx.Done()
		return provisionResult{}, ctx.Err()
	}}}
	start := time.Now()
	err := server.waitNodes(context.Background(), testProvision(), "configure", start.Add(30*time.Millisecond))
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, time.Since(start), 500*time.Millisecond)
}

func TestProvisionThreeAttemptsAndCleanup(t *testing.T) {
	store := &memoryProvisions{}
	var mu sync.Mutex
	cleared := make(map[string]int)
	server := &SiteServer{provisions: store, provisionClient: fakeProvisionClient{call: func(_ context.Context, n provisionNode) (provisionResult, error) {
		mu.Lock()
		defer mu.Unlock()
		if n.Action == "cancel" {
			cleared[n.NodeID]++
			return provisionResult{Cancelled: true, Cleared: true}, nil
		}
		if n.Action == "release" {
			return provisionResult{Cleared: true}, nil
		}
		t.Errorf("unexpected action %s", n.Action)
		return provisionResult{}, errors.New("unexpected action")
	}}}
	op := testProvision()
	for attempt := 1; attempt <= 3; attempt++ {
		op.Deadline = time.Now().Add(-time.Second)
		require.NoError(t, server.runProvision(context.Background(), op))
		require.Equal(t, "cancelling", op.Phase)
		require.Equal(t, attempt, op.Attempt)
		require.NoError(t, server.runProvision(context.Background(), op))
	}
	require.Equal(t, "failed", op.Phase)
	require.Equal(t, 3, op.Attempt)
	for _, n := range op.Nodes {
		require.Equal(t, 3, cleared[n])
	}
}

func TestProvisionIncompleteCleanupDoesNotRetry(t *testing.T) {
	store := &memoryProvisions{}
	server := &SiteServer{provisions: store, provisionClient: fakeProvisionClient{call: func(ctx context.Context, n provisionNode) (provisionResult, error) {
		<-ctx.Done()
		return provisionResult{}, ctx.Err()
	}}}
	op := testProvision()
	op.Phase = "cancelling"
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	require.Error(t, server.runProvision(ctx, op))
	require.Equal(t, "cancelling", op.Phase)
	require.Equal(t, 1, op.Attempt)
}

func TestProvisionCallerCancellationDoesNotEraseWork(t *testing.T) {
	op := testProvision()
	store := &memoryProvisions{op: *op}
	store.begin = func(*db.Site, []string) (*db.SiteProvision, error) { return op, nil }
	server := &SiteServer{provisions: store}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := server.addProvisionedSite(ctx, &op.Site, op.Nodes)
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, "configuring", store.op.Phase)
}

func TestProvisionSiteSaveFailureRollsBack(t *testing.T) {
	op := testProvision()
	op.Phase = "creating"
	store := &memoryProvisions{op: *op}
	sites := &mocks.SiteRepo{}
	sites.On("Add", mock.Anything, mock.Anything).Return(errors.New("site insert failed")).Once()
	server := &SiteServer{provisions: store, siteRepo: sites, provisionClient: fakeProvisionClient{call: func(_ context.Context, node provisionNode) (provisionResult, error) {
		return provisionResult{Cancelled: true, Cleared: true}, nil
	}}}
	require.NoError(t, server.runProvision(context.Background(), op))
	require.Equal(t, "cancelling", op.Phase)
	require.True(t, op.Stop)
	require.NoError(t, server.runProvision(context.Background(), op))
	require.Equal(t, "failed", op.Phase)
	require.Equal(t, 1, op.Attempt, "site persistence failure does not reconfigure nodes")
	sites.AssertExpectations(t)
}
