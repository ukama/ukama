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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	ukamapb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/health/mocks"
	pb "github.com/ukama/ukama/systems/node/health/pb/gen"
	"github.com/ukama/ukama/systems/node/health/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const testOrgName = "test-org"

var testNode = ukama.NewVirtualNodeId("HomeNode")

func TestHealthServerStoreHealthReport(t *testing.T) {
	hRepo := &mocks.HealthRepo{}

	reported := time.Date(2026, 4, 21, 10, 0, 0, 0, time.UTC)
	payload, _ := json.Marshal(map[string]string{"k": "v"})
	req := &pb.StoreHealthReportRequest{
		NodeId:        testNode.String(),
		NodeType:      string(ukama.NODE_TYPE_HOMENODE),
		SchemaVersion: "1",
		ReportedAt:    timestamppb.New(reported),
		Payload:       payload,
	}

	hRepo.On("StoreHealthReport", mock.MatchedBy(func(r *db.HealthReport) bool {
		return r != nil &&
			r.NodeID == testNode.StringLowercase() &&
			r.NodeType == ukama.NodeType(req.GetNodeType()) &&
			r.SchemaVersion == "1" &&
			r.ReportedAt.Equal(reported) &&
			string(r.Payload) == string(payload)
	}), mock.MatchedBy(func(ts time.Time) bool {
		return !ts.IsZero()
	})).Return(nil).Once()

	s := NewHealthServer(testOrgName, hRepo, false)
	resp, err := s.StoreHealthReport(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.GetReportId())
	hRepo.AssertExpectations(t)
}

func TestHealthServerStoreHealthReportInvalidNodeId(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	s := NewHealthServer(testOrgName, hRepo, false)

	resp, err := s.StoreHealthReport(context.Background(), &pb.StoreHealthReportRequest{
		NodeId:       "invalid-node",
		NodeType:   string(ukama.NODE_TYPE_HOMENODE),
		ReportedAt: timestamppb.New(time.Now()),
	})

	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	}
	hRepo.AssertExpectations(t)
}

func TestHealthServerStoreHealthReportMissingReportedAt(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	s := NewHealthServer(testOrgName, hRepo, false)

	resp, err := s.StoreHealthReport(context.Background(), &pb.StoreHealthReportRequest{
		NodeId:     testNode.String(),
		NodeType: string(ukama.NODE_TYPE_HOMENODE),
	})

	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	}
	hRepo.AssertExpectations(t)
}

func TestHealthServerStoreHealthReportRepoError(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	req := &pb.StoreHealthReportRequest{
		NodeId:       testNode.String(),
		NodeType:   string(ukama.NODE_TYPE_HOMENODE),
		ReportedAt: timestamppb.New(time.Now().UTC()),
	}

	hRepo.On("StoreHealthReport", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	s := NewHealthServer(testOrgName, hRepo, false)
	resp, err := s.StoreHealthReport(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	hRepo.AssertExpectations(t)
}

func TestHealthServerListLatest(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	reportID := uuid.NewV4()
	reported := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	received := reported.Add(time.Minute)
	raw := json.RawMessage(`{}`)

	report := &db.HealthReport{
		ID:            reportID,
		NodeID:        testNode.String(),
		NodeType:      ukama.NODE_TYPE_HOMENODE,
		SchemaVersion: "1",
		ReportedAt:    reported,
		ReceivedAt:    received,
		Payload:       raw,
	}

	hRepo.On("List", "", testNode.String(), (*time.Time)(nil), ukama.FilterTimeframesTypeLatest).Return([]*db.HealthReport{report}, nil).Once()
	s := NewHealthServer(testOrgName, hRepo, false)

	resp, err := s.List(context.Background(), &pb.ListRequest{
		NodeId:    testNode.String(),
		Timeframe: ukamapb.FilterTimeframesType_LATEST,
	})

	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.Reports, 1) {
		r := resp.Reports[0]
		assert.Equal(t, reportID.String(), r.Id)
		assert.Equal(t, testNode.String(), r.NodeId)
		assert.Equal(t, string(ukama.NODE_TYPE_HOMENODE), r.NodeType)
	}
	hRepo.AssertExpectations(t)
}

func TestHealthServerListAll(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	id1 := uuid.NewV4()
	id2 := uuid.NewV4()
	ts1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	ts2 := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	reports := []*db.HealthReport{
		{ID: id1, NodeID: testNode.String(), ReportedAt: ts1, SchemaVersion: "1", Payload: json.RawMessage(`{}`)},
		{ID: id2, NodeID: testNode.String(), ReportedAt: ts2, SchemaVersion: "1", Payload: json.RawMessage(`{}`)},
	}

	hRepo.On("List", "", testNode.String(), (*time.Time)(nil), ukama.FilterTimeframesTypeAll).Return(reports, nil).Once()
	s := NewHealthServer(testOrgName, hRepo, false)

	resp, err := s.List(context.Background(), &pb.ListRequest{
		NodeId:    testNode.String(),
		Timeframe: ukamapb.FilterTimeframesType_ALL,
	})

	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.Reports, 2) {
		assert.Equal(t, ts1.Unix(), resp.Reports[0].ReportedAt.AsTime().Unix())
		assert.Equal(t, ts2.Unix(), resp.Reports[1].ReportedAt.AsTime().Unix())
	}
	hRepo.AssertExpectations(t)
}

func TestHealthServerListMissingIdAndNodeId(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	s := NewHealthServer(testOrgName, hRepo, false)

	resp, err := s.List(context.Background(), &pb.ListRequest{
		Timeframe: ukamapb.FilterTimeframesType_LATEST,
	})

	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	}
	hRepo.AssertExpectations(t)
}

func TestHealthServerListRepoError(t *testing.T) {
	hRepo := &mocks.HealthRepo{}

	hRepo.On("List", "", testNode.String(), (*time.Time)(nil), ukama.FilterTimeframesTypeLatest).Return(([]*db.HealthReport)(nil), assert.AnError).Once()
	s := NewHealthServer(testOrgName, hRepo, false)

	resp, err := s.List(context.Background(), &pb.ListRequest{
		NodeId:    testNode.String(),
		Timeframe: ukamapb.FilterTimeframesType_LATEST,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	hRepo.AssertExpectations(t)
}
