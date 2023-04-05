package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RateClient struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	host    string
	client  pb.RateServiceClient
}

func NewRateClient(rateHost string, timeout time.Duration) *RateClient {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, rateHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewRateServiceClient(conn)

	return &RateClient{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    rateHost,
	}
}

func (r *RateClient) Close() {
	r.conn.Close()
}

func (r *RateClient) GetRate(req *pb.GetRateRequest) (*pb.GetRateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetRate(ctx, req)
}

func (r *RateClient) UpdateDefaultMarkup(req *pb.UpdateDefaultMarkupRequest) (*pb.UpdateDefaultMarkupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.UpdateDefaultMarkup(ctx, req)
}

func (r *RateClient) GetDefaultMarkup(req *pb.GetDefaultMarkupRequest) (*pb.GetDefaultMarkupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetDefaultMarkup(ctx, req)
}

func (r *RateClient) GetDefaultMarkupHistory(req *pb.GetDefaultMarkupHistoryRequest) (*pb.GetDefaultMarkupHistoryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.GetDefaultMarkupHistory(ctx, req)
}

func (r *RateClient) UpdateMarkup(req *pb.UpdateMarkupRequest) (*pb.UpdateMarkupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.UpdateMarkup(ctx, req)
}

func (r *RateClient) DeleteMarkup(req *pb.DeleteMarkupRequest) (*pb.DeleteMarkupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.client.DeleteMarkup(ctx, req)
}

func (r *RateClient) GetMarkup(req *pb.GetMarkupRequest) (*pb.GetMarkupResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.client.GetMarkup(ctx, req)
}

func (r *RateClient) GetMarkupHistory(req *pb.GetMarkupHistoryRequest) (*pb.GetMarkupHistoryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.client.GetMarkupHistory(ctx, req)
}
