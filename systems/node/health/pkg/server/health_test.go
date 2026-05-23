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
)

const testOrgName = "test-org"

var testNode = ukama.NewVirtualNodeId("HomeNode")
var testCNode = ukama.NewVirtualNodeId("ctrlnode")

func newTestHealthServer(hRepo *mocks.HealthRepo) *HealthServer {
	return NewHealthServer(testOrgName, hRepo, false, nil)
}

func TestHealthServerStoreHealthReport(t *testing.T) {
	hRepo := &mocks.HealthRepo{}

	reported := time.Date(2026, 4, 21, 10, 0, 0, 0, time.UTC)
	reportedUnix := reported.Unix()
	payload := []byte(`{"k":"v","nodeType":"hnode","schemaVersion":"1","reportedAt":"` + jsonNumber(reportedUnix) + `"}`)
	req := &pb.StoreHealthReportRequest{
		NodeId:  testNode.String(),
		Payload: payload,
	}

	hRepo.On("StoreHealthReport", mock.MatchedBy(func(r *db.HealthReport) bool {
		return r != nil &&
			r.NodeID == testNode.StringLowercase() &&
			r.NodeType == ukama.NODE_TYPE_HOMENODE &&
			r.SchemaVersion == "1" &&
			r.ReportedAt.Equal(reported) &&
			string(r.Payload) == string(payload)
	}), mock.MatchedBy(func(ts time.Time) bool {
		return !ts.IsZero()
	})).Return(nil).Once()

	s := newTestHealthServer(hRepo)
	resp, err := s.StoreHealthReport(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.GetReportId())
	hRepo.AssertExpectations(t)
}

func TestHealthServerStoreHealthReportUnixReportedAt(t *testing.T) {
	hRepo := &mocks.HealthRepo{}

	const reportedUnix int64 = 1779534357
	reported := time.Unix(reportedUnix, 0).UTC()
	payload := []byte(`{"nodeType":"hnode","schemaVersion":"1","reportedAt":"1779534357"}`)
	req := &pb.StoreHealthReportRequest{
		NodeId:  testNode.String(),
		Payload: payload,
	}

	hRepo.On("StoreHealthReport", mock.MatchedBy(func(r *db.HealthReport) bool {
		return r != nil &&
			r.NodeID == testNode.StringLowercase() &&
			r.NodeType == ukama.NODE_TYPE_HOMENODE &&
			r.SchemaVersion == "1" &&
			r.ReportedAt.Equal(reported) &&
			string(r.Payload) == string(payload)
	}), mock.MatchedBy(func(ts time.Time) bool {
		return !ts.IsZero()
	})).Return(nil).Once()

	s := newTestHealthServer(hRepo)
	resp, err := s.StoreHealthReport(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.GetReportId())
	hRepo.AssertExpectations(t)
}

func TestHealthServerStoreHealthReportInvalidNodeId(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	s := newTestHealthServer(hRepo)

	reportedUnix := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	payload := []byte(`{"nodeType":"hnode","schemaVersion":"1","reportedAt":"` + jsonNumber(reportedUnix) + `"}`)
	resp, err := s.StoreHealthReport(context.Background(), &pb.StoreHealthReportRequest{
		NodeId:  "invalid-node",
		Payload: payload,
	})

	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	}
	hRepo.AssertExpectations(t)
}

func TestHealthServerStoreHealthReportMissingReportedAt(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	s := newTestHealthServer(hRepo)

	payload, _ := json.Marshal(map[string]string{
		"nodeType": string(ukama.NODE_TYPE_HOMENODE),
	})
	resp, err := s.StoreHealthReport(context.Background(), &pb.StoreHealthReportRequest{
		NodeId:  testNode.String(),
		Payload: payload,
	})

	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	}
	hRepo.AssertExpectations(t)
}

func TestHealthServerStoreHealthReportRepoError(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	reportedUnix := time.Now().UTC().Unix()
	payload := []byte(`{"nodeType":"hnode","schemaVersion":"1","reportedAt":"` + jsonNumber(reportedUnix) + `"}`)
	req := &pb.StoreHealthReportRequest{
		NodeId:  testNode.String(),
		Payload: payload,
	}

	hRepo.On("StoreHealthReport", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	s := newTestHealthServer(hRepo)
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
	s := newTestHealthServer(hRepo)

	resp, err := s.ListReports(context.Background(), &pb.ListReportsRequest{
		NodeId:    testNode.String(),
		Timeframe: ukamapb.FilterTimeframesType_LATEST,
	})

	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.Reports, 1) {
		r := resp.Reports[0]
		assert.Equal(t, reportID.String(), r.Id)
		assert.Equal(t, testNode.String(), r.NodeId)
		assert.Equal(t, string(ukama.NODE_TYPE_HOMENODE), r.NodeType)
		assert.Equal(t, reported.Unix(), r.ReportedAt)
		assert.True(t, r.ReceivedAt.AsTime().Equal(received))
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
		{ID: id1, NodeID: testNode.String(), NodeType: ukama.NODE_TYPE_HOMENODE, ReportedAt: ts1, SchemaVersion: "1", Payload: json.RawMessage(`{}`)},
		{ID: id2, NodeID: testNode.String(), NodeType: ukama.NODE_TYPE_HOMENODE, ReportedAt: ts2, SchemaVersion: "1", Payload: json.RawMessage(`{}`)},
	}

	hRepo.On("List", "", testNode.String(), (*time.Time)(nil), ukama.FilterTimeframesTypeAll).Return(reports, nil).Once()
	s := newTestHealthServer(hRepo)

	resp, err := s.ListReports(context.Background(), &pb.ListReportsRequest{
		NodeId:    testNode.String(),
		Timeframe: ukamapb.FilterTimeframesType_ALL,
	})

	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.Reports, 2) {
		assert.Equal(t, ts1.Unix(), resp.Reports[0].ReportedAt)
		assert.Equal(t, ts2.Unix(), resp.Reports[1].ReportedAt)
	}
	hRepo.AssertExpectations(t)
}

func TestHealthServerListMissingIdAndNodeId(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	s := newTestHealthServer(hRepo)

	resp, err := s.ListReports(context.Background(), &pb.ListReportsRequest{
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
	s := newTestHealthServer(hRepo)

	resp, err := s.ListReports(context.Background(), &pb.ListReportsRequest{
		NodeId:    testNode.String(),
		Timeframe: ukamapb.FilterTimeframesType_LATEST,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	hRepo.AssertExpectations(t)
}

func TestHealthServerListApps(t *testing.T) {
	hRepo := &mocks.HealthRepo{}
	reportID := uuid.NewV4()
	reported := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	received := reported.Add(time.Minute)
	payload, err := json.Marshal(map[string]interface{}{
		"apps": []map[string]string{
			{
				"name":    "test-app",
				"version": "1.0.0",
				"tag":     "latest",
				"state":   "active",
			},
		},
	})
	assert.NoError(t, err)
	raw := json.RawMessage(payload)

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
	s := newTestHealthServer(hRepo)

	resp, err := s.ListApps(context.Background(), &pb.ListAppsRequest{
		NodeId: testNode.String(),
	})

	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.Apps, 1) {
		assert.Equal(t, "test-app", resp.Apps[0].Name)
		assert.Equal(t, "1.0.0", resp.Apps[0].Version)
		assert.Equal(t, "latest", resp.Apps[0].Tag)
		assert.Equal(t, "active", resp.Apps[0].Status)
	}
	hRepo.AssertExpectations(t)
}

func TestHealthServerListInterfaces(t *testing.T) {
	t.Run("invalid_payload_returns_error", func(t *testing.T) {
		hRepo := &mocks.HealthRepo{}
		hRepo.On("List", "", testNode.String(), (*time.Time)(nil), ukama.FilterTimeframesTypeLatest).
			Return([]*db.HealthReport{{Payload: json.RawMessage("not-json")}}, nil).Once()

		s := newTestHealthServer(hRepo)
		resp, err := s.ListInterfaces(context.Background(), &pb.ListInterfacesRequest{
			NodeId: testNode.String(),
		})

		assert.Nil(t, resp)
		if assert.Error(t, err) {
			assert.Equal(t, codes.InvalidArgument, status.Code(err))
		}
		hRepo.AssertExpectations(t)
	})

	t.Run("no_reports_returns_empty_interfaces", func(t *testing.T) {
		hRepo := &mocks.HealthRepo{}
		hRepo.On("List", "", testCNode.String(), (*time.Time)(nil), ukama.FilterTimeframesTypeLatest).
			Return([]*db.HealthReport{}, nil).Once()

		s := newTestHealthServer(hRepo)
		resp, err := s.ListInterfaces(context.Background(), &pb.ListInterfacesRequest{
			NodeId: testCNode.String(),
		})

		assert.NoError(t, err)
		if assert.NotNil(t, resp) && assert.NotNil(t, resp.Interfaces) {
			assert.Nil(t, resp.Interfaces.Switch)
		}
		hRepo.AssertExpectations(t)
	})

	t.Run("maps_policy_and_ports_from_latest_report", func(t *testing.T) {
		hRepo := &mocks.HealthRepo{}
		reportID := uuid.NewV4()
		reported := time.Date(2026, 1, 3, 10, 0, 0, 0, time.UTC)
		received := reported.Add(time.Minute)

		payload, err := json.Marshal(map[string]interface{}{
			"interfaces": map[string]interface{}{
				"gps": map[string]interface{}{
					"available":   true,
					"lock":        true,
					"coordinates": "-114.0719,51.0447",
					"time":        "2026-05-10T18:29:55Z",
				},
				"switch": map[string]interface{}{
					"state": "active",
					"policy": map[string]interface{}{
						"state":  "applied",
						"hash":   "hash-123",
						"source": "controller",
						"error":  "",
					},
					"ports": []map[string]interface{}{
						{
							"id":         1,
							"name":       "ge1",
							"present":    true,
							"adminState": "up",
							"linkState":  "up",
							"poeState":   "on",
						},
					},
				},
			},
		})
		assert.NoError(t, err)

		report := &db.HealthReport{
			ID:            reportID,
			NodeID:        testCNode.String(),
			NodeType:      ukama.NODE_TYPE_CNODE,
			SchemaVersion: "1",
			ReportedAt:    reported,
			ReceivedAt:    received,
			Payload:       json.RawMessage(payload),
		}

		hRepo.On("List", "", testCNode.String(), (*time.Time)(nil), ukama.FilterTimeframesTypeLatest).
			Return([]*db.HealthReport{report}, nil).Once()

		s := newTestHealthServer(hRepo)
		resp, err := s.ListInterfaces(context.Background(), &pb.ListInterfacesRequest{
			NodeId: testCNode.String(),
		})

		assert.NoError(t, err)
		if assert.NotNil(t, resp) && assert.NotNil(t, resp.Interfaces) {
			if assert.NotNil(t, resp.Interfaces.Gps) {
				assert.Equal(t, "2026-05-10T18:29:55Z", resp.Interfaces.Gps.Time)
			}
			if assert.NotNil(t, resp.Interfaces.Switch) {
			assert.Equal(t, "active", resp.Interfaces.Switch.State)
			assert.Equal(t, "hash-123", resp.Interfaces.Switch.Policy.Hash)
			assert.Equal(t, "controller", resp.Interfaces.Switch.Policy.Source)
			if assert.Len(t, resp.Interfaces.Switch.Ports, 1) {
				p := resp.Interfaces.Switch.Ports[0]
				assert.Equal(t, int64(1), p.Id)
				assert.Equal(t, "ge1", p.Name)
				assert.Equal(t, true, p.Present)
				assert.Equal(t, "up", p.AdminState)
				assert.Equal(t, "up", p.LinkState)
				assert.Equal(t, "on", p.PoeState)
			}
			}
		}
		hRepo.AssertExpectations(t)
	})
}

func jsonNumber(v int64) string {
	b, _ := json.Marshal(v)
	return string(b)
}
