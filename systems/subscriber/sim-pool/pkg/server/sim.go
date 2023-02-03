package server

import (
<<<<<<< HEAD
=======
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"

>>>>>>> subscriber-sys_sim-manager
	"context"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
<<<<<<< HEAD
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg"
=======
>>>>>>> subscriber-sys_sim-manager
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/utils"
)

<<<<<<< HEAD
type SimPoolServer struct {
	simRepo        db.SimRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedSimServiceServer
}

func NewSimPoolServer(simRepo db.SimRepo, msgBus mb.MsgBusServiceClient) *SimPoolServer {
	return &SimPoolServer{simRepo: simRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName)}
}

func (p *SimPoolServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	logrus.Infof("GetSim isPhysical: %v, simType: %v", req.GetIsPhysicalSim(), req.GetSimType())
=======
type SimServer struct {
	simRepo db.SimRepo
	pb.UnimplementedSimServiceServer
}

func NewSimServer(simRepo db.SimRepo) *SimServer {
	return &SimServer{simRepo: simRepo}
}

func (p *SimServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	logrus.Infof("GetSim : %v", req.GetIsPhysicalSim())
>>>>>>> subscriber-sys_sim-manager

	sim, err := p.simRepo.Get(req.GetIsPhysicalSim(), req.GetSimType().String())
	if err != nil {
		logrus.Error("error fetching a sim " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	return &pb.GetResponse{Sim: dbSimToPbSim(sim)}, nil
}

<<<<<<< HEAD
func (p *SimPoolServer) GetByIccid(ctx context.Context, req *pb.GetByIccidRequest) (*pb.GetByIccidResponse, error) {
=======
func (p *SimServer) GetByIccid(ctx context.Context, req *pb.GetByIccidRequest) (*pb.GetByIccidResponse, error) {
>>>>>>> subscriber-sys_sim-manager
	logrus.Infof("GetSimByIccid : %v", req.GetIccid())

	sim, err := p.simRepo.GetByIccid(req.GetIccid())
	if err != nil {
		logrus.Error("error fetching a sim " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	return &pb.GetByIccidResponse{Sim: dbSimToPbSim(sim)}, nil
}

<<<<<<< HEAD
func (p *SimPoolServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
=======
func (p *SimServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
>>>>>>> subscriber-sys_sim-manager
	logrus.Infof("GetSimStats : %v ", req.GetSimType())
	simType := req.SimType.String()
	if req.GetSimType() == pb.SimType_ANY {
		simType = ""
	}
<<<<<<< HEAD
	sim, err := p.simRepo.GetSimsByType(simType)
=======
	sim, err := p.simRepo.GetStats(simType)
>>>>>>> subscriber-sys_sim-manager
	if err != nil {
		logrus.Error("error getting a sim pool stats" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}
	resp := utils.PoolStats(sim)

	return resp, nil
}

<<<<<<< HEAD
func (p *SimPoolServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
=======
func (p *SimServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
>>>>>>> subscriber-sys_sim-manager
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

<<<<<<< HEAD
func (p *SimPoolServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
=======
func (p *SimServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
>>>>>>> subscriber-sys_sim-manager
	logrus.Infof("Upload Sims to pool")
	a, _ := utils.ParseBytesToRawSim(req.SimData)
	s := utils.RawSimToPb(a, req.GetSimType().String())
	err := p.simRepo.Add(s)
	if err != nil {
		logrus.Error("error while Upload sims data" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}
<<<<<<< HEAD
	route := p.baseRoutingKey.SetAction("upload").SetObject("sim").MustBuild()
	err = p.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	return &pb.UploadResponse{Sim: dbSimsToPbSim(s)}, nil
}

func (p *SimPoolServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
=======
	return &pb.UploadResponse{Sim: dbSimsToPbSim(s)}, nil
}

func (p *SimServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
>>>>>>> subscriber-sys_sim-manager
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
<<<<<<< HEAD
		IsAllocated:    p.IsAllocated,
=======
		IsAllocated:    p.Is_allocated,
>>>>>>> subscriber-sys_sim-manager
		SmDpAddress:    p.SmDpAddress,
		ActivationCode: p.ActivationCode,
		CreatedAt:      p.CreatedAt.String(),
		UpdatedAt:      p.UpdatedAt.String(),
		DeletedAt:      p.DeletedAt.Time.String(),
<<<<<<< HEAD
		SimType:        pb.SimType(pb.SimType_value[p.SimType]),
=======
		SimType:        pb.SimType(pb.SimType_value[p.Sim_type]),
>>>>>>> subscriber-sys_sim-manager
	}
}
