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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/rest/client"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	healthpb "github.com/ukama/ukama/systems/node/health/pb/gen"
	opmonpb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
	"github.com/ukama/ukama/systems/node/software/mocks"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// ========== Shared test data (declare once, use in all tests) ==========

const (
	testOrgName          = "test-org"
	testAppName          = "test-app"
	testAppSpace         = "default-space"
	testAppNotes         = "test notes"
	testNodeId           = "UK-SA2156-HNODE-A1-XXXX" // the server lowercases it
	testNodeIdNormalized = "uk-sa2156-hnode-a1-xxxx" // what the repo receives
	testAppNameForUpdate = "myapp"
	testTagVersion       = "1.2.0"
	testCurrentVersion   = "1.0.0"
	testDesiredVersion   = "1.2.0"
	testNodeGwIP         = "192.168.0.1"
	testHubURL           = "http://192.168.0.1:8080"
	testStarterPath      = "/starter/v1/update"
	testSoftwareRoute    = "request.cloud.local.test-org.node.software.nodefeeder.publish"
	testTarget           = "test-org...uk-sa2156-hnode-a1-xxxx"
)

const errMsgDB = "db error"

var (
	testMetricsKeys    = []string{"cpu", "memory"}
	testGetAppListReq  = &pb.GetAppListRequest{}
	successCreateMsg   = "App created successfully"
	successUpdateMsg   = "Software update dispatched"
	alreadyUpToDateMsg = "Software is already up to date"
)

var testCreateAppReq = &pb.CreateAppRequest{Name: testAppName, Space: testAppSpace,
	Notes: testAppNotes, MetricsKeys: testMetricsKeys}

// ========== Helpers to build server with mocks ==========

func newTestServer(sRepo *mocks.SoftwareRepo, appRepo *mocks.AppRepo,
	nodeRepo *mocks.NodeRepo, msgBus *mbmocks.MsgBusServiceClient) *SoftwareServer {
	nodeClient := &mbmocks.NodeClient{}
	nodeClient.On("Get", mock.Anything).
		Return(&creg.NodeInfo{Status: creg.NodeStatusInfo{Connectivity: "Online"}}, nil).Maybe()

	s := NewSoftwareServer(testOrgName, sRepo, appRepo, nodeRepo, nil, nil, nil, msgBus,
		false, []string{testNodeGwIP},
		nil,
		nil,
		0,
		0,
		nodeClient,
	)
	s.startWatch = func(uuid.UUID, string, string, string) {}
	return s
}

type watchCall struct {
	recordID uuid.UUID
	nodeID   string
	appName  string
	version  string
}

func recordWatches(s *SoftwareServer) *[]watchCall {
	calls := &[]watchCall{}
	s.startWatch = func(recordID uuid.UUID, nodeID, appName, version string) {
		*calls = append(*calls, watchCall{recordID, nodeID, appName, version})
	}
	return calls
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
		ChangeLogs:     []string{"initial"},
		CurrentVersion: testCurrentVersion,
		DesiredVersion: testDesiredVersion,
		ReleaseDate:    releaseDate,
		CreatedAt:      now,
		UpdatedAt:      now,
		Status:         ukama.UpdateAvailable,
		App: db.App{Name: testAppNameForUpdate, Space: testAppSpace, Notes: testAppNotes,
			MetricsKeys: testMetricsKeys},
	}
}

func matchSoftwareRoute() interface{} {
	return mock.MatchedBy(func(route string) bool {
		return route == testSoftwareRoute
	})
}

func matchStarterUpdateMessage(t *testing.T) interface{} {
	return mock.MatchedBy(func(msg interface{}) bool {
		m, ok := msg.(*epb.NodeFeederMessage)
		if !ok {
			t.Logf("unexpected msgbus message type: %T", msg)
			return false
		}

		if m.Target != testTarget {
			t.Logf("unexpected target: got=%s want=%s", m.Target, testTarget)
			return false
		}

		if m.HttpMethod != "POST" {
			t.Logf("unexpected method: got=%s want=POST", m.HttpMethod)
			return false
		}

		if m.Path != testStarterPath {
			t.Logf("unexpected path: got=%s want=%s", m.Path, testStarterPath)
			return false
		}

		if m.NodeId != testNodeIdNormalized {
			t.Logf("unexpected node id: got=%s want=%s", m.NodeId, testNodeIdNormalized)
			return false
		}

		var body struct {
			Name string   `json:"name"`
			Tag  string   `json:"tag"`
			Hub  []string `json:"hub"`
		}
		if err := json.Unmarshal(m.Msg, &body); err != nil {
			t.Logf("failed to decode message body: %v", err)
			return false
		}

		if body.Name != testAppNameForUpdate {
			t.Logf("unexpected body.name: got=%s want=%s", body.Name, testAppNameForUpdate)
			return false
		}

		if body.Tag != testTagVersion {
			t.Logf("unexpected body.tag: got=%s want=%s", body.Tag, testTagVersion)
			return false
		}

		if len(body.Hub) != 1 || body.Hub[0] != testHubURL {
			t.Logf("unexpected body.hub: got=%v want=[%s]", body.Hub, testHubURL)
			return false
		}

		var raw map[string]json.RawMessage
		if err := json.Unmarshal(m.Msg, &raw); err != nil {
			t.Logf("failed to decode raw message body: %v", err)
			return false
		}
		if _, ok := raw["host"]; ok {
			t.Logf("body must not contain deprecated host field")
			return false
		}

		return true
	})
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

		s := newTestServer(mocks.NewSoftwareRepo(t), appRepo, mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
		resp, err := s.CreateApp(ctx, testCreateAppReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, successCreateMsg, resp.Message)
		appRepo.AssertExpectations(t)
	})

	t.Run("repo_error", func(t *testing.T) {
		appRepo := mocks.NewAppRepo(t)
		appRepo.On("Create", mock.Anything).Return(errors.New(errMsgDB))

		s := newTestServer(mocks.NewSoftwareRepo(t), appRepo, mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
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

		s := newTestServer(mocks.NewSoftwareRepo(t), appRepo, mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
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

		s := newTestServer(mocks.NewSoftwareRepo(t), appRepo, mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
		resp, err := s.GetAppList(ctx, testGetAppListReq)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Empty(t, resp.Apps)

		appRepo.AssertExpectations(t)
	})

	t.Run("repo_error", func(t *testing.T) {
		appRepo := mocks.NewAppRepo(t)
		appRepo.On("GetAll").Return(nil, errors.New(errMsgDB))

		s := newTestServer(mocks.NewSoftwareRepo(t), appRepo, mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
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
		sRepo.On("List", testNodeIdNormalized, ukama.SoftwareStatusType(0), "").
			Return([]*db.Software{sw}, nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
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
		s := newTestServer(mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
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

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
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
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).
			Return([]*db.Software{sw}, nil)
		sRepo.On("Update", mock.MatchedBy(func(s *db.Software) bool {
			return s.CurrentVersion == testCurrentVersion &&
				s.Status == ukama.UpdateInProgress &&
				len(s.ChangeLogs) > 0
		})).Return(nil)

		msgBus := mbmocks.NewMsgBusServiceClient(t)
		msgBus.On("PublishRequest",
			matchSoftwareRoute(),
			matchStarterUpdateMessage(t),
		).Return(nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t), msgBus)
		watches := recordWatches(s)
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate,
			Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, successUpdateMsg, resp.Message)
		assert.Equal(t, "node:"+testNodeIdNormalized, resp.ResourceKey)
		assert.Equal(t, copr.StatusRunning, resp.Status)
		assert.Equal(t, []watchCall{{sw.Id, testNodeIdNormalized, testAppNameForUpdate,
			testTagVersion}}, *watches)

		sRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})

	t.Run("invalid_node_id", func(t *testing.T) {
		s := newTestServer(mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
		req := &pb.UpdateSoftwareRequest{NodeId: "bad", Name: testAppNameForUpdate,
			Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("invalid_tag", func(t *testing.T) {
		sw := dbSoftwareFixture()
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).
			Return([]*db.Software{sw}, nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate,
			Tag: "not-a-version"}
		resp, err := s.UpdateSoftware(ctx, req)

		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))

		sRepo.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).
			Return([]*db.Software{}, nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate,
			Tag: testTagVersion}
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
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).
			Return([]*db.Software{sw}, nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
		watches := recordWatches(s)
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate,
			Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		assert.Nil(t, resp)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		assert.Contains(t, status.Convert(err).Message(), "2.0.0")
		assert.Empty(t, *watches)

		sRepo.AssertExpectations(t)
	})

	t.Run("already_up_to_date", func(t *testing.T) {
		sw := dbSoftwareFixture()
		sw.CurrentVersion = testTagVersion

		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).
			Return([]*db.Software{sw}, nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
			mbmocks.NewMsgBusServiceClient(t))
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate,
			Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, alreadyUpToDateMsg, resp.Message)

		sRepo.AssertExpectations(t)
	})

	t.Run("publish_error", func(t *testing.T) {
		sw := dbSoftwareFixture()
		sRepo := mocks.NewSoftwareRepo(t)
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).
			Return([]*db.Software{sw}, nil)

		msgBus := mbmocks.NewMsgBusServiceClient(t)
		msgBus.On("PublishRequest",
			matchSoftwareRoute(),
			matchStarterUpdateMessage(t),
		).Return(errors.New("publish failed"))

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t), msgBus)
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate,
			Tag: testTagVersion}
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
		sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).
			Return([]*db.Software{sw}, nil)
		sRepo.On("Update", mock.Anything).Return(errors.New("db update failed"))

		msgBus := mbmocks.NewMsgBusServiceClient(t)
		msgBus.On("PublishRequest",
			matchSoftwareRoute(),
			matchStarterUpdateMessage(t),
		).Return(nil)

		s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t), msgBus)
		watches := recordWatches(s)
		req := &pb.UpdateSoftwareRequest{NodeId: testNodeId, Name: testAppNameForUpdate,
			Tag: testTagVersion}
		resp, err := s.UpdateSoftware(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, successUpdateMsg, resp.Message)
		assert.Equal(t, copr.StatusRunning, resp.Status)
		assert.Equal(t, []watchCall{{sw.Id, testNodeIdNormalized, testAppNameForUpdate,
			testTagVersion}}, *watches)

		sRepo.AssertExpectations(t)
		msgBus.AssertExpectations(t)
	})
}

func TestUpdateSoftware_OfflineNode(t *testing.T) {
	ctx := context.Background()

	offlineNode := &creg.NodeInfo{
		Id:     testNodeIdNormalized,
		Status: creg.NodeStatusInfo{Connectivity: "Offline"},
	}

	t.Run("rejected_with_not_found", func(t *testing.T) {
		nodeClient := &mbmocks.NodeClient{}
		nodeClient.On("Get", testNodeIdNormalized).Return(offlineNode, nil).Once()

		opMgr := &mbmocks.ManagerClient{}

		s := NewSoftwareServer(testOrgName, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t),
			mocks.NewNodeRepo(t), nil, nil, nil, mbmocks.NewMsgBusServiceClient(t), false,
			[]string{testNodeGwIP}, opMgr, nil, 0, 0, nodeClient)

		resp, err := s.UpdateSoftware(ctx, &pb.UpdateSoftwareRequest{
			NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion})

		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.NotFound, status.Code(err))
		nodeClient.AssertExpectations(t)
	})

	t.Run("takes_no_lock", func(t *testing.T) {
		nodeClient := &mbmocks.NodeClient{}
		nodeClient.On("Get", testNodeIdNormalized).Return(offlineNode, nil).Once()

		opMgr := &mbmocks.ManagerClient{}

		s := NewSoftwareServer(testOrgName, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t),
			mocks.NewNodeRepo(t), nil, nil, nil, mbmocks.NewMsgBusServiceClient(t), false,
			[]string{testNodeGwIP}, opMgr, nil, 0, 0, nodeClient)

		_, err := s.UpdateSoftware(ctx, &pb.UpdateSoftwareRequest{
			NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion})

		require.Error(t, err)
		opMgr.AssertNotCalled(t, "Start", mock.Anything)
	})

	t.Run("registry_unreachable_is_internal", func(t *testing.T) {
		nodeClient := &mbmocks.NodeClient{}
		nodeClient.On("Get", testNodeIdNormalized).Return(nil, assert.AnError).Once()

		opMgr := &mbmocks.ManagerClient{}

		s := NewSoftwareServer(testOrgName, mocks.NewSoftwareRepo(t), mocks.NewAppRepo(t),
			mocks.NewNodeRepo(t), nil, nil, nil, mbmocks.NewMsgBusServiceClient(t), false,
			[]string{testNodeGwIP}, opMgr, nil, 0, 0, nodeClient)

		_, err := s.UpdateSoftware(ctx, &pb.UpdateSoftwareRequest{
			NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion})

		require.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		opMgr.AssertNotCalled(t, "Start", mock.Anything)
	})
}

func TestSoftwareServer_failOperation(t *testing.T) {
	op := &copr.OperationInfo{Id: "op-update", FencingToken: 7}

	t.Run("FailsWithFencingToken", func(t *testing.T) {
		opMgr := &mbmocks.ManagerClient{}
		opMgr.On("Fail", "op-update", uint64(7), pkg.ServiceName, "publish failed").
			Return(&copr.OperationInfo{}, nil).Once()

		s := &SoftwareServer{opManager: opMgr}
		s.failOperation(op, "UpdateSoftware", "publish failed")

		opMgr.AssertExpectations(t)
	})

	t.Run("CancelsIntentAfterFail", func(t *testing.T) {
		opMgr := &mbmocks.ManagerClient{}
		opMon := &mocks.OperationMonitor{}
		opMgr.On("Fail", "op-update", uint64(7), pkg.ServiceName, "publish failed").
			Return(&copr.OperationInfo{}, nil).Once()
		opMon.On("Cancel", &opmonpb.CancelIntentRequest{OperationId: "op-update"}).
			Return(nil, assert.AnError).Once()

		s := &SoftwareServer{opManager: opMgr, opMonitor: opMon}
		s.failOperation(op, "UpdateSoftware", "publish failed")

		opMon.AssertExpectations(t)
	})

	t.Run("FailFailureKeepsIntent", func(t *testing.T) {
		opMgr := &mbmocks.ManagerClient{}
		opMon := &mocks.OperationMonitor{}
		opMgr.On("Fail", "op-update", uint64(7), pkg.ServiceName, "publish failed").
			Return(nil, assert.AnError).Once()

		s := &SoftwareServer{opManager: opMgr, opMonitor: opMon}
		s.failOperation(op, "UpdateSoftware", "publish failed")

		opMon.AssertNotCalled(t, "Cancel", mock.Anything)
	})

	t.Run("ManagerErrorIsNotFatal", func(t *testing.T) {
		opMgr := &mbmocks.ManagerClient{}
		opMgr.On("Fail", "op-update", uint64(7), pkg.ServiceName, "publish failed").
			Return(nil, assert.AnError).Once()

		s := &SoftwareServer{opManager: opMgr}
		s.failOperation(op, "UpdateSoftware", "publish failed")

		opMgr.AssertExpectations(t)
	})

	t.Run("NoOperationIsNoop", func(t *testing.T) {
		opMgr := &mbmocks.ManagerClient{}

		s := &SoftwareServer{opManager: opMgr}
		s.failOperation(nil, "UpdateSoftware", "publish failed")

		opMgr.AssertNotCalled(t, "Fail", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything)
	})
}

func TestUpdateSoftware_RegistersTargetVersionRule(t *testing.T) {
	sw := dbSoftwareFixture()
	sRepo := mocks.NewSoftwareRepo(t)
	sRepo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).
		Return([]*db.Software{sw}, nil)

	nodeClient := &mbmocks.NodeClient{}
	nodeClient.On("Get", mock.Anything).
		Return(&creg.NodeInfo{Status: creg.NodeStatusInfo{Connectivity: "Online"}}, nil)

	op := &copr.OperationInfo{Id: uuid.NewV4().String(), FencingToken: 7}
	opMgr := &mbmocks.ManagerClient{}
	opMgr.On("Start", mock.Anything).Return(&copr.StartResponse{
		Operations: []*copr.OperationInfo{op},
	}, nil).Once()
	opMgr.On("Fail", op.Id, uint64(7), pkg.ServiceName,
		mock.MatchedBy(func(reason string) bool {
			return strings.HasPrefix(reason, "register intent failed")
		})).Return(op, nil).Once()

	opMon := mocks.NewOperationMonitor(t)
	opMon.On("Register", mock.MatchedBy(func(r *opmonpb.RegisterIntentRequest) bool {
		return r.ActionType == "UpdateSoftware" &&
			r.ResourceKey == "node:"+testNodeIdNormalized &&
			r.CompletionRule == "app."+testAppNameForUpdate+"="+testTagVersion
	})).Return(nil, assert.AnError).Once()

	s := NewSoftwareServer(testOrgName, sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
		nil, nil, nil, mbmocks.NewMsgBusServiceClient(t), false, []string{testNodeGwIP},
		opMgr, opMon, 0, 0, nodeClient)

	_, err := s.UpdateSoftware(context.Background(), &pb.UpdateSoftwareRequest{
		NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion})

	require.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
	opMon.AssertNotCalled(t, "Cancel", mock.Anything)
	opMon.AssertExpectations(t)
	opMgr.AssertExpectations(t)
}

const (
	updateSiteID = "site-1"
	siblingTower = "uk-sa2156-tnode-a1-xxxx"
)

type updateFixture struct {
	nodes *mbmocks.NodeClient
	mgr   *mbmocks.ManagerClient
	mon   *mocks.OperationMonitor
	bus   *mbmocks.MsgBusServiceClient
	repo  *mocks.SoftwareRepo
}

func newUpdateFixture(t *testing.T) *updateFixture {
	f := &updateFixture{
		nodes: &mbmocks.NodeClient{},
		mgr:   &mbmocks.ManagerClient{},
		mon:   &mocks.OperationMonitor{},
		bus:   &mbmocks.MsgBusServiceClient{},
		repo:  mocks.NewSoftwareRepo(t),
	}
	f.repo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).
		Return([]*db.Software{dbSoftwareFixture()}, nil).Once()
	f.nodes.On("Get", testNodeIdNormalized).Return(&creg.NodeInfo{
		Status: creg.NodeStatusInfo{Connectivity: "Online"},
		Site:   creg.NodeSiteInfo{SiteId: updateSiteID},
	}, nil).Once()
	return f
}

func (f *updateFixture) siteNodes(err error) {
	if err != nil {
		f.nodes.On("GetNodesBySite", updateSiteID).Return(nil, err).Once()
		return
	}
	f.nodes.On("GetNodesBySite", updateSiteID).Return(&creg.NodesBySite{
		Nodes: []creg.NodeInfo{{Id: testNodeIdNormalized}, {Id: strings.ToUpper(siblingTower)}},
	}, nil).Once()
}

func (f *updateFixture) update(t *testing.T) error {
	_, err := f.updateTag(t, testTagVersion)
	return err
}

func (f *updateFixture) updateTag(t *testing.T, tag string) (*pb.UpdateSoftwareResponse,
	error) {
	s := NewSoftwareServer(testOrgName, f.repo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
		nil, nil, nil, f.bus, false, []string{testNodeGwIP}, f.mgr, f.mon, 360, 300,
		f.nodes)
	s.startWatch = func(uuid.UUID, string, string, string) {}
	return s.UpdateSoftware(context.Background(), &pb.UpdateSoftwareRequest{
		NodeId: testNodeId, Name: testAppNameForUpdate, Tag: tag})
}

func TestUpdateSoftware_SiteLock(t *testing.T) {
	t.Run("SiteNodesAreSentAsMustBeFree", func(t *testing.T) {
		f := newUpdateFixture(t)
		f.siteNodes(nil)
		f.mgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
			return reflect.DeepEqual(req.ResourceKeys,
				[]string{"node:" + testNodeIdNormalized}) &&
				reflect.DeepEqual(req.MustBeFree, []string{
					"node:" + testNodeIdNormalized, "node:" + siblingTower,
				}) && req.LeaseSeconds == 360
		})).Return(nil, &copr.BusyError{
			Reason: "node:" + siblingTower +
				" is busy with RestartNode operation op-restart",
		}).Once()

		err := f.update(t)

		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		assert.Contains(t, status.Convert(err).Message(), "op-restart")
		f.mgr.AssertExpectations(t)
		f.mgr.AssertNotCalled(t, "Fail", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything)
		f.mon.AssertNotCalled(t, "Register", mock.Anything)
		f.bus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
	})

	t.Run("EmptyStartResponseIsUnavailable", func(t *testing.T) {
		f := newUpdateFixture(t)
		f.siteNodes(nil)
		f.mgr.On("Start", mock.Anything).Return(&copr.StartResponse{}, nil).Once()

		err := f.update(t)

		assert.Equal(t, codes.Unavailable, status.Code(err))
		f.mon.AssertNotCalled(t, "Register", mock.Anything)
	})

	t.Run("ManagerFailureIsUnavailable", func(t *testing.T) {
		f := newUpdateFixture(t)
		f.siteNodes(nil)
		f.mgr.On("Start", mock.Anything).Return(nil, &client.ErrorStatus{
			StatusCode: http.StatusServiceUnavailable, RawError: errors.New("down"),
		}).Once()

		err := f.update(t)

		assert.Equal(t, codes.Unavailable, status.Code(err))
	})

	t.Run("MarkRunningRejectedIsAborted", func(t *testing.T) {
		f := newUpdateFixture(t)
		f.siteNodes(nil)
		op := &copr.OperationInfo{Id: "op-update", FencingToken: 9,
			ResourceKey: "node:" + testNodeIdNormalized}
		f.mgr.On("Start", mock.Anything).
			Return(&copr.StartResponse{Operations: []*copr.OperationInfo{op}}, nil).
			Once()
		f.mon.On("Register", mock.Anything).
			Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
		f.mgr.On("MarkRunning", "op-update", uint64(9), pkg.ServiceName).
			Return(nil, fmt.Errorf("MarkRunning failure: %w: lease expired",
				copr.ErrNotRunnable)).Once()
		f.mgr.On("Fail", "op-update", uint64(9), pkg.ServiceName, mock.Anything).
			Return(op, nil).Once()
		f.mon.On("Cancel", mock.Anything).
			Return(&opmonpb.CancelIntentResponse{}, nil).Once()

		err := f.update(t)

		assert.Equal(t, codes.Aborted, status.Code(err))
		f.mgr.AssertExpectations(t)
		f.bus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
	})

	t.Run("SuccessReturnsOperation", func(t *testing.T) {
		f := newUpdateFixture(t)
		f.siteNodes(nil)
		op := &copr.OperationInfo{Id: "op-update", FencingToken: 9,
			ResourceKey: "node:" + testNodeIdNormalized}
		f.mgr.On("Start", mock.Anything).
			Return(&copr.StartResponse{Operations: []*copr.OperationInfo{op}}, nil).
			Once()
		f.mon.On("Register", mock.Anything).
			Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
		f.mgr.On("MarkRunning", "op-update", uint64(9), pkg.ServiceName).
			Return(op, nil).Once()
		f.bus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		f.repo.On("Update", mock.Anything).Return(errors.New("db down")).Once()

		resp, err := f.updateTag(t, testTagVersion)

		require.NoError(t, err)
		assert.Equal(t, "op-update", resp.OperationId)
		assert.Equal(t, "node:"+testNodeIdNormalized, resp.ResourceKey)
		assert.Equal(t, copr.StatusRunning, resp.Status)
		f.mgr.AssertNotCalled(t, "Fail", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything)
	})

	t.Run("MismatchedTagTakesNoLock", func(t *testing.T) {
		f := newUpdateFixture(t)

		_, err := f.updateTag(t, "9.9.9")

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		f.mgr.AssertNotCalled(t, "Start", mock.Anything)
		f.bus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
	})

	t.Run("SiteReadFailureIsUnavailable", func(t *testing.T) {
		f := newUpdateFixture(t)
		f.siteNodes(errors.New("registry down"))

		err := f.update(t)

		assert.Equal(t, codes.Unavailable, status.Code(err))
		f.mgr.AssertNotCalled(t, "Start", mock.Anything)
	})
}

func TestMarkRunningCode(t *testing.T) {
	for name, tc := range map[string]struct {
		err  error
		code codes.Code
	}{
		"NotRunnableIsAborted": {
			fmt.Errorf("MarkRunning failure: %w", copr.ErrNotRunnable), codes.Aborted,
		},
		"UnavailableIsUnavailable": {
			fmt.Errorf("MarkRunning failure: %w", copr.ErrUnavailable),
			codes.Unavailable,
		},
		"OtherIsInternal": {errors.New("MarkRunning deserialize"), codes.Internal},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.code, markRunningCode(tc.err))
		})
	}
}

type softwareStatusHealth struct {
	healthpb.HealthServiceClient
	app *healthpb.App
}

func (h *softwareStatusHealth) ListApps(context.Context, *healthpb.ListAppsRequest,
	...grpc.CallOption) (*healthpb.ListAppsResponse, error) {
	return &healthpb.ListAppsResponse{Apps: []*healthpb.App{h.app}}, nil
}

type softwareStatusHealthProvider struct {
	client healthpb.HealthServiceClient
}

func (h softwareStatusHealthProvider) GetClient() (healthpb.HealthServiceClient, error) {
	return h.client, nil
}

func TestSoftwarePromotionDuringUpdate(t *testing.T) {
	const target = "1.2.3-accdefgh"
	const next = "1.2.4-accdefgh"
	sw := dbSoftwareFixture()
	sw.NodeId = testNodeIdNormalized
	sw.DesiredVersion = target
	sw.Status = ukama.UpdateInProgress

	sRepo := mocks.NewSoftwareRepo(t)
	// SetDesiredForApp skips UpdateInProgress rows in SQL, so this row is not touched.
	sRepo.On("SetDesiredForApp", sw.AppName, next).Return(int64(0), nil).Once()
	sRepo.On("Get", sw.Id).Return(*sw, nil).Once()
	var saved db.Software
	sRepo.On("Update", mock.Anything).Run(func(args mock.Arguments) {
		saved = *args.Get(0).(*db.Software)
	}).Return(nil).Once()

	releases := mocks.NewReleaseRepo(t)
	releases.On("Upsert", mock.Anything).Return(nil).Once()
	releases.On("SetDesired", mock.MatchedBy(func(d *db.AppDesiredRelease) bool {
		return d.Name == sw.AppName && d.DesiredVersion == next
	})).Return(nil).Once()
	releases.On("GetDesired", sw.AppName, "app").
		Return(&db.AppDesiredRelease{DesiredVersion: next}, nil).Once()

	s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
		mbmocks.NewMsgBusServiceClient(t))
	s.releaseRepo = releases
	s.healthClient = softwareStatusHealthProvider{client: &softwareStatusHealth{
		app: &healthpb.App{Name: sw.AppName, Version: target, Status: "Active"},
	}}
	resp, err := s.PromoteRelease(context.Background(),
		&pb.PromoteReleaseRequest{Name: sw.AppName, Version: next})
	require.NoError(t, err)
	assert.Equal(t, next, resp.DesiredVersion)
	assert.Equal(t, target, sw.DesiredVersion)
	assert.Equal(t, ukama.UpdateInProgress, sw.Status)
	sRepo.AssertNotCalled(t, "Update", mock.Anything)

	// A competing request cannot select the in-progress row for dispatch.
	sRepo.On("List", sw.NodeId, ukama.UpdateAvailable, sw.AppName).
		Return([]*db.Software{}, nil).Once()
	_, err = s.UpdateSoftware(context.Background(),
		&pb.UpdateSoftwareRequest{NodeId: sw.NodeId, Name: sw.AppName, Tag: next})
	require.Error(t, err)

	// Use the persisted target, as ResumeInProgressUpdates does after a restart.
	s.watchSoftwareUpdate(sw.Id, sw.NodeId, sw.AppName, sw.DesiredVersion,
		time.Now().Add(time.Second), time.Millisecond)
	assert.Equal(t, target, saved.CurrentVersion)
	assert.Equal(t, next, saved.DesiredVersion)
	assert.Equal(t, ukama.UpdateAvailable, saved.Status)
	assert.Contains(t, saved.ChangeLogs, "Software successfully updated to version "+target)
}

func TestPersistSoftwareStatusUsesConfirmedVersion(t *testing.T) {
	const target = "1.2.3-accdefgh"
	const next = "1.2.4-accdefgh"
	for _, tc := range []struct {
		name          string
		storedDesired string
		promoted      string
		confirmed     string
		outcome       ukama.SoftwareStatusType
		wantCurrent   string
		wantDesired   string
		wantStatus    ukama.SoftwareStatusType
	}{
		{"target remains desired", target, target, target, ukama.UpToDate, target, target,
			ukama.UpToDate},
		{"new promotion", target, next, target, ukama.UpToDate, target, next,
			ukama.UpdateAvailable},
		{"already changed desired row", next, next, target, ukama.UpToDate, target, next,
			ukama.UpdateAvailable},
		{"no catalog desired", target, "", target, ukama.UpToDate, target, target,
			ukama.UpToDate},
		{"timeout preserves installed version", target, next, "", ukama.UpdateFailed,
			testCurrentVersion, next, ukama.UpdateFailed},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sw := dbSoftwareFixture()
			sw.DesiredVersion = tc.storedDesired
			sw.Status = ukama.UpdateInProgress
			sRepo := mocks.NewSoftwareRepo(t)
			sRepo.On("Get", sw.Id).Return(*sw, nil).Once()
			var saved db.Software
			sRepo.On("Update", mock.Anything).Run(func(args mock.Arguments) {
				saved = *args.Get(0).(*db.Software)
			}).Return(nil).Once()
			releases := mocks.NewReleaseRepo(t)
			var desired *db.AppDesiredRelease
			if tc.promoted != "" {
				desired = &db.AppDesiredRelease{DesiredVersion: tc.promoted}
			}
			releases.On("GetDesired", sw.AppName, "app").Return(desired, nil).Once()
			s := &SoftwareServer{sRepo: sRepo, releaseRepo: releases}

			require.True(t, s.persistSoftwareStatus(sw.Id, sw.NodeId, sw.AppName, tc.confirmed,
				tc.outcome, "update result"))
			assert.Equal(t, tc.wantCurrent, saved.CurrentVersion)
			assert.Equal(t, tc.wantDesired, saved.DesiredVersion)
			assert.Equal(t, tc.wantStatus, saved.Status)
		})
	}
}

func TestWatchSoftwareUpdateRetriesPersistence(t *testing.T) {
	const target = "1.2.3-accdefgh"
	const next = "1.2.4-accdefgh"
	for _, failure := range []string{"load", "catalog", "save"} {
		t.Run(failure, func(t *testing.T) {
			sw := dbSoftwareFixture()
			sw.DesiredVersion = target
			sw.Status = ukama.UpdateInProgress
			sRepo := mocks.NewSoftwareRepo(t)
			releases := mocks.NewReleaseRepo(t)
			transient := errors.New("temporary database failure")
			if failure == "load" {
				sRepo.On("Get", sw.Id).Return(db.Software{}, transient).Once()
				sRepo.On("Get", sw.Id).Return(*sw, nil).Once()
			} else {
				sRepo.On("Get", sw.Id).Return(*sw, nil).Twice()
			}
			desired := &db.AppDesiredRelease{DesiredVersion: next}
			if failure == "catalog" {
				releases.On("GetDesired", sw.AppName, "app").Return(nil, transient).Once()
			}
			calls := 1
			if failure == "save" {
				calls = 2
				sRepo.On("Update", mock.Anything).Return(transient).Once()
			}
			releases.On("GetDesired", sw.AppName, "app").Return(desired, nil).Times(calls)
			var saved db.Software
			sRepo.On("Update", mock.Anything).Run(func(args mock.Arguments) {
				saved = *args.Get(0).(*db.Software)
			}).Return(nil).Once()
			s := &SoftwareServer{sRepo: sRepo, releaseRepo: releases,
				healthClient: softwareStatusHealthProvider{client: &softwareStatusHealth{
					app: &healthpb.App{Name: sw.AppName, Version: target, Status: "Active"},
				}}}

			// A failed save after confirmation must not turn success into a timeout.
			s.watchSoftwareUpdate(sw.Id, sw.NodeId, sw.AppName, target,
				time.Now().Add(-time.Second), time.Millisecond)
			assert.Equal(t, target, saved.CurrentVersion)
			assert.Equal(t, next, saved.DesiredVersion)
			assert.Equal(t, ukama.UpdateAvailable, saved.Status)
			sRepo.AssertNumberOfCalls(t, "Get", 2)
			assert.Equal(t, []string{"initial",
				"Software successfully updated to version " + target}, saved.ChangeLogs)
		})
	}
}

func TestPromotionStillUpdatesIdleSoftware(t *testing.T) {
	sw := dbSoftwareFixture()
	sRepo := mocks.NewSoftwareRepo(t)
	sRepo.On("SetDesiredForApp", sw.AppName, "1.2.4-accdefgh").Return(int64(1), nil).Once()
	s := &SoftwareServer{sRepo: sRepo}
	s.recomputeDesiredForApp(sw.AppName, "1.2.4-accdefgh")
	sRepo.AssertNotCalled(t, "List", mock.Anything, mock.Anything, mock.Anything)
	sRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestPersistSoftwareStatusStopsWhenRecordDeleted(t *testing.T) {
	sw := dbSoftwareFixture()
	sRepo := mocks.NewSoftwareRepo(t)
	sRepo.On("Get", sw.Id).Return(db.Software{}, gorm.ErrRecordNotFound).Once()
	s := &SoftwareServer{sRepo: sRepo}

	assert.True(t, s.persistSoftwareStatus(sw.Id, sw.NodeId, sw.AppName, "",
		ukama.UpdateFailed, "update result"))
	sRepo.AssertNotCalled(t, "Update", mock.Anything)
}
