package interceptor

import (
	"context"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/utils"
	tapb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	getSimRPCSuffix          = "SimManagerService/GetSim"
	toggleSimStatusRPCSuffix = "SimManagerService/ToggleSimStatus"
	deleteSimRPCSuffix       = "SimManagerService/DeleteSim"

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
			if err := utils.ParseTestUUID(rq.SimID); err == nil {
				log.Infof("Calling %q RPC for vSim: %q", info.FullMethod, rq.SimID)

				return f.getSimHandler(ctx, rq.SimID)
			}
		}

	case strings.HasSuffix(info.FullMethod, toggleSimStatusRPCSuffix):
		if rq, ok := req.(*pb.ToggleSimStatusRequest); ok {
			if err := utils.ParseTestUUID(rq.SimID); err == nil {
				log.Infof("Calling %q RPC for vSim: %q", info.FullMethod, rq.SimID)

				return f.toggleSimStatusHandler(ctx, rq.SimID, rq.Status)
			}
		}

	case strings.HasSuffix(info.FullMethod, deleteSimRPCSuffix):
		if rq, ok := req.(*pb.DeleteSimRequest); ok {
			if err := utils.ParseTestUUID(rq.SimID); err == nil {
				log.Infof("Calling %q RPC for vSim: %q", info.FullMethod, rq.SimID)

				return f.deleteSimHandler(ctx, rq.SimID)
			}
		}
	}

	return rpcHandler(ctx, req)
}

func (f *FakeSimInterceptor) getSimHandler(ctx context.Context, simID string) (any, error) {
	fakeIccid, err := utils.GetIccidFromTestUUID(simID)
	if err != nil {
		return nil, err
	}

	resp, err := f.testAgentAdapter.GetSim(ctx, fakeIccid)
	if err != nil {
		return nil, err
	}

	if simInfo, ok := resp.(*tapb.GetSimResponse); ok {
		return &pb.GetSimResponse{Sim: &pb.Sim{
			Id:     simID,
			Iccid:  simInfo.SimInfo.Iccid,
			Status: simInfo.SimInfo.Status,
			Imsi:   simInfo.SimInfo.Imsi,
		}}, nil
	}

	return nil, status.Errorf(codes.Internal, "an unexpected error has occured. Error")
}

func (f *FakeSimInterceptor) toggleSimStatusHandler(ctx context.Context, simID, simStatus string) (any, error) {
	fakeIccid, err := utils.GetIccidFromTestUUID(simID)
	if err != nil {
		return nil, err
	}

	switch simStatus {
	case statusActive:
		return nil, f.testAgentAdapter.ActivateSim(ctx, fakeIccid)
	case statusInactive:
		return nil, f.testAgentAdapter.DeactivateSim(ctx, fakeIccid)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "status %q not supported for operation ", simStatus)
	}
}

func (f *FakeSimInterceptor) deleteSimHandler(ctx context.Context, simID string) (any, error) {
	fakeIccid, err := utils.GetIccidFromTestUUID(simID)
	if err != nil {
		return nil, err
	}

	return nil, f.testAgentAdapter.TerminateSim(ctx, fakeIccid)
}
