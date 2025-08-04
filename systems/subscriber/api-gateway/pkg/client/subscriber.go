/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
)

type Registry struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.RegistryServiceClient
	host    string
}

func NewRegistry(host string, timeout time.Duration) *Registry {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Subscriber Registry Service: %v", err)
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
	if sub.conn != nil {
		if err := sub.conn.Close(); err != nil {
			log.Warnf("Failed to gracefully close Subscriber Registry connection: %v", err)
		}
	}
}

func (sub *Registry) GetSubscriber(sid string) (*pb.GetSubscriberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.Get(ctx, &pb.GetSubscriberRequest{SubscriberId: sid})
}

func (sub *Registry) GetSubscriberByEmail(sEmail string) (*pb.GetSubscriberByEmailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sub.timeout)
	defer cancel()
	return sub.client.GetByEmail(ctx, &pb.GetSubscriberByEmailRequest{Email: sEmail})
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
		Name:                  subscriber.Name,
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
