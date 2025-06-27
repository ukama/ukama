/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	pb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/sim-pool/pkg/utils"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type SimPoolServer struct {
	orgName        string
	simRepo        db.SimRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedSimServiceServer
}

func NewSimPoolServer(orgName string, simRepo db.SimRepo, msgBus mb.MsgBusServiceClient) *SimPoolServer {
	return &SimPoolServer{
		orgName:        orgName,
		simRepo:        simRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (p *SimPoolServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("GetSim isPhysical: %v, simType: %v", req.GetIsPhysicalSim(), req.GetSimType())

	sim, err := p.simRepo.Get(req.GetIsPhysicalSim(), ukama.ParseSimType(req.GetSimType()))
	if err != nil {
		log.Error("error fetching a sim " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	return &pb.GetResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (p *SimPoolServer) GetByIccid(ctx context.Context, req *pb.GetByIccidRequest) (*pb.GetByIccidResponse, error) {
	log.Infof("GetSimByIccid : %v", req.GetIccid())

	sim, err := p.simRepo.GetByIccid(req.GetIccid())
	if err != nil {
		log.Error("error fetching a sim " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	return &pb.GetByIccidResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (p *SimPoolServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	log.Infof("GetSimStats : %v ", req.GetSimType())

	sim, err := p.simRepo.GetSimsByType(ukama.ParseSimType(req.SimType))
	if err != nil {
		log.Error("error getting a sim pool stats" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}
	resp := utils.PoolStats(sim)

	return resp, nil
}

func (p *SimPoolServer) GetSims(ctx context.Context, req *pb.GetSimsRequest) (*pb.GetSimsResponse, error) {
	log.Infof("GetSims : %v ", req.GetSimType())

	sims, err := p.simRepo.GetSims(ukama.ParseSimType(req.SimType))
	if err != nil {
		log.Error("error getting a sim pool stats" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}
	resp := dbSimsToPbSim(sims)

	return &pb.GetSimsResponse{
		Sims: resp,
	}, nil
}

func (p *SimPoolServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	log.Infof("Add Sims : %v ", req.Sim)
	result := utils.PbParseToModel(req.Sim)
	err := p.simRepo.Add(result)
	if err != nil {
		log.Error("error adding a sims" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	iccids := make([]string, len(result))
	for _, u := range result {
		iccids = append(iccids, u.Iccid)
	}

	route := p.baseRoutingKey.SetAction("upload").SetObject("sim").MustBuild()
	_ = p.PublishEventMessage(route, &epb.SimUploaded{
		Iccid: iccids,
	})

	resp := &pb.AddResponse{Sim: dbSimsToPbSim(result)}
	return resp, nil
}

func (p *SimPoolServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	log.Infof("Upload Sims to pool")
	a, _ := utils.ParseBytesToRawSim(req.SimData)
	s := utils.RawSimToPb(a, req.GetSimType())
	err := p.simRepo.Add(s)
	log.Info("ADDING SIMS: ", s, err)
	if err != nil {
		log.Error("error while Upload sims data" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	iccids := make([]string, len(s))
	for _, u := range s {
		iccids = append(iccids, u.Iccid)
	}

	if p.msgbus != nil {
		route := p.baseRoutingKey.SetAction("upload").SetObject("sim").MustBuild()
		_ = p.PublishEventMessage(route, &epb.SimUploaded{
			Iccid: iccids,
		})
	}

	return &pb.UploadResponse{Iccid: iccids}, nil
}

func (p *SimPoolServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Infof("Delete Sims: %v", req.GetId())
	err := p.simRepo.Delete(req.GetId())
	if err != nil {
		log.Error("error while delete sims data" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-pool")
	}

	route := p.baseRoutingKey.SetActionDelete().SetObject("sim").MustBuild()
	_ = p.PublishEventMessage(route, &epb.SimRemoved{
		Id: req.GetId(),
	})

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

func (p *SimPoolServer) PublishEventMessage(route string, msg protoreflect.ProtoMessage) error {

	err := p.msgbus.PublishRequest(route, msg)
	if err != nil {
		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", msg, route, err.Error())
	}
	return err

}
