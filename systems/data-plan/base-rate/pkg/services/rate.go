package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/models"
	utils "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/utils"
	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"
)

type BaseRateServer struct {
	BaseRate db.Handler
	pb.UnimplementedBaseRatesServiceServer
}

func (s *BaseRateServer) GetBaseRates(ctx context.Context, req *pb.GetBaseRatesRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("Get all rates %v", req.GetCountry())
	simType := validations.ReqPbToStr(req.GetSimType())

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
		if result := s.BaseRate.Where("sim_type = ? ", simType).Find(&rateList.Rates); result.Error != nil {
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
		return &pb.GetBaseRateResponse{
			Status: http.StatusBadRequest,
			Error:  "Please provide a valid rateId",
		}, nil
	}

	if !validations.IsRequestEmpty(rateId) {
		if result := s.BaseRate.First(&rate, req.RateId); result.Error != nil {
			logrus.Error("error getting the rate :" + result.Error.Error())
			return &pb.GetBaseRateResponse{
				Status: http.StatusNotFound,
				Error:  fmt.Sprintf("No rate find with Id : %s", rateId),
			}, nil
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
		SimType:     rate.Sim_type,
	}

	return &pb.GetBaseRateResponse{
		Rate:   data,
		Status: http.StatusOK,
	}, nil
}

func (s *BaseRateServer) UploadBaseRates(ctx context.Context, req *pb.UploadBaseRatesRequest) (*pb.UploadBaseRatesResponse, error) {
	logrus.Infof("Upload base rates %v", req.GetFileURL())

	if validations.IsRequestEmpty(req.GetFileURL()) ||
		validations.IsRequestEmpty(req.GetEffectiveAt()) ||
		validations.IsRequestEmpty(req.GetSimType().String()) {
		logrus.Infof("Please supply valid fileURL: %s, effectiveAt: %s and simType: %s.",
			req.GetFileURL(), req.GetEffectiveAt(), req.GetSimType().String())
		return &pb.UploadBaseRatesResponse{
			Status: http.StatusBadRequest,
			Error:  "Please supply valid fileURL, effectiveAt and simType.",
		}, nil
	}

	fileUrl := req.GetFileURL()
	effectiveAt := req.GetEffectiveAt()
	destinationFileName := "temp.csv"
	simType := req.GetSimType()

	fde := utils.FetchData(fileUrl, destinationFileName)
	if fde != nil {
		logrus.Infof("Error fetching data: %v", fde.Error())
		return &pb.UploadBaseRatesResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("Error fetching data from URL: %v", fde.Error()),
		}, nil
	}

	f, err := os.Open(destinationFileName)
	if err != nil {
		logrus.Infof("Error opening destination file: %v", err.Error())
		return &pb.UploadBaseRatesResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("Error opening destination file: %v", err.Error()),
		}, nil
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		logrus.Infof("Error opening file: %v", err.Error())
		return &pb.UploadBaseRatesResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("Error reading destination file: %s", err.Error()),
		}, nil
	}

	query := utils.CreateQuery(data, effectiveAt, simType)

	dfe := utils.DeleteFile(destinationFileName)
	if fde != nil {
		logrus.Infof("Error while deleting temp file: %s", dfe.Error())
		return &pb.UploadBaseRatesResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("Error while deleting temp file %s", dfe.Error()),
		}, nil
	}

	result := s.BaseRate.Exec(query)
	if result.Error != nil {
		logrus.Error(result.Error)
		return &pb.UploadBaseRatesResponse{
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error inserting data in DB: %v", result.Error),
		}, nil

	}

	var rateList *pb.UploadBaseRatesResponse = &pb.UploadBaseRatesResponse{}

	if result := s.BaseRate.Find(&rateList.Rate); result.Error != nil {
		logrus.Infof("Error while reading data from DB, %s", result.Error)
		return &pb.UploadBaseRatesResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("Error while reading data from DB, %s", result.Error),
		}, nil
	}

	return rateList, nil
}
