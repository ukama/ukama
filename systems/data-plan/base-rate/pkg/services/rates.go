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
		SmsMo:       rate.SmsMo,
		SmsMt:       rate.SmsMt,
		Data:        rate.Data,
		X2G:         rate.X2g,
		X3G:         rate.X3g,
		Lte:         rate.Lte,
		LteM:        rate.LteM,
		Apn:         rate.Apn,
		CreatedAt:   rate.CreatedAt,
		EffectiveAt: rate.EffectiveAt,
		EndAt:       rate.EndAt,
	}
	return &pb.RateResponse{
		Rate: data,
	}, nil
}

func (s *Server) UploadBaseRates(ctx context.Context, req *pb.UploadBaseRatesRequest) (*pb.UploadBaseRatesResponse, error) {
	simType := req.SimType.String()
	fileUrl := req.FileURL
	effectiveAt := req.EffectiveAt
	destinationFileName := "temp.csv"
	utils.FetchData(fileUrl, destinationFileName)

	f, err := os.Open(destinationFileName)
	utils.Check(err)
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	utils.Check(err)

	query := utils.CreateQuery(data, effectiveAt, simType)

	utils.DeleteFile(destinationFileName)

	s.H.DB.Exec(query)

	var rateList *pb.UploadBaseRatesResponse = &pb.UploadBaseRatesResponse{}

	if result := s.H.DB.Find(&rateList.Rate); result.Error != nil {
		fmt.Println(result.Error)
	}

	rates := pb.Rate{}
	rateList.Rate = append(rateList.Rate, &rates)

	return rateList, nil
}
