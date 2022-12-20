package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/api-gateway/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SubscriberRegistry struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.SubscriberRegistryServiceClient
	host    string
}

func NewSubscriberRegistry(host string, timeout time.Duration) *SubscriberRegistry {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewSubscriberRegistryServiceClient(conn)

	return &SubscriberRegistry{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewSubscriberRegistryFromClient(SubscriberRegistryClient pb.SubscriberRegistryServiceClient) *SubscriberRegistry {
	return &SubscriberRegistry{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  SubscriberRegistryClient,
	}
}

func (sr *SubscriberRegistry) Close() {
	sr.conn.Close()
}

func (sr *SubscriberRegistry) GetSubscriber(req *pb.GetSubscriberRequest) (*pb.GetSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sr.timeout)
	defer cancel()

	return sr.client.Get(ctx, req)
}

func (sr *SubscriberRegistry) AddSubscriber(req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sr.timeout)
	defer cancel()

	return sr.client.Add(ctx, req)
}

func (sr *SubscriberRegistry) DeleteSubscriber(req *pb.DeleteSubscriberRequest) (*pb.DeleteSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sr.timeout)
	defer cancel()

	return sr.client.Delete(ctx, req)
}
