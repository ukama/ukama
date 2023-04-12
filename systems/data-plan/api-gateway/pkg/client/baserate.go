package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BaseRateClient struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	host    string
	client  pb.BaseRatesServiceClient
}

func NewBaseRateClient(baserateHost string, timeout time.Duration) *BaseRateClient {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, baserateHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewBaseRatesServiceClient(conn)

	return &BaseRateClient{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    baserateHost,
	}
}

func NewBaseRateClientFromClient(client pb.BaseRatesServiceClient) *BaseRateClient {
	return &BaseRateClient{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  client,
	}
}

func (b *BaseRateClient) Close() {
	b.conn.Close()
}

func (b *BaseRateClient) GetBaseRatesById(req *pb.GetBaseRatesByIdRequest) (*pb.GetBaseRatesByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetBaseRatesById(ctx, req)
}

func (b *BaseRateClient) GetBaseRatesByCountry(req *pb.GetBaseRatesByCountryRequest) (*pb.GetBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetBaseRatesByCountry(ctx, req)
}

func (b *BaseRateClient) GetBaseRatesHistoryByCountry(req *pb.GetBaseRatesByCountryRequest) (*pb.GetBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetBaseRatesHistoryByCountry(ctx, req)
}

func (b *BaseRateClient) GetBaseRatesForPeriod(req *pb.GetBaseRatesByPeriodRequest) (*pb.GetBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetBaseRatesForPeriod(ctx, req)
}

func (b *BaseRateClient) UploadBaseRates(req *pb.UploadBaseRatesRequest) (*pb.UploadBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.UploadBaseRates(ctx, req)
}
