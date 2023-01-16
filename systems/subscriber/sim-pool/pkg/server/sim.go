package server

import (
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"

	"context"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/utils"
)

type SimServer struct {
	simRepo db.SimRepo
	pb.UnimplementedSimServiceServer
}

func NewSimServer(simRepo db.SimRepo) *SimServer {
	return &SimServer{simRepo: simRepo}
}

func (p *SimServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	logrus.Infof("GetSim : %v", req.GetIsPhysicalSim())

	sim, err := p.simRepo.Get(req.GetIsPhysicalSim(), req.GetSimType().String())
	if err != nil {
		logrus.Error("error fetching a sim " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	return &pb.GetResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (p *SimServer) GetByIccid(ctx context.Context, req *pb.GetByIccidRequest) (*pb.GetByIccidResponse, error) {
	logrus.Infof("GetSimByIccid : %v", req.GetIccid())

	sim, err := p.simRepo.GetByIccid(req.GetIccid())
	if err != nil {
		logrus.Error("error fetching a sim " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	return &pb.GetByIccidResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (p *SimServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	logrus.Infof("GetSimStats : %v ", req.GetSimType())
	simType := req.SimType.String()
	if req.GetSimType() == pb.SimType_ANY {
		simType = ""
	}
	sim, err := p.simRepo.GetStats(simType)
	if err != nil {
		logrus.Error("error getting a sim pool stats" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}
	resp := utils.PoolStats(sim)

	return resp, nil
}

func (p *SimServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	logrus.Infof("Add Sims : %v ", req.Sim)
	result := utils.PbParseToModel(req.Sim)
	err := p.simRepo.Add(result)
	if err != nil {
		logrus.Error("error adding a sims" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}
	resp := &pb.AddResponse{Sim: dbSimsToPbSim(result)}
	return resp, nil
}

func (p *SimServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	logrus.Infof("Upload Sims to pool")
	a, _ := utils.ParseBytesToRawSim(req.SimData)
	s := utils.RawSimToPb(a, req.GetSimType().String())
	err := p.simRepo.Add(s)
	if err != nil {
		logrus.Error("error while Upload sims data" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}
	return &pb.UploadResponse{Sim: dbSimsToPbSim(s)}, nil
}

func (p *SimServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	logrus.Infof("Delete Sims: %v", req.GetId())
	err := p.simRepo.Delete(req.GetId())
	if err != nil {
		logrus.Error("error while delete sims data" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}
	return &pb.DeleteResponse{Id: req.GetId()}, nil
}

func dbSimsToPbSim(packages []db.Sim) []*pb.Sim {
	res := []*pb.Sim{}
	for _, u := range packages {
		res = append(res, dbSimToPbSim(&u))
	}
	return res
}

func dbSimToPbSim(p *db.Sim) *pb.Sim {
	return &pb.Sim{
		Id:             uint64(p.ID),
		Iccid:          p.Iccid,
		Msisdn:         p.Msisdn,
		IsAllocated:    p.Is_allocated,
		SmDpAddress:    p.SmDpAddress,
		ActivationCode: p.ActivationCode,
		CreatedAt:      p.CreatedAt.String(),
		UpdatedAt:      p.UpdatedAt.String(),
		DeletedAt:      p.DeletedAt.Time.String(),
		SimType:        pb.SimType(pb.SimType_value[p.Sim_type]),
	}
}
