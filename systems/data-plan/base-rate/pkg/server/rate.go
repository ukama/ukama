package server

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/utils"
	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"
)

type BaseRateServer struct {
	baseRateRepo   db.BaseRateRepo
	pb.UnimplementedBaseRatesServiceServer
}


func NewBaseRateServer(baseRateRepo db.BaseRateRepo) *BaseRateServer {
	return &BaseRateServer{baseRateRepo: baseRateRepo,
	}

}

func (b *BaseRateServer) GetBaseRate(ctx context.Context, req *pb.GetBaseRateRequest) (*pb.GetBaseRateResponse, error) {
	logrus.Infof("Get rate %v", req.GetRateId())

	rateId := req.GetRateId()

	rate, err := b.baseRateRepo.GetBaseRate(rateId)
	if err != nil {
		logrus.Error("error while getting rate" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}
	resp := &pb.GetBaseRateResponse{
		Rate:dbRatesToPbRates(rate),
	}

	return resp, nil
}

func (b *BaseRateServer) GetBaseRates(ctx context.Context, req *pb.GetBaseRatesRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("GetBaseRates where country =  %s and network =%s and simType =%s", req.GetCountry(),req.GetProvider(),req.GetSimType())
	country := req.GetCountry()
	network := req.GetProvider()
	simType := validations.ReqPbToStr(req.GetSimType())
	rates, err := b.baseRateRepo.GetBaseRates(country, network, simType)
	if err != nil {
		logrus.Error("error while getting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rates")
	}

	rateList := &pb.GetBaseRatesResponse{
		Rates:dbratesToPbRates(rates),
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
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}

	query := utils.CreateQuery(data, effectiveAt, simType)

	dfe := utils.DeleteFile(destinationFileName)
	if fde != nil {
		logrus.Infof("Error while deleting temp file: %s", dfe.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}
	err = b.baseRateRepo.UploadBaseRates(query)
	if err != nil {
		logrus.Error("error getting the rate" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}

	// rates, err := b.baseRateRepo.GetAllBaseRates(effectiveAt)
	// if err != nil {
	// 	logrus.Error("error getting the rates" + err.Error())
	// 	return nil, grpc.SqlErrorToGrpc(err, "rate")
	// }

	rateList := &pb.UploadBaseRatesResponse{
		Rate: []*pb.Rate{},
	}

	return rateList, nil
}

func dbratesToPbRates(rates []db.Rate) []*pb.Rate {
	res := []*pb.Rate{}
	for _, u := range rates {
		res = append(res, dbRatesToPbRates(&u))
	}
	return res
}

func dbRatesToPbRates(r *db.Rate) *pb.Rate {
	return &pb.Rate{
		Id:          int64(r.ID),
		X2G:         r.X2g,
		X3G:         r.X3g,
		X5G:         r.X5g,
		Lte:         r.Lte,
		Apn:         r.Apn,
		Vpmn:        r.Vpmn,
		Imsi:        r.Imsi,
		Data:        r.Data,
		LteM:        r.Lte_m,
		SmsMo:       r.Sms_mo,
		SmsMt:       r.Sms_mt,
		Network:     r.Network,
		Country:     r.Country,
		SimType:     r.Sim_type,
		EndAt:       r.End_at.String(),
		CreatedAt:   r.CreatedAt.String(),
		UpdatedAt:   r.UpdatedAt.String(),
		EffectiveAt: r.Effective_at.String(),
		DeletedAt:   r.DeletedAt.Time.String(),
	}
}