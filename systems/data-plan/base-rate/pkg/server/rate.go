package server

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/utils"
	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"
)

type BaseRateServer struct {
	baseRateRepo   db.BaseRateRepo
	baseRoutingKey msgbus.RoutingKeyBuilder
	nameGenerator  namegenerator.Generator
	pb.UnimplementedBaseRatesServiceServer
}
type GetRatesParams struct {
	country, network, simType string
}

func NewBaseRateServer(baseRateRepo db.BaseRateRepo) *BaseRateServer {
	seed := time.Now().UTC().UnixNano()
	return &BaseRateServer{baseRateRepo: baseRateRepo,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
		nameGenerator:  namegenerator.NewNameGenerator(seed),
	}

}

func (b *BaseRateServer) GetBaseRate(ctx context.Context, req *pb.GetBaseRateRequest) (*pb.GetBaseRateResponse, error) {
	logrus.Infof("Get rate  %v", req.GetRateId())

	rateId := req.GetRateId()

	rate, err := b.baseRateRepo.GetBaseRate(rateId)
	if err != nil {
		logrus.Error("error getting the rate" + err.Error())
		return &pb.GetBaseRateResponse{
			Status: http.StatusBadRequest,
			Error:  "Rate ID is required.",
		}, nil

	}
	resp := &pb.GetBaseRateResponse{
		Rate:   rate.ToPbRate(),
		Status: http.StatusAccepted,
	}

	return resp, nil
}

func (b *BaseRateServer) GetBaseRates(ctx context.Context, req *pb.GetBaseRatesRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("GetBaseRates by country or network : %s", req.GetCountry())
	country := req.GetCountry()
	network := req.GetProvider()
	simType := validations.ReqPbToStr(req.GetSimType())
	if validations.IsRequestEmpty(country) {
		return &pb.GetBaseRatesResponse{
			Status: http.StatusBadRequest,
			Error:  "Country name is required!",
		}, nil
	}
	rates, err := b.baseRateRepo.GetBaseRates(country, network, simType)
	if err != nil {
		logrus.Error("error getting the rate" + err.Error())
		return &pb.GetBaseRatesResponse{
			Status: http.StatusBadRequest,
			Error:  "Please provide required params!.",
		}, nil
	}

	rateList := &pb.GetBaseRatesResponse{
		Rates: rates.ToPbRates(),
	}

	return rateList, nil
}

func (b *BaseRateServer) UploadBaseRates(ctx context.Context, req *pb.UploadBaseRatesRequest) (*pb.UploadBaseRatesResponse, error) {

	fileUrl := req.GetFileURL()
	effectiveAt := req.GetEffectiveAt()
	simType := validations.ReqPbToStr(req.GetSimType())

	if validations.IsRequestEmpty(fileUrl) ||
		validations.IsRequestEmpty(effectiveAt) ||
		validations.IsRequestEmpty(simType) {
		logrus.Infof("Please supply valid fileURL: %s, effectiveAt: %s and simType: %s.",
			fileUrl, effectiveAt, simType)
		return &pb.UploadBaseRatesResponse{
			Status: http.StatusBadRequest,
			Error:  "Please supply valid fileURL, effectiveAt and simType.",
		}, nil
	}

	if !utils.IsFutureDate(effectiveAt) {
		logrus.Infof("Date you provided is not a future date.",
			fileUrl, effectiveAt, simType)
		return &pb.UploadBaseRatesResponse{
			Status: http.StatusBadRequest,
			Error:  "Date you provided is not a future date.",
		}, nil
	}

	destinationFileName := "temp.csv"
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
	err = b.baseRateRepo.UploadBaseRates(query)
	if err != nil {
		logrus.Error("error getting the rate" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}

	rates, err := b.baseRateRepo.GetAllBaseRates(effectiveAt)
	if err != nil {
		logrus.Error("error getting the rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}

	rateList := &pb.UploadBaseRatesResponse{
		Rate: rates.ToPbRates(),
	}

	return rateList, nil
}
