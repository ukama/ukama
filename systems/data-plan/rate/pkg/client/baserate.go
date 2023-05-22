package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	bpb "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BaseRate struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  bpb.BaseRatesServiceClient
}

type BaseRateSrvc interface {
	GetBaseRates(req *bpb.GetBaseRatesByPeriodRequest) (*bpb.GetBaseRatesResponse, error)
	GetBaseRate(req *bpb.GetBaseRatesByIdRequest) (*bpb.GetBaseRatesByIdResponse, error)
}

func NewBaseRate(baseRate string, timeout time.Duration) (*BaseRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, baseRate, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Errorf("Failed to connect to base rate service at %s. Error %s", baseRate, err.Error())
		return nil, err
	}
	client := bpb.NewBaseRatesServiceClient(conn)

	return &BaseRate{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *BaseRate) Close() {
	c.conn.Close()
}

func (c *BaseRate) GetBaseRates(req *bpb.GetBaseRatesByPeriodRequest) (*bpb.GetBaseRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.GetBaseRatesForPeriod(ctx, req)
}

func (c *BaseRate) GetBaseRate(req *bpb.GetBaseRatesByIdRequest) (*bpb.GetBaseRatesByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.GetBaseRatesById(ctx, req)
}
