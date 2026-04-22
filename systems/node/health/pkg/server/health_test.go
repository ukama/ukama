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
	ukamapb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/health/mocks"
	pb "github.com/ukama/ukama/systems/node/health/pb/gen"
	"github.com/ukama/ukama/systems/node/health/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const testOrgName = "test-org"
const testTimestamp = "2026-04-21T10:00:00Z"

var testNode = ukama.NewVirtualNodeId("HomeNode")

func TestHealthServerStoreRunningAppsInfo(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	hRepo := &mocks.HealthRepo{}

	req := &pb.StoreRunningAppsInfoRequest{
		NodeId:    testNode.String(),
		Timestamp: testTimestamp,
		System: []*pb.System{
			{Name: "cpu", Value: "35"},
		},
		Capps: []*pb.Capps{
			{
				Space:  "core",
				Name:   "agent",
				Tag:    "v1",
				Status: pb.Status_ACTIVE,
				Resources: []*pb.Resource{
					{Name: "mem", Value: "100mb"},
				},
			},
		},
	}

	hRepo.On("StoreRunningAppsInfo", mock.MatchedBy(func(health *db.Health) bool {
		return health != nil &&
			health.NodeId == testNode.StringLowercase() &&
			health.TimeStamp == req.Timestamp &&
			len(health.System) == 1 &&
			health.System[0].Name == "cpu" &&
			health.System[0].Value == "35" &&
			len(health.Capps) == 1 &&
			health.Capps[0].Name == "agent" &&
			health.Capps[0].Status == db.Status(pb.Status_ACTIVE) &&
			len(health.Capps[0].Resources) == 1 &&
			health.Capps[0].Resources[0].Name == "mem" &&
			health.Capps[0].Resources[0].Value == "100mb"
	}), mock.Anything).Return(nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)
	resp, err := s.StoreRunningAppsInfo(context.TODO(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	hRepo.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}

func TestHealthServerStoreRunningAppsInfoInvalidNodeId(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	hRepo := &mocks.HealthRepo{}
	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)

	resp, err := s.StoreRunningAppsInfo(context.TODO(), &pb.StoreRunningAppsInfoRequest{
		NodeId: "invalid-node",
	})

	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	}
	hRepo.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}

func TestHealthServerStoreRunningAppsInfoRepoError(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	hRepo := &mocks.HealthRepo{}
	req := &pb.StoreRunningAppsInfoRequest{
		NodeId:    testNode.String(),
		Timestamp: testTimestamp,
	}

	hRepo.On("StoreRunningAppsInfo", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)
	resp, err := s.StoreRunningAppsInfo(context.TODO(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	hRepo.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}

func TestHealthServerStoreRunningAppsInfoPublishError(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	hRepo := &mocks.HealthRepo{}
	req := &pb.StoreRunningAppsInfoRequest{
		NodeId:    testNode.String(),
		Timestamp: testTimestamp,
	}

	hRepo.On("StoreRunningAppsInfo", mock.Anything, mock.Anything).Return(nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)
	resp, err := s.StoreRunningAppsInfo(context.TODO(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	hRepo.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}

func TestHealthServerListLatest(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	hRepo := &mocks.HealthRepo{}
	id := uuid.NewV4()
	cappID := uuid.NewV4()

	health := db.Health{
		Id:        id,
		NodeId:    testNode.String(),
		TimeStamp: "test",
		System: []db.System{
			{
				Id:       uuid.NewV4(),
				HealthID: id,
				Name:     "cpu",
				Value:    "40",
			},
		},
		Capps: []db.Capp{
			{
				Id:       cappID,
				HealthID: id,
				Space:    "core",
				Name:     "agent",
				Tag:      "v1",
				Status:   db.Status(1),
				Resources: []db.Resource{
					{
						Id:     uuid.NewV4(),
						CappID: cappID,
						Name:   "memory",
						Value:  "100mb",
					},
				},
			},
		},
	}

	hRepo.On("List", "", testNode.String(), "", ukama.FilterTimeframesTypeLatest).Return([]*db.Health{&health}, nil).Once()
	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)

	resp, err := s.List(context.TODO(), &pb.ListRequest{
		NodeId: testNode.String(),
		Filter: ukamapb.FilterTimeframesType_LATEST,
	})

	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.Healths, 1) {
		assert.Equal(t, health.NodeId, resp.Healths[0].NodeId)
		assert.Equal(t, health.TimeStamp, resp.Healths[0].Timestamp)
		if assert.Len(t, resp.Healths[0].System, 1) {
			assert.Equal(t, health.System[0].Name, resp.Healths[0].System[0].Name)
			assert.Equal(t, health.System[0].Value, resp.Healths[0].System[0].Value)
		}
		if assert.Len(t, resp.Healths[0].Capps, 1) {
			assert.Equal(t, health.Capps[0].Name, resp.Healths[0].Capps[0].Name)
			if assert.Len(t, resp.Healths[0].Capps[0].Resources, 1) {
				assert.Equal(t, health.Capps[0].Resources[0].Name, resp.Healths[0].Capps[0].Resources[0].Name)
				assert.Equal(t, health.Capps[0].Resources[0].Value, resp.Healths[0].Capps[0].Resources[0].Value)
			}
		}
	}

	hRepo.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}

func TestHealthServerListAll(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	hRepo := &mocks.HealthRepo{}
	id1 := uuid.NewV4()
	id2 := uuid.NewV4()
	healths := []*db.Health{
		{Id: id1, NodeId: testNode.String(), TimeStamp: "ts-1"},
		{Id: id2, NodeId: testNode.String(), TimeStamp: "ts-2"},
	}

	hRepo.On("List", "", testNode.String(), "", ukama.FilterTimeframesTypeAll).Return(healths, nil).Once()
	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)

	resp, err := s.List(context.TODO(), &pb.ListRequest{
		NodeId: testNode.String(),
		Filter: ukamapb.FilterTimeframesType_ALL,
	})

	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.Healths, 2) {
		assert.Equal(t, "ts-1", resp.Healths[0].Timestamp)
		assert.Equal(t, "ts-2", resp.Healths[1].Timestamp)
	}

	hRepo.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}

func TestHealthServerListMissingIdAndNodeId(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	hRepo := &mocks.HealthRepo{}
	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)

	resp, err := s.List(context.TODO(), &pb.ListRequest{
		Filter: ukamapb.FilterTimeframesType_LATEST,
	})

	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	}

	hRepo.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}

func TestHealthServerListLatestRepoError(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	hRepo := &mocks.HealthRepo{}

	hRepo.On("List", "", testNode.String(), "", ukama.FilterTimeframesTypeLatest).Return(([]*db.Health)(nil), assert.AnError).Once()
	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)

	resp, err := s.List(context.TODO(), &pb.ListRequest{
		NodeId: testNode.String(),
		Filter: ukamapb.FilterTimeframesType_LATEST,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)

	hRepo.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}

func TestHealthServerListAllRepoError(t *testing.T) {
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	hRepo := &mocks.HealthRepo{}

	hRepo.On("List", "", testNode.String(), "", ukama.FilterTimeframesTypeAll).Return(([]*db.Health)(nil), assert.AnError).Once()
	s := NewHealthServer(testOrgName, hRepo, msgclientRepo, false)

	resp, err := s.List(context.TODO(), &pb.ListRequest{
		NodeId: testNode.String(),
		Filter: ukamapb.FilterTimeframesType_ALL,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)

	hRepo.AssertExpectations(t)
	msgclientRepo.AssertExpectations(t)
}

