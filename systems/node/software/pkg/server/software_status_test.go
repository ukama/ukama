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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/ukama"
	healthpb "github.com/ukama/ukama/systems/node/health/pb/gen"
	"github.com/ukama/ukama/systems/node/software/mocks"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
	"google.golang.org/grpc"
)

type softwareStatusHealth struct {
	healthpb.HealthServiceClient
	app *healthpb.App
}

func (h *softwareStatusHealth) ListApps(context.Context, *healthpb.ListAppsRequest, ...grpc.CallOption) (*healthpb.ListAppsResponse, error) {
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
	sRepo.On("List", "", ukama.Unknown, sw.AppName).Return([]*db.Software{sw}, nil).Once()
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
	releases.On("GetDesired", sw.AppName, "app").Return(&db.AppDesiredRelease{DesiredVersion: next}, nil).Once()

	s := newTestServer(sRepo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t), mbmocks.NewMsgBusServiceClient(t))
	s.releaseRepo = releases
	s.healthClient = softwareStatusHealthProvider{client: &softwareStatusHealth{
		app: &healthpb.App{Name: sw.AppName, Version: target, Status: "Active"},
	}}
	resp, err := s.PromoteRelease(context.Background(), &pb.PromoteReleaseRequest{Name: sw.AppName, Version: next})
	require.NoError(t, err)
	assert.Equal(t, next, resp.DesiredVersion)
	assert.Equal(t, target, sw.DesiredVersion)
	assert.Equal(t, ukama.UpdateInProgress, sw.Status)
	sRepo.AssertNotCalled(t, "Update", mock.Anything)

	// A competing request cannot select the in-progress row for dispatch.
	sRepo.On("List", sw.NodeId, ukama.UpdateAvailable, sw.AppName).Return([]*db.Software{}, nil).Once()
	_, err = s.UpdateSoftware(context.Background(), &pb.UpdateSoftwareRequest{NodeId: sw.NodeId, Name: sw.AppName, Tag: next})
	require.Error(t, err)

	// Use the persisted target, as ResumeInProgressUpdates does after a restart.
	s.watchSoftwareUpdate(sw.Id, sw.NodeId, sw.AppName, sw.DesiredVersion, time.Now().Add(time.Second), time.Millisecond)
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
		{"target remains desired", target, target, target, ukama.UpToDate, target, target, ukama.UpToDate},
		{"new promotion", target, next, target, ukama.UpToDate, target, next, ukama.UpdateAvailable},
		{"already changed desired row", next, next, target, ukama.UpToDate, target, next, ukama.UpdateAvailable},
		{"no catalog desired", target, "", target, ukama.UpToDate, target, target, ukama.UpToDate},
		{"timeout preserves installed version", target, next, "", ukama.UpdateFailed, testCurrentVersion, next, ukama.UpdateFailed},
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

			require.True(t, s.persistSoftwareStatus(sw.Id, sw.NodeId, sw.AppName, tc.confirmed, tc.outcome, "update result"))
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
			s := &SoftwareServer{sRepo: sRepo, releaseRepo: releases, healthClient: softwareStatusHealthProvider{
				client: &softwareStatusHealth{app: &healthpb.App{Name: sw.AppName, Version: target, Status: "Active"}},
			}}

			// A failed save after confirmation must not turn success into a timeout.
			s.watchSoftwareUpdate(sw.Id, sw.NodeId, sw.AppName, target, time.Now().Add(-time.Second), time.Millisecond)
			assert.Equal(t, target, saved.CurrentVersion)
			assert.Equal(t, next, saved.DesiredVersion)
			assert.Equal(t, ukama.UpdateAvailable, saved.Status)
			sRepo.AssertNumberOfCalls(t, "Get", 2)
			assert.Equal(t, []string{"initial", "Software successfully updated to version " + target}, saved.ChangeLogs)
		})
	}
}

func TestPromotionStillUpdatesIdleSoftware(t *testing.T) {
	sw := dbSoftwareFixture()
	sRepo := mocks.NewSoftwareRepo(t)
	sRepo.On("List", "", ukama.Unknown, sw.AppName).Return([]*db.Software{sw}, nil).Once()
	sRepo.On("Update", mock.MatchedBy(func(row *db.Software) bool {
		return row.CurrentVersion == testCurrentVersion && row.DesiredVersion == "1.2.4-accdefgh" && row.Status == ukama.UpdateAvailable
	})).Return(nil).Once()
	s := &SoftwareServer{sRepo: sRepo}
	s.recomputeDesiredForApp(sw.AppName, "1.2.4-accdefgh")
}
