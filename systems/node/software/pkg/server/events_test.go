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

// Routing keys (must match server logic).
var (
	routeChunkReady = msgbus.PrepareRoute(testOrgName, evt.NodeEventToEventConfig[evt.NodeAppChunkReady].RoutingKey)
	routeNodeOnline = msgbus.PrepareRoute(testOrgName, evt.NodeStateEventRoutingKey[evt.NodeStateEventOnline])
	routeUploaded   = msgbus.PrepareRoute(testOrgName, "event.cloud.global.{{ .Org}}.hub.artifactmanager.app.uploaded")
)

// newEventServer builds an event server with permissive reconcile-path mocks
// (health returns no apps, so node-online reconcile is a no-op).
func newEventServer(t *testing.T, sRepo *mocks.SoftwareRepo, releaseRepo *mocks.ReleaseRepo) *SoftwareUpdateEventServer {
	t.Helper()

	nodeRepo := mocks.NewNodeRepo(t)
	nodeRepo.On("Create", mock.Anything).Return(nil).Maybe()

	appRepo := mocks.NewAppRepo(t)
	appRepo.On("GetAll").Return([]db.App{}, nil).Maybe()

	health := mocks.NewHealthClientProvider(t)
	hc := healthmocks.NewHealthServiceClient(t)
	health.On("GetClient").Return(hc, nil).Maybe()
	hc.On("ListApps", mock.Anything, mock.Anything).Return(&hpb.ListAppsResponse{}, nil).Maybe()

	swServer := NewSoftwareServer(testOrgName, sRepo, appRepo, nodeRepo, releaseRepo, nil, health,
		mbmocks.NewMsgBusServiceClient(t), false, []string{testNodeGwIP}, nil, nil, 0, 0)
	return NewSoftwareEventServer(testOrgName, swServer)
}

func uploadedEvent(t *testing.T, name, version string) *epb.Event {
	t.Helper()
	msg, err := anypb.New(&epb.EventArtifactUploaded{Name: name, Version: version})
	require.NoError(t, err)
	return &epb.Event{RoutingKey: routeUploaded, Msg: msg}
}

func chunkReadyEvent(t *testing.T, name, version string) *epb.Event {
	t.Helper()
	msg, err := anypb.New(&epb.EventArtifactChunkReady{Name: name, Version: version})
	require.NoError(t, err)
	return &epb.Event{RoutingKey: routeChunkReady, Msg: msg}
}

func nodeOnlineEvent(t *testing.T, nodeId string) *epb.Event {
	t.Helper()
	msg, err := anypb.New(&epb.NodeOnlineEvent{NodeId: nodeId})
	require.NoError(t, err)
	return &epb.Event{RoutingKey: routeNodeOnline, Msg: msg}
}

func TestEventNotification(t *testing.T) {
	ctx := context.Background()

	t.Run("unknown_routing_key", func(t *testing.T) {
		s := newEventServer(t, mocks.NewSoftwareRepo(t), mocks.NewReleaseRepo(t))
		resp, err := s.EventNotification(ctx, &epb.Event{RoutingKey: "unknown.route", Msg: nil})
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	// uploaded records availability only — no node state is touched.
	t.Run("artifact_uploaded_records_catalog", func(t *testing.T) {
		releaseRepo := mocks.NewReleaseRepo(t)
		releaseRepo.On("Upsert", mock.MatchedBy(func(r *db.ReleaseCatalog) bool {
			return r.Name == testAppNameForUpdate && r.Version == testTagVersion && r.Available
		})).Return(nil)

		s := newEventServer(t, mocks.NewSoftwareRepo(t), releaseRepo)
		resp, err := s.EventNotification(ctx, uploadedEvent(t, testAppNameForUpdate, testTagVersion))

		require.NoError(t, err)
		require.NotNil(t, resp)
		releaseRepo.AssertExpectations(t)
	})

	t.Run("artifact_uploaded_unmarshal_error", func(t *testing.T) {
		s := newEventServer(t, mocks.NewSoftwareRepo(t), mocks.NewReleaseRepo(t))
		resp, err := s.EventNotification(ctx, &epb.Event{RoutingKey: routeUploaded, Msg: nil})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	// chunkready only flags the catalog row chunked — it does not set desired versions.
	t.Run("chunk_ready_marks_chunked", func(t *testing.T) {
		releaseRepo := mocks.NewReleaseRepo(t)
		releaseRepo.On("Upsert", mock.Anything).Return(nil)
		releaseRepo.On("SetChunked", testAppNameForUpdate, "app", testTagVersion).Return(nil)

		s := newEventServer(t, mocks.NewSoftwareRepo(t), releaseRepo)
		resp, err := s.EventNotification(ctx, chunkReadyEvent(t, testAppNameForUpdate, testTagVersion))

		require.NoError(t, err)
		require.NotNil(t, resp)
		releaseRepo.AssertExpectations(t)
	})

	t.Run("chunk_ready_unmarshal_error", func(t *testing.T) {
		s := newEventServer(t, mocks.NewSoftwareRepo(t), mocks.NewReleaseRepo(t))
		resp, err := s.EventNotification(ctx, &epb.Event{RoutingKey: routeChunkReady, Msg: nil})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	// node online creates the node and reconciles synchronously (no 120s delay);
	// with no installed apps and no desired release it is a no-op that still succeeds.
	t.Run("node_online_succeeds", func(t *testing.T) {
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", mock.Anything, ukama.Unknown, "").Return([]*db.Software{}, nil).Maybe()

		s := newEventServer(t, sRepo, mocks.NewReleaseRepo(t))
		resp, err := s.EventNotification(ctx, nodeOnlineEvent(t, testNodeId))

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("node_online_unmarshal_error", func(t *testing.T) {
		s := newEventServer(t, mocks.NewSoftwareRepo(t), mocks.NewReleaseRepo(t))
		resp, err := s.EventNotification(ctx, &epb.Event{RoutingKey: routeNodeOnline, Msg: nil})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
