package services

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateServer struct {
	RateRepo db.Handler
	pb.UnimplementedBaseRatesServiceServer
}

func (r *RateServer) GetBaseRates(ctx context.Context, req *pb.GetBaseRatesRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("Get all rates %v", req.GetCountry())
	simType := reqSimTypeToPb(req.SimType.String())

	var rate_list *pb.GetBaseRatesResponse = &pb.GetBaseRatesResponse{}
	if !isRequestEmpty(req.GetCountry(), *req.Provider) {
		getRateLog := fmt.Sprintf("Get rates from %s where provider=%s", req.Country, *req.Provider)
		logrus.Infof(getRateLog)

		if result := r.RateRepo.Where("Country = ? AND Network = ?", req.Country, req.Provider).Find(&rate_list.Rates); result.Error != nil {
			logrus.Error(result.Error)
			return nil, result.Error

		}

	} else if !isRequestEmpty(req.GetCountry()) {

		if result := r.RateRepo.Where("Country = ? ", req.Country).Find(&rate_list.Rates); result.Error != nil {
			logrus.Error(result.Error)
			return nil, result.Error
		}
	} else {

		fmt.Println(req.SimType)
		if result := r.RateRepo.Where("sim_type = ? ", simType).Find(&rate_list.Rates); result.Error != nil {
			logrus.Error(result.Error)
			return nil, result.Error
		}

	}

	return rate_list, nil

}

func (r *RateServer) GetBaseRate(ctx context.Context, req *pb.GetBaseRateRequest) (*pb.GetBaseRateResponse, error) {
	logrus.Infof("Get rate by Id : %s", req.GetRateId())
	rateId := req.GetRateId()
	var rate models.Rate
	if len(req.GetRateId()) == 0 {
		logrus.Infof("Rate Id is not valid: %s", rateId)
		return &pb.GetBaseRateResponse{}, status.Error(codes.InvalidArgument, "Please supply valid rateId")
	}

	if !isRequestEmpty(rateId) {
		if result := r.RateRepo.First(&rate, req.RateId); result.Error != nil {
			logrus.Error("error getting the rate :" + result.Error.Error())
			return nil, status.Errorf(codes.NotFound, result.Error.Error())
		}
	} else {
		return nil, fmt.Errorf("invalid arguments as RateId")
	}
	data := &pb.Rate{
		Country:     rate.Country,
		Network:     rate.Network,
		Vpmn:        rate.Vpmn,
		Imsi:        rate.Imsi,
		SmsMo:       rate.Sms_mo,
		SmsMt:       rate.Sms_mt,
		Data:        rate.Data,
		X2G:         rate.X2g,
		X3G:         rate.X3g,
		X5G:         rate.X5g,
		Lte:         rate.Lte,
		LteM:        rate.Lte_m,
		Apn:         rate.Apn,
		CreatedAt:   rate.Created_at,
		EffectiveAt: rate.Effective_at,
		EndAt:       rate.End_at,
		SimType:     pb.SimType(rate.SimType),
	}

	return &pb.GetBaseRateResponse{
		Rate: data,
	}, nil

}
func isRequestEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func reqSimTypeToPb(simType string) models.SimType {
	var pbSimType models.SimType
	switch simType {
	case "INTER_MNO_ALL":
		pbSimType = 2
	case "INTER_MNO_DATA":
		pbSimType = 1
	case "INTER_NONE":
		pbSimType = 0
	case "INTER_UKAMA_ALL":
		pbSimType = 3
	default:
		pbSimType = 0
	}
	return pbSimType
}
