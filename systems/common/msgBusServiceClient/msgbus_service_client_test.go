/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package msgBusServiceClient

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	pb "github.com/ukama/ukama/systems/common/pb/gen/msgclient"
	"google.golang.org/grpc"
)

type fakeMsgClient struct {
	pb.MsgClientServiceClient
	failures int32
	interval uint32
	calls    atomic.Int32
}

func (f *fakeMsgClient) RegisterService(ctx context.Context, in *pb.RegisterServiceReq, opts ...grpc.CallOption) (*pb.RegisterServiceResp, error) {
	if f.calls.Add(1) <= f.failures {
		return nil, errors.New("unavailable")
	}
	return &pb.RegisterServiceResp{State: pb.REGISTRAION_STATUS_REGISTERED, ServiceUuid: "uuid", RefreshInterval: f.interval}, nil
}

func newTestClient(f *fakeMsgClient) *msgBusServiceClient {
	return &msgBusServiceClient{client: f, timeout: time.Second, stop: make(chan struct{})}
}

func TestRegisterRetriesUntilRegistered(t *testing.T) {
	registerRetryInterval = time.Millisecond
	f := &fakeMsgClient{failures: 2, interval: 5}
	m := newTestClient(f)

	if err := m.Register(); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if got := f.calls.Load(); got != 3 {
		t.Fatalf("expected 3 register calls, got %d", got)
	}
	if m.uuid != "uuid" || m.refreshInterval != 5*time.Second {
		t.Fatalf("unexpected uuid %q or interval %s", m.uuid, m.refreshInterval)
	}
}

func TestKeepRegisteredRefreshesUntilStopped(t *testing.T) {
	fallbackRefreshInterval = 5 * time.Millisecond
	f := &fakeMsgClient{failures: 1}
	m := newTestClient(f)
	m.refreshInterval = 5 * time.Millisecond

	done := make(chan struct{})
	go func() {
		m.keepRegistered()
		close(done)
	}()

	deadline := time.Now().Add(time.Second)
	for f.calls.Load() < 3 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if f.calls.Load() < 3 {
		t.Fatalf("expected refresh to continue after a failure, got %d calls", f.calls.Load())
	}

	m.stopOnce.Do(func() { close(m.stop) })
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("refresh loop did not stop")
	}
}

func TestRefreshIntervalFallsBackForOldMsgClient(t *testing.T) {
	fallbackRefreshInterval = 30 * time.Second
	if got := refreshInterval(&pb.RegisterServiceResp{RefreshInterval: 2}); got != 2*time.Second {
		t.Fatalf("expected 2s, got %s", got)
	}
	if got := refreshInterval(&pb.RegisterServiceResp{}); got != 30*time.Second {
		t.Fatalf("expected fallback 30s, got %s", got)
	}
}
