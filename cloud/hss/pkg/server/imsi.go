package server

import (
	"context"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	"github.com/ukama/ukamaX/cloud/hss/pkg"
	"github.com/ukama/ukamaX/cloud/hss/pkg/db"
	"github.com/ukama/ukamaX/common/grpc"
)

type ImsiService struct {
	pb.UnimplementedImsiServiceServer
	imsiRepo db.ImsiRepo
	queue    pkg.HssQueue
}

func NewImsiService(hssRepo db.ImsiRepo, queue pkg.HssQueue) *ImsiService {
	return &ImsiService{imsiRepo: hssRepo,
		queue: queue}
}

func (s *ImsiService) Get(c context.Context, r *pb.GetImsiRequest) (*pb.GetImsiResponse, error) {
	sub, err := s.imsiRepo.GetByImsi(r.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}
	resp := &pb.GetImsiResponse{Imsi: &pb.ImsiRecord{
		Imsi:           sub.Imsi,
		Key:            sub.Key,
		Amf:            sub.Amf,
		Op:             sub.Op,
		DefaultApnName: sub.DefaultApnName,
		AuthVector: &pb.AuthVector{
			Response:     sub.AuthVector.Response,
			CipherKey:    sub.AuthVector.CipherKey,
			Token:        sub.AuthVector.Token,
			IntegrityKey: sub.AuthVector.IntegrityKey,
		},
	}}

	return resp, nil
}

func (s *ImsiService) Add(c context.Context, a *pb.AddImsiRequest) (*pb.AddImsiResponse, error) {
	sub := a.Imsi

	dbSub := grpcImsiToDb(sub, a.Org)
	err := s.imsiRepo.Add(a.Org, dbSub)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}
	s.queue.SendImsiAddedEvent(a.Imsi.Imsi)
	return &pb.AddImsiResponse{}, err
}

func (s *ImsiService) Update(c context.Context, req *pb.UpdateImsiRequest) (*pb.UpdateImsiResponse, error) {
	dbSub := grpcImsiToDb(req.Imsi, "")
	err := s.imsiRepo.Update(req.ImsiToUpdate, dbSub)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}
	return &pb.UpdateImsiResponse{}, nil
}

func (s *ImsiService) Delete(c context.Context, req *pb.DeleteImsiRequest) (*pb.DeleteImsiResponse, error) {
	err := s.imsiRepo.Delete(req.Imsi)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}
	return &pb.DeleteImsiResponse{}, nil
}

func grpcImsiToDb(sub *pb.ImsiRecord, orgName string) *db.Imsi {

	dbSub := &db.Imsi{
		Imsi:           sub.Imsi,
		DefaultApnName: sub.DefaultApnName,
		Key:            sub.Key,
		Amf:            sub.Amf,
		Op:             sub.Op,
		Org: &db.Org{
			Name: orgName,
		},
	}

	if sub.AuthVector != nil {
		dbSub.AuthVector = &db.AuthVector{
			CipherKey:    sub.AuthVector.CipherKey,
			Token:        sub.AuthVector.Token,
			IntegrityKey: sub.AuthVector.IntegrityKey,
			Response:     sub.AuthVector.Response,
		}
	}

	return dbSub
}
