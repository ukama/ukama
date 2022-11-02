package server

import (
	"context"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"
)

  
type BaseRateServer struct {
	baseRateRepo       db.BaseRateRepo
    queuePub       msgbus.QPub
	baseRoutingKey msgbus.RoutingKeyBuilder
	nameGenerator  namegenerator.Generator
	pb.UnimplementedBaseRatesServiceServer
}
func NewBaseRateServer(baseRateRepo db.BaseRateRepo, queuePub msgbus.QPub) *BaseRateServer {
	seed := time.Now().UTC().UnixNano()
	return &BaseRateServer{baseRateRepo: baseRateRepo,
		queuePub:       queuePub,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		nameGenerator:  namegenerator.NewNameGenerator(seed),
	}


}

func (b *BaseRateServer) GetBaseRate(ctx context.Context, req *pb.GetBaseRateRequest) (*pb.GetBaseRateResponse, error) {
	logrus.Infof("Get rate  %v", req.GetRateId())

rateId:=req.GetRateId()
	rate, err := b.baseRateRepo.GetBaseRate(rateId)
	if err != nil {
		logrus.Error("error getting the rate" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}
	resp := &pb.GetBaseRateResponse{
		Rate: dbRateToPbRate(rate),
	}

	return resp, nil
}

func dbRateToPbRate(dbr *db.Rate) *pb.Rate {
	rate := &pb.Rate{
		Country:     dbr.Country,
		Network:     dbr.Network,
		Vpmn:        dbr.Vpmn,
		Imsi:        dbr.Imsi,
		SmsMo:       dbr.Sms_mo,
		SmsMt:       dbr.Sms_mt,
		Data:        dbr.Data,
		X2G:         dbr.X2g,
		X3G:         dbr.X3g,
		Lte:         dbr.Lte,
		LteM:        dbr.Lte_m,
		Apn:         dbr.Apn,
		// CreatedAt:   dbr.Created_at.String(),
		EffectiveAt: dbr.Effective_at.String(),
		EndAt:       dbr.End_at.String(),
		SimType:     dbr.Sim_type,

	}

	return rate
}

func (b *BaseRateServer) GetBaseRates(ctx context.Context, req *pb.GetBaseRatesRequest) (*pb.GetBaseRatesResponse, error) {

country:=req.GetCountry()
network:=req.GetProvider()
simType:=validations.ReqPbToStr(req.GetSimType())

	rates, err := b.baseRateRepo.GetBaseRates(country,network,simType)
	if err != nil {
		logrus.Error("error getting the rate" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}
	rateList := &pb.GetBaseRatesResponse{
		Rates: rates,
	}

	return rateList, nil
}

func (b *BaseRateServer) UploadBaseRates(ctx context.Context, req *pb.UploadBaseRatesRequest) (*pb.UploadBaseRatesResponse, error) {

	fileUrl:=req.GetFileURL()
	effectiveAt:=req.GetEffectiveAt()
	simType:=validations.ReqPbToStr(req.GetSimType())
	
		rates, err := b.baseRateRepo.UploadBaseRates(fileUrl,effectiveAt,simType)
		if err != nil {
			logrus.Error("error getting the rate" + err.Error())
			return nil, grpc.SqlErrorToGrpc(err, "rate")
		}
		rateList := &pb.UploadBaseRatesResponse{
			Rate: rates,
		}
	
		return rateList, nil
	}


