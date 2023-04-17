package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/data-plan/rate/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Rate struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.RateServiceClient
}

type RateService interface {
	GetRates(req *pb.GetRatesRequest) (*pb.GetRatesResponse, error)
	GetRate(id string) (*pb.GetRateResponse, error)
}

func NewRate(rate string, timeout time.Duration) (*Rate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, rate, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Errorf("Failed to connect to rate service at %s. Error %s", rate, err.Error())
		return nil, err
	}
	client := pb.NewRateServiceClient(conn)

	return &Rate{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *Rate) Close() {
	c.conn.Close()
}

func (c *Rate) GetRates(req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.GetRates(ctx, req)
}

func (c *Rate) GetRate(id string) (*pb.GetRateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.client.GetRate(ctx, &pb.GetRateRequest{Uuid: id})
}
