package client

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Registry struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.RegistryServiceClient
	host    string
}

func NewRegistry(host string, timeout time.Duration) *Registry {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	client := pb.NewRegistryServiceClient(conn)

	return &Registry{
		conn:    conn,
		client:  client,
		timeout: timeout,
		host:    host,
	}
}

func NewRegistryFromClient(RegistryClient pb.RegistryServiceClient) *Registry {
	return &Registry{
		host:    "localhost",
		timeout: 1 * time.Second,
		conn:    nil,
		client:  RegistryClient,
	}
}

func (sub *Registry) Close() {
	sub.conn.Close()
}

func (sub *Registry) GetSubscriber(sid string) (*pb.GetSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.Get(ctx, &pb.GetSubscriberRequest{SubscriberId: sid})
}

func (sub *Registry) AddSubscriber(req *pb.AddSubscriberRequest) (*pb.AddSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.Add(ctx, req)
}

func (sub *Registry) DeleteSubscriber(sid string) (*pb.DeleteSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.Delete(ctx, &pb.DeleteSubscriberRequest{SubscriberId: sid})
}

func (sub *Registry) UpdateSubscriber(subscriber *pb.UpdateSubscriberRequest) (*pb.UpdateSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.Update(ctx, &pb.UpdateSubscriberRequest{
		SubscriberId:          subscriber.SubscriberId,
		Email:                 subscriber.Email,
		PhoneNumber:           subscriber.PhoneNumber,
		Address:               subscriber.Address,
		IdSerial:              subscriber.IdSerial,
		ProofOfIdentification: subscriber.ProofOfIdentification,
	})
}

func (sub *Registry) GetByNetwork(networkId string) (*pb.GetByNetworkResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.GetByNetwork(ctx, &pb.GetByNetworkRequest{NetworkId: networkId})
}
