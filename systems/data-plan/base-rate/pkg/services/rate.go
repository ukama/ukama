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
	pb.UnimplementedBaseRatesServiceServer
}

func (s *BaseRateServer) GetBaseRates(ctx context.Context, req *pb.GetBaseRatesRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("Get all rates %v", req.GetCountry())
	simType := validations.ReqStrTopb(req.SimType.String())

	var rateList *pb.GetBaseRatesResponse = &pb.GetBaseRatesResponse{}

	if !validations.IsRequestEmpty(req.GetCountry(), *req.Provider) {
		getRateLog := fmt.Sprintf("Get rates from %s where provider=%s", req.Country, *req.Provider)
		logrus.Infof(getRateLog)

		if result := s.BaseRate.Where("country = ? AND network = ?", req.Country, req.Provider).Find(&rateList.Rates); result.Error != nil {
			logrus.Error(result.Error)
			return nil, result.Error

		}
	} else if !validations.IsRequestEmpty(req.GetCountry()) {
		if result := s.BaseRate.Where("country = ? ", req.Country).Find(&rateList.Rates); result.Error != nil {
			logrus.Error(result.Error)
			return nil, result.Error
		}
	} else {
		if result := s.BaseRate.Where("sim_type = ? ", validations.ReqPbToStr(simType)).Find(&rateList.Rates); result.Error != nil {
			logrus.Error(result.Error)
			return nil, result.Error
		}
	}

	return rateList, nil

}

func (s *BaseRateServer) GetBaseRate(ctx context.Context, req *pb.GetBaseRateRequest) (*pb.GetBaseRateResponse, error) {
	logrus.Infof("Get rate by Id : %s", req.GetRateId())
	rateId := req.GetRateId()
	var rate models.Rate
	if len(req.GetRateId()) == 0 {
		logrus.Infof("Rate Id is not valid: %s", rateId)
		return &pb.GetBaseRateResponse{}, status.Error(codes.InvalidArgument, "Please supply valid rateId")
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
		SimType:     validations.ReqStrTopb(rate.Sim_type),
	}

	return &pb.GetBaseRateResponse{
		Rate: data,
	}, nil
}

func (s *BaseRateServer) UploadBaseRates(ctx context.Context, req *pb.UploadBaseRatesRequest) (*pb.UploadBaseRatesResponse, error) {
	logrus.Infof("Upload base rates %v", req.GetFileURL())

	if validations.IsRequestEmpty(req.GetFileURL()) ||
		validations.IsRequestEmpty(req.GetEffectiveAt()) ||
		validations.IsRequestEmpty(req.GetSimType().String()) {
		err := status.Errorf(codes.InvalidArgument, "Please supply valid fileURL, effectiveAt and simType.")
		return nil, err
	}

	fileUrl := req.GetFileURL()
	effectiveAt := req.GetEffectiveAt()
	destinationFileName := "temp.csv"
	simType := req.GetSimType()

	utils.FetchData(fileUrl, destinationFileName)

	f, err := os.Open(destinationFileName)
	if err != nil {
		e := status.Errorf(codes.Internal, "Error opening file: %v", err)
		return nil, e
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		e := status.Errorf(codes.Internal, "Error opening file: %v", err)
		return nil, e
	}

	query := utils.CreateQuery(data, effectiveAt, simType)

	utils.DeleteFile(destinationFileName)

	result := s.BaseRate.Exec(query)
	if result.Error != nil {
		logrus.Error(result.Error)
		e := status.Errorf(codes.Internal, "Error inserting data in DB: %v", result.Error)
		return nil, e
	}

	var rateList *pb.UploadBaseRatesResponse = &pb.UploadBaseRatesResponse{}

	if result := s.BaseRate.Find(&rateList.Rate); result.Error != nil {
		utils.Check(result.Error)
	}

	rates := pb.Rate{}
	rateList.Rate = append(rateList.Rate, &rates)

	return rateList, nil
}
