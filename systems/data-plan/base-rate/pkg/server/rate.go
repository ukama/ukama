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

	rateId := req.GetRateId()
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
		CreatedAt:   dbr.Created_at.Format(time.RFC3339),
		UpdatedAt:   dbr.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   dbr.Deleted_at.Format(time.RFC3339),
		EffectiveAt: dbr.Effective_at.String(),
		EndAt:       dbr.End_at.String(),
		SimType:     dbr.Sim_type,
	}

	return rate
}

func (b *BaseRateServer) GetBaseRates(ctx context.Context, req *pb.GetBaseRatesRequest) (*pb.GetBaseRatesResponse, error) {

	country := req.GetCountry()
	network := req.GetProvider()
	simType := validations.ReqPbToStr(req.GetSimType())

	rates, err := b.baseRateRepo.GetBaseRates(country, network, simType)
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
		Rate: rates,
	}

	return rateList, nil
}
