package server

import (
	"context"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/utils"
	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const uuidParsingError = "Error parsing UUID"

type BaseRateServer struct {
	baseRateRepo   db.BaseRateRepo
	msgBus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedBaseRatesServiceServer
}

func NewBaseRateServer(baseRateRepo db.BaseRateRepo, msgBus mb.MsgBusServiceClient) *BaseRateServer {
	return &BaseRateServer{
		baseRateRepo:   baseRateRepo,
		msgBus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
	}
}

func (b *BaseRateServer) GetBaseRateById(ctx context.Context, req *pb.GetBaseRatesByIdRequest) (*pb.GetBaseRatesByIdResponse, error) {
	uuid, err := uuid.FromString(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}
	rate, err := b.baseRateRepo.GetBaseRateById(uuid)

	if err != nil {
		logrus.Error("error while getting rate" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}
	resp := &pb.GetBaseRatesByIdResponse{
		Rate: dbRatesToPbRates(rate),
	}

	return resp, nil
}

func (b *BaseRateServer) GetBaseRatesByNetwork(ctx context.Context, req *pb.GetBaseRatesByNetworkRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("GetBaseRates where country = %s and network = %s and simType = %s", req.GetCountry(), req.GetNetwork(), req.GetSimType())
	ef, err := time.Parse(time.RFC3339, req.GetEffectiveAt())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time format format for effective at "+err.Error())
	}

	rates, err := b.baseRateRepo.GetBaseRatesByNetwork(req.GetCountry(), req.GetNetwork(), ef, db.ParseType(req.GetSimType()))

	if err != nil {
		logrus.Errorf("error while getting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rates")
	}
	rateList := &pb.GetBaseRatesResponse{
		Rates: dbratesToPbRates(rates),
	}

	return rateList, nil
}

func (b *BaseRateServer) GetBaseRatesHistoryByNetwork(ctx context.Context, req *pb.GetBaseRatesByNetworkRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("GetBaseRates where country = %s and network = %s and simType = %s", req.GetCountry(), req.GetNetwork(), req.GetSimType())

	rates, err := b.baseRateRepo.GetBaseRatesHistoryByNetwork(req.GetCountry(), req.GetNetwork())

	if err != nil {
		logrus.Errorf("error while getting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rates")
	}
	rateList := &pb.GetBaseRatesResponse{
		Rates: dbratesToPbRates(rates),
	}

	return rateList, nil
}

func (b *BaseRateServer) GetBaseRatesForPeriod(ctx context.Context, req *pb.GetBaseRatesByPeriodRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("GetBaseRates where country = %s and network = %s and simType = %s and Period From %s To %s ", req.GetCountry(), req.GetNetwork(), req.GetSimType(), req.From, req.To)

	from, err := time.Parse(time.RFC3339, req.GetFrom())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time format for from "+err.Error())
	}

	to, err := time.Parse(time.RFC3339, req.GetTo())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time format for to "+err.Error())
	}

	rates, err := b.baseRateRepo.GetBaseRatesForPeriod(req.GetCountry(), req.GetNetwork(), from, to, db.ParseType(req.GetSimType()))

	if err != nil {
		logrus.Errorf("error while getting rates" + err.Error())
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
	strType := strings.ToLower(req.GetSimType())
	simType := db.ParseType(strType)

	if !validations.IsValidUploadReqArgs(fileUrl, effectiveAt, simType.String()) {
		logrus.Infof("Please supply valid fileURL: %s, effectiveAt: %s and simType: %s.",
			fileUrl, effectiveAt, simType)
		return nil, status.Errorf(codes.InvalidArgument, "Please supply valid fileURL: %q, effectiveAt: %q & simType: %q",
			fileUrl, effectiveAt, simType)
	}
	formattedEffectiveAt, err := validations.ValidateDate(effectiveAt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if err := validations.IsFutureDate(formattedEffectiveAt); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())

	}

	sType := db.ParseType(strType)

	if sType.String() != req.SimType {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid sim type: provided sim type (%s) does not match with package allowed sim type (%s)",
			sType.String(), req.SimType)
	}
	data, err := utils.FetchData(fileUrl)
	if err != nil {
		logrus.Infof("Error fetching data: %v", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rates, err := utils.ParseToModel(data, formattedEffectiveAt, simType.String())
	if err != nil {
		return nil, err
	}
	err = b.baseRateRepo.UploadBaseRates(rates)

	if err != nil {
		logrus.Error("error inserting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rate")
	}

	route := b.baseRoutingKey.SetActionUpdate().SetObject("rate").MustBuild()
	err = b.msgBus.PublishRequest(route, req)
	if err != nil {
		logrus.Errorf("Failed to publish message %+v with key %+v. Errors %s", req, route, err.Error())
	}

	rateList := &pb.UploadBaseRatesResponse{
		Rate: dbratesToPbRates(rates),
	}

	return rateList, nil
}

func dbratesToPbRates(rates []db.BaseRate) []*pb.Rate {
	res := []*pb.Rate{}
	for _, u := range rates {
		res = append(res, dbRatesToPbRates(&u))
	}
	return res
}

func dbRatesToPbRates(r *db.BaseRate) *pb.Rate {
	return &pb.Rate{
		Uuid:        r.Uuid.String(),
		X2G:         r.X2g,
		X3G:         r.X3g,
		X5G:         r.X5g,
		Lte:         r.Lte,
		Apn:         r.Apn,
		Vpmn:        r.Vpmn,
		Imsi:        r.Imsi,
		Data:        r.Data,
		LteM:        r.LteM,
		SmsMo:       r.SmsMo,
		SmsMt:       r.SmsMt,
		Network:     r.Network,
		Country:     r.Country,
		SimType:     r.SimType.String(),
		EffectiveAt: r.EffectiveAt.Format(time.RFC3339),
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   r.DeletedAt.Time.Format(time.RFC3339),
	}
}
