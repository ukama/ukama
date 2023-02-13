package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SubscriberRegistry struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.RegistryServiceClient
	host    string
}

func NewSubscriberRegistry(host string, timeout time.Duration) *SubscriberRegistry {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewRegistryServiceClient(conn)

	return &SubscriberRegistry{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewSubscriberRegistryFromClient(SubscriberRegistryClient pb.RegistryServiceClient) *SubscriberRegistry {
	return &SubscriberRegistry{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  SubscriberRegistryClient,
	}
}

func (sub *SubscriberRegistry) Close() {
	sub.conn.Close()
}

func (sub *SubscriberRegistry) GetSubscriber(sid string) (*pb.GetSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.Get(ctx, &pb.GetSubscriberRequest{SubscriberID: sid})
}

func (sub *SubscriberRegistry) AddSubscriber(req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.Add(ctx, req)
}

func (sub *SubscriberRegistry) DeleteSubscriber(sid string) (*pb.DeleteSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.Delete(ctx, &pb.DeleteSubscriberRequest{SubscriberID: sid})
}

func (sub *SubscriberRegistry) UpdateSubscriber(subscriber *pb.UpdateSubscriberRequest) (*pb.UpdateSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.Update(ctx, &pb.UpdateSubscriberRequest{
		SubscriberID:          subscriber.SubscriberID,
		Email:                 subscriber.Email,
		PhoneNumber:           subscriber.PhoneNumber,
		Address:               subscriber.Address,
		IdSerial:              subscriber.IdSerial,
		ProofOfIdentification: subscriber.ProofOfIdentification,
	})
}

func (sub *SubscriberRegistry) GetByNetwork(networkId string) (*pb.GetByNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.GetByNetwork(ctx, &pb.GetByNetworkRequest{NetworkID: networkId})
}
