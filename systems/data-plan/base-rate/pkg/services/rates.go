package services

import (
	"context"
	"fmt"

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
	logrus.Infof("Get all rates  %v", req.GetCountry())

	var rate_list *pb.RatesResponse = &pb.RatesResponse{}

	if result := s.H.DB.Find(&rate_list.Rates); result.Error != nil {
		fmt.Println(result.Error)
	}

	rates := pb.Rate{}
	rate_list.Rates = append(rate_list.Rates, &rates)

	return rate_list, nil

}

func (s *Server) GetRate(ctx context.Context, req *pb.RateRequest) (*pb.RateResponse, error) {
	logrus.Infof("Get rate  %v", req.GetRateId())

	var rate models.Rate

	if result := s.H.DB.First(&rate, req.RateId); result.Error != nil {
		fmt.Println(result.Error)
	}

	data := &pb.Rate{
		Country:           rate.Country,
		CountryOnCronus:   rate.Country_on_cronus,
		Network:           rate.Network,
		NetworkIdOnCronus: rate.Network_id_on_cronus,
		Vpmn:              rate.Vpmn,
		Imsi:              rate.Imsi,
		SmsMo:             rate.Sms_mo,
		SmsMt:             rate.Sms_mt,
		Data:              rate.Data,
		X2G:               rate.X2g,
		X3G:               rate.X3g,
		Lte:               rate.Lte,
		LteM:              rate.Lte_m,
		Apn:               rate.Apn,
		CreatedAt:         rate.Created_at,
		EffectiveAt:       rate.Effective_at,
		EndAt:             rate.End_at,
	}
	return &pb.RateResponse{
		Rate: data,
	}, nil

}

func (s *Server) UploadRates(ctx context.Context, req *pb.UploadRatesRequest) (*pb.UploadRatesResponse, error) {

	var rate_list *pb.RatesResponse = &pb.RatesResponse{}

	if result := s.H.DB.Find(&rate_list.Rates); result.Error != nil {
		fmt.Println(result.Error)
	}

	rates := pb.Rate{}

	fmt.Println(rates)

	return &pb.UploadRatesResponse{}, nil
}
