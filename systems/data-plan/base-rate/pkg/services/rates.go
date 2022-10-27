package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/models"
	utils "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/utils"
)

type Server struct {
	H db.Handler
	pb.UnimplementedRatesServiceServer
}

type Result struct {
	Id           int64
	Country      string
	Network      string
	Vpmn         string
	Imsi         string
	Sms_mo       string
	Sms_mt       string
	Data         string
	X2g          string
	X3g          string
	X5g          string
	Lte          string
	Lte_m        string
	Apn          string
	Created_at   string
	Effective_at string
	End_at       string
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

func (s *Server) UploadBaseRates(ctx context.Context, req *pb.UploadBaseRatesRequest) (*pb.UploadBaseRatesResponse, error) {
	sim_type := req.SimType.String()
	fileUrl := req.FileURL
	ratesApplicableFrom := req.EffectiveAt
	destinationFileName := "temp.csv"
	utils.FetchData(fileUrl, destinationFileName)

	f, err := os.Open(destinationFileName)
	utils.Check(err)
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	utils.Check(err)

	query := utils.CreateQuery(data, ratesApplicableFrom, sim_type)

	utils.DeleteFile(destinationFileName)

	s.H.DB.Exec(query)

	var rate_list *pb.UploadBaseRatesResponse = &pb.UploadBaseRatesResponse{}

	if result := s.H.DB.Find(&rate_list.Rate); result.Error != nil {
		fmt.Println(result.Error)
	}

	rates := pb.Rate{}
	rate_list.Rate = append(rate_list.Rate, &rates)

	return rate_list, nil
}
