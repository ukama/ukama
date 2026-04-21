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

func TestHealthClientStoreRunningAppsInfo(t *testing.T) {
	mc := &pbmocks.HealhtServiceClient{}

	req := &pb.StoreRunningAppsInfoRequest{
		NodeId:    testNodeID,
		Timestamp: "2026-04-21T10:00:00Z",
		System: []*pb.System{
			{Name: "cpu", Value: "30"},
		},
		Capps: []*pb.Capps{
			{
				Space:  "core",
				Name:   "agent",
				Tag:    "v1.0.0",
				Status: pb.Status_ACTIVE,
				Resources: []*pb.Resource{
					{Name: "mem", Value: "100mb"},
				},
			},
		},
	}

	expectedReq := &pb.StoreRunningAppsInfoRequest{
		NodeId:    req.NodeId,
		Timestamp: req.Timestamp,
		System: []*pb.System{
			{Name: "cpu", Value: "30"},
		},
		Capps: []*pb.Capps{
			{
				Space:  "core",
				Name:   "agent",
				Tag:    "v1.0.0",
				Status: pb.Status_ACTIVE,
				Resources: []*pb.Resource{
					{Name: "mem", Value: "100mb"},
				},
			},
		},
	}

	mc.On("StoreRunningAppsInfo", mock.Anything, expectedReq).Return(&pb.StoreRunningAppsInfoResponse{}, nil).Once()

	c := client.NewHealthFromClient(mc)
	resp, err := c.StoreRunningAppsInfo(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mc.AssertExpectations(t)
}

func TestHealthClientList(t *testing.T) {
	mc := &pbmocks.HealhtServiceClient{}
	expectedReq := &pb.ListRequest{
		NodeId: testNodeID,
		Filter: ukamapb.FilterTimestampType_ALL,
	}
	expectedResp := &pb.ListResponse{
		Healths: []*pb.Health{
			{NodeId: testNodeID},
		},
	}

	mc.On("List", mock.Anything, expectedReq).Return(expectedResp, nil).Once()

	c := client.NewHealthFromClient(mc)
	resp, err := c.List(expectedReq)

	assert.NoError(t, err)
	if assert.NotNil(t, resp) && assert.Len(t, resp.Healths, 1) {
		assert.Equal(t, testNodeID, resp.Healths[0].NodeId)
	}
	mc.AssertExpectations(t)
}

func TestNewHealthFromClientDefaults(t *testing.T) {
	mc := &pbmocks.HealhtServiceClient{}

	c := client.NewHealthFromClient(mc)
	assert.NotNil(t, c)
	c.Close()
}

func TestHealthCloseWithNilConn(t *testing.T) {
	mc := &pbmocks.HealhtServiceClient{}
	c := client.NewHealthFromClient(mc)

	// Ensure Close is safe when no real gRPC connection exists.
	c.Close()
}

func TestNewHealthCreatesClient(t *testing.T) {
	c := client.NewHealth("localhost:65535", 2*time.Second)
	assert.NotNil(t, c)
	c.Close()
}
