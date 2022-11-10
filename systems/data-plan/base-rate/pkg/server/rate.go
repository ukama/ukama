package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/utils"
	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BaseRateServer struct {
	baseRateRepo db.BaseRateRepo
	pb.UnimplementedBaseRatesServiceServer
}

func NewBaseRateServer(baseRateRepo db.BaseRateRepo) *BaseRateServer {
	return &BaseRateServer{baseRateRepo: baseRateRepo}

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
		Rate: dbRatesToPbRates(rate),
	}

	return resp, nil
}

func (b *BaseRateServer) GetBaseRates(ctx context.Context, req *pb.GetBaseRatesRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("GetBaseRates where country =  %s and network =%s and simType =%s", req.GetCountry(), req.GetProvider(), req.GetSimType())
	country := req.GetCountry()
	network := req.GetProvider()
	effectiveAt := req.GetEffectiveAt()
	simType := validations.ReqPbToStr(req.GetSimType())
	rates, err := b.baseRateRepo.GetBaseRates(country, network, simType, effectiveAt)
	if err != nil {
		logrus.Error("error while getting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rates")
	}

	rateList := &pb.GetBaseRatesResponse{
		Rates: dbratesToPbRates(rates),
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
		return nil, status.Errorf(codes.InvalidArgument, "Please supply valid fileURL: %q, effectiveAt: %q & simType: %q",
			fileUrl, effectiveAt, simType)
	}

	if !utils.IsFutureDate(effectiveAt) {
		logrus.Infof("Date you provided is not a valid future date.",
			fileUrl, effectiveAt, simType)

		return nil, status.Errorf(codes.InvalidArgument, "date you provided is not a valid future date %qs", effectiveAt)
	}

	data, err := utils.FetchData(fileUrl)
	if err != nil {
		logrus.Infof("Error fetching data: %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	res := utils.ParseToModel(data, effectiveAt, simType)

	err = b.baseRateRepo.UploadBaseRates(res)

	if err != nil {
		logrus.Error("error inserting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}

	rates, err := b.baseRateRepo.GetBaseRates("", "", effectiveAt, "")
	if err != nil {
		logrus.Error("error fetching rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}

	rateList := &pb.UploadBaseRatesResponse{
		Rate: dbratesToPbRates(rates),
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
		Id:          uint64(r.ID),
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
		EndAt:       r.End_at,
		Network:     r.Network,
		Country:     r.Country,
		SimType:     r.Sim_type,
		EffectiveAt: r.Effective_at,
		CreatedAt:   r.CreatedAt.String(),
		UpdatedAt:   r.UpdatedAt.String(),
		DeletedAt:   r.DeletedAt.Time.String(),
	}
}
