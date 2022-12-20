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
	logrus.Infof("GetPoolStats : %v ", req.GetOrgId())
	_, err := p.simPoolRepo.GetStats(req.GetOrgId(), req.GetSimType().String())

	if err != nil {
		logrus.Error("error getting a simPool" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "simPool")
	}

	// resp := &pb.GetStatsResponse{}

	return nil, nil
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
	logrus.Infof("Upload SimPool: %v", req.GetFileUrl())

	var s []db.SimPool
	err := p.simPoolRepo.Add(s)
	if err != nil {
		logrus.Error("error while Upload simPool data" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "simPool")
	}
	return &pb.UploadResponse{SimPool: dbSimPoolsToPbSimPool(s)}, nil
}

func dbSimPoolsToPbSimPool(packages []db.SimPool) []*pb.SimPool {
	res := []*pb.SimPool{}
	for _, u := range packages {
		res = append(res, dbSimPoolToPbSimPool(&u))
	}
	return res
}

func dbSimPoolToPbSimPool(p *db.SimPool) *pb.SimPool {
	return &pb.SimPool{}
}
