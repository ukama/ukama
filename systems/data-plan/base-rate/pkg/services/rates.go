package services

import (
	"context"

	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"
)

type Server struct {
	H db.Handler
	pb.UnimplementedRatesServiceServer
}

func (s *Server) GetRates(ctx context.Context, req *pb.RatesRequest) (*pb.RatesResponse, error) {
	//  var rate models.Rate

	return &pb.RatesResponse{Rates: []*pb.Rate{}}, nil

}

func (s *Server) GetRate(ctx context.Context, req *pb.RateRequest) (*pb.RateResponse, error) {

	return &pb.RateResponse{}, nil

}
