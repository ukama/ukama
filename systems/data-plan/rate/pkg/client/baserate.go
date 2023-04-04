package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BaseRate struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.BaseRatesServiceClient
}

type BaseRateSrvc interface {
	GetBaseRates(req *pb.GetBaseRatesRequest) (*pb.GetBaseRatesResponse, error)
	GetBaseRate(id string) (*pb.GetBaseRateResponse, error)
}

func NewBaseRate(baseRate string, timeout time.Duration) (*BaseRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, baseRate, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Errorf("Failed to connect to base rate service at %s. Error %s", baseRate, err.Error())
		return nil, err
	}
	client := pb.NewBaseRatesServiceClient(conn)

	return &BaseRate{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *BaseRate) Close() {
	c.conn.Close()
}

func (c *BaseRate) GetBaseRates(req *pb.GetBaseRatesRequest) (*pb.GetBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.GetBaseRates(ctx, req)
}

func (c *BaseRate) GetBaseRate(id string) (*pb.GetBaseRateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.GetBaseRate(ctx, &pb.GetBaseRateRequest{Uuid: id})
}
