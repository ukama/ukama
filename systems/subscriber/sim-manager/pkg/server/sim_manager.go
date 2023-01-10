package server

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients"

	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
)

type SimManagerServer struct {
	pb.UnimplementedSimManagerServiceServer
	simRepo      sims.SimRepo
	packageRepo  sims.PackageRepo
	agentFactory *clients.AgentFactory
}

func NewSimManagerServer(simRepo sims.SimRepo, packageRepo sims.PackageRepo, agentFactory *clients.AgentFactory) *SimManagerServer {
	return &SimManagerServer{
		simRepo:      simRepo,
		packageRepo:  packageRepo,
		agentFactory: agentFactory,
	}
}

func (s *SimManagerServer) GetSim(ctx context.Context, req *pb.GetSimRequest) (*pb.GetSimResponse, error) {
	simID, err := uuid.Parse(req.GetSimID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	return &pb.GetSimResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (s *SimManagerServer) GetSimsBySubscriber(ctx context.Context, req *pb.GetSimsBySubscriberRequest) (*pb.GetSimsBySubscriberResponse, error) {
	subID, err := uuid.Parse(req.GetSubscriberID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of subscriber uuid. Error %s", err.Error())
	}

	sims, err := s.simRepo.GetBySubscriber(subID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	resp := &pb.GetSimsBySubscriberResponse{
		SubscriberID: req.GetSubscriberID(),
		Sims:         dbSimsToPbSims(sims),
	}

	return resp, nil
}

func (s *SimManagerServer) GetSimsByNetwork(ctx context.Context, req *pb.GetSimsByNetworkRequest) (*pb.GetSimsByNetworkResponse, error) {
	netID, err := uuid.Parse(req.GetNetworkID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of network uuid. Error %s", err.Error())
	}

	sims, err := s.simRepo.GetByNetwork(netID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sims")
	}

	resp := &pb.GetSimsByNetworkResponse{
		NetworkID: req.GetNetworkID(),
		Sims:      dbSimsToPbSims(sims),
	}

	return resp, nil
}

func (s *SimManagerServer) ActivateSim(ctx context.Context, req *pb.ActivateSimRequest) (*pb.ActivateSimResponse, error) {
	simID, err := uuid.Parse(req.GetSimID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Status != sims.SimStatusInactive {
		return nil, status.Errorf(codes.FailedPrecondition, "sim's state %s is invalid for activation: Error %s", sim.Status, err.Error())
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sim type %q for sim ID %q", sim.Type, req.SimID)
	}

	err = simAgent.ActivateSim(ctx, req.SimID)
	if err != nil {
		return nil, err
	}

	simUpdates := &sims.Sim{
		ID:               sim.ID,
		Status:           sims.SimStatusActive,
		ActivationsCount: sim.ActivationsCount + 1,
		LastActivatedOn:  time.Now(),
	}

	if sim.FirstActivatedOn.IsZero() {
		simUpdates.FirstActivatedOn = simUpdates.LastActivatedOn
	}

	err = s.simRepo.Update(simUpdates, nil)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	return &pb.ActivateSimResponse{}, nil
}

func (s *SimManagerServer) DeactivateSim(ctx context.Context, req *pb.DeactivateSimRequest) (*pb.DeactivateSimResponse, error) {
	simID, err := uuid.Parse(req.GetSimID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Status != sims.SimStatusActive {
		return nil, status.Errorf(codes.FailedPrecondition, "sim's state %s is invalid for deactivation: Error %s", sim.Status, err.Error())
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sim type:%q for sim ID %q", sim.Type, req.SimID)
	}

	err = simAgent.DeactivateSim(ctx, req.SimID)
	if err != nil {
		return nil, err
	}

	simUpdates := &sims.Sim{
		ID:                 sim.ID,
		Status:             sims.SimStatusInactive,
		DeactivationsCount: sim.DeactivationsCount + 1}

	err = s.simRepo.Update(simUpdates, nil)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	return &pb.DeactivateSimResponse{}, nil
}

func (s *SimManagerServer) DeleteSim(ctx context.Context, req *pb.DeleteSimRequest) (*pb.DeleteSimResponse, error) {
	simID, err := uuid.Parse(req.GetSimID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of sim uuid. Error %s", err.Error())
	}

	sim, err := s.simRepo.Get(simID)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	if sim.Status != sims.SimStatusInactive {
		return nil, status.Errorf(codes.FailedPrecondition, "sim's state %s is invalid for deletion: Error %s", sim.Status, err.Error())
	}

	simAgent, ok := s.agentFactory.GetAgentAdapter(sim.Type)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sim type:%q for sim ID %q", sim.Type, req.SimID)
	}

	err = simAgent.TerminateSim(ctx, req.SimID)
	if err != nil {
		return nil, err
	}

	// update sim status & mark sim as deleted
	simUpdates := &sims.Sim{
		ID:           sim.ID,
		Status:       sims.SimStatusTerminated,
		TerminatedAt: gorm.DeletedAt{Time: time.Now(), Valid: true},
	}

	err = s.simRepo.Update(simUpdates, nil)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "sim")
	}

	return &pb.DeleteSimResponse{}, nil
}

func (s *SimManagerServer) RemovePackageForSim(ctx context.Context, req *pb.RemovePackageRequest) (*pb.RemovePackageResponse, error) {
	packageID, err := uuid.Parse(req.GetPackageID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format of package uuid. Error %s", err.Error())
	}

	err = s.packageRepo.Delete(packageID, nil)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "package")
	}

	return &pb.RemovePackageResponse{}, nil
}

func dbSimToPbSim(sim *sims.Sim) *pb.Sim {
	res := &pb.Sim{
		Id:                 sim.ID.String(),
		SubscriberID:       sim.SubscriberID.String(),
		NetworkID:          sim.NetworkID.String(),
		Iccid:              sim.Iccid,
		Msisdn:             sim.Msisdn,
		Imsi:               sim.Imsi,
		Type:               sim.Type.String(),
		Status:             sim.Status.String(),
		IsPhysical:         sim.IsPhysical,
		ActivationsCount:   sim.ActivationsCount,
		DeactivationsCount: sim.DeactivationsCount,
	}

	if sim.Package.ID != uuid.Nil {
		res.Package = dbPackageToPbPackage(&sim.Package)
	}

	if !sim.FirstActivatedOn.IsZero() {
		res.FirstActivatedOn = timestamppb.New(sim.FirstActivatedOn)
	}

	if !sim.LastActivatedOn.IsZero() {
		res.LastActivatedOn = timestamppb.New(sim.LastActivatedOn)
	}

	if sim.AllocatedAt != 0 {
		res.AllocatedAt = timestamppb.New(sim.LastActivatedOn)
	}

	return res
}

func dbPackageToPbPackage(pkg *sims.Package) *pb.Package {
	res := &pb.Package{
		Id: pkg.ID.String(),
	}

	if !pkg.EndDate.IsZero() {
		res.EndDate = timestamppb.New(pkg.EndDate)
	}

	if !pkg.StartDate.IsZero() {
		res.StartDate = timestamppb.New(pkg.StartDate)
	}

	return res
}

func dbSimsToPbSims(sims []sims.Sim) []*pb.Sim {
	res := []*pb.Sim{}

	for _, s := range sims {
		res = append(res, dbSimToPbSim(&s))
	}

	return res
}
