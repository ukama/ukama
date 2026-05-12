/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	ukamapb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	pb "github.com/ukama/ukama/systems/node/health/pb/gen"
	pbmocks "github.com/ukama/ukama/systems/node/health/pb/gen/mocks"
	"github.com/ukama/ukama/systems/node/node-gateway/pkg/client"
)

const testNodeID = "ukma-0000-tnode-0000"

func TestHealthClientStoreHealthReport(t *testing.T) {
	mc := &pbmocks.HealthServiceClient{}

	req := &pb.StoreHealthReportRequest{
		NodeId:  testNodeID,
		Payload: []byte(`{"nodeType":"tower","schemaVersion":"1.0","reportedAt":"2026-04-21T10:00:00Z","ok":true}`),
	}

	mc.On("StoreHealthReport", mock.Anything, req).Return(&pb.StoreHealthReportResponse{ReportId: "rid-1"}, nil).Once()

	c := client.NewHealthFromClient(mc)
	resp, err := c.StoreHealthReport(req)

	assert.NoError(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, "rid-1", resp.GetReportId())
	}
	mc.AssertExpectations(t)
}

func TestHealthClientList(t *testing.T) {
	mc := &pbmocks.HealthServiceClient{}
	expectedReq := &pb.ListRequest{
		NodeId:    testNodeID,
		Timeframe: ukamapb.FilterTimeframesType_ALL,
	}
	expectedResp := &pb.ListResponse{
		Reports: []*pb.HealthReport{{NodeId: testNodeID}},
	}

	mc.On("List", mock.Anything, expectedReq).Return(expectedResp, nil).Once()

	c := client.NewHealthFromClient(mc)
	resp, err := c.List(expectedReq)

	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.Reports, 1) {
		assert.Equal(t, testNodeID, resp.Reports[0].NodeId)
	}
	mc.AssertExpectations(t)
}

func TestNewHealthFromClientDefaults(t *testing.T) {
	mc := &pbmocks.HealthServiceClient{}

	c := client.NewHealthFromClient(mc)
	assert.NotNil(t, c)
	c.Close()
}

func TestHealthCloseWithNilConn(t *testing.T) {
	mc := &pbmocks.HealthServiceClient{}
	c := client.NewHealthFromClient(mc)

	// Ensure Close is safe when no real gRPC connection exists.
	c.Close()
}

func TestNewHealthCreatesClient(t *testing.T) {
	c := client.NewHealth("localhost:65535", 2*time.Second)
	assert.NotNil(t, c)
	c.Close()
}
