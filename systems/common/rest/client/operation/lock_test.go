/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package operation_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/operation"
)

type failCall struct {
	id     string
	token  uint64
	actor  string
	reason string
}

type fakeManager struct {
	operation.ManagerClient
	held      map[string]*operation.OperationInfo
	startErr  map[string]error
	readErr   map[string]error
	failErr   error
	started   []string
	failed    []failCall
	nextToken uint64
}

func newFakeManager() *fakeManager {
	return &fakeManager{
		held:     map[string]*operation.OperationInfo{},
		startErr: map[string]error{},
		readErr:  map[string]error{},
	}
}

func (f *fakeManager) Start(req operation.StartRequest) (*operation.StartResponse,
	error) {
	if err := f.startErr[req.ResourceKey]; err != nil {
		return nil, err
	}
	f.nextToken++
	op := &operation.OperationInfo{
		Id:           fmt.Sprintf("op-%d", f.nextToken),
		Type:         req.Type,
		Status:       operation.StatusPending,
		FencingToken: f.nextToken,
		ResourceKey:  req.ResourceKey,
	}
	f.started = append(f.started, req.ResourceKey)
	return &operation.StartResponse{Operation: op}, nil
}

func (f *fakeManager) GetByResource(key string) (*operation.OperationInfo, error) {
	if err := f.readErr[key]; err != nil {
		return nil, err
	}
	return f.held[key], nil
}

func (f *fakeManager) Fail(id string, token uint64, actor,
	reason string) (*operation.OperationInfo, error) {
	f.failed = append(f.failed, failCall{id: id, token: token, actor: actor,
		reason: reason})
	return &operation.OperationInfo{Id: id}, f.failErr
}

func conflict() error {
	return fmt.Errorf("StartOperation failure: %w", &client.ErrorStatus{
		StatusCode: http.StatusConflict, RawError: errors.New("locked"),
	})
}

func lockRequest(targets, siblings []string) operation.LockRequest {
	return operation.LockRequest{
		Type: "RestartNode", System: "node", RequestedBy: "controller",
		LeaseSeconds: 180, Targets: targets, Siblings: siblings,
	}
}

func TestLockThenCheck(t *testing.T) {
	tnode, anode, cnode := "node:tnode", "node:anode", "node:cnode"
	site := []string{tnode, anode, cnode}

	t.Run("LocksTargetsWhenSiteIsIdle", func(t *testing.T) {
		m := newFakeManager()

		ops, err := operation.LockThenCheck(m,
			lockRequest([]string{tnode, anode}, site))

		assert.NoError(t, err)
		assert.Len(t, ops, 2)
		assert.Equal(t, []string{tnode, anode}, m.started)
		assert.Empty(t, m.failed)
	})

	t.Run("BusySiblingReleasesOwnLocks", func(t *testing.T) {
		m := newFakeManager()
		m.held[cnode] = &operation.OperationInfo{
			Id: "update", Type: "UpdateSoftware", Status: operation.StatusRunning,
		}

		ops, err := operation.LockThenCheck(m,
			lockRequest([]string{tnode, anode}, site))

		assert.Nil(t, ops)
		assert.True(t, errors.Is(err, operation.ErrBusy), err)
		var busy *operation.BusyError
		assert.True(t, errors.As(err, &busy))
		assert.Equal(t, cnode, busy.ResourceKey)
		assert.Equal(t, "update", busy.Holder.Id)
		assert.Len(t, m.failed, 2)
		for _, call := range m.failed {
			assert.Equal(t, "controller", call.actor)
			assert.Contains(t, call.reason, "not started: ")
			assert.Contains(t, call.reason, "UpdateSoftware")
		}
	})

	t.Run("FinishedSiblingIsNotBusy", func(t *testing.T) {
		m := newFakeManager()
		m.held[cnode] = &operation.OperationInfo{
			Id: "old", Status: operation.StatusSuccess,
		}

		ops, err := operation.LockThenCheck(m, lockRequest([]string{tnode}, site))

		assert.NoError(t, err)
		assert.Len(t, ops, 1)
	})

	t.Run("LockedTargetReleasesEarlierTargets", func(t *testing.T) {
		m := newFakeManager()
		m.startErr[anode] = conflict()

		ops, err := operation.LockThenCheck(m,
			lockRequest([]string{tnode, anode}, site))

		assert.Nil(t, ops)
		var busy *operation.BusyError
		assert.True(t, errors.As(err, &busy))
		assert.Equal(t, anode, busy.ResourceKey)
		assert.Nil(t, busy.Holder)
		assert.Len(t, m.failed, 1)
		assert.Equal(t, uint64(1), m.failed[0].token)
	})

	t.Run("OtherStartErrorIsNotBusy", func(t *testing.T) {
		m := newFakeManager()
		m.startErr[tnode] = fmt.Errorf("StartOperation failure: %w",
			&client.ErrorStatus{StatusCode: http.StatusServiceUnavailable,
				RawError: errors.New("down")})

		_, err := operation.LockThenCheck(m, lockRequest([]string{tnode}, site))

		assert.Error(t, err)
		assert.False(t, errors.Is(err, operation.ErrBusy), err)
		assert.Empty(t, m.failed)
	})

	t.Run("SiblingReadErrorReleasesOwnLocks", func(t *testing.T) {
		m := newFakeManager()
		m.readErr[cnode] = errors.New("timeout")

		_, err := operation.LockThenCheck(m, lockRequest([]string{tnode}, site))

		assert.Error(t, err)
		assert.False(t, errors.Is(err, operation.ErrBusy), err)
		assert.Len(t, m.failed, 1)
	})

	t.Run("ReleaseFailureIsReported", func(t *testing.T) {
		m := newFakeManager()
		m.held[cnode] = &operation.OperationInfo{
			Id: "update", Type: "UpdateSoftware", Status: operation.StatusPending,
		}
		m.failErr = errors.New("manager down")

		_, err := operation.LockThenCheck(m, lockRequest([]string{tnode}, site))

		assert.True(t, errors.Is(err, operation.ErrBusy), err)
		assert.Contains(t, err.Error(), "release node:tnode")
	})

	t.Run("TargetsAreNotCheckedAsSiblingsOrLockedTwice", func(t *testing.T) {
		m := newFakeManager()

		_, err := operation.LockThenCheck(m,
			lockRequest([]string{tnode, tnode}, site))

		assert.NoError(t, err)
		assert.Equal(t, []string{tnode}, m.started)
	})
}

func TestIsTerminal(t *testing.T) {
	for _, s := range []string{operation.StatusSuccess, operation.StatusFailed,
		operation.StatusTimeout, operation.StatusCancelled} {
		assert.True(t, operation.IsTerminal(s), s)
	}
	for _, s := range []string{operation.StatusPending, operation.StatusRunning, ""} {
		assert.False(t, operation.IsTerminal(s), s)
	}
}
