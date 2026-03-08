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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/software/mocks"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ========== Shared test data (declare once, use in all tests) ==========

const (
	testOrgName          = "test-org"
	testAppName          = "test-app"
	testAppSpace         = "default-space"
	testAppNotes         = "test notes"
	testNodeId           = "UK-SA2156-HNODE-A1-XXXX" // valid 23-char node ID; server normalizes to lowercase
	testNodeIdNormalized = "uk-sa2156-hnode-a1-xxxx" // value passed to repo after ValidateNodeId().String()
	testAppNameForUpdate = "myapp"
	testTagVersion       = "1.2.0"
	testCurrentVersion   = "1.0.0"
	testDesiredVersion   = "1.2.0"
)

const errMsgDB = "db error"

var (
	testMetricsKeys    = []string{"cpu", "memory"}
	testCreateAppReq   = &pb.CreateAppRequest{Name: testAppName, Space: testAppSpace, Notes: testAppNotes, MetricsKeys: testMetricsKeys}
	testGetAppListReq  = &pb.GetAppListRequest{}
	successCreateMsg   = "App created successfully"
	successUpdateMsg  = "Software updated successfully"
	invalidVersionMsg = "Invalid software version provided"
	alreadyUpToDateMsg = "Software is already up to date"
)

// ========== Helpers to build server with mocks ==========

func newTestServer(sRepo *mocks.SoftwareRepo, appRepo *mocks.AppRepo, msgBus *mbmocks.MsgBusServiceClient) *SoftwareServer {
	return NewSoftwareServer(testOrgName, sRepo, appRepo, msgBus, false)
}

func dbAppFixture() db.App {
	return db.App{
		Id:          uuid.NewV4(),
		Name:        testAppName,
		Space:       testAppSpace,
		Notes:       testAppNotes,
		MetricsKeys: testMetricsKeys,
	}
}

func dbSoftwareFixture() *db.Software {
	releaseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	return &db.Software{
		Id:             uuid.NewV4(),
		NodeId:         testNodeId,
		AppName:        testAppNameForUpdate,
		App:            db.App{Name: testAppNameForUpdate, Space: testAppSpace, Notes: testAppNotes, MetricsKeys: testMetricsKeys},
		ChangeLogs:     []string{"initial"},
		CurrentVersion: testCurrentVersion,
		DesiredVersion: testDesiredVersion,
		ReleaseDate:    releaseDate,
		CreatedAt:      now,
		UpdatedAt:      now,
		Status:         ukama.UpdateAvailable,
	}
}

// ========== CreateApp ==========

func TestCreateApp(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		appRepo := mocks.NewAppRepo(t)
		appRepo.On("Create", mock.MatchedBy(func(a db.App) bool {
			return a.Name == testAppName && a.Space == testAppSpace && a.Notes == testAppNotes &&
				len(a.MetricsKeys) == len(testMetricsKeys)
		})).Return(nil)

		s := newTestServer(mocks.NewSoftwareRepo(t), appRepo, mbmocks.NewMsgBusServiceClient(t))
		resp, err := s.CreateApp(ctx, testCreateAppReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, successCreateMsg, resp.Message)
		appRepo.AssertExpectations(t)
	})

	t.Run("repo_error", func(t *testing.T) {
		appRepo := mocks.NewAppRepo(t)
		appRepo.On("Create", mock.Anything).Return(errors.New(errMsgDB))

		s := newTestServer(mocks.NewSoftwareRepo(t), appRepo, mbmocks.NewMsgBusServiceClient(t))
		resp, err := s.CreateApp(ctx, testCreateAppReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		appRepo.AssertExpectations(t)
	})
}

// ========== GetAppList ==========

func TestGetAppList(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		apps := []db.App{dbAppFixture()}
		appRepo := mocks.NewAppRepo(t)
		appRepo.On("GetAll").Return(apps, nil)

		s := newTestServer(mocks.NewSoftwareRepo(t), appRepo, mbmocks.NewMsgBusServiceClient(t))
		resp, err := s.GetAppList(ctx, testGetAppListReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Apps, 1)
		assert.Equal(t, testAppName, resp.Apps[0].Name)
		assert.Equal(t, testAppSpace, resp.Apps[0].Space)
		assert.Equal(t, testAppNotes, resp.Apps[0].Notes)
		assert.Equal(t, testMetricsKeys, resp.Apps[0].MetricsKeys)
		appRepo.AssertExpectations(t)
	})

	t.Run("empty_list", func(t *testing.T) {
		appRepo := mocks.NewAppRepo(t)
		appRepo.On("GetAll").Return([]db.App{}, nil)

		s := newTestServer(mocks.NewSoftwareRepo(t), appRepo, mbmocks.NewMsgBusServiceClient(t))
		resp, err := s.GetAppList(ctx, testGetAppListReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Empty(t, resp.Apps)
		appRepo.AssertExpectations(t)
	})

	t.Run("repo_error", func(t *testing.T) {
		appRepo := mocks.NewAppRepo(t)
		appRepo.On("GetAll").Return(nil, errors.New(errMsgDB))

		s := newTestServer(mocks.NewSoftwareRepo(t), appRepo, mbmocks.NewMsgBusServiceClient(t))
		resp, err := s.GetAppList(ctx, testGetAppListReq)

		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		appRepo.AssertExpectations(t)
	})
}

// ========== GetSoftwareList ==========

func TestGetSoftwareList(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		sw := dbSoftwareFixture()
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.SoftwareStatusType(0), "").Return([]*db.Software{sw}, nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mbmocks.NewMsgBusServiceClient(t))
		req := &pb.GetSoftwareListRequest{NodeId: testNodeId}
		resp, err := s.GetSoftwareList(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Software, 1)
		assert.Equal(t, sw.Id.String(), resp.Software[0].Id)
		assert.Equal(t, testNodeId, resp.Software[0].NodeId)
		assert.Equal(t, testAppNameForUpdate, resp.Software[0].Name)
		assert.Equal(t, testCurrentVersion, resp.Software[0].CurrentVersion)
		assert.Equal(t, testDesiredVersion, resp.Software[0].DesiredVersion)
		sRepo.AssertExpectations(t)
	})

	t.Run("invalid_node_id", func(t *testing.T) {
		s := newTestServer(mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t), mbmocks.NewMsgBusServiceClient(t))
		req := &pb.GetSoftwareListRequest{NodeId: "short"}
		resp, err := s.GetSoftwareList(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("repo_error", func(t *testing.T) {
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", "", ukama.SoftwareStatusType(0), "").Return(nil, errors.New(errMsgDB))

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mbmocks.NewMsgBusServiceClient(t))
		req := &pb.GetSoftwareListRequest{}
		resp, err := s.GetSoftwareList(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		sRepo.AssertExpectations(t)
	})
}

// ========== UpdateSoftware ==========

func TestUpdateSoftware(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		sw := dbSoftwareFixture()
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).Return([]*db.Software{sw}, nil)
		sRepo.On("Update", mock.MatchedBy(func(s *db.Software) bool {
			return s.CurrentVersion == testTagVersion && s.Status == ukama.UpToDate &&
				len(s.ChangeLogs) > 0
		})).Return(nil)
		msgBus := mbmocks.NewMsgBusServiceClient(t)
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), msgBus)
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, successUpdateMsg, resp.Message)
		sRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("invalid_node_id", func(t *testing.T) {
		s := newTestServer(mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t), mbmocks.NewMsgBusServiceClient(t))
		req := &pb.UpdateSoftwareRequest{NodeId: "bad", Name: testAppNameForUpdate, Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("invalid_tag", func(t *testing.T) {
		s := newTestServer(mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t), mbmocks.NewMsgBusServiceClient(t))
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate, Tag: "not-a-version"}
		resp, err := s.UpdateSoftware(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("not_found", func(t *testing.T) {
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).Return([]*db.Software{}, nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mbmocks.NewMsgBusServiceClient(t))
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
		sRepo.AssertExpectations(t)
	})

	t.Run("version_mismatch", func(t *testing.T) {
		sw := dbSoftwareFixture()
		sw.DesiredVersion = "2.0.0"
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).Return([]*db.Software{sw}, nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mbmocks.NewMsgBusServiceClient(t))
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, invalidVersionMsg, resp.Message)
		sRepo.AssertExpectations(t)
	})

	t.Run("already_up_to_date", func(t *testing.T) {
		sw := dbSoftwareFixture()
		sw.CurrentVersion = testTagVersion
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).Return([]*db.Software{sw}, nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mbmocks.NewMsgBusServiceClient(t))
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, alreadyUpToDateMsg, resp.Message)
		sRepo.AssertExpectations(t)
	})

	t.Run("publish_error", func(t *testing.T) {
		sw := dbSoftwareFixture()
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).Return([]*db.Software{sw}, nil)
		msgBus := mbmocks.NewMsgBusServiceClient(t)
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(errors.New("publish failed"))

		s := newTestServer(sRepo, mocks.NewAppRepo(t), msgBus)
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		sRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("update_repo_error", func(t *testing.T) {
		sw := dbSoftwareFixture()
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).Return([]*db.Software{sw}, nil)
		sRepo.On("Update", mock.Anything).Return(errors.New("db update failed"))
		msgBus := mbmocks.NewMsgBusServiceClient(t)
		msgBus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), msgBus)
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		sRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})
}
