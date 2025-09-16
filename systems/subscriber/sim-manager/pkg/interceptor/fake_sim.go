/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package interceptor

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/utils"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	tapb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
)

const (
	getSimRPCSuffix          = "SimManagerService/GetSim"
	toggleSimStatusRPCSuffix = "SimManagerService/ToggleSimStatus"
	terminateSimRPCSuffix    = "SimManagerService/TerminateSim"

	statusActive   = "active"
	statusInactive = "inactive"
)

type FakeSimInterceptor struct {
	testAgentAdapter adapters.AgentAdapter
}

func NewFakeSimInterceptor(testAgentHost string, timeout time.Duration) *FakeSimInterceptor {
	agent, err := adapters.NewTestAgentAdapter(testAgentHost, timeout)
	if err != nil {
		log.Fatalf("Failed to connect to Agent service at %s. Error: %v", testAgentHost, err)
	}

	return &FakeSimInterceptor{
		testAgentAdapter: agent,
	}
}

func (f *FakeSimInterceptor) UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, rpcHandler grpc.UnaryHandler) (any, error) {
	switch {
	case strings.HasSuffix(info.FullMethod, getSimRPCSuffix):
		if rq, ok := req.(*pb.GetSimRequest); ok {
			if err := utils.ParseTestUUID(rq.SimId); err == nil {
				log.Infof("Calling %q RPC for vSim: %q", info.FullMethod, rq.SimId)

				return f.getSimHandler(ctx, rq.SimId)
			}
		}

	case strings.HasSuffix(info.FullMethod, toggleSimStatusRPCSuffix):
		if rq, ok := req.(*pb.ToggleSimStatusRequest); ok {
			if err := utils.ParseTestUUID(rq.SimId); err == nil {
				log.Infof("Calling %q RPC for vSim: %q", info.FullMethod, rq.SimId)

				return f.toggleSimStatusHandler(ctx, rq.SimId, rq.Status)
			}
		}

	case strings.HasSuffix(info.FullMethod, terminateSimRPCSuffix):
		if rq, ok := req.(*pb.TerminateSimRequest); ok {
			if err := utils.ParseTestUUID(rq.SimId); err == nil {
				log.Infof("Calling %q RPC for vSim: %q", info.FullMethod, rq.SimId)

				return f.deleteSimHandler(ctx, rq.SimId)
			}
		}
	}

	return rpcHandler(ctx, req)
}

func (f *FakeSimInterceptor) getSimHandler(ctx context.Context, simId string) (any, error) {
	fakeIccid, err := utils.GetIccidFromTestUUID(simId)
	if err != nil {
		return nil, err
	}

	resp, err := f.testAgentAdapter.GetSim(ctx, fakeIccid)
	if err != nil {
		return nil, err
	}

	if simInfo, ok := resp.(*tapb.GetSimResponse); ok {
		return &pb.GetSimResponse{Sim: &pb.Sim{
			Id:     simId,
			Iccid:  simInfo.SimInfo.Iccid,
			Status: simInfo.SimInfo.Status,
			Imsi:   simInfo.SimInfo.Imsi,
		}}, nil
	}

	return nil, status.Errorf(codes.Internal, "an unexpected error has occurred.")
}

func (f *FakeSimInterceptor) toggleSimStatusHandler(ctx context.Context, simId, simStatus string) (any, error) {
	fakeIccid, err := utils.GetIccidFromTestUUID(simId)
	if err != nil {
		return nil, err
	}

	switch simStatus {
	case statusActive:
		return nil, f.testAgentAdapter.ActivateSim(ctx, client.AgentRequestData{
			Iccid: fakeIccid,
		})
	case statusInactive:
		return nil, f.testAgentAdapter.DeactivateSim(ctx, client.AgentRequestData{
			Iccid: fakeIccid,
		})
	default:
		return nil, status.Errorf(codes.InvalidArgument, "status %q not supported for operation ", simStatus)
	}
}

func (f *FakeSimInterceptor) deleteSimHandler(ctx context.Context, simId string) (any, error) {
	fakeIccid, err := utils.GetIccidFromTestUUID(simId)
	if err != nil {
		return nil, err
	}

	return nil, f.testAgentAdapter.TerminateSim(ctx, fakeIccid)
}
