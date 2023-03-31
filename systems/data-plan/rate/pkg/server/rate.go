package server

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg/client"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const uuidParsingError = "Error parsing UUID"

type RateServer struct {
	baseRate       *client.BaseRate
	markupRepo     db.MarkupsRepo
	msgBus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedRateServiceServer
}

func NewRateServer(markupRepo db.MarkupsRepo, baseRate string, msgBus mb.MsgBusServiceClient, timeout time.Duration) (*RateServer, error) {
	b, err := client.NewBaseRate(baseRate, timeout)
	if err != nil {
		return nil, err
	}

	return &RateServer{
		baseRate:       b,
		markupRepo:     markupRepo,
		msgBus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetContainer(pkg.ServiceName),
	}, nil
}

func (r *RateServer) GetMarkup(ctx context.Context, req *pb.GetMarkupRequest) (*pb.GetMarkupResponse, error) {
	uuid, err := uuid.FromString(req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	markup, err := r.markupRepo.GetMarkupRate(uuid)
	if err != nil {
		log.Error("error while getting markup" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "markup")
	}

	resp := &pb.GetMarkupResponse{
		OwnerId: req.OwnerId,
		Markup:  markup.Markup,
	}

	return resp, nil
}

func (r *RateServer) UpdateMarkup(ctx context.Context, req *pb.UpdateMarkupRequest) (*pb.UpdateMarkupResponse, error) {
	uuid, err := uuid.FromString(req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	err = r.markupRepo.UpdateMarkupRate(uuid, req.Markup)
	if err != nil {
		log.Error("error while updating markup" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "markup")
	}

	return &pb.UpdateMarkupResponse{}, nil
}

func (r *RateServer) DeleteMarkup(ctx context.Context, req *pb.DeleteMarkupRequest) (*pb.DeleteMarkupResponse, error) {
	uuid, err := uuid.FromString(req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	err = r.markupRepo.DeleteMarkupRate(uuid)
	if err != nil {
		log.Error("error while deleting markup" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "markup")
	}

	return &pb.DeleteMarkupResponse{}, nil
}

func (r *RateServer) GetRate(ctx context.Context, req *pb.GetRateRequest) (*pb.GetRateResponse, error) {

	log.Infof("GetRates where country =  %s and network =%s and simType =%s", req.GetCountry(), req.GetProvider(), req.GetSimType())

	uuid, err := uuid.FromString(req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	markup, err := r.markupRepo.GetMarkupRate(uuid)
	if err != nil {
		log.Error("error while getting markup" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "markup")
	}

	rates, err := r.baseRate.GetBaseRates(&pb.GetBaseRatesRequest{
		Country:     req.Country,
		Provider:    req.Provider,
		To:          req.To,
		From:        req.From,
		SimType:     req.SimType,
		EffectiveAt: req.EffectiveAt,
	})
	if err != nil {
		log.Errorf("error while getting base rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "rates")
	}

	rateList := &pb.GetRateResponse{
		Rates: baseratesToMarkupRates(rates.GetRates(), markup.Markup),
	}

	return rateList, nil
}

func baseratesToMarkupRates(rates []*pb.Rate, markup float64) []*pb.Rate {
	res := []*pb.Rate{}
	for _, rate := range rates {
		res = append(res, baseRateToMarkupRate(rate, markup))
	}
	return res
}

func baseRateToMarkupRate(r *pb.Rate, markup float64) *pb.Rate {
	return &pb.Rate{
		Data:  MarkupRate(r.Data, markup),
		SmsMo: MarkupRate(r.SmsMo, markup),
		SmsMt: MarkupRate(r.SmsMt, markup),
	}
}

func MarkupRate(cost float64, markup float64) float64 {
	return (cost + (markup*cost)/100)
}
