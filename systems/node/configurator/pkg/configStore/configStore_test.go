/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package configstore_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/node/configurator/mocks"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"

	configstore "github.com/ukama/ukama/systems/node/configurator/pkg/configStore"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const OrgName = "testorg"

var (
	testNode1 = "uk-000000-hnode-0000"
	testNode2 = "uk-000000-hnode-0001"
	Service   = "node/configurator"
	TestData  = "node/configurator/test/integration/data"
)

func TestConfigStore_PlanConfigStoreEvent(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	commitRepo := &mocks.CommitRepo{}
	configRepo := &mocks.ConfigRepo{}
	store := &mocks.StoreProvider{}

	cS := configstore.NewConfigStore(msgbusClient, nil, nil, nil, configRepo, commitRepo, OrgName, store, (10 * time.Second))
	t.Run("SameVersion", func(t *testing.T) {
		store.On("GetLatestRemoteConfigs", mock.Anything).Return("000", nil).Once()
		commitRepo.On("GetLatest").Return(&db.Commit{Hash: "000"}, nil).Once()
		store.On("GetRemoteConfigVersion", mock.Anything, mock.Anything).Return(nil).Once()
		plan, err := cS.PlanConfigStoreEvent(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, plan.Nodes())
		commitRepo.AssertExpectations(t)
	})

	t.Run("DifferentVersionButNoChanges", func(t *testing.T) {
		store.On("GetLatestRemoteConfigs", mock.Anything).Return("000", nil).Once()
		commitRepo.On("GetLatest").Return(&db.Commit{Hash: "001"}, nil).Once()
		store.On("GetRemoteConfigVersion", mock.Anything, mock.Anything).Return(nil).Once()
		store.On("GetDiff", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		plan, err := cS.PlanConfigStoreEvent(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, plan.Nodes())
		commitRepo.AssertExpectations(t)
	})

}

func TestConfigStore_PlanAndApply(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	commitRepo := &mocks.CommitRepo{}
	configRepo := &mocks.ConfigRepo{}
	store := &mocks.StoreProvider{}
	cVer := "0.0.0."
	rVer := "0.0.1"
	path, err := os.Getwd()
	assert.NoError(t, err)
	p := strings.Split(path, Service)
	dir := p[0] + TestData
	cS := configstore.NewConfigStore(msgbusClient, nil, nil, nil, configRepo, commitRepo, OrgName, store, (10 * time.Second))

	t.Run("DifferentVersionWithChanges", func(t *testing.T) {
		var node string
		store.On("GetDiff", mock.Anything, mock.Anything, mock.Anything).Return([]string{"networkABC/siteXYZ/uk-000000-hnode-0000/epc/epc.json", "networkABC/siteXYZ/uk-000000-hnode-0000/deviced/deviced.json", "networkABC/siteXYZ/uk-000000-hnode-0001/epc/epc.json"}, nil).Once()
		configRepo.On("Get", mock.MatchedBy(func(n string) bool {
			switch n {
			case testNode1:
				node = testNode1
				return true
			case testNode2:
				node = testNode2
				return true
			default:
				return false
			}
		})).Return(&db.Configuration{NodeId: node}, nil)

		msgbusClient.On("PublishRequest", mock.AnythingOfType("string"), mock.Anything).Return(nil)
		commitRepo.On("Get", mock.AnythingOfType("string"), mock.Anything).Return(nil, sql.ErrNoRows)
		configRepo.On("UpdateLastCommit", mock.Anything, mock.MatchedBy(func(a *db.CommitState) bool { return a != nil && *a == db.Published })).Return(nil)
		files, ldir, err := cS.LookingForChanges(dir, cVer, rVer)
		assert.NoError(t, err)
		plan := cS.PlanCommit(files, cVer, ldir)
		assert.Equal(t, []string{testNode1, testNode2}, plan.Nodes())
		results := cS.Apply(plan)
		assert.Equal(t, map[string]error{testNode1: nil, testNode2: nil}, results)
		configRepo.AssertExpectations(t)
	})

}

func deletedFiles() []configstore.FilesToUpdate {
	return []configstore.FilesToUpdate{
		{Name: "networkABC/siteXYZ/" + testNode1 + "/epc/epc.json",
			Reason: configstore.REASON_DELETED},
		{Name: "networkABC/siteXYZ/" + testNode2 + "/epc/epc.json",
			Reason: configstore.REASON_DELETED},
	}
}

func TestConfigStore_PlanCommit(t *testing.T) {
	t.Run("DropsNodesWithoutConfiguration", func(t *testing.T) {
		configRepo := &mocks.ConfigRepo{}
		cS := configstore.NewConfigStore(&mbmocks.MsgBusServiceClient{}, nil, nil, nil,
			configRepo, &mocks.CommitRepo{}, OrgName, &mocks.StoreProvider{},
			10*time.Second)
		configRepo.On("Get", testNode1).Return(&db.Configuration{NodeId: testNode1}, nil)
		configRepo.On("Get", testNode2).Return(nil, sql.ErrNoRows)

		plan := cS.PlanCommit(deletedFiles(), "0.0.1", "/unused/")

		assert.Equal(t, []string{testNode1}, plan.Nodes())
		assert.Equal(t, "0.0.1", plan.Version)
	})

	t.Run("NoFilesGivesNoNodes", func(t *testing.T) {
		cS := configstore.NewConfigStore(&mbmocks.MsgBusServiceClient{}, nil, nil, nil,
			&mocks.ConfigRepo{}, &mocks.CommitRepo{}, OrgName, &mocks.StoreProvider{},
			10*time.Second)

		plan := cS.PlanCommit(nil, "0.0.1", "/unused/")

		assert.Empty(t, plan.Nodes())
	})
}

func TestConfigStore_Apply(t *testing.T) {
	t.Run("PublishFailureFailsOnlyThatNode", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		configRepo := &mocks.ConfigRepo{}
		commitRepo := &mocks.CommitRepo{}
		cS := configstore.NewConfigStore(msgbusClient, nil, nil, nil, configRepo,
			commitRepo, OrgName, &mocks.StoreProvider{}, 10*time.Second)
		configRepo.On("Get", testNode1).Return(&db.Configuration{NodeId: testNode1}, nil)
		configRepo.On("Get", testNode2).Return(&db.Configuration{NodeId: testNode2}, nil)
		commitRepo.On("Get", "0.0.1").Return(nil, sql.ErrNoRows)
		msgbusClient.On("PublishRequest", mock.Anything,
			mock.MatchedBy(func(m *epb.NodeFeederMessage) bool {
				return strings.HasSuffix(m.Target, testNode1)
			})).Return(nil)
		msgbusClient.On("PublishRequest", mock.Anything,
			mock.MatchedBy(func(m *epb.NodeFeederMessage) bool {
				return strings.HasSuffix(m.Target, testNode2)
			})).Return(errors.New("bus down"))
		configRepo.On("UpdateLastCommit",
			mock.MatchedBy(func(c db.Configuration) bool {
				return c.NodeId == testNode1
			}),
			mock.MatchedBy(func(s *db.CommitState) bool { return *s == db.Published })).
			Return(nil).Once()
		configRepo.On("UpdateLastCommit",
			mock.MatchedBy(func(c db.Configuration) bool {
				return c.NodeId == testNode2
			}),
			mock.MatchedBy(func(s *db.CommitState) bool { return *s == db.Failed })).
			Return(nil).Once()

		plan := cS.PlanCommit(deletedFiles(), "0.0.1", "/unused/")
		results := cS.Apply(plan)

		assert.NoError(t, results[testNode1])
		assert.ErrorContains(t, results[testNode2], "bus down")
		configRepo.AssertExpectations(t)
	})

	t.Run("RecordFailureIsReported", func(t *testing.T) {
		msgbusClient := &mbmocks.MsgBusServiceClient{}
		configRepo := &mocks.ConfigRepo{}
		commitRepo := &mocks.CommitRepo{}
		cS := configstore.NewConfigStore(msgbusClient, nil, nil, nil, configRepo,
			commitRepo, OrgName, &mocks.StoreProvider{}, 10*time.Second)
		configRepo.On("Get", testNode1).Return(&db.Configuration{NodeId: testNode1}, nil)
		configRepo.On("Get", testNode2).Return(nil, sql.ErrNoRows).Once()
		commitRepo.On("Get", "0.0.1").Return(nil, sql.ErrNoRows)
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)
		configRepo.On("UpdateLastCommit", mock.Anything, mock.Anything).
			Return(errors.New("db down")).Once()

		plan := cS.PlanCommit(deletedFiles(), "0.0.1", "/unused/")
		results := cS.Apply(plan)

		assert.Len(t, results, 1)
		assert.ErrorContains(t, results[testNode1], "db down")
	})
}
