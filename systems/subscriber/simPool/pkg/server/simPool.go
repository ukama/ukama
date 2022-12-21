package server

import (
	pb "github.com/ukama/ukama/systems/subscriber/simPool/pb/gen"

	"context"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/subscriber/simPool/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/simPool/pkg/utils"
)

type SimPoolServer struct {
	simPoolRepo db.SimPoolRepo
	pb.UnimplementedSimPoolServiceServer
}

func NewSimPoolServer(simPoolRepo db.SimPoolRepo) *SimPoolServer {
	return &SimPoolServer{simPoolRepo: simPoolRepo}
}

func (p *SimPoolServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	logrus.Infof("GetPoolStats : %v ", req.GetSimType())
	simPool, err := p.simPoolRepo.GetStats(req.GetSimType().String())

	if err != nil {
		logrus.Error("error getting a simPool" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "simPool")
	}
	resp := utils.SimPoolStats(simPool)

	return resp, nil
}

func (p *SimPoolServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	logrus.Infof("Add SimPool : %v ", req.SimPool)
	result := utils.PbParseToModel(req.SimPool)
	err := p.simPoolRepo.Add(result)
	if err != nil {
		logrus.Error("error adding a simPool" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "simPool")
	}
	resp := &pb.AddResponse{SimPool: dbSimPoolsToPbSimPool(result)}
	return resp, nil
}

func (p *SimPoolServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	// logrus.Infof("Upload SimPool: %v", req.GetFileUrl())
	var s []db.SimPool
	err := p.simPoolRepo.Add(s)
	if err != nil {
		logrus.Error("error while Upload simPool data" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "simPool")
	}
	return &pb.UploadResponse{SimPool: dbSimPoolsToPbSimPool(s)}, nil
}

func (p *SimPoolServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	logrus.Infof("Delete SimPool: %v", req.GetId())
	err := p.simPoolRepo.Delete(req.GetId())
	if err != nil {
		logrus.Error("error while delete simPool data" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "simPool")
	}
	return &pb.DeleteResponse{Id: req.GetId()}, nil
}

func dbSimPoolsToPbSimPool(packages []db.SimPool) []*pb.SimPool {
	res := []*pb.SimPool{}
	for _, u := range packages {
		res = append(res, dbSimPoolToPbSimPool(&u))
	}
	return res
}

func dbSimPoolToPbSimPool(p *db.SimPool) *pb.SimPool {
	return &pb.SimPool{
		Id:             uint64(p.ID),
		Iccid:          p.Iccid,
		Msisdn:         p.Msisdn,
		IsAllocated:    p.Is_allocated,
		SmDpAddress:    p.SmDpAddress,
		ActivationCode: p.ActivationCode,
		QrCode:         p.QrCode,
		CreatedAt:      p.CreatedAt.String(),
		UpdatedAt:      p.UpdatedAt.String(),
		DeletedAt:      p.DeletedAt.Time.String(),
		SimType:        pb.SimType(pb.SimType_value[p.Sim_type]),
	}
}
