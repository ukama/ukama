/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/node/configurator/mocks"
	"github.com/ukama/ukama/systems/node/configurator/pkg"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	pb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
	configstore "github.com/ukama/ukama/systems/node/configurator/pkg/configStore"
)

const (
	testOrgName = "testOrg"
	testHash    = "4f6e609"
	testNodeA   = "uk-sa0000-tnode-a1-0001"
	testNodeB   = "uk-sa0000-tnode-a1-0002"
)

type serverDeps struct {
	configStore *mocks.ConfigStoreProvider
	opMgr       *cmocks.ManagerClient
}

func newTestServer() (*ConfiguratorServer, serverDeps) {
	d := serverDeps{
		configStore: &mocks.ConfigStoreProvider{},
		opMgr:       &cmocks.ManagerClient{},
	}
	s := NewConfiguratorServer(&cmocks.MsgBusServiceClient{}, &mocks.ConfigRepo{},
		&mocks.CommitRepo{}, d.configStore, testOrgName, pkg.IsDebugMode, d.opMgr,
		1800)
	return s, d
}

func twoNodePlan() *configstore.CommitPlan {
	return &configstore.CommitPlan{
		Version: testHash,
		NodeFiles: map[string][]string{
			testNodeB: {"net/site/" + testNodeB + "/epc/epc.json"},
			testNodeA: {"net/site/" + testNodeA + "/epc/epc.json"},
		},
	}
}

func twoNodeOps() []*copr.OperationInfo {
	return []*copr.OperationInfo{
		{Id: "op-a", FencingToken: 1, ResourceKey: "node:" + testNodeA},
		{Id: "op-b", FencingToken: 2, ResourceKey: "node:" + testNodeB},
	}
}

func expectStart(d serverDeps, ops []*copr.OperationInfo, err error) {
	d.opMgr.On("Start", mock.MatchedBy(func(r copr.StartRequest) bool {
		return r.Type == applyConfigType && r.System == pkg.SystemName &&
			r.RequestedBy == pkg.ServiceName && r.LeaseSeconds == 1800 &&
			len(r.MustBeFree) == 0 && len(r.ResourceKeys) == 2 &&
			r.ResourceKeys[0] == "node:"+testNodeA &&
			r.ResourceKeys[1] == "node:"+testNodeB
	})).Return(&copr.StartResponse{Operations: ops}, err).Once()
}

func expectMarkRunning(d serverDeps, ops []*copr.OperationInfo) {
	for _, op := range ops {
		d.opMgr.On("MarkRunning", op.Id, op.FencingToken, pkg.ServiceName).
			Return(&copr.OperationInfo{}, nil).Once()
	}
}

func TestConfiguratorServer_ApplyConfig(t *testing.T) {
	req := &pb.ApplyConfigRequest{Hash: testHash}

	t.Run("NoChangesTakesNoLock", func(t *testing.T) {
		s, d := newTestServer()
		d.configStore.On("PlanConfigCommit", mock.Anything, testHash).
			Return(&configstore.CommitPlan{Version: testHash}, nil).Once()

		_, err := s.ApplyConfig(context.Background(), req)

		assert.NoError(t, err)
		d.configStore.AssertExpectations(t)
		d.opMgr.AssertNotCalled(t, "Start", mock.Anything)
		d.configStore.AssertNotCalled(t, "Apply", mock.Anything)
	})

	t.Run("PlanErrorIsInternal", func(t *testing.T) {
		s, d := newTestServer()
		d.configStore.On("PlanConfigCommit", mock.Anything, testHash).
			Return(nil, errors.New("store down")).Once()

		_, err := s.ApplyConfig(context.Background(), req)

		assert.Equal(t, codes.Internal, status.Code(err))
		d.opMgr.AssertNotCalled(t, "Start", mock.Anything)
	})

	t.Run("LocksEachNodeAndCompletesAfterDispatch", func(t *testing.T) {
		s, d := newTestServer()
		plan := twoNodePlan()
		ops := twoNodeOps()
		d.configStore.On("PlanConfigCommit", mock.Anything, testHash).
			Return(plan, nil).Once()
		expectStart(d, ops, nil)
		expectMarkRunning(d, ops)
		d.configStore.On("Apply", plan).
			Return(map[string]error{testNodeA: nil, testNodeB: nil}).Once()
		for _, op := range ops {
			d.opMgr.On("Complete", op.Id, op.FencingToken, pkg.ServiceName,
				"dispatched").Return(&copr.OperationInfo{}, nil).Once()
		}

		_, err := s.ApplyConfig(context.Background(), req)

		assert.NoError(t, err)
		d.configStore.AssertExpectations(t)
		d.opMgr.AssertExpectations(t)
	})

	t.Run("BusyNodeIsAlreadyExistsAndNothingSent", func(t *testing.T) {
		s, d := newTestServer()
		d.configStore.On("PlanConfigCommit", mock.Anything, testHash).
			Return(twoNodePlan(), nil).Once()
		expectStart(d, nil, &copr.BusyError{Reason: "node busy", Err: copr.ErrBusy})

		_, err := s.ApplyConfig(context.Background(), req)

		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		d.configStore.AssertNotCalled(t, "Apply", mock.Anything)
		d.opMgr.AssertNotCalled(t, "MarkRunning", mock.Anything, mock.Anything,
			mock.Anything)
	})

	t.Run("LockFailureIsUnavailable", func(t *testing.T) {
		s, d := newTestServer()
		d.configStore.On("PlanConfigCommit", mock.Anything, testHash).
			Return(twoNodePlan(), nil).Once()
		expectStart(d, nil, copr.ErrUnavailable)

		_, err := s.ApplyConfig(context.Background(), req)

		assert.Equal(t, codes.Unavailable, status.Code(err))
		d.configStore.AssertNotCalled(t, "Apply", mock.Anything)
	})

	t.Run("MarkRunningRejectedFailsAllAndNothingSent", func(t *testing.T) {
		s, d := newTestServer()
		ops := twoNodeOps()
		d.configStore.On("PlanConfigCommit", mock.Anything, testHash).
			Return(twoNodePlan(), nil).Once()
		expectStart(d, ops, nil)
		d.opMgr.On("MarkRunning", "op-a", uint64(1), pkg.ServiceName).
			Return(nil, copr.ErrNotRunnable).Once()
		for _, op := range ops {
			d.opMgr.On("Fail", op.Id, op.FencingToken, pkg.ServiceName,
				mock.Anything).Return(&copr.OperationInfo{}, nil).Once()
		}

		_, err := s.ApplyConfig(context.Background(), req)

		assert.Equal(t, codes.Aborted, status.Code(err))
		d.opMgr.AssertExpectations(t)
		d.configStore.AssertNotCalled(t, "Apply", mock.Anything)
	})

	t.Run("NodeSendFailureFailsOnlyThatNode", func(t *testing.T) {
		s, d := newTestServer()
		plan := twoNodePlan()
		ops := twoNodeOps()
		d.configStore.On("PlanConfigCommit", mock.Anything, testHash).
			Return(plan, nil).Once()
		expectStart(d, ops, nil)
		expectMarkRunning(d, ops)
		d.configStore.On("Apply", plan).Return(map[string]error{
			testNodeA: nil,
			testNodeB: errors.New("publish epc.json: bus down"),
		}).Once()
		d.opMgr.On("Complete", "op-a", uint64(1), pkg.ServiceName, "dispatched").
			Return(&copr.OperationInfo{}, nil).Once()
		d.opMgr.On("Fail", "op-b", uint64(2), pkg.ServiceName,
			"publish epc.json: bus down").Return(&copr.OperationInfo{}, nil).Once()

		_, err := s.ApplyConfig(context.Background(), req)

		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Contains(t, err.Error(), testNodeB)
		assert.NotContains(t, err.Error(), testNodeA)
		d.opMgr.AssertExpectations(t)
	})
}

func TestConfiguratorServer_ConfigEvent(t *testing.T) {
	t.Run("TakesTheSameNodeLocks", func(t *testing.T) {
		s, d := newTestServer()
		plan := twoNodePlan()
		ops := twoNodeOps()
		d.configStore.On("PlanConfigStoreEvent", mock.Anything).Return(plan, nil).Once()
		expectStart(d, ops, nil)
		expectMarkRunning(d, ops)
		d.configStore.On("Apply", plan).
			Return(map[string]error{testNodeA: nil, testNodeB: nil}).Once()
		for _, op := range ops {
			d.opMgr.On("Complete", op.Id, op.FencingToken, pkg.ServiceName,
				"dispatched").Return(&copr.OperationInfo{}, nil).Once()
		}

		_, err := s.ConfigEvent(context.Background(), &pb.ConfigStoreEvent{})

		assert.NoError(t, err)
		d.configStore.AssertExpectations(t)
		d.opMgr.AssertExpectations(t)
	})

	t.Run("BusyNodeIsAlreadyExistsAndNothingSent", func(t *testing.T) {
		s, d := newTestServer()
		d.configStore.On("PlanConfigStoreEvent", mock.Anything).
			Return(twoNodePlan(), nil).Once()
		expectStart(d, nil, &copr.BusyError{Reason: "node busy", Err: copr.ErrBusy})

		_, err := s.ConfigEvent(context.Background(), &pb.ConfigStoreEvent{})

		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		d.configStore.AssertNotCalled(t, "Apply", mock.Anything)
	})

	t.Run("NoChangesTakesNoLock", func(t *testing.T) {
		s, d := newTestServer()
		d.configStore.On("PlanConfigStoreEvent", mock.Anything).Return(nil, nil).Once()

		_, err := s.ConfigEvent(context.Background(), &pb.ConfigStoreEvent{})

		assert.NoError(t, err)
		d.opMgr.AssertNotCalled(t, "Start", mock.Anything)
	})
}

func TestConfiguratorServer_GetConfigVersion(t *testing.T) {
	s, _ := newTestServer()
	configRepo := &mocks.ConfigRepo{}
	s.configRepo = configRepo
	configRepo.On("Get", "node-id").Return(&db.Configuration{
		NodeId:          "node-id",
		State:           db.Published,
		Commit:          db.Commit{Hash: "commit"},
		LastCommitState: db.Published,
		LastCommit:      db.Commit{Hash: "lastcommit"},
	}, nil).Once()

	resp, err := s.GetConfigVersion(context.Background(),
		&pb.ConfigVersionRequest{NodeId: "node-id"})

	assert.NoError(t, err)
	assert.Equal(t, "published", resp.Status)
	assert.Equal(t, "lastcommit", resp.LastCommit)
	configRepo.AssertExpectations(t)
}

func TestMarkRunningCode(t *testing.T) {
	assert.Equal(t, codes.Aborted, markRunningCode(copr.ErrNotRunnable))
	assert.Equal(t, codes.Unavailable, markRunningCode(copr.ErrUnavailable))
	assert.Equal(t, codes.Internal, markRunningCode(errors.New("boom")))
}

func TestConfiguratorServer_TerminateErrorsDoNotFailApply(t *testing.T) {
	s, d := newTestServer()
	plan := twoNodePlan()
	ops := twoNodeOps()
	d.configStore.On("PlanConfigCommit", mock.Anything, testHash).Return(plan, nil).Once()
	expectStart(d, ops, nil)
	expectMarkRunning(d, ops)
	d.configStore.On("Apply", plan).Return(map[string]error{
		testNodeA: nil,
		testNodeB: errors.New("bus down"),
	}).Once()
	d.opMgr.On("Complete", "op-a", uint64(1), pkg.ServiceName, "dispatched").
		Return(nil, copr.ErrUnavailable).Once()
	d.opMgr.On("Fail", "op-b", uint64(2), pkg.ServiceName, "bus down").
		Return(nil, copr.ErrUnavailable).Once()

	_, err := s.ApplyConfig(context.Background(), &pb.ApplyConfigRequest{Hash: testHash})

	assert.Equal(t, codes.Internal, status.Code(err))
	d.opMgr.AssertExpectations(t)
}

func TestConfiguratorServer_ConfigEventPlanError(t *testing.T) {
	s, d := newTestServer()
	d.configStore.On("PlanConfigStoreEvent", mock.Anything).
		Return(nil, errors.New("store down")).Once()

	_, err := s.ConfigEvent(context.Background(), &pb.ConfigStoreEvent{})

	assert.Error(t, err)
	d.opMgr.AssertNotCalled(t, "Start", mock.Anything)
}
