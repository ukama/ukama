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
	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BaseRateServer struct {
	BaseRate db.Handler
	pb.UnimplementedRatesServiceServer
}

func (s *BaseRateServer) GetRates(ctx context.Context, req *pb.RatesRequest) (*pb.RatesResponse, error) {
	logrus.Infof("Get all rates %v", req.GetCountry())
	simType := validations.ReqSimTypeToPb(req.SimType.String())

	var rateList *pb.RatesResponse = &pb.RatesResponse{}

	if !validations.IsRequestEmpty(req.GetCountry(), *req.Provider) {
		getRateLog := fmt.Sprintf("Get rates from %s where provider=%s", req.Country, *req.Provider)
		logrus.Infof(getRateLog)

		if result := s.BaseRate.Where("Country = ? AND Network = ?", req.Country, req.Provider).Find(&rateList.Rates); result.Error != nil {
			logrus.Error(result.Error)
			return nil, result.Error

		}
	} else if !validations.IsRequestEmpty(req.GetCountry()) {
		if result := s.BaseRate.Where("Country = ? ", req.Country).Find(&rateList.Rates); result.Error != nil {
			logrus.Error(result.Error)
			return nil, result.Error
		}
	} else {
		if result := s.BaseRate.Where("simType = ? ", simType).Find(&rateList.Rates); result.Error != nil {
			logrus.Error(result.Error)
			return nil, result.Error
		}
	}

	return rateList, nil

}

func (s *BaseRateServer) GetRate(ctx context.Context, req *pb.RateRequest) (*pb.RateResponse, error) {
	logrus.Infof("Get rate by Id : %s", req.GetRateId())
	rateId := req.GetRateId()
	var rate models.Rate
	if len(req.GetRateId()) == 0 {
		logrus.Infof("Rate Id is not valid: %s", rateId)
		return &pb.RateResponse{}, status.Error(codes.InvalidArgument, "Please supply valid rateId")
	}

	if !validations.IsRequestEmpty(rateId) {
		if result := s.BaseRate.First(&rate, req.RateId); result.Error != nil {
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
		SimType:     validations.ReqSimTypeToPb(rate.SimType),
	}

	return &pb.RateResponse{
		Rate: data,
	}, nil
}

func (s *BaseRateServer) UploadBaseRates(ctx context.Context, req *pb.UploadBaseRatesRequest) (*pb.UploadBaseRatesResponse, error) {
	logrus.Infof("Upload rates %v", req.GetFileURL())

	if validations.IsRequestEmpty(req.GetFileURL()) ||
		validations.IsRequestEmpty(req.GetEffectiveAt()) ||
		validations.IsRequestEmpty(req.GetSimType().String()) {
		logrus.Infof("Invalid arguments")
	}

	fileUrl := req.FileURL
	effectiveAt := req.EffectiveAt
	destinationFileName := "temp.csv"
	simType := validations.ReqSimTypeToPb(req.GetSimType().String())

	utils.FetchData(fileUrl, destinationFileName)

	f, err := os.Open(destinationFileName)
	utils.Check(err)
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	utils.Check(err)

	query := utils.CreateQuery(data, effectiveAt, simType)

	utils.DeleteFile(destinationFileName)

	s.BaseRate.Exec(query)

	var rateList *pb.UploadBaseRatesResponse = &pb.UploadBaseRatesResponse{}

	if result := s.BaseRate.Find(&rateList.Rate); result.Error != nil {
		utils.Check(result.Error)
	}

	rates := pb.Rate{}
	rateList.Rate = append(rateList.Rate, &rates)

	return rateList, nil
}
