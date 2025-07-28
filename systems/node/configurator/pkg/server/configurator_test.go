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

	"github.com/ukama/ukama/systems/node/configurator/mocks"
	"github.com/ukama/ukama/systems/node/configurator/pkg"
	"github.com/ukama/ukama/systems/node/configurator/pkg/db"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	pb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
)

const testOrgName = "testOrg"

func TestConfiguratorServer_ConfigEvent(t *testing.T) {
	// Arrange
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	commitRepo := &mocks.CommitRepo{}
	configRepo := &mocks.ConfigRepo{}
	configStore := &mocks.ConfigStoreProvider{}

	s := NewConfiguratorServer(msgbusClient, configRepo, commitRepo, configStore, testOrgName, pkg.IsDebugMode)

	configStore.On("HandleConfigStoreEvent", mock.Anything, mock.Anything).Return(nil).Once()

	_, err := s.ConfigEvent(context.Background(), &pb.ConfigStoreEvent{})
	assert.NoError(t, err)
	// Assert
	configStore.AssertExpectations(t)
	assert.NoError(t, err)

}

func TestConfiguratorServer_GetConfigVersion(t *testing.T) {
	// Arrange
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	commitRepo := &mocks.CommitRepo{}
	configRepo := &mocks.ConfigRepo{}
	configStore := &mocks.ConfigStoreProvider{}

	s := NewConfiguratorServer(msgbusClient, configRepo, commitRepo, configStore, testOrgName, pkg.IsDebugMode)

	configRepo.On("Get", mock.AnythingOfType("string")).Return(&db.Configuration{
		NodeId:          "node-id",
		State:           db.Published,
		Commit:          db.Commit{Hash: "commit"},
		LastCommitState: db.Published,
		LastCommit:      db.Commit{Hash: "lastcommit"},
	}, nil).Once()

	_, err := s.GetConfigVersion(context.Background(), &pb.ConfigVersionRequest{NodeId: "node-id"})
	assert.NoError(t, err)
	// Assert
	configStore.AssertExpectations(t)
	assert.NoError(t, err)

}
func TestConfiguratorServer_ApplyConfig(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	commitRepo := &mocks.CommitRepo{}
	configRepo := &mocks.ConfigRepo{}
	configStore := &mocks.ConfigStoreProvider{}

	s := NewConfiguratorServer(msgbusClient, configRepo, commitRepo, configStore, testOrgName, pkg.IsDebugMode)
	req := &pb.ApplyConfigRequest{Hash: "4f6e609"}
	configStore.On("HandleConfigCommitReq", mock.Anything, req.Hash).Return(nil).Once()

	_, err := s.ApplyConfig(context.Background(), req)
	assert.NoError(t, err)
	// Assert
	configStore.AssertExpectations(t)
	assert.NoError(t, err)

}
