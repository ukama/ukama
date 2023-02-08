package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/data_plan/base_rate/pb/gen"
	"github.com/ukama/ukama/systems/data_plan/base_rate/pkg/db"
	"github.com/ukama/ukama/systems/data_plan/base_rate/pkg/utils"
	validations "github.com/ukama/ukama/systems/data_plan/base_rate/pkg/validations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const uuidParsingError = "Error parsing UUID"

type BaseRateServer struct {
	baseRateRepo   db.BaseRateRepo
	msgbus         *mb.MsgBusClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedBaseRatesServiceServer
}

func NewBaseRateServer(baseRateRepo db.BaseRateRepo) *BaseRateServer {
	return &BaseRateServer{baseRateRepo: baseRateRepo}

}

func (b *BaseRateServer) GetBaseRate(ctx context.Context, req *pb.GetBaseRateRequest) (*pb.GetBaseRateResponse, error) {
	logrus.Infof("Get rate %v", req.GetRateUuid())
	uuid, err := uuid.Parse(req.RateUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}
	rate, err := b.baseRateRepo.GetBaseRate(uuid)

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
	rates, err := b.baseRateRepo.GetBaseRates(req.GetCountry(), req.GetProvider(), req.GetEffectiveAt(), req.GetSimType().String())

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
	simType := req.GetSimType().String()

	if !validations.IsValidUploadReqArgs(fileUrl, effectiveAt, simType) {
		logrus.Infof("Please supply valid fileURL: %s, effectiveAt: %s and simType: %s.",
			fileUrl, effectiveAt, simType)
		return nil, status.Errorf(codes.InvalidArgument, "Please supply valid fileURL: %q, effectiveAt: %q & simType: %q",
			fileUrl, effectiveAt, simType)
	}

	if !validations.IsFutureDate(effectiveAt) {
		logrus.Infof("Date you provided is not a valid future date. %s", effectiveAt)
		return nil, status.Errorf(codes.InvalidArgument, "date you provided is not a valid future date %qs", effectiveAt)
	}

	data, err := utils.FetchData(fileUrl)
	if err != nil {
		logrus.Infof("Error fetching data: %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rates := utils.ParseToModel(data, effectiveAt, simType)

	err = b.baseRateRepo.UploadBaseRates(rates)

	if err != nil {
		logrus.Error("error inserting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}

	// Publish message to msgbus

	route := b.baseRoutingKey.SetActionUpdate().SetObject("base-rate").MustBuild()
	err = b.msgbus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
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
