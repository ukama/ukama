package services

import (
	"context"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/models"
)

type Server struct {
	H db.Handler
	pb.UnimplementedRatesServiceServer
}

func (s *Server) GetRates(ctx context.Context, req *pb.RatesRequest) (*pb.RatesResponse, error) {
	logrus.Infof("Get all rates %v", req.GetCountry())

	var rate_list *pb.RatesResponse = &pb.RatesResponse{}
	if !isRequestEmpty(req.GetCountry(), *req.Provider) {
		if result := s.H.DB.Where("Country = ? AND Network = ?", req.Country, req.Provider).Find(&rate_list.Rates); result.Error != nil {
			logrus.Error(result.Error)
		}

	} else if !isRequestEmpty(req.GetCountry()) {

		if result := s.H.DB.Where("Country = ? ", req.Country).Find(&rate_list.Rates); result.Error != nil {
			logrus.Error(result.Error)
		}

	} else {
		if result := s.H.DB.Find(&rate_list.Rates); result.Error != nil {
			logrus.Error(result.Error)
		}
	}

	return rate_list, nil

}

func (s *Server) GetRate(ctx context.Context, req *pb.RateRequest) (*pb.RateResponse, error) {
	logrus.Infof("Get rate by Id : %v", req.GetRateId())

	var rate models.Rate
	if !isRequestEmpty(req.GetRateId()) {
		if result := s.H.DB.First(&rate, req.RateId); result.Error != nil {
			logrus.Error(result.Error)
		}
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
		Lte:         rate.Lte,
		LteM:        rate.Lte_m,
		Apn:         rate.Apn,
		CreatedAt:   rate.Created_at,
		EffectiveAt: rate.Effective_at,
		EndAt:       rate.End_at,
	}
	return &pb.RateResponse{
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
