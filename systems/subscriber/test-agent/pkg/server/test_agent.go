package server

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const iccidPrefix = "890000"

type TestAgentServer struct {
	storage storage.Storage
	pb.UnimplementedTestAgentServiceServer
}

func NewTestAgentServer(storage storage.Storage) *TestAgentServer {
	return &TestAgentServer{
		storage: storage,
	}
}

func (s *TestAgentServer) GetSimInfo(ctx context.Context, req *pb.GetSimInfoRequest) (*pb.GetSimInfoResponse, error) {
	logrus.Infof("GetSimInfo: %+v", req)
	if !strings.HasPrefix(req.Iccid, iccidPrefix) {
		return nil, status.Errorf(codes.NotFound, "Sim with iccid %q not found. Test sim iccid should start with: %q", req.Iccid, iccidPrefix)
	}
	iccid := req.Iccid

	sim, err := s.getOrCreateSim(ctx, req, iccid)
	if err != nil {
		return nil, err
	}

	return &pb.GetSimInfoResponse{
		SimInfo: &pb.SimInfo{
			Iccid:  sim.Iccid,
			Imsi:   sim.Imsi,
			Status: sim.Status,
		},
	}, nil
}

func (s *TestAgentServer) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
	logrus.Infof("Activate sim for iccid: %s", req.Iccid)
	sim := s.getSim(ctx, req.Iccid)
	if sim == nil {
		return nil, status.Errorf(codes.NotFound, "sim not found.")
	}

	if sim.Status != storage.SimStatusInactive.String() {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for deletion", sim.Status)
	}

	sim.Status = storage.SimStatusActive.String()

	err := s.storage.Put(req.Iccid, sim)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update sim info in storage: %v", err)
	}

	return &pb.ActivateSimResponse{}, nil
}

func (s *TestAgentServer) DeactivateSim(ctx context.Context, req *pb.DeactivateSimRequest) (*pb.DeactivateSimResponse, error) {
	logrus.Infof("Deactivate sim for iccid: %s", req.Iccid)
	sim := s.getSim(ctx, req.Iccid)
	if sim == nil {
		return nil, status.Errorf(codes.NotFound, "sim not found.")
	}

	if sim.Status != storage.SimStatusActive.String() {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for deletion", sim.Status)
	}

	sim.Status = storage.SimStatusInactive.String()

	err := s.storage.Put(req.Iccid, sim)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update sim info in storage: %v", err)
	}

	return &pb.DeactivateSimResponse{}, nil
}

func (s *TestAgentServer) TerminateSim(ctx context.Context, req *pb.TerminateSimRequest) (*pb.TerminateSimResponse, error) {
	logrus.Infof("Terminate sim for iccid: %s", req.Iccid)
	sim := s.getSim(ctx, req.Iccid)
	if sim == nil {
		return nil, status.Errorf(codes.NotFound, "sim not found.")
	}

	if sim.Status != storage.SimStatusInactive.String() {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for deletion", sim.Status)
	}

	err := s.storage.Delete(req.Iccid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot delete sim info from storage: %v", err)
	}
	return &pb.TerminateSimResponse{}, nil
}

func (s *TestAgentServer) getOrCreateSim(ctx context.Context, req *pb.GetSimInfoRequest, iccid string) (*storage.SimInfo, error) {
	logrus.Infof("Get sim info for iccid: %s", iccid)
	sim := s.getSim(ctx, req.Iccid)
	if sim == nil {

		imsi := req.Iccid[len(iccid)-15:]
		sim = &storage.SimInfo{
			Iccid:  iccid,
			Imsi:   imsi,
			Status: storage.SimStatusInactive.String(),
		}
	}

	err := s.storage.Put(req.Iccid, sim)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update sim info in storage: %v", err)
	}
	return sim, nil
}

func (s *TestAgentServer) getSim(ctx context.Context, iccid string) *storage.SimInfo {
	val, err := s.storage.Get(iccid)
	if err != nil {
		logrus.Errorf("cannot get sim info from storage: %v", err)
		return nil
	}

	var sim *storage.SimInfo
	if val != nil {
		sim = val
	} else {
		sim = nil
	}
	return sim
}
