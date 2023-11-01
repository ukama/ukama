package client

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

func AppendHeadersInContext(ctx context.Context, headers http.Header) context.Context {
	_ctx := metadata.AppendToOutgoingContext(ctx, "X-Session-Token", headers.Get("X-Session-Token"))
	return _ctx
}

func (b *BaseRateClient) GetBaseRatesById(h http.Header, req *pb.GetBaseRatesByIdRequest) (*pb.GetBaseRatesByIdResponse, error) {
	// fmt.Println("Headers", h)
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)

	defer cancel()

	return b.client.GetBaseRatesById(AppendHeadersInContext(ctx, h), req)
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

func (b *BaseRateClient) GetBaseRatesForPackage(req *pb.GetBaseRatesByPeriodRequest) (*pb.GetBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.GetBaseRatesForPackage(ctx, req)
}

func (b *BaseRateClient) UploadBaseRates(req *pb.UploadBaseRatesRequest) (*pb.UploadBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()

	return b.client.UploadBaseRates(ctx, req)
}
