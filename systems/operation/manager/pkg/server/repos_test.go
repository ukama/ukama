/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/operation/manager/mocks"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

const orgName = "ukama"

type testRepos struct {
	ops      *mocks.OperationRepo
	locks    *mocks.LockRepo
	audits   *mocks.AuditRepo
	saved    []db.OperationAudit
	released []uuid.UUID
}

func newTestRepos() *testRepos {
	return &testRepos{
		ops: &mocks.OperationRepo{}, locks: &mocks.LockRepo{}, audits: &mocks.AuditRepo{},
	}
}

func (r *testRepos) repos() Repos {
	return Repos{Operations: r.ops, Locks: r.locks, Audits: r.audits}
}

func (r *testRepos) server() *OperationServer {
	return NewOperationServer(orgName, "", r.repos(), nil)
}

func (r *testRepos) assertExpectations(t *testing.T) {
	t.Helper()
	r.ops.AssertExpectations(t)
	r.locks.AssertExpectations(t)
	r.audits.AssertExpectations(t)
}

func (r *testRepos) assertUntouched(t *testing.T) {
	t.Helper()
	assert.Empty(t, r.ops.Calls)
	assert.Empty(t, r.locks.Calls)
	assert.Empty(t, r.audits.Calls)
}

func (r *testRepos) recordAudits() *testRepos {
	record := func(args mock.Arguments) {
		r.saved = append(r.saved, *args.Get(1).(*db.OperationAudit))
	}
	r.audits.On("AddInTx", mock.Anything, mock.Anything).Run(record).Return(nil).Maybe()
	r.audits.On("Add", mock.Anything, mock.Anything).Run(record).Return(nil).Maybe()
	return r
}

func (r *testRepos) events() []string {
	out := make([]string, 0, len(r.saved))
	for _, a := range r.saved {
		out = append(out, a.Event)
	}
	return out
}

func (r *testRepos) calls(method string) []mock.Call {
	var out []mock.Call
	for _, c := range r.ops.Calls {
		if c.Method == method {
			out = append(out, c)
		}
	}
	return out
}

func (r *testRepos) adds(mustBeFree, lease interface{}) {
	r.ops.On("Add", mock.Anything, mock.Anything, mock.Anything).
		Return(func(_ context.Context, ops []*db.Operation,
			nested func(*gorm.DB) error) error {
			for i, op := range ops {
				op.FencingToken = uint64(i + 1)
			}
			return nested(nil)
		}).Once()
	r.locks.On("Acquire", mock.Anything, mock.Anything, mustBeFree, lease).
		Return(nil).Once()
	r.recordAudits()
}

func (r *testRepos) addRunsNested() {
	r.ops.On("Add", mock.Anything, mock.Anything, mock.Anything).
		Return(func(_ context.Context, _ []*db.Operation,
			nested func(*gorm.DB) error) error {
			return nested(nil)
		}).Once()
}

func (r *testRepos) addFails(err error) {
	r.ops.On("Add", mock.Anything, mock.Anything, mock.Anything).Return(err).Once()
}

func leaderFilter(key string) interface{} {
	return mock.MatchedBy(func(f db.OperationFilter) bool {
		return f.IdempotencyKey == key && f.GroupId == nil
	})
}

func membersFilter(leader uuid.UUID) interface{} {
	return mock.MatchedBy(func(f db.OperationFilter) bool {
		return f.GroupId != nil && *f.GroupId == leader &&
			f.OrderBy == db.OrderByResourceKey && f.IdempotencyKey == ""
	})
}

func (r *testRepos) findsGroup(group []*db.Operation) {
	if len(group) == 0 {
		r.ops.On("List", mock.Anything, leaderFilter("key")).Return(nil, nil).Once()
		return
	}
	r.ops.On("List", mock.Anything, leaderFilter("key")).Return(group[:1], nil).Once()
	r.ops.On("List", mock.Anything, membersFilter(group[0].Id)).
		Return(group[1:], nil).Once()
}

func (r *testRepos) groupLookupFails(err error) {
	r.ops.On("List", mock.Anything, leaderFilter("key")).Return(nil, err).Once()
}

func (r *testRepos) updates(current db.Operation) {
	r.ops.On("Update", mock.Anything, current.Id, mock.Anything).
		Return(func(_ context.Context, _ uuid.UUID,
			nested func(*db.Operation, *gorm.DB) error) (*db.Operation, error) {
			op := current
			if err := nested(&op, nil); err != nil {
				return nil, err
			}
			return &op, nil
		}).Once()
	r.locks.On("Release", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r.released = append(r.released, args.Get(1).(*db.Operation).Id)
		}).Return(nil).Maybe()
	r.recordAudits()
}

func (r *testRepos) updateFails(id uuid.UUID, err error) {
	r.ops.On("Update", mock.Anything, id, mock.Anything).Return(nil, err).Once()
}

func (r *testRepos) deletes(current db.Operation) {
	r.ops.On("Delete", mock.Anything, current.Id, mock.Anything).
		Return(func(_ context.Context, _ uuid.UUID,
			nested func(*db.Operation, *gorm.DB) error) error {
			op := current
			return nested(&op, nil)
		}).Once()
	r.recordAudits()
}
