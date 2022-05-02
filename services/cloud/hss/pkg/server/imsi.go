package server

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/services/cloud/hss/pb/gen"
	"github.com/ukama/ukama/services/cloud/hss/pkg/db"
	"github.com/ukama/ukama/services/common/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ImsiService struct {
	pb.UnimplementedImsiServiceServer
	imsiRepo   db.ImsiRepo
	subscriber HssSubscriber
	gutiRepo   db.GutiRepo
}

func NewImsiService(hssRepo db.ImsiRepo, gutiRepo db.GutiRepo, queue HssSubscriber) *ImsiService {
	return &ImsiService{imsiRepo: hssRepo,
		subscriber: queue,
		gutiRepo:   gutiRepo}
}

func (s *ImsiService) Get(c context.Context, r *pb.GetImsiRequest) (*pb.GetImsiResponse, error) {
	sub, err := s.imsiRepo.GetByImsi(r.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}
	resp := &pb.GetImsiResponse{Imsi: &pb.ImsiRecord{
		Imsi: sub.Imsi,
		Key:  sub.Key,
		Amf:  sub.Amf,
		Op:   sub.Op,
		Apn: &pb.Apn{
			Name: sub.DefaultApnName,
		},
		AlgoType:    sub.AlgoType,
		CsgId:       sub.CsgId,
		CsgIdPrsent: sub.CsgIdPrsent,
		Sqn:         sub.Sqn,
		UeDlAmbrBps: sub.UeDlAmbrBps,
		UeUlAmbrBps: sub.UeDlAmbrBps,
	}}

	return resp, nil
}

func (s *ImsiService) Add(c context.Context, a *pb.AddImsiRequest) (*pb.AddImsiResponse, error) {
	sub := a.Imsi

	dbSub, err := grpcImsiToDb(sub, a.Org)
	if err != nil {
		return nil, err
	}
	err = s.imsiRepo.Add(a.Org, dbSub)

	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}
	s.subscriber.ImsiAdded(a.Org, a.Imsi)
	return &pb.AddImsiResponse{}, err
}

func (s *ImsiService) Update(c context.Context, req *pb.UpdateImsiRequest) (*pb.UpdateImsiResponse, error) {
	imsi, err := s.imsiRepo.GetByImsi(req.ImsiToUpdate)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "error getting imsi")
	}

	dbSub, err := grpcImsiToDb(req.Imsi, imsi.Org.Name)
	if err != nil {
		return nil, err
	}

	err = s.imsiRepo.Update(req.ImsiToUpdate, dbSub)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}
	s.subscriber.ImsiUpdated(imsi.Org.Name, req.Imsi)
	return &pb.UpdateImsiResponse{}, nil
}

func (s *ImsiService) Delete(c context.Context, req *pb.DeleteImsiRequest) (resp *pb.DeleteImsiResponse, err error) {
	var delImsi *db.Imsi
	switch req.IdOneof.(type) {
	case *pb.DeleteImsiRequest_Imsi:

		delImsi, err = s.imsiRepo.GetByImsi(req.GetImsi())
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "imsi")
		}

	case *pb.DeleteImsiRequest_UserId:
		uuid, err := uuid.Parse(req.GetUserId())
		if err != nil {
			logrus.Errorf("Error parsing uuid %s. Error: %s", uuid, err)
			return nil, fmt.Errorf("error parsing uuid")
		}

		imsis, err := s.imsiRepo.GetImsiByUserUuid(uuid)
		if err != nil {
			return nil, grpc.SqlErrorToGrpc(err, "imsi")
		}

		if len(imsis) == 1 {
			delImsi = imsis[0]
		} else if len(imsis) == 0 {
			return nil, status.Error(codes.NotFound, "imsi not found")
		} else {
			return nil, status.Error(codes.Internal, "invalid number of imsis found")
		}
	}

	err = s.imsiRepo.Delete(delImsi.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	s.subscriber.ImsiDeleted(delImsi.Org.Name, delImsi.Imsi)
	return &pb.DeleteImsiResponse{}, nil

}

func (s *ImsiService) AddGuti(c context.Context, req *pb.AddGutiRequest) (*pb.AddGutiResponse, error) {
	imsi, err := s.imsiRepo.GetByImsi(req.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	err = s.gutiRepo.Update(&db.Guti{
		Imsi:            req.Imsi,
		Plmn_id:         req.Guti.PlmnId,
		Mmegi:           req.Guti.Mmegi,
		Mmec:            req.Guti.Mmec,
		MTmsi:           req.Guti.Mtmsi,
		DeviceUpdatedAt: time.Unix(int64(req.UpdatedAt), 0),
	})
	if err != nil {
		if err.Error() == db.GutiNotUpdatedErr {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, grpc.SqlErrorToGrpc(err, "guti")
	}

	s.subscriber.GutiAdded(imsi.Org.Name, req.Imsi, req.Guti)
	return &pb.AddGutiResponse{}, nil
}

func (s *ImsiService) UpdateTai(c context.Context, req *pb.UpdateTaiRequest) (*pb.UpdateTaiResponse, error) {
	imsi, err := s.imsiRepo.GetByImsi(req.Imsi)
	if err != nil {
		return nil, grpc.SqlErrorToGrpc(err, "imsi")
	}

	err = s.imsiRepo.UpdateTai(req.Imsi, db.Tai{
		PlmId:           req.PlmnId,
		Tac:             req.Tac,
		DeviceUpdatedAt: time.Unix(int64(req.UpdatedAt), 0),
	})

	if err != nil {
		if err.Error() == db.TaiNotUpdatedErr {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, grpc.SqlErrorToGrpc(err, "tai")
	}
	s.subscriber.TaiUpdated(imsi.Org.Name, req)
	return &pb.UpdateTaiResponse{}, nil
}

func grpcImsiToDb(sub *pb.ImsiRecord, orgName string) (*db.Imsi, error) {
	userId, err := uuid.Parse(sub.UserId)
	if err != nil {
		logrus.Errorf("Error parsing uuid %s. Error: %s", sub.UserId, err)
		return nil, fmt.Errorf("error parsing uuid")
	}

	dbSub := &db.Imsi{
		Imsi:           sub.Imsi,
		UserUuid:       userId,
		DefaultApnName: sub.Apn.Name,
		Key:            sub.Key,
		Amf:            sub.Amf,
		Op:             sub.Op,
		Org: &db.Org{
			Name: orgName,
		},
	}

	return dbSub, nil
}
