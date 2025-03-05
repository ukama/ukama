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
	pb "github.com/ukama/ukama/testing/services/dummy/dsimfactory/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dsimfactory/pkg"
	"github.com/ukama/ukama/testing/services/dummy/dsimfactory/pkg/db"
	"github.com/ukama/ukama/testing/services/dummy/dsimfactory/pkg/utils"
)

type DsimfactoryServer struct {
	orgName        string
	simRepo        db.SimRepo
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedDsimfactoryServiceServer
}

func NewDsimfactoryServer(orgName string, simRepo db.SimRepo, msgBus mb.MsgBusServiceClient) *DsimfactoryServer {
	return &DsimfactoryServer{
		orgName:        orgName,
		simRepo:        simRepo,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (p *DsimfactoryServer) GetByIccid(ctx context.Context, req *pb.GetByIccidRequest) (*pb.GetByIccidResponse, error) {
	log.Infof("GetSimByIccid : %v", req.GetIccid())

	sim, err := p.simRepo.GetByIccid(req.GetIccid())
	if err != nil {
		log.Error("error fetching a sim " + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-factory")
	}

	return &pb.GetByIccidResponse{Sim: dbSimToPbSim(sim)}, nil
}

func (p *DsimfactoryServer) GetSims(ctx context.Context, req *pb.GetSimsRequest) (*pb.GetSimsResponse, error) {
	log.Info("GetSims")

	sims, err := p.simRepo.GetSims()
	if err != nil {
		log.Error("error getting a sim from factory" + err.Error())

		return nil, grpc.SqlErrorToGrpc(err, "sim-factory")
	}
	resp := dbSimsToPbSim(sims)

	return &pb.GetSimsResponse{
		Sims: resp,
	}, nil
}

func (p *DsimfactoryServer) Upload(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	log.Infof("Upload Sims to factory")
	a, _ := utils.ParseBytesToRawSim(req.SimData)
	s := utils.RawSimToPb(a)
	err := p.simRepo.Add(s)
	log.Info("ADDING SIMS: ", s, err)
	if err != nil {
		log.Error("error while Upload sims data" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "sim-factory")
	}

	iccids := make([]string, len(s))
	for _, u := range s {
		iccids = append(iccids, u.Iccid)
	}

	return &pb.UploadResponse{Iccid: iccids}, nil
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
		SmDpAddress:    p.SmDpAddress,
		ActivationCode: p.ActivationCode,
		IsPhysical:     p.IsPhysical,
		QrCode:         p.QrCode,
		Imsi:           p.Imsi,
	}
}
