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
	"github.com/ukama/ukama/systems/common/validation"
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
	orgName        string
	baseRateRepo   db.BaseRateRepo
	msgBus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedBaseRatesServiceServer
}

func NewBaseRateServer(orgName string, baseRateRepo db.BaseRateRepo, msgBus mb.MsgBusServiceClient) *BaseRateServer {
	return &BaseRateServer{
		orgName:        orgName,
		baseRateRepo:   baseRateRepo,
		msgBus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (b *BaseRateServer) GetBaseRatesById(ctx context.Context, req *pb.GetBaseRatesByIdRequest) (*pb.GetBaseRatesByIdResponse, error) {
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

func (b *BaseRateServer) GetBaseRatesByCountry(ctx context.Context, req *pb.GetBaseRatesByCountryRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("GetBaseRates where country = %s and network = %s and simType = %s", req.GetCountry(), req.GetProvider(), req.GetSimType())

	rates, err := b.baseRateRepo.GetBaseRatesByCountry(req.GetCountry(), req.GetProvider(), db.ParseType(req.GetSimType()))

	if err != nil {
		logrus.Errorf("error while getting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rates")
	}
	rateList := &pb.GetBaseRatesResponse{
		Rates: dbratesToPbRates(rates),
	}

	return rateList, nil
}

func (b *BaseRateServer) GetBaseRatesHistoryByCountry(ctx context.Context, req *pb.GetBaseRatesByCountryRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("GetBaseRates where country = %s and network = %s and simType = %s", req.GetCountry(), req.GetProvider(), req.GetSimType())

	rates, err := b.baseRateRepo.GetBaseRatesHistoryByCountry(req.GetCountry(), req.GetProvider(), db.ParseType(req.GetSimType()))

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
	logrus.Infof("GetBaseRates where country = %s and network = %s and simType = %s and Period From %s To %s ", req.GetCountry(), req.GetProvider(), req.GetSimType(), req.From, req.To)

	from, err := time.Parse(time.RFC3339, req.GetFrom())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time format for from "+err.Error())
	}

	to, err := time.Parse(time.RFC3339, req.GetTo())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time format for to "+err.Error())
	}

	rates, err := b.baseRateRepo.GetBaseRatesForPeriod(req.GetCountry(), req.GetProvider(), from, to, db.ParseType(req.GetSimType()))

	if err != nil {
		logrus.Errorf("error while getting rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rates")
	}
	rateList := &pb.GetBaseRatesResponse{
		Rates: dbratesToPbRates(rates),
	}

	return rateList, nil
}

func (b *BaseRateServer) GetBaseRatesForPackage(ctx context.Context, req *pb.GetBaseRatesByPeriodRequest) (*pb.GetBaseRatesResponse, error) {
	logrus.Infof("GetBaseRatesForPackage where country = %s and network = %s and simType = %s and Period From %s To %s ", req.GetCountry(), req.GetProvider(), req.GetSimType(), req.From, req.To)

	from, err := time.Parse(time.RFC3339, req.GetFrom())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time format for from "+err.Error())
	}

	to, err := time.Parse(time.RFC3339, req.GetTo())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time format for to "+err.Error())
	}

	rates, err := b.baseRateRepo.GetBaseRatesForPackage(req.GetCountry(), req.GetProvider(), from, to, db.ParseType(req.GetSimType()))

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
	endAt := req.GetEndAt()
	strType := strings.ToLower(req.GetSimType())
	simType := db.ParseType(strType)
	logrus.Infof("Upload base rate fileURL: %s, effectiveAt: %s endAt: %s and simType: %s.",
		fileUrl, effectiveAt, endAt, simType)

	if !validations.IsValidUploadReqArgs(fileUrl, effectiveAt, simType.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "Please supply valid fileURL: %q, effectiveAt: %q endAt : %q & simType: %q",
			fileUrl, effectiveAt, endAt, simType)
	}

	formattedEffectiveAt, err := validation.ValidateDate(effectiveAt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if err := validation.IsFutureDate(formattedEffectiveAt); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())

	}

	formattedEndAt, err := validation.ValidateDate(endAt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if err := validation.IsFutureDate(formattedEndAt); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())

	}
	if err := validation.IsAfterDate(formattedEndAt, formattedEffectiveAt); err != nil {
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

	rates, err := utils.ParseToModel(data, formattedEffectiveAt, formattedEndAt, simType.String())
	if err != nil {
		return nil, err
	}
	err = b.baseRateRepo.UploadBaseRates(rates)

	if err != nil {
		logrus.Error("error inserting rates " + err.Error())
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
	var del string

	if r.DeletedAt.Valid {
		del = r.DeletedAt.Time.Format(time.RFC3339)
	}

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
		Provider:    r.Provider,
		Country:     r.Country,
		SimType:     r.SimType.String(),
		EffectiveAt: r.EffectiveAt.Format(time.RFC3339),
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   del,
	}
}
