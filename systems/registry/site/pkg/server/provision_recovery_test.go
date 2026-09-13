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
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/registry/site/mocks"
	"github.com/ukama/ukama/systems/registry/site/pkg/db"
)

type recoveryStore struct {
	memoryProvisions
	getErr  error
	saveErr error
	pending func(context.Context) ([]db.SiteProvision, error)
	run     func(context.Context, string, func(context.Context, *db.SiteProvision) error) error
}

func (s *recoveryStore) Get(ctx context.Context, id string) (*db.SiteProvision, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.memoryProvisions.Get(ctx, id)
}
func (s *recoveryStore) Save(ctx context.Context, op *db.SiteProvision) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	return s.memoryProvisions.Save(ctx, op)
}
func (s *recoveryStore) Pending(ctx context.Context) ([]db.SiteProvision, error) {
	return s.pending(ctx)
}
func (s *recoveryStore) Run(ctx context.Context, id string, run func(context.Context, *db.SiteProvision) error) error {
	return s.run(ctx, id, run)
}

func TestProvisionSuccessfulGroupAndCancellation(t *testing.T) {
	store := &memoryProvisions{}
	server := &SiteServer{provisions: store, provisionClient: fakeProvisionClient{call: func(context.Context, provisionNode) (provisionResult, error) {
		return provisionResult{Completed: true}, nil
	}}}
	op := testProvision()
	require.NoError(t, server.runProvision(context.Background(), op))
	require.Equal(t, "creating", store.op.Phase)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	op = testProvision()
	require.ErrorIs(t, server.runProvision(ctx, op), context.Canceled)
	require.Equal(t, "configuring", op.Phase)
	server.provisionClient = fakeProvisionClient{call: func(context.Context, provisionNode) (provisionResult, error) {
		return provisionResult{Cancelled: true}, nil
	}}
	require.ErrorContains(t, server.waitNode(context.Background(), attemptNode(op, 0, "configure")), "cancelled")
}

func TestProvisionWaitPollsAfterDispatch(t *testing.T) {
	calls := 0
	server := &SiteServer{provisionClient: fakeProvisionClient{call: func(_ context.Context, node provisionNode) (provisionResult, error) {
		calls++
		if calls == 1 {
			require.Equal(t, "configure", node.Action)
			return provisionResult{}, nil
		}
		require.Equal(t, "status", node.Action)
		return provisionResult{Completed: true}, nil
	}}}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	require.NoError(t, server.waitNode(ctx, attemptNode(testProvision(), 0, "configure")))
	require.Equal(t, 2, calls)
}

func TestProvisionCreationRecovery(t *testing.T) {
	for _, mode := range []string{"saved", "read-error", "already-publishing"} {
		t.Run(mode, func(t *testing.T) {
			op := testProvision()
			op.Phase = "creating"
			store := &recoveryStore{memoryProvisions: memoryProvisions{op: *op}}
			sites := &mocks.SiteRepo{}
			if mode == "saved" {
				sites.On("Add", mock.Anything, mock.Anything).Return(nil).Once()
			} else {
				sites.On("Add", mock.Anything, mock.Anything).Return(errors.New("commit response lost")).Once()
				if mode == "read-error" {
					store.getErr = errors.New("read unavailable")
				} else {
					store.op.Phase = "publishing"
				}
			}
			server := &SiteServer{provisions: store, siteRepo: sites}
			err := server.runProvision(context.Background(), op)
			if mode == "read-error" {
				require.ErrorIs(t, err, store.getErr)
			} else {
				require.NoError(t, err)
			}
			require.False(t, op.Stop, "a committed site must not be rolled back due to a lost acknowledgement")
			sites.AssertExpectations(t)
		})
	}
}

func TestProvisionPublicationRecovery(t *testing.T) {
	for _, mode := range []string{"success", "lookup-error", "publish-error", "checkpoint-error"} {
		t.Run(mode, func(t *testing.T) {
			op := testProvision()
			op.Phase = "publishing"
			store := &recoveryStore{memoryProvisions: memoryProvisions{op: *op}}
			sites := &mocks.SiteRepo{}
			bus := &mbmocks.MsgBusServiceClient{}
			var lookupErr, publishErr error
			if mode == "lookup-error" {
				lookupErr = errors.New("site unavailable")
			}
			if mode == "publish-error" {
				publishErr = errors.New("broker unavailable")
			}
			if mode == "checkpoint-error" {
				store.saveErr = errors.New("checkpoint unavailable")
			}
			sites.On("Get", op.Site.Id).Return(&op.Site, lookupErr).Once()
			if lookupErr == nil {
				bus.On("PublishRequest", mock.Anything, mock.Anything).Return(publishErr).Once()
			}
			metrics := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
			defer metrics.Close()
			if mode == "success" {
				sites.On("GetSiteCount", op.Site.NetworkId).Return(int64(1), nil).Once()
			}
			server := NewSiteServer("org", sites, bus, nil, metrics.URL, nil, nil, nil)
			server.provisions = store
			err := server.runProvision(context.Background(), op)
			if mode == "success" {
				require.NoError(t, err)
				require.Equal(t, "active", store.op.Phase)
			} else {
				require.Error(t, err)
				require.Equal(t, "publishing", store.op.Phase, "failed publication/checkpoint must remain retryable")
			}
			sites.AssertExpectations(t)
			bus.AssertExpectations(t)
		})
	}
}

func TestProvisionAddOutcomes(t *testing.T) {
	for _, mode := range []string{"uninitialized", "create-error", "active", "failed", "poll-active", "poll-error"} {
		t.Run(mode, func(t *testing.T) {
			op := testProvision()
			store := &recoveryStore{memoryProvisions: memoryProvisions{op: *op}}
			store.begin = func(*db.Site, []string) (*db.SiteProvision, error) {
				if mode == "create-error" {
					return nil, errors.New("reserved")
				}
				if mode == "active" || mode == "failed" {
					op.Phase = mode
				}
				return op, nil
			}
			if mode == "poll-active" {
				store.op.Phase = "active"
			}
			if mode == "poll-error" {
				store.getErr = errors.New("database offline")
			}
			server := &SiteServer{provisions: store}
			if mode == "uninitialized" {
				server.provisions = nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			site, err := server.addProvisionedSite(ctx, &op.Site, op.Nodes)
			if mode == "active" || mode == "poll-active" {
				require.NoError(t, err)
				require.Equal(t, op.Site.NetworkId, site.NetworkId)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestProvisionWorkerProcessesRecoveredOperation(t *testing.T) {
	op := testProvision()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	completed := make(chan struct{})
	var once sync.Once
	store := &recoveryStore{memoryProvisions: memoryProvisions{op: *op}}
	store.pending = func(context.Context) ([]db.SiteProvision, error) { return []db.SiteProvision{{ID: op.ID}}, nil }
	store.run = func(ctx context.Context, id string, run func(context.Context, *db.SiteProvision) error) error {
		err := run(ctx, op)
		once.Do(func() { close(completed) })
		return err
	}
	server := &SiteServer{provisions: store, provisionClient: fakeProvisionClient{call: func(context.Context, provisionNode) (provisionResult, error) {
		return provisionResult{Completed: true}, nil
	}}}
	done := make(chan struct{})
	go func() { server.provisionWorker(ctx); close(done) }()
	select {
	case <-completed:
	case <-time.After(time.Second):
		t.Fatal("recovered operation was not processed")
	}
	cancel()
	<-done
	stored, err := store.Get(context.Background(), op.ID)
	require.NoError(t, err)
	require.Equal(t, "creating", stored.Phase)
}
