package server

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestAgentServer struct {
	storage storage.Storage
	pb.UnimplementedTestAgentServiceServer
}

func NewTestAgentServer(storage storage.Storage) *TestAgentServer {
	return &TestAgentServer{
		storage: storage,
	}
}

func (s *TestAgentServer) BindSim(ctx context.Context, req *pb.BindSimRequest) (*pb.BindSimResponse, error) {
	log.Infof("BindSim: %+v", req)

	_, err := s.getOrCreateSimInfo(ctx, req.Iccid)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "sim not found.")
		}

		return nil, status.Errorf(codes.Internal, "an unexpected error has occurred")
	}

	return &pb.BindSimResponse{}, nil
}

func (s *TestAgentServer) GetSim(ctx context.Context, req *pb.GetSimRequest) (*pb.GetSimResponse, error) {
	log.Infof("GetSim: %+v", req)

	sim, err := s.getOrCreateSimInfo(ctx, req.Iccid)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "sim not found.")
		}

		return nil, status.Errorf(codes.Internal, "an unexpected error has occurred")
	}

	return &pb.GetSimResponse{
		SimInfo: &pb.SimInfo{
			Iccid:  sim.Iccid,
			Imsi:   sim.Imsi,
			Status: sim.Status,
		},
	}, nil
}

func (s *TestAgentServer) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
	log.Infof("Activate sim for iccid: %s", req.Iccid)

	sim, err := s.getSimInfo(ctx, req.Iccid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "sim not found.")
	}

	if sim.Status != storage.SimStatusInactive.String() {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for operation", sim.Status)
	}

	sim.Status = storage.SimStatusActive.String()

	err = s.storage.Put(req.Iccid, sim)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update sim info in storage: %v", err)
	}

	return &pb.ActivateSimResponse{}, nil
}

func (s *TestAgentServer) DeactivateSim(ctx context.Context, req *pb.DeactivateSimRequest) (*pb.DeactivateSimResponse, error) {
	log.Infof("Deactivate sim for iccid: %s", req.Iccid)

	sim, err := s.getSimInfo(ctx, req.Iccid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "sim not found.")
	}

	if sim.Status != storage.SimStatusActive.String() {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for operation", sim.Status)
	}

	sim.Status = storage.SimStatusInactive.String()

	err = s.storage.Put(req.Iccid, sim)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update sim info in storage: %v", err)
	}

	return &pb.DeactivateSimResponse{}, nil
}

func (s *TestAgentServer) TerminateSim(ctx context.Context, req *pb.TerminateSimRequest) (*pb.TerminateSimResponse, error) {
	log.Infof("Terminate sim for iccid: %s", req.Iccid)

	sim, err := s.getSimInfo(ctx, req.Iccid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "sim not found.")
	}

	if sim.Status != storage.SimStatusInactive.String() {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for deletion", sim.Status)
	}

	err = s.storage.Delete(req.Iccid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot delete sim info from storage: %v", err)
	}
	return &pb.TerminateSimResponse{}, nil
}

func (s *TestAgentServer) getOrCreateSimInfo(ctx context.Context, iccid string) (*storage.SimInfo, error) {
	log.Infof("Get sim info for iccid: %s", iccid)

	sim, err := s.getSimInfo(ctx, iccid)

	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return nil, err
		}

		log.Infof("Sim info for iccid: %s does not exist. Adding new sim info to Storage", iccid)
		imsi := iccid[len(iccid)-15:]

		sim = &storage.SimInfo{
			Iccid:  iccid,
			Imsi:   imsi,
			Status: storage.SimStatusInactive.String(),
		}

		err := s.storage.Put(iccid, sim)
		if err != nil {
			return nil, fmt.Errorf("cannot add sim info into storage: %w", err)
		}
	}

	return sim, nil
}

func (s *TestAgentServer) getSimInfo(ctx context.Context, iccid string) (*storage.SimInfo, error) {
	sim, err := s.storage.Get(iccid)
	if err != nil {
		return nil, fmt.Errorf("cannot get sim info from storage: %w", err)
	}

	return sim, nil
}
