package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/utils"
)

type SimPoolServer struct {
	simRepo        db.SimRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedSimServiceServer
}

func NewSimPoolServer(simRepo db.SimRepo, msgBus mb.MsgBusServiceClient) *SimPoolServer {
	return &SimPoolServer{
		simRepo:        simRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)}
}

func (p *SimPoolServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	logrus.Infof("GetSim isPhysical: %v, simType: %v", req.GetIsPhysicalSim(), req.GetSimType())

	sim, err := p.simRepo.Get(req.GetIsPhysicalSim(), db.ParseType(req.GetSimType()))
	if err != nil {
		logrus.Error("error fetching a sim " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	return &pb.GetResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (p *SimPoolServer) GetByIccid(ctx context.Context, req *pb.GetByIccidRequest) (*pb.GetByIccidResponse, error) {
	logrus.Infof("GetSimByIccid : %v", req.GetIccid())

	sim, err := p.simRepo.GetByIccid(req.GetIccid())
	if err != nil {
		logrus.Error("error fetching a sim " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	return &pb.GetByIccidResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (p *SimPoolServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	logrus.Infof("GetSimStats : %v ", req.GetSimType())

	sim, err := p.simRepo.GetSimsByType(db.ParseType(req.SimType))
	if err != nil {
		logrus.Error("error getting a sim pool stats" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}
	resp := utils.PoolStats(sim)

	return resp, nil
}

func (p *SimPoolServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
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

func (p *SimPoolServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	logrus.Infof("Upload Sims to pool")
	a, _ := utils.ParseBytesToRawSim(req.SimData)
	s := utils.RawSimToPb(a, req.GetSimType())
	err := p.simRepo.Add(s)
	logrus.Info("ADDING SIMS: ", s, err)
	if err != nil {
		logrus.Error("error while Upload sims data" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}
	route := p.baseRoutingKey.SetAction("upload").SetObject("sim").MustBuild()
	err = p.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}
	iccids := []string{}
	for _, u := range s {
		iccids = append(iccids, u.Iccid)
	}

	return &pb.UploadResponse{Iccid: iccids}, nil
}

func (p *SimPoolServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
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
		IsAllocated:    p.IsAllocated,
		SimType:        p.SimType.String(),
		SmDpAddress:    p.SmDpAddress,
		ActivationCode: p.ActivationCode,
		CreatedAt:      p.CreatedAt.String(),
		DeletedAt:      p.DeletedAt.Time.String(),
		UpdatedAt:      p.UpdatedAt.String(),
		IsPhysical:     p.IsPhysical,
		QrCode:         p.QrCode,
		IsFailed:       p.IsFailed,
	}
}
