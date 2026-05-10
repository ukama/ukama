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
	"github.com/stretchr/testify/require"
	evt "github.com/ukama/ukama/systems/common/events"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	hpb "github.com/ukama/ukama/systems/node/health/pb/gen"
	healthmocks "github.com/ukama/ukama/systems/node/health/pb/gen/mocks"
	"github.com/ukama/ukama/systems/node/software/mocks"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
	"google.golang.org/protobuf/types/known/anypb"
)

// Event test routing keys (must match server logic).
var (
	routeNodeAppChunkReady = msgbus.PrepareRoute(testOrgName, evt.NodeEventToEventConfig[evt.NodeAppChunkReady].RoutingKey)
	routeNodeOnline        = msgbus.PrepareRoute(testOrgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline])
)

func newEventServerWithMocks(t *testing.T, sRepo *mocks.SoftwareRepo, appRepo *mocks.AppRepo) *SoftwareUpdateEventServer {
	t.Helper()
	nodeRepo := mocks.NewNodeRepo(t)
	nodeRepo.On("Create", mock.Anything).Return(nil).Maybe()
	nodeRepo.On("List").Return([]db.Node{}, nil).Maybe()
	swServer := NewSoftwareServer(testOrgName, sRepo, appRepo, nodeRepo, nil, mbmocks.NewMsgBusServiceClient(t), false, []string{"192.168.0.1"})
	return NewSoftwareEventServer(testOrgName, swServer)
}

func newEventServerWithAllMocks(t *testing.T, sRepo *mocks.SoftwareRepo, appRepo *mocks.AppRepo, nodeRepo *mocks.NodeRepo, healthProvider *mocks.HealthClientProvider) *SoftwareUpdateEventServer {
	t.Helper()
	swServer := NewSoftwareServer(testOrgName, sRepo, appRepo, nodeRepo, healthProvider, mbmocks.NewMsgBusServiceClient(t), false, []string{"192.168.0.1"})
	return NewSoftwareEventServer(testOrgName, swServer)
}

// mustMarshalNodeOnlineEvent builds an epb.Event for NodeOnlineEvent.
func mustMarshalNodeOnlineEvent(t *testing.T, nodeId string) *epb.Event {
	t.Helper()
	msg, err := anypb.New(&epb.NodeOnlineEvent{NodeId: nodeId})
	require.NoError(t, err)
	return &epb.Event{RoutingKey: routeNodeOnline, Msg: msg}
}

// mustMarshalAppChunkReadyEvent builds an epb.Event for EventArtifactChunkReady.
func mustMarshalAppChunkReadyEvent(t *testing.T, appName, version string) *epb.Event {
	t.Helper()
	msg, err := anypb.New(&epb.EventArtifactChunkReady{Name: appName, Version: version})
	require.NoError(t, err)
	return &epb.Event{RoutingKey: routeNodeAppChunkReady, Msg: msg}
}

func TestEventNotification(t *testing.T) {
	ctx := context.Background()

	t.Run("unknown_routing_key", func(t *testing.T) {
		s := newEventServerWithMocks(t, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t))
		e := &epb.Event{RoutingKey: "unknown.route.key", Msg: nil}

		resp, err := s.EventNotification(ctx, e)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotNil(t, resp)
	})

	t.Run("node_online_creates_software_when_none_exist", func(t *testing.T) {
		s := newEventServerWithMocks(t, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t))
		e := mustMarshalNodeOnlineEvent(t, testNodeId)

		resp, err := s.EventNotification(ctx, e)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("node_online_skips_create_when_software_exists", func(t *testing.T) {
		s := newEventServerWithMocks(t, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t))
		e := mustMarshalNodeOnlineEvent(t, testNodeId)

		resp, err := s.EventNotification(ctx, e)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("node_online_unmarshal_error", func(t *testing.T) {
		s := newEventServerWithMocks(t, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t))
		e := &epb.Event{RoutingKey: routeNodeOnline, Msg: nil}

		resp, err := s.EventNotification(ctx, e)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("node_online_list_error", func(t *testing.T) {
		s := newEventServerWithMocks(t, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t))
		e := mustMarshalNodeOnlineEvent(t, testNodeId)

		resp, err := s.EventNotification(ctx, e)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("node_online_get_all_apps_error", func(t *testing.T) {
		s := newEventServerWithMocks(t, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t))
		e := mustMarshalNodeOnlineEvent(t, testNodeId)

		resp, err := s.EventNotification(ctx, e)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("node_online_create_software_error", func(t *testing.T) {
		s := newEventServerWithMocks(t, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t))
		e := mustMarshalNodeOnlineEvent(t, testNodeId)

		resp, err := s.EventNotification(ctx, e)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("node_online_list_record_not_found_continues", func(t *testing.T) {
		s := newEventServerWithMocks(t, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t))
		e := mustMarshalNodeOnlineEvent(t, testNodeId)

		resp, err := s.EventNotification(ctx, e)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("node_app_chunk_ready_updates_softwares", func(t *testing.T) {
		sw := dbSoftwareFixture()
		softwares := []*db.Software{sw}
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", "", ukama.Unknown, testAppNameForUpdate).Return(softwares, nil)
		sRepo.On("Update", mock.MatchedBy(func(s *db.Software) bool {
			return s.DesiredVersion == testTagVersion && s.Status == ukama.UpdateAvailable &&
				len(s.ChangeLogs) > 0
		})).Return(nil)

		s := newEventServerWithMocks(t, sRepo, mocks.NewAppRepo(t))
		e := mustMarshalAppChunkReadyEvent(t, testAppNameForUpdate, testTagVersion)

		resp, err := s.EventNotification(ctx, e)

		require.NoError(t, err)
		require.NotNil(t, resp)
		sRepo.AssertExpectations(t)
	})

	t.Run("node_app_chunk_ready_unmarshal_error", func(t *testing.T) {
		s := newEventServerWithMocks(t, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t))
		e := &epb.Event{RoutingKey: routeNodeAppChunkReady, Msg: nil}

		resp, err := s.EventNotification(ctx, e)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("node_app_chunk_ready_list_error", func(t *testing.T) {
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", "", ukama.Unknown, testAppNameForUpdate).Return(nil, errors.New("list failed"))

		s := newEventServerWithMocks(t, sRepo, mocks.NewAppRepo(t))
		e := mustMarshalAppChunkReadyEvent(t, testAppNameForUpdate, testTagVersion)

		resp, err := s.EventNotification(ctx, e)

		assert.Error(t, err)
		assert.Nil(t, resp)
		sRepo.AssertExpectations(t)
	})

	t.Run("node_app_chunk_ready_update_error", func(t *testing.T) {
		sw := dbSoftwareFixture()
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", "", ukama.Unknown, testAppNameForUpdate).Return([]*db.Software{sw}, nil)
		sRepo.On("Update", mock.Anything).Return(errors.New("update failed"))

		s := newEventServerWithMocks(t, sRepo, mocks.NewAppRepo(t))
		e := mustMarshalAppChunkReadyEvent(t, testAppNameForUpdate, testTagVersion)

		resp, err := s.EventNotification(ctx, e)

		assert.Error(t, err)
		assert.Nil(t, resp)
		sRepo.AssertExpectations(t)
	})

	t.Run("node_app_chunk_ready_empty_list_succeeds", func(t *testing.T) {
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", "", ukama.Unknown, "nonexistent-app").Return([]*db.Software{}, nil)

		s := newEventServerWithMocks(t, sRepo, mocks.NewAppRepo(t))
		e := mustMarshalAppChunkReadyEvent(t, "nonexistent-app", "1.0.0")

		resp, err := s.EventNotification(ctx, e)

		require.NoError(t, err)
		require.NotNil(t, resp)
		sRepo.AssertExpectations(t)
	})
}

func TestHandleNodeAppChunkReadyEvent_VersionCompareFallback(t *testing.T) {
	ctx := context.Background()
	routeNodeAppChunkReady := msgbus.PrepareRoute(testOrgName, evt.NodeEventToEventConfig[evt.NodeAppChunkReady].RoutingKey)

	sRepo := mocks.NewSoftwareRepo(t)
	nodeRepo := mocks.NewNodeRepo(t)
	appRepo := mocks.NewAppRepo(t)
	healthProvider := mocks.NewHealthClientProvider(t)

	nodeRepo.On("List").Return([]db.Node{}, nil).Once()

	sw := dbSoftwareFixture()
	sw.AppName = testAppNameForUpdate
	sw.CurrentVersion = "old_arch-40ab36dc2"
	sw.ChangeLogs = []string{}

	sRepo.On("List", "", ukama.Unknown, testAppNameForUpdate).Return([]*db.Software{sw}, nil).Once()
	sRepo.On("Update", mock.MatchedBy(func(s *db.Software) bool {
		return s.Id == sw.Id &&
			s.DesiredVersion == testTagVersion &&
			s.Status == ukama.UpdateAvailable &&
			len(s.ChangeLogs) == 1
	})).Return(nil).Once()

	s := newEventServerWithAllMocks(t, sRepo, appRepo, nodeRepo, healthProvider)
	msg, err := anypb.New(&epb.EventArtifactChunkReady{Name: testAppNameForUpdate, Version: testTagVersion})
	require.NoError(t, err)
	e := &epb.Event{RoutingKey: routeNodeAppChunkReady, Msg: msg}

	resp, err := s.EventNotification(ctx, e)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestReconcileCurrentAppVersion_UpdatesCurrentVersionFromHealth(t *testing.T) {
	sRepo := mocks.NewSoftwareRepo(t)
	nodeRepo := mocks.NewNodeRepo(t)
	appRepo := mocks.NewAppRepo(t)
	healthProvider := mocks.NewHealthClientProvider(t)
	healthClient := healthmocks.NewHealhtServiceClient(t)

	nodeRepo.On("List").Return([]db.Node{
		{NodeId: testNodeId},
	}, nil).Once()

	healthProvider.On("GetClient").Return(healthClient, nil)
	healthClient.On("ListApps", mock.Anything, mock.MatchedBy(func(req *hpb.ListAppsRequest) bool {
		return req.NodeId == testNodeId && req.Name == ""
	})).Return(&hpb.ListAppsResponse{
		Capps: []*hpb.Capps{
			{Name: testAppNameForUpdate, Tag: testTagVersion},
		},
	}, nil).Once()

	sw := dbSoftwareFixture()
	sw.NodeId = testNodeId
	sw.AppName = testAppNameForUpdate
	sw.CurrentVersion = "1.0.0"
	sw.ChangeLogs = []string{"keep-me"}

	sRepo.On("List", testNodeId, ukama.Unknown, testAppNameForUpdate).Return([]*db.Software{sw}, nil).Once()
	sRepo.On("Update", mock.MatchedBy(func(updated *db.Software) bool {
		return updated.Id == sw.Id &&
			updated.CurrentVersion == testTagVersion &&
			len(updated.ChangeLogs) == 1 &&
			updated.ChangeLogs[0] == "keep-me"
	})).Return(nil).Once()

	s := newEventServerWithAllMocks(t, sRepo, appRepo, nodeRepo, healthProvider)
	err := s.reconcileCurrentAppVersion()
	assert.NoError(t, err)
}
