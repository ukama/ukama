package server

import (
	"context"
<<<<<<< HEAD
	"strings"

	"github.com/sirupsen/logrus"
=======
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
>>>>>>> subscriber-sys_sim-manager
	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

<<<<<<< HEAD
const iccidPrefix = "890000"

=======
>>>>>>> subscriber-sys_sim-manager
type TestAgentServer struct {
	storage storage.Storage
	pb.UnimplementedTestAgentServiceServer
}

func NewTestAgentServer(storage storage.Storage) *TestAgentServer {
	return &TestAgentServer{
		storage: storage,
	}
}

<<<<<<< HEAD
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
=======
func (s *TestAgentServer) GetSim(ctx context.Context, req *pb.GetSimRequest) (*pb.GetSimResponse, error) {
	log.Infof("GetSim: %+v", req)

	sim, err := s.getOrCreateSimInfo(ctx, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "sim not found.")
		}

		return nil, status.Errorf(codes.Internal, "an unexpected error has occurred")
	}

	return &pb.GetSimResponse{
>>>>>>> subscriber-sys_sim-manager
		SimInfo: &pb.SimInfo{
			Iccid:  sim.Iccid,
			Imsi:   sim.Imsi,
			Status: sim.Status,
		},
	}, nil
}

func (s *TestAgentServer) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
<<<<<<< HEAD
	logrus.Infof("Activate sim for iccid: %s", req.Iccid)
	sim := s.getSim(ctx, req.Iccid)
	if sim == nil {
=======
	log.Infof("Activate sim for iccid: %s", req.Iccid)

	sim, err := s.getSimInfo(ctx, req.Iccid)
	if err != nil {
>>>>>>> subscriber-sys_sim-manager
		return nil, status.Errorf(codes.NotFound, "sim not found.")
	}

	if sim.Status != storage.SimStatusInactive.String() {
<<<<<<< HEAD
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for deletion", sim.Status)
=======
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for operation", sim.Status)
>>>>>>> subscriber-sys_sim-manager
	}

	sim.Status = storage.SimStatusActive.String()

<<<<<<< HEAD
	err := s.storage.Put(req.Iccid, sim)
=======
	err = s.storage.Put(req.Iccid, sim)
>>>>>>> subscriber-sys_sim-manager
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update sim info in storage: %v", err)
	}

	return &pb.ActivateSimResponse{}, nil
}

func (s *TestAgentServer) DeactivateSim(ctx context.Context, req *pb.DeactivateSimRequest) (*pb.DeactivateSimResponse, error) {
<<<<<<< HEAD
	logrus.Infof("Deactivate sim for iccid: %s", req.Iccid)
	sim := s.getSim(ctx, req.Iccid)
	if sim == nil {
=======
	log.Infof("Deactivate sim for iccid: %s", req.Iccid)

	sim, err := s.getSimInfo(ctx, req.Iccid)
	if err != nil {
>>>>>>> subscriber-sys_sim-manager
		return nil, status.Errorf(codes.NotFound, "sim not found.")
	}

	if sim.Status != storage.SimStatusActive.String() {
<<<<<<< HEAD
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for deletion", sim.Status)
=======
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for operation", sim.Status)
>>>>>>> subscriber-sys_sim-manager
	}

	sim.Status = storage.SimStatusInactive.String()

<<<<<<< HEAD
	err := s.storage.Put(req.Iccid, sim)
=======
	err = s.storage.Put(req.Iccid, sim)
>>>>>>> subscriber-sys_sim-manager
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update sim info in storage: %v", err)
	}

	return &pb.DeactivateSimResponse{}, nil
}

func (s *TestAgentServer) TerminateSim(ctx context.Context, req *pb.TerminateSimRequest) (*pb.TerminateSimResponse, error) {
<<<<<<< HEAD
	logrus.Infof("Terminate sim for iccid: %s", req.Iccid)
	sim := s.getSim(ctx, req.Iccid)
	if sim == nil {
=======
	log.Infof("Terminate sim for iccid: %s", req.Iccid)

	sim, err := s.getSimInfo(ctx, req.Iccid)
	if err != nil {
>>>>>>> subscriber-sys_sim-manager
		return nil, status.Errorf(codes.NotFound, "sim not found.")
	}

	if sim.Status != storage.SimStatusInactive.String() {
		return nil, status.Errorf(codes.FailedPrecondition, "invalid sim state %q for deletion", sim.Status)
	}

<<<<<<< HEAD
	err := s.storage.Delete(req.Iccid)
=======
	err = s.storage.Delete(req.Iccid)
>>>>>>> subscriber-sys_sim-manager
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot delete sim info from storage: %v", err)
	}
	return &pb.TerminateSimResponse{}, nil
}

<<<<<<< HEAD
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
=======
func (s *TestAgentServer) getOrCreateSimInfo(ctx context.Context, req *pb.GetSimRequest) (*storage.SimInfo, error) {
	log.Infof("Get sim info for iccid: %s", req.Iccid)

	sim, err := s.getSimInfo(ctx, req.Iccid)

	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return nil, err
		}

		log.Infof("Sim info for iccid: %s does not exist. Adding new sim info to Storage", req.Iccid)
		imsi := req.Iccid[len(req.Iccid)-15:]

		sim = &storage.SimInfo{
			Iccid:  req.Iccid,
			Imsi:   imsi,
			Status: storage.SimStatusInactive.String(),
		}

		err := s.storage.Put(req.Iccid, sim)
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
>>>>>>> subscriber-sys_sim-manager
}
