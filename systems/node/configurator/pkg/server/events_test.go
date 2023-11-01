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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	eCfgPb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/configurator/mocks"
	"github.com/ukama/ukama/systems/node/configurator/pkg"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
)

var testNode = ukama.NewVirtualNodeId("HomeNode")
var orgId = uuid.NewV4()

func TestConfiguratorServer_EventNotification(t *testing.T) {
	// Arrange
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	commitRepo := &mocks.CommitRepo{}
	configRepo := &mocks.ConfigRepo{}
	configStore := &mocks.ConfigStoreProvider{}
	registry := &mocks.RegistryProvider{}

	s := NewConfiguratorServer(msgbusClient, registry, configRepo, commitRepo, configStore, testOrgName, pkg.IsDebugMode)

	eventServer := NewConfiguratorEventServer(testOrgName, s)

	t.Run("AddNode", func(t *testing.T) {
		/* Node Cretaed event */
		evt := epb.NodeCreatedEvent{
			NodeId: testNode.String(),
			Name:   "testnode",
			Type:   "hnode",
			Org:    orgId.String(),
		}

		any, err := anypb.New(&evt)
		assert.NoError(t, err)

		configRepo.On("Add", testNode.String()).Return(nil).Once()

		_, err = eventServer.EventNotification(context.Background(), &epb.Event{
			RoutingKey: "event.cloud.local.testorg.registry.node.node.create",
			Msg:        any,
		})

		assert.NoError(t, err)
		// Assert
		configStore.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("UpdateNodeConfigSuccess", func(t *testing.T) {
		/* Node Cretaed event */
		evt := eCfgPb.NodeConfigUpdateEvent{
			NodeId: testNode.String(),
			Commit: "commithash",
			Status: eCfgPb.UpdateStatus_Success,
		}

		any, err := anypb.New(&evt)
		assert.NoError(t, err)
		cfg := db.Configuration{
			NodeId:          "node-id",
			State:           db.Published,
			Commit:          db.Commit{Hash: "commit"},
			LastCommitState: db.Published,
			LastCommit:      db.Commit{Hash: "lastcommit"},
		}

		cmt := db.Commit{Model: gorm.Model{ID: 2}, Hash: "commithash"}
		configRepo.On("Get", testNode.String()).Return(&cfg, nil).Once()

		commitRepo.On("Get", evt.Commit).Return(&cmt, nil).Once()
		cfg.Commit = cmt

		configRepo.On("UpdateCurrentCommit", mock.Anything, mock.MatchedBy(func(a *db.CommitState) bool {
			return *a == db.CommitState(evt.Status)
		})).Return(nil).Once()

		_, err = eventServer.EventNotification(context.Background(), &epb.Event{
			RoutingKey: "event.node.local.testorg.messaging.mesh.config.create",
			Msg:        any,
		})

		assert.NoError(t, err)
		// Assert
		configStore.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("UpdateNodeConfigFailed", func(t *testing.T) {
		/* Node Cretaed event */
		evt := eCfgPb.NodeConfigUpdateEvent{
			NodeId: testNode.String(),
			Commit: "commithash",
			Status: eCfgPb.UpdateStatus_Failed,
		}

		any, err := anypb.New(&evt)
		assert.NoError(t, err)
		cfg := db.Configuration{
			NodeId:          "node-id",
			State:           db.Published,
			Commit:          db.Commit{Hash: "commit"},
			LastCommitState: db.Published,
			LastCommit:      db.Commit{Hash: "lastcommit"},
		}

		cmt := db.Commit{Model: gorm.Model{ID: 2}, Hash: "commithash"}
		configRepo.On("Get", testNode.String()).Return(&cfg, nil).Once()

		commitRepo.On("Get", evt.Commit).Return(&cmt, nil).Once()

		configRepo.On("UpdateLastCommit", mock.Anything, mock.MatchedBy(func(a *db.CommitState) bool {
			return *a == db.CommitState(evt.Status)
		})).Return(nil).Once()

		_, err = eventServer.EventNotification(context.Background(), &epb.Event{
			RoutingKey: "event.node.local.testorg.messaging.mesh.config.create",
			Msg:        any,
		})

		assert.NoError(t, err)
		// Assert
		configStore.AssertExpectations(t)
		assert.NoError(t, err)
	})
}
