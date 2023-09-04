package server

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/common/validation"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	pb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg/client"
	"github.com/ukama/ukama/systems/data-plan/rate/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const uuidParsingError = "Error parsing UUID"

type RateServer struct {
	orgName        string
	baseRate       client.BaseRateSrvc
	markupRepo     db.MarkupsRepo
	defaultRepo    db.DefaultMarkupRepo
	msgBus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
	pb.UnimplementedRateServiceServer
}

func NewRateServer(orgName string, markupRepo db.MarkupsRepo, defualtMarkupRepo db.DefaultMarkupRepo, baseRate client.BaseRateSrvc, msgBus mb.MsgBusServiceClient) *RateServer {

	return &RateServer{
		orgName:        orgName,
		baseRate:       baseRate,
		markupRepo:     markupRepo,
		defaultRepo:    defualtMarkupRepo,
		msgBus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}

func (r *RateServer) getUserMarkup(uuid uuid.UUID) (float64, error) {
	var rate float64

	markup, err := r.markupRepo.GetMarkupRate(uuid)
	if err != nil {
		log.Errorf("Failed to get markup for user %s: Error %s", uuid.String(), err.Error())

		defMarkup := &db.DefaultMarkup{}
		if sql.IsNotFoundError(err) {
			log.Warn("error while getting specific markup. Error: " + err.Error())
			defMarkup, err = r.defaultRepo.GetDefaultMarkupRate()
		}

		if err != nil {
			log.Error("error while getting markup. Error: " + err.Error())
			return 0, grpc.SqlErrorToGrpc(err, "markup")
		}

		rate = defMarkup.Markup

	} else {
		rate = markup.Markup
	}

	return rate, err
}

func (r *RateServer) GetMarkup(ctx context.Context, req *pb.GetMarkupRequest) (*pb.GetMarkupResponse, error) {

	uuid, err := uuid.FromString(req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	rate, err := r.getUserMarkup(uuid)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetMarkupResponse{
		OwnerId: req.OwnerId,
		Markup:  rate,
	}

	return resp, nil
}

func (r *RateServer) UpdateDefaultMarkup(ctx context.Context, req *pb.UpdateDefaultMarkupRequest) (*pb.UpdateDefaultMarkupResponse, error) {

	err := r.defaultRepo.UpdateDefaultMarkupRate(req.Markup)
	if err != nil {
		log.Error("error while updating default markup" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "default markup")
	}

	if r.msgBus != nil {
		r.PublishDefaultMarkupEvents(req.Markup, msgbus.ACTION_CRUD_UPDATE)
	}

	return &pb.UpdateDefaultMarkupResponse{}, nil
}

func (r *RateServer) GetDefaultMarkup(ctx context.Context, req *pb.GetDefaultMarkupRequest) (*pb.GetDefaultMarkupResponse, error) {

	defMarkup, err := r.defaultRepo.GetDefaultMarkupRate()
	if err != nil {
		log.Error("error while getting default markup" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "default markup")
	}

	resp := &pb.GetDefaultMarkupResponse{
		Markup: defMarkup.Markup,
	}

	return resp, nil
}

func (r *RateServer) GetDefaultMarkupHistory(ctx context.Context, req *pb.GetDefaultMarkupHistoryRequest) (*pb.GetDefaultMarkupHistoryResponse, error) {

	defMarkup, err := r.defaultRepo.GetDefaultMarkupRateHistory()
	if err != nil {
		log.Error("error while getting default markup" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "default markup")
	}

	resp := &pb.GetDefaultMarkupHistoryResponse{
		MarkupRates: defMarkupToPbMarkupRates(defMarkup),
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

	if r.msgBus != nil {
		r.PublishMarkupEvents(req.OwnerId, req.Markup, msgbus.ACTION_CRUD_UPDATE)
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

	if r.msgBus != nil {
		r.PublishMarkupEvents(req.OwnerId, 0, msgbus.ACTION_CRUD_DELETE)
	}

	return &pb.DeleteMarkupResponse{}, nil
}

func (r *RateServer) GetMarkupHistory(ctx context.Context, req *pb.GetMarkupHistoryRequest) (*pb.GetMarkupHistoryResponse, error) {
	uuid, err := uuid.FromString(req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	markup, err := r.markupRepo.GetMarkupRateHistory(uuid)
	if err != nil {
		log.Error("error while getting default markup" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "default markup")
	}

	resp := &pb.GetMarkupHistoryResponse{
		OwnerId:     req.OwnerId,
		MarkupRates: markupToPbMarkupRates(markup),
	}

	return resp, nil
}

func (r *RateServer) GetRate(ctx context.Context, req *pb.GetRateRequest) (*pb.GetRateResponse, error) {

	log.Infof("GetRates where country  %s and network %s and simType %s", req.GetCountry(), req.GetProvider(), req.GetSimType())

	uuid, err := uuid.FromString(req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	markup, err := r.getUserMarkup(uuid)
	if err != nil {
		log.Error("error while getting markup" + err.Error())
		return nil, err
	}

	toT, err := validation.FromString(req.To)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date format for to")
	}

	fromT, err := validation.FromString(req.From)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date format for from")
	}

	to := toT.Format(time.RFC3339)
	from := fromT.Format(time.RFC3339)
	rates, err := r.baseRate.GetBaseRates(&bpb.GetBaseRatesByPeriodRequest{
		Country:  req.Country,
		Provider: req.Provider,
		To:       to,
		From:     from,
		SimType:  req.SimType,
	})
	if err != nil {
		log.Errorf("error while getting base rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "baserates")
	}

	if rates == nil || len(rates.GetRates()) == 0 {
		log.Errorf("no valid base rates found")
		return nil, grpc.SqlErrorToGrpc(fmt.Errorf("no valid base rates found"), "baserates")
	}

	rateList := &pb.GetRateResponse{
		Rates: baseratesToMarkupRates(rates.GetRates(), markup),
	}

	return rateList, nil
}

func (r *RateServer) GetRateById(ctx context.Context, req *pb.GetRateByIdRequest) (*pb.GetRateByIdResponse, error) {

	log.Infof("GetRates for User %sbase rate %s", req.OwnerId, req.BaseRate)

	uuid, err := uuid.FromString(req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, uuidParsingError)
	}

	markup, err := r.getUserMarkup(uuid)
	if err != nil {
		log.Error("error while getting markup" + err.Error())
		return nil, err
	}

	rates, err := r.baseRate.GetBaseRate(&bpb.GetBaseRatesByIdRequest{
		Uuid: req.BaseRate,
	})
	if err != nil {
		log.Errorf("error while getting base rates" + err.Error())
		return nil, grpc.SqlErrorToGrpc(err, "baserates")
	}

	rate := &pb.GetRateByIdResponse{
		Rate: baseRateToMarkupRate(rates.Rate, markup),
	}

	return rate, nil
}

func (r *RateServer) PublishMarkupEvents(ownerId string, markup float64, action string) {
	/* Create event */
	e := &epb.MarkupUpdate{
		OwnerId: ownerId,
		Markup:  markup,
	}

	var route string
	switch action {
	case msgbus.ACTION_CRUD_UPDATE:
		route = r.baseRoutingKey.SetActionUpdate().SetObject("markup").MustBuild()
	case msgbus.ACTION_CRUD_DELETE:
		route = r.baseRoutingKey.SetActionDelete().SetObject("markup").MustBuild()

	}
	err := r.msgBus.PublishRequest(route, e)
	if err != nil {

		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
	}
}

func (r *RateServer) PublishDefaultMarkupEvents(markup float64, action string) {
	/* Create event */
	e := &epb.DefaultMarkupUpdate{
		Markup: markup,
	}

	var route string
	switch action {
	case msgbus.ACTION_CRUD_UPDATE:
		route = r.baseRoutingKey.SetActionUpdate().SetObject("defaultmarkup").MustBuild()
	case msgbus.ACTION_CRUD_DELETE:
		route = r.baseRoutingKey.SetActionDelete().SetObject("defaultmarkup").MustBuild()

	}
	err := r.msgBus.PublishRequest(route, e)
	if err != nil {

		log.Errorf("Failed to publish message %+v with key %+v. Errors %s", e, route, err.Error())
	}
}

func baseratesToMarkupRates(rates []*bpb.Rate, markup float64) []*bpb.Rate {
	res := []*bpb.Rate{}
	for _, rate := range rates {
		res = append(res, baseRateToMarkupRate(rate, markup))
	}
	return res
}

func baseRateToMarkupRate(r *bpb.Rate, markup float64) *bpb.Rate {
	mr := r
	mr.Data = MarkupRate(r.Data, markup)
	mr.SmsMo = MarkupRate(r.SmsMo, markup)
	mr.SmsMt = MarkupRate(r.SmsMt, markup)
	return mr
}

func defMarkupToPbMarkupRates(rates []db.DefaultMarkup) []*pb.MarkupRates {
	res := []*pb.MarkupRates{}
	for _, rate := range rates {
		if rate.DeletedAt.Valid {
			res = append(res, &pb.MarkupRates{
				CreatedAt: rate.CreatedAt.Format(time.RFC3339),
				DeletedAt: rate.DeletedAt.Time.Format(time.RFC3339),
				Markup:    rate.Markup,
			})
		} else {
			res = append(res, &pb.MarkupRates{
				CreatedAt: rate.CreatedAt.Format(time.RFC3339),
				Markup:    rate.Markup,
			})
		}
	}
	return res
}

func markupToPbMarkupRates(rates []db.Markups) []*pb.MarkupRates {
	res := []*pb.MarkupRates{}
	for _, rate := range rates {
		if rate.DeletedAt.Valid {
			res = append(res, &pb.MarkupRates{
				CreatedAt: rate.CreatedAt.Format(time.RFC3339),
				DeletedAt: rate.DeletedAt.Time.Format(time.RFC3339),
				Markup:    rate.Markup,
			})
		} else {
			res = append(res, &pb.MarkupRates{
				CreatedAt: rate.CreatedAt.Format(time.RFC3339),
				Markup:    rate.Markup,
			})
		}

	}
	return res
}

func MarkupRate(cost float64, markup float64) float64 {
	return (cost + (markup*cost)/100)
}
