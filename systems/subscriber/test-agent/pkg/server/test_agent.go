package server

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestAgentServer struct {
	storage storage.Storage
	pb.UnimplementedTestAgentServiceServer
}

func NewTestAgentServer() *TestAgentServer {
	return &TestAgentServer{}
}

func (s *TestAgentServer) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "cannot activate sim %s: method TestAgent.ActivateSim not implemented", req.Iccid)
}

func (s *TestAgentServer) DeactivateSim(ctx context.Context, req *pb.DeactivateSimRequest) (*pb.DeactivateSimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "cannot deactivate sim %s: method TestAgent.DeactivateSim not implemented", req.Iccid)
}

func (s *TestAgentServer) TerminateSim(ctx context.Context, req *pb.TerminateSimRequest) (*pb.TerminateSimResponse, error) {
	logrus.Infof("Terminate sim for iccid: %s", req.Iccid)
	sim := s.getSimInfo(ctx, req.Iccid)
	if sim == nil {
		return nil, status.Errorf(codes.NotFound, "Sim not found.")
	}

	err := s.storage.Delete(req.Iccid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot delete sim info from etcd: %v", err)
	}
	return &pb.TerminateSimResponse{}, nil
}

func (s *TestAgentServer) getSimInfo(ctx context.Context, iccid string) *storage.SimInfo {
	val, err := s.storage.Get(iccid)
	if err != nil {
		logrus.Errorf("Cannot get sim info from etcd: %v", err)
		return nil
	}

	var sim *storage.SimInfo
	if val != nil {
		sim = unmarshalSimInfo(val)
	} else {
		sim = nil
	}
	return sim
}

func marshalSimInfo(info *storage.SimInfo) string {
	b, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		logrus.Errorf("Cannot marshal sim info: %v", err)
	}
	return string(b)
}

func unmarshalSimInfo(b []byte) *storage.SimInfo {
	info := storage.SimInfo{}
	err := json.Unmarshal(b, &info)
	if err != nil {
		logrus.Errorf("Cannot unmarshal sim info: %v", err)
		return nil
	}
	return &info
}
